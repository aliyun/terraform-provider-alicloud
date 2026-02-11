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
