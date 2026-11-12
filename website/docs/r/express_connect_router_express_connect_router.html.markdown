---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_express_connect_router"
description: |-
  Provides a Alicloud Express Connect Router Express Connect Router resource.
---

# alicloud_express_connect_router_express_connect_router

Provides a Express Connect Router Express Connect Router resource. Express Connect Router.

For information about Express Connect Router Express Connect Router and how to use it, see [What is Express Connect Router](https://next.api.alibabacloud.com/api/ExpressConnectRouter/2023-09-01/CreateExpressConnectRouter).

-> **NOTE:** Available since v1.224.0.

-> **NOTE:** If you maintain a narrow custom RAM policy that predates v2.3.0, `refresh` and `plan` will surface a `403 NoPermission` error from whichever supplementary action your role lacks, rather than silently leaving the computed field empty. Grant the two actions listed in [RAM Permissions](#ram-permissions) below — at resource level (`acs:expressconnectrouter::<accountId>:expressconnectrouter/<EcrId>`) or service level (`*`) — so refresh and plan stay green.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_router_express_connect_router&exampleId=dd0fd8dd-7c00-81b4-75c5-2a9465c60c311b2d9f15&activeTab=example&spm=docs.r.express_connect_router_express_connect_router.0.dd0fd8dd7c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_express_connect_router_express_connect_router" "defaultM9YxGW" {
  alibaba_side_asn = "65533"
}
```

### Deleting `alicloud_express_connect_router_express_connect_router` or removing it from your configuration

The `alicloud_express_connect_router_express_connect_router` resource allows you to manage  `ecr_id = ""`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_express_connect_router_express_connect_router&spm=docs.r.express_connect_router_express_connect_router.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `alibaba_side_asn` - (Required, ForceNew) ASN representing resources.
* `description` - (Optional) Represents the description of the leased line gateway.
* `ecr_name` - (Optional) Name of the Gateway representing the leased line.
* `regions` - (Optional) List of regions representing leased line gateways. See [`regions`](#regions) below.
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the ECR instance belongs.
  - A string consisting of letters, numbers, hyphens (-), and underscores (_), and the string length can be 0 to 64 characters.
* `tags` - (Optional, Map) The tag of the resource.

### `regions`

The regions supports the following:
* `region_id` - (Optional) Representative region ID.
* `transit_mode` - (Optional, Computed) Represents the forwarding mode of the current region.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Represents the creation time of the resource.
* `status` - The status of the resource.
* `disabled_route_entries` - (Computed, List, Available since v2.3.0) The route entries that are currently disabled on the Express Connect Router. See [`disabled_route_entries`](#disabled_route_entries) below.
* `granted_instances` - (Computed, List, Available since v2.3.0) The network instances granted to the Express Connect Router. See [`granted_instances`](#granted_instances) below.

### `disabled_route_entries`

The disabled_route_entries supports the following:
* `destination_cidr_block` - (Computed, String, Available since v2.3.0) The destination CIDR block of the disabled route entry.
* `nexthop_instance_id` - (Computed, String, Available since v2.3.0) The ID of the next hop instance of the disabled route entry.

### `granted_instances`

The granted_instances supports the following:
* `node_id` - (Computed, String, Available since v2.3.0) The ID of the granted network instance, for example `vpc-xxx` or `vbr-xxx`.
* `node_type` - (Computed, String, Available since v2.3.0) The type of the granted network instance. Valid values: `VPC`, `VBR`.
* `node_region_id` - (Computed, String, Available since v2.3.0) The region ID of the granted network instance.
* `node_owner_uid` - (Computed, String, Available since v2.3.0) The ID of the Alibaba Cloud account that owns the granted network instance.

## RAM Permissions

The two read-only attributes [`disabled_route_entries`](#disabled_route_entries) and [`granted_instances`](#granted_instances) (available since v2.3.0) are populated by the following supplementary Describe actions. Each is an independent RAM permission on the `expressconnectrouter` service and is queried on every `read`/`refresh`/`plan`:

* `expressconnectrouter:DescribeDisabledExpressConnectRouterRouteEntries` — lists the route entries currently disabled on the ECR.
* `expressconnectrouter:DescribeInstanceGrantedToExpressConnectRouter` — lists the network instances granted to the ECR.

These actions are in addition to the base resource actions (`expressconnectrouter:CreateExpressConnectRouter`, `expressconnectrouter:ModifyExpressConnectRouter`, `expressconnectrouter:DeleteExpressConnectRouter`, `expressconnectrouter:DescribeExpressConnectRouter`, etc.).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Express Connect Router.
* `delete` - (Defaults to 5 mins) Used when delete the Express Connect Router.
* `update` - (Defaults to 5 mins) Used when update the Express Connect Router.

## Import

Express Connect Router Express Connect Router can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_router_express_connect_router.example <id>
```