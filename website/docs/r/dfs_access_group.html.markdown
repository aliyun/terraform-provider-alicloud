---
subcategory: "Apsara File Storage for HDFS (DFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_access_group"
description: |-
  Provides a Alicloud DFS Access Group resource.
---

# alicloud_dfs_access_group

Provides a DFS Access Group resource. 

For information about DFS Access Group and how to use it, see [What is Access Group](https://www.alibabacloud.com/help/en/aibaba-cloud-storage-services/latest/apsara-file-storage-for-hdfs).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dfs_access_group&exampleId=4e4954f0-6ed7-76e3-4a7e-954b8a06333489222da3&activeTab=example&spm=docs.r.dfs_access_group.0.4e4954f06e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_dfs_access_group" "default" {
  access_group_name = "tf-example-${random_integer.default.result}"
  network_type      = "VPC"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_dfs_access_group&spm=docs.r.dfs_access_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `access_group_name` - (Required) The permission group name. The naming rules are as follows: The length is 6~64 characters. Globally unique and cannot be an empty string. English letters are supported and can contain numbers, underscores (_), and dashes (-).
* `description` - (Optional) The permission group description.  No more than 32 characters in length.
* `network_type` - (Required, ForceNew) The permission group type. Only VPC (VPC) is supported.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the permission group resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Group.
* `delete` - (Defaults to 5 mins) Used when delete the Access Group.
* `update` - (Defaults to 5 mins) Used when update the Access Group.

## Import

DFS Access Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_dfs_access_group.example <id>
```