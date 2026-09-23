---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_firewall_acl_engine_mode"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Firewall Acl Engine Mode resource.
---

# alicloud_cloud_firewall_vpc_firewall_acl_engine_mode

Provides a Cloud Firewall Vpc Firewall Acl Engine Mode resource.

VPC boundary firewall engine mode.

For information about Cloud Firewall Vpc Firewall Acl Engine Mode and how to use it, see [What is Vpc Firewall Acl Engine Mode](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/ModifyVpcFirewallAclEngineMode).

-> **NOTE:** Available since v1.269.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_vpc_firewall_acl_engine_mode&exampleId=3d937d60-52eb-324c-219e-8316bd94d4f455f39947&activeTab=example&spm=docs.r.cloud_firewall_vpc_firewall_acl_engine_mode.0.3d937d6052&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}

resource "alicloud_cen_instance" "cen" {
  description       = "yqc-example001"
  cen_instance_name = "yqc-example-CenInstance001"
}

resource "alicloud_cen_transit_router" "TR" {
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_vpc" "vpc1" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-example-001"
}

resource "alicloud_vswitch" "vpc1vsw1" {
  vpc_id     = alicloud_vpc.vpc1.id
  zone_id    = "cn-hangzhou-h"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id     = alicloud_vpc.vpc1.id
  zone_id    = "cn-hangzhou-i"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  vpc_id = alicloud_vpc.vpc1.id
  cen_id = alicloud_cen_instance.cen.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = alicloud_vswitch.vpc1vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
    zone_id    = alicloud_vswitch.vpc1vsw2.zone_id
  }
  transit_router_vpc_attachment_name    = "example"
  transit_router_attachment_description = "111"
  auto_publish_route_enabled            = true
  transit_router_id                     = alicloud_cen_transit_router.TR.transit_router_id
}


resource "alicloud_cloud_firewall_vpc_firewall_acl_engine_mode" "default" {
  strict_mode     = "0"
  vpc_firewall_id = alicloud_cen_instance.cen.id
  member_uid      = "1511928242963727"
}
```

### Deleting `alicloud_cloud_firewall_vpc_firewall_acl_engine_mode` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_firewall_vpc_firewall_acl_engine_mode`. Terraform will remove this resource from the state file, however resources may remain.


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_firewall_vpc_firewall_acl_engine_mode&spm=docs.r.cloud_firewall_vpc_firewall_acl_engine_mode.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `strict_mode` - (Required, Int) The mode of the ACL engine. Possible values are `0`, `1`.
* `vpc_firewall_id` - (Required, ForceNew) The ID of the VPC firewall.
* `member_uid` - (Optional, ForceNew) The ID of member account.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Firewall Acl Engine Mode.
* `update` - (Defaults to 5 mins) Used when update the Vpc Firewall Acl Engine Mode.

## Import

Cloud Firewall Vpc Firewall Acl Engine Mode can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_vpc_firewall_acl_engine_mode.example <vpc_firewall_id>
```