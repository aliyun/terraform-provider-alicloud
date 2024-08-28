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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_brain_industrial_pid_loop&exampleId=6c285160-ca8c-80fa-2fdb-ccd431b48f3a876b32bb&activeTab=example&spm=docs.r.brain_industrial_pid_loop.0.6c285160ca&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
