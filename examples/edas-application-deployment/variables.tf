variable "app_id" {
  description = "The ID of the application that you want to deploy."
}

variable "group_id" {
  description = "The ID of the instance group where the application is going to be deployed. Set this parameter to all if you want to deploy the application to all groups."
}


variable "package_version" {
  description = "The version of the application that you want to deploy. It must be unique for every application. The length cannot exceed 64 characters. We recommended you to use a timestamp. "
}

variable "war_url" {
  description = "The address to store the uploaded web application (WAR) package for application deployment. This parameter is required when the deployType parameter is set as url. We recommend you to set this parameter to the address of an Object Storage Service (OSS) system."
}

