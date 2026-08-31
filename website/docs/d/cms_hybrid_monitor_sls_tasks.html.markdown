---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_hybrid_monitor_sls_tasks"
sidebar_current: "docs-alicloud-datasource-cms-hybrid-monitor-sls-tasks"
description: |-
  Provides a list of Cms Hybrid Monitor Sls Tasks to the user.
---

# alicloud\_cms\_hybrid\_monitor\_sls\_tasks

This data source provides the Cms Hybrid Monitor Sls Tasks of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.179.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_hybrid_monitor_sls_tasks" "ids" {
  ids = ["example_value"]
}
output "cms_hybrid_monitor_sls_task_id_1" {
  value = data.alicloud_cms_hybrid_monitor_sls_tasks.ids.tasks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Hybrid Monitor Sls Task IDs.
* `keyword` - (Optional, ForceNew) The keyword that is used to search for metric import tasks.
* `namespace` - (Optional, ForceNew) The name of the namespace.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `tasks` - A list of Cms Hybrid Monitor Sls Tasks. Each element contains the following attributes:
    * `attach_labels` - The tags of the metric import task.
         * `name` - The key of the tag.
         * `value` - The value of the tag.
    * `collect_interval` - The interval between the cloud monitoring plug-in collecting host monitoring data.
    * `collect_target_endpoint` - The address where the cloudmonitor Plug-In collects the monitoring data of the host.
    * `collect_target_path` - When the cloud monitor Agent collects, the relative path of the collection.
    * `collect_target_type` - The type of the monitoring data. Valid values: Spring, Tomcat, Nginx, Tengine, JVM, Redis, MySQL, and AWS.
    * `collect_timout` - The timeout period for the cloudmonitor plug-in to collect host monitoring data.
    * `create_time` - Create the timestamp of the monitoring task. Unit: milliseconds.
    * `description` - Monitoring task description.
    * `extra_info` - Additional information for the instance.
    * `group_id` - The ID of the application Group.
    * `hybrid_monitor_sls_task_id` - The ID of the monitoring task.
    * `id` - The ID of the Hybrid Monitor Sls Task.
    * `instances` - A list of instances where monitoring data is collected in batches.
    * `log_file_path` - The path where on-premises log data is stored. On-premises log data is stored in the specified path of the host where CloudMonitor is deployed.
    * `log_process` - Local Log Monitoring and calculation method.
    * `log_sample` - The sample on-premises log.
    * `log_split` - The local log data is divided according to different matching patterns.
    * `match_express` - The matching condition of the instance in the application Group.
        * `function` - The method that is used to match the instance name.
        * `name` - The name of the instance.
        * `value` - The keyword that corresponds to the instance name.
    * `match_express_relation` - The filter condition of the instance of the monitoring task.
    * `namespace` - The namespace to which the host belongs.
    * `network_type` - The network type of the host.
    * `sls_process` - The configurations of the logs that are imported from Log Service.
    * `sls_process_config` - The configurations of the logs that are imported from Log Service.
    * `group_by` - The dimension based on which data is aggregated. This parameter is equivalent to the GROUP BY clause in SQL.
        * `alias` - The alias of the aggregation result.
        * `sls_key_name` - The name of the key that is used to aggregate logs imported from Log Service.
    * `statistics` - The method that is used to aggregate logs imported from Log Service.
        * `parameter_one` - The value of the function that is used to aggregate logs imported from Log Service.
        * `parameter_two` - The value of the function that is used to aggregate logs imported from Log Service.
        * `sls_key_name` - The name of the key that is used to aggregate logs imported from Log Service.
        * `alias` - The alias of the aggregation result.
        * `function` - The function that is used to aggregate log data within a statistical period.
    * `express` - The extended fields that specify the results of basic operations that are performed on aggregation results.
        * `alias` - The alias of the extended field that specifies the result of basic operations that are performed on aggregation results.
        * `express` - The extended field that specifies the result of basic operations that are performed on aggregation results.
    * `filter` - The conditions that are used to filter logs imported from Log Service.
        * `filters` - The conditions that are used to filter logs imported from Log Service.
            * `operator` - The method that is used to filter logs imported from Log Service.
            * `sls_key_name` - The name of the key that is used to filter logs imported from Log Service.
            * `value` - The value of the key that is used to filter logs imported from Log Service.
        * `relation` - The relationship between multiple filter conditions.
    * `task_name` - The name of the metric import task.
    * `task_type` - Monitoring Task type.
    * `upload_region` - The region where the host resides.