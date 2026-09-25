---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_aigw_records"
description: |-
  Provides a list of Esa Aigw Records to the user.
---

# alicloud_esa_aigw_records

This data source provides the Esa Aigw Records of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_aigw_records" "default" {
  instance_id = "your-aigw-instance-id"
}

output "first_record" {
  value = data.alicloud_esa_aigw_records.default.records.0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the AI Gateway instance.
* `ids` - (Optional) A list of AIGW Record IDs. The ID is in the format `instance_id:record_name`.
* `record_name_regex` - (Optional) A regex string to filter results by the record name.
* `fuzzy_search_key` - (Optional) The fuzzy search key to filter records.
* `output_file` - (Optional) File name where to save data source results after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the `arguments` above:

* `ids` - A list of AIGW Record IDs.
* `records` - A list of Asa Aigw Records. Each element contains the following attributes:
  * `id` - The resource ID in format `instance_id:record_name`.
  * `instance_id` - The ID of the AI Gateway instance.
  * `record_name` - The domain name of the AIGW record.
  * `site_id` - The ID of the ESA site.
  * `create_time` - The creation time of the resource.
