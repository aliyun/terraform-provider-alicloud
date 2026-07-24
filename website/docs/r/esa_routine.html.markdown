---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_routine"
description: |-
  Provides a Alicloud ESA Routine resource.
---

# alicloud_esa_routine

Provides a ESA Routine resource.



For information about ESA Routine and how to use it, see [What is Routine](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateRoutine).

-> **NOTE:** Available since v1.251.0.

## Example Usage

Basic Usage
<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_routine&exampleId=dccb9f7a-016f-67d9-b149-8f737d5d73e49f363584&activeTab=example&spm=docs.r.esa_routine.0.dccb9f7a01&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_esa_routine" "default" {
  description      = var.name
  name             = var.name
  code             = "addEventListener('fetch', e => e.respondWith(new Response('hello world')))"
  code_description = "initial version"
  deploy_env       = "staging"
}
```

Manage the routine code from a local file:

```terraform
resource "alicloud_esa_routine" "from_file" {
  name = "terraform-example-file"
  code = file("${path.module}/index.js")
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_routine&spm=docs.r.esa_routine.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `name` - (Required, ForceNew) Routine Name, which must be unique in the same account.
* `description` - (Optional, ForceNew) The description of the routine.
* `code` - (Optional) The JavaScript source code of the routine. When set or changed, the code is uploaded as a new staging version and then committed into a formal code version. To manage the code from a local file, use the Terraform built-in `file()` function, e.g. `code = file("index.js")`.
* `code_description` - (Optional) The description attached to the committed code version.
* `deploy_env` - (Optional) The environment whose environment variables are bundled when committing the code version. Valid values: `staging`, `production`. If not set, no environment variables are bundled.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the routine was created.
* `latest_code_version` - The most recent committed code version of the routine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Routine.
* `update` - (Defaults to 5 mins) Used when update the Routine.
* `delete` - (Defaults to 5 mins) Used when delete the Routine.

## Import

ESA Routine can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_routine.example <id>
```