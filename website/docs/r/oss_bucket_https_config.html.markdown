---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_https_config"
description: |-
  Provides a Alicloud OSS Bucket Https Config resource.
---

# alicloud_oss_bucket_https_config

Provides a OSS Bucket Https Config resource. Whether the bucket can only be accessed with specific TLS versions.

For information about OSS Bucket Https Config and how to use it, see [What is Bucket Https Config](https://www.alibabacloud.com/help/en/oss/developer-reference/transport-layer-security).

-> **NOTE:** Available since v1.220.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_https_config&exampleId=8a2350a9-65c6-3aac-f7d3-d9f0af141a49f2978d0c&activeTab=example&spm=docs.r.oss_bucket_https_config.0.8a2350a965&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-${random_integer.default.result}"
}


resource "alicloud_oss_bucket_https_config" "default" {
  tls_versions = ["TLSv1.2"]
  bucket       = alicloud_oss_bucket.CreateBucket.bucket
  enable       = true
}
```

### Deleting `alicloud_oss_bucket_https_config` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_https_config`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oss_bucket_https_config&spm=docs.r.oss_bucket_https_config.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.
* `enable` - (Required) Specifies whether to enable TLS version management for the bucket. Valid values: true, false.
* `tls_versions` - (Optional) Specifies the TLS versions allowed to access this buckets.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Https Config.
* `update` - (Defaults to 5 mins) Used when update the Bucket Https Config.

## Import

OSS Bucket Https Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_https_config.example <id>
```