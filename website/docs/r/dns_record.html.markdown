---
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_record"
sidebar_current: "docs-alicloud-resource-dns-record"
description: |-
  Provides a DNS Record resource.
---

# alicloud\_dns\_record

Provides a DNS Record resource.

## Example Usage

```
# Create a new Domain record
resource "alicloud_dns_record" "record" {
  name = "domainname"
  host_record = "@"
  type = "A"
  value = "192.168.99.99"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `host_record` - (Required) Host record for the domain record. This host_record can have at most 253 characters, and each part split with "." can have at most 63 characters, and must contain only alphanumeric characters or hyphens, such as "-",".","*","@",  and must not begin or end with "-".
* `type` - (Required) The type of domain record. Valid values are `A`,`NS`,`MX`,`TXT`,`CNAME`,`SRV`,`AAAA`,`CAA`, `REDIRECT_URL` and `FORWORD_URL`.
* `value` - (Required) The value of domain record.
* `ttl` - (Optional) The effective time of domain record. Its scope depends on the edition of the cloud resolution. Free is `[600, 86400]`, Basic is `[120, 86400]`, Standard is `[60, 86400]`, Ultimate is `[10, 86400]`, Exclusive is `[1, 86400]`. Default value is `600`.
* `priority` - (Optional) The priority of domain record. Valid values are `[1-10]`. When the `type` is `MX`, this parameter is required.
* `routing` - (Optional) The parsing line of domain record. Valid values are `default`, `telecom`, `unicom`, `mobile`, `oversea` and `edu`. When the `type` is `FORWORD_URL`, this parameter must be `default`. Default value is `default`.

## Attributes Reference

The following attributes are exported:

* `id` - The record id.
* `name` - (Required) The record domain name.
* `type` - (Required) The record type.
* `host_record` - The host record of record.
* `value` - The record value.
* `ttl` - The record effective time.
* `priority` - The record priority.
* `routing` - The record parsing line.
* `status` - The record status. `Enable` or `Disable`.
* `Locked` - The record locked state. `true` or `false`.

## Import

RDS record can be imported using the id, e.g.

```
$ terraform import alicloud_dns_record.example abc123456
```