---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_endpoint"
sidebar_current: "docs-alicloud-resource-poalrdb-endpoint"
description: |-
  Provides a PolarDB instance endpoint resource.
---

# alicloud_polardb_endpoint

Provides a PolarDB endpoint resource to manage custom endpoint of PolarDB cluster.

-> **NOTE:** Available since v1.80.0.

-> **NOTE:** After v1.80.0 and before v1.121.0, you can only use this resource to manage the custom endpoint. Since v1.121.0, you also can import the primary endpoint and the cluster endpoint, to modify their ssl status and so on.

-> **NOTE:** The primary endpoint and the default cluster endpoint can not be created or deleted manually.

-> **NOTE:** A PolarDB cluster has only one default cluster endpoint, and it can be managed either by this resource (after importing it) or by `alicloud_polardb_cluster_endpoint`. Manage one endpoint address from only one Terraform resource. If the same default cluster endpoint is controlled by both resources, a change to an attribute they share (such as `scc_mode`) is sent to that endpoint by both resources: when the two configurations specify different values, the later apply overwrites the earlier one and the other configuration keeps reporting a diff. Different endpoint addresses have different identities and can be managed by different resources: for example, a custom endpoint created by this resource and the default cluster endpoint managed by `alicloud_polardb_cluster_endpoint` can coexist. `scc_mode` (global consistency, high-performance mode) is available only on PolarDB for MySQL Enterprise Edition clusters whose engine minor version is `8.0.2.2.19` or later (MySQL `8.0.2`), `8.0.1.1.29` or later (MySQL `8.0.1`), or `5.7.1.0.26` or later (MySQL `5.7`); on other clusters the setting is not supported. On a cluster that meets these prerequisites, endpoint identity alone does not determine the scope: `scc_mode` can be enabled or disabled independently for each endpoint only when PolarProxy is version `2.9.12` or later and the engine minor version is `8.0.2.2.29` or later (MySQL `8.0.2`), `8.0.1.1.48` or later (MySQL `8.0.1`), or `5.7.1.0.38` or later (MySQL `5.7`); in that case different endpoint addresses can control `scc_mode` separately. If the additional conditions are not met, `scc_mode` is a cluster-level setting: enabling or disabling it on one endpoint also applies it to all other endpoints of the same cluster, so the Terraform resources that manage endpoints of the same cluster must use the same value and cannot expect each endpoint address to control it independently.

-> **NOTE:** The DynamoDB type endpoint does not support deletion. When destroying the Terraform resource, the DynamoDB endpoint will be removed from state but not deleted from the cloud.

-> **NOTE:** The endpoint's advanced configuration must stay compatible with global consistency (high-performance mode) and the endpoint's read/write mode. `ConsistLevel` and `CausalConsistRead` are separate settings. An incompatible combination can be rejected with an `EndpointConfig.Conflict` error; an error requiring `CausalConsistRead=session` does not imply that `ConsistLevel` must be `0`.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_endpoint&exampleId=413eb2b4-ee47-16e0-e2f9-a5506e67470c69a6240d&activeTab=example&spm=docs.r.polardb_endpoint.0.413eb2b4ee&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

data "alicloud_polardb_node_classes" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  category      = "Normal"
  db_node_class = "polar.mysql.x4.medium"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = "terraform-example"
}

resource "alicloud_polardb_endpoint" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  endpoint_type = "Custom"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_endpoint&spm=docs.r.polardb_endpoint.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `endpoint_type` - (Optional, ForceNew) Type of the endpoint. Before v1.121.0, it only can be `Custom`. since v1.121.0, `Custom`, `Cluster`, `Primary` are valid, default to `Custom`. However when creating a new endpoint, it also only can be `Custom`.
* `read_write_mode` - (Optional) Read or write mode. Valid values are `ReadWrite`, `ReadOnly`. When creating a new custom endpoint, default to `ReadOnly`.
* `nodes` - (Optional) Node id list for endpoint configuration. At least 2 nodes if specified, or if the cluster has more than 3 nodes, read-only endpoint is allowed to mount only one node. Default is all nodes.
* `auto_add_new_nodes` - (Optional) Whether the new node automatically joins the default cluster address. Valid values are `Enable`, `Disable`. When creating a new custom endpoint, default to `Disable`.
* `endpoint_config` - (Optional) The advanced settings of the endpoint of Apsara PolarDB clusters are in JSON format. Including the settings of consistency level, transaction splitting, connection pool, and offload reads from primary node. For more details, see the [description of EndpointConfig in the Request parameters table for details](https://www.alibabacloud.com/help/doc-detail/116593.htm).
* `scc_mode` - (Optional, Computed, String, Available since v2.3.0) Specifies whether to enable global consistency (high-performance mode) for the endpoint. Valid values: `on`, `off`. Global consistency (high-performance mode) is supported only on PolarDB for MySQL Enterprise Edition clusters whose engine minor version is `8.0.2.2.19` or later (MySQL `8.0.2`), `8.0.1.1.29` or later (MySQL `8.0.1`), or `5.7.1.0.26` or later (MySQL `5.7`); whether the setting applies to a single endpoint or to the whole cluster depends on the PolarProxy and engine versions described in the resource-level note above. For `endpoint_config.ConsistLevel`, `0` means eventual consistency, `1` means session consistency (the API default), and `2` means global consistency; a `ReadOnly` endpoint supports only `0`. Select the value for the endpoint's read/write mode and required proxy consistency behavior; `0` is not a universal prerequisite for setting `scc_mode` to `on`. If this argument is omitted, the provider does not send `SccMode` when the endpoint is created, and the value reported by the service is read back into the state. Removing this argument from the configuration after it has been set does not send an update request; the value last read from the service is kept in the state.
* `ssl_enabled` - (Optional, Available since v1.121.0) Specifies how to modify the SSL encryption status. Valid values: `Disable`, `Enable`, `Update`.
* `net_type` - (Optional, Available since v1.121.0) The network type of the endpoint address.
* `ssl_auto_rotate` - (Optional, Available since v1.132.0) Specifies whether automatic rotation of SSL certificates is enabled. Valid values: `Enable`,`Disable`.  

  -> **NOTE:** For a PolarDB for MySQL cluster, this parameter is required, and only one connection string in each endpoint can enable the ssl, for other notes, see [Configure SSL encryption](https://www.alibabacloud.com/help/doc-detail/153182.htm).
    For a PolarDB for PostgreSQL cluster or a PolarDB-O cluster, this parameter is not required, by default, SSL encryption is enabled for all endpoints.
* `db_endpoint_description` - (Optional, Available since v1.201.0) The name of the endpoint.
* `connection_prefix` - (Optional, Available since v1.217.0) Prefix of the specified endpoint. The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter.
* `port` - (Optional, Available since v1.217.0) Port of the specified endpoint. Valid values: 3000 to 5999.

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<db_cluster_id>:<db_endpoint_id>`.
* `ssl_connection_string` - (Available since v1.121.0) The SSL connection string.
* `ssl_expire_time` - (Available since v1.121.0) The time when the SSL certificate expires. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `db_endpoint_id` - (Available since v1.161.0) The ID of the cluster endpoint.
* `ssl_certificate_url` - (Available since v1.132.0) Specifies SSL certificate download link.

## Import

PolarDB endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_endpoint.example pc-abc123456:pe-abc123456
```
