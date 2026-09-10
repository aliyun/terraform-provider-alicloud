---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone_records"
sidebar_current: "docs-alicloud-datasource-pvtz-zone-records"
description: |-
  Provides a list of Private Zone Records to the user.
---

# alicloud_pvtz_zone_records

This data source provides the Private Zone Records of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.13.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example.com"
}

resource "alicloud_pvtz_zone" "default" {
  zone_name = var.name
}

resource "alicloud_pvtz_zone_record" "default" {
  zone_id  = alicloud_pvtz_zone.default.id
  rr       = "www"
  type     = "MX"
  value    = var.name
  ttl      = "60"
  priority = 2
  remark   = var.name
}

data "alicloud_pvtz_zone_records" "ids" {
  zone_id = alicloud_pvtz_zone_record.default.zone_id
  ids     = [alicloud_pvtz_zone_record.default.record_id]
}

output "pvtz_zone_records_id_0" {
  value = data.alicloud_pvtz_zone_records.ids.records.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List, Available since v1.53.0) A list of Private Zone Record IDs.
* `zone_id` - (Required, ForceNew) The ID of the private zone.
* `keyword` - (Optional, ForceNew) The keyword for record rr and value.
* `tag` - (Optional, ForceNew, Available since v1.109.0) The tag used to search for DNS records.
* `user_client_ip` - (Optional, ForceNew, Available since v1.109.0) The IP address of the client.
* `status` - (Optional, ForceNew, Available since v1.109.0) The status of the Resolve record. Valid values:
  - `ENABLE`: Enable resolution.
  - `DISABLE`: Pause parsing.
* `search_mode` - (Optional, ForceNew, Available since v1.109.0) The search mode. Default value: `EXACT`. Valid values:
  - `LIKE`: Fuzzy search.
  - `EXACT`: Exact search.
* `lang` - (Optional, ForceNew, Available since v1.109.0) The language of the response. Default value: `en`. Valid values: `en`, `zh`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `records` - A list of Zone Record. Each element contains the following attributes:
  * `id` - The ID of the Private Zone Record.
  * `record_id` - The ID of the Record.
  * `priority` - The priority of the MX record.
  * `remark` - The description of the Private Zone Record.
  * `rr` - The hostname of the Private Zone Record.
  * `resource_record` - The hostname of the Private Zone Record.
  * `ttl` - The time to live (TTL) of the Private Zone Record.
  * `type` - The type of the Private Zone Record.
  * `value` - The value of the Private Zone Record.
  * `status` - The state of the Private Zone Record.
