variable "cluster_id" {
  description = "The cluster Id of your kubernetes cluster."
}

variable "utilization" {
  description = "The utilization option of autoscaler."
  default = "0.5"
}

variable "cool_down_duration" {
  description = "The cool_down_duration options of autoscaler."
  default = "1m"
}

variable "defer_scale_in_duration" {
  description = "The defer_scale_in_duration option of autoscaler."
  default = "1m"
}