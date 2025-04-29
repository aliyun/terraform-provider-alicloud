---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_client_key"
description: |-
  Provides a Alicloud KMS Client Key resource.
---

# alicloud_kms_client_key

Provides a KMS Client Key resource. Client key (of Application Access Point).

For information about KMS Client Key and how to use it, see [What is Client Key](https://www.alibabacloud.com/help/zh/key-management-service/latest/api-createclientkey).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kms_client_key&exampleId=531e3cfb-b06e-d765-cb63-25e7cf45557c262ac784&activeTab=example&spm=docs.r.kms_client_key.0.531e3cfbb0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_kms_application_access_point" "AAP0" {
  policies                      = ["aa"]
  description                   = "aa"
  application_access_point_name = var.name
}

resource "alicloud_kms_client_key" "default" {
  aap_name              = alicloud_kms_application_access_point.AAP0.application_access_point_name
  password              = "YouPassword123!"
  not_before            = "2023-09-01T14:11:22Z"
  not_after             = "2028-09-01T14:11:22Z"
  private_key_data_file = "./private_key_data_file.txt"
}
```

## Argument Reference

The following arguments are supported:
* `aap_name` - (Required, ForceNew) ClientKey's parent Application Access Point name.
* `not_after` - (Optional, ForceNew) The ClientKey expiration time. Example: "2027-08-10 T08:03:30Z".
* `not_before` - (Optional, ForceNew) The valid start time of the ClientKey. Example: "2022-08-10 T08:03:30Z".
* `password` - (Required, ForceNew) To enhance security, set a password for the downloaded Client Key,When an application accesses KMS, you must use the ClientKey content and this password to initialize the SDK client.
* `private_key_data_file` - (Optional, ForceNew) The name of file that can save access key id and access key secret. Strongly suggest you to specified it when you creating access key, otherwise, you wouldn't get its secret ever.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Create timestamp, e.g. "2022-08-10T08:03:30Z".

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Client Key.
* `delete` - (Defaults to 5 mins) Used when delete the Client Key.

## Import

KMS Client Key can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_client_key.example <id>
```

Resource attributes such as `password`, `private_key_data_file` are not available for imported resources as this information cannot be read from the KMS API.