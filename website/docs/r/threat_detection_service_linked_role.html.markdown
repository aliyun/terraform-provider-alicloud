---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_service_linked_role"
description: |-
  Provides an Alicloud Threat Detection Service Linked Role resource.
---

# alicloud_threat_detection_service_linked_role

Provides a Threat Detection Service Linked Role resource.

Service Linked Role.

For information about Threat Detection Service Linked Role and how to use it, see [What is Service Linked Role](https://www.alibabacloud.com/help/en/doc-detail/42302.htm).

-> **NOTE:** Available since v1.283.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_service_linked_role&exampleId=104f63ce-fd71-dac0-3353-8fdd7049ece23a337d91&activeTab=example&spm=docs.r.threat_detection_service_linked_role.0.104f63cefd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_threat_detection_service_linked_role" "default" {
  service_linked_role = "AliyunServiceRoleForSas"
}

```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_threat_detection_service_linked_role&spm=docs.r.threat_detection_service_linked_role.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `service_linked_role` - (Optional, ForceNew, Computed, Available since v1.283.0) The name of the service linked role.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `role_status` - Whether the service linked role exists.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Service Linked Role.
* `delete` - (Defaults to 5 mins) Used when delete the Service Linked Role.

## Import

Threat Detection Service Linked Role can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_service_linked_role.example <service_linked_role>
```
