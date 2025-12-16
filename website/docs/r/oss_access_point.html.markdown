---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_access_point"
description: |-
  Provides a Alicloud OSS Access Point resource.
---

# alicloud_oss_access_point

Provides a OSS Access Point resource.

You can create multiple Access points for buckets and configure different Access control permissions and network control policies for different Access points.

For information about OSS Access Point and how to use it, see [What is Access Point](https://www.alibabacloud.com/help/en/oss/developer-reference/createaccesspoint).

-> **NOTE:** Available since v1.240.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_access_point&exampleId=dbb007ad-389a-f551-a332-43d2c9252c6f486665ad&activeTab=example&spm=docs.r.oss_access_point.0.dbb007ad38&intl_lang=EN_US" target="_blank">
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

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
}


resource "alicloud_oss_access_point" "default" {
  access_point_name = var.name
  bucket            = alicloud_oss_bucket.CreateBucket.bucket
  vpc_configuration {
    vpc_id = "vpc-abcexample"
  }
  network_origin = "vpc"
  public_access_block_configuration {
    block_public_access = true
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oss_access_point&spm=docs.r.oss_access_point.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `access_point_name` - (Required, ForceNew) The name of the access point
* `bucket` - (Required, ForceNew) The Bucket to which the current access point belongs.
* `network_origin` - (Required, ForceNew) Access point network source. The valid values are as follows: 
  - vpc: only the specified VPC ID can be used to access the access point. 
  - internet: the access point can be accessed through both external and internal Endpoint.
* `public_access_block_configuration` - (Optional, List) Configuration of Access Point Blocking Public Access See [`public_access_block_configuration`](#public_access_block_configuration) below.
* `vpc_configuration` - (Optional, ForceNew, List) If the Network Origin is vpc, the VPC source information is saved here. See [`vpc_configuration`](#vpc_configuration) below.

### `public_access_block_configuration`

The public_access_block_configuration supports the following:
* `block_public_access` - (Optional, Computed) Block public access enabled for access point

### `vpc_configuration`

The vpc_configuration supports the following:
* `vpc_id` - (Optional, ForceNew) The vpc ID is required only when the value of NetworkOrigin is VPC.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<bucket>:<access_point_name>`.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Point.
* `delete` - (Defaults to 10 mins) Used when delete the Access Point.
* `update` - (Defaults to 5 mins) Used when update the Access Point.

## Import

OSS Access Point can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_access_point.example <bucket>:<access_point_name>
```