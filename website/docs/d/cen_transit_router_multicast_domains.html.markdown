---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domains"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-multicast-domains"
description: |-
  Provides a list of Cen Transit Router Multicast Domains to the user.
---

# alicloud\_cen\_transit\_router\_multicast\_domains

This data source provides the Cen Transit Router Multicast Domains of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_transit_router_multicast_domains" "ids" {
  ids               = ["example_id"]
  transit_router_id = "your_transit_router_id"
}

output "cen_transit_router_multicast_domain_id_0" {
  value = data.alicloud_cen_transit_router_multicast_domains.ids.domains.0.id
}

data "alicloud_cen_transit_router_multicast_domains" "nameRegex" {
  name_regex        = "^my-name"
  transit_router_id = "your_transit_router_id"
}

output "cen_transit_router_multicast_domain_id_1" {
  value = data.alicloud_cen_transit_router_multicast_domains.nameRegex.domains.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Transit Router Multicast Domain IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Transit Router Multicast Domain name.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router.
* `transit_router_multicast_domain_id` - (Optional, ForceNew) The ID of the multicast domain.
* `status` - (Optional, ForceNew) The status of the multicast domain. Valid Value: `Active`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Transit Router Multicast Domain names.
* `domains` - A list of Cen Transit Router Multicast Domains. Each element contains the following attributes:
	* `id` - The ID of the Transit Router Multicast Domain.
	* `transit_router_id` - The ID of the transit router.
	* `transit_router_multicast_domain_id` - The ID of the Transit Router Multicast Domain.
	* `transit_router_multicast_domain_name` - The name of the Transit Router Multicast Domain.
	* `transit_router_multicast_domain_description` - The description of the Transit Router Multicast Domain.
	* `status` - The status of the Transit Router Multicast Domain.
	