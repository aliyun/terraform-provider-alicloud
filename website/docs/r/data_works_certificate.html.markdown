---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_certificate"
description: |-
  Provides a Alicloud Data Works Certificate resource.
---

# alicloud_data_works_certificate

Provides a Data Works Certificate resource.



For information about Data Works Certificate and how to use it, see [What is Certificate](https://next.api.alibabacloud.com/document/dataworks-public/2024-05-18/ImportCertificate).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shenzhen"
}

resource "alicloud_data_works_project" "default" {
  description      = "Terraform example project"
  project_name     = "terraform-example"
  pai_task_enabled = false
  display_name     = "terraform-example"
}


resource "alicloud_data_works_certificate" "default" {
  project_id       = alicloud_data_works_project.default.id
  name             = "terraform-example"
  certificate_file = "https://example.com/certificate.pem"
}
```

## Argument Reference

The following arguments are supported:
* `certificate_file` - (Required, ForceNew, Sensitive) The certificate file content. Changes require re-import.
* `description` - (Optional, ForceNew) The description of the certificate file.
* `name` - (Required, ForceNew) The name of the certificate file, unique within the workspace.
* `project_id` - (Required, ForceNew, Int) The ID of the DataWorks workspace to which the certificate file belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time in RFC3339 format (converted from the millisecond timestamp returned by the API).
* `create_user` - The creator ID of the certificate file.
* `file_size_in_bytes` - File size in bytes.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Certificate.

## Import

Data Works Certificate can be imported using the composite id `<project_id>:<certificate_id>`, e.g.

```shell
$ terraform import alicloud_data_works_certificate.example <project_id>:<certificate_id>
```