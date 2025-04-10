/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	kubeovniov1 "github.com/harvester/kubeovn-operator/api/v1"
)

var (
	config = &kubeovniov1.Configuration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeovniov1.DefaultConfigurationName,
			Namespace: defaultKubeovnNamespace,
		},
		Spec: kubeovniov1.ConfigurationSpec{
			MasterNodesLabel: "node-role.kubernetes.io/control-plane=true",
			Networking: kubeovniov1.NetworkingSpec{
				NetStack:    "ipv4",
				NetworkType: "geneve",
				TunnelType:  "vxlan",
			},
			IPv4: kubeovniov1.NetworkStackSpec{
				PodCIDR:               "10.42.0.0/16",
				ServiceCIDR:           "10.43.0.0/16",
				PodGateway:            "10.42.0.1",
				JoinCIDR:              "100.64.0.0/16",
				PingerExternalAddress: "1.1.1.1",
				PingerExternalDomain:  "google.com.",
			},
		},
	}

	typedConfig = types.NamespacedName{Name: config.Name, Namespace: config.Namespace}
)
var _ = Describe("Configuration Controller", func() {
	Context("When reconciling a resource", func() {

		ctx := context.Background()
		configuration := &kubeovniov1.Configuration{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind Configuration")
			err := k8sClient.Get(ctx, typedConfig, configuration)
			if err != nil && errors.IsNotFound(err) {
				Expect(k8sClient.Create(ctx, config)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &kubeovniov1.Configuration{}
			err := k8sClient.Get(ctx, typedConfig, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Configuration")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("reconcile configuration object", func() {

			By("checking baseline conditions have been set", func() {
				Eventually(func() error {
					resource := &kubeovniov1.Configuration{}
					err := k8sClient.Get(ctx, typedConfig, resource)
					if err != nil {
						return err
					}
					testSuiteLogger.WithValues("current status", resource.Status).Info("current status")
					if len(resource.Status.Conditions) != 2 {
						return fmt.Errorf("expected to find 2 baseline conditions")
					}
					return nil
				}, "30s", "5s").Should(BeNil())
			})

			By("checking master nodes have been discovered from status", func() {
				Eventually(func() error {
					resource := &kubeovniov1.Configuration{}
					err := k8sClient.Get(ctx, typedConfig, resource)
					if err != nil {
						return err
					}
					testSuiteLogger.WithValues("current status", resource.Status).Info("current status")
					if len(resource.Status.MatchingNodeAddresses) == 0 {
						return fmt.Errorf("expected to find at least one master node")
					}
					return nil
				}, "30s", "5s").Should(BeNil())
			})
		})
	})
})
