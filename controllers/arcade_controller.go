/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	arcadev1alpha1 "github.com/redhat-marketplace/marketplace-games-operator/api/v1alpha1"
)

// PORT which container application will be exposed
const PORT = 8080

// ArcadeReconciler reconciles a Arcade object
type ArcadeReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=arcade.marketplace.redhat.com,resources=arcade,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=arcade.marketplace.redhat.com,resources=arcade/status,verbs=get;update;patch

// Reconcile reconciles the Arcade object
func (r *ArcadeReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("arcade", req.NamespacedName)

	arcade := &arcadev1alpha1.Arcade{}
	err := r.Get(ctx, req.NamespacedName, arcade)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Arcade resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Arcade")
		return ctrl.Result{}, err
	}

	nameSpaceKey := types.NamespacedName{Name: arcade.Name, Namespace: arcade.Namespace}
	arcade.Status.ArcadeStatus = arcadev1alpha1.ArcadeStatusPending

	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, nameSpaceKey, deployment)
	if err != nil && errors.IsNotFound(err) {
		dep := r.deploymentForArcade(arcade)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			reason := fmt.Sprintf("Failed to create Deployment %s", err)
			arcade.Status.ArcadeStatus = arcadev1alpha1.ArcadeStatusFailure
			arcade.Status.Reason = reason
			// Update Status
			r.updateStatus(arcade, nameSpaceKey)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	service := &corev1.Service{}
	err = r.Get(ctx, nameSpaceKey, service)
	if err != nil && errors.IsNotFound(err) {
		svc := r.serviceForArcade(arcade)
		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			reason := fmt.Sprintf("Failed to create Service %s", err)
			arcade.Status.ArcadeStatus = arcadev1alpha1.ArcadeStatusFailure
			arcade.Status.Reason = reason
			r.updateStatus(arcade, nameSpaceKey)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	route := &routev1.Route{}
	err = r.Get(ctx, nameSpaceKey, route)
	if err != nil && errors.IsNotFound(err) {
		rt := r.routeForArcade(arcade)
		log.Info("Creating a new Route", "Route.Namespace", rt.Namespace, "Route.Name", rt.Name)
		err = r.Create(ctx, rt)
		if err != nil {
			log.Error(err, "Failed to create new Route", "Route.Namespace", rt.Namespace, "Route.Name", rt.Name)
			reason := fmt.Sprintf("Failed to create Route %s", err)
			arcade.Status.ArcadeStatus = arcadev1alpha1.ArcadeStatusFailure
			arcade.Status.Reason = reason
			r.updateStatus(arcade, nameSpaceKey)
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Second * 15}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Route")
		return ctrl.Result{}, err
	}

	// Ensure Size matches our given Spec
	size := arcade.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		err = r.Update(ctx, deployment)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return ctrl.Result{}, err
		}
		// Return and requeue once Spec is updated
		return ctrl.Result{Requeue: true}, nil
	}

	// Set OK Status
	arcade.Status.ArcadeStatus = arcadev1alpha1.ArcadeStatusOK
	arcade.Status.Reason = ""
	r.updateStatus(arcade, nameSpaceKey)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the Arcade object
func (r *ArcadeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&arcadev1alpha1.Arcade{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&routev1.Route{}).
		Complete(r)
}

func (r *ArcadeReconciler) updateStatus(arcade *arcadev1alpha1.Arcade, nameSpaceKey types.NamespacedName) {
	pollInteraval := time.Millisecond * 500
	timeout := time.Second * 30
	ctx := context.TODO()

	r.Log.Info("Updating status for ", "Arcade instance", arcade.Name)
	err := wait.Poll(pollInteraval, timeout, func() (bool, error) {
		if updateErr := r.Status().Update(ctx, arcade); errors.IsConflict(updateErr) {
			updatedArcade := &arcadev1alpha1.Arcade{ObjectMeta: metav1.ObjectMeta{
				Name:      arcade.Name,
				Namespace: arcade.Namespace,
			}}
			if err := r.Get(context.TODO(), nameSpaceKey, updatedArcade); err != nil {
				return false, err
			}

			// override only the spec
			updatedArcade.Status = arcade.Status
			arcade = updatedArcade
			return false, nil
		} else if updateErr != nil {
			return false, updateErr
		}
		return true, nil
	})
	if err != nil {
		r.Log.Error(err, "Error while updating status")
	}
}

// deploymentForArcade returns a Arcade Deployment object
func (r *ArcadeReconciler) deploymentForArcade(m *arcadev1alpha1.Arcade) *appsv1.Deployment {
	ls := labelsForArcade(m.Name)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: os.Getenv("REGISTRY_REPO") + "." + os.Getenv("REGISTRY_HOST") + "/rhm-arcade",
						Name:  "arcade",
						Ports: []corev1.ContainerPort{{
							ContainerPort: PORT,
							Name:          "arcade",
						}},
					}},
					ImagePullSecrets: []corev1.LocalObjectReference{{
						Name: "regcreds-arcade-art",
					}},
				},
			},
		},
	}
	// Set Arcade instance as the owner and controller for deployments
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// serviceForArcade returns a Arcade Service object
func (r *ArcadeReconciler) serviceForArcade(m *arcadev1alpha1.Arcade) *corev1.Service {
	ls := labelsForArcade(m.Name)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Type: "ClusterIP",
			Ports: []corev1.ServicePort{{
				Port:       PORT,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromInt(PORT),
			}},
			Selector: ls,
		},
	}

	// Set Arcade instance as the owner and controller for created services
	ctrl.SetControllerReference(m, service, r.Scheme)
	return service
}

// routeForArcade returns the route resource
func (r *ArcadeReconciler) routeForArcade(m *arcadev1alpha1.Arcade) *routev1.Route {
	ls := labelsForArcade(m.Name)
	route := &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind: "Route",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: m.Name,
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromInt(PORT),
			},
			TLS: &routev1.TLSConfig{
				Termination: routev1.TLSTerminationEdge,
			},
		},
	}

	// Set Arcade instance as the owner and controller for created routes
	ctrl.SetControllerReference(m, route, r.Scheme)
	return route
}

// labelsForArcade returns the labels for selecting the resources
// belonging to the given arcade CR name.
func labelsForArcade(name string) map[string]string {
	return map[string]string{"app": "arcade", "arcade_cr": name}
}
