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

package v1

import (
	"testing"

	kubeovnv1 "github.com/harvester/kubeovn-operator/api/v1"
	"github.com/stretchr/testify/require"
)

func Test_ConfigurationDefaults(t *testing.T) {
	config := &kubeovnv1.Configuration{}
	defaulter := &ConfigurationCustomDefaulter{}
	defaulter.ApplyConfigurationDefaults(config)
	assert := require.New(t)
	assert.Equal(config.Spec.OVNCentral, ovnCentralDefaultResourceSpec, "defaults applied")

}
