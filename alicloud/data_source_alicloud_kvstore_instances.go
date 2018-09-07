package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudKVStoreInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKVStoreInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(Normal),
					string(Creating),
					string(Changing),
					string(Inactive),
				}),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Memcache",
					"Redis",
				}),
			},
			"instance_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connection_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudKVStoreInstancesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).rkvconn

	args := r_kvstore.CreateDescribeInstancesRequest()
	args.RegionId = getRegionId(d, meta)
	args.VpcId = d.Get("vpc_id").(string)
	args.VSwitchId = d.Get("vswitch_id").(string)
	args.InstanceType = d.Get("instance_type").(string)
	args.InstanceStatus = d.Get("status").(string)
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var dbi []r_kvstore.KVStoreInstance

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	for {
		resp, err := conn.DescribeInstances(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.Instances.KVStoreInstance) < 1 {
			break
		}

		for _, item := range resp.Instances.KVStoreInstance {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.InstanceName) {
					continue
				}
			}
			dbi = append(dbi, item)
		}

		if len(resp.Instances.KVStoreInstance) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	return kvstoreInstancesDescription(d, dbi)
}

func kvstoreInstancesDescription(d *schema.ResourceData, dbi []r_kvstore.KVStoreInstance) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range dbi {
		mapping := map[string]interface{}{
			"id":                item.InstanceId,
			"name":              item.InstanceName,
			"charge_type":       item.ChargeType,
			"instance_type":     item.InstanceType,
			"instance_class":    item.InstanceClass,
			"region_id":         item.RegionId,
			"create_time":       item.CreateTime,
			"expire_time":       item.EndTime,
			"status":            item.InstanceStatus,
			"availability_zone": item.ZoneId,
			"vpc_id":            item.VpcId,
			"vswitch_id":        item.VSwitchId,
			"private_ip":        item.PrivateIp,
			"port":              item.Port,
			"user_name":         item.UserName,
			"capacity":          item.Bandwidth,
			"bandwidth":         item.Bandwidth,
			"connections":       item.Connections,
			"connection_domain": item.ConnectionDomain,
		}

		ids = append(ids, item.InstanceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
