---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_api_sec_token"
description: |-
  Provides a ESA Api Sec Token resource.
---

# alicloud_esa_api_sec_token

Provides a ESA Api Sec Token resource.

-> **NOTE:** The ApiSecToken APIs (`AddApiSecToken`, `GetApiSecToken`, `SetApiSecToken`, `DeleteApiSecToken`) are **Private**. Confirm availability and the accepted `token_type` enum members, `token_sources` and `credentials` payload shapes with the ESA product owner before use.

## Example Usage

```hcl
variable "name" {
  default = "tf-example-esa-api-sec-token"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "${lower(var.name)}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_api_sec_token" "default" {
  site_id       = alicloud_esa_site.default.id
  token_name    = var.name
  token_type    = "jwt"
  token_sources = jsonencode([{ source_type = "api" }])
  credentials   = jsonencode([{ kty = "oct", k = "dGVzdC1rZXk" }])
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, ForceNew) The ID of the ESA site. Changing this creates a new resource.
* `token_name` - (Required) The name of the API security token.
* `token_type` - (Required) The type of the API security token. Valid values are to be confirmed with the ESA product owner.
* `token_sources` - (Required) The sources of the API security token, as a JSON-encoded array. The inner structure is opaque to the provider; provide the JSON array expected by the ESA API.
* `credentials` - (Required, ForceNew) The credentials of the API security token, as a JSON-encoded array of JWKs. Create-only: `SetApiSecToken` does not accept credentials, so changing this forces a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format `<site_id>:<token_id>`.
* `token_id` - The ID of the API security token assigned by the API.
* `credentials_count` - The number of credentials.
* `api_rule_ids` - The API rule IDs, as a JSON-encoded array.
* `last_update_at` - The last update time of the token.
* `endpoint_ids` - The endpoint IDs, as a JSON-encoded array.
* `identifier_ids` - The identifier IDs, as a JSON-encoded array.

## Import

ESA Api Sec Token can be imported using the id, e.g.

```shell
terraform import alicloud_esa_api_sec_token.example <site_id>:<token_id>
```
