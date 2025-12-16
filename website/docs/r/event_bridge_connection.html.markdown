---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_connection"
sidebar_current: "docs-alicloud-resource-event-bridge-connection"
description: |-
  Provides a Alicloud Event Bridge Connection resource.
---

# alicloud_event_bridge_connection

Provides a Event Bridge Connection resource.

For information about Event Bridge Connection and how to use it, see [What is Connection](https://www.alibabacloud.com/help/en/eventbridge/latest/api-eventbridge-2020-04-01-createconnection).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_event_bridge_connection&exampleId=5f2132f4-75c6-44d6-eab9-2bca3ea2855abef271a3&activeTab=example&spm=docs.r.event_bridge_connection.0.5f2132f475&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = var.region
}

variable "region" {
  default = "cn-chengdu"
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vswitch.default.vpc_id
}

resource "alicloud_event_bridge_connection" "default" {
  connection_name = var.name
  description     = "test-connection-basic-pre"
  network_parameters {
    network_type      = "PublicNetwork"
    vpc_id            = alicloud_vpc.default.id
    vswitche_id       = alicloud_vswitch.default.id
    security_group_id = alicloud_security_group.default.id
  }
  auth_parameters {
    authorization_type = "BASIC_AUTH"
    api_key_auth_parameters {
      api_key_name  = "Token"
      api_key_value = "Token-value"
    }
    basic_auth_parameters {
      username = "admin"
      password = "admin"
    }
    oauth_parameters {
      authorization_endpoint = "http://127.0.0.1:8080"
      http_method            = "POST"
      client_parameters {
        client_id     = "ClientId"
        client_secret = "ClientSecret"
      }
      oauth_http_parameters {
        header_parameters {
          key             = "name"
          value           = "name"
          is_value_secret = "true"
        }
        body_parameters {
          key             = "name"
          value           = "name"
          is_value_secret = "true"
        }
        query_string_parameters {
          key             = "name"
          value           = "name"
          is_value_secret = "true"
        }
      }
    }
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_event_bridge_connection&spm=docs.r.event_bridge_connection.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `connection_name` - (Required, ForceNew) The name of the connection.
* `description` - (Optional) The description of the connection.
* `network_parameters` - (Required, Set) The parameters that are configured for the network. See [`network_parameters`](#network_parameters) below.
* `auth_parameters` - (Optional, Set) The parameters that are configured for authentication. See [`auth_parameters`](#auth_parameters) below.

### `network_parameters`

The network_parameters supports the following:

* `network_type` - (Required) The network type. Valid values: `PublicNetwork`, `PrivateNetwork`. **NOTE:** If you set `network_type` to `PrivateNetwork`, you must configure `vpc_id`, `vswitche_id`, and `security_group_id`.
* `vpc_id` - (Optional) The ID of the VPC.
* `vswitche_id` - (Optional) The ID of the VSwitch.
* `security_group_id` - (Optional) The ID of the security group.

### `auth_parameters`

The auth_parameters supports the following:

* `authorization_type` - (Optional) The type of the authentication. Valid values: `API_KEY_AUTH`, `BASIC_AUTH`, `OAUTH_AUTH`.
* `api_key_auth_parameters` - (Optional, Set) The parameters that are configured for API key authentication. See [`api_key_auth_parameters`](#auth_parameters-api_key_auth_parameters) below.
* `basic_auth_parameters` - (Optional, Set) The parameters that are configured for basic authentication. See [`basic_auth_parameters`](#auth_parameters-basic_auth_parameters) below.
* `oauth_parameters` - (Optional, Set) The parameters that are configured for OAuth authentication. See [`oauth_parameters`](#auth_parameters-oauth_parameters) below.

### `auth_parameters-api_key_auth_parameters`

The api_key_auth_parameters supports the following:

* `api_key_name` - (Optional) The name of the API key.
* `api_key_value` - (Optional) The value of the API key.

### `auth_parameters-basic_auth_parameters`

The basic_auth_parameters supports the following:

* `username` - (Optional) The username for basic authentication.
* `password` - (Optional) The password for basic authentication.

### `auth_parameters-oauth_parameters`

The oauth_parameters supports the following:

* `authorization_endpoint` - (Optional) The IP address of the authorized endpoint.
* `http_method` - (Optional) The HTTP request method. Valid values: `GET`, `POST`, `HEAD`, `DELETE`, `PUT`, `PATCH`.
* `client_parameters` - (Optional, Set) The parameters that are configured for the client. See [`client_parameters`](#auth_parameters-oauth_parameters-client_parameters) below.
* `oauth_http_parameters` - (Optional, Set) The request parameters that are configured for OAuth authentication. See [`oauth_http_parameters`](#auth_parameters-oauth_parameters-oauth_http_parameters) below.

### `auth_parameters-oauth_parameters-client_parameters`

The client_parameters supports the following:

* `client_id` - (Optional) The ID of the client.
* `client_secret` - (Optional) The AccessKey secret of the client.


### `auth_parameters-oauth_parameters-oauth_http_parameters`

The oauth_http_parameters supports the following:

* `header_parameters` - (Optional, Set) The parameters that are configured for the request header. See [`header_parameters`](#auth_parameters-oauth_parameters-oauth_http_parameters-header_parameters) below.
* `body_parameters` - (Optional, Set) The parameters that are configured for the request body. See [`body_parameters`](#auth_parameters-oauth_parameters-oauth_http_parameters-body_parameters) below.
* `query_string_parameters` - (Optional, Set) The parameters that are configured for the request path. See [`query_string_parameters`](#auth_parameters-oauth_parameters-oauth_http_parameters-query_string_parameters) below.

### `auth_parameters-oauth_parameters-oauth_http_parameters-header_parameters`

The header_parameters supports the following:

* `key` - (Optional) The key of the request header.
* `value` - (Optional) The value of the request header.
* `is_value_secret` - (Optional) Specifies whether to enable authentication.

### `auth_parameters-oauth_parameters-oauth_http_parameters-body_parameters`

The body_parameters supports the following:

* `key` - (Optional) The key of the request body.
* `value` - (Optional) The value of the request body.
* `is_value_secret` - (Optional) Specifies whether to enable authentication.

### `auth_parameters-oauth_parameters-oauth_http_parameters-query_string_parameters`

The query_string_parameters supports the following:

* `key` - (Optional) The key of the request path.
* `value` - (Optional) The key of the request path.
* `is_value_secret` - (Optional) Specifies whether to enable authentication.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Connection.
* `create_time` - The creation time of the Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Connection.
* `update` - (Defaults to 5 mins) Used when update the Connection.
* `delete` - (Defaults to 5 mins) Used when delete the Connection.

## Import

Event Bridge Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_connection.example <id>
```
