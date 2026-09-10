---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_object_worm_configuration"
description: |-
  Provides a Alicloud OSS Bucket Object Worm Configuration resource.
---

# alicloud_oss_bucket_object_worm_configuration

Provides a OSS Bucket Object Worm Configuration resource.

Stores the Object-level compliant retention policy configuration for a bucket.

For information about OSS Bucket Object Worm Configuration and how to use it, see [What is Bucket Object Worm Configuration](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutBucketObjectWormConfiguration).

-> **NOTE:** Available since v1.278.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_object_worm_configuration&exampleId=d42deea3-eac1-7c24-35bb-48343c07f71de4f890c2&activeTab=example&spm=docs.r.oss_bucket_object_worm_configuration.0.d42deea3ea&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}

resource "alicloud_oss_bucket" "defaultQf8G0L" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket_versioning" "defaultosxikW" {
  status = "Enabled"
  bucket = alicloud_oss_bucket.defaultQf8G0L.id
}


resource "alicloud_oss_bucket_object_worm_configuration" "default" {
  bucket_name         = alicloud_oss_bucket.defaultQf8G0L.id
  object_worm_enabled = "Enabled"
  rule {
    default_retention {
      mode = "COMPLIANCE"
      days = "1"
    }
  }
}
```

### Deleting `alicloud_oss_bucket_object_worm_configuration` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_object_worm_configuration`. Terraform will remove this resource from the state file, however resources may remain.


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oss_bucket_object_worm_configuration&spm=docs.r.oss_bucket_object_worm_configuration.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `bucket_name` - (Required, ForceNew) Bucket name  
* `object_worm_enabled` - (Required) Specifies whether to enable the object-level compliance retention policy configuration.  
* `rule` - (Optional, Set) Container that stores the list of retention policies.   See [`rule`](#rule) below.

### `rule`

The rule supports the following:
* `default_retention` - (Optional, List) Container for the default retention policy.   See [`default_retention`](#rule-default_retention) below.

### `rule-default_retention`

The rule-default_retention supports the following:
* `days` - (Optional, Int) The number of days for compliant retention. This parameter is mutually exclusive with the Years parameter; only one of them can be specified.
* `mode` - (Optional) Compliance retention mode.  
* `years` - (Optional, Int) Default retention period in years. Valid values: 1 to 100. You can specify either Days or Years, but not both.  

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Object Worm Configuration.
* `update` - (Defaults to 5 mins) Used when update the Bucket Object Worm Configuration.

## Import

OSS Bucket Object Worm Configuration can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_object_worm_configuration.example <bucket_name>
```