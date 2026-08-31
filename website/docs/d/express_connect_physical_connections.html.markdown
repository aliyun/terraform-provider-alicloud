---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_physical_connections"
sidebar_current: "docs-alicloud-datasource-express-connect-physical-connections"
description: |-
  Provides a list of Express Connect Physical Connections to the user.
---

# alicloud\_express\_connect\_physical\_connections

This data source provides the Express Connect Physical Connections of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_express_connect_physical_connections" "ids" {
  ids = ["pc-2345678"]
}
output "express_connect_physical_connection_id_1" {
  value = data.alicloud_express_connect_physical_connections.ids.connections.0.id
}

data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^my-PhysicalConnection"
}
output "express_connect_physical_connection_id_2" {
  value = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Physical Connection IDs.
* `include_reservation_data` - (Optional, ForceNew) The include reservation data.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Physical Connection name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) Resources on Behalf of a State of the Resource Attribute Field. Valid values: `Canceled`, `Enabled`, `Terminated`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Physical Connection names.
* `connections` - A list of Express Connect Physical Connections. Each element contains the following attributes:
	* `access_point_id` - The Physical Leased Line Access Point ID.
	* `ad_location` - To Connect a Device Physical Location.
	* `bandwidth` - On the Bandwidth of the ECC Service and Physical Connection.
	* `business_status` - The Physical Connection to Which the Payment Status: Normal, financiallocked, securitylocked.
	* `circuit_code` - Operators for Physical Connection Circuit Provided Coding.
	* `create_time` - The Representative of the Creation Time Resources Attribute Field.
	* `description` - The Physical Connection to Which the Description.
	* `enabled_time` - The Physical Connection to Which the Activation Time.
	* `end_time` - The Expiration Time.
	* `has_reservation_data` - HasReservationData.
	* `id` - The ID of the Physical Connection.
	* `line_operator` - Provides Access to the Physical Line Operator Value CT: China Telecom, CU: China Unicom, CM: china Mobile, CO: Other Chinese, Equinix:Equinix, Other: Other Overseas.
	* `loa_status` - Loa State.
	* `payment_type` - on Behalf of the Pay-as-You-Type of Resource Attribute Field.
	* `peer_location` - and an on-Premises Data Center Location.
	* `physical_connection_id` - on Behalf of the Resource Level Id of the Resources Property Fields.
	* `physical_connection_name` - on Behalf of the Resource Name of the Resources-Attribute Field.
	* `port_number` - To Connect a Device Port: The Port Number of.
	* `port_type` - The Physical Leased Line Access Port Type Value 100Base-T: Fast Electrical Ports, 1000Base-T (the Default): gigabit Electrical Ports, 1000Base-LX: Gigabit Singlemode Optical Ports (10Km), 10GBase-T: Gigabit Electrical Port, 10GBase-LR: Gigabit Singlemode Optical Ports (10Km).
	* `redundant_physical_connection_id` - Redundant Physical Connection to Which the ID.
	* `reservation_active_time` - The Renewal of the Entry into Force of the Time.
	* `reservation_internet_charge_type` - Renewal Type.
	* `reservation_order_type` - Renewal Order Type.
	* `spec` - The Physical Connection to Which the Specifications.
	* `status` - Resources on Behalf of a State of the Resource Attribute Field.
	* `type` - Physical Private Line of Type. Default Value: VPC.
