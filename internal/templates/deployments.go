package templates

var (
	ovn_central_deployment = `kind: Deployment
apiVersion: apps/v1
metadata:
  name: ovn-central
  namespace: {{ .Values.namespace }}
  annotations:
    kubernetes.io/description: |
      OVN components: northd, nb and sb.
spec:
  replicas: {{ include "kubeovn.nodeCount" . }}
  strategy:
    rollingUpdate:
      maxSurge: 0
      maxUnavailable: 1
    type: RollingUpdate
  selector:
    matchLabels:
      app: ovn-central
  template:
    metadata:
      labels:
        app: ovn-central
        component: network
        type: infra
    spec:
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - effect: NoExecute
          operator: Exists
        - key: CriticalAddonsOnly
          operator: Exists
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchLabels:
                  app: ovn-central
              topologyKey: kubernetes.io/hostname
      priorityClassName: system-cluster-critical
      serviceAccountName: ovn-ovs
      hostNetwork: true
      initContainers:
        - name: hostpath-init
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          command:
            - sh
            - -c
            - "chown -R nobody: /var/run/ovn /etc/ovn /var/log/ovn"
          securityContext:
            allowPrivilegeEscalation: true
            capabfserilities:
              drop:
                - ALL
            privileged: true
            runAsUser: 0
          volumeMounts:
            - mountPath: /var/run/ovn
              name: host-run-ovn
            - mountPath: /etc/ovn
              name: host-config-ovn
            - mountPath: /var/log/ovn
              name: host-log-ovn
      containers:
        - name: ovn-central
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          command: 
          - bash
          - /kube-ovn/start-db.sh
          securityContext:
            runAsUser: {{ include "kubeovn.runAsUser" . }}
            privileged: false
            capabilities:
              add:
                - NET_BIND_SERVICE
                - SYS_NICE
          env:
            - name: ENABLE_SSL
              value: "{{ .Values.networking.ENABLE_SSL }}"
            - name: NODE_IPS
              value: "{{ .Values.MASTER_NODES | default (include "kubeovn.nodeIPs" .) }}"
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_IPS
              valueFrom:
                fieldRef:
                  fieldPath: status.podIPs
            - name: ENABLE_BIND_LOCAL_IP
              value: "{{- .Values.components.ENABLE_BIND_LOCAL_IP }}"
            - name: PROBE_INTERVAL
              value: "{{ .Values.networking.PROBE_INTERVAL }}"
            - name: OVN_NORTHD_PROBE_INTERVAL
              value: "{{ .Values.networking.OVN_NORTHD_PROBE_INTERVAL}}"
            - name: OVN_LEADER_PROBE_INTERVAL
              value: "{{ .Values.networking.OVN_LEADER_PROBE_INTERVAL }}"
            - name: OVN_NORTHD_N_THREADS
              value: "{{ .Values.networking.OVN_NORTHD_N_THREADS }}"
            - name: ENABLE_COMPACT
              value: "{{ .Values.networking.ENABLE_COMPACT }}"
            - name: OVN_VERSION_COMPATIBILITY
              value: '{{ include "kubeovn.ovn.versionCompatibility" . }}'
          resources:
            requests:
              cpu: {{ index .Values "ovn-central" "requests" "cpu" }}
              memory: {{ index .Values "ovn-central" "requests" "memory" }}
            limits:
              cpu: {{ index .Values "ovn-central" "limits" "cpu" }}
              memory: {{ index .Values "ovn-central" "limits" "memory" }}
          volumeMounts:
            - mountPath: /var/run/ovn
              name: host-run-ovn
            - mountPath: /etc/ovn
              name: host-config-ovn
            - mountPath: /var/log/ovn
              name: host-log-ovn
            - mountPath: /etc/localtime
              name: localtime
              readOnly: true
            - mountPath: /var/run/tls
              name: kube-ovn-tls
          readinessProbe:
            exec:
              command:
                - bash
                - /kube-ovn/ovn-healthcheck.sh
            periodSeconds: 15
            timeoutSeconds: 45
          livenessProbe:
            exec:
              command:
                - bash
                - /kube-ovn/ovn-healthcheck.sh
            initialDelaySeconds: 30
            periodSeconds: 15
            failureThreshold: 5
            timeoutSeconds: 45
      nodeSelector:
        kubernetes.io/os: "linux"
        {{- with splitList "=" .Values.MASTER_NODES_LABEL }}
        {{ index . 0 }}: "{{ if eq (len .) 2 }}{{ index . 1 }}{{ end }}"
        {{- end }}
      volumes:
        - name: host-run-ovn
          hostPath:
            path: /run/ovn
        - name: host-config-ovn
          hostPath:
            path: {{ .Values.OVN_DIR }}
        - name: host-log-ovn
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/ovn
        - name: localtime
          hostPath:
            path: /etc/localtime
        - name: kube-ovn-tls
          secret:
            optional: true
            secretName: kube-ovn-tls`

	kube_ovn_controller_deployment = `kind: Deployment
apiVersion: apps/v1
metadata:
  name: kube-ovn-controller
  namespace: {{ .Values.namespace }}
  annotations:
    kubernetes.io/description: |
      kube-ovn controller
spec:
  replicas: {{ include "kubeovn.nodeCount" . }}
  selector:
    matchLabels:
      app: kube-ovn-controller
  strategy:
    rollingUpdate:
      maxSurge: 0%
      maxUnavailable: 100%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: kube-ovn-controller
        component: network
        type: infra
    spec:
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - key: CriticalAddonsOnly
          operator: Exists
      affinity:
        nodeAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - preference:
                matchExpressions:
                  - key: "ovn.kubernetes.io/ic-gw"
                    operator: NotIn
                    values:
                      - "true"
              weight: 100
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchLabels:
                  app: kube-ovn-controller
              topologyKey: kubernetes.io/hostname
      priorityClassName: system-cluster-critical
      serviceAccountName: ovn
      hostNetwork: true
      initContainers:
        - name: hostpath-init
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
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
        - name: kube-ovn-controller
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          args:
          - /kube-ovn/start-controller.sh
          - --default-ls={{ .Values.networking.DEFAULT_SUBNET }}
          - --default-cidr=
          {{- if eq .Values.networking.NET_STACK "dual_stack" -}}
          {{ .Values.dual_stack.POD_CIDR }}
          {{- else if eq .Values.networking.NET_STACK "ipv4" -}}
          {{ .Values.ipv4.POD_CIDR }}
          {{- else if eq .Values.networking.NET_STACK "ipv6" -}}
          {{ .Values.ipv6.POD_CIDR }}
          {{- end }}
          - --default-gateway=
          {{- if eq .Values.networking.NET_STACK "dual_stack" -}}
          {{ .Values.dual_stack.POD_GATEWAY }}
          {{- else if eq .Values.networking.NET_STACK "ipv4" -}}
          {{ .Values.ipv4.POD_GATEWAY }}
          {{- else if eq .Values.networking.NET_STACK "ipv6" -}}
          {{ .Values.ipv6.POD_GATEWAY }}
          {{- end }}
          - --default-gateway-check={{- .Values.components.CHECK_GATEWAY }}
          - --default-logical-gateway={{- .Values.components.LOGICAL_GATEWAY }}
          - --default-u2o-interconnection={{- .Values.components.U2O_INTERCONNECTION }}
          - --default-exclude-ips={{- .Values.networking.EXCLUDE_IPS }}
          - --cluster-router={{ .Values.networking.DEFAULT_VPC }}
          - --node-switch={{ .Values.networking.NODE_SUBNET }}
          - --node-switch-cidr=
          {{- if eq .Values.networking.NET_STACK "dual_stack" -}}
          {{ .Values.dual_stack.JOIN_CIDR }}
          {{- else if eq .Values.networking.NET_STACK "ipv4" -}}
          {{ .Values.ipv4.JOIN_CIDR }}
          {{- else if eq .Values.networking.NET_STACK "ipv6" -}}
          {{ .Values.ipv6.JOIN_CIDR }}
          {{- end }}
          - --service-cluster-ip-range=
          {{- if eq .Values.networking.NET_STACK "dual_stack" -}}
          {{ .Values.dual_stack.SVC_CIDR }}
          {{- else if eq .Values.networking.NET_STACK "ipv4" -}}
          {{ .Values.ipv4.SVC_CIDR }}
          {{- else if eq .Values.networking.NET_STACK "ipv6" -}}
          {{ .Values.ipv6.SVC_CIDR }}
          {{- end }}
          - --network-type={{- .Values.networking.NETWORK_TYPE }}
          - --default-provider-name={{ .Values.networking.vlan.PROVIDER_NAME }}
          - --default-interface-name={{- .Values.networking.vlan.VLAN_INTERFACE_NAME }}
          - --default-exchange-link-name={{- .Values.networking.EXCHANGE_LINK_NAME }}
          - --default-vlan-name={{- .Values.networking.vlan.VLAN_NAME }}
          - --default-vlan-id={{- .Values.networking.vlan.VLAN_ID }}
          - --ls-dnat-mod-dl-dst={{- .Values.components.LS_DNAT_MOD_DL_DST }}
          - --ls-ct-skip-dst-lport-ips={{- .Values.components.LS_CT_SKIP_DST_LPORT_IPS }}
          - --pod-nic-type={{- .Values.networking.POD_NIC_TYPE }}
          - --enable-lb={{- .Values.components.ENABLE_LB }}
          - --enable-np={{- .Values.components.ENABLE_NP }}
          - --enable-eip-snat={{- .Values.networking.ENABLE_EIP_SNAT }}
          - --enable-external-vpc={{- .Values.components.ENABLE_EXTERNAL_VPC }}
          - --enable-ecmp={{- .Values.networking.ENABLE_ECMP }}
          - --logtostderr=false
          - --alsologtostderr=true
          - --gc-interval={{- .Values.performance.GC_INTERVAL }}
          - --inspect-interval={{- .Values.performance.INSPECT_INTERVAL }}
          - --log_file=/var/log/kube-ovn/kube-ovn-controller.log
          - --log_file_max_size=200
          - --enable-lb-svc={{- .Values.components.ENABLE_LB_SVC }}
          - --keep-vm-ip={{- .Values.components.ENABLE_KEEP_VM_IP }}
          - --enable-metrics={{- .Values.networking.ENABLE_METRICS }}
          - --node-local-dns-ip={{- .Values.networking.NODE_LOCAL_DNS_IP }}
          - --secure-serving={{- .Values.components.SECURE_SERVING }}
          - --enable-ovn-ipsec={{- .Values.components.ENABLE_OVN_IPSEC }}
          - --enable-anp={{- .Values.components.ENABLE_ANP }}
          - --ovsdb-con-timeout={{- .Values.components.OVSDB_CON_TIMEOUT }}
          - --ovsdb-inactivity-timeout={{- .Values.components.OVSDB_INACTIVITY_TIMEOUT }}
          securityContext:
            runAsUser: {{ include "kubeovn.runAsUser" . }}
            privileged: false
            capabilities:
              add:
                - NET_BIND_SERVICE
                - NET_RAW
          env:
            - name: ENABLE_SSL
              value: "{{ .Values.networking.ENABLE_SSL }}"
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: KUBE_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: KUBE_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: OVN_DB_IPS
              value: "{{ .Values.MASTER_NODES | default (include "kubeovn.nodeIPs" .) }}"
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: POD_IPS
              valueFrom:
                fieldRef:
                  fieldPath: status.podIPs
            - name: ENABLE_BIND_LOCAL_IP
              value: "{{- .Values.components.ENABLE_BIND_LOCAL_IP }}"
          volumeMounts:
            - mountPath: /etc/localtime
              name: localtime
              readOnly: true
            - mountPath: /var/log/kube-ovn
              name: kube-ovn-log
            # ovn-ic log directory
            - mountPath: /var/log/ovn
              name: ovn-log
            - mountPath: /var/run/tls
              name: kube-ovn-tls
          readinessProbe:
            httpGet:
              port: 10660
              path: /readyz
              scheme: '{{ ternary "HTTPS" "HTTP" .Values.components.SECURE_SERVING }}'
            periodSeconds: 3
            timeoutSeconds: 5
          livenessProbe:
            httpGet:
              port: 10660
              path: /livez
              scheme: '{{ ternary "HTTPS" "HTTP" .Values.components.SECURE_SERVING }}'
            initialDelaySeconds: 300
            periodSeconds: 7
            failureThreshold: 5
            timeoutSeconds: 5
          resources:
            requests:
              cpu: {{ index .Values "kube-ovn-controller" "requests" "cpu" }}
              memory: {{ index .Values "kube-ovn-controller" "requests" "memory" }}
            limits:
              cpu: {{ index .Values "kube-ovn-controller" "limits" "cpu" }}
              memory: {{ index .Values "kube-ovn-controller" "limits" "memory" }}
      nodeSelector:
        kubernetes.io/os: "linux"
      volumes:
        - name: localtime
          hostPath:
            path: /etc/localtime
        - name: kube-ovn-log
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/kube-ovn
        - name: ovn-log
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/ovn
        - name: kube-ovn-tls
          secret:
            optional: true
            secretName: kube-ovn-tls`

	ovn_ic_controller_deployment = `{{- if .Values.components.ENABLE_IC }}
kind: Deployment
apiVersion: apps/v1
metadata:
  name: ovn-ic-controller
  namespace: kube-system
  annotations:
    kubernetes.io/description: |
      OVN IC Client
spec:
  replicas: 1
  strategy:
    rollingUpdate:
      maxSurge: 0
      maxUnavailable: 1
    type: RollingUpdate
  selector:
    matchLabels:
      app: ovn-ic-controller
  template:
    metadata:
      labels:
        app: ovn-ic-controller
        component: network
        type: infra
    spec:
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - effect: NoExecute
          operator: Exists
        - key: CriticalAddonsOnly
          operator: Exists
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchLabels:
                  app: ovn-ic-controller
              topologyKey: kubernetes.io/hostname
      priorityClassName: system-cluster-critical
      serviceAccountName: ovn
      hostNetwork: true
      initContainers:
        - name: hostpath-init
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          command:
            - sh
            - -c
            - "chown -R nobody: /var/run/ovn /var/log/ovn /var/log/kube-ovn"
          securityContext:
            allowPrivilegeEscalation: true
            capabilities:
              drop:
                - ALL
            privileged: true
            runAsUser: 0
          volumeMounts:
            - mountPath: /var/run/ovn
              name: host-run-ovn
            - mountPath: /var/log/ovn
              name: host-log-ovn
            - name: kube-ovn-log
              mountPath: /var/log/kube-ovn
      containers:
        - name: ovn-ic-controller
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          command: ["/kube-ovn/start-ic-controller.sh"]
          args:
          - --log_file=/var/log/kube-ovn/kube-ovn-ic-controller.log
          - --log_file_max_size=200
          - --logtostderr=false
          - --alsologtostderr=true
          securityContext:
            runAsUser: {{ include "kubeovn.runAsUser" . }}
            privileged: false
            capabilities:
              add:
                - NET_BIND_SERVICE
                - SYS_NICE
          env:
            - name: ENABLE_SSL
              value: "{{ .Values.networking.ENABLE_SSL }}"
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: OVN_DB_IPS
              value: "{{ .Values.MASTER_NODES | default (include "kubeovn.nodeIPs" .) }}"
          resources:
            requests:
              cpu: 300m
              memory: 200Mi
            limits:
              cpu: 3
              memory: 1Gi
          volumeMounts:
            - mountPath: /var/run/ovn
              name: host-run-ovn
            - mountPath: /var/log/ovn
              name: host-log-ovn
            - mountPath: /etc/localtime
              name: localtime
            - mountPath: /var/run/tls
              name: kube-ovn-tls
            - mountPath: /var/log/kube-ovn
              name: kube-ovn-log
      nodeSelector:
        kubernetes.io/os: "linux"
        kube-ovn/role: "master"
      volumes:
        - name: host-run-ovn
          hostPath:
            path: /run/ovn
        - name: host-log-ovn
          hostPath:
            path: /var/log/ovn
        - name: localtime
          hostPath:
            path: /etc/localtime
        - name: kube-ovn-log
          hostPath:
            path: /var/log/kube-ovn
        - name: kube-ovn-tls
          secret:
            optional: true
            secretName: kube-ovn-tls
{{- end }}`

	kube_ovn_monitor_deployment = `kind: Deployment
apiVersion: apps/v1
metadata:
  name: kube-ovn-monitor
  namespace: {{ .Values.namespace }}
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
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
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
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          command: ["/kube-ovn/start-ovn-monitor.sh"]
          args:
          - --secure-serving={{- .Values.components.SECURE_SERVING }}
          - --log_file=/var/log/kube-ovn/kube-ovn-monitor.log
          - --logtostderr=false
          - --alsologtostderr=true
          - --log_file_max_size=200
          - --enable-metrics={{- .Values.networking.ENABLE_METRICS }}
          securityContext:
            runAsUser: {{ include "kubeovn.runAsUser" . }}
            privileged: false
            capabilities:
              add:
                - NET_BIND_SERVICE
          env:
            - name: ENABLE_SSL
              value: "{{ .Values.networking.ENABLE_SSL }}"
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
              value: "{{- .Values.components.ENABLE_BIND_LOCAL_IP }}"
          resources:
            requests:
              cpu: {{ index .Values "kube-ovn-monitor" "requests" "cpu" }}
              memory: {{ index .Values "kube-ovn-monitor" "requests" "memory" }}
            limits:
              cpu: {{ index .Values "kube-ovn-monitor" "limits" "cpu" }}
              memory: {{ index .Values "kube-ovn-monitor" "limits" "memory" }}
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
              scheme: '{{ ternary "HTTPS" "HTTP" .Values.components.SECURE_SERVING }}'
            timeoutSeconds: 5
          readinessProbe:
            failureThreshold: 3
            initialDelaySeconds: 30
            periodSeconds: 7
            successThreshold: 1
            httpGet:
              port: 10661
              path: /readyz
              scheme: '{{ ternary "HTTPS" "HTTP" .Values.components.SECURE_SERVING }}'
            timeoutSeconds: 5
      nodeSelector:
        kubernetes.io/os: "linux"
        {{- with splitList "=" .Values.MASTER_NODES_LABEL }}
        {{ index . 0 }}: "{{ if eq (len .) 2 }}{{ index . 1 }}{{ end }}"
        {{- end }}
      volumes:
        - name: host-run-ovn
          hostPath:
            path: /run/ovn
        - name: host-config-ovn
          hostPath:
            path: {{ .Values.OVN_DIR }}
        - name: host-log-ovn
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/ovn
        - name: localtime
          hostPath:
            path: /etc/localtime
        - name: kube-ovn-tls
          secret:
            optional: true
            secretName: kube-ovn-tls
        - name: kube-ovn-log
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/kube-ovn`

	DeploymentList = []string{ovn_central_deployment, kube_ovn_controller_deployment, ovn_ic_controller_deployment, kube_ovn_monitor_deployment}
)
