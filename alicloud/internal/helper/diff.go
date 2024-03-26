package helper

import "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

func SetAttribute(d *terraform.InstanceDiff, key string, attr *terraform.ResourceAttrDiff) {
	d.Attributes[key] = attr
}

func DelAttribute(d *terraform.InstanceDiff, key string) {
	delete(d.Attributes, key)
}
