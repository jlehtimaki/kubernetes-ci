# Building the drone-kubernetes
FROM --platform=$BUILDPLATFORM golang:1.14-alpine AS builder

ARG BUILDPLATFORM
ARG TARGETARCH

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=$TARGETARCH


WORKDIR /build
COPY . .
RUN cd cmd/drone && go build -o /build/drone-kubernetes

# Kubectl, awscli and kustomize binaries
FROM --platform=$BUILDPLATFORM ubuntu:20.04 AS executables
ARG TARGETARCH
ARG BUILDPLATFORM

WORKDIR /execs
RUN apt update && apt install -y curl unzip bash

# Fetch Kubectl
RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/v1.18.3/bin/linux/$TARGETARCH/kubectl"
RUN chmod +x kubectl

# Fetch awscli installation
RUN if [ "$TARGETARCH" == "arm64" ]; then \
    curl "curl -O 'https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip";  \
    else \
      curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"; \
    fi
RUN unzip awscliv2.zip

# Fetch kustomize
RUN curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash

# Actual drone-kubernetes docker image
FROM --platform=$BUILDPLATFORM ubuntu:20.04

ARG BUILDPLATFORM
ARG TARGETARCH
ENV GOARCH=$TARGETARCH

# Install AWSCli needed binaries & certs
RUN apt update && apt install -y less ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy binaries
COPY --from=executables /execs/kubectl /bin/
COPY --from=executables /execs/aws .
COPY --from=executables /execs/kustomize /bin/

# Install AWSCli
RUN ./install -b /bin/

# Copy drone-kubernetes binary
COPY --from=builder /build/drone-kubernetes /bin/

ENTRYPOINT ["/bin/drone-kubernetes"]