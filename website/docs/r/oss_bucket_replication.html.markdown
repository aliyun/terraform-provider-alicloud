---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_replication"
sidebar_current: "docs-alicloud-resource-oss-bucket-replication"
description: |-
  Provides a OSS bucket replication configuration resource.
---

# alicloud_oss_bucket_replication

Provides an independent replication configuration resource for OSS bucket.

For information about OSS replication and how to use it, see [What is cross-region replication](https://www.alibabacloud.com/help/doc-detail/31864.html) and [What is same-region replication](https://www.alibabacloud.com/help/doc-detail/254865.html).

-> **NOTE:** Available since v1.161.0.

## Example Usage

Set bucket replication configuration

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_replication&exampleId=064a8772-d0d7-d8be-a92e-d62443ea11cc1ae7f652&activeTab=example&spm=docs.r.oss_bucket_replication.0.064a8772d0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket_src" {
  bucket = "example-src-${random_integer.default.result}"
}

resource "alicloud_oss_bucket" "bucket_dest" {
  bucket = "example-dest-${random_integer.default.result}"
}

resource "alicloud_ram_role" "role" {
  name        = "example-role-${random_integer.default.result}"
  document    = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "oss.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  	EOF
  description = "this is a test"
  force       = true
}

resource "alicloud_ram_policy" "policy" {
  policy_name     = "example-policy-${random_integer.default.result}"
  policy_document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"*"
			  ],
			  "Effect": "Allow",
			  "Resource": [
				"*"
			  ]
			}
		  ],
			"Version": "1"
		}
		EOF
  description     = "this is a policy test"
  force           = true
}

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.policy_name
  policy_type = alicloud_ram_policy.policy.type
  role_name   = alicloud_ram_role.role.name
}

resource "alicloud_kms_key" "key" {
  description            = "Hello KMS"
  pending_window_in_days = "7"
  status                 = "Enabled"
}

resource "alicloud_oss_bucket_replication" "cross-region-replication" {
  bucket                        = alicloud_oss_bucket.bucket_src.id
  action                        = "PUT,DELETE"
  historical_object_replication = "enabled"
  prefix_set {
    prefixes = ["prefix1/", "prefix2/"]
  }
  destination {
    bucket   = alicloud_oss_bucket.bucket_dest.id
    location = alicloud_oss_bucket.bucket_dest.location
  }
  sync_role = alicloud_ram_role.role.name
  encryption_configuration {
    replica_kms_key_id = alicloud_kms_key.key.id
  }
  source_selection_criteria {
    sse_kms_encrypted_objects {
      status = "Enabled"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the bucket.
* `prefix_set` - (Optional, ForceNew) The prefixes used to specify the object to replicate. Only objects that match the prefix are replicated to the destination bucket. See [`prefix_set`](#prefix_set) below.
* `destination` - (Required, ForceNew) Specifies the destination for the rule. See [`destination`](#destination) below.
* `action` - (Optional, ForceNew) The operations that can be synchronized to the destination bucket. You can set action to one or more of the following operation types. Valid values: `ALL`(contains PUT, DELETE, and ABORT), `PUT`, `DELETE` and `ABORT`. Defaults to `ALL`.    
* `historical_object_replication` - (Optional, ForceNew) Specifies whether to replicate historical data from the source bucket to the destination bucket before data replication is enabled. Can be `enabled` or `disabled`. Defaults to `enabled`.
* `sync_role` - (Optional, ForceNew) Specifies the role that you authorize OSS to use to replicate data. If SSE-KMS is specified to encrypt the objects replicated to the destination bucket, it must be specified.
* `source_selection_criteria` - (Optional, ForceNew) Specifies other conditions used to filter the source objects to replicate. See [`source_selection_criteria`](#source_selection_criteria) below.
* `encryption_configuration` - (Optional, ForceNew) Specifies the encryption configuration for the objects replicated to the destination bucket. See [`encryption_configuration`](#encryption_configuration) below.
* `progress` - (Optional) Specifies the progress for querying the progress of a data replication task of a bucket.


### `prefix_set`

The prefix_set configuration block supports the following:

* `prefixes` - (Required, ForceNew) The list of object key name prefix identifying one or more objects to which the rule applies.

`NOTE`: The prefix must be less than or equal to 1024 characters in length.

### `destination`

The destination configuration block supports the following:

* `bucket` - (Required, ForceNew) The destination bucket to which the data is replicated.
* `location` - (Required, ForceNew) The region in which the destination bucket is located.
* `transfer_type` - (Optional, ForceNew) The link used to transfer data in data replication.. Can be `internal` or `oss_acc`. Defaults to `internal`.

`NOTE`: You can set transfer_type to oss_acc only when you create cross-region replication (CRR) rules.

### `source_selection_criteria`

The source_selection_criteria configuration block supports the following:

* `sse_kms_encrypted_objects` - (Optional, ForceNew) Filter source objects encrypted by using SSE-KMS. See [`sse_kms_encrypted_objects`](#source_selection_criteria-sse_kms_encrypted_objects) below.

### `source_selection_criteria-sse_kms_encrypted_objects`

The sse_kms_encrypted_objects configuration block supports the following:

* `status` - (Optional, ForceNew) Specifies whether to replicate objects encrypted by using SSE-KMS. Can be `Enabled` or `Disabled`.

### `encryption_configuration`

The encryption_configuration configuration block supports the following:

* `replica_kms_key_id` - (Required, ForceNew) The CMK ID used in SSE-KMS.

`NOTE`: If the status of sse_kms_encrypted_objects is set to Enabled, you must specify the replica_kms_key_id.

## Attributes Reference

The following attributes are exported:

* `id` - The current replication configuration resource ID. Composed of bucket name and rule_id with format `<bucket>:<rule_id>`.
* `rule_id` - The ID of the data replication rule.
* `status` - The status of the data replication task. Can be starting, doing and closing.
* `progress` - Retrieves the progress of the data replication task. This status is returned only when the data replication task is in the doing state.
    * `historical_object` - The percentage of the replicated historical data. This element is valid only when historical_object_replication is set to enabled.
    * `new_object` - The time used to distinguish new data from historical data. Data that is written to the source bucket before the time is replicated to the destination bucket as new data. The value of this element is in GMT.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `delete` - (Defaults to 30 mins) Used when delete a data replication rule (until the data replication task is cleared).

## Import

Oss Bucket Replication can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_replication.example
```
