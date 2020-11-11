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
// +kubebuilder:docs-gen:collapse=Apache License
package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"

	arcadev1alpha1 "github.com/redhat-marketplace/marketplace-games-operator/api/v1alpha1"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("Arcade Controller", func() {
	const (
		ArcadeName            = "test-arcade"
		ArcadeNamespace       = "default"
		ArcadePort      int32 = 4004
		timeout               = time.Second * 10
		interval              = time.Millisecond * 250
	)

	ctx := context.Background()

	Context("When reconciling", func() {
		It("Should successfully attempt an Arcade install", func() {
			By("Creating a new Arcade instance")
			arcade := &arcadev1alpha1.Arcade{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ArcadeName,
					Namespace: ArcadeNamespace,
				},
			}
			Expect(k8sClient.Create(ctx, arcade)).Should(Succeed())

			By("Verifying Arcade created")
			key := types.NamespacedName{Name: ArcadeName, Namespace: ArcadeNamespace}
			createdArcade := &arcadev1alpha1.Arcade{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, key, createdArcade)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdArcade.Name).To(Equal(ArcadeName))
			Expect(createdArcade.Namespace).To(Equal(ArcadeNamespace))

			By("Reconciling")
			result, err := arcadeReconciler.Reconcile(ctrl.Request{
				NamespacedName: key,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(result).ToNot(BeNil())

			By("Verifying deployment was created")
			dep := &appsv1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, key, dep)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(dep.Spec.Template.Spec.Containers[0].Name).To(Equal("arcade"))
			Expect(dep.Spec.Template.Spec.Containers[0].Image).To(ContainSubstring("arcade"))

			By("Verifying service was created")
			svc := &corev1.Service{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, key, svc)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(svc.Spec.Ports[0].Port).To(Equal(ArcadePort))

			By("Verifying route was created")
			rt := &routev1.Route{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, key, rt)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(rt.Spec.Port.TargetPort).To(Equal(intstr.FromString(ArcadeName)))
			Expect(rt.Spec.TLS.Termination).To(Equal(routev1.TLSTerminationEdge))
		})
	})

})
