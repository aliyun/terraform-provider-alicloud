---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_record"
sidebar_current: "docs-alicloud-resource-alidns-record"
description: |-
  Provides a Alidns Domain Record resource.
---

# alicloud\_alidns\_record

Provides a Alidns Record resource. For information about Alidns Domain Record and how to use it, see [What is Resource Alidns Record](https://www.alibabacloud.com/help/en/doc-detail/29772.htm).

-> **NOTE:** Available in v1.85.0+.

-> **NOTE:** When the site is an international site, the `type` neither supports `REDIRECT_URL` nor `REDIRECT_URL`

## Example Usage

```terraform
# Create a new Domain Record
resource "alicloud_alidns_record" "record" {
  domain_name = "domainname"
  rr          = "@"
  type        = "A"
  value       = "192.168.99.99"
  remark      = "Test new alidns record."
  status      = "ENABLE"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Name of the domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `rr` - (Required) Host record for the domain record. This host_record can have at most 253 characters, and each part split with `.` can have at most 63 characters, and must contain only alphanumeric characters or hyphens, such as `-`, `.`, `*`, `@`, and must not begin or end with `-`.
* `type` - (Required) The type of domain record. Valid values: `A`,`NS`,`MX`,`TXT`,`CNAME`,`SRV`,`AAAA`,`CAA`, `REDIRECT_URL` and `FORWORD_URL`.
* `value` - (Required) The value of domain record, When the `type` is `MX`,`NS`,`CNAME`,`SRV`, the server will treat the `value` as a fully qualified domain name, so it's no need to add a `.` at the end.
* `ttl` - (Optional) The effective time of domain record. Its scope depends on the edition of the cloud resolution. Free is `[600, 86400]`, Basic is `[120, 86400]`, Standard is `[60, 86400]`, Ultimate is `[10, 86400]`, Exclusive is `[1, 86400]`. Default value is `600`.
* `priority` - (Optional) The priority of domain record. Valid values: `[1-10]`. When the `type` is `MX`, this parameter is required.
* `line` - (Optional) The resolution line of domain record. When the `type` is `FORWORD_URL`, this parameter must be `default`. Default value is `default`. For checking all resolution lines enumeration please visit [Alibaba Cloud DNS doc](https://www.alibabacloud.com/help/doc-detail/34339.htm) or using alicloud_dns_resolution_lines in data source to get the value. 
* `lang` - (Optional) User language. 
* `remark` - (Optional) The remark of the domain record. 
* `status` - (Optional) The status of the domain record. Valid values: `ENABLE`,`DISABLE`. 
* `user_client_ip` - (Optional) The IP address of the client. 

## Attributes Reference

The following attributes are exported:

* `id` - The id of Domain Record.

### Timeouts

-> **NOTE:** Available in 1.99.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `Create` - (Defaults to 4 mins) Used when terminating the Alidns record instance.  
* `delete` - (Defaults to 6 mins) Used when terminating the Alidns record instance.  


## Import

Alidns Domain Record can be imported using the id, e.g.

```
$ terraform import alicloud_alidns_record.example abc123456
```
