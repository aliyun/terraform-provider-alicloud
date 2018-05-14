package alicloud

import (
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudRdsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(MySQL),
					string(SQLServer),
					string(PPAS),
					string(PostgreSQL),
				}),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				// please follow the link below to see more details on available statusesplease follow the link below to see more details on available statuses
				// 实例状态表
				// https://help.aliyun.com/document_detail/26315.html
				ValidateFunc: validateAllowedStringValue([]string{
					"Creating",
					"Running",
					"Deleting",
					"Rebooting",
					"DBInstanceClassChanging",
					"TRANSING",
					"EngineVersionUpgrading",
					"TransingToOthers",
					"GuardDBInstanceCreating",
					"Restoring",
					"Importing",
					"ImportingFromOthers",
					"DBInstanceNetTypeChanging",
					"GuardSwitching",
					"INS_CLONING",
				}),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Primary",
					"Readonly",
					"Guard",
					"Temp",
				}),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Standard",
					"Safe",
				}),
			},
			"tags": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateJsonString,
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
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_mode": {
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
						"master_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guard_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"readonly_instance_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRdsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).rdsconn

	args := rds.CreateDescribeDBInstancesRequest()

	args.RegionId = string(getRegion(d, meta))
	args.Engine = d.Get("engine").(string)
	args.DBInstanceStatus = d.Get("status").(string)
	args.DBInstanceType = d.Get("instance_type").(string)
	args.VpcId = d.Get("vpc_id").(string)
	args.VSwitchId = d.Get("vswitch_id").(string)
	args.ConnectionMode = d.Get("connection_mode").(string)
	args.Tags = d.Get("tags").(string)
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var dbi []rds.DBInstance

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	for {
		resp, err := conn.DescribeDBInstances(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.Items.DBInstance) < 1 {
			break
		}

		for _, item := range resp.Items.DBInstance {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.DBInstanceDescription) {
					continue
				}
			}
			dbi = append(dbi, item)
		}

		if len(resp.Items.DBInstance) < PageSizeLarge {
			break
		}

		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}

	return rdsInstancesDescription(d, dbi)
}

func rdsInstancesDescription(d *schema.ResourceData, dbi []rds.DBInstance) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range dbi {
		mapping := map[string]interface{}{
			"id":                    item.DBInstanceId,
			"name":                  item.DBInstanceDescription,
			"charge_type":           item.PayType,
			"instance_type":         item.DBInstanceType,
			"region_id":             item.RegionId,
			"create_time":           item.CreateTime,
			"expire_time":           item.ExpireTime,
			"status":                item.DBInstanceStatus,
			"engine":                item.Engine,
			"engine_version":        item.EngineVersion,
			"net_type":              item.DBInstanceNetType,
			"connection_mode":       item.ConnectionMode,
			"instance_class":        item.DBInstanceClass,
			"availability_zone":     item.ZoneId,
			"master_instance_id":    item.MasterInstanceId,
			"guard_instance_id":     item.GuardDBInstanceId,
			"temp_instance_id":      item.TempDBInstanceId,
			"readonly_instance_ids": item.ReadOnlyDBInstanceIds.ReadOnlyDBInstanceId,
			"vpc_id":                item.VpcId,
			"vswitch_id":            item.VSwitchId,
		}

		log.Printf("alicloud_rds_instances - adding rds instance: %v", mapping)
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
