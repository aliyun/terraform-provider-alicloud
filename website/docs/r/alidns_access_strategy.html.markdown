---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_access_strategy"
sidebar_current: "docs-alicloud-resource-alidns-access-strategy"
description: |-
  Provides a Alicloud DNS Access Strategy resource.
---

# alicloud\_alidns\_access\_strategy

Provides a DNS Access Strategy resource.

For information about DNS Access Strategy and how to use it, see [What is Access Strategy](https://www.alibabacloud.com/help/doc-detail/189620.html).

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

data "alicloud_alidns_gtm_instances" "default" {}

resource "alicloud_alidns_gtm_instance" "default" {
  count                   = length(data.alicloud_alidns_gtm_instances.default.ids) > 0 ? 0 : 1
  instance_name           = var.name
  payment_type            = "Subscription"
  period                  = 1
  renewal_status          = "ManualRenewal"
  package_edition         = "ultimate"
  health_check_task_count = 100
  sms_notification_count  = 1000
  public_cname_mode       = "SYSTEM_ASSIGN"
  ttl                     = 60
  cname_type              = "PUBLIC"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  alert_group             = [alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
  public_user_domain_name = var.domain_name
  alert_config {
    sms_notice      = true
    notice_type     = "ADDR_ALERT"
    email_notice    = true
    dingtalk_notice = true
  }
}

locals {
  gtm_instance_id = length(data.alicloud_alidns_gtm_instances.default.ids) > 0 ? data.alicloud_alidns_gtm_instances.default.ids[0] : concat(alicloud_alidns_gtm_instance.default.*.id, [""])[0]
}

resource "alicloud_alidns_address_pool" "ipv4" {
  count             = 2
  address_pool_name = var.name
  instance_id       = local.gtm_instance_id
  lba_strategy      = "RATIO"
  type              = "IPV4"
  address {
    attribute_info = "{\"lineCodeRectifyType\":\"RECTIFIED\",\"lineCodes\":[\"os_namerica_us\"]}"
    remark         = "address_remark"
    address        = "1.1.1.1"
    mode           = "SMART"
    lba_weight     = 1
  }
}

resource "alicloud_alidns_access_strategy" "default" {
  strategy_name                  = var.name
  strategy_mode                  = "GEO"
  instance_id                    = local.gtm_instance_id
  default_addr_pool_type         = "IPV4"
  default_lba_strategy           = "RATIO"
  default_min_available_addr_num = 1
  default_addr_pools {
    lba_weight   = 1
    addr_pool_id = alicloud_alidns_address_pool.ipv4.0.id
  }
  failover_addr_pool_type         = "IPV4"
  failover_lba_strategy           = "RATIO"
  failover_min_available_addr_num = 1
  failover_addr_pools {
    lba_weight   = 1
    addr_pool_id = alicloud_alidns_address_pool.ipv4.1.id
  }
  lines {
    line_code = "default"
  }
}

```
## Argument Reference

The following arguments are supported:

* `access_mode` - (Optional, Computed) The primary/secondary switchover policy for address pool groups. Valid values: `AUTO`, `DEFAULT`, `FAILOVER`.
* `default_addr_pool_type` - (Required) The type of the primary address pool. Valid values: `IPV4`, `IPV6`, `DOMAIN`.
* `default_addr_pools` - (Required) List of primary address pool collections. See the following `Block default_addr_pools`.
* `default_latency_optimization` - (Optional) Specifies whether to enable scheduling optimization for latency resolution for the primary address pool group. Valid values: `OPEN`, `CLOSE`.
* `default_lba_strategy` - (Optional) The load balancing policy of the primary address pool group. Valid values: `ALL_RR`, `RATIO`. **NOTE:** The `default_lba_strategy` is required under the condition that `strategy_mode` is `GEO`.
* `default_max_return_addr_num` - (Optional) The maximum number of addresses returned by the primary address pool set. **NOTE:** The `default_max_return_addr_num` is required under the condition that `strategy_mode` is `LATENCY`.
* `default_min_available_addr_num` - (Required) The minimum number of available addresses for the primary address pool set.
* `failover_addr_pool_type` - (Optional) The type of the secondary address pool. Valid values: `IPV4`, `IPV6`, `DOMAIN`.
* `failover_addr_pools` - (Optional) List of backup address pool sets. See the following `Block failover_addr_pools`.
* `failover_latency_optimization` - (Optional) Specifies whether to enable scheduling optimization for latency resolution for the secondary address pool group. Valid values: `OPEN`, `CLOSE`.
* `failover_lba_strategy` - (Optional) The load balancing policy of the secondary address pool group. Valid values: `ALL_RR`, `RATIO`.
* `failover_max_return_addr_num` - (Optional) The maximum number of returned addresses in the standby address pool.
* `failover_min_available_addr_num` - (Optional) The minimum number of available addresses in the standby address pool.
* `instance_id` - (Required, ForceNew) The Id of the associated instance.
* `lang` - (Optional) The lang.
* `strategy_mode` - (Required) The type of the access policy. Valid values: `GEO` or `LATENCY`. `GEO`: based on geographic location. `LATENCY`: Based on delay.
* `strategy_name` - (Required) The name of the access policy.
* `lines` - (Optional) The source regions. See the following `Block lines`. **NOTE:** The `lines` is required under the condition that `strategy_mode` is `GEO`.

#### Block failover_addr_pools

The failover_addr_pools supports the following: 

* `addr_pool_id` - (Optional) The ID of the address pool in the secondary address pool group.
* `lba_weight` - (Optional) The weight of the address pool in the secondary address pool group.

#### Block default_addr_pools

The default_addr_pools supports the following: 

* `addr_pool_id` - (Required) The ID of the address pool in the primary address pool group.
* `lba_weight` - (Optional) The weight of the address pool in the primary address pool group.

#### Block lines

The lines supports the following: 

* `line_code` - (Optional) The line code of the source region.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Access Strategy.

## Import

DNS Access Strategy can be imported using the id, e.g.

```
$ terraform import alicloud_alidns_access_strategy.example <id>
```