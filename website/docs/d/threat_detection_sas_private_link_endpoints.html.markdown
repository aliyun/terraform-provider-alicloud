---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_sas_private_link_endpoints"
sidebar_current: "docs-alicloud-datasource-threat-detection-sas-private-link-endpoints"
description: |-
  Provides a list of Threat Detection Sas Private Link Endpoints to the user.
---

# alicloud_threat_detection_sas_private_link_endpoints

This data source provides a list of Threat Detection Sas Private Link Endpoints of current Alibaba Cloud user.

-> **NOTE:** Available since v1.235.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_sas_private_link_endpoints" "default" {
  ids            = ["${alicloud_threat_detection_sas_private_link_endpoint.default.id}"]
  enable_details = true
}

output "first_endpoint_id" {
  value = data.alicloud_threat_detection_sas_private_link_endpoints.default.endpoints.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional) A list of Sas Private Link Endpoint IDs.
* `name_regex` - (Optional) A regex string to filter results by node name.
* `node_name` - (Optional) The name of the private link endpoint node.
* `enable_details` - (Optional) Default to false. Set it to true can get more details.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:
* `ids` - A list of Sas Private Link Endpoint IDs.
* `endpoints` - A list of Sas Private Link Endpoints. Each element contains the following attributes:
  * `id` - The ID of the Sas Private Link Endpoint.
  * `node_name` - The name of the private link endpoint node.
  * `security_group_id` - The security group ID.
  * `vpc_id` - The VPC ID.
  * `region_id` - The region ID.
  * `status` - The status of the endpoint.
  * `update_domain` - The update domain.
  * `jsrv_domain` - The jsrv domain.
  * `zones` - The availability zone information.
    * `v_switch_id` - The VSwitch ID.
    * `zone_id` - The Zone ID.
