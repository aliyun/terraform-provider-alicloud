---
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
    name = "tf-testAccSlbServerCertificate"
    server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
    private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0knDrlNdiys******ErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
  }

  # create a https listener with the server certificate above.
  resource "alicloud_slb" "instance" {
    name                 = "${var.slb_name}"
    internet_charge_type = "${var.internet_charge_type}"
    internet             = "${var.internet}"
  }

  resource "alicloud_slb_listener" "https" {
    load_balancer_id          = "${alicloud_slb.instance.id}"
    backend_port              = 80
    frontend_port             = 443
    protocol                  = "https"
    sticky_session            = "on"
    sticky_session_type       = "insert"
    cookie                    = "testslblistenercookie"
    cookie_timeout            = 86400
    health_check              = "on"
    health_check_uri          = "/cons"
    health_check_connect_port = 20
    healthy_threshold         = 8
    unhealthy_threshold       = 8
    health_check_timeout      = 8
    health_check_interval     = 5
    health_check_http_code    = "http_2xx,http_3xx"
    bandwidth                 = 10
    ssl_certificate_id        = "${alicloud_slb_server_certificate.foo.id}"
  }
  variable "slb_name" {
    default = "slb_htts_server_certificate"
  }

  variable "internet_charge_type" {
    default = "PayByTraffic"
  }

  variable "internet" {
    default = true
  }
```

* using server_certificate/private file example

```
  # create a server certificate
  resource "alicloud_slb_server_certificate" "foo" {
    name = "tf-testAccSlbServerCertificate"
    server_certificate_file = "server_certificate.pem"
    private_key_file = "private_key.pem"
  }

  # create a https listener with the server certificate above.
  resource "alicloud_slb" "instance" {
    name                 = "${var.slb_name}"
    internet_charge_type = "${var.internet_charge_type}"
    internet             = "${var.internet}"
  }

  resource "alicloud_slb_listener" "https" {
    load_balancer_id          = "${alicloud_slb.instance.id}"
    backend_port              = 80
    frontend_port             = 443
    protocol                  = "https"
    sticky_session            = "on"
    sticky_session_type       = "insert"
    cookie                    = "testslblistenercookie"
    cookie_timeout            = 86400
    health_check              = "on"
    health_check_uri          = "/cons"
    health_check_connect_port = 20
    healthy_threshold         = 8
    unhealthy_threshold       = 8
    health_check_timeout      = 8
    health_check_interval     = 5
    health_check_http_code    = "http_2xx,http_3xx"
    bandwidth                 = 10
    ssl_certificate_id        = "${alicloud_slb_server_certificate.foo.id}"
  }
  variable "slb_name" {
    default = "slb_htts_server_certificate"
  }

  variable "internet_charge_type" {
    default = "PayByTraffic"
  }

  variable "internet" {
    default = true
  }
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the Server Certificate.
* `server_certificate` - (Optional, ForceNew) the content of the ssl certificate. where `alicloud_certificate_id` is null, it is required, otherwise it is ignored.
* `private_key` - (Optional, ForceNew) the content of privat key of the ssl certificate specified by `server_certificate`. where `alicloud_certificate_id` is null, it is required, otherwise it is ignored.
* `server_certificate_file` - (Optional, ForceNew) the file of the ssl certificate. where  both `alicloud_certificate_id` and `server_certificate` is null, it is required, otherwise it is ignored.
* `private_key_file` - (Optional, ForceNew) the file of privat key of the ssl certificate specified by `server_certificate`. where both `alicloud_certificate_id` and `private_key` is null, it is required, otherwise it is ignored.
* `alicloud_certificate_id` - (Optional) an id of server certificate ssued/proxied by alibaba cloud. but it is not supported on the international site  of alibaba cloud now.
* `alicloud_certificate_name`- (Optional) the name of the certificate specified by `alicloud_certificate_id`.but it is not supported on the international site  of alibaba cloud now.

## Attributes Reference

The following attributes are exported:

* `id` - The Id of Server Certificate (SSL Certificate).

## Import

Server Load balancer Server Certificate can be imported using the id, e.g.

```
$ terraform import alicloud_slb_server_certificate.example abc123456
```
