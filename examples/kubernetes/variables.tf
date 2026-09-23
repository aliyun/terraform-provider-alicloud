variable "k8s_number" {
  description = "The number of kubernetes cluster."
  default     = 1
}

variable "availability_zone" {
  description = "The availability zones of vswitches."
  default     = ["cn-hangzhou-b", "cn-hangzhou-e", "cn-hangzhou-f"]
}

# leave it to empty would create a new one
variable "vpc_id" {
  description = "Existing vpc id used to create several vswitches and other resources."
  default     = ""
}

variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "10.0.0.0/8"
}

# leave it to empty then terraform will create several vswitches
variable "vswitch_ids" {
  description = "List of existing vswitch id."
  type        = list(string)
  default     = []
}


variable "vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'vswitch_ids' is not specified."
  type        = list(string)
  default     = ["10.1.0.0/16", "10.2.0.0/16", "10.3.0.0/16"]
}

variable "new_nat_gateway" {
  description = "Whether to create a new nat gateway. In this template, a new nat gateway will create a nat gateway, eip and server snat entries."
  default     = "true"
}

# 3 masters is default settings,so choose three appropriate instance types in the availability zones above.
variable "master_instance_types" {
  description = "The ecs instance types used to launch master nodes."
  default     = ["ecs.n4.xlarge", "ecs.n4.xlarge", "ecs.sn1ne.xlarge"]
}

variable "worker_instance_types" {
  description = "The ecs instance types used to launch worker nodes."
  default     = ["ecs.sn1ne.xlarge", "ecs.n4.xlarge"]
}

# options: between 24-28
variable "node_cidr_mask" {
  description = "The node cidr block to specific how many pods can run on single node."
  default     = 24
}

variable "enable_ssh" {
  description = "Enable login to the node through SSH."
  default     = true
}

variable "install_cloud_monitor" {
  description = "Install cloud monitor agent on ECS."
  default     = true
}

# options: none|static
variable "cpu_policy" {
  description = "kubelet cpu policy.default: none."
  default     = "none"
}

# options: ipvs|iptables
variable "proxy_mode" {
  description = "Proxy mode is option of kube-proxy."
  default     = "ipvs"
}

variable "password" {
  description = "The password of ECS instance."
  default     = "Just4Test"
}

variable "worker_number" {
  description = "The number of worker nodes in kubernetes cluster."
  default     = 4
}

# k8s_pod_cidr is only for flannel network
variable "pod_cidr" {
  description = "The kubernetes pod cidr block. It cannot be equals to vpc's or vswitch's and cannot be in them."
  default     = "172.20.0.0/16"
}

variable "service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  default     = "172.21.0.0/20"
}


variable "cluster_addons" {
  description = "Addon components in kubernetes cluster"

  type = list(object({
    name   = string
    config = string
  }))

  default = [
    {
      "name"   = "flannel",
      "config" = "",
    },
    {
      "name"   = "flexvolume",
      "config" = "",
    },
    {
      "name"   = "alicloud-disk-controller",
      "config" = "",
    },
    {
      "name"   = "logtail-ds",
      "config" = "{\"IngressDashboardEnabled\":\"true\"}",
    },
    {
      "name"   = "nginx-ingress-controller",
      "config" = "{\"IngressSlbNetworkType\":\"internet\"}",
    },
  ]
}

