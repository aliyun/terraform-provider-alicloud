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


variable "integer_type" {
  default = "Integer"
}

variable "string_type" {
  default = "String"
}

variable "binary_type" {
  default = "Binary"
}

variable "boolean_type" {
  default = "Boolean"
}

variable "double_type" {
  default = "Double"
}

variable "defined_column_1_name" {
  default = "col1"
}

variable "defined_column_2_name" {
  default = "col2"
}

variable "defined_column_3_name" {
  default = "col3"
}

variable "time_to_live" {
  default = -1
}

variable "max_version" {
  default = 1
}

variable "secondary_index_name" {
  default = "sec_index_1"
}

variable "secondary_index_type" {
  default = "Global"
}

variable "secondary_index_include_base_data" {
  default = true
}
variable "secondary_index_pks" {
  default = ["pk3", "pk2"]
}

variable "index_defined_cols" {
  default = ["col2", "col3"]
}


