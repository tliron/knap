*Work in progress*

Knap
====

The Kubernetes Network Attachment Provider enables "network-as-a-service" for Kubernetes.

This [Kubernetes](https://kubernetes.io/)
[operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
introduces a new custom resource, "Network", that is used to create and manage
["NetworkAttachmentDefinition" resources](https://github.com/k8snetworkplumbingwg/network-attachment-definition-client/blob/master/artifacts/networks-crd.yaml),
which are in turn used by
[Multus CNI](https://github.com/intel/multus-cni)
to attach extra network interfaces to pods and
[KubeVirt](https://kubevirt.io/) virtual machines.


Rationale
---------

Kubernetes intentionally has little to say about networking, requiring only a single cluster-wide
network for the management platform to be able to communicate with pods, i.e. "the control plane".
While this single network can often be good enough for other communication needs, such as
connections to databases, serving HTTP to a loadbalancer, and can even be used in multi-cluster
scenarios, for many use cases it is insufficient. Additional networks might be required for access
to high performance features such as hardware acceleration, as well as enforcing isolation and
better security. Moreover, network functions (CNFs and VNFs) running in Kubernetes might be used
to connect such networks together via routing, switching, etc., i.e. "the data plane".

Multus is a great solution for adding such extra networks. However, using it directly requires
administrative knowledge of the cluster's and even the particular node's infrastructure in order to
write valid
[CNI](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/)
configurations. Moreover, CNI configurations within a cluster must be validated against each other
to ensure that there are no conflicts. This makes it challenging to create portable applications and
deploy them at scale.

Knap bridges this gap by separating the concerns. Knap itself contains the contextual administrative
knowledge and can generate CNI configurations for Multus as needed. Application developers need only
request a network according to the features they require.


Plugins
-------

Because one size cannot fit all in networking, Knap is designed with a plugin architecture, allowing
custom plugins to be created for specific and proprietary network management needs.

Note that a Knap plugin could be backed by network functions running within the cluster, and that
those network functions could in turn be using Knap (via other plugins) to gain access to other
networks.


FAQ
---

### Is this "Neutron for Kubernetes"?

If the description of Knap reads a bit like a description of the
[Neutron component in OpenStack](https://docs.openstack.org/neutron/latest/), it is intentional.
The goal is to bring Kubernetes application development in line with what we have been doing in
"legacy" cloud platforms while staying true to cloud-native orchestration principles.

For example, we absolutely do not want to mimic Neutron's isolation-centric opinions. In OpenStack
any application can create any IP subnet and allocation pool and indeed must specify it. This
requires Neutron plugins to have isolation support built in, e.g. via overlay networks, VLAN
bridging, etc. By contrast, with Knap you may very well not have any isolation features for a
particular plugin and must query for information such as IP ranges and gateways after the network
is provided.

### How do you pronounce "Knap"?

Like "nap". The "k" is silent. If emphasis of the "k" is necessary, try saying: "kay-nap".

According to the [Merriam-Webster Dictionary](https://www.merriam-webster.com/dictionary/knap)
"knap" as a noun means "a crest of a small hill" and as a verb means to "to break with a quick
blow". If you break things please fix them!
