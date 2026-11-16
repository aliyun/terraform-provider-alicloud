package alicloud

// Retained FC test helpers. The migrated FC test files moved into
// alicloud/service/fc, but these shared helpers are still referenced by
// in-package root test consumers that cannot import the FC subpackage
// without a cycle (data_source_alicloud_file_crc64_checksum_test.go,
// provider_test.go, resource_alicloud_fcv2_function_test.go,
// resource_alicloud_fcv3_custom_domain_test.go and the ESA certificate
// tests). They are retained here 1:1 from the pre-migration source to keep
// those root tests compiling and behaving unchanged. The migrated copies
// live independently in alicloud/service/fc.

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/aliyun/fc-go-sdk"
)

func resourceFCFunctionConfigDependence(name string) string {
	dir, err := ioutil.TempDir(os.TempDir(), name)
	if err != nil {
		log.Printf("Failed to create temp directory: %s. Error: %v", dir, err)
		return ""
	}
	result := `'use strict';
	           exports.initializer = (context, callback) => {
		           console.log('hello init');
		           callback(null, 'hello init');
			   }
	           exports.handler = (event, context, callback) => {
		           console.log(event.toString())
		           callback(null, 'hello world');
			   }`
	filePath := filepath.Join(dir, "hello.js")
	err = ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		log.Printf("Failed to write file: %s. Error: %v", filePath, err)
		return ""
	}
	// Create the zip file.
	zipped := &bytes.Buffer{}
	err = fc.ZipDir(dir, zipped)
	if err != nil {
		return ""
	}
	zipFilePath := filepath.Join(os.TempDir(), name+".zip")
	err = ioutil.WriteFile(zipFilePath, zipped.Bytes(), 0644)
	if err != nil {
		log.Printf("Failed to write zip file: %s. Error: %v", zipFilePath, err)
		return ""
	}

	return fmt.Sprintf(`
variable "name" {
    default = "%v"
}

// Custom-container images must be in the same region as the FC function (FC v1).
// Use a data source so the image URL tracks the credentials' region rather than
// being hard-coded.
data "alicloud_regions" "current" {
  current = true
}

// After serveral hours of investigation, finally figure out how to escape the double quotes.
// https://github.com/hashicorp/terraform/issues/17144
// https://discuss.hashicorp.com/t/how-can-i-escape-double-quotes-in-a-variable-value/4697
locals {
	container_command = "[\"python\", \"server.py\"]"
	container_args = "[\"a1\", \"a2\"]"
}

output "container_command" {
	value = "${local.container_command}"
}

output "container_args" {
	value = "${local.container_args}"
}

resource "alicloud_log_project" "default" {
  name = "${var.name}"
  description = "tf unit test"
  lifecycle {
    ignore_changes = [
      policy
    ]
  }
}

resource "alicloud_log_store" "default" {
  project = "${alicloud_log_project.default.name}"
  name = "${var.name}"
  retention_period = "3000"
  shard_count = 1
}
resource "alicloud_fc_service" "default" {
    name = "${var.name}"
    description = "tf unit test"
    log_config {
	project = "${alicloud_log_project.default.name}"
	logstore = "${alicloud_log_store.default.name}"
    }
    role = "${alicloud_ram_role.default.arn}"
    depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "default" {
  bucket = "${alicloud_oss_bucket.default.id}"
  key = "fc/hello.zip"
  source = "%s"
}

resource "alicloud_ram_role" "default" {
  name = "${var.name}"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}
resource "alicloud_ram_role_policy_attachment" "default" {
  role_name = "${alicloud_ram_role.default.name}"
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}
resource "alicloud_ram_role_policy_attachment" "acr" {
  role_name = "${alicloud_ram_role.default.name}"
  policy_name = "AliyunContainerRegistryReadOnlyAccess"
  policy_type = "System"
}
`, name, zipFilePath, testFCRoleTemplate)
}

var testFCRoleTemplate = `
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "fc.aliyuncs.com"
        ]
      }
    }
  ],
  "Version": "1"
}
`
