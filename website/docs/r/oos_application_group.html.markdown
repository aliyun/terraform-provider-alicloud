---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_application_group"
sidebar_current: "docs-alicloud-resource-oos-application-group"
description: |-
  Provides a Alicloud OOS Application Group resource.
---

# alicloud_oos_application_group

Provides a OOS Application Group resource.

For information about OOS Application Group and how to use it, see [What is Application Group](https://www.alibabacloud.com/help/en/operation-orchestration-service/latest/api-oos-2019-06-01-createapplicationgroup).

-> **NOTE:** Available since v1.146.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oos_application_group&exampleId=81fd02c9-75a4-9cf3-ff63-671f74e22682543633d7&activeTab=example&spm=docs.r.oos_application_group.0.81fd02c975&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_application" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  application_name  = "${var.name}-${random_integer.default.result}"
  description       = var.name
  tags = {
    Created = "TF"
  }
}
data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_oos_application_group" "default" {
  application_group_name = var.name
  application_name       = alicloud_oos_application.default.id
  deploy_region_id       = data.alicloud_regions.default.regions.0.id
  description            = var.name
  import_tag_key         = "example_key"
  import_tag_value       = "example_value"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oos_application_group&spm=docs.r.oos_application_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `application_group_name` - (Required, ForceNew) The name of the Application group.
* `application_name` - (Required, ForceNew) The name of the Application.
* `deploy_region_id` - (Required, ForceNew) The region ID of the deployment.
* `description` - (Optional, ForceNew) Application group description information.
* `import_tag_key` - (Optional, ForceNew) The tag key must be passed in at the same time as the tag value (import_tag_value) or none, not just one. If both `import_tag_key` and `import_tag_value` are left empty, the default is app-{ApplicationName} (application name).
* `import_tag_value` - (Optional, ForceNew) The tag value must be passed in at the same time as the tag key (import_tag_key) or none, not just one. If both `import_tag_key` and `import_tag_value` are left empty, the default is application group name.
.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Application Group. The value formats as `<application_name>:<application_group_name>`.

## Import

OOS Application Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_oos_application_group.example <application_name>:<application_group_name>
```