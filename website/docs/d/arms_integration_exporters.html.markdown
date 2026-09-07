---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_integration_exporters"
sidebar_current: "docs-alicloud-datasource-arms-integration-exporters"
description: |-
  Provides a list of Arms Integration Exporters to the user.
---

# alicloud\_arms\_integration\_exporters

This data source provides the Arms Integration Exporters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.203.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_arms_integration_exporters" "ids" {
  ids              = ["example_id"]
  cluster_id       = "your_cluster_id"
  integration_type = "kafka"
}

output "arms_integration_exporters_id_1" {
  value = data.alicloud_arms_integration_exporters.ids.integration_exporters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Integration Exporter IDs.
* `cluster_id` - (Required, ForceNew) The ID of the Prometheus instance.
* `integration_type` - (Required, ForceNew) The type of prometheus integration.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `integration_exporters` - A list of Integration Exporters. Each element contains the following attributes:
  * `id` - The ID of the Integration Exporter. It formats as `<cluster_id>:<integration_type>:<instance_id>`.
  * `cluster_id` - The ID of the Prometheus instance.
  * `integration_type` - The type of prometheus integration.
  * `instance_id` - The ID of the Integration Exporter instance.
  * `param` - Exporter configuration parameter json string.
  * `instance_name` - The name of the instance.
  * `exporter_type` - Integration Exporter Type.
  * `target` - Monitor the target address.
  * `version` - The version information.
