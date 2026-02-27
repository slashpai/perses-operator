# perses-operator Helm Chart

A Helm chart for deploying the [Perses Operator](https://github.com/perses/perses-operator) on Kubernetes.

## Prerequisites

- Kubernetes 1.27+
- Helm 3.x
- [cert-manager](https://cert-manager.io/) installed in the cluster (required for webhook TLS certificates)

## Installation

```bash
helm install perses-operator charts/perses-operator \
  --namespace perses-operator-system \
  --create-namespace
```

## Configuration

See [values.yaml](values.yaml) for the full list of configurable parameters.

Key configuration options:

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `manager.replicas` | Number of operator replicas | `1` |
| `manager.image.repository` | Operator image repository | `docker.io/persesdev/perses-operator` |
| `manager.image.tag` | Operator image tag | `v0.2.0` |
| `manager.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `manager.resources` | Manager container resources | CPU: 10m-500m, Memory: 64Mi-128Mi |
| `kubeRbacProxy.image.repository` | kube-rbac-proxy image | `gcr.io/kubebuilder/kube-rbac-proxy` |
| `kubeRbacProxy.image.tag` | kube-rbac-proxy image tag | `v0.13.1` |
| `crd.enable` | Install CRDs with the chart | `true` |
| `crd.keep` | Keep CRDs when uninstalling | `true` |
| `certManager.enable` | Enable cert-manager integration | `true` |
| `webhook.enable` | Enable conversion webhooks | `true` |
| `metrics.enable` | Enable metrics endpoint | `true` |
| `prometheus.enable` | Enable ServiceMonitor | `false` |

## Uninstallation

```bash
helm uninstall perses-operator --namespace perses-operator-system
```

> **Note:** CRDs are retained by default (`crd.keep: true`). To remove them manually:
>
> ```bash
> kubectl delete crd perses.perses.dev persesdashboards.perses.dev persesdatasources.perses.dev persesglobaldatasources.perses.dev
> ```

## Testing

### Lint and Template Validation

```bash
make helm-lint       # Lint the chart
make helm-template   # Render templates locally
```

### Local Deployment Test

Requires [Go](https://go.dev/doc/install), [kind](https://kind.sigs.k8s.io/), [Docker](https://docs.docker.com/get-docker/) or [Podman](https://podman.io/), and [kubectl](https://kubernetes.io/docs/tasks/tools/).

```bash
make helm-test-setup     # Create cluster, install cert-manager, build image, deploy via Helm
make helm-test-cleanup   # Tear down the kind cluster
```

| Variable | Default | Description |
| -------- | ------- | ----------- |
| `HELM_KIND_CLUSTER` | `perses-helm-test` | Name of the kind cluster |
| `HELM_TEST_IMG` | `$(IMAGE_TAG_BASE):helm-test` | Operator image for testing |
| `HELM_NAMESPACE` | `perses-operator-system` | Namespace for the Helm release |
| `HELM_RELEASE` | `perses-operator` | Name of the Helm release |

### Debugging

```bash
make helm-status    # Check Helm release status

# Operator logs
kubectl -n perses-operator-system logs deployment/perses-operator-controller-manager -f

# Uninstall and redeploy
make helm-uninstall
make helm-test-setup
```