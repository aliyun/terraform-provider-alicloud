---
subcategory: "Governance"
layout: "alicloud"
page_title: "Alicloud: alicloud_governance_service_activation"
description: |-
  Provides a Alicloud Governance Service Activation resource.
---

# alicloud_governance_service_activation

Provides a Governance Service Activation resource.

Manages Cloud Governance Center activation for the current account. Cloud Governance Center is activated once per account in a single activation region, and `region` selects that region. Creating the resource activates the service and waits until its status becomes `Enabled`. Deleting it calls `DecommissionGovernance`, which closes Cloud Governance Center for the account, and waits until its status becomes `Disabled`.

For information about Governance Service Activation and how to use it, see [What is Service Activation](https://next.api.aliyun.com/document/governance/2021-01-20/OpenGovernanceService).

-> **NOTE:** Available since v2.1.1.

-> **NOTE:** This resource manages the account-level activation switch. Destroying it deactivates Cloud Governance Center rather than only removing the Terraform state entry. Use a single instance per account; do not manage the same activation with multiple configurations or concurrently with other tooling.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_governance_service_activation" "default" {
  region = "cn-hangzhou"
}
```

## Argument Reference

The following arguments are supported:
* `region` - (Required, ForceNew, String) The region in which Cloud Governance Center is activated for the account.
  - It must be the activation region: `OpenGovernanceService`, `DecommissionGovernance` and `GetGovernanceServiceStatus` all require `RegionId` to name the region where the service is activated, and requests naming another region fail.
  - It must be the same as the `region` of the provider configuration, because the activation is requested through that region.
  - It is used as the resource ID, so changing it forces a new resource.

## Attributes Reference

The following attributes are exported:
* `id` - (String) The ID of the resource supplied above. It equals `region`.
* `status` - (String) The activation status of Cloud Governance Center for the current account and region, returned by `GetGovernanceServiceStatus`. Valid values: `Enabled`, `Disabled`, `EnableFailed`, `DisableFailed`, `Enabling`, `Disabling`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 20 mins) Used when create the Service Activation.
* `delete` - (Defaults to 30 mins) Used when delete the Service Activation.

## Import

Governance Service Activation can be imported using the region, e.g.

```shell
$ terraform import alicloud_governance_service_activation.example cn-hangzhou
```

The import ID is the activation region, so the provider, or the provider alias applied to this resource, must be configured with that same region. If the configuration currently names another region, set it back to the activation region before importing: the activation cannot move between regions, and an import whose ID is the activation region fails while the provider targets a different region. Moving the resource to another provider alias or configuration keeps the same `region` value.
