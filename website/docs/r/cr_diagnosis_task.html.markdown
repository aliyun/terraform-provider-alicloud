---
subcategory: "unreleased"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_diagnosis_task"
description: |-
  Provides a Alicloud CR Diagnosis Task resource.
---

# alicloud_cr_diagnosis_task

Provides a CR Diagnosis Task resource.

A task that diagnoses whether an instance has abnormal events.  .

For information about CR Diagnosis Task and how to use it, see [What is Diagnosis Task](https://next.api.alibabacloud.com/document/cr/2018-12-01/CreateDiagnosisTask).

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Basic"
  instance_name  = "tf-example"
}

resource "alicloud_cr_diagnosis_task" "default" {
  instance_id    = alicloud_cr_ee_instance.default.id
  diagnosis_type = "INSTANCE_DIAGNOSIS"
}
```

### Deleting `alicloud_cr_diagnosis_task` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cr_diagnosis_task`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `diagnosis_type` - (Required, ForceNew) Diagnosis type.  
* `instance_id` - (Required, ForceNew) Instance ID.  
* `targets` - (Optional, List) Diagnosis targets. See [`targets`](#targets) below.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.


### `targets`

The targets supports the following:
* `end_time` - (Optional) End time of the diagnosis period.
* `namespace` - (Optional) The namespace name.
* `related_id` - (Optional) The ID of the diagnosed item.
* `repository` - (Optional) Repository name.
* `start_time` - (Optional) Start time of the diagnosis period.
* `tag` - (Optional) Image tag.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<instance_id>:<diagnosis_task_id>:<diagnosis_type>`. `diagnosis_type` is carried in the id because the `GetDiagnosisTask` API requires it as a request parameter even though it is not part of the resource's primary identity.
* `diagnosis_task_id` - A resource property field that represents the primary resource ID.
* `diagnosis_time` - The time when the diagnosis task was triggered.
* `status` - A resource property field that represents the resource status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Diagnosis Task.
* `update` - (Defaults to 5 mins) Used when update the Diagnosis Task.

## Import

CR Diagnosis Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_diagnosis_task.example <instance_id>:<diagnosis_task_id>:<diagnosis_type>
```