---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_prometheus_monitorings"
sidebar_current: "docs-alicloud-datasource-arms-prometheus-monitorings"
description: |-
  Provides a list of Arms Prometheus Monitorings to the user.
---

# alicloud_arms_prometheus_monitorings

This data source provides the Arms Prometheus Monitorings of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_arms_prometheus_monitorings" "ids" {
  cluster_id = "your_cluster_id"
  ids        = ["example_id"]
}

output "arms_prometheus_monitorings_id_1" {
  value = data.alicloud_arms_prometheus_monitorings.ids.prometheus_monitorings.0.id
}

data "alicloud_arms_prometheus_monitorings" "nameRegex" {
  cluster_id = "your_cluster_id"
  name_regex = "tf-example"
}

output "arms_prometheus_monitorings_id_2" {
  value = data.alicloud_arms_prometheus_monitorings.nameRegex.prometheus_monitorings.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Prometheus Monitoring IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Prometheus Monitoring name.
* `cluster_id` - (Required, ForceNew) The ID of the Prometheus instance.
* `type` - (Optional, ForceNew) The type of the monitoring configuration. Valid values: `serviceMonitor`, `podMonitor`, `customJob`, `probe`.
* `status` - (Optional, ForceNew) The status of the monitoring configuration. Valid values: `run`, `stop`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Prometheus Monitoring names.
* `prometheus_monitorings` - A list of Prometheus Monitorings. Each element contains the following attributes:
  * `id` - The ID of the Prometheus Monitoring. It formats as `<cluster_id>:<monitoring_name>:<type>`.
  * `cluster_id` - The ID of the Prometheus instance.
  * `monitoring_name` - The name of the monitoring configuration.
  * `type` - The type of the monitoring configuration.
  * `config_yaml` - The monitoring configuration. The value is a YAML string.
  * `status` - The status of the monitoring configuration.
