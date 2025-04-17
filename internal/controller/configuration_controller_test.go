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
	"os"

	kubeovniov1 "github.com/harvester/kubeovn-operator/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/yaml"
)

var (
	config      = &kubeovniov1.Configuration{}
	typedConfig = types.NamespacedName{}
)

var _ = Describe("Configuration Controller", func() {
	Context("When reconciling a resource", func() {

		ctx := context.Background()
		configuration := &kubeovniov1.Configuration{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind Configuration")
			content, err := os.ReadFile("../../config/samples/kubeovn.io_v1_configuration.yaml")
			Expect(err).ShouldNot(HaveOccurred())
			err = yaml.Unmarshal(content, config)
			Expect(err).Should(BeNil())
			typedConfig = types.NamespacedName{Name: config.GetName(), Namespace: config.GetNamespace()}
			err = k8sClient.Get(ctx, typedConfig, configuration)
			if err != nil && errors.IsNotFound(err) {
				Expect(k8sClient.Create(ctx, config)).To(Succeed())
			}
		})

		AfterEach(func() {
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

			By("checking status has been reconcilled to deployed", func() {
				Eventually(func() error {
					resource := &kubeovniov1.Configuration{}
					err := k8sClient.Get(ctx, typedConfig, resource)
					if err != nil {
						return err
					}
					testSuiteLogger.WithValues("current status", resource.Status).Info("current status")
					if resource.Status.Status != kubeovniov1.ConfigurationStatusDeployed {
						return fmt.Errorf("expected to find configuration status to be %s but got %s", kubeovniov1.ConfigurationStatusDeployed, resource.Status.Status)
					}
					fmt.Println(resource.Spec)
					return nil
				}, "30s", "5s").Should(BeNil())
			})

			// ensure all objects defined in templates are actually present in the apiserver
			By("expected objects have been created", func() {

			})

			// trigger upgrade

			// validate new deployments and daemonsets contain the updated images

			// delete a master node to validate cleanup

			// add a new node to ensure ovn db is reconcilled
		})
	})
})
