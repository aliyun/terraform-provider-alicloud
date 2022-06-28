---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_service"
sidebar_current: "docs-alicloud-datasource-cloud-sso-service"
description: |-
  Provides a datasource to open Cloud Sso Service automatically.
---

# alicloud\_cloud\_sso\_service

Using this data source can open Cloud Sso Service automatically.

For information about Cloud SSO and how to use it, see [What is Cloud SSO](https://www.alibabacloud.com/help/en/doc-detail/262819.html).

-> **NOTE:** Available in v1.148.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_sso_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required) Setting the value to `On` to enable the service. Valid values: `On` or `Off`. 

-> **NOTE:** Setting `enable = "On"` to open the Cloud Sso service that means you have read and agreed the [Cloud Sso Terms of Service](https://help.aliyun.com/document_detail/299998.html). When there is no directory in Cloud SSO, you can set `enable = "Off"` to turn off Cloud SSO as needed. After it is closed, you can also open it at any time.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
