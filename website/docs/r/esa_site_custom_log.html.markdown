---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_site_custom_log"
description: |-
  Provides a Alicloud ESA Site Custom Log resource.
---

# alicloud_esa_site_custom_log

Provides a ESA Site Custom Log resource.

Site Custom Log is a site-level singleton: each ESA site owns at most one custom log configuration, keyed by `site_id`. It selects which request headers, response headers and cookies are recorded as custom log fields for the site.

For information about ESA Site Custom Log and how to use it, see [What is Site Custom Log](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateSiteCustomLog).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** ESA does not expose a `DeleteSiteCustomLog` API. Destroying `alicloud_esa_site_custom_log` is a no-op: Terraform removes the resource from state, while the server-side custom log configuration remains. Updating the resource replaces the configured fields in place.

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

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site_custom_log" "default" {
  site_id = data.alicloud_esa_sites.default.sites.0.id

  request_headers = [
    "x-acs-request-id",
    "x-real-ip",
  ]
  response_headers = [
    "X_Swift_Error",
    "Server",
  ]
  cookies = [
    "sessionId",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, ForceNew, String) The ID of the ESA site. The site ID can be obtained by calling [ListSites](https://next.api.alibabacloud.com/document/ESA/2024-09-10/ListSites). The custom log configuration is a singleton scoped to this site.
* `request_headers` - (Optional, List, String) The request header fields to record. Each element is a header name, e.g. `x-acs-request-id`.
* `response_headers` - (Optional, List, String) The response header fields to record. Each element is a header name, e.g. `X_Swift_Error`.
* `cookies` - (Optional, List, String) The cookie fields to record. Each element is a cookie name, e.g. `sessionId`.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed, String) The resource ID in terraform of Site Custom Log. The value is the `site_id`.
* `is_exist` - (Computed, Bool) Whether the custom log configuration exists.
* `config_id` - (Computed, String) The configuration ID of the custom log fields.

## Import

ESA Site Custom Log can be imported using the id, e.g. the `site_id`, e.g.

```shell
$ terraform import alicloud_esa_site_custom_log.example <site_id>
```
