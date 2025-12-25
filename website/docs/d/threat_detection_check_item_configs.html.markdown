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
}

output "alicloud_threat_detection_check_item_config_example_check_id" {
  value = data.alicloud_threat_detection_check_item_configs.default.configs.0.check_id
}
```

## Argument Reference

The following arguments are supported:
* `lang` - (ForceNew, Optional) The language of the content within the request and response. Default value: **zh**. Valid value:*   **zh**: Chinese*   **en**: English
* `page_number` - (ForceNew, Optional) Current page number.
* `page_size` - (ForceNew, Optional) Number of records per page.
* `task_sources` - (ForceNew, Optional) List of task sources.
* `ids` - (Optional, ForceNew, Computed) A list of Check Item Config IDs. 
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Check Item Config IDs.
* `configs` - A list of Check Item Config Entries. Each element contains the following attributes:
  * `check_id` - The ID of the check item
  * `check_show_name` - The name of the check item.
  * `check_type` - The source type of the Situation Awareness check item. Value:- **CUSTOM**: user-defined- **SYSTEM**: Predefined by the situational awareness platform
  * `custom_configs` - The custom configuration items of the check item.
    * `default_value` - The default value of the custom configuration item. The value is a string.
    * `name` - The name of the custom configuration item, which is unique in a check item.
    * `show_name` - The display name of the custom configuration item for internationalization.
    * `type_define` - The type of the custom configuration item. The value is a JSON string.
    * `value` - The value of the custom configuration item. The value is a string.
  * `description` - The description of the check item.
    * `type` - The type of the description of the check item. Valid value:*   **text**.
    * `value` - The content of the description for the check item when the Type parameter is text.
  * `estimated_count` - The estimated quota that will be consumed by this check item.
  * `instance_sub_type` - The asset subtype of the cloud service. Valid values:*   If **InstanceType** is set to **ECS**, this parameter supports the following valid values:    *   **INSTANCE**    *   **DISK**    *   **SECURITY_GROUP***   If **InstanceType** is set to **ACR**, this parameter supports the following valid values:    *   **REPOSITORY_ENTERPRISE**    *   **REPOSITORY_PERSON***   If **InstanceType** is set to **RAM**, this parameter supports the following valid values:    *   **ALIAS**    *   **USER**    *   **POLICY**    *   **GROUP***   If **InstanceType** is set to **WAF**, this parameter supports the following valid value:    *   **DOMAIN***   If **InstanceType** is set to other values, this parameter supports the following valid values:    *   **INSTANCE**
  * `instance_type` - The asset type of the cloud service. Valid values:*   **ECS**: Elastic Compute Service (ECS).*   **SLB**: Server Load Balancer (SLB).*   **RDS**: ApsaraDB RDS.*   **MONGODB**: ApsaraDB for MongoDB (MongoDB).*   **KVSTORE**: ApsaraDB for Redis (Redis).*   **ACR**: Container Registry.*   **CSK**: Container Service for Kubernetes (ACK).*   **VPC**: Virtual Private Cloud (VPC).*   **ACTIONTRAIL**: ActionTrail.*   **CDN**: Alibaba Cloud CDN (CDN).*   **CAS**: Certificate Management Service (formerly SSL Certificates Service).*   **RDC**: Apsara Devops.*   **RAM**: Resource Access Management (RAM).*   **DDOS**: Anti-DDoS.*   **WAF**: Web Application Firewall (WAF).*   **OSS**: Object Storage Service (OSS).*   **POLARDB**: PolarDB.*   **POSTGRESQL**: ApsaraDB RDS for PostgreSQL.*   **MSE**: Microservices Engine (MSE).*   **NAS**: File Storage NAS (NAS).*   **SDDP**: Sensitive Data Discovery and Protection (SDDP).*   **EIP**: Elastic IP Address (EIP).
  * `risk_level` - The risk level of the check item. Valid values:*   **HIGH***   **MEDIUM***   **LOW**
  * `section_ids` - The IDs of the sections associated with the check items.
  * `vendor` - The type of the cloud asset. Valid values:*   **0**: an asset provided by Alibaba Cloud.*   **1**: an asset outside Alibaba Cloud.*   **2**: an asset in a data center.*   **3**, **4**, **5**, and **7**: other cloud asset.*   **8**: a simple application server.
