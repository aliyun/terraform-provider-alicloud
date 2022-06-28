---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_server_certificate"
sidebar_current: "docs-alicloud-resource-slb-server-certificate"
description: |-
  Provides a Load Banlancer Server Certificate resource.
---

# alicloud\_slb\_server\_certificate

A Load Balancer Server Certificate is an ssl Certificate used by the listener of the protocol https.

For information about slb and how to use it, see [What is Server Load Balancer](https://www.alibabacloud.com/help/doc-detail/27539.htm).

For information about Server Certificate and how to use it, see [Configure Server Certificate](https://www.alibabacloud.com/help/doc-detail/85968.htm).


## Example Usage

* using server_certificate/private content as string example

```
# create a server certificate
resource "alicloud_slb_server_certificate" "foo" {
  name               = "slbservercertificate"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0knDrlNdiys******ErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
```

* using server_certificate/private file example

```
# create a server certificate
resource "alicloud_slb_server_certificate" "foo" {
  name               = "slbservercertificate"
  server_certificate = file("${path.module}/server_certificate.pem")
  private_key        = file("${path.module}/private_key.pem")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the Server Certificate.
* `server_certificate` - (Optional, ForceNew) the content of the ssl certificate. where `alicloud_certificate_id` is null, it is required, otherwise it is ignored.
* `private_key` - (Optional, ForceNew) the content of privat key of the ssl certificate specified by `server_certificate`. where `alicloud_certificate_id` is null, it is required, otherwise it is ignored.
* `alicloud_certificate_id` - (Optional, ForceNew) an id of server certificate ssued/proxied by alibaba cloud. but it is not supported on the international site of alibaba cloud now.
* `alicloud_certificate_name` - (Optional, ForceNew) the name of the certificate specified by `alicloud_certificate_id`.but it is not supported on the international site of alibaba cloud now.
* `alicloud_certificate_region_id` - (Optional, ForceNew, Available in 1.69.0+) the region of the certificate specified by `alicloud_certificate_id`. but it is not supported on the international site of alibaba cloud now.
* `resource_group_id` - (Optional, ForceNew, Available in 1.58.0+) The Id of resource group which the slb server certificate belongs.
* `tags` - (Optional, Available in v1.66.0+) A mapping of tags to assign to the resource.
## Attributes Reference

The following attributes are exported:

* `id` - The Id of Server Certificate (SSL Certificate).

## Import

Server Load balancer Server Certificate can be imported using the id, e.g.

```
$ terraform import alicloud_slb_server_certificate.example abc123456
```
