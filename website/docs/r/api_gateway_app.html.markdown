---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_app"
sidebar_current: "docs-alicloud-resource-api-gateway-app"
description: |-
  Provides a Alicloud Api Gateway App Resource.
---

# alicloud_api_gateway_app

Provides an app resource.It must create an app before calling a third-party API because the app is the identity used to call the third-party API.

For information about Api Gateway App and how to use it, see [Create An APP](https://www.alibabacloud.com/help/en/api-gateway/latest/api-cloudapi-2016-07-14-createapp)

-> **NOTE:** Available since v1.22.0.

-> **NOTE:** Terraform will auto build api app while it uses `alicloud_api_gateway_app` to build api app.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_app&exampleId=e2e295a8-1ae1-cb57-3a24-cf8efcfaa1b12242a599&activeTab=example&spm=docs.r.api_gateway_app.0.e2e295a81a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_api_gateway_app" "example" {
  name        = "tf_example"
  description = "tf_example"
  extend      = "extend-info"
  app_code    = "tf-app-code-1234"
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_api_gateway_app&spm=docs.r.api_gateway_app.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the app.
* `description` - (Optional) The description of the app. Defaults to null.
* `app_code` - (Optional, Computed) The AppCode of the app. It is settable at create and modifiable post-create via the ResetAppCode API. Once changed, the provider issues ResetAppCode with the previous value as identifier and the new value as the target AppCode.
* `extend` - (Optional, Computed) The extend info of the app. Settable at create and update.
* `disabled` - (Optional, Computed) Whether the app is disabled. Modifiable via the ModifyApp API; not accepted by CreateApp, so the provider applies it via ModifyApp right after creation when set.
* `app_secret_reset` - (Optional) Trigger field. Any non-empty change invokes the ResetAppSecret API, after which the provider refreshes `app_key` and `app_secret` into state. Use this to rotate credentials without recreating the app.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the app of api gateway.
* `app_key` - The AppKey returned by the service.
* `app_secret` - The AppSecret returned by the service (sensitive).
* `app_code` - The AppCode currently associated with the app.

## Import

Api gateway app can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_app.example "7379660"
```
