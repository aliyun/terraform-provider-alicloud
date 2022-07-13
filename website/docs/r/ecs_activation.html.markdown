---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_activation"
sidebar_current: "docs-alicloud-resource-ecs-activation"
description: |-	
	 Provides a Alicloud ECS Activation resource.
---

# alicloud\_ecs\_activation

Provides a ECS Activation resource.

For information about ECS Activation and how to use it, see [What is Activation](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/createactivation#doc-api-Ecs-CreateActivation).

-> **NOTE:** Available in v1.177.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_activation" "example" {
  description           = var.name
  instance_count        = 10
  instance_name         = var.name
  ip_address_range      = "0.0.0.0/0"
  time_to_live_in_hours = 4
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of the activation code. The description can be 1 to 100 characters in length and cannot start with `http://` or `https://`.
* `instance_count` - (Optional, ForceNew, Computed) The maximum number of times that the activation code can be used to register managed instances. Valid values: `1` to `1000`. Default value: `10`.
* `instance_name` - (Optional, ForceNew) The default instance name prefix. The instance name prefix must be 1 to 50 characters in length. It must start with a letter and cannot start with `http://` or `https://`. The instance name prefix can contain only letters, digits, periods (.), underscores (_), hyphens (-), and colons (:).
		- If you use the activation code created by the CreateActivation operation to register managed instances, the instances are assigned sequential names that are prefixed by the value of this parameter. You can also specify a new instance name to override the assigned sequential name when you register a managed instance.
		- If you specify InstanceName when you register a managed instance, an instance name in the format of `<InstanceName>-<Number>` is generated. The number of digits in the <Number> value is determined by that in the InstanceCount value. Example: 001. If you do not specify InstanceName, the hostname (Hostname) is used as the instance name.
* `ip_address_range` - (Optional, ForceNew, Computed) The IP addresses of hosts that are allowed to use the activation code. The value can be IPv4 addresses, IPv6 addresses, or CIDR blocks.
* `time_to_live_in_hours` - (Optional, ForceNew, Computed) The validity period of the activation code. The activation code cannot be used to register new instances after the validity period expires. Unit: hours. Valid values: `1` to `24`. Default value: `4`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Activation.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Activation.
* `delete` - (Defaults to 1 mins) Used when delete the Activation.

## Import

ECS Activation can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_activation.example <id>
```