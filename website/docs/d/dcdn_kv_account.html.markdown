---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_kv_account"
sidebar_current: "docs-alicloud-datasource-dcdn-kv-account"
description: |-
  Provides a datasource to open the DCDN kv account automatically.
---

# alicloud_dcdn_kv_account

This data source provides DCDN kv account available to the user.[What is DCDN Kv Account](https://www.alibabacloud.com/help/en/dcdn/developer-reference/api-dcdn-2018-01-15-describedcdnkvaccount)

-> **NOTE:** Available since v1.198.0.

## Example Usage

```terraform
data "alicloud_dcdn_kv_account" "status" {
  status = "online"
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Optional, Computed, ForceNew) The status of the KV feature for your account. Valid values: `online`, `offline`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current kv account enable status. 
