package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func dataSourceAlicloudEnhancedNatAvailableZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEnhancedNatAvailableZonesRead,

		Schema: map[string]*schema.Schema{
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
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEnhancedNatAvailableZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateListEnhanhcedNatGatewayAvailableZonesRequest()
	request.RegionId = string(client.Region)
	invoker := NewInvoker()
	var err error
	var raw interface{}
	if err := invoker.Run(func() error {
		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ListEnhanhcedNatGatewayAvailableZones(request)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return err
	}); err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_enhanced_nat_available_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response := raw.(*vpc.ListEnhanhcedNatGatewayAvailableZonesResponse)
	var s []map[string]interface{}
	var ids []string
	if len(response.Zones) > 0 {
		for _, val := range response.Zones {
			mapping := map[string]interface{}{
				"zone_id":    val.ZoneId,
				"local_name": val.LocalName,
			}
			s = append(s, mapping)
			ids = append(ids, val.ZoneId)
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
