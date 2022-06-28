---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_simple_office_site"
sidebar_current: "docs-alicloud-resource-ecd-simple-office-site"
description: |-
  Provides a Alicloud ECD Simple Office Site resource.
---

# alicloud\_ecd\_simple\_office\_site

Provides a ECD Simple Office Site resource.

For information about ECD Simple Office Site and how to use it, see [What is Simple Office Site](https://help.aliyun.com/document_detail/188382.html).

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  bandwidth           = 5
  desktop_access_type = "Internet"
  office_site_name    = "site_name"
}

```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Deprecated from 1.142.0) The Internet Bandwidth Peak. It has been deprecated from version 1.142.0 and can be found in the new resource alicloud_ecd_network_package.
* `cen_id` - (Optional, ForceNew) Cloud Enterprise Network Instance ID.
* `cen_owner_id` - (Optional) The cen owner id.
* `cidr_block` - (Required, ForceNew) Workspace Corresponds to the Security Office Network of IPv4 Segment.
* `desktop_access_type` - (Optional, Computed) Connect to the Cloud Desktop Allows the Use of the Access Mode of. Valid values: `Any`, `Internet`, `VPC`.
* `enable_admin_access` - (Optional, ForceNew) Whether to Use Cloud Desktop User Empowerment of Local Administrator Permissions.
* `enable_cross_desktop_access` - (Optional) Enable Cross-Desktop Access.
* `enable_internet_access` - (Deprecated from 1.142.0) Whether the Open Internet Access Function.
* `mfa_enabled` - (Optional) Whether to Enable Multi-Factor Authentication MFA.
* `office_site_name` - (Optional) The office site name.
* `sso_enabled` - (Optional) Whether to Enable Single Sign-on (SSO) for User-Based SSO.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Simple Office Site.
* `status` - Workspace State. Valid Values: `REGISTERED`,`REGISTERING`.

## Import

ECD Simple Office Site can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_simple_office_site.example <id>
```
