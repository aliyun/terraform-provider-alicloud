package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"sort"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudSlbZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbZonesRead,

		Schema: map[string]*schema.Schema{
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
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
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_zone_ids": {
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

func dataSourceAlicloudSlbZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	if ipVersion, ok := d.GetOk("available_slb_address_ip_version"); ok {
		request.AddressIPVersion = ipVersion.(string)
	}
	if addressType, ok := d.GetOk("available_slb_address_type"); ok {
		request.AddressType = addressType.(string)
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAvailableResource(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeAvailableResourceResponse)

	if len(response.AvailableResources.AvailableResource) <= 0 {
		return WrapError(fmt.Errorf("[ERROR] There is no available region for slb."))
	}
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	zones := make(map[string][]string)

	for _, resource := range response.AvailableResources.AvailableResource {
		slaveIds := zones[resource.MasterZoneId]
		slaveIds = append(slaveIds, resource.SlaveZoneId)
		zones[resource.MasterZoneId] = slaveIds
	}
	var ids []string
	var zoneMaps []interface{}
	for masterId, slaveIds := range zones {
		if len(slaveIds) > 0 {
			sort.Strings(slaveIds)
		}
		zone := map[string]interface{}{
			"zone_id":        masterId,
			"slave_zone_ids": slaveIds,
		}
		ids = append(ids, masterId)
		zoneMaps = append(zoneMaps, zone)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("zones", zoneMaps); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), zoneMaps)
	}

	return nil
}
