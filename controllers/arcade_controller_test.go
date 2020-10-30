/*
Ideally, we should have one `<kind>_conroller_test.go` for each controller scaffolded and called in the `test_suite.go`.
So, let's write our example test for the CronJob controller (`cronjob_controller_test.go.`)
*/

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

/*
As usual, we start with the necessary imports. We also define some utility variables.
*/
package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	arcadev1alpha1 "github.com/redhat-marketplace/marketplace-games-operator/api/v1alpha1"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("Reconciler", func() {
	defer GinkgoRecover()

	const (
		ArcadeName      = "test-arcade"
		ArcadeNamespace = "default"
		timeout         = time.Second * 10
		interval        = time.Millisecond * 250
	)

	ctx := context.Background()

	BeforeEach(func() {
		// reconciled = make(chan reconcile.Request)
		Expect(cfg).NotTo(BeNil())
	})

	Context("Arcade", func() {
		It("Should successfully create Arcade CustomResource (CR)", func() {
			By("Creating a new Arcade CR")
			// Create Arcade
			arcade := &arcadev1alpha1.Arcade{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "game.marketplace.redhat.com/v1alpha1",
					Kind:       "Arcade",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ArcadeName,
					Namespace: ArcadeNamespace,
				},
				Spec: arcadev1alpha1.ArcadeSpec{
					Foo: "bar",
				},
			}
			Expect(k8sClient.Create(ctx, arcade)).Should(Succeed())

			By("Validate Arcade CR was created")
			// Look up Arcade CR
			lookupKey := types.NamespacedName{Name: ArcadeName, Namespace: ArcadeNamespace}
			createdArcade := &arcadev1alpha1.Arcade{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdArcade)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdArcade.Name).To(Equal(arcade.Name))
			Expect(createdArcade.Namespace).To(Equal(arcade.Namespace))
			Expect(createdArcade.Spec.Foo).To(Equal(arcade.Spec.Foo))
		})
	})

})
