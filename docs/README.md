# Usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

    helm repo add forward https://udhos.github.io/forward

Update files from repo:

    helm repo update

Search forward:

    helm search repo forward -l --version ">=0.0.0"
    NAME                        CHART VERSION	APP VERSION	DESCRIPTION
    forward/forward          	0.1.0        	0.0.0     	A Helm chart for installing Forward into Kubernetes

To install the charts:

    helm install my-forward forward/forward
    #            ^          ^       ^
    #            |          |        \______ chart
    #            |          |
    #            |           \______________ repo
    #            |
    #             \_________________________ release (chart instance installed in cluster)

To uninstall the charts:

    helm uninstall my-forward

# Source

<https://github.com/udhos/forward>
