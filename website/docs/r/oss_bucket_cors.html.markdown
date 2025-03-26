---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_cors"
description: |-
  Provides a Alicloud OSS Bucket Cors resource.
---

# alicloud_oss_bucket_cors

Provides a OSS Bucket Cors resource. Cross-Origin Resource Sharing (CORS) allows web applications to access resources in other regions.

For information about OSS Bucket Cors and how to use it, see [What is Bucket Cors](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketcors).

-> **NOTE:** Available since v1.223.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_cors&exampleId=efee4706-8728-22cc-95cf-07c5f321a8a2e5a05f33&activeTab=example&spm=docs.r.oss_bucket_cors.0.efee470687&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_uuid" "default" {

}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-${random_uuid.default.result}"
  lifecycle {
    ignore_changes = [
      cors_rule,
    ]
  }
}


resource "alicloud_oss_bucket_cors" "default" {
  bucket        = alicloud_oss_bucket.CreateBucket.bucket
  response_vary = true
  cors_rule {
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
    allowed_headers = ["x-oss-test", "x-oss-abc"]
    expose_header   = ["x-oss-request-id"]
    max_age_seconds = "1000"
  }
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the Bucket.
* `cors_rule` - (Required) The Cross-Origin Resource Sharing (CORS) configuration of the Bucket. See [`cors_rule`](#cors_rule) below.
* `response_vary` - (Optional, Computed) Specifies whether to return the Vary: Origin header. Valid values: true: returns the Vary: Origin header, regardless of whether the request is a cross-origin request or whether the cross-origin request succeeds. false: does not return the Vary: Origin header. This element is valid only when at least one CORS rule is configured.

### `cors_rule`

The cors_rule supports the following:
* `allowed_headers` - (Optional) Specifies whether the headers specified by Access-Control-Request-Headers in the OPTIONS preflight request are allowed. You can use only one asterisk (*) as the wildcard for allowed header. .
* `allowed_methods` - (Required) The cross-origin request method that is allowed. Valid values: GET, PUT, DELETE, POST, and HEAD.
* `allowed_origins` - (Optional) The origins from which cross-origin requests are allowed. .
* `expose_header` - (Optional) The response headers for allowed access requests from applications, such as an XMLHttpRequest object in JavaScript. .
* `max_age_seconds` - (Optional) The period of time within which the browser can cache the response to an OPTIONS preflight request for the specified resource. Unit: seconds.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Cors.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Cors.
* `update` - (Defaults to 5 mins) Used when update the Bucket Cors.

## Import

OSS Bucket Cors can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_cors.example <id>
```