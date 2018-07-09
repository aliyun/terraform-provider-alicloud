package alicloud

import (
	"fmt"
	"strings"
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/go-uuid"
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
				ValidateFunc: validateAllowedStringValue([]string{string(MySQL), string(SQLServer), string(PostgreSQL), string(PPAS)}),
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
				ValidateFunc: validateAllowedStringValue([]string{string(Postpaid), string(Prepaid)}),
				Optional:     true,
				ForceNew:     true,
				Default:      Postpaid,
			},

			"period": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: rdsPostPaidDiffSuppressFunc,
			},

			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"multi_az": &schema.Schema{
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'multi_az' has been deprecated from provider version 1.8.1. Please use field 'zone_id' to specify multiple availability zone.",
			},

			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"instance_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBInstanceName,
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
				Type:     schema.TypeSet,
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
							Type:     schema.TypeString,
							Required: true,
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

	request, err := buildDBCreateRequest(d, meta)
	if err != nil {
		return err
	}

	resp, err := conn.CreateDBInstance(request)

	if err != nil {
		return fmt.Errorf("Error creating Alicloud db instance: %#v", err)
	}

	d.SetId(resp.DBInstanceId)

	// wait instance status change from Creating to running
	if err := client.WaitForDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}

	return resourceAlicloudDBInstanceUpdate(d, meta)
}

func resourceAlicloudDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rdsconn
	d.Partial(true)

	if d.HasChange("security_ips") && !d.IsNewResource() {
		ipList := expandStringList(d.Get("security_ips").(*schema.Set).List())

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
	request := rds.CreateModifyDBInstanceSpecRequest()
	request.DBInstanceId = d.Id()
	request.PayType = string(Postpaid)

	if d.HasChange("instance_type") && !d.IsNewResource() {
		request.DBInstanceClass = d.Get("instance_type").(string)
		update = true
		d.SetPartial("instance_type")
	}

	if d.HasChange("instance_storage") && !d.IsNewResource() {
		request.DBInstanceStorage = requests.NewInteger(d.Get("instance_storage").(int))
		update = true
		d.SetPartial("instance_storage")
	}

	if update {
		// wait instance status is running before modifying
		if err := client.WaitForDBInstance(d.Id(), Running, 500); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
		if _, err := conn.ModifyDBInstanceSpec(request); err != nil {
			return err
		}
		// wait instance status is running after modifying
		if err := client.WaitForDBInstance(d.Id(), Running, 1800); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
	}

	if d.HasChange("instance_name") {
		request := rds.CreateModifyDBInstanceDescriptionRequest()
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = d.Get("instance_name").(string)

		if _, err := conn.ModifyDBInstanceDescription(request); err != nil {
			return fmt.Errorf("ModifyDBInstanceDescription got an error: %#v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudDBInstanceRead(d, meta)
}

func resourceAlicloudDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	instance, err := client.DescribeDBInstanceById(d.Id())
	if err != nil {
		if NotFoundDBInstance(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err)
	}

	ips, err := client.GetSecurityIps(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Describe DB security ips error: %#v", err)
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
	d.Set("instance_name", instance.DBInstanceDescription)

	return nil
}

func resourceAlicloudDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	instance, err := client.DescribeDBInstanceById(d.Id())
	if err != nil {
		if NotFoundDBInstance(err) {
			return nil
		}
		return fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err)
	}
	if PayType(instance.PayType) == Prepaid {
		return fmt.Errorf("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically.")
	}

	request := rds.CreateDeleteDBInstanceRequest()
	request.DBInstanceId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.rdsconn.DeleteDBInstance(request)

		if err != nil {
			if NotFoundDBInstance(err) {
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

func buildDBCreateRequest(d *schema.ResourceData, meta interface{}) (*rds.CreateDBInstanceRequest, error) {
	client := meta.(*AliyunClient)
	request := rds.CreateCreateDBInstanceRequest()
	request.RegionId = string(getRegion(d, meta))
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.Engine = Trim(d.Get("engine").(string))
	request.DBInstanceStorage = requests.NewInteger(d.Get("instance_storage").(int))
	request.DBInstanceClass = Trim(d.Get("instance_type").(string))
	request.DBInstanceNetType = string(Intranet)

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))

	request.InstanceNetworkType = string(Classic)

	if vswitchId != "" {
		request.VSwitchId = vswitchId
		request.InstanceNetworkType = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		vsw, err := client.DescribeVswitch(vswitchId)
		if err != nil {
			return nil, fmt.Errorf("DescribeVSwitche got an error: %#v.", err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, fmt.Errorf("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId)
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, fmt.Errorf("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId)
		}

		request.VPCId = vsw.VpcId
	}

	request.PayType = Trim(d.Get("instance_charge_type").(string))

	// if charge type is postpaid, the commodity code must set to bards
	//args.CommodityCode = rds.Bards
	// At present, API supports two charge options about 'Prepaid'.
	// 'Month': valid period ranges [1-9]; 'Year': valid period range [1-3]
	// This resource only supports to input Month period [1-9, 12, 24, 36] and the values need to be converted before using them.
	if PayType(request.PayType) == Prepaid {

		period := d.Get("period").(int)
		request.UsedTime = strconv.Itoa(period)
		request.Period = string(Month)
		if period > 9 {
			request.UsedTime = strconv.Itoa(period / 12)
			request.Period = string(Year)
		}
	}

	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		uuid = resource.UniqueId()
	}
	request.ClientToken = fmt.Sprintf("Terraform-Alicloud-%d-%s", time.Now().Unix(), uuid)

	return request, nil
}
