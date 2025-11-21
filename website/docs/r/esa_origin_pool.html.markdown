---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_origin_pool"
description: |-
  Provides a Alicloud ESA Origin Pool resource.
---

# alicloud_esa_origin_pool

Provides a ESA Origin Pool resource.



For information about ESA Origin Pool and how to use it, see [What is Origin Pool](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateOriginPool).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_origin_pool&exampleId=9e425e76-fb68-2057-33b6-502f60913fc8fcaeab69&activeTab=example&spm=docs.r.esa_origin_pool.0.9e425e76fb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_esa_site" "default" {
  site_name   = "gositecdn-${random_integer.default.result}.cn"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}


resource "alicloud_esa_origin_pool" "default" {
  origins {
    type    = "OSS"
    address = "example.oss-cn-beijing.aliyuncs.com"
    header  = "{\"Host\":[\"example.oss-cn-beijing.aliyuncs.com\"]}"
    enabled = "true"
    auth_conf {
      secret_key = "<SecretKeyId>"
      auth_type  = "private_cross_account"
      access_key = "<AccessKeyId>"
    }

    weight = "50"
    name   = "origin1"
  }
  origins {
    address = "example.s3.com"
    header  = "{\"Host\": [\"example1.com\"]}"
    enabled = "true"
    auth_conf {
      version    = "v2"
      region     = "us-east-1"
      auth_type  = "private"
      access_key = "<AccessKeyId>"
      secret_key = "<SecretKeyId>"
    }

    weight = "50"
    name   = "origin2"
    type   = "S3"
  }
  origins {
    type    = "S3"
    address = "example1111.s3.com"
    header  = "{\"Host\":[\"example1111.com\"]}"
    enabled = "true"
    auth_conf {
      secret_key = "<SecretKeyId>"
      version    = "v2"
      region     = "us-east-1"
      auth_type  = "private"
      access_key = "<AccessKeyId>"
    }

    weight = "30"
    name   = "origin3"
  }

  site_id          = alicloud_esa_site.default.id
  origin_pool_name = "exampleoriginpool"
  enabled          = "true"
}
```

## Argument Reference

The following arguments are supported:
* `enabled` - (Optional) Whether the source address pool is enabled:
  - `true`: Enabled;
  - `false`: Not enabled.
* `origin_pool_name` - (Required, ForceNew) The source address pool name.
* `origins` - (Optional, Set) The Source station information added to the source address pool. Multiple Source stations use arrays to transfer values. See [`origins`](#origins) below.
* `site_id` - (Required, ForceNew) The site ID.

### `origins`

The origins supports the following:
* `address` - (Optional) Origin Address.
* `auth_conf` - (Optional, List) The authentication information. When the source Station is an OSS or S3 and other source stations need to be authenticated, the authentication-related configuration information needs to be transmitted. See [`auth_conf`](#origins-auth_conf) below.
* `enabled` - (Optional) Whether the source station is enabled:
  - `true`: Enabled;
  - `false`: Not enabled.
* `header` - (Optional) The request header that is sent when returning to the source. Only Host is supported.
* `name` - (Optional) Origin Name.
* `type` - (Optional) Source station type:
ip_domain: ip or domain name type origin station;
  - `OSS`:OSS address source station;
  - `S3`:AWS S3 Source station.
* `weight` - (Optional, Int) Weight, 0-100.

### `origins-auth_conf`

The origins-auth_conf supports the following:
* `access_key` - (Optional) The AccessKey to be passed when AuthType is set to private_cross_account or private.
* `auth_type` - (Optional) Authentication type.
  - `public`: public read/write, which is used when the source station is OSS or S3 and is public read/write;
  - `private_same_account`: Used when the same account is private, the source station is OSS, and the authentication type is private authentication of the same account;
  - `private_cross_account`: private cross-account, used when the origin station is OSS and the authentication type is cross-account private authentication;
  - `private`: Used when the source station is S3 and the authentication type is private.
* `region` - (Optional) The Region of the source station to be transmitted when the source station is AWS S3.
* `secret_key` - (Optional) The SecretKey to be passed when AuthType is set to private_cross_account or private.
* `version` - (Optional) The signature version to be transmitted when the source station is AWS S3.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<origin_pool_id>`.
* `origin_pool_id` - OriginPool Id
* `origins` - The Source station information added to the source address pool. Multiple Source stations use arrays to transfer values.
  * `origin_id` - Origin ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Origin Pool.
* `delete` - (Defaults to 5 mins) Used when delete the Origin Pool.
* `update` - (Defaults to 5 mins) Used when update the Origin Pool.

## Import

ESA Origin Pool can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_origin_pool.example <site_id>:<origin_pool_id>
```