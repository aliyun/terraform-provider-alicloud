package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSlbs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"master_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"slave_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"slbs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := &SlbService{client}

	request := slb.CreateDescribeLoadBalancersRequest()
	request.RegionId = client.RegionId
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	if v, ok := d.GetOk("master_availability_zone"); ok && v.(string) != "" {
		request.MasterZoneId = v.(string)
	}
	if v, ok := d.GetOk("slave_availability_zone"); ok && v.(string) != "" {
		request.SlaveZoneId = v.(string)
	}
	if v, ok := d.GetOk("network_type"); ok && v.(string) != "" {
		request.NetworkType = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request.VpcId = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.VSwitchId = v.(string)
	}
	if v, ok := d.GetOk("address"); ok && v.(string) != "" {
		request.Address = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		var tags []Tag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, Tag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tags = toSlbTagsString(tags)
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allLoadBalancers []slb.LoadBalancer
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeLoadBalancers(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_slbs", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*slb.DescribeLoadBalancersResponse)
		if len(response.LoadBalancers.LoadBalancer) < 1 {
			break
		}

		allLoadBalancers = append(allLoadBalancers, response.LoadBalancers.LoadBalancer...)

		if len(response.LoadBalancers.LoadBalancer) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredLoadBalancersTemp []slb.LoadBalancer

	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, balancer := range allLoadBalancers {
			if r != nil && !r.MatchString(balancer.LoadBalancerName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[balancer.LoadBalancerId]; !ok {
					continue
				}
			}

			filteredLoadBalancersTemp = append(filteredLoadBalancersTemp, balancer)
		}
	} else {
		filteredLoadBalancersTemp = allLoadBalancers
	}

	return slbsDescriptionAttributes(d, filteredLoadBalancersTemp, slbService)
}

func slbsDescriptionAttributes(d *schema.ResourceData, loadBalancers []slb.LoadBalancer, slbService *SlbService) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, loadBalancer := range loadBalancers {
		tags, _ := slbService.DescribeTags(loadBalancer.LoadBalancerId, nil, TagResourceInstance)
		mapping := map[string]interface{}{
			"id":                       loadBalancer.LoadBalancerId,
			"region_id":                loadBalancer.RegionId,
			"master_availability_zone": loadBalancer.MasterZoneId,
			"slave_availability_zone":  loadBalancer.SlaveZoneId,
			"status":                   loadBalancer.LoadBalancerStatus,
			"name":                     loadBalancer.LoadBalancerName,
			"network_type":             loadBalancer.NetworkType,
			"vpc_id":                   loadBalancer.VpcId,
			"vswitch_id":               loadBalancer.VSwitchId,
			"address":                  loadBalancer.Address,
			"internet":                 loadBalancer.AddressType == strings.ToLower(string(Internet)),
			"creation_time":            loadBalancer.CreateTime,
			"tags":                     slbService.tagsToMap(tags),
		}

		ids = append(ids, loadBalancer.LoadBalancerId)
		names = append(names, loadBalancer.LoadBalancerName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slbs", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
