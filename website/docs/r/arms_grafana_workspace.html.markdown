---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_grafana_workspace"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Grafana Workspace resource.
---

# alicloud_arms_grafana_workspace

Provides a Application Real-Time Monitoring Service (ARMS) Grafana Workspace resource.



For information about Application Real-Time Monitoring Service (ARMS) Grafana Workspace and how to use it, see [What is Grafana Workspace](https://next.api.alibabacloud.com/document/ARMS/2019-08-08/ListGrafanaWorkspace).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_arms_grafana_workspace" "default" {
  grafana_version           = "9.0.x"
  description               = var.name
  resource_group_id         = data.alicloud_resource_manager_resource_groups.default.ids.0
  grafana_workspace_edition = "standard"
  grafana_workspace_name    = var.name
  tags = {
    Created = "tf"
    For     = "example"
  }
}
```

## Argument Reference

The following arguments are supported:
* `account_number` - (Optional, Available since v1.247.0) Value Description:
GrafanaWorkspaceEdition is standard, this parameter is invalid.
GrafanaWorkspaceEdition is personal_edition. This parameter is invalid. Default value: 1.
The value of GrafanaWorkspaceEdition is experts_edition. The values are respectively 10, 30, and 50. The default value is 10.
The value of GrafanaWorkspaceEdition is advanced_edition. This parameter is invalid. The default value is 100.
* `aliyun_lang` - (Optional, Available since v1.247.0) Language environment (if not filled in, default is zh):
  - zh
  - en
* `auto_renew` - (Optional, Available since v1.247.0) Whether to automatically renew. Value range:
  - true: Automatic renewal. Default value: true.
  - false: Do not renew automatically.
* `custom_account_number` - (Optional, Available since v1.247.0) The number of additional user-defined accounts. Value Description:
  - GrafanaWorkspaceEdition is standard, this parameter is invalid.
  - GrafanaWorkspaceEdition is personal_edition, this parameter is invalid.
  - GrafanaWorkspaceEdition is experts_edition, this parameter is invalid.
  - GrafanaWorkspaceEdition is advanced_edition. The value range is 0 to 2000 and is a multiple of 10. The default value is 0.
* `description` - (Optional) Description
* `duration` - (Optional, Available since v1.247.0) The time of the instance package. Valid values:
  - PricingCycle is Month, indicating monthly payment. The value range is 1 to 9.
  - PricingCycle is set to Year, indicating annual payment. The value range is 1 to 3. Default value: 1.
* `grafana_version` - (Optional) Grafana version
* `grafana_workspace_edition` - (Optional, ForceNew) The edition. **Valid values:**
  - standard: `Beta Edition(For internal testing only) `
  - personal_edition: Developer Edition
  - experts_edition: Pro Edition
  - advanced_edition: Advanced Edition
* `grafana_workspace_name` - (Required) The name of the resource
* `password` - (Optional, Available since v1.247.0) The password of the instance. It is 8 to 30 characters in length and must contain three types of characters: uppercase and lowercase letters, numbers, and special symbols. Special symbols can be:()'~! @#$%^& *-_+ =
* `pricing_cycle` - (Optional, Available since v1.247.0) The billing cycle of the package year and Month. Value: Month (default): purchase by Month. Year: Purchased by Year.
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `tags` - (Optional, Map) The tag of the resource

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - The region ID of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Grafana Workspace.
* `delete` - (Defaults to 10 mins) Used when delete the Grafana Workspace.
* `update` - (Defaults to 5 mins) Used when update the Grafana Workspace.

## Import

Application Real-Time Monitoring Service (ARMS) Grafana Workspace can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_grafana_workspace.example <id>
```