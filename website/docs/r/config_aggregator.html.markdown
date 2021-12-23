---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregator"
sidebar_current: "docs-alicloud-resource-config-aggregator"
description: |-
  Provides a Alicloud Cloud Config Aggregator resource.
---

# alicloud\_config\_aggregator

Provides a Cloud Config Aggregator resource.

For information about Cloud Config Aggregate Config Rule and how to use it, see [What is Aggregator](https://www.alibabacloud.com/help/en/doc-detail/211197.html).

-> **NOTE:** Available in v1.124.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_config_aggregator" "example" {
  aggregator_accounts {
    account_id   = "123968452689****"
    account_name = "tf-testacc1234"
    account_type = "ResourceDirectory"
  }
  aggregator_name = "tf-testaccConfigAggregator1234"
  description     = "tf-testaccConfigAggregator1234"
}

```

## Argument Reference

The following arguments are supported:

* `aggregator_accounts` - (Optional) The information of account in aggregator. If the aggregator_type is RD, it is optional and means add all members in the resource directory to the account group. **NOTE:** the field `aggregator_accounts` is not required from version 1.148.0.
* `aggregator_name` - (Required) The name of aggregator.
* `aggregator_type` - (Optional, ForceNew) The type of aggregator. Valid values: `CUSTOM`, `RD`. The Default value: `CUSTOM`.
  * `CUSTOM` - The custom account group.
  * `RD` - The global account group.
* `description` - (Required) The description of aggregator.

#### Block aggregator_accounts

The aggregator_accounts supports the following: 

* `account_id` - (Required) Aggregator account Uid.
* `account_name` - (Required) Aggregator account name.
* `account_type` - (Required) Aggregator account source type. Valid values: `ResourceDirectory`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Aggregator.
* `status` - The status of the resource. Valid values: `0`: creating `1`: normal `2`: deleting.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Aggregator.

## Import

Cloud Config Aggregator can be imported using the id, e.g.

```
$ terraform import alicloud_config_aggregator.example <id>
```
