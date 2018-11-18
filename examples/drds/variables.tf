variable "description" {
  default = "DRDS instance for RDS"
}
 variable "type" {
  default = "PRIVATE"
}

 variable "specification" {
  default = "drds.sn1.8c16g.16C32G"
}
 variable "pay_type" {
  default = "drdsPost"
}

variable "instance_series" {
	default = "drds.sn1.4c8g"
}

variable "name" {
	default = "tf-testaccDrdsdatabase_vpc"
}