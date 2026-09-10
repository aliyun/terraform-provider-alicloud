---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_command"
sidebar_current: "docs-alicloud-resource-ecd-command"
description: |-
  Provides a Alicloud ECD Command resource.
---

# alicloud_ecd_command

Provides a ECD Command resource.

For information about ECD Command and how to use it, see [What is Command](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-ecd-2020-09-30-runcommand).

-> **NOTE:** Available since v1.146.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecd_command&exampleId=4a0a6ea4-6841-14d6-2744-0a7541bbe40cedc43f98&activeTab=example&spm=docs.r.ecd_command.0.4a0a6ea468&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-beijing"
}

variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  enable_admin_access = true
  desktop_access_type = "Internet"
  office_site_name    = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard         = "read"
  local_drive       = "read"
  usb_redirect      = "off"
  watermark         = "off"

  authorize_access_policy_rules {
    description = var.name
    cidr_ip     = "1.2.3.45/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = var.name
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "1.2.3.4/24"
  }
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
  bundle_id   = ["bundle_eds_enterprise_office_4c8g_s8d2_win2019_edu", "bundle_eds_enterprise_office_8c16g_s8d2_win2019_edu"]
}

resource "alicloud_ecd_desktop" "default" {
  office_site_id  = alicloud_ecd_simple_office_site.default.id
  policy_group_id = alicloud_ecd_policy_group.default.id
  bundle_id       = data.alicloud_ecd_bundles.default.bundles.0.id
  desktop_name    = var.name
}

resource "alicloud_ecd_command" "default" {
  command_content = "ipconfig"
  command_type    = "RunPowerShellScript"
  desktop_id      = alicloud_ecd_desktop.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecd_command&spm=docs.r.ecd_command.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `command_content` - (Required, ForceNew) The Contents of the Script to Base64 Encoded Transmission.
* `command_type` - (Required, ForceNew) The Script Type. Valid values: `RunBatScript`, `RunPowerShellScript`.
* `content_encoding` - (Optional, Computed) That Returns the Data Encoding Method. Valid values: `Base64`, `PlainText`.
* `desktop_id` - (Required, ForceNew) The desktop id of the Desktop.
* `timeout` - (Optional) The timeout period for script execution the unit is seconds. Default to: `60`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Command.
* `status` - Script Is Executed in the Overall Implementation of the State. Valid values: `Pending`, `Failed`, `PartialFailed`, `Running`, `Stopped`, `Stopping`, `Finished`, `Success`.

## Import

ECD Command can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_command.example <id>
```