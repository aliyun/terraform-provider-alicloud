---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_sls_group"
sidebar_current: "docs-alicloud-resource-cms-sls-group"
description: |-
  Provides a Alicloud Cloud Monitor Service Sls Group resource.
---

# alicloud\_cms\_sls\_group

Provides a Cloud Monitor Service Sls Group resource.

For information about Cloud Monitor Service Sls Group and how to use it, see [What is Sls Group](https://www.alibabacloud.com/help/doc-detail/28608.htm).

-> **NOTE:** Available in v1.171.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_account" "this" {}

resource "alicloud_cms_sls_group" "default" {
  sls_group_config {
    sls_user_id  = data.alicloud_account.this.id
    sls_logstore = "Logstore-ECS"
    sls_project  = "aliyun-project"
    sls_region   = "cn-hangzhou"
  }
  sls_group_description = var.name
  sls_group_name        = var.name
}
```
## Argument Reference

The following arguments are supported:

* `sls_group_config` - (Required) The Config of the Sls Group. You can specify up to 25 Config. See the following `Block sls_group_config`.
* `sls_group_description` - (Optional) The Description of the Sls Group.
* `sls_group_name` - (Required, ForceNew) The name of the resource. The name must be `2` to `32` characters in length, and can contain letters, digits and underscores (_). It must start with a letter.

#### Block sls_group_config

The sls_group_config supports the following: 

* `sls_logstore` - (Required) The name of the Log Store.
* `sls_project` - (Required) The name of the Project.
* `sls_region` - (Required) The Sls Region.
* `sls_user_id` - (Optional) The ID of the Sls User.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Sls Group. Its value is same as `sls_group_name`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Sls Group.
* `delete` - (Defaults to 2 mins) Used when delete the Sls Group.
* `update` - (Defaults to 2 mins) Used when update the Sls Group.

## Import

Cloud Monitor Service Sls Group can be imported using the id, e.g.

```
$ terraform import alicloud_cms_sls_group.example <sls_group_name>
```