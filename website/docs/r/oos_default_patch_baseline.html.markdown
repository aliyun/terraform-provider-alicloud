---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_default_patch_baseline"
sidebar_current: "docs-alicloud-resource-oos-default-patch-baseline"
description: |-
  Provides a Alicloud Oos Default Patch Baseline resource.
---

# alicloud_oos_default_patch_baseline

Provides a Oos Default Patch Baseline resource.

For information about Oos Default Patch Baseline and how to use it, see [What is Default Patch Baseline](https://www.alibabacloud.com/help/en/operation-orchestration-service/latest/api-oos-2019-06-01-registerdefaultpatchbaseline).

-> **NOTE:** Available since v1.203.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oos_default_patch_baseline&exampleId=af10cc02-6ad1-def3-e332-f5c2afd8af66cd876fb4&activeTab=example&spm=docs.r.oos_default_patch_baseline.0.af10cc026a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_oos_patch_baseline" "default" {
  operation_system    = "Windows"
  patch_baseline_name = "terraform-example"
  description         = "terraform-example"
  approval_rules      = "{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Key\":\"PatchSet\",\"Values\":[\"OS\"]},{\"Key\":\"ProductFamily\",\"Values\":[\"Windows\"]},{\"Key\":\"Product\",\"Values\":[\"Windows 10\",\"Windows 7\"]},{\"Key\":\"Classification\",\"Values\":[\"Security Updates\",\"Updates\",\"Update Rollups\",\"Critical Updates\"]},{\"Key\":\"Severity\",\"Values\":[\"Critical\",\"Important\",\"Moderate\"]}],\"ApproveAfterDays\":7,\"EnableNonSecurity\":true,\"ComplianceLevel\":\"Medium\"}]}"
}
resource "alicloud_oos_default_patch_baseline" "default" {
  patch_baseline_name = alicloud_oos_patch_baseline.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oos_default_patch_baseline&spm=docs.r.oos_default_patch_baseline.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `patch_baseline_name` - (Required,ForceNew) The name of the patch baseline.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `patch_baseline_id` - The ID of the patch baseline.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Default Patch Baseline.
* `delete` - (Defaults to 5 mins) Used when delete the Default Patch Baseline.

## Import

Oos Default Patch Baseline can be imported using the id, e.g.

```shell
$ terraform import alicloud_oos_default_patch_baseline.example <id>
```