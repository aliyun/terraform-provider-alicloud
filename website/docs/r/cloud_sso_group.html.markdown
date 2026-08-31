---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_group"
sidebar_current: "docs-alicloud-resource-cloud-sso-group"
description: |-
  Provides a Alicloud Cloud Sso Group resource.
---

# alicloud_cloud_sso_group

Provides a Cloud SSO Group resource.

For information about Cloud SSO Group and how to use it, see [What is Group](https://www.alibabacloud.com/help/en/cloudsso/latest/api-cloudsso-2021-05-15-creategroup).

-> **NOTE:** Available since v1.138.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_sso_group&exampleId=2b2bd897-d9e0-6a42-a736-7a825c303b52bdcda8d2&activeTab=example&spm=docs.r.cloud_sso_group.0.2b2bd897d9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
provider "alicloud" {
  region = "cn-shanghai"
}
data "alicloud_cloud_sso_directories" "default" {}

resource "alicloud_cloud_sso_directory" "default" {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}

locals {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}

resource "alicloud_cloud_sso_group" "default" {
  directory_id = local.directory_id
  group_name   = var.name
  description  = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_sso_group&spm=docs.r.cloud_sso_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The Description of the group. The description can be up to `1024` characters long.
* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `group_name` - (Required) The Name of the group. The name must be `1` to `128` characters in length and can contain letters, digits, periods (.), underscores (_), and hyphens (-).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Group. The value formats as `<directory_id>:<group_id>`.
* `group_id` - The GroupId of the group.

## Import

Cloud SSO Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_group.example <directory_id>:<group_id>
```
