package templates

var (
	vpc_dnses_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: vpc-dnses.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: vpc-dnses
    singular: vpc-dns
    shortNames:
      - vpc-dns
    kind: VpcDns
    listKind: VpcDnsList
  scope: Cluster
  versions:
    - additionalPrinterColumns:
        - jsonPath: .status.active
          name: Active
          type: boolean
        - jsonPath: .spec.vpc
          name: Vpc
          type: string
        - jsonPath: .spec.subnet
          name: Subnet
          type: string
      name: v1
      served: true
      storage: true
      subresources:
        status: {}
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                vpc:
                  type: string
                subnet:
                  type: string
                replicas:
                  type: integer
                  minimum: 1
                  maximum: 3
            status:
              type: object
              properties:
                active:
                  type: boolean
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string`
	switch_lb_rules_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: switch-lb-rules.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: switch-lb-rules
    singular: switch-lb-rule
    shortNames:
      - slr
    kind: SwitchLBRule
    listKind: SwitchLBRuleList
  scope: Cluster
  versions:
    - additionalPrinterColumns:
        - jsonPath: .spec.vip
          name: vip
          type: string
        - jsonPath: .status.ports
          name: port(s)
          type: string
        - jsonPath: .status.service
          name: service
          type: string
        - jsonPath: .metadata.creationTimestamp
          name: age
          type: date
      name: v1
      served: true
      storage: true
      subresources:
        status: {}
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                namespace:
                  type: string
                vip:
                  type: string
                sessionAffinity:
                  type: string
                ports:
                  items:
                    properties:
                      name:
                        type: string
                      port:
                        type: integer
                        minimum: 1
                        maximum: 65535
                      protocol:
                        type: string
                      targetPort:
                        type: integer
                        minimum: 1
                        maximum: 65535
                    type: object
                  type: array
                selector:
                  items:
                    type: string
                  type: array
                endpoints:
                  items:
                    type: string
                  type: array
            status:
              type: object
              properties:
                ports:
                  type: string
                service:
                  type: string`

	vpc_nat_gateways_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: vpc-nat-gateways.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: vpc-nat-gateways
    singular: vpc-nat-gateway
    shortNames:
      - vpc-nat-gw
    kind: VpcNatGateway
    listKind: VpcNatGatewayList
  scope: Cluster
  versions:
    - additionalPrinterColumns:
        - jsonPath: .spec.vpc
          name: Vpc
          type: string
        - jsonPath: .spec.subnet
          name: Subnet
          type: string
        - jsonPath: .spec.lanIp
          name: LanIP
          type: string
      name: v1
      served: true
      storage: true
      subresources:
        status: {}
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                externalSubnets:
                  items:
                    type: string
                  type: array
                selector:
                  type: array
                  items:
                    type: string
                qosPolicy:
                  type: string
                tolerations:
                  type: array
                  items:
                    type: object
                    properties:
                      key:
                        type: string
                      operator:
                        type: string
                        enum:
                          - Equal
                          - Exists
                      value:
                        type: string
                      effect:
                        type: string
                        enum:
                          - NoExecute
                          - NoSchedule
                          - PreferNoSchedule
                      tolerationSeconds:
                        type: integer
                affinity:
                  properties:
                    nodeAffinity:
                      properties:
                        preferredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              preference:
                                properties:
                                  matchExpressions:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                  matchFields:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                type: object
                              weight:
                                format: int32
                                type: integer
                            required:
                              - preference
                              - weight
                            type: object
                          type: array
                        requiredDuringSchedulingIgnoredDuringExecution:
                          properties:
                            nodeSelectorTerms:
                              items:
                                properties:
                                  matchExpressions:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                  matchFields:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                type: object
                              type: array
                          required:
                            - nodeSelectorTerms
                          type: object
                      type: object
                    podAffinity:
                      properties:
                        preferredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              podAffinityTerm:
                                properties:
                                  labelSelector:
                                    properties:
                                      matchExpressions:
                                        items:
                                          properties:
                                            key:
                                              type: string
                                              x-kubernetes-patch-strategy: merge
                                              x-kubernetes-patch-merge-key: key
                                            operator:
                                              type: string
                                            values:
                                              items:
                                                type: string
                                              type: array
                                          required:
                                            - key
                                            - operator
                                          type: object
                                        type: array
                                      matchLabels:
                                        additionalProperties:
                                          type: string
                                        type: object
                                    type: object
                                  namespaces:
                                    items:
                                      type: string
                                    type: array
                                  topologyKey:
                                    type: string
                                required:
                                  - topologyKey
                                type: object
                              weight:
                                format: int32
                                type: integer
                            required:
                              - podAffinityTerm
                              - weight
                            type: object
                          type: array
                        requiredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              labelSelector:
                                properties:
                                  matchExpressions:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                          x-kubernetes-patch-strategy: merge
                                          x-kubernetes-patch-merge-key: key
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                  matchLabels:
                                    additionalProperties:
                                      type: string
                                    type: object
                                type: object
                              namespaces:
                                items:
                                  type: string
                                type: array
                              topologyKey:
                                type: string
                            required:
                              - topologyKey
                            type: object
                          type: array
                      type: object
                    podAntiAffinity:
                      properties:
                        preferredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              podAffinityTerm:
                                properties:
                                  labelSelector:
                                    properties:
                                      matchExpressions:
                                        items:
                                          properties:
                                            key:
                                              type: string
                                              x-kubernetes-patch-strategy: merge
                                              x-kubernetes-patch-merge-key: key
                                            operator:
                                              type: string
                                            values:
                                              items:
                                                type: string
                                              type: array
                                          required:
                                            - key
                                            - operator
                                          type: object
                                        type: array
                                      matchLabels:
                                        additionalProperties:
                                          type: string
                                        type: object
                                    type: object
                                  namespaces:
                                    items:
                                      type: string
                                    type: array
                                  topologyKey:
                                    type: string
                                required:
                                  - topologyKey
                                type: object
                              weight:
                                format: int32
                                type: integer
                            required:
                              - podAffinityTerm
                              - weight
                            type: object
                          type: array
                        requiredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              labelSelector:
                                properties:
                                  matchExpressions:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                          x-kubernetes-patch-strategy: merge
                                          x-kubernetes-patch-merge-key: key
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                  matchLabels:
                                    additionalProperties:
                                      type: string
                                    type: object
                                type: object
                              namespaces:
                                items:
                                  type: string
                                type: array
                              topologyKey:
                                type: string
                            required:
                              - topologyKey
                            type: object
                          type: array
                      type: object
                  type: object
            spec:
              type: object
              properties:
                lanIp:
                  type: string
                subnet:
                  type: string
                externalSubnets:
                  items:
                    type: string
                  type: array
                vpc:
                  type: string
                selector:
                  type: array
                  items:
                    type: string
                qosPolicy:
                  type: string
                bgpSpeaker:
                  type: object
                  properties:
                    enabled:
                      type: boolean
                    asn:
                      type: integer
                    remoteAsn:
                      type: integer
                    neighbors:
                      type: array
                      items:
                        type: string
                    holdTime:
                      type: string
                    routerId:
                      type: string
                    password:
                      type: string
                    enableGracefulRestart:
                      type: boolean
                    extraArgs:
                      type: array
                      items:
                        type: string
                tolerations:
                  type: array
                  items:
                    type: object
                    properties:
                      key:
                        type: string
                      operator:
                        type: string
                        enum:
                          - Equal
                          - Exists
                      value:
                        type: string
                      effect:
                        type: string
                        enum:
                          - NoExecute
                          - NoSchedule
                          - PreferNoSchedule
                      tolerationSeconds:
                        type: integer
                affinity:
                  properties:
                    nodeAffinity:
                      properties:
                        preferredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              preference:
                                properties:
                                  matchExpressions:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                  matchFields:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                type: object
                              weight:
                                format: int32
                                type: integer
                            required:
                              - preference
                              - weight
                            type: object
                          type: array
                        requiredDuringSchedulingIgnoredDuringExecution:
                          properties:
                            nodeSelectorTerms:
                              items:
                                properties:
                                  matchExpressions:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                  matchFields:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                type: object
                              type: array
                          required:
                            - nodeSelectorTerms
                          type: object
                      type: object
                    podAffinity:
                      properties:
                        preferredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              podAffinityTerm:
                                properties:
                                  labelSelector:
                                    properties:
                                      matchExpressions:
                                        items:
                                          properties:
                                            key:
                                              type: string
                                              x-kubernetes-patch-strategy: merge
                                              x-kubernetes-patch-merge-key: key
                                            operator:
                                              type: string
                                            values:
                                              items:
                                                type: string
                                              type: array
                                          required:
                                            - key
                                            - operator
                                          type: object
                                        type: array
                                      matchLabels:
                                        additionalProperties:
                                          type: string
                                        type: object
                                    type: object
                                  namespaces:
                                    items:
                                      type: string
                                    type: array
                                  topologyKey:
                                    type: string
                                required:
                                  - topologyKey
                                type: object
                              weight:
                                format: int32
                                type: integer
                            required:
                              - podAffinityTerm
                              - weight
                            type: object
                          type: array
                        requiredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              labelSelector:
                                properties:
                                  matchExpressions:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                          x-kubernetes-patch-strategy: merge
                                          x-kubernetes-patch-merge-key: key
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                  matchLabels:
                                    additionalProperties:
                                      type: string
                                    type: object
                                type: object
                              namespaces:
                                items:
                                  type: string
                                type: array
                              topologyKey:
                                type: string
                            required:
                              - topologyKey
                            type: object
                          type: array
                      type: object
                    podAntiAffinity:
                      properties:
                        preferredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              podAffinityTerm:
                                properties:
                                  labelSelector:
                                    properties:
                                      matchExpressions:
                                        items:
                                          properties:
                                            key:
                                              type: string
                                              x-kubernetes-patch-strategy: merge
                                              x-kubernetes-patch-merge-key: key
                                            operator:
                                              type: string
                                            values:
                                              items:
                                                type: string
                                              type: array
                                          required:
                                            - key
                                            - operator
                                          type: object
                                        type: array
                                      matchLabels:
                                        additionalProperties:
                                          type: string
                                        type: object
                                    type: object
                                  namespaces:
                                    items:
                                      type: string
                                    type: array
                                  topologyKey:
                                    type: string
                                required:
                                  - topologyKey
                                type: object
                              weight:
                                format: int32
                                type: integer
                            required:
                              - podAffinityTerm
                              - weight
                            type: object
                          type: array
                        requiredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              labelSelector:
                                properties:
                                  matchExpressions:
                                    items:
                                      properties:
                                        key:
                                          type: string
                                          x-kubernetes-patch-strategy: merge
                                          x-kubernetes-patch-merge-key: key
                                        operator:
                                          type: string
                                        values:
                                          items:
                                            type: string
                                          type: array
                                      required:
                                        - key
                                        - operator
                                      type: object
                                    type: array
                                  matchLabels:
                                    additionalProperties:
                                      type: string
                                    type: object
                                type: object
                              namespaces:
                                items:
                                  type: string
                                type: array
                              topologyKey:
                                type: string
                            required:
                              - topologyKey
                            type: object
                          type: array
                      type: object
                  type: object`

	iptables_eips_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: iptables-eips.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: iptables-eips
    singular: iptables-eip
    shortNames:
      - eip
    kind: IptablesEIP
    listKind: IptablesEIPList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
      - jsonPath: .status.ip
        name: IP
        type: string
      - jsonPath: .spec.macAddress
        name: Mac
        type: string
      - jsonPath: .status.nat
        name: Nat
        type: string
      - jsonPath: .spec.natGwDp
        name: NatGwDp
        type: string
      - jsonPath: .status.ready
        name: Ready
        type: boolean
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                ready:
                  type: boolean
                ip:
                  type: string
                nat:
                  type: string
                redo:
                  type: string
                qosPolicy:
                  type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                v4ip:
                  type: string
                v6ip:
                  type: string
                macAddress:
                  type: string
                natGwDp:
                  type: string
                qosPolicy:
                  type: string
                externalSubnet:
                  type: string`

	iptables_fip_rules_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: iptables-fip-rules.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: iptables-fip-rules
    singular: iptables-fip-rule
    shortNames:
      - fip
    kind: IptablesFIPRule
    listKind: IptablesFIPRuleList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
      - jsonPath: .spec.eip
        name: Eip
        type: string
      - jsonPath: .status.v4ip
        name: V4ip
        type: string
      - jsonPath: .spec.internalIp
        name: InternalIp
        type: string
      - jsonPath: .status.v6ip
        name: V6ip
        type: string
      - jsonPath: .status.ready
        name: Ready
        type: boolean
      - jsonPath: .status.natGwDp
        name: NatGwDp
        type: string
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                ready:
                  type: boolean
                v4ip:
                  type: string
                v6ip:
                  type: string
                natGwDp:
                  type: string
                redo:
                  type: string
                internalIp:
                  type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                eip:
                  type: string
                internalIp:
                  type: string`

	iptables_dnat_rules_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: iptables-dnat-rules.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: iptables-dnat-rules
    singular: iptables-dnat-rule
    shortNames:
      - dnat
    kind: IptablesDnatRule
    listKind: IptablesDnatRuleList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
      - jsonPath: .spec.eip
        name: Eip
        type: string
      - jsonPath: .spec.protocol
        name: Protocol
        type: string
      - jsonPath: .status.v4ip
        name: V4ip
        type: string
      - jsonPath: .status.v6ip
        name: V6ip
        type: string
      - jsonPath: .spec.internalIp
        name: InternalIp
        type: string
      - jsonPath: .spec.externalPort
        name: ExternalPort
        type: string
      - jsonPath: .spec.internalPort
        name: InternalPort
        type: string
      - jsonPath: .status.natGwDp
        name: NatGwDp
        type: string
      - jsonPath: .status.ready
        name: Ready
        type: boolean
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                ready:
                  type: boolean
                v4ip:
                  type: string
                v6ip:
                  type: string
                natGwDp:
                  type: string
                redo:
                  type: string
                protocol:
                  type: string
                internalIp:
                  type: string
                internalPort:
                  type: string
                externalPort:
                  type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                eip:
                  type: string
                externalPort:
                  type: string
                protocol:
                  type: string
                internalIp:
                  type: string
                internalPort:
                  type: string`

	iptables_snat_rules_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: iptables-snat-rules.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: iptables-snat-rules
    singular: iptables-snat-rule
    shortNames:
      - snat
    kind: IptablesSnatRule
    listKind: IptablesSnatRuleList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
      - jsonPath: .spec.eip
        name: EIP
        type: string
      - jsonPath: .status.v4ip
        name: V4ip
        type: string
      - jsonPath: .status.v6ip
        name: V6ip
        type: string
      - jsonPath: .spec.internalCIDR
        name: InternalCIDR
        type: string
      - jsonPath: .status.natGwDp
        name: NatGwDp
        type: string
      - jsonPath: .status.ready
        name: Ready
        type: boolean
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                ready:
                  type: boolean
                v4ip:
                  type: string
                v6ip:
                  type: string
                natGwDp:
                  type: string
                redo:
                  type: string
                internalCIDR:
                  type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                eip:
                  type: string
                internalCIDR:
                  type: string`
	ovn_eips_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ovn-eips.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: ovn-eips
    singular: ovn-eip
    shortNames:
      - oeip
    kind: OvnEip
    listKind: OvnEipList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
      - jsonPath: .status.v4Ip
        name: V4IP
        type: string
      - jsonPath: .status.v6Ip
        name: V6IP
        type: string
      - jsonPath: .status.macAddress
        name: Mac
        type: string
      - jsonPath: .status.type
        name: Type
        type: string
      - jsonPath: .status.nat
        name: Nat
        type: string
      - jsonPath: .status.ready
        name: Ready
        type: boolean
      - jsonPath: .spec.externalSubnet
        name: ExternalSubnet
        type: string
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                type:
                  type: string
                nat:
                  type: string
                ready:
                  type: boolean
                v4Ip:
                  type: string
                v6Ip:
                  type: string
                macAddress:
                  type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                externalSubnet:
                  type: string
                type:
                  type: string
                v4Ip:
                  type: string
                v6Ip:
                  type: string
                macAddress:
                  type: string`

	ovn_fips_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ovn-fips.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: ovn-fips
    singular: ovn-fip
    shortNames:
      - ofip
    kind: OvnFip
    listKind: OvnFipList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
      - jsonPath: .status.vpc
        name: Vpc
        type: string
      - jsonPath: .status.v4Eip
        name: V4Eip
        type: string
      - jsonPath: .status.v6Eip
        name: V6Eip
        type: string
      - jsonPath: .status.v4Ip
        name: V4Ip
        type: string
      - jsonPath: .status.v6Ip
        name: V6Ip
        type: string
      - jsonPath: .status.ready
        name: Ready
        type: boolean
      - jsonPath: .spec.ipType
        name: IpType
        type: string
      - jsonPath: .spec.ipName
        name: IpName
        type: string
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                ready:
                  type: boolean
                v4Eip:
                  type: string
                v6Eip:
                  type: string
                v4Ip:
                  type: string
                v6Ip:
                  type: string
                vpc:
                  type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                ovnEip:
                  type: string
                ipType:
                  type: string
                ipName:
                  type: string
                vpc:
                  type: string
                v4Ip:
                  type: string
                v6Ip:
                  type: string`

	ovn_snat_rules_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ovn-snat-rules.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: ovn-snat-rules
    singular: ovn-snat-rule
    shortNames:
      - osnat
    kind: OvnSnatRule
    listKind: OvnSnatRuleList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
      - jsonPath: .status.vpc
        name: Vpc
        type: string
      - jsonPath: .status.v4Eip
        name: V4Eip
        type: string
      - jsonPath: .status.v6Eip
        name: V6Eip
        type: string
      - jsonPath: .status.v4IpCidr
        name: V4IpCidr
        type: string
      - jsonPath: .status.v6IpCidr
        name: V6IpCidr
        type: string
      - jsonPath: .status.ready
        name: Ready
        type: boolean
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                ready:
                  type: boolean
                v4Eip:
                  type: string
                v6Eip:
                  type: string
                v4IpCidr:
                  type: string
                v6IpCidr:
                  type: string
                vpc:
                  type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                ovnEip:
                  type: string
                vpcSubnet:
                  type: string
                ipName:
                  type: string
                vpc:
                  type: string
                v4IpCidr:
                  type: string
                v6IpCidr:
                  type: string`

	ovn_dnat_rules_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ovn-dnat-rules.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: ovn-dnat-rules
    singular: ovn-dnat-rule
    shortNames:
      - odnat
    kind: OvnDnatRule
    listKind: OvnDnatRuleList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
        - jsonPath: .status.vpc
          name: Vpc
          type: string
        - jsonPath: .spec.ovnEip
          name: Eip
          type: string
        - jsonPath: .status.protocol
          name: Protocol
          type: string
        - jsonPath: .status.v4Eip
          name: V4Eip
          type: string
        - jsonPath: .status.v6Eip
          name: V6Eip
          type: string
        - jsonPath: .status.v4Ip
          name: V4Ip
          type: string
        - jsonPath: .status.v6Ip
          name: V6Ip
          type: string
        - jsonPath: .status.internalPort
          name: InternalPort
          type: string
        - jsonPath: .status.externalPort
          name: ExternalPort
          type: string
        - jsonPath: .spec.ipName
          name: IpName
          type: string
        - jsonPath: .status.ready
          name: Ready
          type: boolean
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                ready:
                  type: boolean
                v4Eip:
                  type: string
                v6Eip:
                  type: string
                v4Ip:
                  type: string
                v6Ip:
                  type: string
                vpc:
                  type: string
                externalPort:
                  type: string
                internalPort:
                  type: string
                protocol:
                  type: string
                ipName:
                  type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                ovnEip:
                  type: string
                ipType:
                  type: string
                ipName:
                  type: string
                externalPort:
                  type: string
                internalPort:
                  type: string
                protocol:
                  type: string
                vpc:
                  type: string
                v4Ip:
                  type: string
                v6Ip:
                  type: string`

	vpcs_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: vpcs.kubeovn.io
spec:
  group: kubeovn.io
  versions:
    - additionalPrinterColumns:
        - jsonPath: .status.enableExternal
          name: EnableExternal
          type: boolean
        - jsonPath: .status.enableBfd
          name: EnableBfd
          type: boolean
        - jsonPath: .status.standby
          name: Standby
          type: boolean
        - jsonPath: .status.subnets
          name: Subnets
          type: string
        - jsonPath: .status.extraExternalSubnets
          name: ExtraExternalSubnets
          type: string
        - jsonPath: .spec.namespaces
          name: Namespaces
          type: string
        - jsonPath: .status.defaultLogicalSwitch
          name: DefaultSubnet
          type: string
      name: v1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              properties:
                defaultSubnet:
                  type: string
                enableExternal:
                  type: boolean
                enableBfd:
                  type: boolean
                namespaces:
                  items:
                    type: string
                  type: array
                extraExternalSubnets:
                  items:
                    type: string
                  type: array
                staticRoutes:
                  items:
                    properties:
                      policy:
                        type: string
                      cidr:
                        type: string
                      nextHopIP:
                        type: string
                      ecmpMode:
                        type: string
                      bfdId:
                        type: string
                      routeTable:
                        type: string
                    type: object
                  type: array
                policyRoutes:
                  items:
                    properties:
                      priority:
                        type: integer
                      action:
                        type: string
                      match:
                        type: string
                      nextHopIP:
                        type: string
                    type: object
                  type: array
                vpcPeerings:
                  items:
                    properties:
                      remoteVpc:
                        type: string
                      localConnectIP:
                        type: string
                    type: object
                  type: array
              type: object
            status:
              properties:
                conditions:
                  items:
                    properties:
                      lastTransitionTime:
                        type: string
                      lastUpdateTime:
                        type: string
                      message:
                        type: string
                      reason:
                        type: string
                      status:
                        type: string
                      type:
                        type: string
                    type: object
                  type: array
                default:
                  type: boolean
                defaultLogicalSwitch:
                  type: string
                router:
                  type: string
                standby:
                  type: boolean
                enableExternal:
                  type: boolean
                enableBfd:
                  type: boolean
                subnets:
                  items:
                    type: string
                  type: array
                extraExternalSubnets:
                  items:
                    type: string
                  type: array
                vpcPeerings:
                  items:
                    type: string
                  type: array
                tcpLoadBalancer:
                  type: string
                tcpSessionLoadBalancer:
                  type: string
                udpLoadBalancer:
                  type: string
                udpSessionLoadBalancer:
                  type: string
                sctpLoadBalancer:
                  type: string
                sctpSessionLoadBalancer:
                  type: string
              type: object
          type: object
      served: true
      storage: true
      subresources:
        status: {}
  names:
    kind: Vpc
    listKind: VpcList
    plural: vpcs
    shortNames:
      - vpc
    singular: vpc
  scope: Cluster`

	ips_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ips.kubeovn.io
spec:
  group: kubeovn.io
  versions:
    - name: v1
      served: true
      storage: true
      additionalPrinterColumns:
      - name: V4IP
        type: string
        jsonPath: .spec.v4IpAddress
      - name: V6IP
        type: string
        jsonPath: .spec.v6IpAddress
      - name: Mac
        type: string
        jsonPath: .spec.macAddress
      - name: Node
        type: string
        jsonPath: .spec.nodeName
      - name: Subnet
        type: string
        jsonPath: .spec.subnet
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                podName:
                  type: string
                namespace:
                  type: string
                subnet:
                  type: string
                attachSubnets:
                  type: array
                  items:
                    type: string
                nodeName:
                  type: string
                ipAddress:
                  type: string
                v4IpAddress:
                  type: string
                v6IpAddress:
                  type: string
                attachIps:
                  type: array
                  items:
                    type: string
                macAddress:
                  type: string
                attachMacs:
                  type: array
                  items:
                    type: string
                containerID:
                  type: string
                podType:
                  type: string
  scope: Cluster
  names:
    plural: ips
    singular: ip
    kind: IP
    shortNames:
      - ip`

	vips_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: vips.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: vips
    singular: vip
    shortNames:
      - vip
    kind: Vip
    listKind: VipList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      additionalPrinterColumns:
      - name: V4IP
        type: string
        jsonPath: .status.v4ip
      - name: V6IP
        type: string
        jsonPath: .status.v6ip
      - name: Mac
        type: string
        jsonPath: .status.mac
      - name: PMac
        type: string
        jsonPath: .spec.parentMac
      - name: Subnet
        type: string
        jsonPath: .spec.subnet
      - jsonPath: .status.ready
        name: Ready
        type: boolean
      - jsonPath: .status.type
        name: Type
        type: string
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                type:
                  type: string
                ready:
                  type: boolean
                v4ip:
                  type: string
                v6ip:
                  type: string
                mac:
                  type: string
                pv4ip:
                  type: string
                pv6ip:
                  type: string
                pmac:
                  type: string
                selector:
                  type: array
                  items:
                    type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                namespace:
                  type: string
                subnet:
                  type: string
                type:
                  type: string
                attachSubnets:
                  type: array
                  items:
                    type: string
                v4ip:
                  type: string
                macAddress:
                  type: string
                v6ip:
                  type: string
                parentV4ip:
                  type: string
                parentMac:
                  type: string
                parentV6ip:
                  type: string
                selector:
                  type: array
                  items:
                    type: string`

	subnets_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: subnets.kubeovn.io
spec:
  group: kubeovn.io
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
      - name: Provider
        type: string
        jsonPath: .spec.provider
      - name: Vpc
        type: string
        jsonPath: .spec.vpc
      - name: Vlan
        type: string
        jsonPath: .spec.vlan
      - name: Protocol
        type: string
        jsonPath: .spec.protocol
      - name: CIDR
        type: string
        jsonPath: .spec.cidrBlock
      - name: Private
        type: boolean
        jsonPath: .spec.private
      - name: NAT
        type: boolean
        jsonPath: .spec.natOutgoing
      - name: Default
        type: boolean
        jsonPath: .spec.default
      - name: GatewayType
        type: string
        jsonPath: .spec.gatewayType
      - name: V4Used
        type: number
        jsonPath: .status.v4usingIPs
      - name: V4Available
        type: number
        jsonPath: .status.v4availableIPs
      - name: V6Used
        type: number
        jsonPath: .status.v6usingIPs
      - name: V6Available
        type: number
        jsonPath: .status.v6availableIPs
      - name: ExcludeIPs
        type: string
        jsonPath: .spec.excludeIps
      - name: U2OInterconnectionIP
        type: string
        jsonPath: .status.u2oInterconnectionIP
      schema:
        openAPIV3Schema:
          type: object
          properties:
            metadata:
              type: object
              properties:
                name:
                  type: string
                  pattern: ^[^0-9]
            status:
              type: object
              properties:
                v4availableIPs:
                  type: number
                v4usingIPs:
                  type: number
                v6availableIPs:
                  type: number
                v6usingIPs:
                  type: number
                activateGateway:
                  type: string
                dhcpV4OptionsUUID:
                  type: string
                dhcpV6OptionsUUID:
                  type: string
                u2oInterconnectionIP:
                  type: string
                u2oInterconnectionMAC:
                  type: string
                u2oInterconnectionVPC:
                  type: string
                mcastQuerierIP:
                  type: string
                mcastQuerierMAC:
                  type: string
                v4usingIPrange:
                  type: string
                v4availableIPrange:
                  type: string
                v6usingIPrange:
                  type: string
                v6availableIPrange:
                  type: string
                natOutgoingPolicyRules:
                  type: array
                  items:
                    type: object
                    properties:
                      ruleID:
                        type: string
                      action:
                        type: string
                        enum:
                          - nat
                          - forward
                      match:
                        type: object
                        properties:
                          srcIPs:
                            type: string
                          dstIPs:
                            type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                vpc:
                  type: string
                default:
                  type: boolean
                protocol:
                  type: string
                  enum:
                    - IPv4
                    - IPv6
                    - Dual
                cidrBlock:
                  type: string
                namespaces:
                  type: array
                  items:
                    type: string
                gateway:
                  type: string
                provider:
                  type: string
                excludeIps:
                  type: array
                  items:
                    type: string
                vips:
                  type: array
                  items:
                    type: string
                gatewayType:
                  type: string
                allowSubnets:
                  type: array
                  items:
                    type: string
                gatewayNode:
                  type: string
                natOutgoing:
                  type: boolean
                externalEgressGateway:
                  type: string
                policyRoutingPriority:
                  type: integer
                  minimum: 1
                  maximum: 32765
                policyRoutingTableID:
                  type: integer
                  minimum: 1
                  maximum: 2147483647
                  not:
                    enum:
                      - 252 # compat
                      - 253 # default
                      - 254 # main
                      - 255 # local
                mtu:
                  type: integer
                  minimum: 68
                  maximum: 65535
                private:
                  type: boolean
                vlan:
                  type: string
                logicalGateway:
                  type: boolean
                disableGatewayCheck:
                  type: boolean
                disableInterConnection:
                  type: boolean
                enableDHCP:
                  type: boolean
                dhcpV4Options:
                  type: string
                dhcpV6Options:
                  type: string
                enableIPv6RA:
                  type: boolean
                ipv6RAConfigs:
                  type: string
                allowEWTraffic:
                  type: boolean
                acls:
                  type: array
                  items:
                    type: object
                    properties:
                      direction:
                        type: string
                        enum:
                          - from-lport
                          - to-lport
                      priority:
                        type: integer
                        minimum: 0
                        maximum: 32767
                      match:
                        type: string
                      action:
                        type: string
                        enum:
                          - allow-related
                          - allow-stateless
                          - allow
                          - drop
                          - reject
                natOutgoingPolicyRules:
                  type: array
                  items:
                    type: object
                    properties:
                      action:
                        type: string
                        enum:
                          - nat
                          - forward
                      match:
                        type: object
                        properties:
                          srcIPs:
                            type: string
                          dstIPs:
                            type: string
                u2oInterconnection:
                  type: boolean
                u2oInterconnectionIP:
                  type: string
                enableLb:
                  type: boolean
                enableEcmp:
                  type: boolean
                enableMulticastSnoop:
                  type: boolean
                routeTable:
                  type: string
                namespaceSelectors:
                  type: array
                  items:
                    type: object
                    properties:
                      matchLabels:
                        type: object
                        additionalProperties:
                          type: string
                      matchExpressions:
                        type: array
                        items:
                          type: object
                          properties:
                            key:
                              type: string
                            operator:
                              type: string
                            values:
                              type: array
                              items:
                                type: string
  scope: Cluster
  names:
    plural: subnets
    singular: subnet
    kind: Subnet
    shortNames:
      - subnet`

	ippools_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ippools.kubeovn.io
spec:
  group: kubeovn.io
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
      - name: Subnet
        type: string
        jsonPath: .spec.subnet
      - name: IPs
        type: string
        jsonPath: .spec.ips
      - name: V4Used
        type: number
        jsonPath: .status.v4UsingIPs
      - name: V4Available
        type: number
        jsonPath: .status.v4AvailableIPs
      - name: V6Used
        type: number
        jsonPath: .status.v6UsingIPs
      - name: V6Available
        type: number
        jsonPath: .status.v6AvailableIPs
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                subnet:
                  type: string
                  x-kubernetes-validations:
                    - rule: "self == oldSelf"
                      message: "This field is immutable."
                namespaces:
                  type: array
                  x-kubernetes-list-type: set
                  items:
                    type: string
                ips:
                  type: array
                  minItems: 1
                  x-kubernetes-list-type: set
                  items:
                    type: string
                    anyOf:
                      - format: ipv4
                      - format: ipv6
                      - format: cidr
                      - pattern: ^(?:(?:[01]?\d{1,2}|2[0-4]\d|25[0-5])\.){3}(?:[01]?\d{1,2}|2[0-4]\d|25[0-5])\.\.(?:(?:[01]?\d{1,2}|2[0-4]\d|25[0-5])\.){3}(?:[01]?\d{1,2}|2[0-4]\d|25[0-5])$
                      - pattern: ^((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|:)))\.\.((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|:)))$
              required:
                - subnet
                - ips
            status:
              type: object
              properties:
                v4AvailableIPs:
                  type: number
                v4UsingIPs:
                  type: number
                v6AvailableIPs:
                  type: number
                v6UsingIPs:
                  type: number
                v4AvailableIPRange:
                  type: string
                v4UsingIPRange:
                  type: string
                v6AvailableIPRange:
                  type: string
                v6UsingIPRange:
                  type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
  scope: Cluster
  names:
    plural: ippools
    singular: ippool
    kind: IPPool
    shortNames:
      - ippool`

	vlans_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: vlans.kubeovn.io
spec:
  group: kubeovn.io
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                id:
                  type: integer
                  minimum: 0
                  maximum: 4095
                provider:
                  type: string
                vlanId:
                  type: integer
                  description: Deprecated in favor of id
                providerInterfaceName:
                  type: string
                  description: Deprecated in favor of provider
              required:
                - provider
            status:
              type: object
              properties:
                subnets:
                  type: array
                  items:
                    type: string
      additionalPrinterColumns:
      - name: ID
        type: string
        jsonPath: .spec.id
      - name: Provider
        type: string
        jsonPath: .spec.provider
  scope: Cluster
  names:
    plural: vlans
    singular: vlan
    kind: Vlan
    shortNames:
      - vlan`

	provider_networks_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: provider-networks.kubeovn.io
spec:
  group: kubeovn.io
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      schema:
        openAPIV3Schema:
          type: object
          properties:
            metadata:
              type: object
              properties:
                name:
                  type: string
                  maxLength: 12
                  not:
                    enum:
                      - int
            spec:
              type: object
              properties:
                defaultInterface:
                  type: string
                  maxLength: 15
                  pattern: '^[^/\s]+$'
                customInterfaces:
                  type: array
                  items:
                    type: object
                    properties:
                      interface:
                        type: string
                        maxLength: 15
                        pattern: '^[^/\s]+$'
                      nodes:
                        type: array
                        items:
                          type: string
                exchangeLinkName:
                  type: boolean
                excludeNodes:
                  type: array
                  items:
                    type: string
              required:
                - defaultInterface
            status:
              type: object
              properties:
                ready:
                  type: boolean
                readyNodes:
                  type: array
                  items:
                    type: string
                notReadyNodes:
                  type: array
                  items:
                    type: string
                vlans:
                  type: array
                  items:
                    type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      node:
                        type: string
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
      additionalPrinterColumns:
      - name: DefaultInterface
        type: string
        jsonPath: .spec.defaultInterface
      - name: Ready
        type: boolean
        jsonPath: .status.ready
  scope: Cluster
  names:
    plural: provider-networks
    singular: provider-network
    kind: ProviderNetwork
    listKind: ProviderNetworkList`

	security_groups_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: security-groups.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: security-groups
    singular: security-group
    shortNames:
      - sg
    kind: SecurityGroup
    listKind: SecurityGroupList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                ingressRules:
                  type: array
                  items:
                    type: object
                    properties:
                      ipVersion:
                        type: string
                      protocol:
                        type: string
                      priority:
                        type: integer
                      remoteType:
                        type: string
                      remoteAddress:
                        type: string
                      remoteSecurityGroup:
                        type: string
                      portRangeMin:
                        type: integer
                      portRangeMax:
                        type: integer
                      policy:
                        type: string
                egressRules:
                  type: array
                  items:
                    type: object
                    properties:
                      ipVersion:
                        type: string
                      protocol:
                        type: string
                      priority:
                        type: integer
                      remoteType:
                        type: string
                      remoteAddress:
                        type: string
                      remoteSecurityGroup:
                        type: string
                      portRangeMin:
                        type: integer
                      portRangeMax:
                        type: integer
                      policy:
                        type: string
                allowSameGroupTraffic:
                  type: boolean
            status:
              type: object
              properties:
                portGroup:
                  type: string
                allowSameGroupTraffic:
                  type: boolean
                ingressMd5:
                  type: string
                egressMd5:
                  type: string
                ingressLastSyncSuccess:
                  type: boolean
                egressLastSyncSuccess:
                  type: boolean
      subresources:
        status: {}
  conversion:
    strategy: None`

	qos_policies_crd = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: qos-policies.kubeovn.io
spec:
  group: kubeovn.io
  names:
    plural: qos-policies
    singular: qos-policy
    shortNames:
      - qos
    kind: QoSPolicy
    listKind: QoSPolicyList
  scope: Cluster
  versions:
    - name: v1
      served: true
      storage: true
      subresources:
        status: {}
      additionalPrinterColumns:
      - jsonPath: .spec.shared
        name: Shared
        type: string
      - jsonPath: .spec.bindingType
        name: BindingType
        type: string
      schema:
        openAPIV3Schema:
          type: object
          properties:
            status:
              type: object
              properties:
                shared:
                  type: boolean
                bindingType:
                  type: string
                bandwidthLimitRules:
                  type: array
                  items:
                    type: object
                    properties:
                      name:
                        type: string
                      interface:
                        type: string
                      rateMax:
                        type: string
                      burstMax:
                        type: string
                      priority:
                        type: integer
                      direction:
                        type: string
                      matchType:
                        type: string
                      matchValue:
                        type: string
                conditions:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                      status:
                        type: string
                      reason:
                        type: string
                      message:
                        type: string
                      lastUpdateTime:
                        type: string
                      lastTransitionTime:
                        type: string
            spec:
              type: object
              properties:
                shared:
                  type: boolean
                bindingType:
                  type: string
                bandwidthLimitRules:
                  type: array
                  items:
                    type: object
                    properties:
                      name:
                        type: string
                      interface:
                        type: string
                      rateMax:
                        type: string
                      burstMax:
                        type: string
                      priority:
                        type: integer
                      direction:
                        type: string
                      matchType:
                        type: string
                      matchValue:
                        type: string
                    required:
                      - name
                  x-kubernetes-list-map-keys:
                    - name
                  x-kubernetes-list-type: map`

	CRDList = []string{vpc_dnses_crd, switch_lb_rules_crd, vpc_nat_gateways_crd, iptables_eips_crd, iptables_fip_rules_crd, iptables_dnat_rules_crd, iptables_snat_rules_crd, ovn_eips_crd, ovn_fips_crd, ovn_snat_rules_crd, ovn_dnat_rules_crd, vpcs_crd, ips_crd, vips_crd, subnets_crd, ippools_crd, vlans_crd, provider_networks_crd, security_groups_crd, qos_policies_crd}
)
