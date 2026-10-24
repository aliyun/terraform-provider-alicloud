---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_alert_resource"
sidebar_current: "docs-alicloud-datasource-log-alert-resource"
description: |-
    Provides a datasource to init SLS Alert resources automatically.
---

# alicloud_log_alert_resource

Using this data source can init SLS Alert resources automatically.

For information about SLS Alert and how to use it, see [SLS Alert Overview](https://www.alibabacloud.com/help/en/doc-detail/209202.html)

-> **DEPRECATED:**  This resource  has been deprecated from version `1.219.0`. Please use new resource [alicloud_log_alert_resource](https://www.terraform.io/docs/providers/alicloud/r/log_alert_resource).

-> **NOTE:** Available since v1.161.0.

## Example Usage

```terraform
data "alicloud_log_alert_resource" "example_user" {
  type = "user"
  lang = "cn"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) The type of alert resources, must be user or project, 'user' for init the user-level alert center notification resources, including alert policies, action policies and content templates, it does not create the project named sls-alert-{uid}-{region}, which can only be initialized via the SLS console; 'project' for init the project's alert history logstore named internal-alert-history.
* `lang` - (Optional) The lang of alert center resource when type is user.
* `project` - (Optional) The project of alert resource when type is project.
