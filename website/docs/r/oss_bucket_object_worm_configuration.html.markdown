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