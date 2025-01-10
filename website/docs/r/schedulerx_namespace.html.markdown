---
subcategory: "Schedulerx"
layout: "alicloud"
page_title: "Alicloud: alicloud_schedulerx_namespace"
description: |-
  Provides a Alicloud Schedulerx Namespace resource.
---

# alicloud_schedulerx_namespace

Provides a Schedulerx Namespace resource.



For information about Schedulerx Namespace and how to use it, see [What is Namespace](https://www.alibabacloud.com/help/en/schedulerx/schedulerx-serverless/developer-reference/api-schedulerx2-2019-04-30-listnamespaces).

-> **NOTE:** Available since v1.173.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_schedulerx_namespace&exampleId=204fe6a4-76f7-53ee-03c8-cb092c6f945757d24220&activeTab=example&spm=docs.r.schedulerx_namespace.0.204fe6a476&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_schedulerx_namespace" "default" {
  namespace_name = var.name
  description    = var.name
}
```

### Deleting `alicloud_schedulerx_namespace` or removing it from your configuration

Terraform cannot destroy resource `alicloud_schedulerx_namespace`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) Namespace description.
* `namespace_name` - (Required, ForceNew) Namespace name.
* `namespace_uid` - (Optional, ForceNew, Computed,  Available since v1.240.0) Namespace uid.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Namespace.

## Import

Schedulerx Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_schedulerx_namespace.example <id>
```