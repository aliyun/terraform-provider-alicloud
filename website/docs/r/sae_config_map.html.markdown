---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_config_map"
sidebar_current: "docs-alicloud-resource-sae-config-map"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Config Map resource.
---

# alicloud\_sae\_config\_map

Provides a Serverless App Engine (SAE) Config Map resource.

For information about Serverless App Engine (SAE) Config Map and how to use it, see [What is Config Map](https://help.aliyun.com/).

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_sae_config_map" "example" {
  data         = jsonencode({"env.home": "/root", "env.shell": "/bin/sh"})
  name         = "Terraform Name"
  namespace_id = "cn-hangzhou:Namespace"
}

```

## Argument Reference

The following arguments are supported:

* `data` - (Required) ConfigMap instance data.
* `description` - (Optional) The Description of ConfigMap.
* `name` - (Required, ForceNew) ConfigMap instance name.
* `namespace_id` - (Required, ForceNew) The NamespaceId of ConfigMap.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Config Map.

## Import

Serverless App Engine (SAE) Config Map can be imported using the id, e.g.

```
$ terraform import alicloud_sae_config_map.example <id>
```