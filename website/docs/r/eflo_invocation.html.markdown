---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_invocation"
description: |-
  Provides a Alicloud Eflo Invocation resource.
---

# alicloud_eflo_invocation

Provides a Eflo Invocation resource.

Cloud assistant command execution on the node.

For information about Eflo Invocation and how to use it, see [What is Invocation](https://next.api.alibabacloud.com/document/eflo-controller/2022-12-15/RunCommand).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
# Before executing this example, you need to confirm with the product team whether the resources are sufficient or you will get an error message with "Failure to check order before create instance"
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "create_vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = "cluster-resoure-example"
}

resource "alicloud_vswitch" "create_vswitch" {
  vpc_id       = alicloud_vpc.create_vpc.id
  zone_id      = "cn-hangzhou-b"
  cidr_block   = "192.168.0.0/24"
  vswitch_name = "cluster-resoure-example"
}

resource "alicloud_security_group" "create_security_group" {
  description         = "sg"
  security_group_name = "cluster-resoure-example"
  security_group_type = "normal"
  vpc_id              = alicloud_vpc.create_vpc.id
}

resource "alicloud_eflo_cluster" "default" {
  cluster_description  = "cluster-resource-example"
  open_eni_jumbo_frame = "false"
  hpn_zone             = "B1"
  nimiz_vswitches = [
    "1111"
  ]
  ignore_failed_node_tasks = "true"
  resource_group_id        = data.alicloud_resource_manager_resource_groups.default.ids.1
  node_groups {
    image_id               = "i198448731735114628708"
    zone_id                = "cn-hangzhou-b"
    node_group_name        = "cluster-resource-example"
    node_group_description = "cluster-resource-example"
    machine_type           = "efg2.C48cA3sen"
  }

  networks {
    tail_ip_version = "ipv4"
    new_vpd_info {
      monitor_vpc_id     = alicloud_vpc.create_vpc.id
      monitor_vswitch_id = alicloud_vswitch.create_vswitch.id
      cen_id             = "11111"
      cloud_link_id      = "1111"
      vpd_cidr           = "111"
      vpd_subnets {
        zone_id     = "1111"
        subnet_cidr = "111"
        subnet_type = "111"
      }
      cloud_link_cidr = "169.254.128.0/23"
    }

    security_group_id = alicloud_security_group.create_security_group.id
    vswitch_zone_id   = "cn-hangzhou-b"
    vpc_id            = alicloud_vpc.create_vpc.id
    vswitch_id        = alicloud_vswitch.create_vswitch.id
    vpd_info {
      vpd_id = "111"
      vpd_subnets = [
        "111"
      ]
    }
    ip_allocation_policy {
      bond_policy {
        bond_default_subnet = "111"
        bonds {
          name   = "111"
          subnet = "111"
        }
      }
      machine_type_policy {
        bonds {
          name   = "111"
          subnet = "111"
        }
        machine_type = "111"
      }
      node_policy {
        bonds {
          name   = "111"
          subnet = "111"
        }
        node_id = "111"
      }
    }
  }

  cluster_name = "tfacceflo7165"
  cluster_type = "Lite"
}

resource "alicloud_eflo_node" "default" {
  period           = "36"
  discount_level   = "36"
  billing_cycle    = "1month"
  classify         = "gpuserver"
  zone             = "cn-hangzhou-b"
  product_form     = "instance"
  payment_ratio    = "0"
  hpn_zone         = "B1"
  server_arch      = "bmserver"
  computing_server = "efg1.nvga1n"
  stage_num        = "36"
  renewal_status   = "AutoRenewal"
  renew_period     = "36"
  status           = "Unused"
}

resource "alicloud_eflo_node_group" "default" {
  nodes {
    node_id        = alicloud_eflo_node.default.id
    vpc_id         = alicloud_vpc.create_vpc.id
    vswitch_id     = alicloud_vswitch.create_vswitch.id
    hostname       = "jxyhostname"
    login_password = "Alibaba@2025"
  }

  ignore_failed_node_tasks = "true"
  cluster_id               = alicloud_eflo_cluster.default.id
  image_id                 = "i195048661660874208657"
  zone_id                  = "cn-hangzhou-b"
  vpd_subnets = [
    "example"
  ]
  user_data       = "YWxpLGFsaSxhbGliYWJh"
  vswitch_zone_id = "cn-hangzhou-b"
  ip_allocation_policy {
    bond_policy {
      bond_default_subnet = "example"
      bonds {
        name   = "example"
        subnet = "example"
      }
    }
    machine_type_policy {
      bonds {
        name   = "example"
        subnet = "example"
      }
      machine_type = "example"
    }
    node_policy {
      node_id = alicloud_eflo_node.default.id
      bonds {
        name   = "example"
        subnet = "example"
      }
    }
  }
  machine_type           = "efg1.nvga1"
  az                     = "cn-hangzhou-b"
  node_group_description = "resource-example1"
  node_group_name        = "tfacceflo63657_update"
}

resource "alicloud_eflo_invocation" "default" {
  description      = "example"
  content_encoding = "Base64"
  name             = "resource-example"
  repeat_mode      = "Once"
  parameters {
    name = "example"
  }
  node_id_list     = [alicloud_eflo_node.default.id]
  timeout          = "68"
  command_content  = "ZWNobyAxMjM="
  working_dir      = "/home/"
  username         = "root"
  enable_parameter = false
  termination_mode = "ProcessTree"
}
```

### Deleting `alicloud_eflo_invocation` or removing it from your configuration

Terraform cannot destroy resource `alicloud_eflo_invocation`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `command_content` - (Optional) The command content. You need to pay attention:
  - Specify the parameter 'EnableParameter = true' to enable the custom parameter feature in the command content.
  - Define custom parameters in the form of {{}} inclusion, and spaces and line breaks before and after the parameter name in '{{}}' are ignored.
  - The number of custom parameters cannot exceed 20.
  - Custom parameter names can a-zA-Z0-9 a combination of-_. Other characters are not supported. Parameter names are not case-sensitive.
  - A single custom parameter name cannot exceed 64 bytes.
* `command_id` - (Optional) Command ID
* `content_encoding` - (Optional) The encoding of the script content. Value range:
  - PlainText: no encoding, using PlainText transmission.
  - Base64:Base64 encoding.

Default value: PlainText. If you fill it randomly or wrongly, the value will be treated as a PlainText.
* `description` - (Optional) The command description.
* `enable_parameter` - (Optional) Whether custom parameters are included in the command.
Default value: false.
* `frequency` - (Optional) The execution time of the scheduled execution command. Currently, three scheduled execution methods are supported: fixed interval execution (based on Rate expression), only once at a specified time, and timed execution based on clock (based on Cron expression).
  - Fixed time interval execution: Based on the Rate expression, the command is executed at the set time interval. Time intervals can be selected by seconds (s), minutes (m), hours (h), and days (d), which is suitable for scenarios where tasks are executed at fixed time intervals. The format is rate( ). If the execution is performed every 5 minutes, the format is rate(5m). Executing with a fixed time interval has the following limitations:
  - The set time interval is no more than 7 days and no less than 60 seconds, and must be greater than the timeout period of the scheduled task.
  - The execution interval is based only on a fixed frequency, independent of the time the task actually takes to execute. For example, if the command is executed every 5 minutes and the task takes 2 minutes to complete, the next round will be executed 3 minutes after the task is completed.
  - The task is not executed immediately when it is created. For example, if a command is executed every 5 minutes, the command is not executed immediately when a task is created, but is executed 5 minutes after the task is created.
  - Execute only once at the specified time: Execute the command once according to the set time zone and execution time point. The format is at(yyyy-MM-dd HH:mm:ss ), that is, at (year-month-day time: minute: Second ). If you do not specify a time zone, the default is the UTC time zone. Time zones can be in the following three formats: the full name of the time zone, such as Asia/Shanghai (China/Shanghai time), America/los_angles (United States/Los Angeles time), and so on. The offset of the time zone relative to Greenwich Mean Time: E.G. GMT +8:00 (East Zone 8), GMT-7 (West Zone 7), etc. When using the GMT format, the hour bit does not support adding leading zeros. Time zone abbreviation: Only UTC (Coordinated Universal Time) is supported.
If it is specified to be executed once 13:15:30 June 06, 2022, China/Shanghai time, the format is at (Asia/Shanghai, 2022-06-06 13:15:30); If it is specified to be executed once 13:15:30 June 06, 2022, the format is at(2022-06-06 13:15:30 GMT-7:00).
  - Timing based on clock (based on Cron expression): Based on Cron expression, commands are executed according to the set timing task. The format is        , that is,  . In the specified time zone, calculate the execution time of the scheduled task based on the Cron expression and execute it. If no time zone is specified, the default time zone is the internal time zone of the scheduled task instance. For more information about Cron expressions, see Cron Expressions. Time zones support the following three forms:
  - Full time zone name: such as Asia/Shanghai (China/Shanghai time), America/los_angles (US/Los Angeles time), etc.
  - The offset of the time zone relative to Greenwich Mean Time: E.G. GMT +8:00 (East Zone 8), GMT-7 (West Zone 7), etc. When using the GMT format, the hour bit does not support adding leading zeros.
  - Time zone abbreviation: Only UTC (Coordinated Universal Time) is supported.

For example, in China/Shanghai time, the command will be executed once every day at 10:15 am in 2022 in the format 0 15 10? * * 2022 Asia/Shanghai; In the eastern 8th District time, it will be executed every half hour from 10:00 a.m. to 11:30 a.m. every day in 2022, in the format of 0 0/30 10-11 * *? 2022 GMT +8:00; In UTC time, starting from 2022, it will be executed every 5 minutes from 14:00 P.M. to 14:55 p. M. Every two years in October, in the format of 0 0/5 14*10? 2022/2 UTC.
* `launcher` - (Optional) The bootstrapper for script execution. The length cannot exceed 1KB.
* `name` - (Optional) The command name.
* `node_id_list` - (Optional, ForceNew, List) A list of nodes.
* `parameters` - (Optional, Map) When the command contains custom parameters, the key-value pair of the custom parameters passed in when the command is executed. For example, if the command content is 'echo {{name}}', the key-value pair'{"name":"Jack"}'can be passed through the 'Parameter' parameter'. The custom parameter will automatically replace the variable value 'name' to get a new command that actually executes 'echo Jack '.

The number of custom parameters ranges from 0 to 10, and you need to pay attention:
  - The key is not allowed to be an empty string and supports a maximum of 64 characters.
  - The value is allowed to be an empty string.
  - After the custom parameters and the original command content are encoded in Base64, if the command is saved, the size of the command content after Base64 encoding cannot exceed 18KB. If the command is not saved, the size of the command content after Base64 encoding cannot exceed 24KB. You can set whether to keep the command through 'KeepCommand.
  - The set of custom parameter names must be a subset of the parameter set defined when the command is created. For parameters that are not passed in, you can use an empty string instead.

The default value is empty, which means that the parameter is unset and the custom parameter is disabled.
* `repeat_mode` - (Optional) Sets the way the command is executed. Value range:
  - Once: Execute the command immediately.
  - Period: executes the command regularly. When the value of this parameter is 'Period', the 'Frequency' parameter must also be specified.
  - NextRebootOnly: Automatically execute the command when the instance is next started.
  - EveryReboot: The command is automatically executed every time the instance is started.

Default:
  - When the'frequency' parameter is not specified, the default value is'once '.
  - When the'frequency' parameter is specified, regardless of whether the parameter value has been set or not, it will be processed according to'period.
* `termination_mode` - (Optional) The mode when the task is stopped (manually stopped or execution time-out interrupted). Possible values:
Process: Stops the current script Process.
ProcessTree: Stops the current process tree (the script process and the collection of all child processes it created)
* `timeout` - (Optional, Int) The timeout period for command execution. Unit: seconds. A timeout occurs when a command cannot be run due to a process, a missing module, or a missing cloud assistant Agent. After the timeout, the command process is forcibly terminated. Default value: 60.
* `username` - (Optional) The name of the user who executed the command in the instance. The length must not exceed 255 characters.
The instance of the Linux system. By default, the root user runs commands.
* `working_dir` - (Optional) You can customize the command execution path. The default path is as follows:
Linux instance: the execution path is in the/home directory of the root user by default.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Invocation.
* `update` - (Defaults to 5 mins) Used when update the Invocation.

## Import

Eflo Invocation can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_invocation.example <id>
```