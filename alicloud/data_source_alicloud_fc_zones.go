package alicloud

import (
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"sort"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudFCZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFCZonesRead,

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
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"multi_zone_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudFCZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var clientInfo *fc.Client
	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		clientInfo = fcClient
		return fcClient.GetAccountSettings(fc.NewGetAccountSettingsInput())
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_fc_zones", "GetAccountSettings", AlibabaCloudSdkGoERROR)
	}
	addDebug("GetAccountSettings", raw, clientInfo)
	response, _ := raw.(*fc.GetAccountSettingsOutput)
	sort.Strings(response.AvailableAZs)
	return zoneIdsDescriptionAttributes(d, response.AvailableAZs)
}
