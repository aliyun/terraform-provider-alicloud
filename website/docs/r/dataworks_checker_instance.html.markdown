---
subcategory: "DataWorks"
layout: "alicloud"
page_title: "Alicloud: alicloud_dataworks_checker_instance"
description: |-
  Provides a DataWorks Checker Instance resource.
---

# alicloud_dataworks_checker_instance

Provides a DataWorks Checker Instance resource.

-> **NOTE:** Available since v1.294.0.

When a file submitted in the DataStudio UI enters the publish-check state,
DataWorks emits a deployment-check event to the customer (carrying a
`CheckerInstanceId`). The customer decides whether the file is allowed to
continue the publish check and calls back the
[CheckFileDeployment](https://next.api.alibabacloud.com/document/dataworks-public/2020-05-18/CheckFileDeployment)
operation with that `CheckerInstanceId` and the chosen `Status`. This
Terraform resource records that decision in Terraform state.

-> **NOTE:** The `CheckFileDeployment` operation is a callback-response API.
DataWorks does not expose a separate Create / Read / Delete / List API for the
CheckerInstance resource itself. As a result this resource is a
command-style resource:
* `Create` / `Update` invokes `CheckFileDeployment` to record the decision.
* `Read` is a no-op — the configuration is the source of truth because the
  product does not expose a way to query a persisted decision afterwards.
* `Delete` is a no-op — removing the resource from Terraform state does not
  retract the decision; DataWorks retains the last decision it received for
  a `CheckerInstanceId` until the next `CheckFileDeployment` call replaces
  it.

-> **NOTE:** `CheckFileDeployment` requires DataWorks Enterprise Edition or
Flagship Edition; the request fails with `Forbidden.Access` otherwise.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_dataworks_checker_instance" "default" {
  checker_instance_id = var.name
  status              = "OK"
  check_detail_url    = "https://example.invalid/check/${var.name}"
}
```

## Argument Reference

The following arguments are supported:
* `checker_instance_id` - (Required, ForceNew) ID of the checker instance. DataWorks assigns this ID when it emits the corresponding deployment-check event; reuse it unchanged when reporting the decision back.
* `status` - (Required) Customer decision on whether the file is allowed to continue the publish check. Valid values: `OK`, `FAIL`, `WARN`.
* `check_detail_url` - (Optional) URL pointing to the detail of the check result. The customer can provide this so DataWorks users can navigate to the check details.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is the same as `checker_instance_id`.
* `region_id` - Region ID of the DataWorks workspace that emitted the deployment-check event.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax/operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when invoke `CheckFileDeployment` to record the decision.
* `update` - (Defaults to 5 mins) Used when re-invoke `CheckFileDeployment` to update the decision.

## Import

DataWorks Checker Instance can be imported using the id, which is the
`CheckerInstanceId`, e.g.

```shell
$ terraform import alicloud_dataworks_checker_instance.example <checker_instance_id>
```

Note that `status` and `check_detail_url` cannot be queried from the server
(see the NOTE above); after import you must set them in the configuration to
match the decision you intend to maintain going forward.
