package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudActiontrailGlobalEventsStorageRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudActiontrailGlobalEventsStorageRegionRead,
		Schema: map[string]*schema.Schema{
			"storage_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudActiontrailGlobalEventsStorageRegionRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId("GlobalEventsStorageRegion")
	client := meta.(*connectivity.AliyunClient)

	actiontrailService := ActiontrailService{client}
	object, err := actiontrailService.DescribeActiontrailGlobalEventsStorageRegion(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("storage_region", object["StorageRegion"])

	return nil
}
