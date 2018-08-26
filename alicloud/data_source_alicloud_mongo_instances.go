package alicloud

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudMongoInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMongoInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"sharding",
					"replicate",
				}),
			},
			"instance_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replication": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mongos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"class": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"shards": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"class": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"storage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMongoInstancesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ddsconn

	args := dds.CreateDescribeDBInstancesRequest()

	args.RegionId = getRegionId(d, meta)
	args.DBInstanceType = d.Get("instance_type").(string)
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	var dbi []dds.DBInstance

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	var instClass string
	if v, ok := d.GetOk("instance_class"); ok {
		instClass = strings.ToLower(v.(string))
	}

	var az string
	if v, ok := d.GetOk("availability_zone"); ok {
		az = strings.ToLower(v.(string))
	}

	for {
		resp, err := conn.DescribeDBInstances(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.DBInstances.DBInstance) < 1 {
			break
		}

		for _, item := range resp.DBInstances.DBInstance {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.DBInstanceDescription) {
					continue
				}
			}

			if len(instClass) > 0 && instClass != strings.ToLower(string(item.DBInstanceClass)) {
				continue
			}

			if len(az) > 0 && az != strings.ToLower(string(item.ZoneId)) {
				continue
			}

			dbi = append(dbi, item)
		}

		if len(resp.DBInstances.DBInstance) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}

	}

	return mongoInstancesDescription(d, dbi)
}

func formatRFC3339(layout string, input string, location *time.Location) (output string) {
	t, err := time.ParseInLocation(layout, input, location)
	if err != nil {
		log.Printf("[ERROR] formatRFC3339 got error: %#v", err)
		return string("n/a")
	}
	return t.UTC().Format(time.RFC3339)
}

func mongoInstancesDescription(d *schema.ResourceData, dbi []dds.DBInstance) error {
	var ids []string
	var s []map[string]interface{}
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)

	for _, item := range dbi {
		mapping := map[string]interface{}{
			"id":                item.DBInstanceId,
			"name":              item.DBInstanceDescription,
			"charge_type":       item.ChargeType,
			"instance_type":     item.DBInstanceType,
			"region_id":         item.RegionId,
			"creation_time":     formatRFC3339("2006-01-02 15:04:05.0", item.CreationTime, beijing),
			"expiration_time":   formatRFC3339("2006-01-02T15:04Z", item.ExpireTime, time.UTC),
			"status":            item.DBInstanceStatus,
			"replication":       item.ReplicationFactor,
			"engine":            item.Engine,
			"engine_version":    item.EngineVersion,
			"network_type":      item.NetworkType,
			"lock_mode":         item.LockMode,
			"instance_class":    item.DBInstanceClass,
			"storage":           item.DBInstanceStorage,
			"mongos":            item.MongosList.MongosAttribute,
			"shards":            item.ShardList.ShardAttribute,
			"availability_zone": item.ZoneId,
		}

		ids = append(ids, item.DBInstanceId)
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
