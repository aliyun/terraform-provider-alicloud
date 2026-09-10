---
subcategory: "Intelligent Media Management (IMM)"
layout: "alicloud"
page_title: "Alicloud: alicloud_imm_project"
sidebar_current: "docs-alicloud-resource-imm-project"
description: |-
  Provides a Alicloud Intelligent Media Management Project resource.
---

# alicloud_imm_project

Provides a Intelligent Media Management Project resource.

For information about Intelligent Media Management Project and how to use it, see [What is Project](https://www.alibabacloud.com/help/en/network-intelligence-service/latest/user-overview).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_imm_project&exampleId=fd26d7e6-95c4-eade-3aae-38af4a8744083af28020&activeTab=example&spm=docs.r.imm_project.0.fd26d7e695&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}
resource "alicloud_ram_role" "role" {
  name        = var.name
  document    = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "imm.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
  description = "this is a role test."
  force       = true
}
resource "alicloud_imm_project" "example" {
  project      = var.name
  service_role = alicloud_ram_role.role.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_imm_project&spm=docs.r.imm_project.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The name of Project.
* `service_role` - (Optional) The service role authorized to the Intelligent Media Management service to access other cloud resources. Default value: `AliyunIMMDefaultRole`. You can also create authorization  roles through the `alicloud_ram_role`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Project. Its value is same as `project`.

## Import

Intelligent Media Management Project can be imported using the id, e.g.

```shell
$ terraform import alicloud_imm_project.example <project>
```
