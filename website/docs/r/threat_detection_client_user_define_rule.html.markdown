---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_client_user_define_rule"
description: |-
  Provides a Alicloud Threat Detection Client User Define Rule resource.
---

# alicloud_threat_detection_client_user_define_rule

Provides a Threat Detection Client User Define Rule resource. Malicious Behavior Defense Custom Rules.

For information about Threat Detection Client User Define Rule and how to use it, see [What is Client User Define Rule](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-addclientuserdefinerule).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_client_user_define_rule&exampleId=31182f3c-0096-591c-e624-ce1d34823f18273d36cb&activeTab=example&spm=docs.r.threat_detection_client_user_define_rule.0.31182f3c00&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_threat_detection_client_user_define_rule" "default" {
  action_type                  = "0"
  platform                     = "windows"
  registry_content             = "123"
  client_user_define_rule_name = var.name

  parent_proc_path = "/root/bash"
  type             = "5"
  cmdline          = "bash"
  proc_path        = "/root/bash"
  parent_cmdline   = "bash"
  registry_key     = "123"
}
```

## Argument Reference

The following arguments are supported:
* `action_type` - (Required) The operation type. Value:
  - **0**: plus White
  - **1**: Plus Black.
* `client_user_define_rule_name` - (Required) The custom rule name.
* `cmdline` - (Optional) Command line. When the value of the Type attribute is 2, 3, 4, 5, 6, or 7, the command line field is required.
* `file_path` - (Optional) The file path. When the value of the Type attribute is 4 or 6, 7, the FilePath field is required.
* `hash` - (Optional) Process hash list. When the value of the Type attribute is 1, the Hash attribute is required.
* `ip` - (Optional) IP address. When the value of the Type attribute is 3, the Ip attribute is required.
* `new_file_path` - (Optional) The new file path to rename the file. When the value of the Type attribute is 7, the NewFilePath attribute is required.
* `parent_cmdline` - (Optional) The parent command line.
* `parent_proc_path` - (Optional) Parent process path.
* `platform` - (Required) The operating system type. Value:
  - **windows**:widows
  - **linux**:linux
  - **all**: all.
* `port_str` - (Optional, Computed) The port number. When the value of the Type attribute is 3, the PortStr attribute is required. Value range: **1-65535**.
* `proc_path` - (Optional) The process path. When the Type attribute is set to 2, 3, 4, 5, 6, or 7, the ProcPath attribute is required.
* `registry_content` - (Optional) The registry value. When the value of the Type attribute is 5, the RegistryKey attribute is required.
* `registry_key` - (Optional) The registry key. When the value of the Type attribute is 5, the RegistryKey attribute is required.
* `type` - (Required, ForceNew) The rule type. Value:
  - **1**: Process hash
  - **2**: command line
  - **3**: Process network
  - **4**: File reading and writing
  - **5**: Operate the registry
  - **6**: Load Dynamic Link Library
  - **7**: File Rename.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Client User Define Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Client User Define Rule.
* `update` - (Defaults to 5 mins) Used when update the Client User Define Rule.

## Import

Threat Detection Client User Define Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_client_user_define_rule.example <id>
```