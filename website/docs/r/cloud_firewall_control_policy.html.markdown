---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_control_policy"
sidebar_current: "docs-alicloud-resource-cloud-firewall-control-policy"
description: |-
  Provides a Alicloud Cloud Firewall Control Policy resource.
---

# alicloud_cloud_firewall_control_policy

Provides a Cloud Firewall Control Policy resource.

For information about Cloud Firewall Control Policy and how to use it, see [What is Control Policy](https://www.alibabacloud.com/help/doc-detail/138867.htm).

-> **NOTE:** Available since v1.129.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_control_policy&exampleId=9466770b-c784-5cd7-73b1-db97948c811470394f57&activeTab=example&spm=docs.r.cloud_firewall_control_policy.0.9466770bc7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cloud_firewall_control_policy" "default" {
  direction        = "in"
  application_name = "ANY"
  description      = var.name
  acl_action       = "accept"
  source           = "127.0.0.1/32"
  source_type      = "net"
  destination      = "127.0.0.2/32"
  destination_type = "net"
  proto            = "ANY"
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, ForceNew) The direction of the traffic to which the access control policy applies. Valid values: `in`, `out`.
* `description` - (Required) The description of the access control policy.
* `acl_action` - (Required) The action that Cloud Firewall performs on the traffic. Valid values: `accept`, `drop`, `log`.
* `source` - (Required) The source address in the access control policy.
* `source_type` - (Required) The type of the source address in the access control policy. Valid values: `net`, `group`, `location`.
* `destination` - (Required) The destination address in the access control policy.
* `destination_type` - (Required) The type of the destination address in the access control policy. Valid values: `net`, `group`, `domain`, `location`.
* `proto` - (Required) The protocol type supported by the access control policy. Valid values: `ANY`, ` TCP`, `UDP`, `ICMP`.
* `application_name` - (Optional) The application type supported by the access control policy. Valid values: `ANY`, `HTTP`, `HTTPS`, `MQTT`, `Memcache`, `MongoDB`, `MySQL`, `RDP`, `Redis`, `SMTP`, `SMTPS`, `SSH`, `SSL`, `VNC`.
-> **NOTE:** If `proto` is set to `TCP`, you can set `application_name` to any valid value. If `proto` is set to `UDP`, `ICMP`, or `ANY`, you can only set `application_name` to `ANY`.
* `dest_port` - (Optional) The destination port in the access control policy. **Note:** If `dest_port_type` is set to `port`, you must specify `dest_port`.
* `dest_port_group` - (Optional) The name of the destination port address book in the access control policy. **Note:** If `dest_port_type` is set to `group`, you must specify `dest_port_group`.
* `dest_port_type` - (Optional) The type of the destination port in the access control policy. Valid values: `port`, `group`.
* `ip_version` - (Optional, ForceNew) The IP version supported by the access control policy. Default value: `4`. Valid values:
  - `4`: IPv4.
  - `6`: IPv6.
* `domain_resolve_type` - (Optional, Available since v1.232.0) The domain name resolution method of the access control policy. Valid values:
  - `FQDN`: Fully qualified domain name (FQDN)-based resolution.
  - `DNS`: DNS-based dynamic resolution.
  - `FQDN_AND_DNS`: FQDN and DNS-based dynamic resolution.
* `repeat_type` - (Optional, Available since v1.232.0) The recurrence type for the access control policy to take effect. Default value: `Permanent`. Valid values:
  - `Permanent`: The policy always takes effect.
  - `None`: The policy takes effect for only once.
  - `Daily`: The policy takes effect on a daily basis.
  - `Weekly`: The policy takes effect on a weekly basis.
  - `Monthly`: The policy takes effect on a monthly basis.
* `start_time` - (Optional, Int, Available since v1.232.0) The time when the access control policy starts to take effect. The value is a UNIX timestamp. Unit: seconds. The value must be on the hour or on the half hour, and at least 30 minutes earlier than the end time.
* `end_time` - (Optional, Int, Available since v1.232.0) The time when the access control policy stops taking effect. The value is a UNIX timestamp. Unit: seconds. The value must be on the hour or on the half hour, and at least 30 minutes later than the start time.
-> **NOTE:** If `repeat_type` is set to `None`, `Daily`, `Weekly`, or `Monthly`, `start_time` and `end_time` must be set.
* `repeat_start_time` - (Optional, Available since v1.232.0) The point in time when the recurrence starts. Example: `08:00`. The start time must be on the hour or on the half hour, and at least 30 minutes earlier than the end time.
* `repeat_end_time` - (Optional, Available since v1.232.0) The point in time when the recurrence ends. Example: `23:30`. The end time must be on the hour or on the half hour, and at least 30 minutes later than the start time.
-> **NOTE:** If `repeat_type` is set to `Daily`, `Weekly`, or `Monthly`, `repeat_start_time` and `repeat_end_time` must be set.
* `repeat_days` - (Optional, List, Available since v1.232.0) The days of a week or of a month on which the access control policy takes effect. Valid values:
  - If `repeat_type` is set to `Weekly`. Valid values: `0` to `6`.
  - If `repeat_type` is set to `Monthly`. Valid values: `1` to `31`.
-> **NOTE:** If `repeat_type` is set to `Weekly`, or `Monthly`, `repeat_days` must be set.
* `application_name_list` - (Optional, List, Available since v1.232.0) The application types supported by the access control policy.
-> **NOTE:** If `proto` is set to `TCP`, you can set `application_name_list` to any valid value. If `proto` is set to `UDP`, `ICMP`, or `ANY`, you can only set `application_name_list` to `["ANY"]`. From version 1.232.0, You must specify at least one of the `application_name_list` and `application_name`. If you specify both `application_name_list` and `application_name`, only the `application_name_list` takes effect.
* `release` - (Optional) The status of the access control policy. Valid values: `true`, `false`.
* `source_ip` - (Optional) The source IP address of the request.
* `lang` - (Optional) The language of the content within the request and response. Valid values: `zh`, `en`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Control Policy. It formats as `<acl_uuid>:<direction>`.
* `acl_uuid` - (Available since v1.148.0) The unique ID of the access control policy.
* `create_time` - (Available since v1.232.0) The time when the access control policy was created.

## Timeouts

-> **NOTE:** Available since v1.232.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Control Policy.
* `update` - (Defaults to 5 mins) Used when update the Control Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Control Policy.

## Import

Cloud Firewall Control Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_control_policy.example <acl_uuid>:<direction>
```
