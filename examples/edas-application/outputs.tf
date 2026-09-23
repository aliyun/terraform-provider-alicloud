output "application_name" {
  value       = alicloud_edas_application.default.application_name
  description = "Name of your EDAS application. Only letters '-' '_' and numbers are allowed. The length cannot exceed 36 characters"
}

output "cluster_id" {
  value       = alicloud_edas_application.default.cluster_id
  description = "The ID of the cluster that you want to create the application. The default cluster will be used if you do not specify this parameter. You can call the ListClusterOperation interface to query the cluster ID."
}

output "package_type" {
  value       = alicloud_edas_application.default.package_type
  description = "The type of the package for the deployment of the application that you want to create. The valid values are: WAR and JAR. We strongly recommend you to set this parameter when creating the application."
}

output "build_pack_id" {
  value       = alicloud_edas_application.default.build_pack_id
  description = "The package ID of Enterprise Distributed Application Service (EDAS) Container, which can be retrieved by calling container version list interface ListBuildPack or the \"Pack ID\" column in container version list. When creating High-speed Service Framework (HSF) application, this parameter is required."
}

output "descriotion" {
  value       = alicloud_edas_application.default.descriotion
  description = "The description of the application that you want to create."
}

output "health_check_url" {
  value       = alicloud_edas_application.default.health_check_url
  description = "The URL for health checking of the application."
}

output "logical_region_id" {
  value       = alicloud_edas_application.default.logical_region_id
  description = "The ID of the namespace where you want to create the application. You can call the ListUserDefineRegion operation to query the namespace ID."
}

output "component_ids" {
  value       = alicloud_edas_application.default.component_ids
  description = "The ID of the component in the container where the application is going to be deployed. If the runtime environment is not specified when the application is created and the application is not deployed, you can set the parameter as fellow: when deploying a native Dubbo or Spring Cloud application using a WAR package for the first time, you must specify the version of the Apache Tomcat component based on the deployed application. You can call the ListClusterOperation interface to query the components. When deploying a non-native Dubbo or Spring Cloud application using a WAR package for the first time, you can leave this parameter empty."
}

output "ecu_info" {
  value       = alicloud_edas_application.default.ecu_info
  description = "The IDs of the Elastic Compute Unit (ECU) where you want to deploy the application."
}

output "group_id" {
  value       = alicloud_edas_application.default.group_id
  description = "The ID of the instance group where the application is going to be deployed. Set this parameter to all if you want to deploy the application to all groups."
}

output "package_version" {
  value       = alicloud_edas_application.default.package_version
  description = "The version of the application that you want to deploy. It must be unique for every application. The length cannot exceed 64 characters. We recommended you to use a timestamp. "
}

output "war_url" {
  value       = alicloud_edas_application.default.war_url
  description = "The address to store the uploaded web application (WAR) package for application deployment. This parameter is required when the deployType parameter is set as url. We recommend you to set this parameter to the address of an Object Storage Service (OSS) system."
}