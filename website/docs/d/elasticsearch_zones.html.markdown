---
subcategory: "Elasticsearch"
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch_zones"
sidebar_current: "docs-alicloud-datasource-elasticsearch-zones"
description: |-
    Provides a list of availability zones for Elasticsearch that can be used by an Alibaba Cloud account.
---

# alicloud\_elasticsearch\_zones

This data source provides availability zones for Elasticsearch that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in v1.73.0+.

## Example Usage

```
# Declare the data source
data "alicloud_elasticsearch_zones" "zones_ids" {}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch Elasticsearch instances.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.
