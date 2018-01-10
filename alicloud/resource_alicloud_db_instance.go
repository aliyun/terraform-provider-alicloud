package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBInstanceCreate,
		Read:   resourceAlicloudDBInstanceRead,
		Update: resourceAlicloudDBInstanceUpdate,
		Delete: resourceAlicloudDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"engine": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{"MySQL", "SQLServer", "PostgreSQL", "PPAS"}),
				ForceNew:     true,
				Required:     true,
			},
			"engine_version": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{"5.5", "5.6", "5.7", "2008r2", "2012", "9.4", "9.3"}),
				ForceNew:     true,
				Required:     true,
			},
			"db_instance_class": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'db_instance_class' has been deprecated from provider version 1.5.0. New field 'instance_type' replaces it.",
			},
			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_storage": &schema.Schema{
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'db_instance_storage' has been deprecated from provider version 1.5.0. New field 'instance_storage' replaces it.",
			},

			"instance_storage": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"instance_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(rds.Postpaid), string(rds.Prepaid)}),
				Optional:     true,
				ForceNew:     true,
				Default:      rds.Postpaid,
			},

			"period": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: rdsPostPaidDiffSuppressFunc,
			},

			"zone_id": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: zoneIdDiffSuppressFunc,
			},

			"multi_az": &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"zone_id"},
			},

			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"connection_string": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"db_instance_net_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'db_instance_net_type' has been deprecated from provider version 1.5.0.",
			},
			"allocate_public_connection": &schema.Schema{
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'allocate_public_connection' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_connection' replaces it.",
			},

			"instance_network_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'instance_network_type' has been deprecated from provider version 1.5.0.",
			},

			"master_user_name": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'master_user_name' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_account' field 'name' replaces it.",
			},

			"master_user_password": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'master_user_password' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_account' field 'password' replaces it.",
			},

			"preferred_backup_period": &schema.Schema{
				Type:       schema.TypeList,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Optional:   true,
				Deprecated: "Field 'preferred_backup_period' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_backup_policy' field 'backup_period' replaces it.",
			},

			"preferred_backup_time": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'preferred_backup_time' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_backup_policy' field 'backup_time' replaces it.",
			},

			"backup_retention_period": &schema.Schema{
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'backup_retention_period' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_backup_policy' field 'retention_period' replaces it.",
			},

			"security_ips": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},

			"connections": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_string": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_address": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'connections' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_connection' replaces it.",
			},

			"db_mappings": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"character_set_name": &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validateAllowedStringValue(rds.CHARACTER_SET_NAME),
							Required:     true,
						},
						"db_description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'db_mappings' has been deprecated from provider version 1.5.0. New resource 'alicloud_db_database' replaces it.",
			},
		},
	}
}

func resourceAlicloudDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rdsconn

	args, err := buildDBCreateOrderArgs(d, meta)
	if err != nil {
		return err
	}

	resp, err := conn.CreateOrder(args)

	if err != nil {
		return fmt.Errorf("Error creating Alicloud db instance: %#v", err)
	}

	d.SetId(resp.DBInstanceId)

	// wait instance status change from Creating to running
	if err := conn.WaitForInstanceAsyn(d.Id(), rds.Running, defaultLongTimeout); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", rds.Running, err)
	}

	return resourceAlicloudDBInstanceUpdate(d, meta)
}

func resourceAlicloudDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rdsconn
	d.Partial(true)

	if d.HasChange("security_ips") {
		ipList := expandStringList(d.Get("security_ips").([]interface{}))

		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := client.ModifyDBSecurityIps(d.Id(), ipstr); err != nil {
			return fmt.Errorf("Moodify DB security ips %s got an error: %#v", ipstr, err)
		}
		d.SetPartial("security_ips")
	}

	update := false
	args := rds.ModifyDBInstanceSpecArgs{
		DBInstanceId: d.Id(),
		PayType:      rds.Postpaid,
	}

	if d.HasChange("instance_type") && !d.IsNewResource() {
		args.DBInstanceClass = d.Get("instance_type").(string)
		update = true
		d.SetPartial("instance_type")
	}

	if d.HasChange("instance_storage") && !d.IsNewResource() {
		args.DBInstanceStorage = strconv.Itoa(d.Get("instance_storage").(int))
		update = true
		d.SetPartial("instance_storage")
	}

	if update {
		if _, err := conn.ModifyDBInstanceSpec(&args); err != nil {
			return err
		}
		//// wait instance status change from Creating to running
		//if err := conn.WaitForInstanceAsyn(d.Id(), rds.Running, defaultLongTimeout); err != nil {
		//	return fmt.Errorf("WaitForInstance %s got error: %#v", rds.Running, err)
		//}
	}

	d.Partial(false)
	return resourceAlicloudDBInstanceRead(d, meta)
}

func resourceAlicloudDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	instance, err := client.DescribeDBInstanceById(d.Id())
	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, InvalidDBInstanceNameNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err)
	}

	ips, err := client.GetSecurityIps(d.Id(), d.Get("security_ips"))
	if err != nil {
		log.Printf("Describe DB security ips error: %#v", err)
	}
	d.Set("security_ips", ips)

	d.Set("engine", instance.Engine)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("instance_type", instance.DBInstanceClass)
	d.Set("port", instance.Port)
	d.Set("instance_storage", instance.DBInstanceStorage)
	d.Set("zone_id", instance.ZoneId)
	d.Set("instance_charge_type", instance.PayType)
	d.Set("period", d.Get("period"))
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("connection_string", instance.ConnectionString)

	return nil
}

func resourceAlicloudDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	instance, err := client.DescribeDBInstanceById(d.Id())
	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, InvalidDBInstanceNameNotFound) {
			return nil
		}
		return fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err)
	}
	if instance.PayType == rds.Prepaid {
		return fmt.Errorf("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically.")
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := client.rdsconn.DeleteInstance(d.Id())

		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidDBInstanceNameNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete DB instance timeout and got an error: %#v.", err))
		}

		instance, err := client.DescribeDBInstanceById(d.Id())
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidDBInstanceNameNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err))
		}
		if instance == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Delete DB instance timeout and got an error: %#v.", err))
	})
}

func buildDBCreateOrderArgs(d *schema.ResourceData, meta interface{}) (*rds.CreateOrderArgs, error) {
	client := meta.(*AliyunClient)
	args := &rds.CreateOrderArgs{
		RegionId: getRegion(d, meta),
		// we does not expose this param to user,
		// because create prepaid instance progress will be stopped when set auto_pay to false,
		// then could not get instance info, cause timeout error
		AutoPay:           "true",
		EngineVersion:     d.Get("engine_version").(string),
		Engine:            rds.Engine(d.Get("engine").(string)),
		DBInstanceStorage: d.Get("instance_storage").(int),
		DBInstanceClass:   d.Get("instance_type").(string),
		Quantity:          DEFAULT_INSTANCE_COUNT,
		Resource:          rds.DefaultResource,
		DBInstanceNetType: common.Intranet,
	}

	bussStr, err := json.Marshal(DefaultBusinessInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to translate bussiness info %#v from json to string", DefaultBusinessInfo)
	}

	args.BusinessInfo = string(bussStr)

	if zone, ok := d.GetOk("zone_id"); ok && zone.(string) != "" {
		args.ZoneId = zone.(string)
	}

	multiAZ := d.Get("multi_az").(bool)
	if multiAZ {
		azs, err := client.DescribeMultiIZByRegion()
		if err != nil {
			return nil, fmt.Errorf("DescribeMultiIZByRegion got an error: %#v.", err)
		}

		if len(azs) < 1 {
			return nil, fmt.Errorf("Current region does not support multiple availability zones. Please change to other regions.")
		}

		args.ZoneId = azs[0]
	}

	vswitchId := d.Get("vswitch_id").(string)

	args.InstanceNetworkType = common.Classic

	if vswitchId != "" {
		args.VSwitchId = vswitchId
		args.InstanceNetworkType = common.VPC

		// check vswitchId in zone
		vsws, err := client.QueryVswitches(&ecs.DescribeVSwitchesArgs{
			RegionId:  getRegion(d, meta),
			VSwitchId: vswitchId,
		})
		if err != nil {
			return nil, fmt.Errorf("DescribeVSwitches got an error: %#v.", err)
		}

		if len(vsws) < 1 {
			return nil, fmt.Errorf("VSwitch %s is not found in the region %s.", vswitchId, getRegion(d, meta))
		}

		args.ZoneId = vsws[0].ZoneId
		args.VPCId = vsws[0].VpcId
	}

	args.PayType = rds.DBPayType(d.Get("instance_charge_type").(string))

	// if charge type is postpaid, the commodity code must set to bards
	args.CommodityCode = rds.Bards
	// At present, API supports two charge options about 'Prepaid'.
	// 'Month': valid period ranges [1-9]; 'Year': valid period range [1-3]
	// This resource only supports to input Month period [1-9, 12, 24, 36] and the values need to be converted before using them.
	if args.PayType == rds.Prepaid {
		args.CommodityCode = rds.Rds

		period := d.Get("period").(int)
		args.UsedTime = period
		args.TimeType = common.Month
		if period > 9 {
			args.UsedTime = period / 12
			args.TimeType = common.Year
		}
	}

	return args, nil
}
