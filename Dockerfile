FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS builder
WORKDIR /go/src/github.com/mfojtik/rebase-operator
COPY . .
ENV GO_PACKAGE github.com/mfojtik/rebase-operator
RUN make build --warn-undefined-variables

FROM registry.svc.ci.openshift.org/ocp/4.2:base
COPY --from=builder /go/src/github.com/mfojtik/rebase-operator/bugzilla-operator /usr/bin/

