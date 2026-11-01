---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_certificate"
description: |-
  Provides an Alicloud Data Works Certificate resource.
---

# alicloud_data_works_certificate

Provides a Data Works Certificate resource. The certificate file is imported into a DataWorks workspace so that data sources such as RDS can authenticate with SSL.

For information about Data Works Certificate and how to use it, see [What is Certificate](https://next.api.alibabacloud.com/document/dataworks-public/2024-05-18/ImportCertificate).

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** `certificate_file` is the URL of the certificate file, not its content. The object must be readable by DataWorks; typically it is an object in an OSS bucket owned by the current account.

-> **NOTE:** DataWorks has no update API for certificate files. Changing `certificate_file`, `name`, `description` or `project_id` forces a new resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "random_integer" "randint" {
  max = 999
  min = 1
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "default" {
  status                  = "Available"
  description             = "tf_desc"
  project_name            = "${var.name}${random_integer.randint.id}"
  pai_task_enabled        = "false"
  display_name            = "tf_new_api_display"
  dev_role_disabled       = "true"
  dev_environment_enabled = "false"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}${random_integer.randint.id}"
  acl    = "public-read"
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "example.crt"
  acl     = "public-read"
  content = <<-EOT
-----BEGIN CERTIFICATE-----
MIICqDCCAZACCQCbNs38rCeB6TANBgkqhkiG9w0BAQsFADAWMRQwEgYDVQQDDAt0
Zi1hY2MtdGVzdDAeFw0yNjA5MDMwMDE4MDFaFw0zNjA4MzEwMDE4MDFaMBYxFDAS
BgNVBAMMC3RmLWFjYy10ZXN0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAwyKgyXSNG+NfmE6Ukz7oX4AgDbTaHAyJ3BWHfTK/DJNY0+0tgYE5Zle8bNui
rksdOpydcGAJeaoFNpCGwvd5FO9HklcYwLGdzASRXEHEN0XtG2oeHNnyH82yE2SH
hINbhYcM98WCEuIU5cIsohaQUXhdgBzQcUdc7q5mG/U87LQk4ZGq+JuqUhOAB2xG
8phoguTZ3ODk4+x9MEycpXfpTVkzz9y4JmL6jh3dEZvpyMmC94NYjrGRbioI9CDe
mpY1BsHRtSgxJ7Pvl5yYy7Tj2KC3jRv4FEP7C8vg9tqqV37yLtZsXw/UeVhn9nOm
tOQcqaS3o4BTp//Db7/ReoE2fwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBn+Hi1
Lq6pHiV0PPLN+4cQK8sLFSk1nHvJ1ONk9FsBmsnkrbky6I5fNQiBnzOm/rNl72Ao
JvsyROYVyyxX/M5LnZQwahUwDTaz4cNasUiPt33lNtueiLLhzC7UpB8NzH/v29MT
YnEmgc7msKrSnj8SfvZu6y9vbmLhJ0CbrpBTmcG6cXr/Em+l781p4hFdJmpOA7o8
oX+vfX/riPSP8bGgYlmsr2Iv4Uo/b8q6ZpF5p76M/MJirligYjc5fZG8UmsBu6eX
d5LC+4/r6W5uyIS/7fF7FLdzHfRhZknuu8Wzx+HXI+oKv0r/z8MGa1DsvJ6DHQIE
qWre020mtWvGRU1t
-----END CERTIFICATE-----
EOT
}

resource "alicloud_data_works_certificate" "default" {
  project_id       = alicloud_data_works_project.default.id
  name             = "${var.name}${random_integer.randint.id}"
  description      = "tf certificate"
  certificate_file = "https://${alicloud_oss_bucket.default.id}.${alicloud_oss_bucket.default.extranet_endpoint}/${alicloud_oss_bucket_object.default.key}"
}
```

## Argument Reference

The following arguments are supported:
* `project_id` - (Required, ForceNew, Int) The ID of the DataWorks workspace to which the certificate file belongs.
* `name` - (Required, ForceNew, String) The name of the certificate file. It must be unique within the workspace.
* `certificate_file` - (Required, ForceNew, String, Sensitive) The URL of the certificate file. The object must be readable by DataWorks. Changing this value replaces the certificate.
* `description` - (Optional, ForceNew, String) The description of the certificate file.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<project_id>:<certificate_id>`.
* `create_time` - The time when the certificate file was created, in RFC3339 format.
* `create_user` - The ID of the user who created the certificate file.
* `file_size_in_bytes` - The size of the certificate file, in bytes.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Certificate.

## Import

Data Works Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_certificate.example <project_id>:<certificate_id>
```

-> **NOTE:** `certificate_file` is write-only: the API never returns it, so it cannot be recovered on import. Set the same value in configuration before applying an imported certificate, otherwise the next apply replaces the certificate.
