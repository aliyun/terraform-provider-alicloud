variable "ots_instance_name" {
  default = "tf-test"
}

variable "table_name" {
  default = "ots_table"
}

variable "primary_key_1_name" {
  default = "pk1"
}

variable "primary_key_2_name" {
  default = "pk2"
}

variable "primary_key_3_name" {
  default = "pk3"
}

variable "primary_key_4_name" {
  default = "pk4"
}

variable "primary_key_integer_type" {
  default = "Integer"
}

variable "primary_key_string_type" {
  default = "String"
}

variable "time_to_live" {
  default = -1
}

variable "max_version" {
  default = 1
}
