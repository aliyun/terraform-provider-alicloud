---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_domain"
sidebar_current: "docs-alicloud-resource-direct-mail-domain"
description: |-
  Provides a Alicloud Direct Mail Domain resource.
---

# alicloud_direct_mail_domain

Provides a Direct Mail Domain resource.

For information about Direct Mail Domain and how to use it, see [What is Domain](https://www.alibabacloud.com/help/en/doc-detail/29414.htm).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_direct_mail_domain&exampleId=8cf26e8e-e78e-4738-8871-c25bdf29a899db03147b&activeTab=example&spm=docs.r.direct_mail_domain.0.8cf26e8ee7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}
provider "alicloud" {
  region = "cn-hangzhou"
}
resource "alicloud_direct_mail_domain" "example" {
  domain_name = "alicloud-provider-${random_integer.default.result}.online"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_direct_mail_domain&spm=docs.r.direct_mail_domain.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Domain, length `1` to `50`, including numbers or capitals or lowercase letters or `.` or `-`

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Domain.
* `status` - The status of the domain name. Valid values:`0` to `4`. `0`:Available, Passed. `1`: Unavailable, No passed. `2`: Available, cname no passed, icp no passed. `3`: Available, icp no passed. `4`: Available, cname no passed.

## Import

Direct Mail Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_direct_mail_domain.example <id>
```
