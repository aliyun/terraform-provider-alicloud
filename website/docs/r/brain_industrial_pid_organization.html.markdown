---
subcategory: "Brain Industrial"
layout: "alicloud"
page_title: "Alicloud: alicloud_brain_industrial_pid_organization"
sidebar_current: "docs-alicloud-resource-brain-industrial-pid-organization"
description: |-
  Provides a Alicloud Brain Industrial Pid Organization resource.
---

# alicloud_brain_industrial_pid_organization

Provides a Brain Industrial Pid Organization resource.

-> **NOTE:** Available since v1.113.0.

-> **DEPRECATED:**  This resource has been deprecated from version `1.222.0`.

## Example Usage

Basic Usage

```terraform
resource "alicloud_brain_industrial_pid_organization" "example" {
  pid_organization_name = "tf-testAcc"
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_brain_industrial_pid_organization&spm=docs.r.brain_industrial_pid_organization.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `parent_pid_organization_id` - (Optional, ForceNew) The ID of parent pid organization.
* `pid_organization_name` - (Required) The name of pid organization.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Pid Organization.

## Import

Brain Industrial Pid Organization can be imported using the id, e.g.

```shell
$ terraform import alicloud_brain_industrial_pid_organization.example <id>
```
