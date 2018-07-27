---
layout: "alicloud"
page_title: "Provider: alicloud"
sidebar_current: "docs-alicloud-index"
description: |-
  The Alicloud provider is used to interact with many resources supported by Alicloud. The provider needs to be configured with the proper credentials before it can be used.
---

# Alicloud Provider

The Alicloud provider is used to interact with the
many resources supported by [Alicloud](https://www.aliyun.com). The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

-> **Note:** When you use terraform on the `Windowns` computer, please install [golang](https://golang.org/dl/) in your computer.
Otherwise, you may happen the issue from version 1.8.1 and the issue details can refer to [Crash Error](https://github.com/alibaba/terraform-provider/issues/469).


## Example Usage

```hcl
# Configure the Alicloud Provider
provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "alicloud_instance_types" "2c4g" {
  cpu_core_count = 2
  memory_size = 4
}

# Create a web server
resource "alicloud_instance" "web" {
  # cn-beijing
  image_id          = "ubuntu_140405_32_40G_cloudinit_20161115.vhd"
  internet_charge_type  = "PayByBandwidth"

  instance_type        = "${data.alicloud_instance_types.2c4g.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  security_groups      = ["${alicloud_security_group.default.id}"]
  instance_name        = "web"
  vswitch_id = "vsw-abc12345"
}

# Create security group
resource "alicloud_security_group" "default" {
  name        = "default"
  description = "default"
  vpc_id = "vpc-abc12345"
}
```

## Authentication

The Alicloud provider offers a flexible means of providing credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables

### Static credentials ###

Static credentials can be provided by adding an `access_key` `secret_key` and `region` in-line in the
alicloud provider block:

Usage:

```hcl
provider "alicloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}
```


###Environment variables

You can provide your credentials via `ALICLOUD_ACCESS_KEY` and `ALICLOUD_SECRET_KEY`,
environment variables, representing your Alicloud Access Key and Secret Key, respectively.
`ALICLOUD_REGION` is also used, if applicable:

```hcl
provider "alicloud" {}
```

Usage:

```shell
$ export ALICLOUD_ACCESS_KEY="anaccesskey"
$ export ALICLOUD_SECRET_KEY="asecretkey"
$ export ALICLOUD_REGION="cn-beijing"
$ terraform plan
```


## Argument Reference

The following arguments are supported:

* `access_key` - This is the Alicloud access key. It must be provided, but
  it can also be sourced from the `ALICLOUD_ACCESS_KEY` environment variable.

* `secret_key` - This is the Alicloud secret key. It must be provided, but
  it can also be sourced from the `ALICLOUD_SECRET_KEY` environment variable.

* `region` - This is the Alicloud region. It must be provided, but
  it can also be sourced from the `ALICLOUD_REGION` environment variables.

* `security_token` - Alicloud [Security Token Service](https://www.alibabacloud.com/help/doc-detail/66222.html).
It can be sourced from the `ALICLOUD_SECURITY_TOKEN`.

* `account_id` - (Optional) Alibaba Cloud Account ID. It is required for Function Compute Service.
It can be sourced from the `ALICLOUD_ACCOUNT_ID`

Nested `endpoints` block supports the following:

* `log_endpoint` - (Optional) The self-defined endpoint of log service, referring to [Service Endpoints](https://www.alibabacloud.com/help/doc-detail/29008.html).
It can be sourced from the `LOG_ENDPOINT`

* `fc` - (Optional) Use this to override the default endpoint
  URL constructed from the `region`. Referring to [Function Compute Service Endpoints](https://www.alibabacloud.com/help/doc-detail/52984.htm).
  It's typically used to connect to
  custom Function Compute Service endpoints.
  It can be sourced from the `FC_ENDPOINT`

## Testing

Credentials must be provided via the `ALICLOUD_ACCESS_KEY`, and `ALICLOUD_SECRET_KEY` environment variables in order to run acceptance tests.
