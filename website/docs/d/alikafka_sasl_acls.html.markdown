---
subcategory: "Alikafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_sasl_acls"
sidebar_current: "docs-alicloud-datasource-alikafka-sasl-acls"
description: |-
    Provides a list of alikafka sasl acls available to the user.
---

# alicloud\_alikafka\_sasl\_acls

This data source provides a list of ALIKAFKA Sasl acls in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.66.0+

## Example Usage

```
data "alicloud_alikafka_sasl_acls" "sasl_acls_ds" {
  instance_id = "xxx"
  username = "username"
  acl_resource_type = "Topic"
  acl_resource_name = "testTopic"
  output_file = "saslAcls.txt"
}

output "first_sasl_acl_username" {
  value = "${data.alicloud_alikafka_sasl_acls.sasl_acls_ds.acls.0.username}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ALIKAFKA Instance that owns the sasl acls.
* `username` - (Required) Get results for the specified username. 
* `acl_resource_type` - (Required) Get results for the specified resource type. 
* `acl_resource_name` - (Required) Get results for the specified resource name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `acls` - A list of sasl acls. Each element contains the following attributes:
  * `username` - The username of the sasl acl.
  * `acl_resource_type` - The resource type of the sasl acl.
  * `acl_resource_name` - The resource name of the sasl acl.
  * `acl_resource_pattern_type` - The resource pattern type of the sasl acl.
  * `host` - The host of the sasl acl.
  * `acl_operation_type` - The operation type of the sasl acl.
