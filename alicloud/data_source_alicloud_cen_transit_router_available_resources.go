package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCenTransitRouterAvailableResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCenTransitRouterAvailableResourcesRead,
		Schema: map[string]*schema.Schema{
			"support_multicast": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"support_multicast": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"master_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"slave_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"available_zones": {
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

func dataSourceAliCloudCenTransitRouterAvailableResourcesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTransitRouterAvailableResource"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOkExists("support_multicast"); ok {
		request["SupportMulticast"] = v
	}

	var response map[string]interface{}
	var err error

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_available_resources", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$", response)
	}

	result, _ := resp.(map[string]interface{})
	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"support_multicast": result["SupportMulticast"],
		"master_zones":      result["MasterZones"],
		"slave_zones":       result["SlaveZones"],
		"available_zones":   result["AvailableZones"],
	}

	s = append(s, mapping)

	d.SetId(tea.ToString(hashcode.String(fmt.Sprint(s))))

	if err := d.Set("resources", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
