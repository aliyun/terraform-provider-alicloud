---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_firewall_control_policies"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-vpc-firewall-control-policies"
description: |-
  Provides a list of Cloud Firewall Vpc Firewall Control Policies to the user.
---

# alicloud_cloud_firewall_vpc_firewall_control_policies

This data source provides the Cloud Firewall Vpc Firewall Control Policies of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_firewall_vpc_firewall_control_policies" "ids" {
  vpc_firewall_id = "example_value"
  ids             = ["example_value-1", "example_value-2"]
}
output "alicloud_cloud_firewall_vpc_firewall_control_policies_id_1" {
  value = data.alicloud_cloud_firewall_vpc_firewall_control_policies.ids.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `acl_action` - (Optional, ForceNew) The action that Cloud Firewall performs on the traffic. Valid values: `accept`, `drop`, `log`.
* `acl_uuid` - (Optional, ForceNew) Access control over VPC firewalls strategy unique identifier.
* `description` - (Optional, ForceNew) Access control over VPC firewalls description of the strategy information.
* `destination` - (Optional, ForceNew) Access control over VPC firewalls strategy the destination address in.
* `lang` - (Optional, ForceNew) The language of the content within the request and response. Valid values: `zh`, `en`.
* `member_uid` - (Optional, ForceNew) The UID of the member account of the current Alibaba cloud account.
* `proto` - (Optional, ForceNew) Access control over VPC firewalls strategy access traffic of the protocol type.
* `release` - (Optional, ForceNew) The enabled status of the access control policy. The policy is enabled by default after it is created. Value:
  - **true**: Enable access control policies
  - **false**: does not enable access control policies.
* `source` - (Optional, ForceNew) Access control over VPC firewalls strategy in the source address.
* `vpc_firewall_id` - (Required, ForceNew) The ID of the VPC firewall instance. Value:
  - When the VPC firewall protects traffic between two VPCs connected through the cloud enterprise network, the policy group ID uses the cloud enterprise network instance ID.
  - When the VPC firewall protects traffic between two VPCs connected through the express connection, the policy group ID uses the ID of the VPC firewall instance.
* `ids` - (Optional, ForceNew, Computed)  A list of Vpc Firewall Control Policy IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `policies` - A list of Cloud Firewall Vpc Firewall Control Policies. Each element contains the following attributes:
  * `acl_action` - Access control over VPC firewalls are set in the access traffic via Alibaba cloud firewall way (ACT).
  * `acl_uuid` - Access control over VPC firewalls strategy unique identifier.
  * `application_id` - Policy specifies the application ID.
  * `application_name` - Access control over VPC firewalls policies support the application types.
  * `description` - Access control over VPC firewalls description of the strategy information.
  * `dest_port` - Access control over VPC firewalls strategy access traffic of the destination port.
  * `dest_port_group` - Access control policy in the access traffic of the destination port address book name.
  * `dest_port_group_ports` - Port Address Book port list.
  * `dest_port_type` - Access control over VPC firewalls strategy access traffic of the destination port type.
  * `destination` - Access control over VPC firewalls strategy the destination address in.
  * `destination_group_cidrs` - Destination address book defined in the address list.
  * `destination_group_type` - The destination address book type in the access control policy. Value: `ip`, `domain`.
  * `destination_type` - Access control over VPC firewalls strategy in the destination address of the type.
  * `hit_times` - Control strategy of hits per second.
  * `member_uid` - The UID of the member account of the current Alibaba cloud account.
  * `order` - Access control over VPC firewalls policies will go into effect of priority. The priority value starts from 1, the smaller the priority number, the higher the priority. -1 represents the lowest priority.
  * `proto` - Access control over VPC firewalls strategy access traffic of the protocol type.
  * `release` - The enabled status of the access control policy. The policy is enabled by default after it is created. Value:
  * `source` - Access control over VPC firewalls strategy in the source address.
  * `source_group_cidrs` - SOURCE address of the address list.
  * `source_group_type` - The source address type in the access control policy. Unique value: **ip**. The IP address book contains one or more IP address segments.
  * `source_type` - Access control over VPC firewalls policy source address type.
  * `vpc_firewall_id` - The ID of the VPC firewall instance.
  * `id` - The ID of the Cloud Firewall Vpc Firewall Control Policy.