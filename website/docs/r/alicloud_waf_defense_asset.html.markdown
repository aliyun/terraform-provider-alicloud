---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_defense_asset"
description: |-
  Provides a Alicloud WAF Defense Asset resource.
---

# alicloud_waf_defense_asset

Provides a WAF Defense Asset resource.

A defense asset is a protected target fronted by a WAF instance — an AI prompt endpoint, an API-sec API, an account-protector login surface, an SSO endpoint, or a client-injection target — together with the match conditions that select the traffic to protect and a type-dependent configuration.

For information about WAF Defense Asset and how to use it, see [CreateDefenseAsset](https://next.api.alibabacloud.com/document/waf-openapi/2021-10-01/CreateDefenseAsset).

-> **NOTE:** Available since v1.257.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tfexample"
}

data "alicloud_wafv3_instances" "default" {
}

data "alicloud_wafv3_domains" "default" {
  instance_id    = data.alicloud_wafv3_instances.default.ids.0
  enable_details = true
}

resource "alicloud_waf_defense_asset" "default" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  asset_type   = "AI"
  # `resource` must reference a real protection object (a domain added to the
  # WAF instance); resolve it from the alicloud_wafv3_domains data source.
  resource = data.alicloud_wafv3_domains.default.domains.0.domain_id
  conditions = jsonencode([{
    Domain     = data.alicloud_wafv3_domains.default.domains.0.domain
    UrlPath    = "/"
    HttpMethod = "GET"
  }])
  asset_config = jsonencode({
    PromptPosition   = "request"
    InjectionRatio   = 1
    IsPayment        = false
    SensitiveLevel   = "N"
    OperationPurpose = "tf-example"
  })
  asset_name   = "tf-example-defense-asset"
  description  = "Managed by Terraform"
  asset_status = "1"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Web Application Firewall (WAF) instance that owns the defense asset. Required on every DefenseAsset call; changing it recreates the resource.
* `asset_type` - (Required, ForceNew) The type of the protective asset. CreateDefenseAsset validates this against a fixed enum and rejects unknown values. Changing it recreates the resource. Valid values: `AI`, `Client_Injection`, `Apisec`, `Account_Protector`, `SSO`.
* `resource` - (Required, ForceNew) The protected resource identifier (a domain added to the WAF instance) bound to the asset at creation. Required by CreateDefenseAsset and not accepted by ModifyDefenseAsset, so changing it recreates the resource.
* `conditions` - (Required, ForceNew) A JSON array of match conditions that select the traffic to protect. Each element is an object with `Domain` (string, required), `UrlPath` (string, required) and `HttpMethod` (string, optional). Required by CreateDefenseAsset and not accepted by ModifyDefenseAsset, so changing it recreates the resource. The value is stored as canonical JSON; formatting-only differences are not treated as changes.
* `asset_config` - (Required) The asset configuration as a JSON object. The shape depends on `asset_type`; it is required by both CreateDefenseAsset and ModifyDefenseAsset and is updated in place. The value is stored as canonical JSON; formatting-only differences are not treated as changes. For `asset_type = "AI"`, the following properties are accepted:
  - `PromptPosition` - (Optional) The position of the prompt to inspect.
  - `ResponseContents` - (Optional) The response content rules.
  - `InjectionRatio` - (Optional, Int) The injection detection ratio. Valid values: 1 to 100.
  - `IsPayment` - (Optional, Bool) Whether the asset is a payment scenario.
  - `SensitiveLevel` - (Optional) The sensitive data level. Valid values: `N`, `L1`, `L2`, `L3`.
  - `OperationPurpose` - (Optional) The purpose of the operation.
  - `SuccessConditions` - (Optional) The success match conditions.
  - `FailureConditions` - (Optional) The failure match conditions.
  - `ExtractAccounts` - (Optional) The account extraction policies.
  - `CookieDomain` - (Optional) The cookie domain.
  - `Swagger` - (Optional) The OpenAPI/Swagger specification (used when `asset_type = "Apisec"`).
* `asset_name` - (Optional) A human-readable name for the asset. Updatable through ModifyDefenseAsset.
* `description` - (Optional) A description for the asset. Updatable through ModifyDefenseAsset.
* `asset_status` - (Optional, Computed) The protection status of the asset. `1` means enabled, `0` means disabled. Defaults to `1` (enabled) on creation. Toggled through ModifyDefenseAssetStatus.

## Attributes Reference

The following attributes are exported:

* `id` - The AssetId of the defense asset, assigned by CreateDefenseAsset. The import id format is `<InstanceId>:<AssetId>`.

## Timeouts

The `timeouts` block allows you to specify timeouts for the `create`, `update` and `delete` operations:

* `create` - (Defaults to 5 minutes) Used when creating the defense asset.
* `update` - (Defaults to 5 minutes) Used when updating the defense asset.
* `delete` - (Defaults to 5 minutes) Used when deleting the defense asset.

## Import

WAF Defense Asset can be imported using the id with format `<InstanceId>:<AssetId>`, e.g.

```shell
$ terraform import alicloud_waf_defense_asset.example waf_v2_public_cn-xxx:asset-xxx
```

-> **NOTE:** `asset_config` and `conditions` are opaque caller-owned JSON and are not refreshed on import; supply them in config after import to avoid a diff, or accept the generated diff.
