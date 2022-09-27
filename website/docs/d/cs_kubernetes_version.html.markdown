---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_version"
sidebar_current: "docs-alicloud-datasource-cs-kubernetes-version"
description: |-
  Provides a data source to query the details of the Kubernetes version supported by ACK.

---

# alicloud\_cs\_kubernetes\_version

This data source provides the details of the Kubernetes version supported by ACK.

-> **NOTE:** Available in 1.170.0+.

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

## Argument Reference

The following arguments are supported.

* `cluster_type` - (Required) The type of cluster. Its valid value are `Kubernetes` and `ManagedKubernetes`.
* `kubernetes_version` - (Optional) The ACK released kubernetes version. 
* `profile` - (Optional) The profile of cluster. Its valid value are `Default`, `Serverless` and `Edge`.

## Attributes Reference

The following attributes are exported.

* `metadata` - A list of metadata of kubernetes version.
  * `version` - The ACK released kubernetes version. 
  * `runtime` - The list of supported runtime.
    * `name` - The runtime name.
    * `version` - The runtime version.

    