---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Cen Instance resource.
---

# alicloud_cen_instance

Provides a Cloud Enterprise Network (CEN) Cen Instance resource.



For information about Cloud Enterprise Network (CEN) Cen Instance and how to use it, see [What is Cen Instance](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createcen).

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cen_instance&spm=docs.r.cen_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cen_instance_name` - (Optional, Available since v1.98.0) The name of the CEN instance.
* `description` - (Optional) The description of the CEN instance.
* `protection_level` - (Optional, Available since v1.76.0) The level of CIDR block overlapping. Valid values:  REDUCED: Overlapped CIDR blocks are allowed. However, the overlapped CIDR blocks cannot be the same.
* `resource_group_id` - (Optional, Computed, Available since v1.232.0) The ID of the resource group
* `tags` - (Optional, Map, Available since v1.80.0) The tags of the CEN instance.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.98.0). Field 'name' has been deprecated from provider version 1.246.0. New field 'cen_instance_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the CEN instance was created.
* `status` - The state of the CEN instance.   Creating: The CEN instance is being created. Active: The CEN instance is running. Deleting: The CEN instance is being deleted.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cen Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Cen Instance.
* `update` - (Defaults to 5 mins) Used when update the Cen Instance.

## Import

Cloud Enterprise Network (CEN) Cen Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_instance.example <id>
```