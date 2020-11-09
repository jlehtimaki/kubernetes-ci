# kubernetes-ci
This project aims to provide support for easy kubernetes deployments for (almost) every system out there.

# Supported types of Kubernetes
| Architecture  | Type      | Kustomize | Other notes   |
| ------------  |:----:     |:---------:|:-----------:  |
| AMD64         | EKS       | Yes       | --            |
| ARMv8         | EKS       | No        | --            |
| AMD64         | Baremetal | Yes       | --            |
| ARMv8         | Baremetal | No        | Can use just plain kubectl commands due to ARM restrictions, will be fixed in the future| 

# Supported CI
- Drone

# Examples
## Drone pipeline with AWS
```yaml
kind: pipeline
type: kubernetes
name: Drone example pipeline

steps:
  - name: Deploy test app
    image: lehtux/drone-kubernetes
    environment:
      AWS_REGION: "eu-west-1"
      AWS_ACCESS_KEY_ID:
        from_secret: access_key
      AWS_SECRET_ACCESS_KEY:
        from_secret: access_key_secret
    settings:
      type: EKS
      assume_role: arn:aws:iam::xxxxxx:role/EKS
      actions: ["apply"]
      namespace: default
      kubectl_version: v1.16.6
      manifest_dir: deployments/deployment.yml

```

# Build
Build the binary with the following commands:

```shell script
export CGO_ENABLED=0
cd cmd/drone && go build
```

# Docker

Build the docker image with:
```
docker buildx build -f docker/drone/Dockerfile --platform linux/amd64,linux/arm64/v8 -t lehtux/drone-kubernetes --push .
```

# Usage
```
docker run --rm -it -e AWS_ACCESS_KEY=.... -e PLUGIN_ASSUME_ROLE=.... -e AWS_SECRET_KEY=.... 
-e PLUGIN_ACTIONS="apply" -e PLUGIN_MANIFEST_DIR="manifests/" lehtux/drone-kubernetes
```