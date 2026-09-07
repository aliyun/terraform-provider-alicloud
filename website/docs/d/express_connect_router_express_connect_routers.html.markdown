---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_express_connect_routers"
description: |-
  Provides a list of Express Connect Router Express Connect Routers to the user.
---

# alicloud_express_connect_router_express_connect_routers

This data source provides the Express Connect Router Express Connect Router of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.291.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_express_connect_router_express_connect_router" "default" {
  alibaba_side_asn = "65533"
}

data "alicloud_express_connect_router_express_connect_routers" "ids" {
  ids = [alicloud_express_connect_router_express_connect_router.default.id]
}

output "express_connect_router_express_connect_routers_id_0" {
  value = data.alicloud_express_connect_router_express_connect_routers.ids.routers.0.id
}
```

## Argument Reference

The following attributes are exported:

* `ids` - (Optional, List) A list of Express Connect Router IDs.
* `name_regex` - (Optional) A regex string to filter results by Express Connect Router name.
* `ecr_id` - (Optional) The ID of ECR.
* `ecr_name` - (Optional) The name of ECR.
* `resource_group_id` - (Optional) The ID of the resource group to which the ECR belongs.
* `status` - (Optional) The deployment status of the service instance. Valid values: `ACTIVE`, `UPDATING`, `ASSOCIATING`, `DISSOCIATING`, `LOCKED_ATTACHING`, `LOCKED_DETACHING`, `RECLAIMING`, `DELETING`.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Express Connect Router names.
* `routers` - A list of Express Connect Routers. Each element contains the following attributes:
  * `id` - The ID of the Express Connect Router.
  * `ecr_id` - The ID of ECR.
  * `ecr_name` - The name of ECR.
  * `description` - The description of the ECR instance.
  * `owner_id` - The ID of the Alibaba Cloud account to which the ECR belongs.
  * `resource_group_id` - The ID of the resource group to which the ECR instance belongs.
  * `alibaba_side_asn` - The ASN of the ECR instance.
  * `biz_status` - The business status of the service instance.
  * `status` - The deployment status of the service instance.
  * `create_time` - The time when the ECR was created.
  * `modify_time` - The time when the ECR was modified.
  * `tags` - A mapping of tags to assign to the resource.
