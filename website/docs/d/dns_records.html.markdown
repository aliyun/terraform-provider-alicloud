---
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_records"
sidebar_current: "docs-alicloud-datasource-dns-records"
description: |-
    Provides a list of records available to the dns.
---

# alicloud\_dns\_records

The Dns Domain Records data source provides a list of Alicloud Dns Domain Records in an Alicloud account according to the specified filters.

## Example Usage

```
data "alicloud_dns_records" "record" {
  domain_name = "xiaozhu.top"
  is_locked = false
  type = "A"
  host_record_regex = "^@"
  output_file = "records.txt"
}

```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) A domain name which is the necessary parameter for the records query.
* `host_record_regex` - (Optional) Limit search to provide host record regex. 
* `value_regex` - (Optional) Limit search to provide host record value regex. 
* `type` - (Optional) Limit search to specific record type. Valid items are `A`, `NS`, `MX`, `TXT`, `CNAME`, `SRV`, `AAAA`, `REDIRECT_URL`, `FORWORD_URL` .
* `line` - (Optional) Limit search to specific parsing line. Valid items are `default`, `telecom`, `unicom`, `mobile`, `oversea`, `edu`.
* `status` - (Optional) Limit search to specific record status. Valid items are `ENABLE` and `DISABLE`.
* `is_locked` - (Optional, type: bool) Limit search to specific record lock status.
* `output_file` - (Optional) The name of file that can save records data source after running `terraform plan`.


## Attributes Reference

A list of records will be exported and its every element contains the following attributes:

* `record_id` - ID of the record.
* `domain_name` - Name of the domain which the record belong to.
* `host_record` - Host record of the record.
* `value` - Host record value of the record.
* `type` - Type of the record.
* `ttl` - TTL of the record.
* `priority` - Priority of the `MX` record.
* `line` - Parsing line of the record. 
* `status` - Status of the record.
* `locked` - Indicates whether the record is locked.