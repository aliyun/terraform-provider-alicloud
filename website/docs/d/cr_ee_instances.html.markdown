---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_instances"
sidebar_current: "docs-alicloud-datasource-cr-ee-instances"
description: |-
  Provides a list of Container Registry Enterprise Edition instances.
---

# alicloud\_cr\_ee\_instances

This data source provides a list Container Registry Enterprise Edition instances on Alibaba Cloud.

-> **NOTE:** Available in v1.86.0+

## Example Usage

```
# Declare the data source
data "alicloud_cr_ee_instances" "my_instances" {
  name_regex  = "my-instances"
  output_file = "my-instances-json"
}

output "output" {
  value = "${data.alicloud_cr_ee_instances.my_instances.instances}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ids to filter results by instance id.
* `name_regex` - (Optional) A regex string to filter results by instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional, Available in 1.132.0+) Default to `true`. Set it to true can output instance authorization token.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Container Registry Enterprise Edition instances. Its element is an instance uuid.
* `names` - A list of instance names.
* `instances` - A list of matched Container Registry Enterprise Editioninstances. Each element contains the following attributes:
  * `id` - ID of Container Registry Enterprise Edition instance.
  * `name` - Name of Container Registry Enterprise Edition instance.
  * `region` - Region of Container Registry Enterprise Edition instance.
  * `specification` - Specification of Container Registry Enterprise Edition instance.
  * `namespace_quota` - The max number of namespaces that an instance can create.
  * `namespace_usage` - The number of namespaces already created.
  * `repo_quota` - The max number of repos that an instance can create.
  * `repo_usage` - The number of repos already created.
  * `vpc_endpoints` - A list of domains for access on vpc network.
  * `public_endpoints` - A list of domains for access on internet network.
  * `authorization_token` - The password that was used to log on to the registry.
  * `temp_username` - The username that was used to log on to the registry.
  
