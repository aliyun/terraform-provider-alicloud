---
subcategory: "Brain Industrial"
layout: "alicloud"
page_title: "Alicloud: alicloud_brain_industrial_pid_organization"
sidebar_current: "docs-alicloud-resource-brain-industrial-pid-organization"
description: |-
  Provides a Alicloud Brain Industrial Pid Organization resource.
---

# alicloud\_brain\_industrial\_pid\_organization

Provides a Brain Industrial Pid Organization resource.

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_brain_industrial_pid_organization" "example" {
  pid_organization_name = "tf-testAcc"
}

```

## Argument Reference

The following arguments are supported:

* `parent_pid_organization_id` - (Optional, ForceNew) The ID of parent pid organization.
* `pid_organization_name` - (Required) The name of pid organization.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Pid Organization.

## Import

Brain Industrial Pid Organization can be imported using the id, e.g.

```
$ terraform import alicloud_brain_industrial_pid_organization.example <id>
```
