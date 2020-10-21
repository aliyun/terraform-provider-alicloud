package alicloud

import (
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOnsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOnsInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 2, 5, 7}),
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_internal_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http_internet_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http_internet_secure_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
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
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"release_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"tcp_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudOnsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ons.CreateOnsInstanceInServiceListRequest()
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]ons.OnsInstanceInServiceListTag, len(v.(map[string]interface{})))
		i := 0
		for key, value := range v.(map[string]interface{}) {
			tags[i] = ons.OnsInstanceInServiceListTag{
				Key:   key,
				Value: value.(string),
			}
			i++
		}
		request.Tag = &tags
	}
	var objects []ons.InstanceVO
	var instanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		instanceNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOkExists("status")
	var response *ons.OnsInstanceInServiceListResponse
	raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsInstanceInServiceList(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ons_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*ons.OnsInstanceInServiceListResponse)

	for _, item := range response.Data.InstanceVO {
		if instanceNameRegex != nil {
			if !instanceNameRegex.MatchString(item.InstanceName) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[item.InstanceId]; !ok {
				continue
			}
		}
		if statusOk && status != item.InstanceStatus {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"independent_naming": object.IndependentNaming,
			"id":                 object.InstanceId,
			"instance_id":        object.InstanceId,
			"instance_name":      object.InstanceName,
			"instance_type":      object.InstanceType,
			"release_time":       time.Unix(int64(object.ReleaseTime)/1000, 0).Format("2006-01-02 03:04:05"),
			"instance_status":    object.InstanceStatus,
			"status":             object.InstanceStatus,
		}
		ids = append(ids, object.InstanceId)
		tags := make(map[string]string)
		for _, t := range object.Tags.Tag {
			tags[t.Key] = t.Value
		}
		mapping["tags"] = tags
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			names = append(names, object.InstanceName)
			s = append(s, mapping)
			continue
		}

		request := ons.CreateOnsInstanceBaseInfoRequest()
		request.RegionId = client.RegionId
		request.InstanceId = object.InstanceId
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsInstanceBaseInfo(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ons_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		responseGet, _ := raw.(*ons.OnsInstanceBaseInfoResponse)
		mapping["http_internal_endpoint"] = responseGet.InstanceBaseInfo.Endpoints.HttpInternalEndpoint
		mapping["http_internet_endpoint"] = responseGet.InstanceBaseInfo.Endpoints.HttpInternetEndpoint
		mapping["http_internet_secure_endpoint"] = responseGet.InstanceBaseInfo.Endpoints.HttpInternetSecureEndpoint
		mapping["remark"] = responseGet.InstanceBaseInfo.Remark
		mapping["tcp_endpoint"] = responseGet.InstanceBaseInfo.Endpoints.TcpEndpoint
		names = append(names, object.InstanceName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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
