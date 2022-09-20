terraform {
  backend "oss" {}
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}