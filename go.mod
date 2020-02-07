module github.com/dgoodwin/argocd-appgenerator

go 1.13

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/argoproj/argo-cd v1.4.2
	github.com/argoproj/pkg v0.0.0-20200102163130-2dd1f3f6b4de // indirect
	github.com/go-logr/logr v0.1.0
	github.com/gogo/protobuf v1.3.1 // manual
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/robfig/cron v1.2.0 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	k8s.io/kubernetes v1.17.2 // indirect
	sigs.k8s.io/controller-runtime v0.4.0
)

replace golang.org/x/net => github.com/golang/net v0.0.0-20200202094626-16171245cfb2

replace k8s.io/api => k8s.io/api v0.17.2

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.2

replace k8s.io/apimachinery => k8s.io/apimachinery v0.17.3-beta.0

replace k8s.io/apiserver => k8s.io/apiserver v0.17.2

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.2

replace k8s.io/client-go => k8s.io/client-go v0.17.2

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.2

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.2

replace k8s.io/code-generator => k8s.io/code-generator v0.17.3-beta.0

replace k8s.io/component-base => k8s.io/component-base v0.17.2

replace k8s.io/cri-api => k8s.io/cri-api v0.17.3-beta.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.2

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.2

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.2

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.2

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.2

replace k8s.io/kubectl => k8s.io/kubectl v0.17.2

replace k8s.io/kubelet => k8s.io/kubelet v0.17.2

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.2

replace k8s.io/metrics => k8s.io/metrics v0.17.2

replace k8s.io/node-api => k8s.io/node-api v0.17.2

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.2

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.17.2

replace k8s.io/sample-controller => k8s.io/sample-controller v0.17.2
