---
subcategory: "Open Api Explorer"
layout: "alicloud"
page_title: "Alicloud: alicloud_open_api_explorer_api_mcp_server"
description: |-
  Provides a Alicloud Open Api Explorer Api Mcp Server resource.
---

# alicloud_open_api_explorer_api_mcp_server

Provides a Open Api Explorer Api Mcp Server resource.

API MCP Server.

For information about Open Api Explorer Api Mcp Server and how to use it, see [What is Api Mcp Server](https://next.api.alibabacloud.com/document/OpenAPIExplorer/2024-11-30/CreateApiMcpServer).

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_open_api_explorer_api_mcp_server" "default" {
  system_tools = ["FetchRamActionDetails"]
  description  = "Create"
  prompts {
    description = "Obtain user customization information description"
    content     = "Prompt body,{{name}}"
    arguments {
      description = "Name information"
      required    = true
      name        = "name"
    }
    name = "Obtain user customization information"
  }
  prompts {
    description = "Obtain user customization information description"
    content     = "Prompt text, {{name }}, {{age }}, {{description}}"
    arguments {
      description = "Name information"
      required    = true
      name        = "name"
    }
    arguments {
      description = "Age information"
      required    = true
      name        = "age"
    }
    arguments {
      description = "Description Information"
      required    = true
      name        = "description"
    }
    name = "Obtain user customization information 1"
  }
  oauth_client_id = "123456"
  apis {
    api_version = "2014-05-26"
    product     = "Ecs"
    selectors   = ["DescribeAvailableResource", "DescribeRegions", "DescribeZones"]
  }
  apis {
    api_version = "2017-03-21"
    product     = "vod"
    selectors   = ["CreateUploadVideo"]
  }
  apis {
    api_version = "2014-05-15"
    product     = "Slb"
    selectors   = ["DescribeAvailableResource"]
  }
  instructions = "Describes the role of the entire MCP Server"
  additional_api_descriptions {
    api_version          = "2014-05-26"
    enable_output_schema = true
    api_name             = "DescribeAvailableResource"
    const_parameters {
      value = "cn-hangzhou"
      key   = "x_mcp_region_id"
    }
    const_parameters {
      value = "B1"
      key   = "a1"
    }
    const_parameters {
      value = "b2"
      key   = "a2"
    }
    api_override_json   = jsonencode({ "summary" : "This operation supports querying the list of instances based on different request conditions and associating the query instance details. " })
    product             = "Ecs"
    execute_cli_command = false
  }
  additional_api_descriptions {
    api_version          = "2014-05-26"
    enable_output_schema = true
    api_name             = "DescribeRegions"
    product              = "Ecs"
    execute_cli_command  = true
  }
  additional_api_descriptions {
    api_version          = "2014-05-26"
    enable_output_schema = true
    api_name             = "DescribeZones"
    product              = "Ecs"
    execute_cli_command  = true
  }
  vpc_whitelists           = ["vpc-examplea", "vpc-exampleb", "vpc-examplec"]
  name                     = "my-name"
  language                 = "ZH_CN"
  enable_assume_role       = true
  assume_role_extra_policy = jsonencode({ "Version" : "1", "Statement" : [{ "Effect" : "Allow", "Action" : ["ecs:Describe*", "vpc:Describe*", "vpc:List*"], "Resource" : "*" }] })
  terraform_tools {
    description    = "Terraform as tool example"
    async          = true
    destroy_policy = "NEVER"
    code           = <<EOF
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "default" {
  ipv6_isp    = "BGP"
  description = "example"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name
  enable_ipv6 = true
}
  EOF
    name           = "tfexample"
  }
  terraform_tools {
    description    = "Terraform as tool example"
    async          = true
    destroy_policy = "NEVER"
    code           = <<EOF
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "default" {
  ipv6_isp    = "BGP"
  description = "example"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name
  enable_ipv6 = true
}
  EOF
    name           = "tfexample2"
  }
  terraform_tools {
    description    = "Terraform as tool example"
    async          = true
    destroy_policy = "NEVER"
    code           = <<EOF
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "default" {
  ipv6_isp    = "BGP"
  description = "example"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name
  enable_ipv6 = true
}
  EOF
    name           = "tfexample3"
  }
  assume_role_name            = "default-role"
  public_access               = "on"
  enable_custom_vpc_whitelist = true
}
```

## Argument Reference

The following arguments are supported:
* `additional_api_descriptions` - (Optional, Set) Additional OpenAPI description information that can override the default behavior of APIs, including:
  - API name  
  - Modification or removal of API parameter names  
  - Whether to exclude the API from the output API response structure definition  
  - Whether to return a CLI execution command instead of directly executing the API  
  - Configuration of constant values for API parameters; parameters set as constants will not have their definitions returned in the tool list   See [`additional_api_descriptions`](#additional_api_descriptions) below.
* `apis` - (Required, Set) The list of APIs to be included in the API MCP Server. See [`apis`](#apis) below.
* `assume_role_extra_policy` - (Optional, JsonString) When multi-account access is enabled, this field defines an additional policy for role assumption. If specified, this policy overrides the original permissions defined for the role, and the assumed role’s permissions are determined solely by this policy.
* `assume_role_name` - (Optional) The name of the RAM role in the target account to assume when enabling multi-account access for cross-account operations.
* `description` - (Optional) Description of the API MCP service.
* `enable_assume_role` - (Optional) Specifies whether to enable multi-account access. When enabled, the MCP Server exposes the x_assume_account_id parameter by default. When this parameter is provided, the MCP Server switches to the specified account to perform operations.
* `enable_custom_vpc_whitelist` - (Optional) Whether to enable a custom VPC whitelist. If disabled, the configuration follows the account-level setting.
* `instructions` - (Optional) MCP instructions that guide the large language model on how to use this MCP. The client must support the Instructions field defined in the MCP standard protocol.  
* `language` - (Optional) Documentation language for the API MCP service. You can select either Chinese or English API documentation. The choice of language may affect the AI's response quality due to differences in prompt wording. Supported values are EN_US and ZH_CN.
* `name` - (Required, ForceNew) Name of the MCP Server. It can contain digits, English letters, and hyphens (-).
* `oauth_client_id` - (Optional) The custom OAuth Client ID when selecting a custom OAuth configuration.
`Supported only for Web/Native applications, and the OAuth scope must include /acs/mcp-server.`
* `prompts` - (Optional, Set) List of prompts supported by the MCP Server. For the MCP protocol, clients retrieve this list through the prompts/list RPC call. See [`prompts`](#prompts) below.
* `public_access` - (Optional) Whether to enable public network access. This setting takes precedence over the account-level configuration and supports the following options:
  - on: enables public network access;
  - off: disables public network access;
  - follow: inherits the account-level configuration.
* `system_tools` - (Optional, Set) Enabled system services.
* `terraform_tools` - (Optional, Set) A list of Terraform Tools. The MCP Server allows using Terraform HCL code as a complete tool to improve the determinism of orchestration. See [`terraform_tools`](#terraform_tools) below.
* `vpc_whitelists` - (Optional, Set) When public network access is disabled, this field specifies the VPC whitelist that restricts source VPCs. If not set or left empty, no restriction is applied to the source.

### `additional_api_descriptions`

The additional_api_descriptions supports the following:
* `api_name` - (Optional) The API name, such as ListApiMcpServers.  
* `api_override_json` - (Optional, JsonString) API structure definition information. You can use this parameter to directly modify the API description and parameter list. You can obtain the API definition information from an API endpoint such as https://api.aliyun.com/meta/v1/products/Ecs/versions/2014-05-26/apis/DescribeInstances/api.json.  

-> **NOTE:** Note that required parameters must not be removed; otherwise, calls by the large model will continuously fail due to missing required parameters.>  

* `api_version` - (Optional) API version information, typically in date format, such as 2014-05-26.  
* `const_parameters` - (Optional, Set) Constant configuration information. When the MCP Server needs to fix certain tool parameters to specific values, you can configure this parameter to enforce those fixed values.  
Parameters configured as constants will not be returned as tool parameters through the MCP protocol. Large models cannot define these parameters. During execution, the MCP Server merges these constant values into the API call parameters.   See [`const_parameters`](#additional_api_descriptions-const_parameters) below.
* `enable_output_schema` - (Optional) By default, this feature is disabled, and the MCP Server returns only the structure definition of input parameters. When enabled, the MCP Server returns the output parameter structure definition via the MCP protocol.  

-> **NOTE:** The output parameter structure may be complex. Enabling this feature significantly increases the MCP context size. Use this feature with caution.>  

* `execute_cli_command` - (Optional) Call interception. When this parameter is enabled, the MCP Server returns the complete CLI command name instead of directly executing the API call. Use this option when the API call is long-running or requires interaction with local files. The MCP Server enforces theoretical time limits for single-tool invocations:  
  - SSE protocol: up to 30 minutes  
  - Streamable HTTP protocol: up to 1 minute  

For tools whose single API execution exceeds 30 minutes, we recommend enabling this parameter. Install the CLI and complete account authentication on the machine initiating the call, then combine it with this tool for optimal results.  

-> **NOTE:** The identity used to execute the CLI differs from the identity used by the MCP Server. Pay attention to the associated security risks.>  

* `product` - (Optional) The name of the cloud product, such as Ecs.  

### `additional_api_descriptions-const_parameters`

The additional_api_descriptions-const_parameters supports the following:
* `key` - (Optional) Parameter location. Currently, except for ROA-style body parameters (which support up to two levels), nested parameter configurations beyond two levels are not supported. If you need to configure a composite data structure, set the Value to a JSON object.  

For RPC-style APIs, examples include:  
  - Name: sets the Name parameter to a fixed value.  

For ROA-style APIs, examples include:  
  - Name: sets a query or path parameter named Name to a fixed value;  
  - body.Name: sets the Name field within the request body to a fixed value.  

Configurations such as body.Name.Sub are not supported. If you need to set body.Name as a composite structure, specify the Value as a JSON object—for example, {"Sub": "xxx"}.  

-> **NOTE:** x_mcp_region_id is a built-in MCP parameter used to control the region and can also be configured as a fixed value to invoke services in a specified region.>  

* `value` - (Optional) This property does not have a description in the spec, please add it before generating code.

### `apis`

The apis supports the following:
* `api_version` - (Required) API version information, typically in date format—for example, the version for ECS is 2014-05-26.
* `product` - (Required) Product code, such as Ecs.
* `selectors` - (Required, Set) Selectors in array format, where each item is an API name—for example, GetApiDefinition or ListApiDefinitions. You can obtain the complete list of supported APIs from the Alibaba Cloud Developer Portal.

### `prompts`

The prompts supports the following:
* `arguments` - (Optional, Set) Parameters for the prompt. See [`arguments`](#prompts-arguments) below.
* `content` - (Optional) Full content of the prompt, supporting dynamic parameters. Parameters must be defined in Arguments, using the format {{ARG}}, where ARG supports English characters. Example: My name is: {{name}}.
* `description` - (Optional) Description of the prompt.
* `name` - (Optional) Name of the prompt.

### `prompts-arguments`

The prompts-arguments supports the following:
* `description` - (Optional) Description of the prompt parameter.
* `name` - (Optional) Name of the prompt parameter.
* `required` - (Optional) Indicates whether the prompt parameter is required.

### `terraform_tools`

The terraform_tools supports the following:
* `async` - (Optional) Specifies whether execution is asynchronous. If enabled, the system immediately proceeds to the next task after initiating a task, without waiting for each resource operation to complete.
* `code` - (Optional) Terraform Tool code. [Overview of the HCL Language](https://help.aliyun.com/zh/terraform/terraform-configuration-and-hcl-language-overview)
* `description` - (Optional) Description of the Terraform Tool. This description will be used as the description for the MCP tool.
* `destroy_policy` - (Optional) The cleanup policy applied to temporary resources after task completion, based on the task execution status:
  - NEVER: Do not delete any created resources, regardless of whether the task succeeds or fails.
  - ALWAYS: Immediately destroy all related resources upon task completion, regardless of success or failure.
  - ON_FAILURE: Delete related resources only if the task fails; retain them if the task succeeds.
* `name` - (Optional) The name of the Terraform Tool, which supports letters (a–z, A–Z) and digits (0–9).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - MCP Server creation time in China Standard Time (CST), for example, 2025-12-04 19:46:52.  

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Api Mcp Server.
* `delete` - (Defaults to 5 mins) Used when delete the Api Mcp Server.
* `update` - (Defaults to 5 mins) Used when update the Api Mcp Server.

## Import

Open Api Explorer Api Mcp Server can be imported using the id, e.g.

```shell
$ terraform import alicloud_open_api_explorer_api_mcp_server.example <id>
```