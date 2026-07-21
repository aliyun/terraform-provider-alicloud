---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_enterprise_snapshot_policy_attachment"
description: |-
  Provides a Alicloud EBS Enterprise Snapshot Policy Attachment resource.
---

# alicloud_ebs_enterprise_snapshot_policy_attachment

Provides a EBS Enterprise Snapshot Policy Attachment resource. Enterprise-level snapshot policy cloud disk binding relationship.

For information about EBS Enterprise Snapshot Policy Attachment and how to use it, see [What is Enterprise Snapshot Policy Attachment](https://next.api.aliyun.com/api/ebs/2021-07-30/BindEnterpriseSnapshotPolicy).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ebs_enterprise_snapshot_policy_attachment&exampleId=643a03f0-92b5-eeec-e55c-ee20a6897859d54529ba&activeTab=example&spm=docs.r.ebs_enterprise_snapshot_policy_attachment.0.643a03f092&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ecs_disk" "defaultJkW46o" {
  category          = "cloud_essd"
  description       = "esp-attachment-test"
  zone_id           = "cn-hangzhou-i"
  performance_level = "PL1"
  size              = "20"
  disk_name         = var.name
}

resource "alicloud_ebs_enterprise_snapshot_policy" "defaultPE3jjR" {
  status = "DISABLED"
  desc   = "DESC"
  schedule {
    cron_expression = "0 0 0 1 * ?"
  }
  enterprise_snapshot_policy_name = var.name

  target_type = "DISK"
  retain_rule {
    time_interval = "120"
    time_unit     = "DAYS"
  }
}


resource "alicloud_ebs_enterprise_snapshot_policy_attachment" "default" {
  policy_id = alicloud_ebs_enterprise_snapshot_policy.defaultPE3jjR.id
  disk_id   = alicloud_ecs_disk.defaultJkW46o.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ebs_enterprise_snapshot_policy_attachment&spm=docs.r.ebs_enterprise_snapshot_policy_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `disk_id` - (Optional, ForceNew, Computed) Cloud Disk ID.
* `policy_id` - (Required, ForceNew) the enterprise snapshot policy id.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<policy_id>:<disk_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Enterprise Snapshot Policy Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise Snapshot Policy Attachment.

## Import

EBS Enterprise Snapshot Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_enterprise_snapshot_policy_attachment.example <policy_id>:<disk_id>
```