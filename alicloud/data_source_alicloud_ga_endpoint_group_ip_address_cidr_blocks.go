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

func dataSourceAliCloudGaEndpointGroupIpAddressCidrBlocks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudGaEndpointGroupIpAddressCidrBlocksRead,
		Schema: map[string]*schema.Schema{
			"endpoint_group_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accelerator_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_group_ip_address_cidr_blocks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_group_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address_cidr_blocks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudGaEndpointGroupIpAddressCidrBlocksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListEndpointGroupIpAddressCidrBlocks"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["EndpointGroupRegion"] = d.Get("endpoint_group_region")

	if v, ok := d.GetOk("accelerator_id"); ok {
		request["AcceleratorId"] = v
	}

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_endpoint_group_ip_address_cidr_blocks", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$", response)
	}

	result, _ := resp.(map[string]interface{})
	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"endpoint_group_region":  fmt.Sprint(result["EndpointGroupRegion"]),
		"ip_address_cidr_blocks": result["IpAddressCidrBlocks"],
		"status":                 result["State"],
	}

	s = append(s, mapping)

	d.SetId(tea.ToString(hashcode.String(fmt.Sprint(request["EndpointGroupRegion"]))))

	if err := d.Set("endpoint_group_ip_address_cidr_blocks", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
