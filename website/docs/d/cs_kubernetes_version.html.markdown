---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_version"
sidebar_current: "docs-alicloud-datasource-cs-kubernetes-version"
description: |-
Provides a list of available kubernetes version.

---

# alicloud\_cs\_kubernetes\_version

This data source provides a list of  support of Kubernetes version of the detailed information.

-> **NOTE:** Available in 1.163.0+.

## Example Usage

```terraform
data "alicloud_cs_kubernetes_version" "default" {
  cluster_type = "Kubernetes"
}

output "kubernetes_versions" {
  value = data.alicloud_cs_kubernetes_version.default.kubernetes_versions
}
```

## Argument Reference

The following arguments are supported.

* `Kubernetes` - (Required) The id of kubernetes cluster type (Kubernetes、 ManagedKubernetes、 Ask、 ExternalKubernetes).



## Attributes Reference

* `cluster_type` - The type of kubernetes cluster.

* `kubernetes_versions` - A list of support kubernetes version.

    * `runtimes` - The list of ACK runtime info.

        * `name` - Container runtime name

        * `version` - Container runtime version

    * `version` - The ACK released kubernetes version. 

    