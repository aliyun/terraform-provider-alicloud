---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_namespace"
sidebar_current: "docs-alicloud-resource-edas-namespace"
description: |-
  Provides a Alicloud EDAS Namespace resource.
---

# alicloud_edas_namespace

Provides a EDAS Namespace resource.

For information about EDAS Namespace and how to use it, see [What is Namespace](https://www.alibabacloud.com/help/en/enterprise-distributed-application-service/latest/insertorupdateregion).

-> **NOTE:** Available since v1.173.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_edas_namespace&exampleId=34281039-bffb-a43d-3670-ce75c36528dc9c56a834&activeTab=example&spm=docs.r.edas_namespace.0.34281039bf&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = var.region
}

variable "region" {
  default = "cn-hangzhou"
}
variable "name" {
  default = "tfexample"
}

resource "alicloud_edas_namespace" "default" {
  debug_enable         = false
  description          = var.name
  namespace_logical_id = "${var.region}:${var.name}"
  namespace_name       = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_edas_namespace&spm=docs.r.edas_namespace.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `debug_enable` - (Optional) Specifies whether to enable remote debugging.
* `description` - (Optional) The description of the namespace, The description can be up to `128` characters in length.
* `namespace_logical_id` - (Required, ForceNew) The ID of the namespace.
  - The ID of a custom namespace is in the `region ID:namespace identifier` format. An example is `cn-beijing:tdy218`.
  - The ID of the default namespace is in the `region ID` format. An example is cn-beijing.
* `namespace_name` - (Required) The name of the namespace, The name can be up to `63` characters in length.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Namespace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Namespace.
* `delete` - (Defaults to 1 mins) Used when delete the Namespace.
* `update` - (Defaults to 1 mins) Used when update the Namespace.

## Import

EDAS Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_edas_namespace.example <id>
```