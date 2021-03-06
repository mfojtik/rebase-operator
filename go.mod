module github.com/mfojtik/rebase-operator

go 1.14

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/google/go-github/v32 v32.0.0
	github.com/gorilla/handlers v1.4.2
	github.com/openshift/build-machinery-go v0.0.0-20200512074546-3744767c4131
	github.com/openshift/library-go v0.0.0-20200512120242-21a1ff978534
	github.com/prometheus/client_golang v1.6.0 // indirect
	github.com/shomali11/commander v0.0.0-20191122162317-51bc574c29ba
	github.com/shomali11/proper v0.0.0-20190608032528-6e70a05688e7
	github.com/slack-go/slack v0.6.4
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.18.2
	k8s.io/component-base v0.18.2
	k8s.io/klog v1.0.0
)

replace github.com/eparis/bugzilla => github.com/sttts/bugzilla v0.0.0-20200525151909-b7660389ebf3
