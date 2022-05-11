---
subcategory: "Apsara Agile Live (IMP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_imp_app_template"
sidebar_current: "docs-alicloud-resource-imp-app-template"
description: |-
  Provides a Alicloud Apsara Agile Live (IMP) App Template resource.
---

# alicloud\_imp\_app\_template

Provides a Apsara Agile Live (IMP) App Template resource.

For information about Apsara Agile Live (IMP) App Template and how to use it, see [What is App Template](https://help.aliyun.com/document_detail/270121.html).

-> **NOTE:** Available in v1.137.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_imp_app_template" "example" {
  app_template_name = "example_value"
  component_list    = ["component.live", "component.liveRecord"]
  integration_mode  = "paasSDK"
  scene             = "business"
}

```

## Argument Reference

The following arguments are supported:

* `app_template_name` - (Required) The name of the resource.
* `component_list` - (Required, ForceNew) List of components. Its element valid values: ["component.live","component.liveRecord","component.liveBeauty","component.rtc","component.rtcRecord","component.im","component.whiteboard","component.liveSecurity","component.chatSecurity"].
* `config_list` - (Optional, Computed) Configuration list. It have several default configs after the resource is created. See the following `Block config_list`.
* `integration_mode` - (Optional, ForceNew) Integration mode. Valid values:
  * paasSDK: Integrated SDK.
  * standardRoom: Model Room.
  
* `scene` - (Optional, ForceNew) Application Template scenario. Valid values: ["business", "classroom"].

#### Block config_list

The config_list supports the following: 

* `key` - (Optional) Configuration item key. Valid values: ["config.appCallbackAuthKey","config.appCallbackUrl","config.callbackClass.live","config.callbackClass.user","config.livePullDomain","config.livePushDomain","config.multipleClientsLogin","config.regionId","config.streamChangeCallbackUrl"].
* `value` - (Optional) Configuration item content.

-> **NOTE:**  By default, the attribute `config_list` will return all of nine keys with empty value. If you want to set one or more of the key's value, you had better also set other keys, otherwise, there will be a diff.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of App Template.
* `status` - Application template usage status.

## Import

Apsara Agile Live (IMP) App Template can be imported using the id, e.g.

```
$ terraform import alicloud_imp_app_template.example <id>
```
