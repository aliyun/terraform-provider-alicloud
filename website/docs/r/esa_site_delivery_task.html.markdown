---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_site_delivery_task"
description: |-
  Provides a Alicloud ESA Site Delivery Task resource.
---

# alicloud_esa_site_delivery_task

Provides a ESA Site Delivery Task resource.



For information about ESA Site Delivery Task and how to use it, see [What is Site Delivery Task](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateSiteDeliveryTask).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_site_delivery_task&exampleId=9dd75504-7387-5da4-c633-0a4f018d075f263936bb&activeTab=example&spm=docs.r.esa_site_delivery_task.0.9dd7550473&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_http_example" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_site_delivery_task" "default" {
  http_delivery {
    standard_auth_param {
      private_key  = "***"
      url_path     = "v1/log/upload"
      expired_time = "300"
    }

    transform_timeout = "10"
    max_retry         = "3"
    max_batch_mb      = "5"
    compress          = "gzip"
    log_body_suffix   = "cdnVersion:1.0"
    standard_auth_on  = "false"
    log_body_prefix   = "cdnVersion:1.0"
    dest_url          = "http://11.177.129.13:8081"
    max_batch_size    = "1000"
  }

  data_center   = "global"
  discard_rate  = "0.0"
  task_name     = "dcdn-example-task"
  business_type = "dcdn_log_access_l1"
  field_name    = "ConsoleLog,CPUTime,Duration,ErrorCode,ErrorMessage,ResponseSize,ResponseStatus,RoutineName,ClientRequestID,LogTimestamp,FetchStatus,SubRequestID"
  delivery_type = "http"
  site_id       = alicloud_esa_site.resource_Site_http_example.id
}
```

## Argument Reference

The following arguments are supported:
* `business_type` - (Required) Real-time log type. Valid values:
  - `dcdn_log_access_l1 (default)`: access log.
  - `dcdn_log_er`: edge function log.
  - `dcdn_log_waf`: security protection log.
  - `dcdn_log_ipa`: 4 layer acceleration log.
* `data_center` - (Required, ForceNew) Data Center. Values:
  - `cn`: Mainland China.
  - `sg`: Global (excluding Mainland China).
* `delivery_type` - (Required, ForceNew) Delivery Type:
  - `sls`: Alibaba Cloud Log Service.
  - `http`: Http service.
  - `aws3`: Amazon s3 service.
  - `oss`: Alibaba Cloud Object Storage Service.
  - `kafka`: Kafka service.
  - `aws3cmpt`: Amazon s3 Compatible Service.
* `discard_rate` - (Optional, Float) If the discard rate is not filled, the default value is 0.
* `field_name` - (Required) The list of delivery fields to be modified, separated by commas.
* `http_delivery` - (Optional, List) HTTP delivery configuration parameters. See [`http_delivery`](#http_delivery) below.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `kafka_delivery` - (Optional, List) Kafka delivery configuration parameters. See [`kafka_delivery`](#kafka_delivery) below.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `oss_delivery` - (Optional, List) OSS delivery configuration. See [`oss_delivery`](#oss_delivery) below.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `s3_delivery` - (Optional, List) S3/S3 compatible delivery configuration parameters. See [`s3_delivery`](#s3_delivery) below.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the [ListSites](https://help.aliyun.com/document_detail/2850189.html) interface.
* `sls_delivery` - (Optional, List) SLS delivery configuration. See [`sls_delivery`](#sls_delivery) below.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `status` - (Optional, Computed) Task status, value:
  - `online`: push in.
  - `offline`: deactivated.
* `task_name` - (Required, ForceNew) The task name.

### `http_delivery`

The http_delivery supports the following:
* `compress` - (Optional) 
* `dest_url` - (Optional) 
* `header_param` - (Optional, Map) 
* `log_body_prefix` - (Optional) 
* `log_body_suffix` - (Optional) 
* `max_batch_mb` - (Optional, Int) 
* `max_batch_size` - (Optional, Int) 
* `max_retry` - (Optional, Int) 
* `query_param` - (Optional, Map) 
* `standard_auth_on` - (Optional) 
* `standard_auth_param` - (Optional, List)  See [`standard_auth_param`](#http_delivery-standard_auth_param) below.
* `transform_timeout` - (Optional, Int) 

### `http_delivery-standard_auth_param`

The http_delivery-standard_auth_param supports the following:
* `expired_time` - (Optional, Int) 
* `private_key` - (Optional) 
* `url_path` - (Optional) 

### `kafka_delivery`

The kafka_delivery supports the following:
* `balancer` - (Optional) 
* `brokers` - (Optional, List) 
* `compress` - (Optional) The compression method. By default, data is not compressed.
* `machanism_type` - (Optional) 
* `password` - (Optional) 
* `topic` - (Optional) 
* `user_auth` - (Optional) 
* `user_name` - (Optional) 

### `oss_delivery`

The oss_delivery supports the following:
* `aliuid` - (Optional) 
* `bucket_name` - (Optional) 
* `prefix_path` - (Optional) 
* `region` - (Optional) The region ID of the service.

### `s3_delivery`

The s3_delivery supports the following:
* `access_key` - (Optional) 
* `bucket_path` - (Optional) 
* `endpoint` - (Optional) 
* `prefix_path` - (Optional) 
* `region` - (Optional) 
* `s3_cmpt` - (Optional) 
* `secret_key` - (Optional) 
* `server_side_encryption` - (Optional) Server-side encryption
* `vertify_type` - (Optional) Authentication Type

### `sls_delivery`

The sls_delivery supports the following:
* `sls_log_store` - (Optional) 
* `sls_project` - (Optional) 
* `sls_region` - (Optional) 

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<task_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Site Delivery Task.
* `delete` - (Defaults to 5 mins) Used when delete the Site Delivery Task.
* `update` - (Defaults to 5 mins) Used when update the Site Delivery Task.

## Import

ESA Site Delivery Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_site_delivery_task.example <site_id>:<task_name>
```