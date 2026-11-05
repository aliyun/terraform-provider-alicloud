---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_package"
description: |-
  Provides a Alicloud Realtime Compute Package resource.
---

# alicloud_realtime_compute_package

Provides a Realtime Compute Package resource.

A package is a reusable artifact (JAR, archive, script, or Python archive) stored in OSS that a Realtime Compute for Apache Flink (FOAS) project references from its jobs and deployments. The package belongs to a FOAS project and is uniquely identified within the project by its name.

For information about Realtime Compute Package and how to use it, see [What is Package](https://next.api.alibabacloud.com/document/foasconsole/2018-11-11/CreatePackage).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The package lifecycle APIs are exposed on the legacy `foas` `2018-11-11` ROA version and operate on a FOAS project (not a VVP workspace/namespace). The FOAS project must already exist and is referenced by `project_name`; Terraform does not manage the project itself.

-> **NOTE:** `project_name`, `package_name` and `type` cannot be changed without recreating the resource. `terraform destroy` deletes the package from the project; deleting a package that is still referenced by a running job or deployment may fail, so remove the references first.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_realtime_compute_package&exampleId=realtime-compute-package-basic&activeTab=example&spm=docs.r.realtime_compute_package.0.basic&intl_lang=EN_US" target="_blank">
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

variable "name" {
  default = "terraform-example"
}

variable "project_name" {
  default = "my-foas-project"
}

data "alicloud_account" "current" {}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "packages/example.jar"
  content = "dummy-jar-content"
}

resource "alicloud_realtime_compute_package" "default" {
  project_name = var.project_name
  package_name = var.name
  oss_bucket   = alicloud_oss_bucket.default.id
  oss_endpoint = "oss-cn-hangzhou.aliyuncs.com"
  oss_owner    = data.alicloud_account.current.id
  oss_path     = alicloud_oss_bucket_object.default.key
  type         = "JAR"
  description  = "An example JAR package for Realtime Compute."
}
```


## Argument Reference

The following arguments are supported:

* `description` - (Optional, String) The description of the package.
* `md5` - (Optional, String) The MD5 checksum of the package artifact. When provided, the service validates that the uploaded OSS object matches this checksum.
* `origin_name` - (Optional, String) The alias (origin name) of the package artifact.
* `oss_bucket` - (Required, String) The OSS bucket that stores the package artifact.
* `oss_endpoint` - (Required, String) The OSS endpoint used to access the bucket that stores the package artifact, for example `oss-cn-hangzhou.aliyuncs.com`.
* `oss_owner` - (Required, String) The owner (UID) of the OSS bucket that stores the package artifact. This is the Alibaba Cloud account ID that owns the bucket.
* `oss_path` - (Required, String) The OSS object path (key) of the package artifact.
* `package_name` - (Required, ForceNew, String) The name of the package, unique within the project. Part of the composite resource identity; changing it recreates the resource.
* `project_name` - (Required, ForceNew, String) The name of the FOAS project that owns the package. The project must already exist. Part of the composite resource identity; changing it recreates the resource.
* `tag` - (Optional, String) The tag of the package.
* `type` - (Required, ForceNew, String) The package type. Valid values: `JAR`, `DICTIONARY`, `SCRIPT`, `PYTHON`. Changing it recreates the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, the composite key `<project_name>:<package_name>`.
* `create_time` - (Computed, Int) The creation time of the package.
* `creator` - (Computed, String) The creator of the package.
* `modify_time` - (Computed, Int) The latest modification time of the package.
* `modifier` - (Computed, String) The latest modifier of the package.

## Import

Realtime Compute Package can be imported using the composite id `<project_name>:<package_name>`, e.g.

```shell
$ terraform import alicloud_realtime_compute_package.example my-foas-project:example-package
```
