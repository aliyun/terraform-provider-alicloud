package alicloud

import (
	"log"
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudSlbs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
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
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: validateDataSourceSlbsTagsNum,
				ForceNew:     true,
			},
			"output_file": {
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

	args := slb.CreateDescribeLoadBalancersRequest()

	if v, ok := d.GetOk("master_availability_zone"); ok && v.(string) != "" {
		args.MasterZoneId = v.(string)
	}
	if v, ok := d.GetOk("slave_availability_zone"); ok && v.(string) != "" {
		args.SlaveZoneId = v.(string)
	}
	if v, ok := d.GetOk("network_type"); ok && v.(string) != "" {
		args.NetworkType = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		args.VpcId = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		args.VSwitchId = v.(string)
	}
	if v, ok := d.GetOk("address"); ok && v.(string) != "" {
		args.Address = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		var tags []Tag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, Tag{
				Key:   key,
				Value: value.(string),
			})
		}
		args.Tags = toSlbTagsString(tags)
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allLoadBalancers []slb.LoadBalancer
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeLoadBalancers(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*slb.DescribeLoadBalancersResponse)
		if resp == nil || len(resp.LoadBalancers.LoadBalancer) < 1 {
			break
		}

		allLoadBalancers = append(allLoadBalancers, resp.LoadBalancers.LoadBalancer...)

		if len(resp.LoadBalancers.LoadBalancer) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
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
	var s []map[string]interface{}
	for _, loadBalancer := range loadBalancers {
		tags, _ := slbService.describeTags(loadBalancer.LoadBalancerId)
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
			"tags":                     slbService.slbTagsToMap(tags),
		}

		log.Printf("[DEBUG] alicloud_slbs - adding slb mapping: %v", mapping)
		ids = append(ids, loadBalancer.LoadBalancerId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slbs", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
