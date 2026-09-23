---
subcategory: "Brain Industrial"
layout: "alicloud"
page_title: "Alicloud: alicloud_brain_industrial_pid_loop"
sidebar_current: "docs-alicloud-resource-brain-industrial-pid-loop"
description: |-
  Provides a Alicloud Brain Industrial Pid Loop resource.
---

# alicloud_brain_industrial_pid_loop

Provides a Brain Industrial Pid Loop resource.

-> **NOTE:** Available since v1.117.0.

-> **DEPRECATED:**  This resource has been deprecated from version `1.229.1`.

## Example Usage

Basic Usage

```terraform
resource "alicloud_brain_industrial_pid_loop" "example" {
  pid_loop_configuration = "YourLoopConfiguration"
  pid_loop_dcs_type      = "standard"
  pid_loop_is_crucial    = true
  pid_loop_name          = "tf-testAcc"
  pid_loop_type          = "0"
  pid_project_id         = "856c6b8f-ca63-40a4-xxxx-xxxx"
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_brain_industrial_pid_loop&spm=docs.r.brain_industrial_pid_loop.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `pid_loop_configuration` - (Required) The Pid Loop Configuration.
* `pid_loop_dcs_type` - (Required, ForceNew) The dcs type of Pid Loop. Valid values: `standard`.
* `pid_loop_desc` - (Optional) The desc of Pid Loop.
* `pid_loop_is_crucial` - (Required) Whether is crucial Pid Loop.
* `pid_loop_name` - (Required) The name of Pid Loop.
* `pid_loop_type` - (Required) The type of Pid Loop. Valid values: `0`, `1`, `2`, `3`, `4`, `5`.
* `pid_project_id` - (Required) The pid project id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Pid Loop.
* `status` - The status of Pid Loop.

## Import

Brain Industrial Pid Loop can be imported using the id, e.g.

```shell
$ terraform import alicloud_brain_industrial_pid_loop.example <id>
```
