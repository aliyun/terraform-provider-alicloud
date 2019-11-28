package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudGpdbZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGpdbZonesRead,

		Schema: map[string]*schema.Schema{
			"multi": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      common.PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAlicloudGpdbZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := gpdb.CreateDescribeRegionsRequest()
	request.RegionId = client.RegionId
	raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
		return gpdbClient.DescribeRegions(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_gpdb_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*gpdb.DescribeRegionsResponse)
	if len(response.Regions.Region) <= 0 {
		return WrapError(fmt.Errorf("[ERROR] There is no available region for gpdb."))
	}

	var zoneIds []string

	for _, r := range response.Regions.Region {
		for _, zone := range r.Zones.Zone {
			if d.Get("multi").(bool) {
				if strings.Contains(zone.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
					zoneIds = append(zoneIds, zone.ZoneId)
				}
			} else {
				if !strings.Contains(zone.ZoneId, MULTI_IZ_SYMBOL) {
					zoneIds = append(zoneIds, zone.ZoneId)
				}
			}
		}
	}

	d.SetId(dataResourceIdHash(zoneIds))

	if err := d.Set("ids", zoneIds); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), zoneIds)
	}

	return nil
}
