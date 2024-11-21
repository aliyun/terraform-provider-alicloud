---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_er"
sidebar_current: "docs-alicloud-resource-dcdn-er"
description: |-
  Provides a Alicloud DCDN Er resource.
---

# alicloud_dcdn_er

Provides a DCDN Er resource.

For information about DCDN Er and how to use it, see [What is Er](https://www.alibabacloud.com/help/en/dcdn/developer-reference/api-dcdn-2018-01-15-createroutine).

-> **NOTE:** Available since v1.201.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dcdn_er&exampleId=945bf35c-6a86-ef0e-3deb-1eb09c4a7299b9328d44&activeTab=example&spm=docs.r.dcdn_er.0.945bf35c6a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_dcdn_er" "default" {
  er_name     = var.name
  description = var.name
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
* `env_conf` - (Optional) The configurations of the specified environment. See [`env_conf`](#env_conf) below.

### `env_conf`

The env_conf supports the following:

* `staging` - (Optional) The configuration of a staging environment. See [`staging`](#env_conf-staging) below.
* `production` - (Optional) The configuration of a production environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_anhui` - (Optional) The configuration of a presetCanaryAnhui environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_beijing` - (Optional) The configuration of a presetCanaryBeijing environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_chongqing` - (Optional) The configuration of a presetCanaryChongqing environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_fujian` - (Optional) The configuration of a presetCanaryFujian environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_gansu` - (Optional) The configuration of a presetCanaryGansu environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_guangdong` - (Optional) The configuration of a presetCanaryGuangdong environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_guangxi` - (Optional) The configuration of a presetCanaryGuangxi environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_guizhou` - (Optional) The configuration of a presetCanaryGuizhou environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_hainan` - (Optional) The configuration of a presetCanaryHainan environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_hebei` - (Optional) The configuration of a presetCanaryHebei environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_heilongjiang` - (Optional) The configuration of a presetCanaryHeilongjiang environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_henan` - (Optional) The configuration of a presetCanaryHenan environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_hong_kong` - (Optional) The configuration of a presetCanaryHongKong environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_hubei` - (Optional) The configuration of a presetCanaryHubei environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_hunan` - (Optional) The configuration of a presetCanaryHunan environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_jiangsu` - (Optional) The configuration of a presetCanaryJiangsu environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_jiangxi` - (Optional) The configuration of a presetCanaryJiangxi environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_jilin` - (Optional) The configuration of a presetCanaryJilin environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_liaoning` - (Optional) The configuration of a presetCanaryLiaoning environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_macau` - (Optional) The configuration of a presetCanaryMacau environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_neimenggu` - (Optional) The configuration of a presetCanaryNeimenggu environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_ningxia` - (Optional) The configuration of a presetCanaryNingxia environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_qinghai` - (Optional) The configuration of a presetCanaryQinghai environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_shaanxi` - (Optional) The configuration of a presetCanaryShaanxi environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_shandong` - (Optional) The configuration of a presetCanaryShandong environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_shanghai` - (Optional) The configuration of a presetCanaryShanghai environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_shanxi` - (Optional) The configuration of a presetCanaryShanxi environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_sichuan` - (Optional) The configuration of a presetCanarySichuan environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_taiwan` - (Optional) The configuration of a presetCanaryTaiwan environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_tianjin` - (Optional) The configuration of a presetCanaryTianjin environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_xinjiang` - (Optional) The configuration of a presetCanaryXinjiang environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_xizang` - (Optional) The configuration of a presetCanaryXizang environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_yunnan` - (Optional) The configuration of a presetCanaryYunnan environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_zhejiang` - (Optional) The configuration of a presetCanaryZhejiang environment. See [`staging`](#env_conf-staging) below.
* `preset_canary_overseas` - (Optional) The configuration of a presetCanaryOverseas environment. See [`staging`](#env_conf-staging) below.

### `env_conf-staging`

The staging supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Er.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Er.
* `update` - (Defaults to 5 mins) Used when update the Er.
* `delete` - (Defaults to 5 mins) Used when delete the Er.

## Import

DCDN Er can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_er.example <id>
```
