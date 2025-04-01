package templates

var (
	system_ovn_clusterrole = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.k8s.io/system-only: "true"
  name: system:ovn
rules:
  - apiGroups:
      - "kubeovn.io"
    resources:
      - vpcs
      - vpcs/status
      - vpc-nat-gateways
      - vpc-nat-gateways/status
      - subnets
      - subnets/status
      - ippools
      - ippools/status
      - ips
      - vips
      - vips/status
      - vlans
      - vlans/status
      - provider-networks
      - provider-networks/status
      - security-groups
      - security-groups/status
      - iptables-eips
      - iptables-fip-rules
      - iptables-dnat-rules
      - iptables-snat-rules
      - iptables-eips/status
      - iptables-fip-rules/status
      - iptables-dnat-rules/status
      - iptables-snat-rules/status
      - ovn-eips
      - ovn-fips
      - ovn-snat-rules
      - ovn-eips/status
      - ovn-fips/status
      - ovn-snat-rules/status
      - ovn-dnat-rules
      - ovn-dnat-rules/status
      - switch-lb-rules
      - switch-lb-rules/status
      - vpc-dnses
      - vpc-dnses/status
      - qos-policies
      - qos-policies/status
    verbs:
      - "*"
  - apiGroups:
      - ""
    resources:
      - pods
      - namespaces
    verbs:
      - get
      - list
      - patch
      - watch
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - pods/exec
    verbs:
      - create
  - apiGroups:
      - "k8s.cni.cncf.io"
    resources:
      - network-attachment-definitions
    verbs:
      - get
  - apiGroups:
      - ""
      - networking.k8s.io
    resources:
      - networkpolicies
      - configmaps
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - apps
    resources:
      - daemonsets
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - services
      - services/status
    verbs:
      - get
      - list
      - update
      - patch
      - create
      - delete
      - watch
  - apiGroups:
      - ""
    resources:
      - endpoints
    verbs:
      - create
      - update
      - get
      - list
      - watch
  - apiGroups:
      - apps
    resources:
      - statefulsets
      - deployments
      - deployments/scale
    verbs:
      - get
      - list
      - create
      - delete
      - update
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
  - apiGroups:
      - coordination.k8s.io
    resources:
      - leases
    verbs:
      - "*"
  - apiGroups:
      - "kubevirt.io"
    resources:
      - virtualmachines
      - virtualmachineinstances
    verbs:
      - get
      - list
  - apiGroups:
      - "policy.networking.k8s.io"
    resources:
      - adminnetworkpolicies
      - baselineadminnetworkpolicies
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - authentication.k8s.io
    resources:
      - tokenreviews
    verbs:
      - create
  - apiGroups:
      - authorization.k8s.io
    resources:
      - subjectaccessreviews
    verbs:
      - create
  - apiGroups: 
      - "certificates.k8s.io"
    resources: 
      - "certificatesigningrequests"
    verbs: 
      - "get"
      - "list"
      - "watch"
  - apiGroups:
    - certificates.k8s.io
    resources:
    - certificatesigningrequests/status
    - certificatesigningrequests/approval
    verbs:
    - update
  - apiGroups:
    - ""
    resources:
    - secrets
    verbs:
    - get
    - create
  - apiGroups:
    - certificates.k8s.io
    resourceNames:
    - kubeovn.io/signer
    resources:
    - signers
    verbs:
    - approve
    - sign
  - apiGroups:
      - kubevirt.io
    resources:
      - virtualmachineinstancemigrations
    verbs:
      - "list"
      - "watch"
      - "get"
  - apiGroups:
      - apiextensions.k8s.io
    resources:
      - customresourcedefinitions
    verbs:
      - get`

	system_ovn_ovs_clusterrole = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.k8s.io/system-only: "true"
  name: system:ovn-ovs
rules:
  - apiGroups:
      - ""
    resources:
      - pods
    verbs:
      - get
      - patch
  - apiGroups:
      - ""
    resources:
      - services
      - endpoints
    verbs:
      - get
  - apiGroups:
      - apps
    resources:
      - controllerrevisions
    verbs:
      - get
      - list`

	system_kube_ovn_cni_clusterrole = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.k8s.io/system-only: "true"
  name: system:kube-ovn-cni
rules:
  - apiGroups:
      - "kubeovn.io"
    resources:
      - subnets
      - vlans
      - provider-networks
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - "kubeovn.io"
    resources:
      - ovn-eips
      - ovn-eips/status
      - nodes
      - nodes/status
      - pods
    verbs:
      - get
      - list
      - patch
      - watch
  - apiGroups:
      - "kubeovn.io"
    resources:
      - ips
    verbs:
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
  - apiGroups:
      - ""
    resources:
      - configmaps
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - authentication.k8s.io
    resources:
      - tokenreviews
    verbs:
      - create
  - apiGroups:
      - authorization.k8s.io
    resources:
      - subjectaccessreviews
    verbs:
      - create
  - apiGroups: 
      - "certificates.k8s.io"
    resources: 
      - "certificatesigningrequests"
    verbs: 
      - "create" 
      - "get"
      - "list"
      - "watch"
      - "delete"
  - apiGroups:
      - ""
    resources:
      - "secrets"
    verbs:
      - "get"`

	system_kube_ovn_app_clusterrole = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.k8s.io/system-only: "true"
  name: system:kube-ovn-app
rules:
  - apiGroups:
      - ""
    resources:
      - pods
      - nodes
    verbs:
      - get
      - list
  - apiGroups:
      - apps
    resources:
      - daemonsets
    verbs:
      - get
  - apiGroups:
      - authentication.k8s.io
    resources:
      - tokenreviews
    verbs:
      - create
  - apiGroups:
      - authorization.k8s.io
    resources:
      - subjectaccessreviews
    verbs:
      - create`

	ClusterRoleList = []string{system_ovn_clusterrole, system_ovn_ovs_clusterrole, system_kube_ovn_cni_clusterrole, system_kube_ovn_app_clusterrole}
)
