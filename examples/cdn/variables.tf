variable "domain_name" {
  default = "www.example.com"
}

variable "cdn_type" {
  default = "web"
}

variable "sources" {
  type = list(string)

  default = [
    "example.com",
    "example.net",
    "example.cn",
  ]
}

variable "source_type" {
  default = "domain"
}

variable "enable" {
  default = "off"
}

variable "page_type" {
  default = "other"
}

variable "refer_type" {
  default = "block"
}

variable "auth_type" {
  default = "type_a"
}

variable "block_ips" {
  type = list(string)

  default = [
    "1.2.3.4",
    "111.222.111.111",
  ]
}

variable "hash_key_args" {
  type = list(string)

  default = [
    "youyouyou",
    "checkitout",
  ]
}

variable "refer_list" {
  type = list(string)

  default = [
    "www.example.com",
    "www.example.net",
  ]
}
