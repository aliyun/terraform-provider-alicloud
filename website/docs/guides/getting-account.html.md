---
layout: "alicloud"
page_title: "Alibaba Cloud Account Guide"
sidebar_current: "docs-alicloud-guide-alicloud-account-guide"
description: |-
  Sign up Alibaba Cloud and distinguish its type.
---

# Getting an Alibaba Cloud Account

The Alibaba Cloud has two accounts: International-Site Account and China-Site Account.
In most cases the account type makes no difference when creating Alibaba Cloud resources.
However, based on some internal reason, when using Terraform to manage cloud resources
a few products and resources have different limitations when using different account types.
We will introduce the limitations gradually to help you avoid some needless errors.

## Sign Up for an Alibaba Cloud International-Site Account

-> **Warning:** At present, Terraform can not use an international-site account to open `Subscription`
resources which instance charge type is "PrePaid"

To sign up for an International-Site account, visit [Alibaba Cloud International-Site Website](https://www.alibabacloud.com/). For more account registration details, see [Sign up with Alibaba Cloud](https://www.alibabacloud.com/help/doc-detail/50482.html)

## Sign Up for an Alibaba Cloud China-Site Account

The China-Site website has a different URL. To sign up for a China-Site account, visit
[Alibaba Cloud China-Site Website](https://www.aliyun.com/).
For more account registration details, see [Sign up with Alibaba Cloud](https://help.aliyun.com/knowledge_detail/37195.html)

## How to Distinguish Account Types

There is a simple method to check whether an Alibaba Cloud account is an International-Site or China-Site account:
An account can only access its corresponding site. That is, an International-Site account can only login to the [International-Site Website](https://www.alibabacloud.com/) and
a China-Site account to the [China-Site Website](https://www.aliyun.com/).

