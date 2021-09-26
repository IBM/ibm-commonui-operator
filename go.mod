module github.com/IBM/ibm-commonui-operator

go 1.17

require (
	github.com/Azure/go-autorest/autorest v0.11.18 // indirect
	github.com/jetstack/cert-manager v0.10.1
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/openshift/api v3.9.0+incompatible
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.9.7
)

replace (
	k8s.io/api => k8s.io/api v0.21.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.4
	k8s.io/client-go => k8s.io/client-go v0.21.4
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.9.7
)
