# kubernetes-ci
This project aims to provide support for easy kubernetes deployments for (almost) every system out there.
Kubernetes-CI uses Environment variables to do `kubectl` tasks automatically with almost every Cloud provider there is.

## Supported types of Kubernetes
| Architecture  | Type      | Kustomize | Other notes   |
| ------------  |:----:     |:---------:|:-----------:  |
| AMD64         | EKS       | Yes       | --            |
| ARMv8         | EKS       | No        | --            |
| AMD64         | Baremetal | Yes       | --            |
| ARMv8         | Baremetal | No        | Can use just plain kubectl commands due to ARM restrictions, will be fixed in the future| 

## Supported CI's
Here's a list with tested CI's and examples how to use kubernetes-ci with them.
- [Drone](docs/drone/README.md)

## Supported environment paramenters
| Paramenter            | Description                   |Required       | Default Value | Allowed Values |
| -------------         |:-------------:                |:-------------:|:-------------:|:-------------: |
| AWS_ACCESS_KEY        | AWS Access key                | YES           | -             | -              |
| AWS_SECRET_KEY        | AWS Access key secret         | YES           | -             | -              |
| AWS_REGION            | AWS Region                    | NO            | eu-west-1     | -              |
| ASSUME_ROLE           | AWS Assume role               | NO            | -             | Role ARN       |
| ACTIONS               | AWS Client command to be run  | YES           | test          | apply/delete/diff|
| KUBECTL_VERSION       | Kubectl version to be installed| NO           | v1.8.3        | -              |
| NAMESPACE             | Kubernetes namespace          | NO            | default       | -              |
| CLUSTER_NAME          | EKS Cluster name              | NO            | EKS-Cluster   | -              |
| MANIFEST_DIR          | Directory holding the manifests| NO           | ./            | -              |
| KUSTOMIZE             | Use Kustomize                 | NO            | false         | true / false   |
| IMAGE_VERSION         | Version to deploy             | NO            | -             | -              |
| IMAGE                 | Image name of the deployment. Used with Kustomize | NO | -    | -              |
| TYPE                  | Type of Kubernetes deployment | NO            | Baremetal     | EKS / Baremetal|
| TOKEN                 | Kubernetes authentication token| NO           | -             | -              |
| CA                    | Kubernetes CA certificate     | NO            | -             | -              |
| K8S_SERVER            | Kubernetes server address     | NO            | -             | -              |
| K8S_USER              | Kubernetes authentication username | NO       | default       | -              |
| ROLLOUT_CHECK         | Check rollout success         | NO            | true          | true / false   |
| ROLLOUT_TIMEOUT       | Rollout check timeout         | NO            | 1m            | Xs, Xm, Xh ... |

## Build
Build the binary with the following commands:

```shell script
export CGO_ENABLED=0
cd cmd/kubernetes-ci && go build
```

## Docker

Kubernetes-ci project uses multiarch builds with dockers buildx tool. \
Build the docker image with:
```
docker buildx build -f docker/Dockerfile --platform linux/amd64,linux/arm64/v8 -t quay.io/jlehtimaki/kubernetes-ci --push .
```

## Usage
```
docker run --rm -it -e AWS_ACCESS_KEY=.... -e PLUGIN_ASSUME_ROLE=.... -e AWS_SECRET_KEY=.... 
-e PLUGIN_ACTIONS="apply" -e PLUGIN_MANIFEST_DIR="manifests/" lehtux/drone-kubernetes
```