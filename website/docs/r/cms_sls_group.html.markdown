---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_sls_group"
sidebar_current: "docs-alicloud-resource-cms-sls-group"
description: |-
  Provides a Alicloud Cloud Monitor Service Sls Group resource.
---

# alicloud_cms_sls_group

Provides a Cloud Monitor Service Sls Group resource.

For information about Cloud Monitor Service Sls Group and how to use it, see [What is Sls Group](https://www.alibabacloud.com/help/doc-detail/28608.htm).

-> **NOTE:** Available since v1.171.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_sls_group&exampleId=2bcfe684-0054-dde6-466c-1c3de8d06c1648e8ad59&activeTab=example&spm=docs.r.cms_sls_group.0.2bcfe68400&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_account" "default" {}
data "alicloud_regions" "default" {
  current = true
}
resource "random_uuid" "default" {
}
resource "alicloud_log_project" "default" {
  project_name = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.project_name
  logstore_name         = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_cms_sls_group" "default" {
  sls_group_config {
    sls_user_id  = data.alicloud_account.default.id
    sls_logstore = alicloud_log_store.default.logstore_name
    sls_project  = alicloud_log_project.default.project_name
    sls_region   = data.alicloud_regions.default.regions.0.id
  }
  sls_group_description = var.name
  sls_group_name        = var.name
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_sls_group&spm=docs.r.cms_sls_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `sls_group_config` - (Required) The Config of the Sls Group. You can specify up to 25 Config. See [`sls_group_config`](#sls_group_config) below. 
* `sls_group_description` - (Optional) The Description of the Sls Group.
* `sls_group_name` - (Required, ForceNew) The name of the resource. The name must be `2` to `32` characters in length, and can contain letters, digits and underscores (_). It must start with a letter.

### `sls_group_config`

The sls_group_config supports the following: 

* `sls_logstore` - (Required) The name of the Log Store.
* `sls_project` - (Required) The name of the Project.
* `sls_region` - (Required) The Sls Region.
* `sls_user_id` - (Optional) The ID of the Sls User.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Sls Group. Its value is same as `sls_group_name`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Sls Group.
* `delete` - (Defaults to 2 mins) Used when delete the Sls Group.
* `update` - (Defaults to 2 mins) Used when update the Sls Group.

## Import

Cloud Monitor Service Sls Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_sls_group.example <sls_group_name>
```
