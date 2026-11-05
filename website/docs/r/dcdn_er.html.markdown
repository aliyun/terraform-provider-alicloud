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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_dcdn_er&spm=docs.r.dcdn_er.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `er_name` - (Required, ForceNew) The name of the routine. The name must be unique among the routines that belong to the same Alibaba Cloud account.
* `description` - (Optional) Routine The description of the routine.
* `env_conf` - (Optional) The configurations of the specified environment. See [`env_conf`](#env_conf) below.

### `env_conf`

The env_conf supports the following:

* `staging` - (Optional) The configuration of a staging environment. See [`staging`](#env_conf-staging) below.
* `production` - (Optional) The configuration of a production environment. See [`production`](#env_conf-production) below.
* `preset_canary_anhui` - (Optional) The configuration of a presetCanaryAnhui environment. See [`preset_canary_anhui`](#env_conf-preset_canary_anhui) below.
* `preset_canary_beijing` - (Optional) The configuration of a presetCanaryBeijing environment. See [`preset_canary_beijing`](#env_conf-preset_canary_beijing) below.
* `preset_canary_chongqing` - (Optional) The configuration of a presetCanaryChongqing environment. See [`preset_canary_chongqing`](#env_conf-preset_canary_chongqing) below.
* `preset_canary_fujian` - (Optional) The configuration of a presetCanaryFujian environment. See [`preset_canary_fujian`](#env_conf-preset_canary_fujian) below.
* `preset_canary_gansu` - (Optional) The configuration of a presetCanaryGansu environment. See [`preset_canary_gansu`](#env_conf-preset_canary_gansu) below.
* `preset_canary_guangdong` - (Optional) The configuration of a presetCanaryGuangdong environment. See [`preset_canary_guangdong`](#env_conf-preset_canary_guangdong) below.
* `preset_canary_guangxi` - (Optional) The configuration of a presetCanaryGuangxi environment. See [`preset_canary_guangxi`](#env_conf-preset_canary_guangxi) below.
* `preset_canary_guizhou` - (Optional) The configuration of a presetCanaryGuizhou environment. See [`preset_canary_guizhou`](#env_conf-preset_canary_guizhou) below.
* `preset_canary_hainan` - (Optional) The configuration of a presetCanaryHainan environment. See [`preset_canary_hainan`](#env_conf-preset_canary_hainan) below.
* `preset_canary_hebei` - (Optional) The configuration of a presetCanaryHebei environment. See [`preset_canary_hebei`](#env_conf-preset_canary_hebei) below.
* `preset_canary_heilongjiang` - (Optional) The configuration of a presetCanaryHeilongjiang environment. See [`preset_canary_heilongjiang`](#env_conf-preset_canary_heilongjiang) below.
* `preset_canary_henan` - (Optional) The configuration of a presetCanaryHenan environment. See [`preset_canary_henan`](#env_conf-preset_canary_henan) below.
* `preset_canary_hong_kong` - (Optional) The configuration of a presetCanaryHongKong environment. See [`preset_canary_hong_kong`](#env_conf-preset_canary_hong_kong) below.
* `preset_canary_hubei` - (Optional) The configuration of a presetCanaryHubei environment. See [`preset_canary_hubei`](#env_conf-preset_canary_hubei) below.
* `preset_canary_hunan` - (Optional) The configuration of a presetCanaryHunan environment. See [`preset_canary_hunan`](#env_conf-preset_canary_hunan) below.
* `preset_canary_jiangsu` - (Optional) The configuration of a presetCanaryJiangsu environment. See [`preset_canary_jiangsu`](#env_conf-preset_canary_jiangsu) below.
* `preset_canary_jiangxi` - (Optional) The configuration of a presetCanaryJiangxi environment. See [`preset_canary_jiangxi`](#env_conf-preset_canary_jiangxi) below.
* `preset_canary_jilin` - (Optional) The configuration of a presetCanaryJilin environment. See [`preset_canary_jilin`](#env_conf-preset_canary_jilin) below.
* `preset_canary_liaoning` - (Optional) The configuration of a presetCanaryLiaoning environment. See [`preset_canary_liaoning`](#env_conf-preset_canary_liaoning) below.
* `preset_canary_macau` - (Optional) The configuration of a presetCanaryMacau environment. See [`preset_canary_macau`](#env_conf-preset_canary_macau) below.
* `preset_canary_neimenggu` - (Optional) The configuration of a presetCanaryNeimenggu environment. See [`preset_canary_neimenggu`](#env_conf-preset_canary_neimenggu) below.
* `preset_canary_ningxia` - (Optional) The configuration of a presetCanaryNingxia environment. See [`preset_canary_ningxia`](#env_conf-preset_canary_ningxia) below.
* `preset_canary_qinghai` - (Optional) The configuration of a presetCanaryQinghai environment. See [`preset_canary_qinghai`](#env_conf-preset_canary_qinghai) below.
* `preset_canary_shaanxi` - (Optional) The configuration of a presetCanaryShaanxi environment. See [`preset_canary_shaanxi`](#env_conf-preset_canary_shaanxi) below.
* `preset_canary_shandong` - (Optional) The configuration of a presetCanaryShandong environment. See [`preset_canary_shandong`](#env_conf-preset_canary_shandong) below.
* `preset_canary_shanghai` - (Optional) The configuration of a presetCanaryShanghai environment. See [`preset_canary_shanghai`](#env_conf-preset_canary_shanghai) below.
* `preset_canary_shanxi` - (Optional) The configuration of a presetCanaryShanxi environment. See [`preset_canary_shanxi`](#env_conf-preset_canary_shanxi) below.
* `preset_canary_sichuan` - (Optional) The configuration of a presetCanarySichuan environment. See [`preset_canary_sichuan`](#env_conf-preset_canary_sichuan) below.
* `preset_canary_taiwan` - (Optional) The configuration of a presetCanaryTaiwan environment. See [`preset_canary_taiwan`](#env_conf-preset_canary_taiwan) below.
* `preset_canary_tianjin` - (Optional) The configuration of a presetCanaryTianjin environment. See [`preset_canary_tianjin`](#env_conf-preset_canary_tianjin) below.
* `preset_canary_xinjiang` - (Optional) The configuration of a presetCanaryXinjiang environment. See [`preset_canary_xinjiang`](#env_conf-preset_canary_xinjiang) below.
* `preset_canary_xizang` - (Optional) The configuration of a presetCanaryXizang environment. See [`preset_canary_xizang`](#env_conf-preset_canary_xizang) below.
* `preset_canary_yunnan` - (Optional) The configuration of a presetCanaryYunnan environment. See [`preset_canary_yunnan`](#env_conf-preset_canary_yunnan) below.
* `preset_canary_zhejiang` - (Optional) The configuration of a presetCanaryZhejiang environment. See [`preset_canary_zhejiang`](#env_conf-preset_canary_zhejiang) below.
* `preset_canary_overseas` - (Optional) The configuration of a presetCanaryOverseas environment. See [`preset_canary_overseas`](#env_conf-preset_canary_overseas) below.

### `env_conf-staging`

The staging supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-production`

The production supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_anhui`

The preset_canary_anhui supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_beijing`

The preset_canary_beijing supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_chongqing`

The preset_canary_chongqing supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_fujian`

The preset_canary_fujian supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_gansu`

The preset_canary_gansu supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_guangdong`

The preset_canary_guangdong supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_guangxi`

The preset_canary_guangxi supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_guizhou`

The preset_canary_guizhou supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_hainan`

The preset_canary_hainan supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_hebei`

The preset_canary_hebei supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_heilongjiang`

The preset_canary_heilongjiang supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_henan`

The preset_canary_henan supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_hong_kong`

The preset_canary_hong_kong supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_hubei`

The preset_canary_hubei supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_hunan`

The preset_canary_hunan supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_jiangsu`

The preset_canary_jiangsu supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_jiangxi`

The preset_canary_jiangxi supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_jilin`

The preset_canary_jilin supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_liaoning`

The preset_canary_liaoning supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_macau`

The preset_canary_macau supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_neimenggu`

The preset_canary_neimenggu supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_ningxia`

The preset_canary_ningxia supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_qinghai`

The preset_canary_qinghai supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_shaanxi`

The preset_canary_shaanxi supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_shandong`

The preset_canary_shandong supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_shanghai`

The preset_canary_shanghai supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_shanxi`

The preset_canary_shanxi supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_sichuan`

The preset_canary_sichuan supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_taiwan`

The preset_canary_taiwan supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_tianjin`

The preset_canary_tianjin supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_xinjiang`

The preset_canary_xinjiang supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_xizang`

The preset_canary_xizang supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_yunnan`

The preset_canary_yunnan supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_zhejiang`

The preset_canary_zhejiang supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

### `env_conf-preset_canary_overseas`

The preset_canary_overseas supports the following:

* `spec_name` - (Optional) The specification of the CPU time slice. Valid values: `5ms`, `50ms`, `100ms`.
* `code_rev` - (Optional) The version number of the code.
* `allowed_hosts` - (Optional) Allowed DCDN domain names.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Er.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Er.
* `update` - (Defaults to 5 mins) Used when update the Er.
* `delete` - (Defaults to 5 mins) Used when delete the Er.

## Import

DCDN Er can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_er.example <id>
```
