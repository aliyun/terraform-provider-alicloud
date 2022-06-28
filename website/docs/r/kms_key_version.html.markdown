---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_key_version"
sidebar_current: "docs-alicloud-resource-kms-key-version"
description: |-
  Provides a Alikms key version resource.
---

# alicloud\_kms\_key\_version

Provides a Alikms Key Version resource. For information about Alikms Key Version and how to use it, see [What is Resource Alikms Key Version](https://www.alibabacloud.com/help/doc-detail/133838.htm).

-> **NOTE:** Available in v1.85.0+.

## Example Usage

Basic Usage

```
resource "alicloud_kms_key" "this" {}

resource "alicloud_kms_key_version" "keyversion" {
  key_id = alicloud_kms_key.this.id
}
```
## Argument Reference

The following arguments are supported:

* `key_id` - (Required, ForceNew) The id of the master key (CMK).

-> **NOTE:** The minimum interval for creating a Alikms key version is 7 days.


## Attributes Reference

* `creation_date` - (Removed from v1.124.4) The date and time (UTC time) when the Alikms key version was created.
* `key_id` - The id of the master key (CMK).
* `key_version_id` - The id of the Alikms key version.


## Import

Alikms key version can be imported using the id, e.g.

```
$ terraform import alicloud_kms_key_version.example 72da539a-2fa8-4f2d-b854-*****	
```
