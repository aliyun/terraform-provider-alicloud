---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_service"
sidebar_current: "docs-alicloud-resource-fc"
description: |-
  Provides a Alicloud Function Compute Service resource. The resource is the base of launching Function and Trigger configuration.
---

# alicloud\_fc\_service

Provides a Alicloud Function Compute Service resource. The resource is the base of launching Function and Trigger configuration.
 For information about Service and how to use it, see [What is Function Compute](https://www.alibabacloud.com/help/doc-detail/52895.htm).

-> **NOTE:** The resource requires a provider field 'account_id'. [See account_id](https://www.terraform.io/docs/providers/alicloud/index.html#account_id).

-> **NOTE:** If you happen the error "Argument 'internetAccess' is not supported", you need to log on web console and click button "Apply VPC Function"
which is in the upper of [Function Service Web Console](https://fc.console.aliyun.com/) page.

-> **NOTE:** Currently not all regions support Function Compute Service.
For more details supported regions, see [Service endpoints](https://www.alibabacloud.com/help/doc-detail/52984.htm)

## Example Usage

Basic Usage

```
variable "name" {
  default = "tf-testaccalicloudfcservice"
}

resource "alicloud_log_project" "foo" {
  name = var.name
}

resource "alicloud_log_store" "foo" {
  project = alicloud_log_project.foo.name
  name    = var.name
}

resource "alicloud_ram_role" "role" {
  name     = var.name
  document = <<DEFINITION
  {
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "fc.aliyuncs.com"
        ]
      }
    }
  ],
  "Version": "1"
}
  DEFINITION
  description = "this is a test"
  force = true
}

resource "alicloud_ram_role_policy_attachment" "attac" {
  role_name = alicloud_ram_role.role.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

resource "alicloud_fc_service" "foo" {
  name = var.name
  description = "tf unit test"
  role = alicloud_ram_role.role.arn
  depends_on = ["alicloud_ram_role_policy_attachment.attac"]
}
```

## Module Support

You can use to the existing [fc module](https://registry.terraform.io/modules/terraform-alicloud-modules/fc/alicloud) 
to create a service and a function quickly and then set several triggers for it.

## Argument Reference

The following arguments are supported:

* `name` - (ForceNew) The Function Compute service name. It is the only in one Alicloud account and is conflict with "name_prefix".
* `name_prefix` - (ForceNew) Setting a prefix to get a only name. It is conflict with "name".
* `description` - (Optional) The function compute service description.
* `internet_access` - (Optional) Whether to allow the service to access Internet. Default to "true".
* `role` - (Optional) RAM role arn attached to the Function Compute service. This governs both who / what can invoke your Function, as well as what resources our Function has access to. See [User Permissions](https://www.alibabacloud.com/help/doc-detail/52885.htm) for more details.
* `log_config` - (Optional) Provide this to store your FC service logs. Fields documented below. See [Create a Service](https://www.alibabacloud.com/help/doc-detail/51924.htm).
* `vpc_config` - (Optional) Provide this to allow your FC service to access your VPC. Fields documented below. See [Function Compute Service in VPC](https://www.alibabacloud.com/help/faq-detail/72959.htm).

**log_config** requires the following:

* `project` - (Required) The project name of Logs service.
* `logstore` - (Required) The log store name of Logs service.

-> **NOTE:** If both `project` and `logstore` are empty, log_config is considered to be empty or unset.

**vpc_config** requires the following:

* `vswitch_ids` - (Required) A list of vswitch IDs associated with the FC service.
* `security_group_id` - (Required) A security group ID associated with the FC service.

-> **NOTE:** If both `vswitch_ids` and `security_group_id` are empty, vpc_config is considered to be empty or unset.

## Attributes Reference

The following arguments are exported:

* `id` - The ID of the FC service. The value is same as name.
* `service_id` - The Function Compute service ID.
* `last_modified` - The date this resource was last modified.

## Import

Function Compute Service can be imported using the id or name, e.g.

```
$ terraform import alicloud_fc_service.foo my-fc-service
```
