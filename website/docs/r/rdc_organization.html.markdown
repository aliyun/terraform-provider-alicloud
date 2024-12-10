---
subcategory: "Apsara Devops(RDC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_rdc_organization"
sidebar_current: "docs-alicloud-resource-rdc-organization"
description: |-
  Provides a Alicloud RDC Organization resource.
---

# alicloud_rdc_organization

Provides a RDC Organization resource.

For information about RDC Organization and how to use it, see [What is Organization](https://help.aliyun.com/document_detail/51678.html).

-> **NOTE:** Available since v1.137.0.

-> **DEPRECATED:** This resource has been deprecated from version `1.238.0`.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rdc_organization&exampleId=a3011801-d6f5-97cf-5394-61233647189078396e0d&activeTab=example&spm=docs.r.rdc_organization.0.a3011801d6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_rdc_organization" "example" {
  organization_name = "example_value"
  source            = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `desired_member_count` - (Optional) The desired member count.
* `organization_name` - (Required, ForceNew, ForceNew) Company name.
* `real_pk` - (Optional) User pk, not required, only required when the ak used by the calling interface is inconsistent with the user pk
* `source` - (Required) This is organization source information

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Organization.

## Import

RDC Organization can be imported using the id, e.g.

```shell
$ terraform import alicloud_rdc_organization.example <id>
```
