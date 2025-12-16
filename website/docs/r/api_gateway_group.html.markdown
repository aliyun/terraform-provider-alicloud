---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_group"
sidebar_current: "docs-alicloud-resource-api-gateway-group"
description: |-
  Provides a Alicloud Api Gateway Group Resource.
---

# alicloud_api_gateway_group

Provides an api group resource.To create an API, you must firstly create a group which is a basic attribute of the API.

For information about Api Gateway Group and how to use it, see [Create An Api Group](https://www.alibabacloud.com/help/en/api-gateway/latest/api-cloudapi-2016-07-14-createapigroup)

-> **NOTE:** Available since v1.19.0.

-> **NOTE:** Terraform will auto build api group while it uses `alicloud_api_gateway_group` to build api group.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_group&exampleId=2980c309-cddc-7ed2-31fd-ed7e60173a3b74267fb2&activeTab=example&spm=docs.r.api_gateway_group.0.2980c309cd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_api_gateway_group" "default" {
  name        = "tf_example"
  description = "tf_example"
  base_path   = "/"
  user_log_config {
    request_body     = true
    response_body    = true
    query_string     = "*"
    request_headers  = "*"
    response_headers = "*"
    jwt_claims       = "*"
  }
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_api_gateway_group&spm=docs.r.api_gateway_group.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the api gateway group. Defaults to null.
* `description` - (Optional) The description of the api gateway group. Defaults to null.
* `instance_id` - (Optional, ForceNew, Available in 1.179.0+)	The id of the api gateway.
* `base_path` - (Optional, Computed, Available since v1.228.0) The base path of the api gateway group. Defaults to `/`.
* `user_log_config` - (Optional, Available since v1.246.0) user_log_config defines the config of user log of the group. See [`user_log_config`](#user_log_config) below.
* `vpc_intranet_enable` - (Optional, Available since v1.247.0) Whether to enable `vpc_domain`. Defaults to `false`.

### `user_log_config`

The user_log_config mapping supports the following:

* `request_body` - (Optional, Type: bool) Whether to record the request body.
* `response_body` - (Optional, Type: bool) Whether to record the response body.
* `query_string` - (Optional) The query params to be record, support multi query params split by `,`. Set `*` to record all.
* `request_headers` - (Optional) The request headers to be record, support multi request headers split by `,`. Set `*` to record all.
* `response_headers` - (Optional) The response headers to be record, support multi response headers split by `,`. Set `*` to record all.
* `jwt_claims` - (Optional) The jwt claims to be record, support multi jwt claims split by `,`. Set `*` to record all.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the api group of api gateway.
* `sub_domain` - (Available in 1.69.0+)	Second-level domain name automatically assigned to the API group.
* `vpc_domain` - (Available in 1.69.0+)	Second-level VPC domain name automatically assigned to the API group.

## Import

Api gateway group can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_group.example "ab2351f2ce904edaa8d92a0510832b91"
```
