---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_user_tenants"
sidebar_current: "docs-alicloud-datasource-dms-user-tenants"
description: |-
    Provides a list of available DMS Enterprise User Tenants.
---

# alicloud\_dms\_user\_tenants

This data source provides a list of DMS User Tenants in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.161.0+

## Example Usage

```terraform
# Declare the data source
data "alicloud_dms_user_tenants" "default" {
  status = "ACTIVE"
}

output "tid" {
  value = data.alicloud_dms_user_tenants.default.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Optional) The status of the user tenant.
* `ids` - (Optional) A list of DMS User Tenant IDs (TID).
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of DMS User Tenant IDs (UID).
* `names` - A list of DMS User Tenant names.
* `tenants` - A list of DMS User Tenants. Each element contains the following attributes:
  * `id` - The user tenant id.
  * `tid` - The user tenant id. Same as id.
  * `tenant_name` - The name of the user tenant.
  * `status` - The status of the user tenant.
