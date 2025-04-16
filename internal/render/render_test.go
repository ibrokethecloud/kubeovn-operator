package render

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"sigs.k8s.io/yaml"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsv1 "k8s.io/api/apps/v1"

	ovnoperatorv1 "github.com/harvester/kubeovn-operator/api/v1"
	sourcetemplate "github.com/harvester/kubeovn-operator/internal/templates"
)

var config = &ovnoperatorv1.Configuration{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "sample",
		Namespace: "kube-system",
	},
	Spec: ovnoperatorv1.ConfigurationSpec{
		Global: ovnoperatorv1.GlobalSpec{
			Registry: ovnoperatorv1.RegistrySpec{
				ImagePullSecrets: []string{"registry-one-secret", "registry-two-secret"},
			},
		},
	},
}

const (
	sampleConfigurationJSON = `{
    "apiVersion": "kubeovn.io/v1",
    "kind": "Configuration",
    "metadata": {
        "creationTimestamp": "2025-04-16T02:10:18Z",
        "generation": 1,
        "name": "kubeovn",
        "namespace": "kube-system",
        "resourceVersion": "2371970",
        "uid": "d22666a1-00f6-4487-959f-b524a30de3ca"
    },
    "spec": {
        "cniConf": {
            "cniBinDir": "/opt/cni/bin",
            "cniConfFile": "/kube-ovn/01-kube-ovn.conflist",
            "cniConfigDir": "/etc/cni/net.d",
            "cniConfigPriority": "01",
            "localBinDir": "/usr/local/bin"
        },
        "components": {
            "OVSDBConTimeout": 3,
            "OVSDBInactivityTimeout": 10,
            "checkGateway": true,
            "enableANP": false,
            "enableBindLocalIP": true,
            "enableExternalVPC": true,
            "enableIC": false,
            "enableKeepVMIP": true,
            "enableLB": true,
            "enableLBSVC": false,
            "enableLiveMigrationOptimize": true,
            "enableNATGateway": true,
            "enableNP": true,
            "enableOVNIPSec": false,
            "enableTProxy": false,
            "hardwareOffload": false,
            "logicalGateway": false,
            "lsCtSkipOstLportIPS": true,
            "lsDnatModDlDst": true,
            "secureServing": false,
            "setVLANTxOff": false,
            "u2oInterconnection": false
        },
        "debug": {
            "mirrorInterface": "mirror0"
        },
        "dpdkCPU": "0",
        "dpdkMEMORY": "0",
        "dpdkVersion": "19.11",
        "dualStack": {
            "joinCIDR": "fd00:100:64::/112",
            "pingerExternalAddress": "2606:4700:4700::1111",
            "pingerExternalDomain": "google.com.",
            "podCIDR": "fd00:10:16::/112",
            "podGateway": "fd00:10:16::1",
            "serviceCIDR": "fd00:10:96::/112"
        },
        "global": {
            "images": {
                "kubeovn": {
                    "dpdkRepository": "kube-ovn-dpdk",
                    "repository": "kube-ovn",
                    "supportArm": true,
                    "tag": "v1.14.0",
                    "thirdParty": true,
                    "vpcRepository": "vpc-nat-gateway"
                }
            },
            "registry": {
                "address": "docker.io/kubeovn"
            }
        },
        "hugePages": "0",
        "hugepageSizeType": "hugepages-2Mi",
        "imagePullPolicy": "IfNotPresent",
        "ipv4": {
            "joinCIDR": "100.64.0.0/16",
            "pingerExternalAddress": "1.1.1.1",
            "pingerExternalDomain": "google.com.",
            "podCIDR": "10.42.0.0/16",
            "podGateway": "10.42.0.1",
            "serviceCIDR": "10.43.0.0/16"
        },
        "ipv6": {
            "joinCIDR": "fd00:100:64::/112",
            "pingerExternalAddress": "2606:4700:4700::1111",
            "pingerExternalDomain": "google.com.",
            "podCIDR": "fd00:10:16::/112",
            "podGateway": "fd00:10:16::1",
            "serviceCIDR": "fd00:10:96::/112"
        },
        "kubeOvnCNI": {
            "limits": {
                "cpu": "0",
                "memory": "0"
            },
            "requests": {
                "cpu": "0",
                "memory": "0"
            }
        },
        "kubeOvnController": {
            "limits": {
                "cpu": "0",
                "memory": "0"
            },
            "requests": {
                "cpu": "0",
                "memory": "0"
            }
        },
        "kubeOvnMonitor": {
            "limits": {
                "cpu": "0",
                "memory": "0"
            },
            "requests": {
                "cpu": "0",
                "memory": "0"
            }
        },
        "kubeOvnPinger": {
            "limits": {
                "cpu": "0",
                "memory": "0"
            },
            "requests": {
                "cpu": "0",
                "memory": "0"
            }
        },
        "kubeletConfig": {
            "kubeletDir": "/var/lib/kubelet"
        },
        "logConfig": {
            "logDir": "/var/log"
        },
        "masterNodesLabel": "node-role.kubernetes.io/control-plane=true",
        "networking": {
            "defaultSubnet": "ovn-default",
            "defaultVPC": "ovn-cluster",
            "enableECMP": false,
            "enableEIPSNAT": true,
            "enableMetrics": true,
            "enableSSL": false,
            "netStack": "ipv4",
            "networkType": "geneve",
            "nodeSubnet": "join",
            "ovnLeaderProbeInterval": 5,
            "ovnNorthdNThreads": 1,
            "ovnNorthdProbeInterval": 5000,
            "ovnRemoteOpenflowInterval": 10,
            "ovnRemoteProbeInterval": 10000,
            "podNicType": "veth-pair",
            "probeInterval": 180000,
            "tunnelType": "vxlan",
            "vlan": {
                "providerName": "provider",
                "vlanId": 1,
                "vlanName": "ovn-vlan"
            }
        },
        "openVSwitchDir": "/etc/origin/openvswitch",
        "ovnCentral": {
            "limits": {
                "cpu": "0",
                "memory": "0"
            },
            "requests": {
                "cpu": "0",
                "memory": "0"
            }
        },
        "ovnDir": "/etc/origin/ovn",
        "ovsOVN": {
            "limits": {
                "cpu": "0",
                "memory": "0"
            },
            "requests": {
                "cpu": "0",
                "memory": "0"
            }
        },
        "performance": {
            "gcInterval": 360,
            "inspectInterval": 20,
            "ovsVSCtlConcurrency": 100
        }
    },
    "status": {
        "conditions": [
            {
                "lastTransitionTime": "2025-04-16T05:33:01Z",
                "message": "",
                "observedGeneration": 1,
                "reason": "Unknown",
                "status": "Unknown",
                "type": "erroredObjects"
            },
            {
                "lastTransitionTime": "2025-04-16T05:33:01Z",
                "message": "",
                "observedGeneration": 1,
                "reason": "Unknown",
                "status": "True",
                "type": "waitingForMatchignNodes"
            }
        ],
        "matchingNodeAddresses": [
            "172.18.0.2",
            "172.18.0.3",
            "172.18.0.4"
        ],
        "status": "Deployed"
    }
}`
)

func Test_GenerateSAObjects(t *testing.T) {
	sa := corev1.ServiceAccount{}
	assert := require.New(t)
	returnedObjects, err := GenerateObjects(sourcetemplate.ServiceAccountList, config, &sa, nil)
	assert.NoError(err)
	assert.Len(returnedObjects, 4)

}

func Test_GenerateCRD(t *testing.T) {
	crd := apiextensions.CustomResourceDefinition{}
	assert := require.New(t)
	returnedObjects, err := GenerateObjects(sourcetemplate.CRDList, config, &crd, nil)
	assert.NoError(err)
	assert.Equal(len(returnedObjects), len(sourcetemplate.CRDList))
}

func Test_GenerateSecret(t *testing.T) {
	secret := corev1.Secret{}
	assert := require.New(t)
	_, err := GenerateObjects(sourcetemplate.SecretList, config, &secret, nil)
	assert.NoError(err)
}

func Test_generateMap(t *testing.T) {
	assert := require.New(t)
	_, err := generateMap(config)
	assert.NoError(err)
}

func Test_deployments(t *testing.T) {
	assert := require.New(t)
	c, err := generateConfigObject()
	assert.NoError(err, "expected no error while generating config object")
	deployment := appsv1.Deployment{}
	returnedObjects, err := GenerateObjects(sourcetemplate.DeploymentList, c, &deployment, nil)
	assert.NoError(err)
	for _, v := range returnedObjects {
		fmt.Println(v.GetName())
		if v.GetName() == "kube-ovn-monitor" {
			fmt.Println(v)
			var out []byte
			out, err = yaml.Marshal(v)
			assert.NoError(err)
			fmt.Println(string(out))
		}
	}
}

func Test_daemonsets(t *testing.T) {
	assert := require.New(t)
	c, err := generateConfigObject()
	assert.NoError(err, "expected no error while generating config object")
	ds := appsv1.DaemonSet{}
	returnedObjects, err := GenerateObjects(sourcetemplate.DaemonsetList, c, &ds, nil)
	assert.NoError(err)
	for _, v := range returnedObjects {
		if v.GetName() == "kube-ovn-pinger" {
			var out []byte
			out, err = yaml.Marshal(v)
			assert.NoError(err)
			fmt.Println(string(out))
		}
	}
}
func generateConfigObject() (*ovnoperatorv1.Configuration, error) {
	config := &ovnoperatorv1.Configuration{}
	err := json.Unmarshal([]byte(sampleConfigurationJSON), config)
	return config, err
}

func Test_ObjectRendering(t *testing.T) {
	assert := require.New(t)
	c, err := generateConfigObject()
	assert.NoError(err, "expected no error while generating config object")
	for objectType, objectList := range sourcetemplate.OrderedObjectList {
		returnedObjects, err := GenerateObjects(objectList, c, objectType, nil)
		assert.NoError(err, "expected no error while generating object", objectType)
		for _, object := range returnedObjects {
			assert.NotEmpty(object.GetName())
		}
	}
}

var sampleDeployment = `kind: Deployment
apiVersion: apps/v1
metadata:
  name: kube-ovn-monitor
  namespace: kube-system
  annotations:
    kubernetes.io/description: |
      Metrics for OVN components: northd, nb and sb.
spec:
  replicas: 1
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  selector:
    matchLabels:
      app: kube-ovn-monitor
  template:
    metadata:
      labels:
        app: kube-ovn-monitor
        component: network
        type: infra
    spec:
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - key: CriticalAddonsOnly
          operator: Exists
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchLabels:
                  app: kube-ovn-monitor
              topologyKey: kubernetes.io/hostname
      priorityClassName: system-cluster-critical
      serviceAccountName: kube-ovn-app
      hostNetwork: true
      initContainers:
        - name: hostpath-init
          image: docker.io/kubeovn/kube-ovn:v1.14.0
          imagePullPolicy: IfNotPresent
          command:
            - sh
            - -c
            - "chown -R nobody: /var/log/kube-ovn"
          securityContext:
            allowPrivilegeEscalation: true
            capabilities:
              drop:
                - ALL
            privileged: true
            runAsUser: 0
          volumeMounts:
            - name: kube-ovn-log
              mountPath: /var/log/kube-ovn
      containers:
        - name: kube-ovn-monitor
          image: docker.io/kubeovn/kube-ovn:v1.14.0
          imagePullPolicy: IfNotPresent
          command: ["/kube-ovn/start-ovn-monitor.sh"]
          args:
          - --secure-serving=false
          - --log_file=/var/log/kube-ovn/kube-ovn-monitor.log
          - --logtostderr=false
          - --alsologtostderr=true
          - --log_file_max_size=200
          - --enable-metrics=<no value>
          securityContext:
            runAsUser: 
            privileged: false
            capabilities:
              add:
                - NET_BIND_SERVICE
          env:
            - name: ENABLE_SSL
              value: "false"
            - name: KUBE_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: POD_IPS
              valueFrom:
                fieldRef:
                  fieldPath: status.podIPs
            - name: ENABLE_BIND_LOCAL_IP
              value: "true"
          resources:
            requests:
              cpu: 0
              memory: 0
            limits:
              cpu: 0
              memory: 0
          volumeMounts:
            - mountPath: /var/run/ovn
              name: host-run-ovn
            - mountPath: /etc/ovn
              name: host-config-ovn
            - mountPath: /var/log/ovn
              name: host-log-ovn
              readOnly: true
            - mountPath: /etc/localtime
              name: localtime
              readOnly: true
            - mountPath: /var/run/tls
              name: kube-ovn-tls
            - mountPath: /var/log/kube-ovn
              name: kube-ovn-log
          livenessProbe:
            failureThreshold: 3
            initialDelaySeconds: 30
            periodSeconds: 7
            successThreshold: 1
            httpGet:
              port: 10661
              path: /livez
              scheme: 'HTTP'
            timeoutSeconds: 5
          readinessProbe:
            failureThreshold: 3
            initialDelaySeconds: 30
            periodSeconds: 7
            successThreshold: 1
            httpGet:
              port: 10661
              path: /readyz
              scheme: 'HTTP'
            timeoutSeconds: 5
      nodeSelector:
        kubernetes.io/os: "linux"
        node-role.kubernetes.io/control-plane: "true"
      volumes:
        - name: host-run-ovn
          hostPath:
            path: /run/ovn
        - name: host-config-ovn
          hostPath:
            path: /etc/origin/ovn
        - name: host-log-ovn
          hostPath:
            path: /var/log/ovn
        - name: localtime
          hostPath:
            path: /etc/localtime
        - name: kube-ovn-tls
          secret:
            optional: true
            secretName: kube-ovn-tls
        - name: kube-ovn-log
          hostPath:
            path: /var/log/kube-ovn
`

func Test_renderDeployment(t *testing.T) {
	assert := require.New(t)
	d := &appsv1.Deployment{}
	assert.NoError(yaml.Unmarshal([]byte(sampleDeployment), d))
	t.Log(d)
	out, err := yaml.Marshal(d)
	assert.NoError(err)
	t.Log(string(out))
}
