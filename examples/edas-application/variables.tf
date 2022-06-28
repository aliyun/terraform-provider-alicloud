variable "application_name" {
  description = "Name of your EDAS application. Only letters '-' '_' and numbers are allowed. The length cannot exceed 36 characters"
}

variable "cluster_id" {
  description = "The ID of the cluster that you want to create the application. The default cluster will be used if you do not specify this parameter."
}

variable "package_type" {
  description = "The type of the package for the deployment of the application that you want to create. The valid values are: WAR and JAR. We strongly recommend you to set this parameter when creating the application."
}

variable "build_pack_id" {
  type        = number
  description = "The package ID of Enterprise Distributed Application Service (EDAS) Container, which can be retrieved by calling container version list interface ListBuildPack or the \"Pack ID\" column in container version list. When creating High-speed Service Framework (HSF) application, this parameter is required."
}

variable "descriotion" {
  description = "The description of the application that you want to create."
}

variable "health_check_url" {
  description = "The URL for health checking of the application."
}

variable "logical_region_id" {
  description = "The ID of the namespace where you want to create the application."
}

variable "component_ids" {
  description = "The ID of the component in the container where the application is going to be deployed. If the runtime environment is not specified when the application is created and the application is not deployed, you can set the parameter as fellow: when deploying a native Dubbo or Spring Cloud application using a WAR package for the first time, you must specify the version of the Apache Tomcat component based on the deployed application. You can call the ListClusterOperation interface to query the components. When deploying a non-native Dubbo or Spring Cloud application using a WAR package for the first time, you can leave this parameter empty."
}

variable "ecu_info" {
  type        = list(string)
  description = "The IDs of the Elastic Compute Unit (ECU) where you want to deploy the application."
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
