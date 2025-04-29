---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone_record"
sidebar_current: "docs-alicloud-resource-pvtz-zone-record"
description: |-
  Provides a Alicloud Private Zone Record resource.
---

# alicloud_pvtz_zone_record

Provides a Private Zone Record resource.

-> **NOTE:** Terraform will auto Create a Private Zone Record while it uses `alicloud_pvtz_zone_record` to build a Private Zone Record resource.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pvtz_zone_record&exampleId=fc2efd0a-3bf4-6e20-1008-129e2b25871d0e6af25f&activeTab=example&spm=docs.r.pvtz_zone_record.0.fc2efd0a3b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_pvtz_zone" "zone" {
  zone_name = "foo.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
  zone_id = alicloud_pvtz_zone.zone.id
  rr      = "www"
  type    = "CNAME"
  value   = "bbb.test.com"
  ttl     = 60
}
```
## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) The name of the Private Zone Record.
* `lang` - (Optional, Available since v1.109.0) User language.
* `resource_record` - (Optional, ForceNew, Deprecated since v1.109.0) The resource record of the Private Zone Record.
* `rr` - (Optional, ForceNew) The rr of the Private Zone Record.
* `type` - (Required) The type of the Private Zone Record. Valid values: A, CNAME, TXT, MX, PTR, SRV.
* `value` - (Required) The value of the Private Zone Record.
* `ttl` - (Optional) The ttl of the Private Zone Record. Default to `60`.
* `priority` - (Optional) The priority of the Private Zone Record. At present, only can "MX" record support it. Valid values: [1-99]. Default to 1.
* `remark` - (Optional, Available since v1.103.2) The remark of the Private Zone Record.
* `status` - (Optional, Available since v1.109.0) Resolve record status. Value:
    - ENABLE: enable resolution.
    - DISABLE: pause parsing.

## Timeouts

-> **NOTE:** Available since v1.109.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Private Zone Record.
* `update` - (Defaults to 1 mins) Used when updating the Private Zone Record.
* `delete` - (Defaults to 1 mins) Used when terminating the Private Zone Record. 

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is formate as `<record_id>:<zone_id>`.
* `record_id` - The Private Zone Record ID.

## Import

Private Zone Record can be imported using the id, e.g.

```shell
$ terraform import alicloud_pvtz_zone_record.example abc123456
```

