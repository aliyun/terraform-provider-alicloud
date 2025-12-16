---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_group"
description: |-
  Provides a Alicloud NAS Access Group resource.
---

# alicloud_nas_access_group

Provides a NAS Access Group resource. File system Access Group.

In NAS, the permission group acts as a whitelist that allows you to restrict file system access. You can allow specified IP addresses or CIDR blocks to access the file system, and assign different levels of access permission to different IP addresses or CIDR blocks by adding rules to the permission group.
For information about NAS Access Group and how to use it, see [What is NAS Access Group](https://www.alibabacloud.com/help/en/nas/developer-reference/api-nas-2017-06-26-createaccessgroup)

-> **NOTE:** Available since v1.33.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_access_group&exampleId=9a1da6e1-e347-3eb1-1191-1d83d83d21fa508c4bbf&activeTab=example&spm=docs.r.nas_access_group.0.9a1da6e1e3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_nas_access_group" "foo" {
  access_group_name = "terraform-example-${random_integer.default.result}"
  access_group_type = "Vpc"
  description       = "terraform-example"
  file_system_type  = "extreme"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_nas_access_group&spm=docs.r.nas_access_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `access_group_name` - (Optional, ForceNew) The name of the permission group.
* `access_group_type` - (Optional, ForceNew) Permission group types, including Vpc.
* `description` - (Optional) Permission group description information.
* `file_system_type` - (Optional, ForceNew, Computed) File system type. Value:
  - standard (default): Universal NAS
  - extreme: extreme NAS
The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.218.0) Field 'name' has been deprecated from provider version 1.218.0. New field 'access_group_name' instead.
* `type` - (Deprecated since v1.218.0) Field 'type' has been deprecated from provider version 1.218.0. New field 'access_group_type' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<access_group_name>:<file_system_type>`.
* `create_time` - (Available since v1.218.0) Creation time.
* `region_id` - (Available since v1.256.0) The region ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Group.
* `delete` - (Defaults to 5 mins) Used when delete the Access Group.
* `update` - (Defaults to 5 mins) Used when update the Access Group.

## Import

NAS Access Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_access_group.example <access_group_name>:<file_system_type>
```
