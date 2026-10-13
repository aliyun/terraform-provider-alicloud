---
subcategory: "IMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ims_default_domain"
sidebar_current: "docs-alicloud-datasource-ims-default-domain"
description: |-
    Provides the default domain of the Alibaba Cloud account.
---

# alicloud_ims_default_domain

This data source provides the default domain of the Alibaba Cloud account. Every RAM user's
logon name is `<user_principal_name>@<default_domain>`.

-> **NOTE:** Available since v2.0.0-beta3.

## Example Usage

```terraform
data "alicloud_ims_default_domain" "default" {
}

output "default_domain" {
  value = data.alicloud_ims_default_domain.default.default_domain
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

The following attributes are exported:

* `id` - The default domain, same as `default_domain`.
* `default_domain` - The default domain of the Alibaba Cloud account, for example `examplecompany.onaliyun.com`.
