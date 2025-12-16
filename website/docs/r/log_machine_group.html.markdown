---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_machine_group"
sidebar_current: "docs-alicloud-resource-log-machine-group"
description: |-
  Provides a Alicloud log tail machine group resource.
---

# alicloud\_log\_machine\_group

Log Service manages all the ECS instances whose logs need to be collected by using the Logtail client in the form of machine groups.
 [Refer to details](https://www.alibabacloud.com/help/doc-detail/28966.htm)

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_machine_group&exampleId=7cb4e15a-0700-f4ee-2227-e3c857285ff16f30f0d9&activeTab=example&spm=docs.r.log_machine_group.0.7cb4e15a07&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  project_name = "terraform-example-${random_integer.default.result}"
  description  = "terraform-example"
}

resource "alicloud_log_machine_group" "example" {
  project       = alicloud_log_project.example.project_name
  name          = "terraform-example"
  identify_type = "ip"
  topic         = "terraform"
  identify_list = ["10.0.0.1", "10.0.0.2"]
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_log_machine_group&spm=docs.r.log_machine_group.example&intl_lang=EN_US)

## Module Support

You can use the existing [sls-logtail module](https://registry.terraform.io/modules/terraform-alicloud-modules/sls-logtail/alicloud) 
to create logtail config, machine group, install logtail on ECS instances and join instances into machine group one-click.

## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the machine group belongs.
* `name` - (Required, ForceNew) The machine group name, which is unique in the same project.
* `identify_type` - (Optional) The machine identification type, including IP and user-defined identity. Valid values are "ip" and "userdefined". Default to "ip".
* `identify_list`- (Required) The specific machine identification, which can be an IP address or user-defined identity.
* `topic` - (Optional) The topic of a machine group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log machine group. It formats of `<project>:<name>`.
* `project` - The project name.
* `name` - The machine group name.
* `identify_type` - The machine identification type.
* `identify_list` - The machine identification.
* `topic` - The machine group topic.

## Import

Log machine group can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_machine_group.example tf-log:tf-machine-group
```
