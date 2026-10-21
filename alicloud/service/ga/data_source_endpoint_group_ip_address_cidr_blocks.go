package ga

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEndpointGroupIpAddressCidrBlocks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEndpointGroupIpAddressCidrBlocksRead,
		Schema: map[string]*schema.Schema{
			"endpoint_group_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"accelerator_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourceEndpointGroupIpAddressCidrBlocksRead(d *schema.ResourceData, meta interface{}) error {
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ga_endpoint_group_ip_address_cidr_blocks", action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$", response)
	}

	result, _ := resp.(map[string]interface{})
	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"endpoint_group_region":  fmt.Sprint(result["EndpointGroupRegion"]),
		"ip_address_cidr_blocks": result["IpAddressCidrBlocks"],
		"status":                 result["State"],
	}

	s = append(s, mapping)

	d.SetId(tea.ToString(helper.Hashcode(fmt.Sprint(request["EndpointGroupRegion"]))))

	if err := d.Set("endpoint_group_ip_address_cidr_blocks", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
