# Drone Kubernetes-CI

Drone can use Kubernetes-CI really easily using `settings` parameters. \
Drone converts these to PLUGIN_ prefixes. Kubernetes-CI then fetches those environment variables and does tasks
accordingly.  

## Examples
### Drone pipeline with AWS
```yaml
kind: pipeline
type: kubernetes
name: Drone example pipeline

steps:
  - name: Deploy test app
    image: quay.io/jlehtimaki/kubernetes-ci
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
      kustomize: true
      image: foobar
      image_version: ${DRONE_COMMIT_SHA}
      manifest_dir: deployments/deployment.yml
```

### Drone pipeline with Baremetal
```yaml
kind: pipeline
type: kubernetes
name: Drone example pipeline

steps:
  - name: Deploy test app
    image: quay.io/jlehtimaki/kubernetes-ci
    settings:
      type: Baremetal
      token:
        from_secret: token
      ca:
        from_secret: ca
      k8s_server: https://1.2.3.4:6666
      k8s_user: kubernetes-admin
      actions: ["apply"]
      namespace: default
      kubectl_version: v1.16.6
      manifest_dir: deployments/deployment.yml