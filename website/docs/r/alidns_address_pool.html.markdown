---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_address_pool"
sidebar_current: "docs-alicloud-resource-alidns-address-pool"
description: |-
  Provides a Alicloud Alidns Address Pool resource.
---

# alicloud_alidns_address_pool

Provides a Alidns Address Pool resource.

For information about Alidns Address Pool and how to use it, see [What is Address Pool](https://www.alibabacloud.com/help/doc-detail/189621.html).

-> **NOTE:** Available since v1.152.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alidns_address_pool&exampleId=66c68e52-2208-5c71-a3af-b408d85f87661eefaeea&activeTab=example&spm=docs.r.alidns_address_pool.0.66c68e5222&intl_lang=EN_US" target="_blank">
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
  address_pool_name = var.name
  instance_id       = alicloud_alidns_gtm_instance.default.id
  lba_strategy      = "RATIO"
  type              = "IPV4"
  address {
    attribute_info = <<EOF
    {
      "lineCodeRectifyType": "RECTIFIED",
      "lineCodes": ["os_namerica_us"]
    }
    EOF
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
* `address` - (Required) The address lists of the Address Pool. See [`address`](#address) below for details.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `lba_strategy` - (Required)The load balancing policy of the address pool. Valid values:`ALL_RR` or `RATIO`. `ALL_RR`: returns all addresses. `RATIO`: returns addresses by weight.
* `type` - (Required, ForceNew) The type of the address pool. Valid values: `IPV4`, `IPV6`, `DOMAIN`.

### `address`

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

```shell
$ terraform import alicloud_alidns_address_pool.example <id>
```