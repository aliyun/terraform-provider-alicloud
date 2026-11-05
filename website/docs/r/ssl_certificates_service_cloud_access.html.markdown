---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_cloud_access"
description: |-
  Provides a Certificate Management Service (Original SSL Certificate) Cloud Access resource.
---

# alicloud_ssl_certificates_service_cloud_access

Provides a Certificate Management Service (Original SSL Certificate) Cloud Access resource.

For information about Certificate Management Service (Original SSL Certificate) Cloud Access and how to use it, see [What is Cloud Access](https://next.api.alibabacloud.com/document/cas/2020-04-07/AddCloudAccess).

-> **NOTE:** Available since v2.2.0.

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

variable "name" {
  default = "terraform-example"
}

resource "alicloud_ssl_certificates_service_cloud_access" "default" {
  cloud_name = "Tencent"
  secret_id  = "your-tencent-cloud-secret-id"
  secret_key = "your-tencent-cloud-secret-key"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_name` - (Required, ForceNew) The cloud service provider. Currently, only Tencent Cloud (`Tencent`) is supported.
* `secret_id` - (Required, ForceNew) The SecretId used to access the external cloud.
* `secret_key` - (Required, ForceNew, Sensitive) The SecretKey used to access the external cloud.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. It is the AccessId of the cloud access authorization.
* `access_id` - (Computed, String) The AccessId of the cloud access authorization.
* `service_status` - (Computed, String) The status of the cloud access service.

## Import

Certificate Management Service (Original SSL Certificate) Cloud Access can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_cloud_access.example <access_id>
```
