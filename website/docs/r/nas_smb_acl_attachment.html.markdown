---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_smb_acl_attachment"
sidebar_current: "docs-alicloud-resource-nas-smb-acl-attachment"
description: |-
  Provides a Alicloud Nas Smb Acl resource.
---

# alicloud\_nas_smb_acl

Provides a Nas Smb Acl resource.

Alibaba Cloud SMB protocol file storage service supports user authentication based on AD domain system and permission access control at the file system level. Connecting and accessing the SMB file system as a domain user can implement the requirements for access control at the file and directory level in the SMB protocol file system. The current Alibaba Cloud SMB protocol file storage service does not support multi-user file and directory-level permission access control, and only provides file system-level authentication and access based on the whitelist mechanism that supports cloud accounts and source IP permission groups control.
-> **NOTE:** Available in 1.186.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_smb_acl_attachment&exampleId=ac126602-e234-c0e4-c3d0-ae14095df14e111f001a&activeTab=example&spm=docs.r.nas_smb_acl_attachment.0.ac126602e2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_nas_zones" "example" {
  file_system_type = "standard"
}

resource "alicloud_nas_file_system" "example" {
  protocol_type    = "SMB"
  storage_type     = "Capacity"
  description      = "terraform-example"
  encrypt_type     = "0"
  file_system_type = "standard"
  zone_id          = data.alicloud_nas_zones.example.zones[0].zone_id
}

resource "alicloud_nas_smb_acl_attachment" "example" {
  file_system_id = alicloud_nas_file_system.example.id
  keytab         = "BQIAAABHAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAAQAIqIx6v7p11oUAAABHAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAAwAIqIx6v7p11oUAAABPAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAFwAQnQZWB3RAPHU7PMIJyBWePAAAAF8AAgANQUxJQURURVNULkNPTQAEY2lmcwAZc21ic2VydmVyMjQuYWxpYWR0ZXN0LmNvbQAAAAEAAAAAAQASACAGJ7F0s+bcBjf6jD5HlvlRLmPSOW+qDZe0Qk0lQcf8WwAAAE8AAgANQUxJQURURVNULkNPTQAEY2lmcwAZc21ic2VydmVyMjQuYWxpYWR0ZXN0LmNvbQAAAAEAAAAAAQARABDdFmanrSIatnDDhoOXYadj"
  keytab_md5     = "E3CCF7E2416DF04FA958AA4513EA29E8"
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `keytab` - (Required) The string that is generated after the system encodes the keytab file by using Base64.
* `keytab_md5` - (Required) RThe string that is generated after the system encodes the keytab file by using MD5.
* `enable_anonymous_access` - (Optional) Specifies whether to allow anonymous access. Valid values:
                                         true: The file system allows anonymous access.
                                         false: The file system denies anonymous access. Default value: false.
* `encrypt_data` - (Optional) Specifies whether to enable encryption in transit. Valid values:
                                          true: enables encryption in transit.
                                          false: disables encryption in transit. Default value: false.
* `reject_unencrypted_access` - (Optional) Specifies whether to deny access from non-encrypted clients. Valid values:
                                           true: The file system denies access from non-encrypted clients.
                                           false: The file system allows access from non-encrypted clients. Default value: false.
* `super_admin_sid` - (Optional) The ID of a super admin. The ID must meet the following requirements:
                                         The ID starts with S and does not contain letters except S.
                                         The ID contains at least three hyphens (-) as delimiters.
                                         Example: S-1-5-22 and S-1-5-22-23.
* `home_dir_path` - (Optional) The home directory of each user. Each user-specific home directory must meet the following requirements:    
                                       Each segment starts with a forward slash (/) or a backslash (\).
                                       Each segment does not contain the following special characters: <>":?*.
                                       Each segment is 0 to 255 characters in length.
                                       The total length is 0 to 32,767 characters.
                                       For example, if you create a user named A and the home directory is /home, the file system automatically creates a directory named /home/A when User A logs on to the file system. If the /home/A directory already exists, the file system does not create the directory.
                                  

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is formate as `<file_system_id>`.
* `enabled` - Specifies whether to enable the ACL feature.
                  true: enables the ACL feature.
                  false: disables the ACL feature.
* `auth_method` - The method that is used to authenticate network identities.
