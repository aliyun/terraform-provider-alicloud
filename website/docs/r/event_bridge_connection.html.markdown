---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_connection"
description: |-
  Provides a Alicloud Event Bridge Connection resource.
---

# alicloud_event_bridge_connection

Provides a Event Bridge Connection resource. 

For information about Event Bridge Connection and how to use it, see [What is Connection](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.207.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_event_bridge_connection" "default" {
  connection_name = var.name
  description     = "test-connection-basic-pre"
  network_parameters {
    network_type      = "PublicNetwork"
    vpc_id            = "eb-cn-huhehaote/vpc-hp3bdy0vbee0vb87fq2i6"
    vswitche_id       = "vsw-hp3uinuttt9qbl27482v9"
    security_group_id = "/sg-hp37abv61w1mpsuc0zco"
  }
  auth_parameters {
    api_key_auth_parameters {
      api_key_name  = "Token"
      api_key_value = "Token-value"
    }
    basic_auth_parameters {
      password = "admin"
      username = "admin"
    }
    oauth_parameters {
      authorization_endpoint = "http://127.0.0.1:8080"
      client_parameters {
        client_secret = "ClientSecret"
        client_id     = "ClientId"
      }
      http_method = "POST"
      oauth_http_parameters {
        body_parameters {
          is_value_secret = "true"
          key             = "name"
          value           = "name"
        }
        header_parameters {
          is_value_secret = "true"
          key             = "name"
          value           = "name"
        }
        query_string_parameters {
          is_value_secret = "true"
          key             = "name"
          value           = "name"
        }
      }
    }
    authorization_type = "BASIC_AUTH"
  }
}
```

## Argument Reference

The following arguments are supported:
* `auth_parameters` - (Optional) Authentication Data Structure. See [`auth_parameters`](#auth_parameters) below.
* `description` - (Optional) The connection configuration description. The maximum length is 255 characters.
* `network_parameters` - (Required) Network Configuration Data Structure. See [`network_parameters`](#network_parameters) below.

### `api_key_auth_parameters`

The api_key_auth_parameters supports the following:
* `api_key_name` - (Optional) The key value of Api key.
* `api_key_value` - (Optional) Value of Api key.

### `basic_auth_parameters`

The basic_auth_parameters supports the following:
* `password` - (Optional) Password for basic authentication.
* `username` - (Optional) Username for basic authentication.

### `client_parameters`

The client_parameters supports the following:
* `client_id` - (Optional) The ID of the client.
* `client_secret` - (Optional) Application's client key secret.

### `body_parameters`

The body_parameters supports the following:
* `is_value_secret` - (Optional) Whether it is Authentication.
* `key` - (Optional) The key of the body request parameter.
* `value` - (Optional) The value of the body request parameter.

### `header_parameters`

The header_parameters supports the following:
* `is_value_secret` - (Optional) Whether it is Authentication.
* `key` - (Optional) The parameter key of the request header.
* `value` - (Optional) Request header parameter value.

### `query_string_parameters`

The query_string_parameters supports the following:
* `is_value_secret` - (Optional) Whether it is Authentication.
* `key` - (Optional) Request path parameter key.
* `value` - (Optional) Request path parameter value.

### `oauth_http_parameters`

The oauth_http_parameters supports the following:
* `body_parameters` - (Optional) Body request parameter data structure List. See [`body_parameters`](#body_parameters) below.
* `header_parameters` - (Optional) Parameter list of request header. See [`header_parameters`](#header_parameters) below.
* `query_string_parameters` - (Optional) Data structure of request path parameters. See [`query_string_parameters`](#query_string_parameters) below.

### `oauth_parameters`

The oauth_parameters supports the following:
* `authorization_endpoint` - (Optional) Authorized endpoint address. The maximum length is 127 characters.
* `client_parameters` - (Optional) Customer Parameter Data Structure. See [`client_parameters`](#client_parameters) below.
* `http_method` - (Optional) The method of the probe type. Value:
  - GET
  - POST
  - HEAD
  - DELETE
  - PUT
  - PATCH.
* `oauth_http_parameters` - (Optional) Request parameters for Oauth Authentication. See [`oauth_http_parameters`](#oauth_http_parameters) below.

### `auth_parameters`

The auth_parameters supports the following:
* `api_key_auth_parameters` - (Optional, ForceNew) API KEY data structure. See [`api_key_auth_parameters`](#api_key_auth_parameters) below.
* `authorization_type` - (Optional) Authentication type:

BASIC:BASIC_AUTH

Introduction: This authorization method is the basic authorization method implemented by the browser in compliance with the http protocol. In the process of communication with the HTTP protocol, the HTTP protocol defines the method by which the basic authentication allows the HTTP server to carry out user identity cards on the client. In the request header, add Authorization :Basic empty Base64 encryption (user name: password) fixed format.


1. Username and Password are required

API KEY :API_KEY_AUTH

Introduction:
Fixed format Add Token :Token value in request header
  - ApiKeyName and ApiKeyValue are required

OAUTH :OAUTH_AUTH

Introduction:
OAuth2.0 is an authorization mechanism. Normally, for systems that do not use authorization mechanisms such as OAuth2.0, the client can directly Access the resources of the Resource Server. In order for users to Access data safely, an Access Token mechanism is added in the middle of the Access. Clients need to carry Access tokens to Access protected resources. Therefore, OAuth2.0 ensures that resources are not accessed by malicious clients, thus improving the security of the system.
  - AuthorizationEndpoint, oauthttpparameters, and HttpMethod are required.
* `basic_auth_parameters` - (Optional) Basic authentication data structure. See [`basic_auth_parameters`](#basic_auth_parameters) below.
* `oauth_parameters` - (Optional) OAuth Authentication parameter data structure. See [`oauth_parameters`](#oauth_parameters) below.

### `network_parameters`

The network_parameters supports the following:
* `network_type` - (Required) Public network: PublicNetwork

Private network: PrivateNetwork
  - Tip: When selecting a VPC, VpcId, VswitcheId, and SecurityGroupId are required.
* `security_group_id` - (Optional) Security group ID.
* `vpc_id` - (Optional) The ID of the VPC.
* `vswitche_id` - (Optional) Switch ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Connection.
* `delete` - (Defaults to 5 mins) Used when delete the Connection.
* `update` - (Defaults to 5 mins) Used when update the Connection.

## Import

Event Bridge Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_connection.example <id>
```