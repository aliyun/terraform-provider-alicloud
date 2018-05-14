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
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_network_type": {
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
			"connection_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner_account": {
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
						"pay_type": {
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
						"db_instance_net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_cloud_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"multi_or_single": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"master_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guard_db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"readonly_db_instance_ids": {
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
						"replicate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
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
	args.InstanceNetworkType = d.Get("instance_network_type").(string)
	args.VpcId = d.Get("vpc_id").(string)
	args.VSwitchId = d.Get("vswitch_id").(string)
	args.ConnectionMode = d.Get("connection_mode").(string)
	args.Tags = d.Get("tags").(string)
	args.OwnerAccount = d.Get("owner_account").(string)
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
			"id":                       item.DBInstanceId,
			"name":                     item.DBInstanceDescription,
			"pay_type":                 item.PayType,
			"instance_type":            item.DBInstanceType,
			"region_id":                item.RegionId,
			"create_time":              item.CreateTime,
			"expire_time":              item.ExpireTime,
			"status":                   item.DBInstanceStatus,
			"engine":                   item.Engine,
			"engine_version":           item.EngineVersion,
			"db_instance_net_type":     item.DBInstanceNetType,
			"connection_mode":          item.ConnectionMode,
			"lock_mode":                item.LockMode,
			"lock_reason":              item.LockReason,
			"db_instance_class":        item.DBInstanceClass,
			"instance_network_type":    item.InstanceNetworkType,
			"vpc_cloud_instance_id":    item.VpcCloudInstanceId,
			"zone_id":                  item.ZoneId,
			"multi_or_single":          item.MutriORsignle,
			"master_instance_id":       item.MasterInstanceId,
			"guard_db_instance_id":     item.GuardDBInstanceId,
			"temp_db_instance_id":      item.TempDBInstanceId,
			"readonly_db_instance_ids": item.ReadOnlyDBInstanceIds.ReadOnlyDBInstanceId,
			"vpc_id":                   item.VpcId,
			"vswitch_id":               item.VSwitchId,
			"replicate_id":             item.ReplicateId,
			"resource_group_id":        item.ResourceGroupId,
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
