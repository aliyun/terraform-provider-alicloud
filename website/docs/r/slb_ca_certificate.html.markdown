---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_ca_certificate"
sidebar_current: "docs-alicloud-resource-slb-ca-certificate"
description: |-
  Provides a Load Banlancer CA Certificate resource.
---

# alicloud\_slb\_ca\_certificate

A Load Balancer CA Certificate is used by the listener of the protocol https.

For information about slb and how to use it, see [What is Server Load Balancer](https://www.alibabacloud.com/help/doc-detail/27539.htm).

For information about CA Certificate and how to use it, see [Configure CA Certificate](https://www.alibabacloud.com/help/doc-detail/85968.htm).


## Example Usage

* using CA certificate content

```terraform
# create a CA certificate
resource "alicloud_slb_ca_certificate" "foo" {
  ca_certificate_name = "tf-testAccSlbCACertificate"
  ca_certificate      = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJnI******90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}
```

* using CA certificate file

```terraform
resource "alicloud_slb_ca_certificate" "foo-file" {
  ca_certificate_name = "tf-testAccSlbCACertificate"
  ca_certificate      = file("${path.module}/ca_certificate.pem")
}
```

## Argument Reference

The following arguments are supported:

* `ca_certificate_name` - (Optional, Available in 1.123.1+) Name of the CA Certificate.
* `ca_certificate` - (Required, ForceNew) the content of the CA certificate.
* `resource_group_id` - (Optional, ForceNew, Available in 1.58.0+) The Id of resource group which the slb_ca certificate belongs.
* `tags` - (Optional, Available in v1.66.0+) A mapping of tags to assign to the resource.
* `name` - (Deprecated) Field `name` has been deprecated from provider version 1.123.1. New field `ca_certificate_name` instead

## Attributes Reference

The following attributes are exported:

* `id` - The Id of CA Certificate .

### Timeouts

-> **NOTE:** -Available in 1.123.1+

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `delete` - (Defaults to 5 mins) Used when delete the SLB CA Certificate.

## Import

Server Load balancer CA Certificate can be imported using the id, e.g.

```
$ terraform import alicloud_slb_ca_certificate.example abc123456
```
