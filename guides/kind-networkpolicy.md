# Kind NetworkPolicy

The default CNI plugin for [Kind](https://github.com/aojea/kindnet) is [Kindnet](https://github.com/aojea/kindnet) its focus on simplicty might be the reason it doesn't currently support NetworkPolicy but fortunately Kind is flexible enough to work with other CNI providers such as [Project Calico CNI](https://github.com/projectcalico/cni-plugin).

## Project Calico CNI - Kind Installation
First step to setup Kind with Project Calico CNI is to create a Kind cluster without the default Kindnet plugin which Kind supports via its cluster yaml configuration file by specifying the setting `disableDefaultCNI: true` under the `networking` configuration, as an example this is how it would look like in the current KONK `cluster.yaml` file

```
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  image: kindest/node:v1.21.2
  extraPortMappings:
  - containerPort: 31080 # expose port 31380 of the node to port 80 on the host, later to be use by kourier or contour ingress
    listenAddress: 127.0.0.1
    hostPort: 80
networking:
  disableDefaultCNI: true # disable kindnet
```
and this is the command that would use the file `kind create cluster --config=cluster.yaml`

You can check the absence of Kubenet in the cluster by checking the pods in the `kube-system` namespace as follows `kubectl get po -n kube-system`, no pod with the name `kindnet` should be present.

### Project Calico CNI
The current version of Project Calico CNI can be installed with the following command `kubectl apply -f https://docs.projectcalico.org/v3.18/manifests/calico.yaml` next a couple of Calico pods should be present in `kube-system` namely `calico-kube-controllers` and `calico-node`.

At this point the rest of the usual Knative installation steps can be performed as usual.

To check that the NetworkPolicy works you can deploy one that denies all pod traffic like this:
```
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

As long as the NetwokPolicy is in place you won't be able to access any service backed by pods in the namespace the NetworkPolicy was created.
