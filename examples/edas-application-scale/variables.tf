variable "app_id" {
  description = "The ID of the application that you want to deploy."
}

variable "deploy_group" {
  description = "The ID of the instance group to which you want to add ECS instances to scale out the application."
}

variable "ecu_info" {
  type        = list(string)
  description = "The ID of the Elastic Compute Unit (ECU) where you want to deploy the application."
}

variable "force_status" {
  type        = bool
  description = "This parameter specifies whether to forcibly remove an ECS instance where the application is deployed. It is set as true only after the ECS instance expires. In normal cases, this parameter do not need to be specified."
}
