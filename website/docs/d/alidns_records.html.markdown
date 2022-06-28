---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_records"
sidebar_current: "docs-alicloud-datasource-alidns-records"
description: |-
    Provides a list of domain records available to the Alidns.
---

# alicloud\_alidns\_records

This data source provides a list of Alidns Domain Records in an Alibaba Cloud account according to the specified filters.

-> **NOTE:**  Available in 1.86.0+.

## Example Usage

```terraform 
data "alicloud_alidns_records" "records_ds" {
  domain_name = "xiaozhu.top"
  ids         = ["1978593525779****"]
  type        = "A"
  output_file = "records.txt"
}

output "first_record_id" {
  value = "${data.alicloud_alidns_records.records_ds.records.0.record_id}"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) The domain name associated to the records. 
* `type` - (Optional) Record type. Valid values: `A`, `NS`, `MX`, `TXT`, `CNAME`, `SRV`, `AAAA`, `REDIRECT_URL`, `FORWORD_URL` .
* `line` - (Optional) ISP line. For checking all resolution lines enumeration please visit [Alibaba Cloud DNS doc](https://www.alibabacloud.com/help/doc-detail/34339.htm) 
* `status` - (Optional) Record status. Valid values: `ENABLE` and `DISABLE`.
* `is_locked` - (Optional, type: bool) Whether the record is locked or not.
* `ids` - (Optional) A list of record IDs.
* `direction` - (Optional) Sorting direction. Valid values: `DESC`,`ASC`. Default to `AESC`.
* `group_id` - (Optional) Domain name group ID.
* `key_word` - (Optional) Keywords.
* `lang` - (Optional) User language.
* `order_by` - (Optional) Sort by. Sort from newest to oldest according to the time added by resolution.
* `rr_key_word` - (Optional) The keywords recorded by the host are searched according to the `%RRKeyWord%` mode, and are not case sensitive.
* `search_mode` - (Optional) Search mode, Valid values: `LIKE`, `EXACT`, `ADVANCED`, `LIKE` (fuzzy), `EXACT` (accurate) search supports KeyWord field, `ADVANCED` (advanced) mode supports other fields.
* `type_key_word` - (Optional) Analyze type keywords, search by full match, not case sensitive.
* `value_key_word` - (Optional) The keywords of the recorded value are searched according to the `%ValueKeyWord%` mode, and are not case sensitive.
* `rr_regex` - (Optional) Host record regex. 
* `value_regex` - (Optional) Host record value regex. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of record IDs. 
* `records` - A list of records. Each element contains the following attributes:
  * `id` - ID of the resource.
  * `record_id` - ID of the record.
  * `domain_name` - Name of the domain record belongs to.
  * `rr` - Host record of the domain.
  * `value` - Host record value of the domain.
  * `type` - Type of the record.
  * `ttl` - TTL of the record.
  * `priority` - Priority of the `MX` record.
  * `line` - ISP line of the record. 
  * `status` - Status of the record.
  * `locked` - Indicates whether the record is locked.
  * `remark` - The remark of the domain record.  **NOTE:** Available in 1.144.0+.
