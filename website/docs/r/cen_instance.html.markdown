---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance"
sidebar_current: "docs-alicloud-resource-cen-instance"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Instance resource.
---

# alicloud_cen_instance

Provides a Cloud Enterprise Network (CEN) Instance resource.

For information about Cloud Enterprise Network (CEN) Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createcen).

-> **NOTE:** Available since v1.15.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_instance&exampleId=33971f6a-499d-04cb-58e9-dba3e6bf2b27cbc2c2f5&activeTab=example&spm=docs.r.cen_instance.0.33971f6a49&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}
```
## Argument Reference

The following arguments are supported:

* `protection_level` - (Optional, Available since v1.76.0) The level of CIDR block overlapping. Default value: `REDUCE`.
* `resource_group_id` - (Optional, Available since v1.232.0) The ID of the resource group. **Note:** Once you set a value of this property, you cannot set it to an empty string anymore.
* `cen_instance_name` - (Optional, Available since v1.98.0) The name of the CEN Instance. The name can be empty or `1` to `128` characters in length and cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the CEN Instance. The description can be empty or `1` to `256` characters in length and cannot start with `http://` or `https://`.
* `tags` - (Optional, Available since v1.80.0) A mapping of tags to assign to the resource.
* `name` - (Optional, Deprecated since v1.98.0) Field `name` has been deprecated from provider version 1.98.0. New field `cen_instance_name` instead.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `status` - The status of the Instance.

## Timeouts

-> **NOTE:** Available since v1.48.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.
* `delete` - (Defaults to 10 mins) Used when delete the Instance.

## Import

Cloud Enterprise Network (CEN) Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_instance.example <id>
```
