---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_addon_metadata"
sidebar_current: "docs-alicloud-datasource-cs-kubernetes-addon_metadata"
description: |-
  Provide metadata of kubernetes cluster addons.
---

# alicloud\_cs\_kubernetes\_addon\_metadata

This data source provides metadata of kubernetes cluster addons.

-> **NOTE:** Available in 1.166.0+.

## Example Usage

```terraform
data "alicloud_cs_kubernetes_addon_metadata" "default" {
  cluster_id = var.cluster_id
  name       = "nginx-ingress-controller"
  version    = "v1.1.2-aliyun.2"
}

// Output addon configuration that can be customized
output "addon_config_schema" {
  value = data.alicloud_cs_kubernetes_addons.default.config_schema
}
```

## Argument Reference

The following arguments are supported.
* `cluster_id` - (Required) The id of kubernetes cluster.
* `name` - (Required) The name of the cluster addon. You can get a list of available addons that the cluster can install by using data source `alicloud_cs_kubernetes_addons`.
* `version` - (Required) The version of the cluster addon.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `config_schema` - The addon configuration that can be customized. The returned format is the standard json schema. If return empty, it means that the addon does not support custom configuration yet.