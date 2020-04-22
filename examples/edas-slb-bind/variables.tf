variable "app_id" {
  description = "The ID of the applicaton to which you want to bind an SLB instance."
}

variable "slb_id" {
  description = "The ID of the SLB instance that is going to be bound."
}

variable "slb_ip" {
  description = "The IP address that is allocated to the bound SLB instance."
}

variable "type" {
  description = "The type of the bound SLB instance."
}

variable "listener_port" {
  description = "The listening port for the bound SLB instance."
}

variable "vserver_group_id" {
  description = "The ID of the virtual server (VServer) group associated with the intranet SLB instance."
}

