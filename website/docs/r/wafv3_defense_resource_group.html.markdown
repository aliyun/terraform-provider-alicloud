---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_wafv3_defense_resource_group"
description: |-
  Provides a Alicloud WAFV3 Defense Resource Group resource.
---

# alicloud_wafv3_defense_resource_group

Provides a WAFV3 Defense Resource Group resource.



For information about WAFV3 Defense Resource Group and how to use it, see [What is Defense Resource Group](https://next.api.alibabacloud.com/document/waf-openapi/2021-10-01/CreateDefenseResourceGroup).

-> **NOTE:** Available since v1.263.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_wafv3_defense_resource_group&exampleId=423bc80f-43b8-55d2-f69a-1466afc0734eee1713ac&activeTab=example&spm=docs.r.wafv3_defense_resource_group.0.423bc80f43&intl_lang=EN_US" target="_blank">
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

variable "region_id" {
  default = "cn-hangzhou"
}

resource "alicloud_wafv3_instance" "defaultHaF1fD" {
}

resource "alicloud_wafv3_domain" "defaultHVcskT" {
  instance_id = alicloud_wafv3_instance.defaultHaF1fD.id
  listen {
    http_ports = ["80"]
  }
  redirect {
    backends    = ["6.36.36.36"]
    loadbalance = "iphash"
  }
  domain      = "1511928242963727_1.wafqax.top"
  access_type = "share"
}

resource "alicloud_wafv3_domain" "defaultEH4CwO" {
  instance_id = alicloud_wafv3_instance.defaultHaF1fD.id
  listen {
    http_ports = ["80"]
  }
  redirect {
    backends    = ["6.36.36.36"]
    loadbalance = "iphash"
  }
  domain      = "1511928242963727_2.wafqax.top"
  access_type = "share"
}

resource "alicloud_wafv3_domain" "defaultY0ge1N" {
  instance_id = alicloud_wafv3_instance.defaultHaF1fD.id
  listen {
    http_ports = ["80"]
  }
  redirect {
    backends    = ["6.36.36.36"]
    loadbalance = "iphash"
  }
  domain      = "1511928242963727_3.wafqax.top"
  access_type = "share"
}


resource "alicloud_wafv3_defense_resource_group" "default" {
  group_name    = "examplefromTF"
  resource_list = ["${alicloud_wafv3_domain.defaultHVcskT.domain_id}"]
  description   = "example"
  instance_id   = alicloud_wafv3_instance.defaultHaF1fD.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_wafv3_defense_resource_group&spm=docs.r.wafv3_defense_resource_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the protected object group. 
* `group_name` - (Required, ForceNew) The name of the protected object group. The name must be 1 to 255 characters long and can contain Chinese characters, letters, digits, underscores (_), periods (.), and hyphens (-)
* `instance_id` - (Required, ForceNew) The ID of the WAF instance.
* `resource_list` - (Optional, Set) The names of the protected objects that are added to the protected object group.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<group_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Defense Resource Group.
* `delete` - (Defaults to 5 mins) Used when delete the Defense Resource Group.
* `update` - (Defaults to 5 mins) Used when update the Defense Resource Group.

## Import

WAFV3 Defense Resource Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_wafv3_defense_resource_group.example <instance_id>:<group_name>
```