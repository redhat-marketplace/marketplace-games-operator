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

	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	arcadev1alpha1 "github.com/redhat-marketplace/marketplace-games-operator/api/v1alpha1"
)

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

	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: arcade.Name, Namespace: arcade.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		dep := r.deploymentForArcade(arcade)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	service := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: arcade.Name, Namespace: arcade.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		svc := r.serviceForArcade(arcade)
		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	route := &routev1.Route{}
	err = r.Get(ctx, types.NamespacedName{Name: arcade.Name, Namespace: arcade.Namespace}, route)
	if err != nil && errors.IsNotFound(err) {
		rt := r.routeForArcade(arcade)
		log.Info("Creating a new Route", "Route.Namespace", rt.Namespace, "Route.Name", rt.Name)
		err = r.Create(ctx, rt)
		if err != nil {
			log.Error(err, "Failed to create new Route", "Route.Namespace", rt.Namespace, "Route.Name", rt.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

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
						Image: "arcade:1.0.0",
						Name:  "arcade",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 4004,
							Name:          "arcade",
						}},
					}},
				},
			},
		},
	}
	// Set Arcade instance as the owner and controller
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
				Port:       4004,
				Protocol:   "TCP",
				TargetPort: intstr.FromInt(4004),
			}},
			Selector: ls,
		},
	}

	return service
}

// routeForArcade returns the route resource
func (r *ArcadeReconciler) routeForArcade(m *arcadev1alpha1.Arcade) *routev1.Route {
	ls := labelsForArcade(m.Name)
	route := &routev1.Route{
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
				TargetPort: intstr.FromString(m.Name),
			},
			TLS: &routev1.TLSConfig{
				Termination: routev1.TLSTerminationEdge,
			},
		},
	}

	// Set MobileSecurityService mss as the owner and controller
	ctrl.SetControllerReference(m, route, r.Scheme)
	return route
}

// labelsForArcade returns the labels for selecting the resources
// belonging to the given arcade CR name.
func labelsForArcade(name string) map[string]string {
	return map[string]string{"app": "arcade", "arcade_cr": name}
}
