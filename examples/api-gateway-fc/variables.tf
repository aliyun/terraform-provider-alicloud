variable "apigateway_group_name" {
  default = "TerraformApiGatewayGroupFc"
}

variable "apigateway_group_description" {
  default = "api gateway description"
}

variable "apigateway_app_name_test" {
  default = "tf_testAccAppResource"
}

variable "apigateway_app_description_test" {
  default = "tf_testAcc api gateway description"
}

variable "fc_region" {
  description = "The region to launch resources."
  default     = "cn-beijing"
}

variable "service_name" {
  description = "The Function Compute service name."
  default     = "my-fc-service"
}

variable "service_description" {
  description = "The Function Compute service description."
  default     = "created by terraform"
}

variable "service_internet_access" {
  description = "Whether to allow the service to access Internet. Default to true."
  default     = false
}

variable "function_name" {
  description = "The Function Compute function name."
  default     = "hello-world"
}

variable "function_description" {
  description = "The Function Compute function description."
  default     = "created by terraform"
}

variable "function_filename" {
  description = "The path to the function's deployment package within the local filesystem. It is conflict with the oss_-prefixed options.."
  default     = "./index.zip"
}

variable "function_memory_size" {
  description = "Amount of memory in MB your Function can use at runtime. Defaults to 128. Limits to [128, 32768]."
  default     = "512"
}

variable "function_runtime" {
  description = "The Function Compute function runtime type."
  default     = "nodejs6"
}

variable "function_handler" {
  description = "The function entry point in your code."
  default     = "index.handler"
}

variable "trigger_name" {
  description = "The Function Compute trigger name.."
  default     = "trigger-for-fc"
}

variable "ram_role_name" {
  description = "The Function Compute trigger name.."
  default     = "AliyunApiGatewayAccessingFCRoleTest"
}

variable "ram_policy_name" {
  description = "The Function Compute trigger name.."
  default     = "AliyunApiGatewayAccessingFCRolePolicyTest"
}

