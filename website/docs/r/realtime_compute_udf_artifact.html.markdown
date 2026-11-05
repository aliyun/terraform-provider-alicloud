---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_udf_artifact"
description: |-
  Provides a Alicloud Realtime Compute Udf Artifact resource.
---

# alicloud_realtime_compute_udf_artifact

Provides a Realtime Compute Udf Artifact resource.

A UDF (User-Defined Function) artifact registers a Java jar or Python file stored in OSS
with a Realtime Compute for Apache Flink namespace. The service derives the UDF classes
from the uploaded artifact so they can be referenced by deployments and SQL jobs.

For information about Realtime Compute Udf Artifact and how to use it, see
[What is Udf Artifact](https://next.api.alibabacloud.com/document/ververica/2022-07-18/CreateUdfArtifact).

-> **NOTE:** Available since v2.3.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "tfexample-udf-artifact"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_oss_bucket" "default" {
}

resource "alicloud_oss_bucket_object" "udf" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "udf/UpperUdf.py"
  content = <<-EOF
    class UpperUdf:
        def eval(self, s):
            if s is None:
                return None
            return str(s).upper()
  EOF
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.id
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = [alicloud_vswitch.default.id]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id      = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_udf_artifact" "default" {
  workspace     = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace     = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  name          = var.name
  artifact_type = "ARTIFACT_TYPE_PYTHON"
  jar_url       = "oss://${alicloud_oss_bucket.default.id}/udf/UpperUdf.py"
}
```

## Argument Reference

The following arguments are supported:
* `artifact_type` - (Required, String) The artifact type. Valid values: `ARTIFACT_TYPE_JAVA`, `ARTIFACT_TYPE_PYTHON`.
* `dependency_jar_uris` - (Optional, List) The list of additional dependency jar storage paths, in the `oss://<bucket>/<key>` format.
* `jar_url` - (Required, String) The OSS storage path of the UDF artifact, in the `oss://<bucket>/<key>` format. The bucket must be the one bound to the Realtime Compute instance so the Flink service can read the artifact.
* `name` - (Required, ForceNew, String) The UDF artifact name (the jar or python file name). Changing this creates a new resource.
* `namespace` - (Required, ForceNew, String) The project namespace name. Changing this creates a new resource.
* `workspace` - (Required, ForceNew, String) The workspace ID of the Realtime Compute (Ververica) instance. Changing this creates a new resource.

## Attributes Reference

The following attributes are exported:
* `id` - (String) The ID of the resource. The value is formulated as `<workspace>:<namespace>:<udfArtifactName>`.
* `created_at` - (String) The creation time of the artifact (Unix seconds).
* `creator` - (String) The creator of the artifact.
* `modified_at` - (String) The last modification time of the artifact (Unix seconds).
* `udf_classes` - (List) The UDF classes the service derives from the uploaded artifact.
  * `class_name` - (String) The class name.
  * `class_type` - (String) The class type.
  * `function_names` - (List) The function names of the class.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Udf Artifact.
* `delete` - (Defaults to 5 mins) Used when delete the Udf Artifact.
* `update` - (Defaults to 5 mins) Used when update the Udf Artifact.

## Import

Realtime Compute Udf Artifact can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_udf_artifact.example <workspace>:<namespace>:<udfArtifactName>
```
