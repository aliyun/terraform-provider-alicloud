---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_hybrid_monitor_fc_task"
sidebar_current: "docs-alicloud-resource-cms-hybrid-monitor-fc-task"
description: |-
  Provides a Alicloud Cloud Monitor Service Hybrid Monitor Fc Task resource.
---

# alicloud_cms_hybrid_monitor_fc_task

Provides a Cloud Monitor Service Hybrid Monitor Fc Task resource.

For information about Cloud Monitor Service Hybrid Monitor Fc Task and how to use it, see [What is Hybrid Monitor Fc Task](https://www.alibabacloud.com/help/en/cloudmonitor/latest/createhybridmonitortask).

-> **NOTE:** Available since v1.179.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cms_hybrid_monitor_fc_task&exampleId=f78c1559-ffec-f89c-4fd3-84a9de1bc4012ba49755&activeTab=example&spm=docs.r.cms_hybrid_monitor_fc_task.0.f78c1559ff&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_account" "default" {}

resource "alicloud_cms_namespace" "default" {
  description   = var.name
  namespace     = var.name
  specification = "cms.s1.large"
}

resource "alicloud_cms_hybrid_monitor_fc_task" "default" {
  namespace      = alicloud_cms_namespace.default.id
  yarm_config    = <<EOF
---
products:
- namespace: "acs_ecs_dashboard"
  metric_info:
  - metric_list:
    - "CPUUtilization"
    - "DiskReadBPS"
    - "InternetOut"
    - "IntranetOut"
    - "cpu_idle"
    - "cpu_system"
    - "cpu_total"
    - "diskusage_utilization"
- namespace: "acs_rds_dashboard"
  metric_info:
  - metric_list:
    - "MySQL_QPS"
    - "MySQL_TPS"
EOF
  target_user_id = data.alicloud_account.default.id
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, ForceNew) The index warehouse where the host belongs.
* `target_user_id` - (Optional, ForceNew) The ID of the member account. If you call API operations by using a management account, you can connect the Alibaba Cloud services that are activated for a member account in Resource Directory to Hybrid Cloud Monitoring. You can use Resource Directory to monitor Alibaba Cloud services across enterprise accounts.
* `yarm_config` - (Required) The configuration file of the Alibaba Cloud service that you want to monitor by using Hybrid Cloud Monitoring.
  - `namespace`: the namespace of the Alibaba Cloud service.
  - `metric_list`: the metrics of the Alibaba Cloud service.
  
## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Hybrid Monitor Fc Task. The value formats as `<hybrid_monitor_fc_task_id>:<namespace>`.
* `hybrid_monitor_fc_task_id` - The ID of the monitoring task.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Hybrid Monitor Fc Task.
* `delete` - (Defaults to 2 mins) Used when delete the Hybrid Monitor Fc Task.
* `update` - (Defaults to 2 mins) Used when update the Hybrid Monitor Fc Task.

## Import

Cloud Monitor Service Hybrid Monitor Fc Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_hybrid_monitor_fc_task.example <hybrid_monitor_fc_task_id>:<namespace>
```