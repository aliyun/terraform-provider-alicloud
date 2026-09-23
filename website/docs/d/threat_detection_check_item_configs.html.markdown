---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_check_item_configs"
sidebar_current: "docs-alicloud-datasource-threat-detection-check-item-configs"
description: |-
  Provides a list of Threat Detection Check Item Config owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_check_item_configs

This data source provides Threat Detection Check Item Config available to the user.[What is Check Item Config](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListCheckItem)

-> **NOTE:** Available since v1.267.0.

## Example Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_threat_detection_check_item_configs" "default" {
  lang        = "zh"
  page_number = 1
  page_size   = 10
}

output "alicloud_threat_detection_check_item_config_example_check_id" {
  value = data.alicloud_threat_detection_check_item_configs.default.configs.0.check_id
}
```

## Argument Reference

The following arguments are supported:
* `lang` - (Optional) The language of the content within the request and response. Default value: `zh`. Valid values: `zh` (Chinese), `en` (English).
* `page_number` - (Optional) The page number. Must be greater than 0.
* `page_size` - (Optional) Number of records per page. Must be greater than 0.
* `task_sources` - (Optional) List of task sources.
* `ids` - (Optional, Computed) A list of Check Item Config IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Check Item Config IDs.
* `configs` - A list of Check Item Config Entries. Each element contains the following attributes:
    * `check_id` - The ID of the check item.
    * `check_show_name` - The name of the check item.
    * `check_type` - The source type of the Situation Awareness check item. Valid values: `CUSTOM` (user-defined), `SYSTEM` (predefined by the situational awareness platform).
    * `custom_configs` - The custom configuration items of the check item.
      * `default_value` - The default value of the custom configuration item. The value is a string.
      * `name` - The name of the custom configuration item, which is unique in a check item.
      * `show_name` - The display name of the custom configuration item for internationalization.
      * `type_define` - The type of the custom configuration item. The value is a JSON string.
      * `value` - The value of the custom configuration item. The value is a string.
    * `description` - The description of the check item.
      * `type` - The type of the description of the check item. Valid value: `text`.
      * `value` - The content of the description for the check item when the Type parameter is text.
    * `estimated_count` - The estimated quota that will be consumed by this check item.
    * `instance_sub_type` - The asset subtype of the cloud service. Valid values depend on `instance_type`:
      * `ECS`: `INSTANCE`, `DISK`, `SECURITY_GROUP`.
      * `ACR`: `REPOSITORY_ENTERPRISE`, `REPOSITORY_PERSON`.
      * `RAM`: `ALIAS`, `USER`, `POLICY`, `GROUP`.
      * `WAF`: `DOMAIN`.
      * Other values: `INSTANCE`.
    * `instance_type` - The asset type of the cloud service. Valid values:
      * `ECS`: Elastic Compute Service.
      * `SLB`: Server Load Balancer.
      * `RDS`: ApsaraDB RDS.
      * `MONGODB`: ApsaraDB for MongoDB.
      * `KVSTORE`: ApsaraDB for Redis.
      * `ACR`: Container Registry.
      * `CSK`: Container Service for Kubernetes.
      * `VPC`: Virtual Private Cloud.
      * `ACTIONTRAIL`: ActionTrail.
      * `CDN`: Alibaba Cloud CDN.
      * `CAS`: Certificate Management Service.
      * `RDC`: Apsara DevOps.
      * `RAM`: Resource Access Management.
      * `DDOS`: Anti-DDoS.
      * `WAF`: Web Application Firewall.
      * `OSS`: Object Storage Service.
      * `POLARDB`: PolarDB.
      * `POSTGRESQL`: ApsaraDB RDS for PostgreSQL.
      * `MSE`: Microservices Engine.
      * `NAS`: File Storage NAS.
      * `SDDP`: Sensitive Data Discovery and Protection.
      * `EIP`: Elastic IP Address.
    * `risk_level` - The risk level of the check item. Valid values: `HIGH`, `MEDIUM`, `LOW`.
    * `section_ids` - The IDs of the sections associated with the check items.
    * `vendor` - The type of the cloud asset. Valid values: `0` (an asset provided by Alibaba Cloud), `1` (an asset outside Alibaba Cloud), `2` (an asset in a data center), `3`/`4`/`5`/`7` (other cloud asset), `8` (a simple application server).
