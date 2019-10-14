---
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch_instance"
sidebar_current: "docs-alicloud-resource-elasticsearch-instance"
description: |-
  Provides a Alicloud Elasticsearch instance resource.
---

# alicloud\_elasticsearch\_instance

Provides a Elasticsearch instance resource. It contains data nodes, dedicated master node(optional) and etc. It can be associated with private IP whitelists and kibana IP whitelist.

-> **NOTE:** Only one operation is supported in a request. So if `data_node_spec` and `data_node_disk_size` are both changed, system will respond error.

-> **NOTE:** At present, `version` can not be modified once instance has been created.

## Example Usage

Basic Usage

```
resource "alicloud_elasticsearch_instance" "instance" {
  instance_charge_type = "PostPaid"
  data_node_amount     = "2"
  data_node_spec       = "elasticsearch.sn2ne.large"
  data_node_disk_size  = "20"
  data_node_disk_type  = "cloud_ssd"
  vswitch_id           = "some vswitch id"
  password             = "Your password"
  version              = "5.5.3_with_X-Pack"
  description          = "description"
  zone_count           = "2"
}
```
## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of instance. It a string of 0 to 30 characters.
* `instance_charge_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`, Default to `PostPaid`.
* `period` - (Optional) The duration that you will buy Elasticsearch instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. Default to 1.
* `data_node_amount` - (Required) The Elasticsearch cluster's data node quantity, between 2 and 50.
* `data_node_spec` - (Required) The data node specifications of the Elasticsearch instance.
* `data_node_disk_size` - (Required) The single data node storage space.
  - `cloud_ssd`: An SSD disk, supports a maximum of 2048 GiB (2 TB).
  - `cloud_efficiency` An ultra disk, supports a maximum of 5120 GiB (5 TB). If the data to be stored is larger than 2048 GiB, an ultra disk can only support the following data sizes (GiB): [`2560`, `3072`, `3584`, `4096`, `4608`, `5120`].
* `data_node_disk_type` - (Required) The data node disk type. Supported values: cloud_ssd, cloud_efficiency.
* `vswitch_id` - (Required, ForceNew) The ID of VSwitch.
* `password` - (Optional, Sensitive) The password of the instance. The password can be 8 to 32 characters in length and must contain three of the following conditions: uppercase letters, lowercase letters, numbers, and special characters (!@#$%^&*()_+-=).
* `kms_encrypted_password` - (Optional, Available in 1.57.1+) An KMS encrypts password used to a instance. It is conflicted with `password`, but you have to specify one of `password` and `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `version` - (Required, ForceNew) Elasticsearch version. Supported values: `5.5.3_with_X-Pack`, `6.3_with_X-Pack` and `6.7_with_X-Pack`.
* `private_whitelist` - (Optional) Set the instance's IP whitelist in VPC network.
* `kibana_whitelist` - (Optional) Set the Kibana's IP whitelist in internet network.
* `master_node_spec` - (Optional) The dedicated master node spec. If specified, dedicated master node will be created.
* `zone_count` - (Optional, Available in 1.44.0+) The Multi-AZ supported for Elasticsearch, between 1 and 3. The `data_node_amount` value must be an integral multiple of the `zone_count` value.

### Timeouts

-> **NOTE:** Available in 1.48.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 120 mins) Used when creating the elasticsearch instance (until it reaches the initial `active` status). 
* `update` - (Defaults to 120 mins) Used when activating the elasticsearch instance when necessary during update - e.g. when changing elasticsearch instance description, whitelist, data node settings, master node spec and password.
* `delete` - (Defaults to 120 mins) Used when terminating the elasticsearch instance. `Note`: There are 5 minutes to sleep to eusure the instance is deleted. It is not in the timeouts configure.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Elasticsearch instance.
* `domain` - Instance connection domain (only VPC network access supported).
* `port` - Instance connection port.
* `kibana_domain` - Kibana console domain (Internet access supported).
* `kibana_port` - Kibana console port.
* `status` - The Elasticsearch instance status. Includes `active`, `activating`, `inactive`. Some operations are denied when status is not `active`.

## Import

Elasticsearch can be imported using the id, e.g.

```
$ terraform import alicloud_elasticsearch.example es-cn-abcde123456
```

