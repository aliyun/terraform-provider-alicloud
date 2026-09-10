---
subcategory: "Live"
layout: "alicloud"
page_title: "Alicloud: alicloud_live_domain"
description: |-
  Provides a Alicloud Live Domain resource.
---

# alicloud_live_domain

Provides a Live Domain resource.

Live domain name.

For information about Live Domain and how to use it, see [What is Domain](https://next.api.alibabacloud.com/document/live/2016-11-01/AddLiveDomain).

-> **NOTE:** Available since v1.272.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_live_domain&exampleId=1f2838fb-dded-56b4-3c02-8151c41f11862dcc6e05&activeTab=example&spm=docs.r.live_domain.0.1f2838fbdd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "domain_name" {
  default = "demo.alicloud.com"
}

resource "alicloud_live_domain" "default" {
  domain_type = "liveVideo"
  scope       = "overseas"
  domain_name = var.domain_name
  region      = "cn-shanghai"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_live_domain&spm=docs.r.live_domain.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `check_url` - (Optional) Health check URL.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `domain_name` - (Required, ForceNew) Fuzzy match filter for domain names.

-> **NOTE:** - If `domain_type` (live streaming domain business type) is set to `liveVideo`, and this parameter is not specified, the system queries information about the user's playback domains by default.

-> **NOTE:** - If `domain_type` is set to `liveEdge`, and this parameter is not specified, the system queries information about the user's ingest domains by default.

* `domain_type` - (Required, ForceNew) Domain business type. Valid values:  
  - `liveVideo`: Playback domain.  
  - `liveEdge`: Edge ingest domain.
* `region` - (Required, ForceNew) Region to which the domain belongs.
* `resource_group_id` - (Optional, Computed) Resource group ID. For more information about resource groups, see [What is a resource group?](https://help.aliyun.com/document_detail/2381067.html).
* `scope` - (Optional, Computed) Acceleration region. This parameter takes effect only for international users and China site users with L3 or higher privileges. Valid values:  
  - `domestic` (default): Mainland China.  
  - `overseas`: Acceleration in regions outside mainland China, including Hong Kong, Macao, and Taiwan.  
  - `global`: Global acceleration.
* `status` - (Optional, Computed) Domain status. Valid values:
  - `online`: Running (indicating that the domain name service is operating normally).
  - `offline`: Stopped.
* `tags` - (Optional, Map) List of tags.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - Creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 15 mins) Used when create the Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Domain.
* `update` - (Defaults to 5 mins) Used when update the Domain.

## Import

Live Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_live_domain.example <domain_name>
```
