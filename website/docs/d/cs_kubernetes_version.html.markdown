---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_version"
sidebar_current: "docs-alicloud-datasource-cs-kubernetes-version"
description: |-
  Provides a data source to query the details of the Kubernetes version supported by ACK.

---

# alicloud_cs_kubernetes_version

This data source provides the details of the Kubernetes version supported by ACK.

-> **NOTE:** Available since v1.170.0.

## Example Usage

```terraform
# Query the managed kubernetes cluster metadata of version 1.22.3-aliyun.1 in the region specified by the client.
data "alicloud_cs_kubernetes_version" "default" {
  cluster_type       = "ManagedKubernetes"
  kubernetes_version = "1.22.3-aliyun.1"
  profile            = "Default"
}
output "metadata" {
  value = data.alicloud_cs_kubernetes_version.default.metadata
}
```

```terraform
# Query the kubernetes cluster metadata of version 1.22.3-aliyun.1 in the region specified by the client.
data "alicloud_cs_kubernetes_version" "default" {
  cluster_type       = "Kubernetes"
  kubernetes_version = "1.22.3-aliyun.1"
  profile            = "Default"
}
output "metadata" {
  value = data.alicloud_cs_kubernetes_version.default.metadata
}
```

```terraform
# Query the serverless kubernetes cluster metadata of version 1.22.3-aliyun.1 in the region specified by the client.
data "alicloud_cs_kubernetes_version" "default" {
  cluster_type       = "ManagedKubernetes"
  kubernetes_version = "1.22.3-aliyun.1"
  profile            = "Serverless"
}
output "metadata" {
  value = data.alicloud_cs_kubernetes_version.default.metadata
}
```

```terraform
# Query the edge kubernetes cluster metadata of version 1.20.11-aliyunedge.1 in the region specified by the client.
data "alicloud_cs_kubernetes_version" "default" {
  cluster_type       = "ManagedKubernetes"
  kubernetes_version = "1.20.11-aliyunedge.1"
  profile            = "Edge"
}
output "metadata" {
  value = data.alicloud_cs_kubernetes_version.default.metadata
}
```

### Query extended support versions

```terraform
data "alicloud_cs_kubernetes_version" "extended" {
  cluster_type   = "ManagedKubernetes"
  profile        = "Default"
  support_policy = "Extended"
}
```

## Argument Reference

The following arguments are supported.

* `cluster_type` - (Required) The type of cluster. Its valid value are `Kubernetes` and `ManagedKubernetes`.
* `kubernetes_version` - (Optional) The ACK released kubernetes version. 
* `support_policy` - (Optional, Available since v2.3.0) Version support policy used to filter metadata. Valid values: `Standard` and `Extended`. When omitted, ACK uses the Standard view. Use `Extended` to query extended support versions for ACK Pro clusters using the Default profile. Independent ACS and Edge version lines are not controlled by this policy.
* `profile` - (Optional) The profile of cluster. Its valid value are `Default`, `Serverless`, `Edge` and `Acs`.

## Attributes Reference

The following attributes are exported.

* `metadata` - A Set of Kubernetes version metadata.
  * `version` - The ACK released kubernetes version. 
  * `support_policy` - (Available since v2.3.0) The support policy required by the version (`Standard` or `Extended`).
  * `runtime` - The list of supported runtime.
    * `name` - The runtime name.
    * `version` - The runtime version.

