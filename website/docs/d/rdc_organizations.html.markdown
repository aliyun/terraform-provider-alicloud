---
subcategory: "Apsara Devops(RDC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_rdc_organizations"
sidebar_current: "docs-alicloud-datasource-rdc-organizations"
description: |-
  Provides a list of Rdc Organizations to the user.
---

# alicloud\_rdc\_organizations

This data source provides the Rdc Organizations of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.137.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testAccOrganizations-Organizations"
}

resource "alicloud_rdc_organization" "default" {
  organization_name = var.name
  source            = var.name
}
data "alicloud_rdc_organizations" "ids" {
  ids = [alicloud_rdc_organization.default.id]
}
output "rdc_organization_id_1" {
  value = data.alicloud_rdc_organizations.ids.id
}

data "alicloud_rdc_organizations" "nameRegex" {
  name_regex = "^my-Organization"
}
output "rdc_organization_id_2" {
  value = data.alicloud_rdc_organizations.nameRegex.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Organization IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Organization name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `real_pk` - (Optional, ForceNew) User pk, not required, only required when the ak used by the calling interface is inconsistent with the user pk

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Organization names.
* `organizations` - A list of Rdc Organizations. Each element contains the following attributes:
	* `id` - The ID of the Organization.
	* `organization_id` - The first ID of the resource.
	* `organization_name` - Company name.
