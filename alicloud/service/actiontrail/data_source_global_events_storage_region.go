package actiontrail

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGlobalEventsStorageRegion() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceGlobalEventsStorageRegionRead,
		Schema: map[string]*schema.Schema{
			"storage_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceGlobalEventsStorageRegionRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId("GlobalEventsStorageRegion")
	client := meta.(*connectivity.AliyunClient)

	actiontrailService := NewActiontrailService(client)
	object, err := actiontrailService.DescribeActiontrailGlobalEventsStorageRegion(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("storage_region", object["StorageRegion"])

	return nil
}
