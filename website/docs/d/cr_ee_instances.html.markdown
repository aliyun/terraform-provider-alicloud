---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_instances"
sidebar_current: "docs-alicloud-datasource-cr-ee-instances"
description: |-
  Provides a list of Container Registry Enterprise instances.
---

# alicloud\_cr_ee\_instances

This data source provides a list Container Registry Enterprise instances on Alibaba Cloud.

-> **NOTE:** Available in v1.85.0+

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

* `name_regex` - (Optional) A regex string to filter results by instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched CR EE instances. Its element is an instance uuid.
* `names` - A list of instance names.
* `instances` - A list of matched CR EE instances. Each element contains the following attributes:
  * `id` - ID of CR EE instance.
  * `name` - Name of CR EE instance.
  * `region` - Region of CR EE instance.
  * `specification` - Specification of CR EE instance.
  * `namespace_quota` - The max number of namespaces that an instance can create.
  * `namespace_usage` - The number of namespaces already created.
  * `repo_quota` - The max number of repos that an instance can create.
  * `repo_usage` - The number of repos already created.
  * `vpc_endpoints` - A list of domains for access on vpc network.
  * `public_endpoints` - A list of domains for access on internet network.