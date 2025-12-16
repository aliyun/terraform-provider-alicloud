---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_firewall_cen"
sidebar_current: "docs-alicloud-resource-cloud-firewall-vpc-firewall-cen"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Firewall Cen resource.
---

# alicloud_cloud_firewall_vpc_firewall_cen

Provides a Cloud Firewall Vpc Firewall Cen resource.

For information about Cloud Firewall Vpc Firewall Cen and how to use it, see [What is Vpc Firewall Cen](https://www.alibabacloud.com/help/en/cloud-firewall/latest/createvpcfirewallcenconfigure).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_vpc_firewall_cen&exampleId=d0c1e272-f582-436c-ab92-39ed385cc0f2ae94a454&activeTab=example&spm=docs.r.cloud_firewall_vpc_firewall_cen.0.d0c1e272f5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# These resource primary keys should be replaced with your actual values.
resource "alicloud_cloud_firewall_vpc_firewall_cen" "default" {
  cen_id = "cen-xxx"
  local_vpc {
    network_instance_id = "vpc-xxx"
  }
  status            = "open"
  member_uid        = "14151*****827022"
  vpc_region        = "cn-hangzhou"
  vpc_firewall_name = "tf-vpc-firewall-name"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_firewall_vpc_firewall_cen&spm=docs.r.cloud_firewall_vpc_firewall_cen.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `vpc_firewall_name` - (Required) The name of the VPC firewall instance.
* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `vpc_region` - (Required, ForceNew) The ID of the region to which the VPC is created.
* `status` - (Required) Firewall switch status.
* `member_uid` - (Optional, ForceNew) The UID of the member account (other Alibaba Cloud account) of the current Alibaba cloud account.
* `lang` - (Optional, ForceNew) The language type of the requested and received messages. Valid values:
  - `zh`: Chinese.
  - `en`: English.
* `local_vpc` - (Required, ForceNew) The details of the VPC. See [`local_vpc`](#local_vpc) below.

### `local_vpc`

The local_vpc supports the following:

* `network_instance_id` - (Required, ForceNew) The ID of the VPC instance that created the VPC firewall.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above.
* `connect_type` - Intercommunication type, value: expressconnect: Express Channel cen: Cloud Enterprise Network
* `local_vpc` - The details of the VPC.
    * `attachment_id` - The connection ID of the network instance.
    * `attachment_name` - The connection name of the network instance.
    * `defend_cidr_list` - The list of network segments protected by the VPC firewall.
    * `eni_list` - List of elastic network cards.
        * `eni_id` - The ID of the instance of the ENI in the VPC.
        * `eni_private_ip_address` - The private IP address of the ENI in the VPC.
    * `manual_vswitch_id` - The ID of the vSwitch specified when the routing mode is manual mode.
    * `network_instance_name` - The name of the network instance.
    * `network_instance_type` - The type of the network instance. Value: **VPC * *.
    * `owner_id` - The UID of the Alibaba Cloud account to which the VPC belongs.
    * `region_no` - The region ID of the VPC.
    * `route_mode` - Routing mode,. Value:-auto: indicates automatic mode.-manual: indicates manual mode.
    * `support_manual_mode` - Whether routing mode supports manual mode. Value:-**1**: Supported.-**0**: Not supported.
    * `transit_router_id` - The ID of the CEN-TR instance.
    * `transit_router_type` - The version of the cloud enterprise network forwarding router (CEN-TR). Value:-**Basic**: Basic Edition.-**Enterprise**: Enterprise Edition.
    * `vpc_cidr_table_list` - The VPC network segment list.
        * `route_entry_list` - The list of route entries in the VPC.
            * `destination_cidr` - The target network segment of the VPC.
            * `next_hop_instance_id` - The ID of the next hop instance in the VPC.
        * `route_table_id` - The ID of the route table of the VPC.
    * `vpc_id` - The ID of the VPC instance.
    * `vpc_name` - The instance name of the VPC.
* `vpc_firewall_id` - VPC firewall ID

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 31 mins) Used when create the Vpc Firewall Cen.
* `update` - (Defaults to 31 mins) Used when update the Vpc Firewall Cen.
* `delete` - (Defaults to 31 mins) Used when delete the Vpc Firewall Cen.

## Import

Cloud Firewall Vpc Firewall Cen can be imported using the id, e.g.

```shell
$terraform import alicloud_cloud_firewall_vpc_firewall_cen.example <id>
```
