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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_brain_industrial_pid_organization&exampleId=be7e5997-a69b-bc7c-e157-46e9c04b8e7443d8a23d&activeTab=example&spm=docs.r.brain_industrial_pid_organization.0.be7e5997a6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

```shell
$ terraform import alicloud_brain_industrial_pid_organization.example <id>
```
