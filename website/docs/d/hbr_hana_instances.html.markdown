---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_hana_instances"
sidebar_current: "docs-alicloud-datasource-hbr-hana-instances"
description: |-
  Provides a list of Hbr Hana Instances to the user.
---

# alicloud\_hbr\_hana\_instances

This data source provides the Hbr Hana Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.178.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_hana_instances" "ids" {
  ids = ["example_id"]
}
output "hbr_hana_instance_id_1" {
  value = data.alicloud_hbr_hana_instances.ids.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Hana Instance IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Hana Instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vault_id` - (Optional, ForceNew) The id of the vault.
* `status` - (Optional, ForceNew) The status of the SAP HANA instance. Valid values:
  - `INITIALIZING`: The instance is being initialized. 
  - `INITIALIZED`: The instance is registered. 
  - `INVALID_HANA_NODE`: The instance is invalid. 
  - `INITIALIZE_FAILED`: The client fails to be installed on the instance.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Hbr Hana Instance names.
* `instances` - A list of Hbr Hana Instances. Each element contains the following attributes:
	* `alert_setting` - The alert settings. Valid value: `INHERITED`, which indicates that the backup client sends alert notifications in the same way as the backup vault.
	* `hana_instance_id` - The ID of the SAP HANA instance.
	* `hana_name` - The name of the SAP HANA instance.
	* `host` - The private or internal IP address of the host where the primary node of the SAP HANA instance resides.
	* `id` - The ID of the Hana Instance. The value formats as `<vault_id>:<hana_instance_id>`.
	* `instance_number` - The instance number of the SAP HANA system.
	* `resource_group_id` - The ID of the resource group.
	* `status` - The status of the SAP HANA instance.
	* `status_message` - The status information.
	* `use_ssl` - Indicates whether the SAP HANA instance is connected over Secure Sockets Layer (SSL).
	* `user_name` - The username of the SYSTEMDB database.
	* `validate_certificate` - Indicates whether the SSL certificate of the SAP HANA instance is verified.
	* `vault_id` - The ID of the backup vault.