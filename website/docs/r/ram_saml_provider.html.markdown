---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_saml_provider"
description: |-
  Provides a Alicloud RAM Saml Provider resource.
---

# alicloud_ram_saml_provider

Provides a RAM Saml Provider resource.



For information about RAM Saml Provider and how to use it, see [What is Saml Provider](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ims-2019-08-15-createsamlprovider).

-> **NOTE:** Available since v1.114.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_saml_provider&exampleId=9af45430-b5a5-519a-941f-47f2d84c5c512e0b4b1a&activeTab=example&spm=docs.r.ram_saml_provider.0.9af45430b5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ram_saml_provider" "example" {
  saml_provider_name            = "terraform-example"
  encodedsaml_metadata_document = <<EOF
  PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciBlbnRpdHlJRD0iaHR0cDovL2V4YW1wbGUuYWxpeXVuLmNvbS9leGFtcGxlLWlkcCIgeG1sbnM6bWQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSI+PG1kOklEUFNTT0Rlc2NyaXB0b3IgV2FudEF1dGhuUmVxdWVzdHNTaWduZWQ9ImZhbHNlIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIj48bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPjxkczpLZXlJbmZvIHhtbG5zOmRzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj48ZHM6WDUwOURhdGE+PGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEL3pDQ0F1ZWdBd0lCQWdJRU1yb0tjakFOQmdrcWhraUc5dzBCQVFzRkFEQ0JnREVuTUNVR0ExVUVBeE1lClFXeHBlWFZ1SUZKQlRTQkZlR0Z0Y0d4bElFTmxjblJwWm1sallYUmxNUkF3RGdZRFZRUUxFd2RCYkdsaVlXSmgKTVJBd0RnWURWUVFLRXdkQmJHbGlZV0poTVJFd0R3WURWUVFIRXdoSVlXNW5XbWh2ZFRFUk1BOEdBMVVFQ0JNSQpXbWhsU21saGJtY3hDekFKQmdOVkJBWVRBa05PTUNBWERUSXpNVEl3TkRFeU1EY3dNRm9ZRHpJd05URXdOREl4Ck1USXdOekF3V2pDQmdERW5NQ1VHQTFVRUF4TWVRV3hwZVhWdUlGSkJUU0JGZUdGdGNHeGxJRU5sY25ScFptbGoKWVhSbE1SQXdEZ1lEVlFRTEV3ZEJiR2xpWVdKaE1SQXdEZ1lEVlFRS0V3ZEJiR2xpWVdKaE1SRXdEd1lEVlFRSApFd2hJWVc1bldtaHZkVEVSTUE4R0ExVUVDQk1JV21obFNtbGhibWN4Q3pBSkJnTlZCQVlUQWtOT01JSUJJakFOCkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQW9KVGVndWc0eXFaalNzQzFQUWpzbGxreUxWZHEKcXR0UGFqNmNOYldQdVFRNThSMkF4ZHNYeng4c05lOElLYUFYbW84azdDTFhDcXFLVzNjNEpzRWtTOUcva3B2NApJWFpBOGFVcDBCeXdQUFBocmFjUXd4cmJ5dkhja0dqVUpkNEZrOEVjbVVjNjRrSE5LbjBCaVJpL0NEZlM3MXBuCjh5T3dDNUZPSUlYWXhWMGtTTnNQMnozV2tBbFBXWm1sVkZSd1dxeHhGS2xCTjVpdVhaVHA4dk5rU2htVndBTW8KcjVpb2VBaFdXd0N1L0pvdUhLa3lnbVJnSDNhRjlSRlkrOGZ4NkMzR2hjZktISUszRTFBbVVtWlpjR3NDUCtxNApXeTBuSFp4QStaZEhTeE1OYUJPMm5JbkxJSHVDWHgza096eWpGV3dUaTVGSTlwdE5vNktBay9wRThRSURBUUFCCm8zMHdlekFQQmdOVkhSTUVDREFHQVFIL0FnRURNQjBHQTFVZEpRUVdNQlFHQ0NzR0FRVUZCd01CQmdnckJnRUYKQlFjREFqQUxCZ05WSFE4RUJBTUNBb1F3SFFZRFZSMFJBUUgvQkJNd0VZSUpiRzlqWVd4b2IzTjBod1IvQUFBQgpNQjBHQTFVZERnUVdCQlQ2TXluMnJjMVhEQTZqQkZYWVBOYitGaldMVmpBTkJna3Foa2lHOXcwQkFRc0ZBQU9DCkFRRUFoWHpUUzJJaHZjY3hzSVNzVVNFcldNNDJiQlZESHhTa05EemhPRmd0eGNtNUxuNHdjWXJvdkM3NHZxS1oKUWdQWmpGcWw3YUJTb1ZyNFdseWFaZlVBdHdNL3pZZytJbklUSVpBQ0dhM1VNK3h5V0NLSVhRNGpJVldnNG9QWQpxTStjNWllLzJFVlE0YmhObEQyL0lYZUVEZFd2TXMzdmFyRTFCUE5PQXJZZ2tZTmNER3lDSnA0ZmQ3d3ladWxhCllEZFFIWDdpdUJ1R0JOZFRBajlCUW5xaTJRcTc5RndMVTBrQkFvdUpVVVBPUjBpMGtwZ24vc2dSbHhvaHY1bHgKVTFwYVhtMEZRWHpUUDEvdjV5Y24wM3NVckFUekg2VkRpVlQ2N0NRQjR4MXJpOTFvUVRkWERXN1RvRkVhOGIrOApPdE8wZERMdDlnbCtNMkxYRzJTWnBZTkJoZz09PC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9tZDpLZXlEZXNjcmlwdG9yPjxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6ZW1haWxBZGRyZXNzPC9tZDpOYW1lSURGb3JtYXQ+PG1kOk5hbWVJREZvcm1hdD51cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6bmFtZWlkLWZvcm1hdDpwZXJzaXN0ZW50PC9tZDpOYW1lSURGb3JtYXQ+PG1kOlNpbmdsZVNpZ25PblNlcnZpY2UgQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUE9TVCIgTG9jYXRpb249Imh0dHA6Ly9leGFtcGxlLmFsaXl1bi5jb20vZXhhbXBsZS1pZHAvc3NvL3NhbWwiLz48bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHA6Ly9leGFtcGxlLmFsaXl1bi5jb20vZXhhbXBsZS1pZHAvc3NvL3NhbWwiLz48L21kOklEUFNTT0Rlc2NyaXB0b3I+PC9tZDpFbnRpdHlEZXNjcmlwdG9yPg==
  EOF
  description                   = "For Terraform Test"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The metadata file which is Base64-encoded.
The file is provided by an IdP that supports Security Assertion Markup Language (SAML) 2.0.
* `encodedsaml_metadata_document` - (Required) The new metadata file.
* `saml_provider_name` - (Required, ForceNew) The name of the IdP.
The name can be up to 128 characters in length. The name can contain letters, digits, `periods (.), hyphens (-), and underscores (_)`. The name cannot start or end with `periods (.), hyphens (-), or underscores (_)`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `arn` - The identity provider's ARN.
* `update_date` - Update time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Saml Provider.
* `delete` - (Defaults to 5 mins) Used when delete the Saml Provider.
* `update` - (Defaults to 5 mins) Used when update the Saml Provider.

## Import

RAM Saml Provider can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_saml_provider.example <id>
```