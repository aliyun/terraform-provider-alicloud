---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_firewall_ips_config"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Firewall Ips Config resource.
---

# alicloud_cloud_firewall_vpc_firewall_ips_config

Provides a Cloud Firewall Vpc Firewall Ips Config resource.

IP configuration of VPC firewall.

For information about Cloud Firewall Vpc Firewall Ips Config and how to use it, see [What is Vpc Firewall Ips Config](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/ModifyVpcFirewallDefaultIPSConfig).

-> **NOTE:** Available since v1.269.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_vpc_firewall_ips_config&exampleId=96479fd1-1d8f-89fd-0062-a2182c616b889990f0ee&activeTab=example&spm=docs.r.cloud_firewall_vpc_firewall_ips_config.0.96479fd11d&intl_lang=EN_US" target="_blank">
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


resource "alicloud_cloud_firewall_vpc_firewall_ips_config" "default" {
  enable_all_patch = "0"
  basic_rules      = "0"
  run_mode         = "0"
  vpc_firewall_id  = "vfw-tr-bb81adb2d8184bc290a5"
  rule_class       = "0"
  lang             = "cn-shenzhen"
  member_uid       = "1094685339207557"
}
```

### Deleting `alicloud_cloud_firewall_vpc_firewall_ips_config` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_firewall_vpc_firewall_ips_config`. Terraform will remove this resource from the state file, however resources may remain.


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_firewall_vpc_firewall_ips_config&spm=docs.r.cloud_firewall_vpc_firewall_ips_config.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `basic_rules` - (Required, Int) Base rule switch. Value:
  - `1`: on.
  - `0`: Off.
* `enable_all_patch` - (Required, Int) Virtual patch switch. Value:
  - `1`: on.
  - `0`: Off.
* `lang` - (Optional) Language

  -> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `member_uid` - (Optional) MemberUid

  -> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `rule_class` - (Optional) IPS rule Group
* `run_mode` - (Required, Int) IPS defense mode. Value:
  - `1`: Intercept mode.
  - `0`: Observation mode.
* `vpc_firewall_id` - (Required, ForceNew) The ID of the VPC firewall instance. Value:
  - When VPC firewall protects the network instances (including VPC, VBR, and CCN) and the specified VPC, the instance ID uses the CEN instance ID. You can call the [DescribeVpcFirewallCenList](~~ 345777 ~~) operation to query the instance ID of CEN.
  - When the VPC firewall protects the traffic between two VPCs connected through the express connection, the instance ID uses the VPC firewall instance ID. You can call the [DescribeVpcFirewallList](~~ 342932 ~~) operation to query the instance ID of the VPC firewall.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Firewall Ips Config.
* `update` - (Defaults to 5 mins) Used when update the Vpc Firewall Ips Config.

## Import

Cloud Firewall Vpc Firewall Ips Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_vpc_firewall_ips_config.example <id>
```