---
subcategory: "Schedulerx"
layout: "alicloud"
page_title: "Alicloud: alicloud_schedulerx_namespace"
sidebar_current: "docs-alicloud-resource-schedulerx-namespace"
description: |- 
    Provides a Alicloud Schedulerx Namespace resource.
---

# alicloud\_schedulerx\_namespace

Provides a Schedulerx Namespace resource.

For information about Schedulerx Namespace and how to use it, see [What is Namespace](https://help.aliyun.com/document_detail/206088.html).

-> **NOTE:** Available in v1.173.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_schedulerx_namespace&exampleId=0f5e16f3-668c-8d7f-70f8-e3adc3f2b9948ccbbcba&activeTab=example&spm=docs.r.schedulerx_namespace.0.0f5e16f366&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_schedulerx_namespace" "example" {
  namespace_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the resource.
* `namespace_name` - (Required) The name of the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Namespace. Its value is same as `namespace_id`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the resource.
* `update` - (Defaults to 1 mins) Used when update the resource.



## Import

Schedulerx Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_schedulerx_namespace.example <id>
```