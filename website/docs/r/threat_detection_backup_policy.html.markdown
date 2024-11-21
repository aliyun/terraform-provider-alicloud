---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_backup_policy"
sidebar_current: "docs-alicloud-resource-threat-detection-backup-policy"
description: |-
  Provides a Alicloud Threat Detection Backup Policy resource.
---

# alicloud\_threat\_detection\_backup\_policy

Provides a Threat Detection Backup Policy resource.

For information about Threat Detection Backup Policy and how to use it, see [What is Backup Policy](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createbackuppolicy).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_backup_policy&exampleId=30e1ca53-b512-dfe7-ad70-0d077ffb503953a6c3ac&activeTab=example&spm=docs.r.threat_detection_backup_policy.0.30e1ca53b5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_threat_detection_assets" "default" {
  machine_types = "ecs"
}
resource "alicloud_threat_detection_backup_policy" "default" {
  backup_policy_name = "tf-example-name"
  policy             = "{\"Exclude\":[\"/bin/\",\"/usr/bin/\",\"/sbin/\",\"/boot/\",\"/proc/\",\"/sys/\",\"/srv/\",\"/lib/\",\"/selinux/\",\"/usr/sbin/\",\"/run/\",\"/lib32/\",\"/lib64/\",\"/lost+found/\",\"/var/lib/kubelet/\",\"/var/lib/ntp/proc\",\"/var/lib/container\"],\"ExcludeSystemPath\":true,\"Include\":[],\"IsDefault\":1,\"Retention\":7,\"Schedule\":\"I|1668703620|PT24H\",\"Source\":[],\"SpeedLimiter\":\"\",\"UseVss\":true}"
  policy_version     = "2.0.0"
  uuid_list          = [data.alicloud_threat_detection_assets.default.ids.0]
}
```

## Argument Reference

The following arguments are supported:

* `backup_policy_name` - (Required) Protection of the Name of the Policy.
* `policy` - (Required) The Specified Protection Policies of the Specific Configuration. see [how to use it](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createbackuppolicy).
* `policy_version` - (Required, ForceNew) Anti-Blackmail Policy Version. Valid values: `1.0.0`, `2.0.0`.
* `uuid_list` - (Required) Specify the Protection of Server UUID List.
* `policy_region_id` - (Optional) The region ID of the non-Alibaba cloud server. You can call the [DescribeSupportRegion](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-describesupportregion) interface to view the region supported by anti-ransomware, and then select the region supported by anti-ransomware according to the region where your non-Alibaba cloud server is located.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Backup Policy.
* `status` - The status of the Backup Policy instance.

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Backup Policy.
* `update` - (Defaults to 5 mins) Used when update the Backup Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Backup Policy.

## Import

Threat Detection Backup Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_backup_policy.example <id>
```
