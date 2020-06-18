---
subcategory: "BGP-Line Anti-DDoS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_instance"
sidebar_current: "docs-alicloud-resource-ddoscoo-instance"
description: |-
  Provides a Alicloud BGP-line Anti-DDoS Pro(Ddoscoo) Instance Resource.
---

# alicloud\_ddoscoo\_instance

BGP-Line Anti-DDoS instance resource. "Ddoscoo" is the short term of this product. See [What is Anti-DDoS Pro](https://www.alibabacloud.com/help/en/doc-detail/69319.htm).

-> **NOTE:** The product region only support cn-hangzhou.

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available in 1.37.0+ .

## Example Usage

Basic Usage

```
provider "alicloud" {
    endpoints {
        bssopenapi = "business.aliyuncs.com"
    }
}
resource "alicloud_ddoscoo_instance" "example" {
    service_partner     ="coop-line-001"    
    edition             ="coop"
    bandwidth           ="50"
    port_count          ="50"
    domain_count        ="50"
    service_bandwidth   ="200"
    base_bandwidth      ="30"
    function_version    ="0"
    period              ="1"
    normal_qps          ="3000"
    remark              ="updatemark"
}
```
## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required, ForceNew) Elastic defend bandwidth of the instance. This value must be larger than the base defend bandwidth. Valid values: 30, 60, 100, 300, 400, 500, 600. The unit is Gbps.
* `base_bandwidth` - (Required, ForceNew) Base defend bandwidth of the instance. Valid values: 30, 60, 100, 300, 400, 500, 600. The unit is Gbps.
* `domain_count` - (Required) Domain retransmission rule count of the instance. At least 50. Increase 5 per step, such as 55, 60, 65.
* `edition` - (Required, Available in v1.87.0+) The instance version. Valid values: coop.
* `function_version` - (Required, Available in v1.87.0+) The function plan of the instance. Valid values:
    `0`: standard function plan.
    `1`: enhanced function plan.
* `modify_type` - (Optional, Available in v1.87.0+) Variant type. Valid values: Upgrade, Downgrade.
* `normal_qps` - (Required, Available in v1.87.0+) Business QPS. Value range is 3000 to 10000, step is 100, default value is 3000.
* `period` - (Optional, ForceNew) The duration that you will buy Ddoscoo instance (in month). Valid values: [1~9], 12, 24, 36. Default to 1. At present, the provider does not support modify "period".
* `port_count` - (Required) Port retransmission rule count of the instance. At least 50. Increase 5 per step, such as 55, 60, 65.
* `remark` - (Optional, Available in v1.87.0+)  The description that you want to use for the instance.
* `renew_period` - (Optional, ForceNew, Available in v1.87.0+) Automatic renewal cycle, the unit is month. When setting `RenewalStatus` to `AutoRenewal`, it must be set.
* `renewal_status` - (Optional, ForceNew, Available in v1.87.0+) Automatic renewal status. Valid values:
    `AutoRenewal`: automatic renewal. 
    `ManualRenewal`: manual renewal. 
* `service_bandwidth` - (Required) Business bandwidth of the instance. At leaset 100. Increased 100 per step, such as 100, 200, 300. The unit is Mbps.
* `service_partner` - (Required, ForceNew, Available in v1.87.0+) Line resources. Corresponding value and description: coop-line-001.
* `tags` - (Optional, Available in v1.87.0+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Ddoscoo.
* `status` - The expiration status of the instance. Valid values:
    `1`: The instance works properly.
    `2`: The instance expires.
* `instance_spec` - The specifications of the instance.
    * `instance_id` - The ID of the instance.
    * `port_limit` - The number of ports that the instance can protect.
    * `qps_limit` - The queries per second (QPS) of services.
    * `site_limit` - The number of websites that the instance can protect.
    * `defense_count` - The number of available advanced mitigation sessions for this month.-1 indicates there is no limit on the number of available advanced mitigation sessions. That is, the instance uses the Unlimited mitigation plan.
    * `domain_limit` - The number of domain names that the instance can protect.
    * `elastic_bandwidth` - The brustable protection bandwidth of the instance. Unit: Gbit/s.
    * `function_version` - The function plan of the instance. Valid values:
        `default`: standard function plan.
        `enhance`: enhanced function plan.
    * `bandwidth_mbps` - The clean bandwidth of the instance. Unit: Mbit/s.
    * `base_bandwidth` - The basic protection bandwidth of the instance. Unit: Gbit/s.

## Import

Ddoscoo instance can be imported using the id, e.g.

```
$ terraform import alicloud_ddoscoo_instance.example ddoscoo-cn-123456
```
