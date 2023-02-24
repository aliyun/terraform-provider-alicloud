---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_er"
sidebar_current: "docs-alicloud-resource-dcdn-er"
description: |-
  Provides a Alicloud DCDN Er resource.
---

# alicloud\_dcdn\_er

Provides a DCDN Er resource.

For information about DCDN Er and how to use it, see [What is Er](https://www.alibabacloud.com/help/en/dynamic-route-for-cdn/latest/createroutine).

-> **NOTE:** Available in v1.201.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_dcdn_er" "default" {
  er_name     = "tf-example-name"
  description = "tf-example-description"
  env_conf {
    staging {
      spec_name     = "5ms"
      allowed_hosts = ["example.com"]
    }
    production {
      spec_name     = "5ms"
      allowed_hosts = ["example.com"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `er_name` - (Required, ForceNew) The name of the routine. The name must be unique among the routines that belong to the same Alibaba Cloud account.
* `description` - (Optional) Routine The description of the routine.
* `env_conf` - (Optional, Computed) The configurations of the specified environment. See the following `Block env_conf`.

### Block env_conf

The env_conf supports the following:

* `staging` - (Optional, Computed) The configuration of a staging environment. See the following `Block staging`.
* `production` - (Optional, Computed) The configuration of a production environment. The `Block production` is same as `Block staging`.
* `preset_canary_anhui` - (Optional, Computed) The configuration of a presetCanaryAnhui environment. The `Block preset_canary_anhui` is same as `Block staging`.
* `preset_canary_beijing` - (Optional, Computed) The configuration of a presetCanaryBeijing environment. The `Block preset_canary_beijing` is same as `Block staging`.
* `preset_canary_chongqing` - (Optional, Computed) The configuration of a presetCanaryChongqing environment. The `Block preset_canary_chongqing` is same as `Block staging`.
* `preset_canary_fujian` - (Optional, Computed) The configuration of a presetCanaryFujian environment. The `Block preset_canary_fujian` is same as `Block staging`.
* `preset_canary_gansu` - (Optional, Computed) The configuration of a presetCanaryGansu environment. The `Block preset_canary_gansu` is same as `Block staging`.
* `preset_canary_guangdong` - (Optional, Computed) The configuration of a presetCanaryGuangdong environment. The `Block preset_canary_guangdong` is same as `Block staging`.
* `preset_canary_guangxi` - (Optional, Computed) The configuration of a presetCanaryGuangxi environment. The `Block preset_canary_guangxi` is same as `Block staging`.
* `preset_canary_guizhou` - (Optional, Computed) The configuration of a presetCanaryGuizhou environment. The `Block preset_canary_guizhou` is same as `Block staging`.
* `preset_canary_hainan` - (Optional, Computed) The configuration of a presetCanaryHainan environment. The `Block preset_canary_hainan` is same as `Block staging`.
* `preset_canary_hebei` - (Optional, Computed) The configuration of a presetCanaryHebei environment. The `Block preset_canary_hebei` is same as `Block staging`.
* `preset_canary_heilongjiang` - (Optional, Computed) The configuration of a presetCanaryHeilongjiang environment. The `Block preset_canary_heilongjiang` is same as `Block staging`.
* `preset_canary_henan` - (Optional, Computed) The configuration of a presetCanaryHenan environment. The `Block preset_canary_henan` is same as `Block staging`.
* `preset_canary_hong_kong` - (Optional, Computed) The configuration of a presetCanaryHongKong environment. The `Block preset_canary_hong_kong` is same as `Block staging`.
* `preset_canary_hubei` - (Optional, Computed) The configuration of a presetCanaryHubei environment. The `Block preset_canary_hubei` is same as `Block staging`.
* `preset_canary_hunan` - (Optional, Computed) The configuration of a presetCanaryHunan environment. The `Block preset_canary_hunan` is same as `Block staging`.
* `preset_canary_jiangsu` - (Optional, Computed) The configuration of a presetCanaryJiangsu environment. The `Block preset_canary_jiangsu` is same as `Block staging`.
* `preset_canary_jiangxi` - (Optional, Computed) The configuration of a presetCanaryJiangxi environment. The `Block preset_canary_jiangxi` is same as `Block staging`.
* `preset_canary_jilin` - (Optional, Computed) The configuration of a presetCanaryJilin environment. The `Block preset_canary_jilin` is same as `Block staging`.
* `preset_canary_liaoning` - (Optional, Computed) The configuration of a presetCanaryLiaoning environment. The `Block preset_canary_liaoning` is same as `Block staging`.
* `preset_canary_macau` - (Optional, Computed) The configuration of a presetCanaryMacau environment. The `Block preset_canary_macau` is same as `Block staging`.
* `preset_canary_neimenggu` - (Optional, Computed) The configuration of a presetCanaryNeimenggu environment. The `Block preset_canary_neimenggu` is same as `Block staging`.
* `preset_canary_ningxia` - (Optional, Computed) The configuration of a presetCanaryNingxia environment. The `Block preset_canary_ningxia` is same as `Block staging`.
* `preset_canary_qinghai` - (Optional, Computed) The configuration of a presetCanaryQinghai environment. The `Block preset_canary_qinghai` is same as `Block staging`.
* `preset_canary_shaanxi` - (Optional, Computed) The configuration of a presetCanaryShaanxi environment. The `Block preset_canary_shaanxi` is same as `Block staging`.
* `preset_canary_shandong` - (Optional, Computed) The configuration of a presetCanaryShandong environment. The `Block preset_canary_shandong` is same as `Block staging`.
* `preset_canary_shanghai` - (Optional, Computed) The configuration of a presetCanaryShanghai environment. The `Block preset_canary_shanghai` is same as `Block staging`.
* `preset_canary_shanxi` - (Optional, Computed) The configuration of a presetCanaryShanxi environment. The `Block preset_canary_shanxi` is same as `Block staging`.
* `preset_canary_sichuan` - (Optional, Computed) The configuration of a presetCanarySichuan environment. The `Block preset_canary_sichuan` is same as `Block staging`.
* `preset_canary_taiwan` - (Optional, Computed) The configuration of a presetCanaryTaiwan environment. The `Block preset_canary_taiwan` is same as `Block staging`.
* `preset_canary_tianjin` - (Optional, Computed) The configuration of a presetCanaryTianjin environment. The `Block preset_canary_tianjin` is same as `Block staging`.
* `preset_canary_xinjiang` - (Optional, Computed) The configuration of a presetCanaryXinjiang environment. The `Block preset_canary_xinjiang` is same as `Block staging`.
* `preset_canary_xizang` - (Optional, Computed) The configuration of a presetCanaryXizang environment. The `Block preset_canary_xizang` is same as `Block staging`.
* `preset_canary_yunnan` - (Optional, Computed) The configuration of a presetCanaryYunnan environment. The `Block preset_canary_yunnan` is same as `Block staging`.
* `preset_canary_zhejiang` - (Optional, Computed) The configuration of a presetCanaryZhejiang environment. The `Block preset_canary_zhejiang` is same as `Block staging`.
* `preset_canary_overseas` - (Optional, Computed) The configuration of a presetCanaryOverseas environment. The `Block preset_canary_overseas` is same as `Block staging`.

#### Block staging

The staging supports the following:

* `spec_name` - (Optional, Computed) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional, Computed) Allowed DCDN domain names.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Er.

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Er.
* `update` - (Defaults to 5 mins) Used when update the Er.
* `delete` - (Defaults to 5 mins) Used when delete the Er.

## Import

DCDN Er can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_er.example <id>
```
