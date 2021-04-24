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

-> **NOTE:** Available in v1.122.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_config_aggregator" "example" {
    aggregator_accounts {
        account_id   = "123968452689****"
        account_name = "tf-testacc1"
        account_type = "ResourceDirectory"
    }
  aggregator_name = "example_value"
  description     = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `aggregator_accounts` - (Required) The information of account in aggregator.
* `aggregator_name` - (Required) The name of aggregator.
* `aggregator_type` - (Optional, ForceNew) The type of aggregator. Valid values: `CUSTOM`, `RD`.
* `description` - (Required) The description of aggregator.

#### Block aggregator_accounts

The aggregator_accounts supports the following: 

* `account_id` - (Required) Aggregator account Uid.
* `account_name` - (Required) Aggregator account name.
* `account_type` - (Required) Aggregator account source type. Valid values: `ResourceDirectory`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Aggregator.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Aggregator.

## Import

Cloud Config Aggregator can be imported using the id, e.g.

```
$ terraform import alicloud_config_aggregator.example <id>
```