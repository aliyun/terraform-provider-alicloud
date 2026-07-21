---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_private_dns"
description: |-
  Provides a Alicloud Cloud Firewall Private Dns resource.
---

# alicloud_cloud_firewall_private_dns

Provides a Cloud Firewall Private Dns resource.

Private DNS Endpoint.

For information about Cloud Firewall Private Dns and how to use it, see [What is Private Dns](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/CreatePrivateDnsEndpoint).

-> **NOTE:** Available since v1.264.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_private_dns&exampleId=f68cdacc-ea89-f884-d0a3-59638c21c7d491f9b290&activeTab=example&spm=docs.r.cloud_firewall_private_dns.0.f68cdaccea&intl_lang=EN_US" target="_blank">
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

data "alicloud_account" "current" {
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-example-vpc"
}

resource "alicloud_vswitch" "vpcvsw1" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = "cn-hangzhou-i"
  cidr_block = "172.16.3.0/24"
}

resource "alicloud_vswitch" "vpcvsw2" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "172.16.4.0/24"
}


resource "alicloud_cloud_firewall_private_dns" "default" {
  region_no            = "cn-hangzhou"
  access_instance_name = var.name
  port                 = "53"
  primary_vswitch_id   = alicloud_vswitch.vpcvsw1.id
  standby_dns          = "4.4.4.4"
  primary_dns          = "8.8.8.8"
  vpc_id               = alicloud_vpc.vpc.id
  private_dns_type     = "Custom"
  firewall_type        = ["internet"]
  ip_protocol          = "UDP"
  standby_vswitch_id   = alicloud_vswitch.vpcvsw2.id
  domain_name_list     = ["www.aliyun.com"]
  primary_vswitch_ip   = "172.16.3.1"
  standby_vswitch_ip   = "172.16.4.1"
  member_uid           = data.alicloud_account.current.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_firewall_private_dns&spm=docs.r.cloud_firewall_private_dns.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `access_instance_name` - (Required) The name of Private DNS instance
* `domain_name_list` - (Optional, Set) Private DNS domain name list
* `firewall_type` - (Required, ForceNew, List) The type of firewall
* `ip_protocol` - (Optional, ForceNew) IP protocol
* `member_uid` - (Optional, ForceNew, Int) The member Uid
* `port` - (Optional, ForceNew, Int) The Port of Private DNS instance
* `primary_dns` - (Optional) Primary DNS IP
* `primary_vswitch_id` - (Optional, ForceNew) Primary zone Switch ID
* `primary_vswitch_ip` - (Optional, ForceNew) Primary zone switch IP
* `private_dns_type` - (Required) The type of Private DNS instance
* `region_no` - (Required, ForceNew) The region ID of Private DNS instance
* `standby_dns` - (Optional) Standby DNS IP
* `standby_vswitch_id` - (Optional, ForceNew) Standby zone switch ID
* `standby_vswitch_ip` - (Optional, ForceNew) Standby zone switch IP address
* `vpc_id` - (Required, ForceNew) The ID of the VPC.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<access_instance_id>:<region_no>`.
* `access_instance_id` - The id of Private DNS instance
* `status` - status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Private Dns.
* `delete` - (Defaults to 5 mins) Used when delete the Private Dns.
* `update` - (Defaults to 5 mins) Used when update the Private Dns.

## Import

Cloud Firewall Private Dns can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_private_dns.example <access_instance_id>:<region_no>
```