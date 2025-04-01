package templates

var (
	ovs_ovn_dpdk_daemonset = `{{- if .Values.HYBRID_DPDK }}
kind: DaemonSet
apiVersion: apps/v1
metadata:
  name: ovs-ovn-dpdk
  namespace: {{ .Values.namespace }}
  annotations:
    kubernetes.io/description: |
      This daemon set launches the openvswitch daemon.
spec:
  selector:
    matchLabels:
      app: ovs-dpdk
  updateStrategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: ovs-dpdk
        component: network
        type: infra
    spec:
      tolerations:
      - operator: Exists
      priorityClassName: system-node-critical
      serviceAccountName: ovn-ovs
      hostNetwork: true
      hostPID: true
      containers:
        - name: openvswitch
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}-dpdk
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          command: ["/kube-ovn/start-ovs-dpdk-v2.sh"]
          securityContext:
            runAsUser: 0
            privileged: true
          env:
            - name: ENABLE_SSL
              value: "{{ .Values.networking.ENABLE_SSL }}"
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: HW_OFFLOAD
              value: "{{- .Values.components.HW_OFFLOAD }}"
            - name: TUNNEL_TYPE
              value: "{{- .Values.networking.TUNNEL_TYPE }}"
            - name: DPDK_TUNNEL_IFACE
              value: "{{- .Values.networking.DPDK_TUNNEL_IFACE }}"
            - name: KUBE_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: OVN_DB_IPS
              value: "{{ .Values.MASTER_NODES | default (include "kubeovn.nodeIPs" .) }}"
            - name: OVN_REMOTE_PROBE_INTERVAL
              value: "{{ .Values.networking.OVN_REMOTE_PROBE_INTERVAL }}"
            - name: OVN_REMOTE_OPENFLOW_INTERVAL
              value: "{{ .Values.networking.OVN_REMOTE_OPENFLOW_INTERVAL }}"
          volumeMounts:
            - mountPath: /opt/ovs-config
              name: host-config-ovs
            - name: shareddir
              mountPath: {{ .Values.kubelet_conf.KUBELET_DIR }}/pods
            - name: hugepage
              mountPath: /dev/hugepages
            - mountPath: /lib/modules
              name: host-modules
              readOnly: true
            - mountPath: /var/run/openvswitch
              name: host-run-ovs
              mountPropagation: HostToContainer
            - mountPath: /var/run/ovn
              name: host-run-ovn
            - mountPath: /sys
              name: host-sys
            - mountPath: /etc/openvswitch
              name: host-config-openvswitch
            - mountPath: /etc/ovn
              name: host-config-ovn
            - mountPath: /var/log/openvswitch
              name: host-log-ovs
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
                - /kube-ovn/ovs-healthcheck.sh
            periodSeconds: 5
            timeoutSeconds: 45
          livenessProbe:
            exec:
              command:
                - bash
                - /kube-ovn/ovs-healthcheck.sh
            initialDelaySeconds: 60
            periodSeconds: 5
            failureThreshold: 5
            timeoutSeconds: 45
          resources:
            requests:
              cpu: {{ index .Values "ovs-ovn" "requests" "cpu" }}
              memory: {{ index .Values "ovs-ovn" "requests" "memory" }}
            limits:
              cpu: {{ index .Values "ovs-ovn" "limits" "cpu" }}
              {{.Values.HUGEPAGE_SIZE_TYPE}}: {{.Values.HUGEPAGES}}
              memory: {{ index .Values "ovs-ovn" "limits" "memory" }}
      nodeSelector:
        kubernetes.io/os: "linux"
        ovn.kubernetes.io/ovs_dp_type: "userspace"
      volumes:
        - name: host-config-ovs
          hostPath:
            path: /opt/ovs-config
            type: DirectoryOrCreate
        - name: shareddir
          hostPath:
            path: {{ .Values.kubelet_conf.KUBELET_DIR }}/pods
            type: ''
        - name: hugepage
          emptyDir:
            medium: HugePages
        - name: host-modules
          hostPath:
            path: /lib/modules
        - name: host-run-ovs
          hostPath:
            path: /run/openvswitch
        - name: host-run-ovn
          hostPath:
            path: /run/ovn
        - name: host-sys
          hostPath:
            path: /sys
        - name: host-config-openvswitch
          hostPath:
            path: {{ .Values.OPENVSWITCH_DIR }}
        - name: host-config-ovn
          hostPath:
            path: {{ .Values.OVN_DIR }}
        - name: host-log-ovs
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/openvswitch
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
{{- end }}`

	kube_ovn_cni_daemonset = `kind: DaemonSet
apiVersion: apps/v1
metadata:
  name: kube-ovn-cni
  namespace: {{ .Values.namespace }}
  annotations:
    kubernetes.io/description: |
      This daemon set launches the kube-ovn cni daemon.
spec:
  selector:
    matchLabels:
      app: kube-ovn-cni
  template:
    metadata:
      labels:
        app: kube-ovn-cni
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
      priorityClassName: system-node-critical
      serviceAccountName: kube-ovn-cni
      hostNetwork: true
      hostPID: true
      initContainers:
      - name: hostpath-init
        image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
        imagePullPolicy: {{ .Values.image.pullPolicy }}
        command:
          - sh
          - -xec
          - {{ if not .Values.DISABLE_MODULES_MANAGEMENT -}}
            iptables -V
            {{- else -}}
            echo "nothing to do"
            {{- end }}
        securityContext:
          allowPrivilegeEscalation: true
          capabilities:
            drop:
              - ALL
          privileged: true
          runAsUser: 0
          runAsGroup: 0
        volumeMounts:
          - name: usr-local-sbin
            mountPath: /usr/local/sbin
          - mountPath: /run/xtables.lock
            name: xtables-lock
            readOnly: false
          - mountPath: /var/run/netns
            name: host-ns
            readOnly: false
          - name: kube-ovn-log
            mountPath: /var/log/kube-ovn
      - name: install-cni
        image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
        imagePullPolicy: {{ .Values.image.pullPolicy }}
        command:
          - /kube-ovn/install-cni.sh
          - --cni-conf-dir={{ .Values.cni_conf.CNI_CONF_DIR }}
          - --cni-conf-file={{ .Values.cni_conf.CNI_CONF_FILE }}
          - --cni-conf-name={{- .Values.cni_conf.CNI_CONFIG_PRIORITY -}}-kube-ovn.conflist
        securityContext:
          runAsUser: 0
          privileged: true
        volumeMounts:
          - mountPath: /opt/cni/bin
            name: cni-bin
          - mountPath: /etc/cni/net.d
            name: cni-conf
          {{- if .Values.cni_conf.MOUNT_LOCAL_BIN_DIR }}
          - mountPath: /usr/local/bin
            name: local-bin
          {{- end }}
      containers:
      - name: cni-server
        image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
        imagePullPolicy: {{ .Values.image.pullPolicy }}
        command:
          - bash
          - /kube-ovn/start-cniserver.sh
        args:
          - --enable-mirror={{- .Values.debug.ENABLE_MIRROR }}
          - --mirror-iface={{- .Values.debug.MIRROR_IFACE }}
          - --node-switch={{ .Values.networking.NODE_SUBNET }}
          - --encap-checksum=true
          - --service-cluster-ip-range=
          {{- if eq .Values.networking.NET_STACK "dual_stack" -}}
          {{ .Values.dual_stack.SVC_CIDR }}
          {{- else if eq .Values.networking.NET_STACK "ipv4" -}}
          {{ .Values.ipv4.SVC_CIDR }}
          {{- else if eq .Values.networking.NET_STACK "ipv6" -}}
          {{ .Values.ipv6.SVC_CIDR }}
          {{- end }}
          {{- if eq .Values.networking.NETWORK_TYPE "vlan" }}
          - --iface=
          {{- else}}
          - --iface={{- .Values.networking.IFACE }}
          {{- end }}
          - --dpdk-tunnel-iface={{- .Values.networking.DPDK_TUNNEL_IFACE }}
          - --network-type={{- .Values.networking.TUNNEL_TYPE }}
          - --default-interface-name={{- .Values.networking.vlan.VLAN_INTERFACE_NAME }}
          - --logtostderr=false
          - --alsologtostderr=true
          - --log_file=/var/log/kube-ovn/kube-ovn-cni.log
          - --log_file_max_size=200
          - --enable-metrics={{- .Values.networking.ENABLE_METRICS }}
          - --kubelet-dir={{ .Values.kubelet_conf.KUBELET_DIR }}
          - --enable-tproxy={{ .Values.components.ENABLE_TPROXY }}
          - --ovs-vsctl-concurrency={{ .Values.performance.OVS_VSCTL_CONCURRENCY }}
          - --secure-serving={{- .Values.components.SECURE_SERVING }}
          - --enable-ovn-ipsec={{- .Values.components.ENABLE_OVN_IPSEC }}
          - --set-vxlan-tx-off={{- .Values.components.SET_VXLAN_TX_OFF }}
        securityContext:
          runAsUser: 0
          privileged: false
          capabilities:
            add:
              - NET_ADMIN
              - NET_BIND_SERVICE
              - NET_RAW
              - SYS_ADMIN
              - SYS_PTRACE
              {{- if not .Values.DISABLE_MODULES_MANAGEMENT }}
              - SYS_MODULE
              {{- end }}
              - SYS_NICE
        env:
          - name: ENABLE_SSL
            value: "{{ .Values.networking.ENABLE_SSL }}"
          - name: POD_IP
            valueFrom:
              fieldRef:
                fieldPath: status.podIP
          - name: KUBE_NODE_NAME
            valueFrom:
              fieldRef:
                fieldPath: spec.nodeName
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
          - name: DBUS_SYSTEM_BUS_ADDRESS
            value: "unix:path=/host/var/run/dbus/system_bus_socket"
        volumeMounts:
          - name: usr-local-sbin
            mountPath: /usr/local/sbin
          - name: host-modules
            mountPath: /lib/modules
            readOnly: true
          - mountPath: /run/xtables.lock
            name: xtables-lock
            readOnly: false
          - name: shared-dir
            mountPath: {{ .Values.kubelet_conf.KUBELET_DIR }}/pods
          - mountPath: /etc/openvswitch
            name: systemid
            readOnly: true
          - mountPath: /run/openvswitch
            name: host-run-ovs
            mountPropagation: HostToContainer
          - mountPath: /run/ovn
            name: host-run-ovn
          - mountPath: /host/var/run/dbus
            name: host-dbus
            mountPropagation: HostToContainer
          - mountPath: /var/run/netns
            name: host-ns
            mountPropagation: HostToContainer
          - mountPath: /var/log/kube-ovn
            name: kube-ovn-log
          - mountPath: /var/log/openvswitch
            name: host-log-ovs
          - mountPath: /var/log/ovn
            name: host-log-ovn
          - mountPath: /etc/localtime
            name: localtime
            readOnly: true
        {{- if .Values.components.ENABLE_OVN_IPSEC }}
          - mountPath: /etc/ovs_ipsec_keys
            name: ovs-ipsec-keys
        {{- end }}
        readinessProbe:
          failureThreshold: 3
          periodSeconds: 7
          successThreshold: 1
          httpGet:
            port: 10665
            path: /readyz
            scheme: '{{ ternary "HTTPS" "HTTP" .Values.components.SECURE_SERVING }}'
          timeoutSeconds: 5
        livenessProbe:
          failureThreshold: 3
          initialDelaySeconds: 30
          periodSeconds: 7
          successThreshold: 1
          httpGet:
            port: 10665
            path: /livez
            scheme: '{{ ternary "HTTPS" "HTTP" .Values.components.SECURE_SERVING }}'
          timeoutSeconds: 5
        resources:
          requests:
            cpu: {{ index .Values "kube-ovn-cni" "requests" "cpu" }}
            memory: {{ index .Values "kube-ovn-cni" "requests" "memory" }}
          limits:
            cpu: {{ index .Values "kube-ovn-cni" "limits" "cpu" }}
            memory: {{ index .Values "kube-ovn-cni" "limits" "memory" }}
      nodeSelector:
        kubernetes.io/os: "linux"
      volumes:
        - name: usr-local-sbin
          emptyDir: {}
        - name: host-modules
          hostPath:
            path: /lib/modules
        - name: xtables-lock
          hostPath:
            path: /run/xtables.lock
            type: FileOrCreate
        - name: shared-dir
          hostPath:
            path: {{ .Values.kubelet_conf.KUBELET_DIR }}/pods
        - name: systemid
          hostPath:
            path: {{ .Values.OPENVSWITCH_DIR }}
        - name: host-run-ovs
          hostPath:
            path: /run/openvswitch
        - name: host-run-ovn
          hostPath:
            path: /run/ovn
        - name: cni-conf
          hostPath:
            path: {{ .Values.cni_conf.CNI_CONF_DIR }}
        - name: cni-bin
          hostPath:
            path: {{ .Values.cni_conf.CNI_BIN_DIR }}
        - name: host-ns
          hostPath:
            path: /var/run/netns
        - name: host-dbus
          hostPath:
            path: /var/run/dbus
        - name: kube-ovn-log
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/kube-ovn
        - name: localtime
          hostPath:
            path: /etc/localtime
        - name: host-log-ovs
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/openvswitch
        - name: host-log-ovn
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/ovn
        {{- if .Values.cni_conf.MOUNT_LOCAL_BIN_DIR }}
        - name: local-bin
          hostPath:
            path: {{ .Values.cni_conf.MOUNT_LOCAL_BIN_DIR }}
        {{- end }}
        {{- if .Values.components.ENABLE_OVN_IPSEC }}
        - name: ovs-ipsec-keys
          hostPath:
            path: /etc/origin/ovs_ipsec_keys
        {{- end }}`

	ovs_ovn_daemonset = `kind: DaemonSet
apiVersion: apps/v1
metadata:
  name: ovs-ovn
  namespace: {{ .Values.namespace }}
  annotations:
    kubernetes.io/description: |
      This daemon set launches the openvswitch daemon.
    chart-version: "{{ .Chart.Name }}-{{ .Chart.Version }}"
spec:
  selector:
    matchLabels:
      app: ovs
  updateStrategy:
    type: {{ include "kubeovn.ovs-ovn.updateStrategy" . }}
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: ovs
        component: network
        type: infra
      annotations:
        chart-version: "{{ .Chart.Name }}-{{ .Chart.Version }}"
    spec:
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - effect: NoExecute
          operator: Exists
        - key: CriticalAddonsOnly
          operator: Exists
      priorityClassName: system-node-critical
      serviceAccountName: ovn-ovs
      hostNetwork: true
      hostPID: true
      initContainers:
        - name: hostpath-init
          {{- if .Values.DPDK }}
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.dpdkRepository }}:{{ .Values.DPDK_VERSION }}-{{ .Values.global.images.kubeovn.tag }}
          {{- else }}
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          {{- end }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          command:
            - sh
            - -xec
            - |
              chown -R nobody: /var/run/ovn /var/log/ovn /etc/openvswitch /var/run/openvswitch /var/log/openvswitch
              {{- if not .Values.DISABLE_MODULES_MANAGEMENT }}
              iptables -V
              {{- else }}
              ln -sf /bin/true /usr/local/sbin/modprobe
              ln -sf /bin/true /usr/local/sbin/modinfo
              ln -sf /bin/true /usr/local/sbin/rmmod
              {{- end }}
          securityContext:
            allowPrivilegeEscalation: true
            capabilities:
              drop:
                - ALL
            privileged: true
            runAsUser: 0
          volumeMounts:
            - mountPath: /usr/local/sbin
              name: usr-local-sbin
            - mountPath: /var/log/ovn
              name: host-log-ovn
            - mountPath: /var/run/ovn
              name: host-run-ovn
            - mountPath: /etc/openvswitch
              name: host-config-openvswitch
            - mountPath: /var/run/openvswitch
              name: host-run-ovs
            - mountPath: /var/log/openvswitch
              name: host-log-ovs
      containers:
        - name: openvswitch
          {{- if .Values.DPDK }}
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.dpdkRepository }}:{{ .Values.DPDK_VERSION }}-{{ .Values.global.images.kubeovn.tag }}
          {{- else }}
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          {{- end }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          {{- if .Values.DPDK }}
          command: ["/kube-ovn/start-ovs-dpdk.sh"]
          {{- else }}
          command: ["/kube-ovn/start-ovs.sh"]
          {{- end }}
          securityContext:
            runAsUser: {{ include "kubeovn.runAsUser" . }}
            privileged: false
            capabilities:
              add:
                - NET_ADMIN
                - NET_BIND_SERVICE
                {{- if not .Values.DISABLE_MODULES_MANAGEMENT }}
                - SYS_MODULE
                {{- end }}
                - SYS_NICE
                - SYS_ADMIN
          env:
            - name: ENABLE_SSL
              value: "{{ .Values.networking.ENABLE_SSL }}"
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
            - name: HW_OFFLOAD
              value: "{{- .Values.components.HW_OFFLOAD }}"
            - name: TUNNEL_TYPE
              value: "{{- .Values.networking.TUNNEL_TYPE }}"
            - name: KUBE_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: OVN_DB_IPS
              value: "{{ .Values.MASTER_NODES | default (include "kubeovn.nodeIPs" .) }}"
            - name: OVN_REMOTE_PROBE_INTERVAL
              value: "{{ .Values.networking.OVN_REMOTE_PROBE_INTERVAL }}"
            - name: OVN_REMOTE_OPENFLOW_INTERVAL
              value: "{{ .Values.networking.OVN_REMOTE_OPENFLOW_INTERVAL }}"
          volumeMounts:
            - mountPath: /usr/local/sbin
              name: usr-local-sbin
            - mountPath: /lib/modules
              name: host-modules
              readOnly: true
            - mountPath: /var/run/openvswitch
              name: host-run-ovs
            - mountPath: /var/run/ovn
              name: host-run-ovn
            - mountPath: /etc/openvswitch
              name: host-config-openvswitch
            - mountPath: /var/log/openvswitch
              name: host-log-ovs
            - mountPath: /var/log/ovn
              name: host-log-ovn
            - mountPath: /etc/localtime
              name: localtime
              readOnly: true
            - mountPath: /var/run/tls
              name: kube-ovn-tls
            - mountPath: /var/run/containerd
              name: cruntime
              readOnly: true
            {{- if .Values.DPDK }}
            - mountPath: /opt/ovs-config
              name: host-config-ovs
            - mountPath: /dev/hugepages
              name: hugepage
            {{- end }}
          readinessProbe:
            exec:
              {{- if .Values.DPDK }}
              command:
                - bash
                - /kube-ovn/ovs-dpdk-healthcheck.sh
              {{- else }}
              command:
                - bash
                - /kube-ovn/ovs-healthcheck.sh
              {{- end }}
            initialDelaySeconds: 10
            periodSeconds: 5
            timeoutSeconds: 45
          livenessProbe:
            exec:
              {{- if .Values.DPDK }}
              command:
                - bash
                - /kube-ovn/ovs-dpdk-healthcheck.sh
              {{- else }}
              command:
                - bash
                - /kube-ovn/ovs-healthcheck.sh
              {{- end }}
            initialDelaySeconds: 60
            periodSeconds: 5
            failureThreshold: 5
            timeoutSeconds: 45
          resources:
            requests:
              {{- if .Values.DPDK }}
              cpu: {{ .Values.DPDK_CPU }}
              memory: {{ .Values.DPDK_MEMORY }}
              {{- else }}
              cpu: {{ index .Values "ovs-ovn" "requests" "cpu" }}
              memory: {{ index .Values "ovs-ovn" "requests" "memory" }}
              {{- end }}
            limits:
              {{- if .Values.DPDK }}
              cpu: {{ .Values.DPDK_CPU }}
              memory: {{ .Values.DPDK_MEMORY }}
              hugepages-1Gi: 1Gi
              {{- else }}
              cpu: {{ index .Values "ovs-ovn" "limits" "cpu" }}
              memory: {{ index .Values "ovs-ovn" "limits" "memory" }}
              {{- end }}
      nodeSelector:
        kubernetes.io/os: "linux"
      volumes:
        - name: usr-local-sbin
          emptyDir: {}
        - name: host-modules
          hostPath:
            path: /lib/modules
        - name: host-run-ovs
          hostPath:
            path: /run/openvswitch
        - name: host-run-ovn
          hostPath:
            path: /run/ovn
        - name: host-config-openvswitch
          hostPath:
            path: {{ .Values.OPENVSWITCH_DIR }}
        - name: host-log-ovs
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/openvswitch
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
        - hostPath:
            path: /var/run/containerd
          name: cruntime
        {{- if .Values.DPDK }}
        - name: host-config-ovs
          hostPath:
            path: /opt/ovs-config
            type: DirectoryOrCreate
        - name: hugepage
          emptyDir:
            medium: HugePages
        {{- end }}`

	kube_ovn_pinger_daemonsets = `kind: DaemonSet
apiVersion: apps/v1
metadata:
  name: kube-ovn-pinger
  namespace: {{ .Values.namespace }}
  annotations:
    kubernetes.io/description: |
      This daemon set launches the openvswitch daemon.
spec:
  selector:
    matchLabels:
      app: kube-ovn-pinger
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: kube-ovn-pinger
        component: network
        type: infra
    spec:
      priorityClassName: system-node-critical
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - effect: NoExecute
          operator: Exists
        - key: CriticalAddonsOnly
          operator: Exists
      serviceAccountName: kube-ovn-app
      hostPID: true
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
        - name: pinger
          image: {{ .Values.global.registry.address }}/{{ .Values.global.images.kubeovn.repository }}:{{ .Values.global.images.kubeovn.tag }}
          command:
          - /kube-ovn/kube-ovn-pinger
          args:
          - --external-address=
          {{- if eq .Values.networking.NET_STACK "dual_stack" -}}
          {{ .Values.dual_stack.PINGER_EXTERNAL_ADDRESS }}
          {{- else if eq .Values.networking.NET_STACK "ipv4" -}}
          {{ .Values.ipv4.PINGER_EXTERNAL_ADDRESS }}
          {{- else if eq .Values.networking.NET_STACK "ipv6" -}}
          {{ .Values.ipv6.PINGER_EXTERNAL_ADDRESS }}
          {{- end }}
          - --external-dns=
          {{- if eq .Values.networking.NET_STACK "dual_stack" -}}
          {{ .Values.dual_stack.PINGER_EXTERNAL_DOMAIN }}
          {{- else if eq .Values.networking.NET_STACK "ipv4" -}}
          {{ .Values.ipv4.PINGER_EXTERNAL_DOMAIN }}
          {{- else if eq .Values.networking.NET_STACK "ipv6" -}}
          {{ .Values.ipv6.PINGER_EXTERNAL_DOMAIN }}
          {{- end }}
          - --ds-namespace={{ .Values.namespace }}
          - --logtostderr=false
          - --alsologtostderr=true
          - --log_file=/var/log/kube-ovn/kube-ovn-pinger.log
          - --log_file_max_size=200
          - --enable-metrics={{- .Values.networking.ENABLE_METRICS }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
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
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: HOST_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.hostIP
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
          volumeMounts:
            - mountPath: /var/run/openvswitch
              name: host-run-ovs
            - mountPath: /var/run/ovn
              name: host-run-ovn
            - mountPath: /etc/openvswitch
              name: host-config-openvswitch
            - mountPath: /var/log/openvswitch
              name: host-log-ovs
              readOnly: true
            - mountPath: /var/log/ovn
              name: host-log-ovn
              readOnly: true
            - mountPath: /var/log/kube-ovn
              name: kube-ovn-log
            - mountPath: /etc/localtime
              name: localtime
              readOnly: true
            - mountPath: /var/run/tls
              name: kube-ovn-tls
          resources:
            requests:
              cpu: {{ index .Values "kube-ovn-pinger" "requests" "cpu" }}
              memory: {{ index .Values "kube-ovn-pinger" "requests" "memory" }}
            limits:
              cpu: {{ index .Values "kube-ovn-pinger" "limits" "cpu" }}
              memory: {{ index .Values "kube-ovn-pinger" "limits" "memory" }}
      nodeSelector:
        kubernetes.io/os: "linux"
      volumes:
        - name: host-run-ovs
          hostPath:
            path: /run/openvswitch
        - name: host-run-ovn
          hostPath:
            path: /run/ovn
        - name: host-config-openvswitch
          hostPath:
            path: {{ .Values.OPENVSWITCH_DIR }}
        - name: host-log-ovs
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/openvswitch
        - name: kube-ovn-log
          hostPath:
            path: {{ .Values.log_conf.LOG_DIR }}/kube-ovn
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

	DaemonsetList = []string{ovs_ovn_dpdk_daemonset, kube_ovn_cni_daemonset, ovs_ovn_daemonset, kube_ovn_pinger_daemonsets}
)
