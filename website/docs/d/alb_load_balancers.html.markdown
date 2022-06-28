---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_load_balancers"
sidebar_current: "docs-alicloud-datasource-alb-load-balancers"
description: |- 
    Provides a list of Alb Load Balancers to the user.
---

# alicloud\_alb\_load\_balancers

This data source provides the Alb Load Balancers of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alb_load_balancers" "ids" {}
output "alb_load_balancer_id_1" {
  value = data.alicloud_alb_load_balancers.ids.balancers.0.id
}

data "alicloud_alb_load_balancers" "nameRegex" {
  name_regex = "^my-LoadBalancer"
}
output "alb_load_balancer_id_2" {
  value = data.alicloud_alb_load_balancers.nameRegex.balancers.0.id
}

```

## Argument Reference

The following arguments are supported:

* `address_type` - (Optional, ForceNew) The type of IP address that the ALB instance uses to provide services. Valid
  values: `Intranet`, `Internet`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Load Balancer IDs.
* `load_balancer_business_status` - (Optional, ForceNew,Available in 1.142.0+) Load Balancing of the Service Status. Valid Values: `Abnormal`and `Normal`.
* `load_balancer_ids` - (Optional, ForceNew) The load balancer ids.
* `load_balancer_name` - (Optional, ForceNew) The name of the resource.
* `status` - (Optional, ForceNew) The load balancer status. Valid values: `Active`, `Configuring`, `CreateFailed`, `Inactive` and `Provisioning`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Load Balancer name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tag` - (Optional, ForceNew) The tag.
* `vpc_id` - (Optional, ForceNew) The ID of the virtual private cloud (VPC) where the ALB instance is deployed.
* `vpc_ids` - (Optional, ForceNew) The vpc ids.
* `zone_id` - (Optional, ForceNew) The zone ID of the resource.
* `load_balancer_bussiness_status` - (Deprecated) Field 'load_balancer_bussiness_status' has been deprecated from provider version 1.142.0. Use 'load_balancer_business_status' replaces it.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Load Balancer names.
* `balancers` - A list of Alb Load Balancers. Each element contains the following attributes:
    * `access_log_config` - The Access Logging Configuration Structure.
        * `log_project` -  The log service that access logs are shipped to.
        * `log_store` - The logstore that access logs are shipped to.
    * `address_allocated_mode` - The method in which IP addresses are assigned. Valid values:  Fixed: The ALB instance
      uses a fixed IP address. Dynamic (default): An IP address is dynamically assigned to each zone of the ALB
      instance.
    * `address_type` - The type of IP address that the ALB instance uses to provide services.
    * `bandwidth_package_id` - The ID of the EIP bandwidth plan which is associated with an ALB instance that uses a
      public IP address.
    * `create_time` - The creation time of the resource.
    * `deletion_protection_config` - Remove the Protection Configuration.
        * `enabled` - Remove the Protection Status.
        * `enabled_time` - Deletion Protection Turn-on Time Use Greenwich Mean Time, in the Format of Yyyy-MM-ddTHH: mm:SSZ.
    * `dns_name` - DNS Domain Name.
    * `id` - The ID of the Load Balancer.
    * `load_balancer_billing_config` - The configuration of the billing method.
        * `pay_type` - The billing method of the ALB instance. Valid value: `PayAsYouGo`.
    * `load_balancer_bussiness_status` - Load Balancing of the Service Status. Valid Values: `Abnormal` and `Normal`.  **NOTE:** Field 'load_balancer_bussiness_status' has been deprecated from provider version 1.142.0.
    * `load_balancer_business_status` - Load Balancing of the Service Status. Valid Values: `Abnormal` and `Normal`. **NOTE:** Available in 1.142.0+
    * `load_balancer_edition` - The edition of the ALB instance.
    * `load_balancer_id` - The first ID of the resource.
    * `load_balancer_name` - The name of the resource.
    * `load_balancer_operation_locks` - The Load Balancing Operations Lock Configuration.
        * `lock_reason` - The Locking of the Reasons. 
        * `lock_type` - The Locking of the Type. Valid Values: `securitylocked`,`relatedresourcelocked`, `financiallocked`, and `residuallocked`.
    * `modification_protection_config` - Modify the Protection Configuration.
      * `status` - Specifies whether to enable the configuration read-only mode for the ALB instance. Valid values: `NonProtection` and `ConsoleProtection`.
        * `NonProtection` - disables the configuration read-only mode. After you disable the configuration read-only mode, you cannot set the ModificationProtectionReason parameter. If the parameter is set, the value is cleared.
        * `ConsoleProtection` - enables the configuration read-only mode. After you enable the configuration read-only mode, you can set the ModificationProtectionReason parameter.
      * `reason` - The reason for modification protection. This parameter must be 2 to 128 characters in length, and can contain letters, digits, periods, underscores, and hyphens. The reason must start with a letter. This parameter is required only if `ModificationProtectionStatus` is set to `ConsoleProtection`.
    * `resource_group_id` - The ID of the resource group.
    * `status` - The The load balancer status. Valid values: `Active`, `Configuring`, `CreateFailed`, `Inactive` and `Provisioning`.
    * `tags` - The tag of the resource. 
        * `tag_key` - The key of the tags. 
        * `tag_value` - The value of the tags.
    * `vpc_id` - The ID of the virtual private cloud (VPC) where the ALB instance is deployed. 
    * `zone_mappings` - The zones and vSwitches. You must specify at least two zones.
       * `vswitch_id` - The ID of the vSwitch that corresponds to the zone. Each zone can use only one vSwitch and subnet.
       * `zone_id` - The ID of the zone to which the ALB instance belongs.
