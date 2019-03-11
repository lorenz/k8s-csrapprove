# K8s CSR Approver
*Auto-approves CSRs from kubelets that are up for renewal*

Contains a wrapper around the Kubernetes in-tree node renewal approver.

## Installation
I'm currently not offering pre-built images. Build the Dockerfile in this repository and push the
image somewhere where your Kubernetes cluster can access it.
Modify the `image` property in the supplied `k8s.yml` to point to the image location and apply it
to your cluster. That's it :)