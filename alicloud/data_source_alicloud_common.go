package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataResourceIdHash moved to package helper (DataResourceIdHash) so service
// packages can reach it without importing the monolith.
var dataResourceIdHash = helper.DataResourceIdHash

func outputInstancesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"image_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vswitch_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func outputShortVpcsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"region_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vpc_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vpc_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
