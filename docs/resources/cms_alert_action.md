---
subcategory: "Cloud Monitor (CMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_action"
description: |-
  Provides a Cms Alert Action resource.
---

# alicloud_cms_alert_action

Provides a Cms Alert Action resource.

Alert actions are integrations that define how alert notifications are delivered. Each alert action has a type (such as WEBHOOK, MNS, SLS, FC, ESS, PAGER_DUTY, EB, FC3) and corresponding configuration parameters.

For information about the Cms Alert Action API, see [List alert actions](https://www.alibabacloud.com/help/en/cloudmonitor/developer-reference/api-cms-2024-03-30-list-alert-actions).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_alert_action" "example" {
  alert_action_name = "tf-test-alert-action"
  type              = "WEBHOOK"
  webhook_param {
    method       = "POST"
    url          = "https://example.com/webhook"
    content_type = "JSON"
    headers = {
      X-Custom-Header = "value"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `alert_action_name` - (Required) The name of the alert action.
* `type` - (Required, ForceNew) The type of the alert action. Valid values: `FC`, `MNS`, `OPEN_API`, `SLS`, `ESS`, `PAGER_DUTY`, `WEBHOOK`, `EB`, `FC3`.
* `webhook_param` - (Optional) The webhook parameters. See `webhook_param` below.
* `mns_param` - (Optional) The message queue parameters. See `mns_param` below.
* `sls_param` - (Optional) The log service parameters. See `sls_param` below.
* `ess_param` - (Optional) The elastic scaling parameters. See `ess_param` below.
* `fc_param` - (Optional) The function compute parameters. See `fc_param` below.
* `pager_duty_param` - (Optional) The PagerDuty parameters. See `pager_duty_param` below.
* `fc3_param` - (Optional) The FC3.0 parameters. See `fc3_param` below.
* `eb_param` - (Optional) The EventBridge parameters. See `eb_param` below.

### webhook_param

* `method` - (Optional) The HTTP request method. Valid values: `GET`, `POST`.
* `url` - (Optional) The webhook URL.
* `content_type` - (Optional) The request content type. Valid values: `JSON`, `FORM`.
* `headers` - (Optional) The request headers.

### mns_param

* `mns_type` - (Optional) The MNS type. Valid values: `queue`, `topic`.
* `name` - (Optional) The queue or topic name.
* `region_id` - (Optional) The MNS region.

### sls_param

* `logstore` - (Optional) The log store name.
* `project` - (Optional) The log project name.
* `region_id` - (Optional) The SLS region.

### ess_param

* `ess_group_id` - (Optional) The auto scaling group ID.
* `ess_rule_id` - (Optional) The auto scaling rule ID.
* `region_id` - (Optional) The ESS region.

### fc_param

* `function` - (Optional) The function compute function name.
* `region_id` - (Optional) The FC region.
* `service` - (Optional) The function compute service name.

### pager_duty_param

* `key` - (Optional) The PagerDuty integration key.
* `url` - (Optional) The PagerDuty URL.

### fc3_param

* `region_id` - (Optional) The FC3.0 region.
* `function` - (Optional) The FC3.0 function name.
* `qualifier` - (Optional) The FC3.0 qualifier.

### eb_param

* `region_id` - (Optional) The EventBridge region.
* `event_bus_name` - (Optional) The event bus name.
* `subject` - (Optional) The event subject.
* `eb_source` - (Optional) The event source.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alert action.
* `alert_action_id` - The ID of the alert action.
* `region_id` - The region ID of the alert action.

## Import

Cms Alert Action can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_alert_action.example <alert-action-id>
```
