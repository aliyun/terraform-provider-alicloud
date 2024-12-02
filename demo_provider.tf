terraform {
  required_providers {
    alicloud = {
      source  = "hashicorp/alicloud"
      version = "1.0.0"
      #            version = "1.232.0"
    }
    #    aws = {
    #      source  = "hashicorp/aws"
    #      version = "5.17.0"
    #    }
  }
}

variable "region" {
  default = "cn-hangzhou"
}

variable "profile_name" {
  default = "quanxi"
}
provider "alicloud" {
  profile = var.profile_name
  region  = var.region
}

provider "alicloud" {
  profile = var.profile_name
  alias   = "hz"
  region  = "cn-hangzhou"
}

provider "alicloud" {
  alias  = "assume"
  region = "cn-qingdao"
  #  profile = "quanxi"
  #  profile = "assumeRole"
  profile = "ak-for-assumerole"
  #  region = ""
  #  configuration_source = "xiaozhu-demo-for-assumerole"
  #  source_ip = "10.1.1.1"
  assume_role {
    role_arn           = "acs:ram::1182725234319447:role/assumeroletest"
    session_expiration = 900
    #    external_id = "terraformTestExternalId"
    #    "Condition": {
    #      "StringEquals": {
    #        "sts:ExternalId": "terraformTestExternalId"
    #      }
    #    },
  }
  #  assume_role_with_oidc {
  #    oidc_provider_arn = "acs:ram::1182725234319447:oidc-provider/ack-rrsa-c3470446df7a64cee9c6bf5ec949ea2ec"
  #    role_arn          = "acs:ram::1182725234319447:role/demo-role-for-rrsa"
  ##    oidc_token_file   = "./rrsa-token"
  ##    role_session_name = "terraform"
  #    oidc_token = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjE5UjRIbWNPekhEVGtJbk9iUVNFcXVQSGw3bURuNG4tLUxhOVZ4VkhUd1EifQ.eyJhdWQiOlsic3RzLmFsaXl1bmNzLmNvbSJdLCJleHAiOjE3MjMwNDIxMzcsImlhdCI6MTcyMzAzODUzNywiaXNzIjoiaHR0cHM6Ly9vaWRjLWFjay1jbi1iZWlqaW5nLm9zcy1jbi1iZWlqaW5nLmFsaXl1bmNzLmNvbS9jMzQ3MDQ0NmRmN2E2NGNlZTljNmJmNWVjOTQ5ZWEyZWMiLCJrdWJlcm5ldGVzLmlvIjp7Im5hbWVzcGFjZSI6InJyc2EtZGVtbyIsInBvZCI6eyJuYW1lIjoiZGVtbyIsInVpZCI6IjA5MjUwNmJkLTdjMjgtNGVjMS1hYjk3LTQyNmU1NTQ3OWRhNSJ9LCJzZXJ2aWNlYWNjb3VudCI6eyJuYW1lIjoiZGVtby1zYSIsInVpZCI6IjczNDRiOWZmLWFjMDktNGQyMy1hMjQ5LTJlNjE5NWJlY2UxZSJ9fSwibmJmIjoxNzIzMDM4NTM3LCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6cnJzYS1kZW1vOmRlbW8tc2EifQ.KRgbCUFhUx27jSjG885ZjECtCoHngccI77JMDwlV_t8O7-40WjHe0LSqTO0eqDDJEztcD63q8kCY5gvimHAd2bXrQ4TrS7P4Y1N095quBAegsQB9tlDU680pSVVUROGfJ30ZYs24q6FZXcFGlCH60TIFc6qgGnbOwPJeJeuio0vExGEtDzvqXN6498v1ifzT72yWpcTPlBFtxueLkAzx8-9OS7SL7jhv9czsrxopn5CC-tHCuZ-fSaCs7PtYqT8kUws9fNM2pvATikpfXDlxNI6vSipN2BN9DYn6YN4YdmFF4PG6P8lDaEFlbYmczrqK0arGeE6rNASQDklYfrUBfA"
  #    session_expiration = 3600
  #  }
  #  ecs_role_name = "xiaozhu"
}
