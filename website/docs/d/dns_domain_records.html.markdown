---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_domain_records"
sidebar_current: "docs-alicloud-datasource-dns-domain-records"
description: |-
    Provides a list of records available to the domain.
---

# alicloud_dns_domain_records

-> **DEPRECATED:** This data source has been deprecated from [v1.3.2](https://github.com/alibaba/terraform-provider/releases/tag/V1.3.2). Please use the data source `alicloud_dns_records` instead.

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) The domain name associated with the records.
* `host_record_regex` - (Optional) A regex string to filter results by host record.
* `value_regex` - (Optional) A regex string to filter results by record value.
* `type` - (Optional) The type of the domain record. Valid values: `A`, `NS`, `MX`, `TXT`, `CNAME`, `SRV`, `AAAA`, `CAA`, `REDIRECT_URL` and `FORWORD_URL`.
* `line` - (Optional) The ISP line of the record.
* `status` - (Optional) The status of the record. Valid values: `enable` and `disable`.
* `is_locked` - (Optional) Whether the record is locked.
* `ids` - (Optional) A list of record IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `urls` - A list of record URLs.
* `records` - A list of records.

---

### `records`

The records supports the following attributes:

* `record_id` - The ID of the record.
* `domain_name` - The domain name associated with the record.
* `host_record` - The host record.
* `type` - The type of the record.
* `value` - The value of the record.
* `ttl` - The TTL of the record.
* `priority` - The priority of the MX record.
* `line` - The ISP line of the record.
* `status` - The status of the record.
* `locked` - Whether the record is locked.
