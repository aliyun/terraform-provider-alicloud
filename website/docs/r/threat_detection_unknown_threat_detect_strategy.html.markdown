---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_unknown_threat_detect_strategy"
sidebar_current: "docs-alicloud-resource-threat-detection-unknown-threat-detect-strategy"
description: |-
  Provides a Alicloud Threat Detection Unknown Threat Detect Strategy resource.
---

# alicloud_threat_detection_unknown_threat_detect_strategy

Provides a Threat Detection Unknown Threat Detect Strategy resource.

智能行为分析策略（Unknown Threat Detect Strategy）用于配置 Agentic-EDR 的未知威胁检测学习模型。

For information about Threat Detection Unknown Threat Detect Strategy and how to use it, see [CreateUnknownThreatDetectStrategy](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createunknownthreaddetectstrategy).

-> **NOTE:** The strategy requires Agentic-EDR to be purchased and activated; otherwise CreateUnknownThreatDetectStrategy returns `EdrSaleAdmissionDenied`.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_unknown_threat_detect_strategy&spm=docs.r.threat_detection_unknown_threat_detect_strategy" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_threat_detection_unknown_threat_detect_strategy" "default" {
  name                   = "tf-example-strategy"
  asset_selection_type   = "UNKNOWN_THREAT_DETECT_CONFIG_example"
  study_mode             = "hash"
  duration_days_after_init  = 1
  duration_days_after_stop  = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the strategy.
* `asset_selection_type` - (Required, ForceNew) The asset selection identifier. Changing this creates a new strategy.
* `study_mode` - (Optional, Computed) The whitelist mode. Valid values: `hash`, `path`.
* `duration_days_after_init` - (Optional, Computed) The number of days to study after the model is initially established before it ends.
* `duration_days_after_stop` - (Optional, Computed) The number of days after which the model ends when no new processes are added.

## Attributes Reference

The following attributes are exported:

* `id` - The strategy ID.
* `machine_count` - The number of servers covered by the strategy.

## Import

Threat Detection Unknown Threat Detect Strategy can be imported using the id, e.g.

```terraform
terraform import alicloud_threat_detection_unknown_threat_detect_strategy.example 12345678
```
