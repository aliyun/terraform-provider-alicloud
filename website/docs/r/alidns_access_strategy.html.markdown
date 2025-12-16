---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_access_strategy"
sidebar_current: "docs-alicloud-resource-alidns-access-strategy"
description: |-
  Provides a Alicloud DNS Access Strategy resource.
---

# alicloud_alidns_access_strategy

Provides a DNS Access Strategy resource.

For information about DNS Access Strategy and how to use it, see [What is Access Strategy](https://www.alibabacloud.com/help/doc-detail/189620.html).

-> **NOTE:** Available since v1.152.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alidns_access_strategy&exampleId=bbd1dfce-b2d7-19ca-d7bf-8d5959cf650a54ca5189&activeTab=example&spm=docs.r.alidns_access_strategy.0.bbd1dfceb2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
variable "domain_name" {
  default = "alicloud-provider.com"
}
data "alicloud_resource_manager_resource_groups" "default" {}
resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

resource "alicloud_alidns_gtm_instance" "default" {
  instance_name           = var.name
  payment_type            = "Subscription"
  period                  = 1
  renewal_status          = "ManualRenewal"
  package_edition         = "standard"
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

resource "alicloud_alidns_address_pool" "default" {
  count             = 2
  address_pool_name = format("${var.name}_%d", count.index + 1)
  instance_id       = alicloud_alidns_gtm_instance.default.id
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
  instance_id                    = alicloud_alidns_gtm_instance.default.id
  default_addr_pool_type         = "IPV4"
  default_lba_strategy           = "RATIO"
  default_min_available_addr_num = 1
  default_addr_pools {
    lba_weight   = 1
    addr_pool_id = alicloud_alidns_address_pool.default.0.id
  }
  failover_addr_pool_type         = "IPV4"
  failover_lba_strategy           = "RATIO"
  failover_min_available_addr_num = 1
  failover_addr_pools {
    lba_weight   = 1
    addr_pool_id = alicloud_alidns_address_pool.default.1.id
  }
  lines {
    line_code = "default"
  }
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alidns_access_strategy&spm=docs.r.alidns_access_strategy.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `access_mode` - (Optional) The primary/secondary switchover policy for address pool groups. Valid values: `AUTO`, `DEFAULT`, `FAILOVER`.
* `default_addr_pool_type` - (Required) The type of the primary address pool. Valid values: `IPV4`, `IPV6`, `DOMAIN`.
* `default_addr_pools` - (Required) List of primary address pool collections. See [`default_addr_pools`](#default_addr_pools) below for details.
* `default_latency_optimization` - (Optional) Specifies whether to enable scheduling optimization for latency resolution for the primary address pool group. Valid values: `OPEN`, `CLOSE`.
* `default_lba_strategy` - (Optional) The load balancing policy of the primary address pool group. Valid values: `ALL_RR`, `RATIO`. **NOTE:** The `default_lba_strategy` is required under the condition that `strategy_mode` is `GEO`.
* `default_max_return_addr_num` - (Optional) The maximum number of addresses returned by the primary address pool set. **NOTE:** The `default_max_return_addr_num` is required under the condition that `strategy_mode` is `LATENCY`.
* `default_min_available_addr_num` - (Required) The minimum number of available addresses for the primary address pool set.
* `failover_addr_pool_type` - (Optional) The type of the secondary address pool. Valid values: `IPV4`, `IPV6`, `DOMAIN`.
* `failover_addr_pools` - (Optional) List of backup address pool sets. See [`failover_addr_pools`](#failover_addr_pools) below for details.
* `failover_latency_optimization` - (Optional) Specifies whether to enable scheduling optimization for latency resolution for the secondary address pool group. Valid values: `OPEN`, `CLOSE`.
* `failover_lba_strategy` - (Optional) The load balancing policy of the secondary address pool group. Valid values: `ALL_RR`, `RATIO`.
* `failover_max_return_addr_num` - (Optional) The maximum number of returned addresses in the standby address pool.
* `failover_min_available_addr_num` - (Optional) The minimum number of available addresses in the standby address pool.
* `instance_id` - (Required, ForceNew) The Id of the associated instance.
* `lang` - (Optional) The lang.
* `strategy_mode` - (Required) The type of the access policy. Valid values: `GEO` or `LATENCY`. `GEO`: based on geographic location. `LATENCY`: Based on delay.
* `strategy_name` - (Required) The name of the access policy.
* `lines` - (Optional) The source regions. See [`lines`](#lines) below for details. **NOTE:** The `lines` is required under the condition that `strategy_mode` is `GEO`.

### `failover_addr_pools`

The failover_addr_pools supports the following: 

* `addr_pool_id` - (Optional) The ID of the address pool in the secondary address pool group.
* `lba_weight` - (Optional) The weight of the address pool in the secondary address pool group.

### `default_addr_pools`

The default_addr_pools supports the following: 

* `addr_pool_id` - (Required) The ID of the address pool in the primary address pool group.
* `lba_weight` - (Optional) The weight of the address pool in the primary address pool group.

### `lines`

The lines supports the following: 

* `line_code` - (Optional) The line code of the source region.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Access Strategy.

## Import

DNS Access Strategy can be imported using the id, e.g.

```shell
$ terraform import alicloud_alidns_access_strategy.example <id>
```