---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone_record"
sidebar_current: "docs-alicloud-resource-pvtz-zone-record"
description: |-
  Provides a Alicloud Private Zone Record resource.
---

# alicloud\_pvtz\_zone\_record

Provides a Private Zone Record resource.

-> **NOTE:** Terraform will auto Create a Private Zone Record while it uses `alicloud_pvtz_zone_record` to build a Private Zone Record resource.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_pvtz_zone_record&exampleId=b6ead031-6af1-72be-b9ee-f679cee4755801bc96f2&activeTab=example&spm=docs.r.pvtz_zone_record.0.b6ead0316a" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
resource "alicloud_pvtz_zone" "zone" {
  name = "foo.test.com"
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
* `lang` - (Optional, Available in 1.109.0+) User language.
* `resource_record` - (Optional, ForceNew, Deprecated from v1.109.0+) The resource record of the Private Zone Record.
* `rr` - (Optional, ForceNew) The rr of the Private Zone Record.
* `type` - (Required) The type of the Private Zone Record. Valid values: A, CNAME, TXT, MX, PTR, SRV.
* `value` - (Required) The value of the Private Zone Record.
* `ttl` - (Optional) The ttl of the Private Zone Record. Default to `60`.
* `priority` - (Optional) The priority of the Private Zone Record. At present, only can "MX" record support it. Valid values: [1-99]. Default to 1.
* `remark` - (Optional, Available in 1.103.2+) The remark of the Private Zone Record.
* `status` - (Optional, Available in 1.109.0+) Resolve record status. Value:
    - ENABLE: enable resolution.
    - DISABLE: pause parsing.

### Timeouts

-> **NOTE:** Available in 1.109.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

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

