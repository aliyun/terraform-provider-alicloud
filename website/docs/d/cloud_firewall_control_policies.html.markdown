---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_control_policies"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-control-policies"
description: |- 
  Provides a list of Cloud Firewall Control Policies to the user.
---

# alicloud\_cloud\_firewall\_control\_policies

This data source provides the Cloud Firewall Control Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```
data "alicloud_cloud_firewall_control_policies" "example" {	
	direction = "in"
}
```

## Argument Reference

The following arguments are supported:

* `acl_action` - (Optional, ForceNew) The action that Cloud Firewall performs on the traffic. Valid values: `accept`, `drop`, `log`.
* `acl_uuid` - (Optional, ForceNew) The unique ID of the access control policy.
* `description` - (Optional, ForceNew) The description of the access control policy.
* `destination` - (Optional, ForceNew) The destination address defined in the access control policy.
* `direction` - (Required, ForceNew) Direction. Valid values: `in`, `out`.
* `ip_version` - (Optional, ForceNew) The ip version.
* `lang` - (Optional, ForceNew) DestPortGroupPorts. Valid values: `en`, `zh`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `proto` - (Optional, ForceNew) The protocol type of traffic to which the access control policy applies. Valid values: If `direction` is  `in`, the valid value is `ANY`. If `direction` is `out`, the valid values are `ANY`, `TCP`, `UDP`, `ICMP`.
* `source` - (Optional, ForceNew) The source address defined in the access control policy.
* `source_ip` - (Optional, ForceNew) The source IP address of the request.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Control Policy IDs.
* `policies` - A list of Cloud Firewall Control Policies. Each element contains the following attributes:
    * `acl_action` - The action that Cloud Firewall performs on the traffic. Valid values: `accept`, `drop`, `log`.
    * `acl_uuid` - The unique ID of the access control policy.
    * `application_name` - The application type that the access control policy supports.If `direction` is `in`, the valid value is `ANY`. If `direction` is `out`, `ANY`, `HTTP`, `HTTPS`, `MQTT`, `Memcache`, `MongoDB`, `MySQL`, `RDP`, `Redis`, `SMTP`, `SMTPS`, `SSH`, `SSL`, `VNC`.
    * `description` - The description of the access control policy.
    * `dest_port` - The destination port defined in the access control policy. 
    * `dest_port_group` - The destination port address book defined in the access control policy.
    * `dest_port_type` - The destination port type defined in the access control policy. Valid values: `group`, `port`.
    * `destination` - The destination address defined in the access control policy. 
    * `destination_type` - The destination address type defined in the access control policy.Valid values: If `direction` is `in`, the valid values are `net`, `group`. If `direction` is `out`, the valid values are `net`, `group`, `domain`, `location`.
    * `direction` - The direction of traffic to which the access control policy applies. Valid values: `in`, `out`.
    * `id` - The ID of the Control Policy.
    * `proto` - The protocol type of traffic to which the access control policy applies. Valid values: If `direction` is `in`, the valid value is `ANY`. If `direction` is `out`, the valid values are `ANY`, `TCP`, `UDP`, `ICMP`.
    * `release` - Specifies whether the access control policy is enabled. By default, an access control policy is enabled after it is created. Valid values: `true`, `false`.
    * `source` - The source address defined in the access control policy.
    * `source_type` - The type of the source address book defined in the access control policy. Valid values: If `direction` is to `in`, the valid values are `net`, `group`, `location`. If `direction` is `out`, the valid values are `net`, `group`.
