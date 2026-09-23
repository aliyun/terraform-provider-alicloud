---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_kv_namespace"
description: |-
  Provides a Alicloud ESA Kv Namespace resource.
---

# alicloud_esa_kv_namespace

Provides a ESA Kv Namespace resource.



For information about ESA Kv Namespace and how to use it, see [What is Kv Namespace](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateKvNamespace).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_kv_namespace&exampleId=700a0f69-b96a-3aa0-0d44-445a4aba053770793022&activeTab=example&spm=docs.r.esa_kv_namespace.0.700a0f69b9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_esa_kv_namespace" "default" {
  description  = "this is a example namespace."
  kv_namespace = "example_namespace"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_kv_namespace&spm=docs.r.esa_kv_namespace.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) The description of the namespace.
* `kv_namespace` - (Required, ForceNew) The name of the namespace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the namespace. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Kv Namespace.
* `delete` - (Defaults to 5 mins) Used when delete the Kv Namespace.

## Import

ESA Kv Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_kv_namespace.example <id>
```