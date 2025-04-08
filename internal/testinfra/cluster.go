package testinfra

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	cliutil "github.com/k3d-io/k3d/v5/cmd/util"
	k3dclient "github.com/k3d-io/k3d/v5/pkg/client"
	"github.com/k3d-io/k3d/v5/pkg/config"
	k3dtypes "github.com/k3d-io/k3d/v5/pkg/config/types"
	k3dconfig "github.com/k3d-io/k3d/v5/pkg/config/v1alpha5"
	k3druntimes "github.com/k3d-io/k3d/v5/pkg/runtimes"
	k3d "github.com/k3d-io/k3d/v5/pkg/types"
)

const (
	defaultImage        = "rancher/k3s:v1.32.3-k3s1"
	defaultClusterName  = "kubeovn-test"
	serverCount         = 3
	defaultNetworkName  = "k3d-kubeovn"
	defaultClusterToken = "k3d-kubeovn"
)

func CreateCluster(ctx context.Context, name string) error {
	cluster := k3dconfig.SimpleConfig{
		ObjectMeta: k3dtypes.ObjectMeta{
			Name: defaultClusterName,
		},
		Servers:      serverCount,
		Image:        defaultImage,
		Network:      defaultNetworkName,
		ClusterToken: defaultClusterToken,
		Options: k3dconfig.SimpleConfigOptions{
			K3dOptions: k3dconfig.SimpleConfigOptionsK3d{
				Wait:                true,
				DisableLoadbalancer: false,
				Timeout:             600 * time.Second,
			},
		},
	}

	var freePort string
	port, err := cliutil.GetFreePort()
	if err != nil {
		return fmt.Errorf("error getting free port during create cluster: %v", err)
	}
	freePort = strconv.Itoa(port)

	cluster.ExposeAPI.HostPort = freePort
	runtime, err := k3druntimes.GetRuntime("docker")
	if err != nil {
		return fmt.Errorf("error fetching docker runtime while creating a cluster: %v", err)
	}

	generatedClusterConfig, err := config.TransformSimpleToClusterConfig(ctx, runtime, cluster, "")
	if err != nil {
		return fmt.Errorf("error transforming simple cluster config: %v", err)
	}
	generatedClusterConfig, err = config.ProcessClusterConfig(*generatedClusterConfig)
	if err != nil {
		return fmt.Errorf("error processing cluster config: %v", err)
	}
	return k3dclient.ClusterRun(ctx, runtime, generatedClusterConfig)
}

func DeleteCluster(ctx context.Context, name string) error {
	runtime, err := k3druntimes.GetRuntime("docker")
	if err != nil {
		return fmt.Errorf("error fetching docker runtime while deleting a cluster: %v", err)
	}
	clusterObj, err := k3dclient.ClusterGet(ctx, runtime, &k3d.Cluster{Name: name})
	if err != nil {
		if errors.Is(err, k3dclient.ClusterGetNoNodesFoundError) {
			return nil
		}
	}

	return k3dclient.ClusterDelete(ctx, runtime, clusterObj, k3d.ClusterDeleteOpts{})
}

func GetKubeconfig(ctx context.Context, name string) ([]byte, error) {
	runtime, err := k3druntimes.GetRuntime("docker")
	if err != nil {
		return nil, fmt.Errorf("error fetching docker runtime while fetching kubeconfig: %v", err)
	}

	clusterList, err := k3dclient.ClusterList(ctx, runtime)
	if err != nil {
		return nil, fmt.Errorf("error listing clusters :%v", err)
	}

	for _, v := range clusterList {
		if v.Name == name {

			clientCfg, err := k3dclient.KubeconfigGet(ctx, runtime, v)
			if err != nil {
				return nil, fmt.Errorf("error getting kubeconfig: %v", err)
			}
			buf := bytes.NewBuffer([]byte{})
			err = k3dclient.KubeconfigWriteToStream(ctx, clientCfg, buf)
			if err != nil {
				return nil, fmt.Errorf("error writing kubeconfig: %v", err)
			}
			return buf.Bytes(), nil
		}
	}
	// no cluster object matching defined name found so return error
	return nil, fmt.Errorf("no cluster name %s found", name)
}
