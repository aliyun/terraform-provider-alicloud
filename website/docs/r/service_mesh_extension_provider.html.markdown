---
subcategory: "Service Mesh"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_mesh_extension_provider"
sidebar_current: "docs-alicloud-resource-service-mesh-extension-provider"
description: |-
  Provides a Alicloud Service Mesh Extension Provider resource.
---

# alicloud_service_mesh_extension_provider

Provides a Service Mesh Extension Provider resource.

For information about Service Mesh Extension Provider and how to use it, see [What is Extension Provider](https://help.aliyun.com/document_detail/461549.html).

-> **NOTE:** Available since v1.191.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_service_mesh_extension_provider&exampleId=90824bc7-c128-f7c7-3ea3-1839e905b5389002b626&activeTab=example&spm=docs.r.service_mesh_extension_provider.0.90824bc7c1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf_example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_vpc" "default" {
  count = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
}

data "alicloud_vswitches" "default" {
  vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
}

resource "alicloud_vswitch" "default" {
  count      = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id     = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_service_mesh_service_mesh" "default" {
  service_mesh_name = "mesh-c50f3fef117ad45b6b26047cdafef65ad"
  version           = "v1.21.6.103-g5ddeaef7-aliyun"
  edition           = "Default"
  network {
    vpc_id        = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
    vswitche_list = [length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id]
  }
}

resource "alicloud_service_mesh_extension_provider" "default" {
  service_mesh_id         = alicloud_service_mesh_service_mesh.default.id
  extension_provider_name = "httpextauth-tf-example"
  type                    = "httpextauth"
  config                  = "{\"headersToDownstreamOnDeny\":[\"content-type\",\"set-cookie\"],\"headersToUpstreamOnAllow\":[\"authorization\",\"cookie\",\"path\",\"x-auth-request-access-token\",\"x-forwarded-access-token\"],\"includeRequestHeadersInCheck\":[\"cookie\",\"x-forward-access-token\"],\"oidc\":{\"clientID\":\"qweqweqwewqeqwe\",\"clientSecret\":\"asdasdasdasdsadas\",\"cookieExpire\":\"1000\",\"cookieRefresh\":\"500\",\"cookieSecret\":\"scxzcxzcxzcxzcxz\",\"issuerURI\":\"qweqwewqeqweqweqwe\",\"redirectDomain\":\"www.alicloud-provider.cn\",\"redirectProtocol\":\"http\",\"scopes\":[\"profile\"]},\"port\":4180,\"service\":\"oauth2proxy-httpextauth-tf-example.istio-system.svc.cluster.local\",\"timeout\":\"10s\"}"
}
```

## Argument Reference

The following arguments are supported:

* `service_mesh_id` - (Required, ForceNew) The ID of the Service Mesh.
* `extension_provider_name` - (Required, ForceNew) The name of the Service Mesh Extension Provider. It must be prefixed with `$type-`, for example `httpextauth-xxx`, `grpcextauth-xxx`.
* `type` - (Required, ForceNew) The type of the Service Mesh Extension Provider. Valid values: `httpextauth`, `grpcextauth`.
* `config` - (Required) The config of the Service Mesh Extension Provider. The `config` format is json.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Extension Provider. The value formats as `<service_mesh_id>:<type>:<extension_provider_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Extension Provider.
* `update` - (Defaults to 3 mins) Used when update the Extension Provider.
* `delete` - (Defaults to 3 mins) Used when delete the Extension Provider.

## Import

Service Mesh Extension Provider can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_mesh_extension_provider.example <service_mesh_id>:<type>:<extension_provider_name>
```
