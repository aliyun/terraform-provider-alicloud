package alicloud

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudNasProtocols() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNasProtocolsRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Capacity",
					"Performance",
				}),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAlicloudNasProtocolsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := nas.CreateDescribeZonesRequest()
	request.RegionId = client.RegionId
	var nasProtocol [][]string
	for {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeZones(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_protocols", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*nas.DescribeZonesResponse)
		if len(response.Zones.Zone) < 1 {
			break
		}
		for _, val := range response.Zones.Zone {
			if v, ok := d.GetOkExists("zone_id"); ok && val.ZoneId != Trim(v.(string)) {
				continue
			}
			if v, ok := d.GetOkExists("type"); ok {
				if Trim(v.(string)) == "Performance" {
					if len(val.Performance.Protocol) == 0 {
						continue
					} else {
						nasProtocol = append(nasProtocol, val.Performance.Protocol)
					}
				}
				if Trim(v.(string)) == "Capacity" {
					if len(val.Capacity.Protocol) == 0 {
						continue
					} else {
						nasProtocol = append(nasProtocol, val.Capacity.Protocol)
					}
				}
			}
		}
		break
	}

	return nasProtocolsDescriptionAttributes(d, nasProtocol)
}

func nasProtocolsDescriptionAttributes(d *schema.ResourceData, nasProtocol [][]string) error {
	var s []string
	var ids []string
	for _, val := range nasProtocol {
		for _, protocol := range val {
			s = append(s, strings.ToUpper(protocol))
			ids = append(ids, protocol)
		}
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("protocols", s); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOkExists("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
