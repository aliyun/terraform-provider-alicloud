---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_alert_resource"
sidebar_current: "docs-alicloud-resource-log-alert-resource"
description: |-
    Provides a resource to init SLS Alert resources automatically.
---

# alicloud_log_alert_resource

Using this resource can init SLS Alert resources automatically.

For information about SLS Alert and how to use it, see [SLS Alert Overview](https://www.alibabacloud.com/help/en/doc-detail/209202.html)

-> **NOTE:** Available since v1.219.0.

## Example Usage

```terraform
resource "alicloud_log_alert_resource" "example_user" {
  type = "user"
  lang = "cn"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, ForceNew) The type of alert resources, must be user or project, 'user' for init aliyuncloud account's alert center resource, including project named sls-alert-{uid}-{region} and some dashboards; 'project' for init project's alert resource, including logstore named internal-alert-history and alert dashboard.
* `lang` - (Optional, ForceNew) The lang of alert center resource when type is user.
* `project` - (Optional, ForceNew) The project of alert resource when type is project.

## Import

Log alert resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_alert_resource.example alert_resource:project:tf-project
```
