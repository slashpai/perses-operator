# perses-operator

A Helm chart for the Perses Operator - manages Perses instances and dashboards on Kubernetes

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

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| certManager.enable | bool | `true` | Enable cert-manager integration (required for webhook certificates) |
| crd.enable | bool | `true` | Install CRDs with the chart |
| crd.keep | bool | `true` | Keep CRDs when uninstalling |
| kubeRbacProxy.image.repository | string | `"gcr.io/kubebuilder/kube-rbac-proxy"` | kube-rbac-proxy image repository |
| kubeRbacProxy.image.tag | string | `"v0.13.1"` | kube-rbac-proxy image tag |
| kubeRbacProxy.resources | object | `{"limits":{"cpu":"500m","memory":"128Mi"},"requests":{"cpu":"5m","memory":"64Mi"}}` | Resource limits and requests for the kube-rbac-proxy container |
| kubeRbacProxy.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]}}` | Security context for the kube-rbac-proxy container |
| manager.affinity | object | `{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"kubernetes.io/arch","operator":"In","values":["amd64","arm64","ppc64le","s390x"]},{"key":"kubernetes.io/os","operator":"In","values":["linux"]}]}]}}}` | Affinity rules for manager pods |
| manager.args | list | `["--leader-elect"]` | Extra arguments passed to the manager container |
| manager.env | list | `[]` | Environment variables for the manager container |
| manager.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| manager.image.repository | string | `"docker.io/persesdev/perses-operator"` | Operator image repository |
| manager.image.tag | string | `"v0.2.0"` | Operator image tag |
| manager.imagePullSecrets | list | `[]` | Image pull secrets |
| manager.nodeSelector | object | `{}` | Node selector for manager pods |
| manager.podSecurityContext | object | `{}` | Pod-level security context |
| manager.replicas | int | `1` | Number of operator replicas |
| manager.resources | object | `{"limits":{"cpu":"500m","memory":"128Mi"},"requests":{"cpu":"10m","memory":"64Mi"}}` | Resource limits and requests for the manager container |
| manager.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]}}` | Container-level security context |
| manager.tolerations | list | `[]` | Tolerations for manager pods |
| metrics.enable | bool | `true` | Enable metrics endpoint with RBAC protection |
| metrics.port | int | `8082` | Metrics server port |
| prometheus.enable | bool | `false` | Enable ServiceMonitor (requires prometheus-operator) |
| rbacHelpers.enable | bool | `false` | Install convenience admin/editor/viewer roles for CRDs |
| webhook.enable | bool | `true` | Enable conversion webhooks |
| webhook.port | int | `9443` | Webhook server port |

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
