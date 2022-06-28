---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_custom_domain"
sidebar_current: "docs-alicloud-resource-fc"
description: |-
  Provides an Alicloud Function Compute Custom Domain resource. 
---

# alicloud\_fc\_custom_domain

Provides an Alicloud Function Compute custom domain resource. 
 For the detailed information, please refer to the [developer guide](https://www.alibabacloud.com/help/doc-detail/90759.htm).

-> **NOTE:** Available in 1.98.0+


## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testaccalicloudfcservice"
}

resource "alicloud_fc_custom_domain" "default" {
  domain_name = "terraform.functioncompute.com"
  protocol    = "HTTP"
  route_config {
    path          = "/login/*"
    service_name  = alicloud_fc_service.default.name
    function_name = alicloud_fc_function.default.name
    qualifier     = "v1"
    methods       = ["GET", "POST"]
  }
  cert_config {
    cert_name   = "your certificate name"
    private_key = "your private key"
    certificate = "your certificate data"
  }
}

resource "alicloud_fc_service" "default" {
  name        = var.name
  description = "${var.name}-description"
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "fc/hello.zip"
  content = <<EOF
		# -*- coding: utf-8 -*-
	def handler(event, context):
		print "hello world"
		return 'hello world'
	EOF
}

resource "alicloud_fc_function" "default" {
  service     = alicloud_fc_service.default.name
  name        = var.name
  oss_bucket  = alicloud_oss_bucket.default.id
  oss_key     = alicloud_oss_bucket_object.default.key
  memory_size = 512
  runtime     = "python2.7"
  handler     = "hello.handler"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) The custom domain name. For example, "example.com".
* `protocol` - (Required) The protocol, `HTTP` or `HTTP,HTTPS`.
* `route_config` - (Optional) The configuration of domain route, mapping the path and Function Compute function.
* `cert_config` - (Optional) The configuration of HTTPS certificate.


**route_config** includes the following arguments:

* `path` - (Required) The path that requests are routed from.
* `serivce_name` - (Required) The name of the Function Compute service that requests are routed to. 
* `function_name` - (Required) The name of the Function Compute function that requests are routed to.
* `qualifier` - (Optional) The version or alias of the Function Compute service that requests are routed to. For example, qualifier v1 indicates that the requests are routed to the version 1 Function Compute service. For detail information about version and alias, please refer to the [developer guide](https://www.alibabacloud.com/help/doc-detail/96464.htm).
* `methods` - (Optional) The requests of the specified HTTP methos are routed from. Valid method: GET, POST, DELETE, HEAD, PUT and PATCH. For example, "GET, HEAD" methods indicate that only requests from GET and HEAD methods are routed.

**cert_config** includes the following arguments:

* `cert_name` - (Required) The name of the certificate, used to distinguish different certificates.
* `private_key` - (Required) Private key of the HTTPS certificates, follow the 'pem' format.
* `certificate` - (Required) Certificate data of the HTTPS certificates, follow the 'pem' format.

## Attributes Reference

The following arguments are exported:

* `id` -The id of the custom domain. It is the same as the domain name.
* `account_id` - The account id.
* `api_version` - The api version of Function Compute.
* `created_time` - The date this resource was created.
* `last_modified_time` - The date this resource was last modified.

## Import

Function Compute custom domain can be imported using the id or the domain name, e.g.

```
$ terraform import alicloud_fc_custom_domain.foo my-fc-custom-domain
```
