---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_certificates"
description: |-
  Provides a list of Data Works Certificate owned by an Alibaba Cloud account.
---

# alicloud_data_works_certificates

This data source provides Data Works Certificate available to the user.[What is Certificate](https://next.api.alibabacloud.com/document/dataworks-public/2024-05-18/ImportCertificate)

-> **NOTE:** Available since v1.287.0.

## Example Usage

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

data "alicloud_data_works_certificates" "default" {
  ids        = ["${alicloud_data_works_certificate.default.id}"]
  project_id = alicloud_data_works_project.default.id
}

output "alicloud_data_works_certificate_example_id" {
  value = data.alicloud_data_works_certificates.default.certificates.0.id
}
```

## Argument Reference

The following arguments are supported:
* `project_id` - (Required) The ID of the DataWorks workspace to which the certificate file belongs.
* `ids` - (Optional, Computed) A list of Certificate IDs.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Certificate IDs.
* `certificates` - A list of Certificate Entries. Each element contains the following attributes:
  * `create_time` - Creation time in RFC3339 format (converted from the millisecond timestamp returned by the API).
  * `create_user` - The creator ID of the certificate file.
  * `description` - The description of the certificate file.
  * `file_size_in_bytes` - File size in bytes.
  * `id` - The unique identifier of the certificate file.
  * `name` - The name of the certificate file, unique within the workspace.
  * `project_id` - **NOTE:** This field is only available when `enable_details` is `true`. The ID of the DataWorks workspace to which the certificate file belongs.
