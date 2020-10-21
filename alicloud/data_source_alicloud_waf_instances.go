package alicloud

import (
	"strconv"

	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudWafInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudWafInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_source": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_date": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"in_debt": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remain_day": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"subscription_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trial": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudWafInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := waf_openapi.CreateDescribeInstanceInfosRequest()
	if v, ok := d.GetOk("instance_source"); ok {
		request.InstanceSource = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	var objects []waf_openapi.InstanceInfo

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response *waf_openapi.DescribeInstanceInfosResponse
	raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
		return waf_openapiClient.DescribeInstanceInfos(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_waf_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*waf_openapi.DescribeInstanceInfosResponse)

	for _, item := range response.InstanceInfos {
		if len(idsMap) > 0 {
			if _, ok := idsMap[item.InstanceId]; !ok {
				continue
			}
		}
		if statusOk && status != "" && status != strconv.Itoa(item.Status) {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"end_date":          object.EndDate,
			"in_debt":           object.InDebt,
			"id":                object.InstanceId,
			"instance_id":       object.InstanceId,
			"remain_day":        object.RemainDay,
			"status":            object.Status,
			"subscription_type": object.SubscriptionType,
			"trial":             object.Trial,
		}
		ids = append(ids, object.InstanceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
