---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_acl_entry_attachment"
sidebar_current: "docs-alicloud-resource-slb-acl-entry-attachment"
description: |-
  Provides a Acl entry attachment resource.
---

# alicloud\_slb\_acl\_entry\_attachment

-> **NOTE:** Available in v1.162.0+.

-> **NOTE:** The maximum number of entries per acl is 300.

For information about acl entry attachment and how to use it, see [Configure an acl entry](https://www.alibabacloud.com/help/en/doc-detail/70023.html).


## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_acl_entry_attachment&exampleId=7c6c4728-07cc-937d-e55e-9ac7e8f9feee94398cfc&activeTab=example&spm=docs.r.slb_acl_entry_attachment.0.7c6c472807&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_slb_acl" "attachment" {
  name       = "forSlbAclEntryAttachment"
  ip_version = "ipv4"
}

resource "alicloud_slb_acl_entry_attachment" "attachment" {
  acl_id  = alicloud_slb_acl.attachment.id
  entry   = "168.10.10.0/24"
  comment = "second"
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) The ID of the Acl.
* `entry` - (Required, ForceNew) The CIDR blocks.
* `comment` - (Optional, ForceNew) The comment of the entry.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value formats as `<acl_id>:<entry>`.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the resource.
* `delete` - (Defaults to 5 mins) Used when delete the resource.

## Import

Acl entry attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_acl_entry_attachment.example <acl_id>:<entry>
```
