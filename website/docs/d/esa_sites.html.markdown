---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_sites"
sidebar_current: "docs-alicloud-datasource-esa-sites"
description: |-
  Provides a list of Esa Site owned by an Alibaba Cloud account.
---

# alicloud_esa_sites

This data source provides Esa Site available to the user.[What is Site](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateSite)

-> **NOTE:** Available since v1.244.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_esa_rate_plan_instance" "defaultIEoDfU" {
  type         = "NS"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "basic"
}


resource "alicloud_esa_site" "default" {
  site_name   = "bcd.com"
  coverage    = "overseas"
  access_type = "NS"
  instance_id = alicloud_esa_rate_plan_instance.defaultIEoDfU.id
}

data "alicloud_esa_sites" "default" {
  ids        = ["${alicloud_esa_site.default.id}"]
  name_regex = alicloud_esa_site.default.site_name
  site_name  = "bcd.com"
}

output "alicloud_esa_site_example_id" {
  value = data.alicloud_esa_sites.default.sites.0.id
}
```

## Argument Reference

The following arguments are supported:

* `access_type` - (Optional) Access type. Value:-**NS**: Managed access via NS.-**CNAME**: access through CNAME.
* `coverage` - (Optional) Acceleration area
* `only_enterprise` - (Optional) Only the Enterprise version, pass **true** when the table only queries the enterprise version of the site.
* `page_number` - (Optional) Current page number.
* `page_size` - (Optional) Number of records per page.
* `plan_subscribe_type` - (Optional) Package subscription type. Value:-**basicplan**: Basic version.-**standardplan**: Standard version.-**advancedplan**: Advanced version.-**enterpriseplan**: Enterprise Edition.
* `resource_group_id` - (Optional) The ID of the resource group
* `site_name` - (Optional) Site Name
* `site_search_type` - (Optional) The search match pattern for the site name. The default value is exact match. Valid values:-**prefix**: matches the prefix.-**suffix**: The suffix matches.-**exact**: exact match.-**fuzzy**: fuzzy match.
* `status` - (Optional) The status of the resource
* `tags` - (Optional) Resource tags
* `ids` - (Optional, Computed) A list of Site IDs.
* `name_regex` - (Optional) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Site IDs.
* `names` - A list of name of Sites.
* `sites` - A list of Site Entries.

---

### `sites`

The sites supports the following attributes:

* `access_type` - Access type. Value:-**NS**: Managed access via NS.-**CNAME**: access through CNAME.
* `coverage` - Acceleration area
* `create_time` - Creation time
* `instance_id` - The ID of the associated package instance.
* `modify_time` - Modification time
* `name_server_list` - Site Resolution Name Server List
* `resource_group_id` - The ID of the resource group
* `site_id` - Site ID
* `site_name` - Site Name
* `status` - The status of the resource
* `id` - The ID of the resource supplied above.
