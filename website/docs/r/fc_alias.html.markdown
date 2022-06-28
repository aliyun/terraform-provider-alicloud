---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_alias"
sidebar_current: "docs-alicloud-resource-fc"
description: |-
  Provides an Alicloud Function Compute Alias resource. 
---

# alicloud\_fc\_alias

Creates a Function Compute service alias. Creates an alias that points to the specified Function Compute service version. 
 For the detailed information, please refer to the [developer guide](https://www.alibabacloud.com/help/doc-detail/171635.htm).

-> **NOTE:** Available in 1.104.0+


## Example Usage

Basic Usage

```terraform
resource "alicloud_fc_alias" "example" {
  alias_name      = "my_alias"
  description     = "a sample description"
  service_name    = "my_service_name"
  service_version = "1"

  routing_config {
    additional_version_weights = {
      "2" = 0.5
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, ForceNew) Name for the alias you are creating. 
* `description` - (Optional) Description of the alias.
* `service_name` - (Required, ForceNew) The Function Compute service name.
* `service_version` - (Required) The Function Compute service version for which you are creating the alias. Pattern: (LATEST|[0-9]+).
* `routing_config` - (Optional) The Function Compute alias' route configuration settings. Fields documented below.

**routing_config** includes the following arguments:

* `additional_version_weights` - (Optional) A map that defines the proportion of events that should be sent to different versions of a Function Compute service.


## Import

Function Compute alias can be imported using the id, e.g.

```
$ terraform import alicloud_fc_alias.example my_alias_id
```
