---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_remote_write"
sidebar_current: "docs-alicloud-resource-arms-remote-write"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Remote Write resource.
---

# alicloud\_arms\_remote\_write

Provides a Application Real-Time Monitoring Service (ARMS) Remote Write resource.

For information about Application Real-Time Monitoring Service (ARMS) Remote Write and how to use it, see [What is Remote Write](https://www.alibabacloud.com/help/en/application-real-time-monitoring-service/latest/api-doc-arms-2019-08-08-api-doc-addprometheusremotewrite).

-> **NOTE:** Available in v1.204.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_arms_remote_write" "default" {
  cluster_id        = "your_cluster_id"
  remote_write_yaml = "remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth: {username: 666, password: '******'}\n  write_relabel_configs:\n  - source_labels: [instance_id]\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the Prometheus instance.
* `remote_write_yaml` - (Required) The details of the Remote Write configuration item. Specify the value in the YAML format.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Remote Write. It formats as `<cluster_id>:<remote_write_name>`.
* `remote_write_name` - The name of the Remote Write configuration item.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Remote Write.
* `update` - (Defaults to 3 mins) Used when update the Remote Write.
* `delete` - (Defaults to 3 mins) Used when delete the Remote Write.

## Import

Application Real-Time Monitoring Service (ARMS) Remote Write can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_remote_write.example <cluster_id>:<remote_write_name>
```