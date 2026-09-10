---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_endpoints"
sidebar_current: "docs-alicloud-datasource-pvtz-endpoints"
description: |-
  Provides a list of Pvtz Endpoints to the user.
---

# alicloud\_pvtz\_endpoints

This data source provides the Pvtz Endpoints of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_pvtz_endpoints" "ids" {
  ids = ["example_id"]
}

output "pvtz_endpoint_id_1" {
  value = data.alicloud_pvtz_endpoints.ids.endpoints.0.id
}

data "alicloud_pvtz_endpoints" "nameRegex" {
  name_regex = "^my-Endpoint"
}

output "pvtz_endpoint_id_2" {
  value = data.alicloud_pvtz_endpoints.nameRegex.endpoints.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Endpoint IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Endpoint name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `CHANGE_FAILED`, `CHANGE_INIT`, `EXCEPTION`, `FAILED`, `INIT`, `SUCCESS`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Endpoint names.
* `endpoints` - A list of Pvtz Endpoints. Each element contains the following attributes:
	* `create_time` - The creation time of the resource.
	* `endpoint_name` - The name of the resource.
	* `ip_configs` - The Ip Configs.
		* `ip` - The IP address within the parameter range of the subnet mask. **NOTE:** It is recommended to use the IP address assigned by the system.
		* `vswitch_id` - The Vswitch id.
		* `zone_id` - The Zone ID.
		* `cidr_block` - The Subnet mask.
	* `security_group_id` - The ID of the Security Group.
	* `status` - The status of the resource. Valid values: `CHANGE_FAILED`, `CHANGE_INIT`, `EXCEPTION`, `FAILED`, `INIT`, `SUCCESS`.
	* `vpc_id` - The VPC ID.
	* `vpc_name` - The name of the VPC.
	* `vpc_region_id` - The Region of the VPC.