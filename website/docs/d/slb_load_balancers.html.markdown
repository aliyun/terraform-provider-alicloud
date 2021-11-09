---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slbs_load_balancers"
sidebar_current: "docs-alicloud-datasource-slb-load-balancers"
description: |-
    Provides a list of server load balancers to the user.
---

# alicloud\_slb\_load\_balancers

This data source provides the server load balancers of the current Alibaba Cloud user.

-> **NOTE:** Available in 1.123.1+

## Example Usage

```terraform
data "alicloud_slb_load_balancers" "example" {
  name_regex = "sample_slb"
  tags = {
    tagKey1 = "tagValue1",
    tagKey2 = "tagValue2"
  }
}

output "first_slb_id" {
  value = data.alicloud_slb_load_balancers.example.balancers[0].id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of SLBs IDs.
* `name_regex` - (Optional) A regex string to filter results by SLB name.
* `network_type` - (Optional) Network type of the SLBs. Valid values: `vpc` and `classic`.
* `vpc_id` - (Optional) ID of the VPC linked to the SLBs.
* `vswitch_id` - (Optional) ID of the VSwitch linked to the SLBs.
* `address` - (Optional) Service address of the SLBs.
* `tags` - (Optional) A map of tags assigned to the SLB instances. The `tags` can have a maximum of 5 tag. It must be in the format:
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The Id of resource group which SLB belongs.
* `address_ip_version` - (Optional, ForceNew) The address ip version. Valid values `ipv4` and `ipv6`.
* `address_type` - (Optional, ForceNew) The address type of the SLB. Valid values `internet` and `intranet`.
* `internet_charge_type` - (Optional, ForceNew) The internet charge type. Valid values `PayByBandwidth` and `PayByTraffic`.
* `load_balancer_name` - (Optional, ForceNew) The name of the SLB.
* `payment_type` - (Optional, ForceNew) The payment type of SLB. Valid values `PayAsYouGo` and `Subscription`.
* `server_id` - (Optional, ForceNew) The server ID.
* `server_intranet_address` - (Optional, ForceNew) The server intranet address.
* `master_zone_id` - (Optional, ForceNew) The master zone id of the SLB.
* `slave_zone_id` - (Optional, ForceNew) The slave zone id of the SLB.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of slb IDs.
* `names` - A list of slb names.
* `balancers` - A list of SLBs. Each element contains the following attributes:
    * `id` - ID of the SLB.
    * `address` - The IP address that the SLB instance uses to provide services.
    * `address_ip_version` - The address ip version.
    * `address_type` - The address type.
    * `auto_release_time` - The auto release time.
    * `backend_servers` - The backend servers of the SLB.
        * `description` - The description of servers.
        * `server_id` - The ID of the Elastic Compute Service (ECS) instance that is specified as a backend server of the CLB instance.
        * `type` - The type of servers.
        * `weight` - The weight of servers.
    * `bandwidth` - The bandwidth of the SLB.
    * `create_time_stamp` - The create time stamp of the SLB.
    * `delete_protection` - Whether the SLB should delete protection.
    * `end_time` - The end time of the SLB.
    * `end_time_stamp` - The end time stamp of the SLB.
    * `internet_charge_type` - The billing method of the Internet-facing SLB instance.
    * `listener_ports_and_protocal` - The listener ports and protocal of the SLB.
        * `listener_port` - The listener port.
        * `listener_protocal` - The listener protoal.
    * `listener_ports_and_protocol` - The listener ports and protocol of the SLB.
        * `description` - The description of protocol.
        * `forward_port` - The forward port.
        * `listener_forward` - The listener forward.
        * `listener_port` - The listener port.
        * `listener_protocol` - The listener protocol.
    * `load_balancer_id` - Thd ID of the SLB.
    * `load_balancer_name` - The name of the SLB.
    * `master_zone_id` - Master availability zone of the SLBs.
    * `modification_protection_reason` - The reason of modification protection.
    * `modification_protection_status` - The status of modification protection.
    * `network_type` -  Network type of the SLB. Possible values: `vpc` and `classic`.
    * `region_id_alias` - Region ID the SLB belongs to.
    * `renewal_cyc_unit` - The renewal cyc unit of the SLB.
    * `renewal_duration` - The renewal duration of the SLB.
    * `renewal_status` - The renewal status of the SLB.
    * `resource_group_id` - The ID of the resource group.
    * `slave_zone_id` - Slave availability zone of the SLBs.
    * `load_balancer_spec` - The specification of the SLB.
    * `status` - SLB current status. Possible values: `inactive`, `active` and `locked`.
    * `tags` - The tags of the SLB.
    * `vswitch_id` - ID of the VSwitch the SLB belongs to.
    * `vpc_id` - ID of the VPC the SLB belongs to.
