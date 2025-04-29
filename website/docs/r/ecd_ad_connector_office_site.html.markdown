---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_ad_connector_office_site"
sidebar_current: "docs-alicloud-resource-ecd-ad-connector-office-site"
description: |-
  Provides a Alicloud ECD Ad Connector Office Site resource.
---

# alicloud_ecd_ad_connector_office_site

Provides a ECD Ad Connector Office Site resource.

For information about ECD Ad Connector Office Site and how to use it, see [What is Ad Connector Office Site](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-ecd-2020-09-30-createadconnectorofficesite).

-> **NOTE:** Available since v1.176.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecd_ad_connector_office_site&exampleId=54b04686-69f8-771d-b0f4-515c3c40112ed648b315&activeTab=example&spm=docs.r.ecd_ad_connector_office_site.0.54b0468669&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_ecd_ad_connector_office_site" "default" {
  ad_connector_office_site_name = var.name
  bandwidth                     = 100
  cen_id                        = alicloud_cen_instance.default.id
  cidr_block                    = "10.0.0.0/12"
  desktop_access_type           = "INTERNET"
  dns_address                   = ["127.0.0.2"]
  domain_name                   = "corp.example.com"
  domain_password               = "Example1234"
  domain_user_name              = "sAMAccountName"
  enable_admin_access           = false
  enable_internet_access        = false
  mfa_enabled                   = false
  sub_domain_dns_address        = ["127.0.0.3"]
  sub_domain_name               = "child.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `ad_connector_office_site_name` - (Required) The name of the workspace. The name must be 2 to 255 characters in length. It must start with a letter and cannot start with `http://` or `https://`. It can contain digits, colons (:), underscores (_), and hyphens (-).
* `ad_hostname` - (Optional) The ad hostname.
* `bandwidth` - (Optional, ForceNew) The maximum public bandwidth value. Valid values: 0 to 200. If you do not specify this parameter or you set this parameter to 0, Internet access is disabled.
* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `cen_owner_id` - (Optional) The cen owner id.
* `cidr_block` - (Required, ForceNew) Workspace Corresponds to the Security Office Network of IPv4 Segment.
* `desktop_access_type` - (Optional, Computed) The method that you use to connect to cloud desktops. **Note:** The VPC connection method is provided by Alibaba Cloud PrivateLink. You are not charged for PrivateLink. When you set this parameter to VPC or Any, PrivateLink is automatically activated. Default value: `INTERNET`. Valid values:
  - `INTERNET`: connects clients to cloud desktops only over the Internet.
  - `VPC`: connects clients to cloud desktops only over a VPC.
  - `ANY`: connects clients to cloud desktops over the Internet or a VPC. You can select a connection method when you use a client to connect to the cloud desktop.
* `dns_address` - (Required) The IP address N of the DNS server of the enterprise AD system. You can specify only one IP address.
* `domain_name` - (Required) The domain name of the enterprise AD system. You can register each domain name only once.
* `domain_password` - (Optional) The password of the domain administrator. The password can be up to 64 characters in length.
* `domain_user_name` - (Optional) The username of the domain administrator. The username can be up to 64 characters in length.
* `enable_admin_access` - (Optional, ForceNew, Computed) Specifies whether to grant the permissions of the local administrator to the desktop users. Default value: true.
* `enable_internet_access` - (Optional, ForceNew, Computed) Specifies whether to enable Internet access.
* `mfa_enabled` - (Optional) Specifies whether to enable multi-factor authentication (MFA).
* `protocol_type` - (Optional) The protocol type. Valid values: `ASP`, `HDX`.
* `specification` - (Optional) The AD Connector specifications. Valid values: `1`, `2`.
* `sub_domain_dns_address` - (Optional) The DNS address N of the enterprise AD subdomain. If you specify a value for the `sub_domain_name` parameter but you do not specify a value for this parameter, the DNS address of the subdomain is the same as the DNS address of the parent domain.
* `sub_domain_name` - (Optional) The domain name of the enterprise AD subdomain.
* `verify_code` - (Optional) The verification code. If the CEN instance that you specify for the CenId parameter belongs to another Alibaba Cloud account, you must call the SendVerifyCode operation to obtain the verification code.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ad Connector Office Site.
* `status` - The resource State.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Ad Connector Office Site.
* `delete` - (Defaults to 1 mins) Used when delete the Ad Connector Office Site.


## Import

ECD Ad Connector Office Site can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_ad_connector_office_site.example <id>
```