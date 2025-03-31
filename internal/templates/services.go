package templates

var (
	kube_ovn_controller_template = `kind: Service
apiVersion: v1
metadata:
  name: kube-ovn-controller
  namespace: {{ .Values.namespace }}
  labels:
    app: kube-ovn-controller
spec:
  selector:
    app: kube-ovn-controller
  ports:
    - port: 10660
      name: metrics
  {{- if eq .Values.networking.NET_STACK "dual_stack" }}
  ipFamilyPolicy: PreferDualStack
  {{- end }}`

	kube_ovn_monitor_template = `kind: Service
apiVersion: v1
metadata:
  name: kube-ovn-monitor
  namespace: {{ .Values.namespace }}
  labels:
    app: kube-ovn-monitor
spec:
  ports:
    - name: metrics
      port: 10661
  type: ClusterIP
  selector:
    app: kube-ovn-monitor
  sessionAffinity: None
  {{- if eq .Values.networking.NET_STACK "dual_stack" }}
  ipFamilyPolicy: PreferDualStack
  {{- end }}`

	nb_nb_template = `kind: Service
apiVersion: v1
metadata:
  name: ovn-nb
  namespace: {{ .Values.namespace }}
spec:
  ports:
    - name: ovn-nb
      protocol: TCP
      port: 6641
      targetPort: 6641
  type: ClusterIP
  {{- if eq .Values.networking.NET_STACK "dual_stack" }}
  ipFamilyPolicy: PreferDualStack
  {{- end }}
  selector:
    app: ovn-central
    ovn-nb-leader: "true"
  sessionAffinity: None`

	ovn_northd_template = `kind: Service
apiVersion: v1
metadata:
  name: ovn-northd
  namespace: {{ .Values.namespace }}
spec:
  ports:
    - name: ovn-northd
      protocol: TCP
      port: 6643
      targetPort: 6643
  type: ClusterIP
  {{- if eq .Values.networking.NET_STACK "dual_stack" }}
  ipFamilyPolicy: PreferDualStack
  {{- end }}
  selector:
    app: ovn-central
    ovn-northd-leader: "true"
  sessionAffinity: None`

	kube_ovn_cni_template = `
kind: Service
apiVersion: v1
metadata:
  name: kube-ovn-cni
  namespace: {{ .Values.namespace }}
  labels:
    app: kube-ovn-cni
spec:
  selector:
    app: kube-ovn-cni
  ports:
    - port: 10665
      name: metrics
  {{- if eq .Values.networking.NET_STACK "dual_stack" }}
  ipFamilyPolicy: PreferDualStack
  {{- end }}`

	kube_ovn_pinger_template = `kind: Service
apiVersion: v1
metadata:
  name: kube-ovn-pinger
  namespace: {{ .Values.namespace }}
  labels:
    app: kube-ovn-pinger
spec:
  selector:
    app: kube-ovn-pinger
  ports:
    - port: 8080
      name: metrics
  {{- if eq .Values.networking.NET_STACK "dual_stack" }}
  ipFamilyPolicy: PreferDualStack
  {{- end }}`

	ovn_sb_template = `kind: Service
apiVersion: v1
metadata:
  name: ovn-sb
  namespace: {{ .Values.namespace }}
spec:
  ports:
    - name: ovn-sb
      protocol: TCP
      port: 6642
      targetPort: 6642
  type: ClusterIP
  {{- if eq .Values.networking.NET_STACK "dual_stack" }}
  ipFamilyPolicy: PreferDualStack
  {{- end }}
  selector:
    app: ovn-central
    ovn-sb-leader: "true"
  sessionAffinity: None`
)
