---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_records"
sidebar_current: "docs-alicloud-datasource-dns-records"
description: |-
    Provides a list of records available to the dns.
---

# alicloud\_dns\_records

This data source provides a list of DNS Domain Records in an Alibaba Cloud account according to the specified filters.

## Example Usage

```
data "alicloud_dns_records" "records_ds" {
  domain_name       = "xiaozhu.top"
  is_locked         = false
  type              = "A"
  host_record_regex = "^@"
  output_file       = "records.txt"
}

output "first_record_id" {
  value = "${data.alicloud_dns_records.records_ds.records.0.record_id}"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) The domain name associated to the records.
* `host_record_regex` - (Optional) Host record regex. 
* `value_regex` - (Optional) Host record value regex. 
* `type` - (Optional) Record type. Valid items are `A`, `NS`, `MX`, `TXT`, `CNAME`, `SRV`, `AAAA`, `REDIRECT_URL`, `FORWORD_URL` .
* `line` - (Optional) ISP line. Valid items are `default`, `telecom`, `unicom`, `mobile`, `oversea`, `edu`, `drpeng`, `btvn`, .etc. For checking all resolution lines enumeration please visit [Alibaba Cloud DNS doc](https://www.alibabacloud.com/help/doc-detail/34339.htm) 
* `status` - (Optional) Record status. Valid items are `ENABLE` and `DISABLE`.
* `is_locked` - (Optional, type: bool) Whether the record is locked or not.
* `ids` - (Optional, Available 1.52.2+) A list of record IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of record IDs. 
* `urls` - A list of entire URLs. Each item format as `<host_record>.<domain_name>`.
* `records` - A list of records. Each element contains the following attributes:
  * `record_id` - ID of the record.
  * `domain_name` - Name of the domain the record belongs to.
  * `host_record` - Host record of the domain.
  * `value` - Host record value of the domain.
  * `type` - Type of the record.
  * `ttl` - TTL of the record.
  * `priority` - Priority of the `MX` record.
  * `line` - ISP line of the record. 
  * `status` - Status of the record.
  * `locked` - Indicates whether the record is locked.
