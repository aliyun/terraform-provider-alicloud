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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_alert_resource&exampleId=763a3da5-7c2f-f903-99ba-d4f5f2b0a1f3a7f6eb46&activeTab=example&spm=docs.r.log_alert_resource.0.763a3da57c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_log_alert_resource" "example_user" {
  type = "user"
  lang = "cn"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_log_alert_resource&spm=docs.r.log_alert_resource.example&intl_lang=EN_US)

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
