---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_address_pool"
sidebar_current: "docs-alicloud-resource-alidns-address-pool"
description: |-
  Provides a Alicloud Alidns Address Pool resource.
---

# alicloud\_alidns\_address\_pool

Provides a Alidns Address Pool resource.

For information about Alidns Address Pool and how to use it, see [What is Address Pool](https://www.alibabacloud.com/help/doc-detail/189621.html).

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = "example_value"
}

data "alicloud_alidns_gtm_instances" "default" {}

resource "alicloud_alidns_gtm_instance" "default" {
  count                   = length(data.alicloud_alidns_gtm_instances.default.ids) > 0 ? 0 : 1
  instance_name           = "example_value"
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
  public_user_domain_name = "example_domain_name"
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

resource "alicloud_alidns_address_pool" "default" {
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
```

## Argument Reference

The following arguments are supported:
* `address_pool_name` - (Required) The name of the address pool.
* `address` - (Required) The address lists of the Address Pool. See the following `Block address`.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `lba_strategy` - (Required)The load balancing policy of the address pool. Valid values:`ALL_RR` or `RATIO`. `ALL_RR`: returns all addresses. `RATIO`: returns addresses by weight.
* `type` - (Required, ForceNew) The type of the address pool. Valid values: `IPV4`, `IPV6`, `DOMAIN`.

#### Block address

The address supports the following:
* `address` - (Required) The address that you want to add to the address pool.
* `attribute_info` - (Required) The source region of the address. expressed as a JSON string. The structure is as follows:
  * `LineCodes`: List of home lineCodes.
  * `lineCodeRectifyType`: The rectification type of the line code. Default value: `AUTO`. Valid values: `NO_NEED`: no need for rectification. `RECTIFIED`: rectified. `AUTO`: automatic rectification.
* `lba_weight` - (Optional) The weight of the address. **NOTE:** The attribute is valid when the attribute `lba_strategy` is `RATIO`.
* `mode` - (Required) The type of the address. Valid values:`SMART`, `ONLINE` and `OFFLINE`.
* `remark` - (Optional) The description of the address.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Address Pool.

## Import

Alidns Address Pool can be imported using the id, e.g.

```
$ terraform import alicloud_alidns_address_pool.example <id>
```