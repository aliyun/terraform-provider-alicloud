---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_sources"
sidebar_current: "docs-alicloud-datasource-apig-sources"
description: |-
  Provides a list of Apig Source owned by an Alibaba Cloud account.
---

# alicloud_apig_sources

This data source provides Apig Source available to the user. [What is Source](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateSource)

-> **NOTE:** Available since v1.292.0.

## Example Usage

```terraform
data "alicloud_apig_sources" "default" {
  gateway_id = "gw-cq7l5s5lhtgi6q***"
}

output "first_source_id" {
  value = data.alicloud_apig_sources.default.sources.0.id
}
```

## Argument Reference

The following arguments are supported:
* `gateway_id` - (Optional) The ID of the gateway instance.
* `resource_group_id` - (Optional) The ID of the resource group.
* `type` - (Optional) The type of the source. Valid values: `K8S`, `MSE_NACOS`.
* `ids` - (Optional, Computed) A list of Source IDs.
* `name_regex` - (Optional) A regex string to filter results by Source name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Source IDs.
* `names` - A list of name of Sources.
* `sources` - A list of Source Entries. Each element contains the following attributes:
  * `association_reason` - The reason for the association status.
  * `association_status` - The association status of the source.
  * `create_time` - The creation timestamp of the source.
  * `gateway_id` - **NOTE:** This field is only available when `enable_details` is `true`. The ID of the gateway instance.
  * `k8s_source_info` - The ACK cluster source information.
    * `cluster_id` - The ID of the ACK cluster.
  * `nacos_source_info` - The MSE Nacos source information.
    * `address` - The access address of the Nacos instance.
    * `cluster_id` - The ID of the Nacos cluster.
    * `instance_id` - The ID of the MSE Nacos instance.
  * `resource_group_id` - The ID of the resource group.
  * `source_id` - The ID of the source.
  * `source_name` - The name of the source.
  * `type` - **NOTE:** This field is only available when `enable_details` is `true`. The type of the source.
  * `update_time` - The update timestamp of the source.
  * `id` - The ID of the resource supplied above.
