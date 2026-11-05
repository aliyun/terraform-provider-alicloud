---
subcategory: "Voice Messaging Service (VMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_vms_skill"
description: |-
  Provides a Alicloud Voice Messaging Service (VMS) Skill resource.
---

# alicloud_vms_skill

Provides a Voice Messaging Service (VMS) Skill resource.

A VMS skill is a routing/skill group inside a Voice Messaging Service enterprise (call center instance), identified by its `enterprise_id`. Each skill has a server-assigned integer `Id` (primary key), a `Name`, an optional `Comment` and a `CreateTime`.

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The Voice Messaging Service exposes no update API for a skill. Every writable argument (`enterprise_id`, `name`, `comment`) forces a replacement: changing any of them destroys the existing skill and creates a new one. Do not concurrently manage the same skill with other tooling, because `Delete` removes the skill from the enterprise.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vms_skill" "default" {
  enterprise_id = "7000002"
  name          = "skill-example"
  comment       = "example skill managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `enterprise_id` - (Required, ForceNew, String) The ID of the Voice Messaging Service enterprise (call center instance) that owns the skill. Changing this recreates the resource.
* `name` - (Required, ForceNew, String) The name of the skill. Changing this recreates the resource.
* `comment` - (Optional, ForceNew, String) The description of the skill. Omitting it leaves the API default; changing the value recreates the resource.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The resource ID, in the format `<enterprise_id>:<skill_id>`. Use this format to import the resource.
* `create_time` - (String) The creation time of the skill, accurate to the second, returned by `CloudCreateSkill`.

## Import

VMS Skill can be imported using the composite ID `<enterprise_id>:<skill_id>`, e.g.

```shell
$ terraform import alicloud_vms_skill.example 7000002:49
```
