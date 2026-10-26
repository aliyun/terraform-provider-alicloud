---
subcategory: "Governance"
layout: "alicloud"
page_title: "Alicloud: alicloud_governance_service"
sidebar_current: "docs-alicloud-datasource-governance-service"
description: |-
  Provides a datasource to open the Cloud Governance Center service automatically.
---

# alicloud_governance_service

Using this data source can open the Cloud Governance Center service automatically. If the service has been opened, it will return opened.[What is Cloud Governance Center](https://next.api.aliyun.com/document/governance/2021-01-20/OpenGovernanceService)

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_governance_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to open the Cloud Governance Center service. If has been opened, return the result. Valid values: `On` or `Off`. Defaults to `Off`.

  -> **NOTE:** Setting `enable = "On"` to open the Cloud Governance Center service means you have read and agreed the Cloud Governance Center Terms of Service. The account must complete enterprise real-name verification before opening; otherwise, the API returns `InvalidEnterpriseRealName.NotFound`. Closing the service is not supported by this data source.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status.
