---
subcategory: "Function Compute Service (FC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_custom_domain"
sidebar_current: "docs-alicloud-resource-fc"
description: |-
  Provides an Alicloud Function Compute Custom Domain resource. 
---

# alicloud_fc_custom_domain

Provides an Alicloud Function Compute custom domain resource. 
 For the detailed information, please refer to the [developer guide](https://www.alibabacloud.com/help/en/fc/developer-reference/api-fc-open-2021-04-06-createcustomdomain).

-> **NOTE:** Available since v1.98.0.


## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fc_custom_domain&exampleId=65a53235-ef37-0b7a-0879-314cf8b590dc5fa63714&activeTab=example&spm=docs.r.fc_custom_domain.0.65a53235ef&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "default" {
  name = "example-value-${random_integer.default.result}"
}

resource "alicloud_log_store" "default" {
  project = alicloud_log_project.default.name
  name    = "example-value"
}

resource "alicloud_ram_role" "default" {
  name        = "fcservicerole-${random_integer.default.result}"
  document    = <<EOF
  {
      "Statement": [
        {
          "Action": "sts:AssumeRole",
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "fc.aliyuncs.com"
            ]
          }
        }
      ],
      "Version": "1"
  }
  EOF
  description = "this is a example"
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

resource "alicloud_fc_service" "default" {
  name        = "example-value-${random_integer.default.result}"
  description = "example-value"
  role        = alicloud_ram_role.default.arn
  log_config {
    project                 = alicloud_log_project.default.name
    logstore                = alicloud_log_store.default.name
    enable_instance_metrics = true
    enable_request_metrics  = true
  }
}

resource "alicloud_oss_bucket" "default" {
  bucket = "terraform-example-${random_integer.default.result}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "index.py"
  content = "import logging \ndef handler(event, context): \nlogger = logging.getLogger() \nlogger.info('hello world') \nreturn 'hello world'"
}

resource "alicloud_fc_function" "default" {
  service     = alicloud_fc_service.default.name
  name        = "terraform-example"
  description = "example"
  oss_bucket  = alicloud_oss_bucket.default.id
  oss_key     = alicloud_oss_bucket_object.default.key
  memory_size = "512"
  runtime     = "python2.7"
  handler     = "hello.handler"
}

resource "alicloud_fc_custom_domain" "default" {
  domain_name = "terraform.functioncompute.com"
  protocol    = "HTTP"
  route_config {
    path          = "/login/*"
    service_name  = alicloud_fc_service.default.name
    function_name = alicloud_fc_function.default.name
    qualifier     = "?query"
    methods       = ["GET", "POST"]
  }
  cert_config {
    cert_name   = "example"
    certificate = "-----BEGIN CERTIFICATE-----\nMIICWD****-----END CERTIFICATE-----"
    private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICX****n-----END RSA PRIVATE KEY-----"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) The custom domain name. For example, "example.com".
* `protocol` - (Required) The protocol, `HTTP` or `HTTP,HTTPS`.
* `route_config` - (Optional) The configuration of domain route, mapping the path and Function Compute function.See [`route_config`](#route_config) below.
* `cert_config` - (Optional) The configuration of HTTPS certificate.See [`cert_config`](#cert_config) below.


### `route_config`

The route_config supports the following:

* `path` - (Required) The path that requests are routed from.
* `service_name` - (Required) The name of the Function Compute service that requests are routed to. 
* `function_name` - (Required) The name of the Function Compute function that requests are routed to.
* `qualifier` - (Optional) The version or alias of the Function Compute service that requests are routed to. For example, qualifier v1 indicates that the requests are routed to the version 1 Function Compute service. For detail information about version and alias, please refer to the [developer guide](https://www.alibabacloud.com/help/doc-detail/96464.htm).
* `methods` - (Optional) The requests of the specified HTTP methos are routed from. Valid method: GET, POST, DELETE, HEAD, PUT and PATCH. For example, "GET, HEAD" methods indicate that only requests from GET and HEAD methods are routed.

### `cert_config`

The cert_config supports the following:

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

```shell
$ terraform import alicloud_fc_custom_domain.foo my-fc-custom-domain
```
