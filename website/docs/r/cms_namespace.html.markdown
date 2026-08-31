---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_namespace"
sidebar_current: "docs-alicloud-resource-cms-namespace"
description: |-
  Provides a Alicloud Cloud Monitor Service Namespace resource.
---

# alicloud_cms_namespace

Provides a Cloud Monitor Service Namespace resource.

For information about Cloud Monitor Service Namespace and how to use it, see [What is Namespace](https://www.alibabacloud.com/help/en/cloudmonitor/latest/createhybridmonitornamespace).

-> **NOTE:** Available since v1.171.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_namespace&exampleId=21824ffe-22ff-339a-af06-13bd61a2ebcee92bf83a&activeTab=example&spm=docs.r.cms_namespace.0.21824ffe22&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_cms_namespace" "example" {
  namespace     = "tf-example"
  specification = "cms.s1.large"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_namespace&spm=docs.r.cms_namespace.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, ForceNew) The name of the namespace. The name can contain lowercase letters, digits, and hyphens (-).
* `specification` - (Optional) The data retention period. Default value: `cms.s1.3xlarge`. Valid values:
  - `cms.s1.large`: Data storage duration is 15 days.
  - `cms.s1.xlarge`: Data storage duration is 32 days.
  - `cms.s1.2xlarge`: Data storage duration 63 days.
  - `cms.s1.3xlarge`: Data storage duration 93 days.
  - `cms.s1.6xlarge`: Data storage duration 185 days.
  - `cms.s1.12xlarge`: Data storage duration 376 days.
* `description` - (Optional) The description of the namespace.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Namespace. Its value is same as `namespace`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 min) Used when create the Namespace.
* `update` - (Defaults to 1 min) Used when update the Namespace.
* `delete` - (Defaults to 1 min) Used when delete the Namespace.

## Import

Cloud Monitor Service Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_namespace.example <id>
```
