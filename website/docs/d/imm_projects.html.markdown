---
subcategory: "Intelligent Media Management (IMM)"
layout: "alicloud"
page_title: "Alicloud: alicloud_imm_projects"
sidebar_current: "docs-alicloud-datasource-imm-projects"
description: |-
  Provides a list of Intelligent Media Management Projects to the user.
---

# alicloud\_imm\_projects

This data source provides the Intelligent Media Management Projects of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_imm_projects" "ids" {
  ids = ["example_id"]
}
output "imm_project_id_1" {
  value = data.alicloud_imm_projects.ids.projects.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Project IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `projects` - A list of Imm Projects. Each element contains the following attributes:

    * `create_time` - The creation time of project.
    * `id` - The ID of project.
    * `modify_time` - The modification time of project.
    * `type` - The type of project.
    * `project` -The name of project.
    * `service_role` - The service role authorized to the Intelligent Media Management service to access other cloud resources.
    * `endpoint` - The service address of project.
    * `billing_type` - The billing type. **Note:** This parameter is deprecated from 2021-04-01.
    * `compute_unit` - The maximum number of requests that can be processed per second. **Note:** This parameter is deprecated from 2021-04-01.
