---
layout: "alicloud"
page_title: "Alibaba Cloud Provider Mirror Guide"
sidebar_current: "docs-alicloud-guide-alicloud-mirror-guide"
description: |-
  Setting alibaba cloud provider mirror configure
---

# Setting a Network Mirror

Alibaba Cloud Provider has two source: `source = hashicorp/alicloud` and `source = aliyun/alicloud`. 
Sometimes, running `terraform init` command to download and install provider will be too slow, 
or even though failed as following:

```text
- Finding aliyun/alicloud versions matching "1.191.0"...
╷
│ Error: Failed to query available provider packages
│
│ Could not retrieve the list of available versions for provider aliyun/alicloud: could not query provider registry for registry.terraform.io/aliyun/alicloud: the request
│ failed after 2 attempts, please try again later: Get "https://registry.terraform.io/v1/providers/aliyun/alicloud/versions": net/http: request canceled (Client.Timeout
│ exceeded while awaiting headers)
╵

- Finding hashicorp/alicloud versions matching "1.191.0"...
╷
│ Error: Failed to query available provider packages
│
│ Could not retrieve the list of available versions for provider hashicorp/alicloud: could not query provider registry for registry.terraform.io/hashicorp/alicloud: the
│ request failed after 2 attempts, please try again later: Get "https://registry.terraform.io/v1/providers/hashicorp/alicloud/versions": context deadline exceeded
│ (Client.Timeout exceeded while awaiting headers)
╵
```

Since Terraform CLI v0.13.2, it provides to set the [network_mirror](https://developer.hashicorp.com/terraform/cli/config/config-file#network_mirror) feature.
In order to fix above downloading `alicloud` provider failed issue, Alibaba Cloud provides mirror service, and you can set 
the following configuration in the [CLI Configuration File](https://developer.hashicorp.com/terraform/cli/config/config-file):
```terraform
provider_installation {
  network_mirror {
    url = "https://mirrors.aliyun.com/terraform/"
    // Setting alicloud from Alibaba Cloud Mirror Service
    include = ["registry.terraform.io/aliyun/alicloud", 
               "registry.terraform.io/hashicorp/alicloud",
              ]   
  }
  direct {
    // setting other providers from Terraform Registry
    exclude = ["registry.terraform.io/aliyun/alicloud", 
               "registry.terraform.io/hashicorp/alicloud",
              ]  
  }
}
```