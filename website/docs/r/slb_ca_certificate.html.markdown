---
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

```
  # create a CA certificate
  resource "alicloud_slb_ca_certificate" "foo" {
    name = "tf-testAccSlbCACertificate"
    ca_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJnI******90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  }
```

* using CA certificate file

```
resource "alicloud_slb_ca_certificate" "foo-file" {
  name = "tf-testAccSlbCACertificate"
  ca_certificate = "${file("${path.module}/ca_certificate.pem")}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (ForceNew) Name of the CA Certificate.
* `ca_certificate` - (Required, ForceNew) the content of the CA certificate.

## Attributes Reference

The following attributes are exported:

* `id` - The Id of CA Certificate .

## Import

Server Load balancer CA Certificate can be imported using the id, e.g.

```
$ terraform import alicloud_slb_ca_certificate.example abc123456
```
