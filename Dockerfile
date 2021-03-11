# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Copy license
RUN mkdir /licenses
COPY LICENSE /licenses/LICENSE

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

### Required OpenShift Labels
LABEL name="RHM Arcade Operator" \
  vendor="Red Hat Marketplace" \
  version="v0.0.1" \
  release="1" \
  summary="This is an example Arcade operator provided by Red Hat Marketplace." \
  description="This operator will deploy the Red Hat Marketplace arcade to the cluster."

ARG REGISTRY_HOST
ARG REGISTRY_REPO
ENV REGISTRY_HOST=${REGISTRY_HOST}
ENV REGISTRY_REPO=${REGISTRY_REPO}

WORKDIR /
COPY --from=builder /workspace/manager .
COPY --from=builder /licenses /licenses
USER 1001:1001

ENTRYPOINT ["/manager"]
