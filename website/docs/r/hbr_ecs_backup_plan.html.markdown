---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ecs_backup_plan"
sidebar_current: "docs-alicloud-resource-hbr-ecs-backup-plan"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Ecs Backup Plan resource.
---

# alicloud\_hbr\_ecs\_backup\_plan

Provides a HBR Ecs Backup Plan resource.

For information about HBR Ecs Backup Plan and how to use it, see [What is Ecs Backup Plan](https://www.alibabacloud.com/help/doc-detail/186574.htm).

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "valut-name"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

data "alicloud_instances" "default" {
  name_regex = "no-deleteing-hbr-ecs-backup-plan"
  status     = "Running"
}

resource "alicloud_hbr_ecs_backup_plan" "example" {
  ecs_backup_plan_name = "example_value"
  instance_id          = data.alicloud_instances.default.instances.0.id
  vault_id             = alicloud_hbr_vault.default.id
  retention            = "1"
  schedule             = "I|1602673264|PT2H"
  backup_type          = "COMPLETE"
  speed_limit          = "0:24:5120"
  path                 = ["/home", "/var"]
  exclude              = <<EOF
  ["/home/exclude"]
  EOF
  include              = <<EOF
  ["/home/include"]
  EOF
}
```

## Argument Reference

The following arguments are support:

* `ecs_backup_plan_name` - (Required) The name of the backup plan. 1~64 characters, the backup plan name of each data source type in a single warehouse required to be unique.
* `vault_id` - (Required, ForceNew) The ID of Backup vault.
* `instance_id` - (Required, ForceNew) The ID of ECS instance. The ecs backup client must have been installed on the host.
* `retention` - (Required) Backup retention days, the minimum is 1.
* `backup_type` - (Required, ForceNew) Backup type. Valid values: `COMPLETE`.
* `schedule` - (Required) Backup strategy. Optional format: `I|{startTime}|{interval}`. It means to execute a backup task every `{interval}` starting from `{startTime}`. The backup task for the elapsed time will not be compensated. If the last backup task has not completed yet, the next backup task will not be triggered.
    * `startTime` Backup start time, UNIX time seconds.
    * `interval` ISO8601 time interval. E.g: `PT1H` means one hour apart. `P1D` means one day apart. 
* `disabled` - (Optional) Whether to disable the backup task. Valid values: `true`, `false`.
* `options` - (Optional) Windows operating system with application consistency using VSS, e.g: `{\"UseVSS\":false}`.
* `speed_limit` - (Optional) Flow control. The format is: `{start}|{end}|{bandwidth}`. Use `|` to separate multiple flow control configurations, multiple flow control configurations not allowed to have overlapping times.
    * `start` starting hour 
    * `end` end hour 
    * `bandwidth` limit rate, in KiB
* `path` - (Optional) List of backup path. e.g. `["/home", "/var"]`. **Note** If `path` is empty, it means that all directories will be backed up.
* `exclude` - (Optional) Exclude path. String of Json list, up to 255 characters. e.g. `"[\"/home/work\"]"`
* `include` - (Optional) Include path. String of Json list, up to 255 characters. e.g. `"[\"/var\"]"`
* `update_paths` - (Optional, Deprecated in v1.139.0+) Attribute update_paths has been deprecated in v1.139.0+, and you do not need to set it anymore.
* `detail` - (Optional) The detail of the backup plan.

## Notice

**About Backup path rules:**
1. If there is no wildcard `*`, you can enter 8 items of path.
2. When using wildcard `*`, only one item of path can be input, and wildcards like `/*/*` are supported.
3. Each item of path only supports absolute paths, for example starting with `/`, `\`, `C:\`, `D:\`.

**About Restrictions:**
1. When using `VSS`: multiple paths, UNC paths, wildcards, and excluded files not supported.
2. When using `UNC`: VSS not supported, wildcards not supported, and files to be excluded are not supported.

**About include/exclude path rules:**
1. Supports up to 8 paths, including paths using wildcards `*`.
2. If the path does not contain `/`, then `*` matches multiple path names or file names, for example `*abc*` will match `/abc/`, `/d/eabcd/`, `/a/abc`; `*.txt` will match all files with an extension `.txt`.
3. If the path contains `/`, each `*` only matches a single-level path or file name. For example, `/a/*/*/` share will match `/a/b/c/share`, but not `/a/d/share`.
4. If the path ends with `/`, it means the folder matches. For example, `*tmp/` will match `/a/b/aaatmp/`, `/tmp/` and so on.
5. The path separator takes Linux system `/` as an example, if it is Windows system, please replace it with `\`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ecs Backup Plan.

## Import

HBR Ecs Backup Plan can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_ecs_backup_plan.example <id>
```
