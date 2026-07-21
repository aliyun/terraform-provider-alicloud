---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_server_group_server_attachments"
sidebar_current: "docs-alicloud-datasource-nlb-server-group-server-attachments"
description: |-
  Provides a list of Nlb Server Group Server Attachments to the user.
---

# alicloud_nlb_server_group_server_attachments

This data source provides the Nlb Server Group Server Attachments of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.192.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_nlb_server_group_server_attachments" "ids" {
  ids = ["example_value"]
}
output "nlb_server_group_server_attachment_id_1" {
  value = data.alicloud_nlb_server_group_server_attachments.ids.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:

* `server_group_id` - (Optional, ForceNew) The ID of the server group.
* `server_ids` - (Optional, ForceNew) The IDs of the servers. You can specify at most 40 server IDs in each call.
* `server_ips` - (Optional, ForceNew) The IP addresses of the servers. You can specify at most 40 server IP addresses in each call.
* `ids` - (Optional, ForceNew, Computed) A list of Server Group Server Attachment IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `attachments` - A list of Nlb Server Group Server Attachments. Each element contains the following attributes:
	* `description` - The description of the backend server.
	* `port` - The port used by the backend server.
	* `server_group_id` - The ID of the server group.
	* `server_id` - The ID of the server.
	* `server_ip` - The IP address of the backend server.
	* `server_type` - The type of the backend server.
	* `status` - Indicates the status of the backend server.
	* `weight` - The weight of the backend server.
	* `zone_id` - The zone ID of the server.
	* `id` - The ID of the server group. The value is formulated as `<server_group_id>:<server_id>:<server_type>:<port>`.