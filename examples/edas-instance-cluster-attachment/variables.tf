variable "cluster_id" {
  description = "The ID of the cluster that you want to create the application. The default cluster will be used if you do not specify this parameter."
}

variable "instance_ids" {
  type        = list(string)
  description = "The ID of instance. e.g. instanceId1, instanceId2."
}