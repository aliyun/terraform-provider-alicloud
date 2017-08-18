## 0.1.1 (Unreleased)

IMPROVEMENTS:

- resource/rds: Add ability to import existing RDS resources [GH-16]
- datasource/alicloud_zones: Add more options for filtering [GH-19]

BUG FIXES:

- resource/disk_attachment: Fix issue attaching multiple disks and set disk_attachment's parameter 'device_name' as deprecated [GH-9]
- resource/rds: Fix diff error about rds security_ips [GH-13]
- resource/security_group_rule: Fix diff error when authorizing security group rules [GH-15]

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
