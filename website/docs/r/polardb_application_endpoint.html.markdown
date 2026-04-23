---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_application_endpoint"
sidebar_current: "docs-alicloud-resource-polardb-application-endpoint"
description: |-
  Provides a PolarDB Application Endpoint resource.
---

# alicloud_polardb_application_endpoint

Provides a PolarDB Application Endpoint resource. This resource is used to manage the network access endpoints (such as Public Network) for a specific PolarDB AI Application.

-> **NOTE:** Available since v1.278.0.

## Example Usage

Create a PolarDB Application Endpoint

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_application_endpoint&exampleId=example-endpoint-01&activeTab=example&spm=docs.r.polardb_application_endpoint.0.example&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Assume you have already created a PolarDB Application
resource "alicloud_polardb_application_endpoint" "default" {
  application_id = "pa-xxx"
  endpoint_id    = "pa-xxx"
  net_type       = "Public"
}
```

### Removing alicloud_polardb_application_endpoint from your configuration

The `alicloud_polardb_application_endpoint` resource allows you to manage the network endpoints of your PolarDB application. Removing this resource from your configuration will remove it from your statefile and management, but may not disable the network access immediately depending on the underlying API behavior. Please verify the status in the PolarDB Console.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_application_endpoint&spm=docs.r.polardb_application_endpoint.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, ForceNew) The ID of the PolarDB Application to which the endpoint belongs.
* `endpoint_id` - (Required, ForceNew) The ID of the endpoint. This is usually obtained from the application details or list interfaces.
* `net_type` - (Required, ForceNew) The network type of the endpoint. Valid value: `Public`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PolarDB Application Endpoint. The value is formatted as `<application_id>:<endpoint_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the polardb application endpoint.
* `delete` - (Defaults to 5 mins) Used when deleting the polardb application endpoint.

## Import

PolarDB Application Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_application_endpoint.example pa-abc12345678:pa-def4567890
```