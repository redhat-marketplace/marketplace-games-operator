/*
Copyright 2020 Red Hat Marketplace contributors.

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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gamev1alpha1 "marketplace.redhat.com/marketplace-games-operator/api/v1alpha1"
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=game.marketplace.redhat.com,resources=instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=game.marketplace.redhat.com,resources=instances/status,verbs=get;update;patch

func (r *InstanceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("instance", req.NamespacedName)

	// TODO: Verify the req object still exists (harder)

	// TODO: Set Status Condition Install to False if not set (harder)

	// TODO: Deploy a deployment - use a dummy docker image for now (easier)

	// pull in a appsv1.Deployment type, initialize it, client.Get to look for it, client.Create if not found
	// name := req.Name
	// namespace := req.Namespace
	// myDeployment := BuildMyDeployment()
	// err := client.Get(types.NamespacedName{name, namespace}, myDeployment)
	// if errors.IsNotFound(err) {
	//   client.Create(myDeployment)
	// }

	// TODO: Deploy a service for the deployment (easier)

	// TODO: Deploy a ingress for the deployment (easier)

	// TODO: Set Status Condition Install to True (harder)

	// TODO: Add resources to the resource array defined on the status (harder)

	return ctrl.Result{}, nil
}

func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gamev1alpha1.Instance{}).
		Complete(r)
}
