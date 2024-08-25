---
subcategory: "Brain Industrial"
layout: "alicloud"
page_title: "Alicloud: alicloud_brain_industrial_pid_project"
sidebar_current: "docs-alicloud-resource-brain-industrial-pid-project"
description: |-
  Provides a Alicloud Brain Industrial Pid Project resource.
---

# alicloud\_brain\_industrial\_pid\_project

Provides a Brain Industrial Pid Project resource.

-> **NOTE:** Available in v1.113.0+.

-> **DEPRECATED:**  This resource has been from version `1.222.0`.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_brain_industrial_pid_project&exampleId=5a1c9696-b2fa-42b1-ff37-2e1171b1c94c194f9a2c&activeTab=example&spm=docs.r.brain_industrial_pid_project.0.5a1c9696b2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_brain_industrial_pid_project" "example" {
  pid_organization_id = "3e74e684-cbb5-xxxx"
  pid_project_name    = "tf-testAcc"
}

```

## Argument Reference

The following arguments are supported:

* `pid_organization_id` - (Required) The ID of Pid Organization.
* `pid_project_desc` - (Optional) The description of Pid Project.
* `pid_project_name` - (Required) The name of Pid Project.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Pid Project.

## Import

Brain Industrial Pid Project can be imported using the id, e.g.

```shell
$ terraform import alicloud_brain_industrial_pid_project.example <id>
```
