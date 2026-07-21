---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_services"
description: |-
  Provides a list of Private Link Vpc Endpoint Services to the user.
---

# alicloud_privatelink_vpc_endpoint_services

This data source provides the Private Link Vpc Endpoint Services of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.109.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_privatelink_vpc_endpoint_service" "default" {
  service_description    = var.name
  auto_accept_connection = true
}

data "alicloud_privatelink_vpc_endpoint_services" "ids" {
  ids = [alicloud_privatelink_vpc_endpoint_service.default.id]
}

output "privatelink_vpc_endpoint_services_id_0" {
  value = data.alicloud_privatelink_vpc_endpoint_services.ids.services.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` (Optional, ForceNew, List)  A list of Vpc Endpoint Service IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Vpc Endpoint Service name.
* `vpc_endpoint_service_name` - (Optional, ForceNew) The name of the endpoint service.
* `auto_accept_connection` - (Optional, ForceNew, Bool) Specifies whether to automatically accept endpoint connection requests. Valid values: : `true`, `false`.
* `service_business_status` - (Optional, ForceNew) The service state of the endpoint service. Default value: `Normal`. Valid values: `Normal`, `FinancialLocked` and `SecurityLocked`.
* `status` - (Optional, ForceNew) The state of the endpoint service. Valid values: `Active`, `Creating`, `Deleted`, `Deleting` and `Pending`.
* `tags` - (Optional, ForceNew,  Available since v1.232.0)  A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Vpc Endpoint Service names.
* `services` - A list of Vpc Endpoint Services. Each element contains the following attributes:
  * `id` - The ID of the Vpc Endpoint Service.
  * `service_id` - The ID of the endpoint service.
  * `vpc_endpoint_service_name` - The name of the endpoint service.
  * `service_description` - The description of the endpoint service.
  * `service_domain` - The domain name of the endpoint service.
  * `connect_bandwidth` - The default maximum bandwidth of the endpoint connection.
  * `auto_accept_connection` - Indicates whether endpoint connection requests are automatically accepted.
  * `service_business_status` - The service state of the endpoint service.
  * `status` - The state of the endpoint service.
  * `tags` - The tags added to the resource.
