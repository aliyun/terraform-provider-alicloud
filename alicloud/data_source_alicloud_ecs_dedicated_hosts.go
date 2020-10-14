package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcsDedicatedHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsDedicatedHostsRead,
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
			"dedicated_host_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"dedicated_host_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dedicated_host_type": {
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
				ValidateFunc: validation.StringInSlice([]string{"Available", "Creating", "PermanentFailure", "Released", "UnderAssessment"}, false),
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_on_maintenance": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_placement": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_release_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"machine_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_gpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sale_cycle": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sockets": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supported_instance_types_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcsDedicatedHostsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateDescribeDedicatedHostsRequest()
	if v, ok := d.GetOk("dedicated_host_id"); ok {
		request.DedicatedHostIds = v.(string)
	}
	if v, ok := d.GetOk("dedicated_host_name"); ok {
		request.DedicatedHostName = v.(string)
	}
	if v, ok := d.GetOk("dedicated_host_type"); ok {
		request.DedicatedHostType = v.(string)
	}
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		request.Status = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]ecs.DescribeDedicatedHostsTag, len(v.(map[string]interface{})))
		i := 0
		for key, value := range v.(map[string]interface{}) {
			tags[i] = ecs.DescribeDedicatedHostsTag{
				Key:   key,
				Value: value.(string),
			}
			i++
		}
		request.Tag = &tags
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []ecs.DedicatedHost
	var dedicatedHostNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		dedicatedHostNameRegex = r
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
	var response *ecs.DescribeDedicatedHostsResponse
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDedicatedHosts(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_dedicated_hosts", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*ecs.DescribeDedicatedHostsResponse)

		for _, item := range response.DedicatedHosts.DedicatedHost {
			if dedicatedHostNameRegex != nil {
				if !dedicatedHostNameRegex.MatchString(item.DedicatedHostName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.DedicatedHostId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.DedicatedHosts.DedicatedHost) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, len(objects))
	names := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"action_on_maintenance":         object.ActionOnMaintenance,
			"auto_placement":                object.AutoPlacement,
			"auto_release_time":             object.AutoReleaseTime,
			"cores":                         object.Cores,
			"id":                            object.DedicatedHostId,
			"dedicated_host_id":             object.DedicatedHostId,
			"dedicated_host_name":           object.DedicatedHostName,
			"dedicated_host_type":           object.DedicatedHostType,
			"description":                   object.Description,
			"expired_time":                  object.ExpiredTime,
			"gpu_spec":                      object.GPUSpec,
			"machine_id":                    object.MachineId,
			"payment_type":                  object.ChargeType,
			"physical_gpus":                 object.PhysicalGpus,
			"resource_group_id":             object.ResourceGroupId,
			"sale_cycle":                    object.SaleCycle,
			"sockets":                       object.Sockets,
			"status":                        object.Status,
			"supported_instance_types_list": object.SupportedInstanceTypesList.SupportedInstanceTypesList,
			"zone_id":                       object.ZoneId,
		}
		ids[i] = object.DedicatedHostId
		tags := make(map[string]string)
		for _, t := range object.Tags.Tag {
			tags[t.TagKey] = t.TagValue
		}
		mapping["tags"] = tags
		names[i] = object.DedicatedHostName
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("hosts", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
