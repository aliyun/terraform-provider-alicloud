---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_routes"
sidebar_current: "docs-alicloud-datasource-apig-routes"
description: |-
  Provides a list of Apig Route owned by an Alibaba Cloud account.
---

# alicloud_apig_routes

This data source provides Apig Route available to the user.[What is Route](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateHttpApiRoute)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_apig_http_api" "route_httpapi_arr" {
  http_api_name = "cspec-route-arr-httpapi"
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "array test httpapi"
  base_path     = "/cspec-route-arr"
}


resource "alicloud_apig_route" "default" {
  backend {
    services {
      version    = "v1.0"
      port       = "8080"
      protocol   = "HTTP"
      weight     = "70"
      service_id = "svc-a"
    }
    services {
      version    = "v1.0"
      port       = "8081"
      protocol   = "HTTP"
      weight     = "30"
      service_id = "svc-b"
    }
    scene = "SingleService"
  }
  description = "array route description"
  route_name  = "cspec-arr-route"
  domain_infos {
  }
  domain_infos {
  }
  http_api_id = alicloud_apig_http_api.route_httpapi_arr.id
  match {
    path {
      type  = "Prefix"
      value = "/arr-path"
    }
    headers {
      type  = "Exact"
      value = "v1"
      name  = "h1"
    }
    headers {
      type  = "Prefix"
      value = "v2"
      name  = "h2"
    }
    query_params {
      type  = "Exact"
      value = "v1"
      name  = "q1"
    }
    query_params {
      type  = "Prefix"
      value = "v2"
      name  = "q2"
    }
    methods         = ["GET", "POST", "PUT"]
    ignore_uri_case = false
  }
}

data "alicloud_apig_routes" "default" {
  ids         = ["${alicloud_apig_route.default.id}"]
  name_regex  = alicloud_apig_route.default.route_name
  http_api_id = alicloud_apig_http_api.route_httpapi_arr.id
  route_name  = "cspec-arr-route"
}

output "alicloud_apig_route_example_id" {
  value = data.alicloud_apig_routes.default.routes.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Computed) A list of Route IDs. The value is formulated as `<http_api_id>:<route_id>`.
* `name_regex` - (Optional) A regex string to filter results by Route name.
* `environment_info` - (Optional) The environment information of the route. See [`environment_info`](#environment_info) below.
* `http_api_id` - (Required) The ID of the HTTP API to which the route belongs.
* `route_name` - (Optional) The name of the route.
* `status` - (Optional) The deployment status of the route.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

---

### `environment_info`

The environment_info supports the following:

* `environment_id` - (Optional) The environment ID.
* `alias` - (Computed) The alias of the environment name.
* `name` - (Computed) The environment name.
* `gateway_info` - (Computed) The gateway instance information corresponding to the environment. See [`gateway_info`](#environment_info-gateway_info) below.
* `sub_domains` - (Computed) The default second-level domain names of the environment. See [`sub_domains`](#environment_info-sub_domains) below.

---

### `environment_info-gateway_info`

The environment_info-gateway_info supports the following:

* `gateway_edition` - (Computed) The edition of the gateway instance.
* `gateway_id` - (Computed) The ID of the Cloud-native API Gateway.
* `name` - (Computed) The name of the gateway.

---

### `environment_info-sub_domains`

The environment_info-sub_domains supports the following:

* `domain_id` - (Computed) The ID of the second-level domain name.
* `network_type` - (Computed) The domain access type, such as Intranet or Internet.
* `protocol` - (Computed) The domain protocol, such as HTTP or HTTPS.
* `name` - (Computed) The name of the second-level domain name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Route IDs.
* `names` - A list of name of Routes.
* `routes` - A list of Route Entries.

---

### `routes`

The routes supports the following attributes:

* `backend` - Backend service.
* `builtin` - Indicates whether the route is a built-in route.
* `create_time` - The creation time in UTC format: yyyy-MM-ddTHH:mm:ssZ.
* `description` - The description of the route.
* `domain_infos` - The domain name information.
* `environment_info` - The environment information of the route.
* `gateway_status` - The publishing status of the route on each gateway.
* `match` - The route match rule.
* `route_id` - The route ID.
* `route_name` - The name of the route.
* `status` - The deployment status of the route.
* `update_time` - The update time in Greenwich Mean Time (GMT).
* `id` - The ID of the resource supplied above.

---

### `routes.backend`

The backend supports the following attributes:

* `scene` - The backend service scenario.
* `services` - Backend service.

---

### `routes.backend.services`

The services supports the following attributes:

* `name` - The name of the service.
* `port` - Service port.
* `protocol` - Service protocol.
* `service_id` - Service ID.
* `version` - Service version.
* `weight` - The percentage value of the traffic ratio.

---

### `routes.domain_infos`

The domain_infos supports the following attributes:

* `domain_id` - The domain name ID.
* `name` - The domain name.
* `protocol` - The domain name protocol: HTTPS or HTTP.

---

### `routes.environment_info`

The environment_info supports the following attributes:

* `alias` - The alias of the environment name.
* `environment_id` - The environment ID.
* `gateway_info` - The gateway instance information corresponding to the environment.
* `name` - The environment name.
* `sub_domains` - The default second-level domain names of the environment.

---

### `routes.environment_info.gateway_info`

The gateway_info supports the following attributes:

* `gateway_edition` - The edition of the gateway instance.
* `gateway_id` - The ID of the Cloud-native API Gateway.
* `name` - The name of the gateway.

---

### `routes.environment_info.sub_domains`

The sub_domains supports the following attributes:

* `domain_id` - The ID of the second-level domain name.
* `name` - The name of the second-level domain name.
* `network_type` - The domain access type, such as Intranet or Internet.
* `protocol` - The domain protocol, such as HTTP or HTTPS.

---

### `routes.match`

The match supports the following attributes:

* `headers` - The list of HTTP request header matching rules.
* `ignore_uri_case` - Specifies whether the path is case-sensitive.
* `methods` - The request method.
* `path` - The path rule.
* `query_params` - The matching rules for query parameters.

---

### `routes.match.headers`

The headers supports the following attributes:

* `name` - The header name.
* `type` - The header matching rule type.
* `value` - The header value.

---

### `routes.match.path`

The path supports the following attributes:

* `type` - The path matching type.
* `value` - The path value.

---

### `routes.match.query_params`

The query_params supports the following attributes:

* `name` - The parameter name.
* `type` - The matching rule for the query parameter.
* `value` - The parameter value.
