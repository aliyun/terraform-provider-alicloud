---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_global_events_storage_region"
sidebar_current: "docs-alicloud-resource-actiontrail-global-events-storage-region"
description: |-
  Provides Alibaba Cloud Actiontrail Global Events Storage Region Resource
---

# alicloud_actiontrail_global_events_storage_region

Provides a Global events storage region resource.

For information about global events storage region and how to use it, see [What is Global Events Storage Region](https://next.api.alibabacloud.com/api/Actiontrail/2020-07-06/UpdateGlobalEventsStorageRegion).

-> **NOTE:** Available since v1.201.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_actiontrail_global_events_storage_region&exampleId=cceefb5a-1f88-b691-05e9-02d07bb3a677a4b81828&activeTab=example&spm=docs.r.actiontrail_global_events_storage_region.0.cceefb5a1f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_actiontrail_global_events_storage_region" "foo" {
  storage_region = "cn-hangzhou"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_actiontrail_global_events_storage_region&spm=docs.r.actiontrail_global_events_storage_region.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `storage_region` - (Optional) Global Events Storage Region.

## Attributes Reference

The following attributes are exported:


## Import

Global events storage region not can be imported.

