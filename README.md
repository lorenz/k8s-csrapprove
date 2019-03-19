# K8s CSR Approver

_Auto-approves CSRs from kubelets that are up for renewal_

The actual approver is based on the original in-tree approver that was taken out and has been extended
to also support kubelet server certificates.

## Validation

The only accepted SANs are the name of the node and any IP where the node can present a valid certificate
with its name as CN. This should be both secure and enough to satisfy both client and server validation.

## Installation

I'm currently not offering pre-built images. Build the Dockerfile in this repository and push the
image somewhere where your Kubernetes cluster can access it.
Modify the `image` property in the supplied `k8s.yml` to point to the image location and apply it
to your cluster. That's it :)
