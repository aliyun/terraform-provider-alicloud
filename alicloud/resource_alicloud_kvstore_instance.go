package alicloud

import (
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strings"
	"time"
)

func resourceAlicloudKvstoreInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKvstoreInstanceCreate,
		Read:   resourceAlicloudKvstoreInstanceRead,
		Update: resourceAlicloudKvstoreInstanceUpdate,
		Delete: resourceAlicloudKvstoreInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: redisPostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntBetween(1, 12),
				DiffSuppressFunc: redisPostPaidAndRenewDiffSuppressFunc,
			},
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.IsNewResource()
				},
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backup_period": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"backup_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"business_info": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"config": {
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"parameters"},
			},
			"connection_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"coupon_no": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "youhuiquan_promotion_option_id_for_blank",
			},
			"db_instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_name"},
			},
			"instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Removed:       "Field 'instance_name' has been removed from version 1.126.0. Use 'db_instance_name' instead.",
				ConflictsWith: []string{"db_instance_name"},
			},
			"dedicated_host_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_backup_log": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
				Default:      0,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"2.8", "4.0", "5.0", "6.0"}, false),
				Default:      "5.0",
			},
			"force_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"global_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.IsNewResource()
				},
			},
			"global_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_release_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Memcache", "Redis"}, false),
				Default:      "Redis",
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modify_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"MASTER_SLAVE", "STAND_ALONE", "double", "single"}, false),
				Removed:      "Field 'node_type' has been removed from version 1.126.1",
			},
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"UPGRADE", "DOWNGRADE"}, false),
				Default:      "UPGRADE",
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"payment_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				ConflictsWith: []string{"instance_charge_type"},
			},
			"instance_charge_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				Removed:       "Field 'instance_charge_type' has been removed from version 1.126.0. Use 'payment_type' instead.",
				ConflictsWith: []string{"payment_type"},
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"1", "12", "2", "24", "3", "36", "4", "5", "6", "7", "8", "9"}, false),
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"private_connection_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_connection_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"qps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"restore_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ssl_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disable", "Enable", "Update"}, false),
			},
			"security_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: redisSecurityGroupIdDiffSuppressFunc,
			},
			"security_ip_group_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_ip_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"srcdb_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_audit_retention": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_audit": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_auth_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Close", "Open"}, false),
				Default:      "Open",
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Removed:       "Field 'availability_zone' has been removed from version 1.126.0. Use 'zone_id' instead.",
				ConflictsWith: []string{"zone_id"},
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"enable_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Removed:  "Field 'enable_public' has been removed from version 1.126.0. Please use resource 'alicloud_kvstore_connection' instead.",
			},
			"connection_string_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'connection_string_prefix' has been removed from version 1.126.0. Please use resource 'alicloud_kvstore_connection' instead.",
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Field 'connection_string' has been removed from version 1.126.0. Please use resource 'alicloud_kvstore_connection' instead.",
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: func(v interface{}) int {
					return hashcode.String(
						v.(map[string]interface{})["name"].(string) + "|" + v.(map[string]interface{})["value"].(string))
				},
				Optional:      true,
				Computed:      true,
				Removed:       "Field 'parameters' has been removed from version 1.126.0. Use 'config' instead.",
				ConflictsWith: []string{"config"},
			},
		},
	}
}
func resourceAlicloudKvstoreInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})
	conn, err := client.NewRedisaClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	if v, ok := d.GetOk("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}
	if v, ok := d.GetOk("backup_id"); ok {
		request["BackupId"] = v
	}
	if v, ok := d.GetOk("business_info"); ok {
		request["BusinessInfo"] = v
	}
	if v, ok := d.GetOk("capacity"); ok {
		request["Capacity"] = v
	}
	if v, ok := d.GetOk("coupon_no"); ok {
		request["CouponNo"] = v
	}
	if v, ok := d.GetOk("db_instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("dedicated_host_group_id"); ok {
		request["DedicatedHostGroupId"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("engine_version"); ok {
		request["EngineVersion"] = v
	}
	if v, ok := d.GetOkExists("global_instance"); ok {
		request["GlobalInstance"] = v
	}
	if v, ok := d.GetOk("global_instance_id"); ok {
		request["GlobalInstanceId"] = v
	}
	if v, ok := d.GetOk("instance_class"); ok {
		request["InstanceClass"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}

	request["NetworkType"] = "CLASSIC"
	if v, ok := d.GetOk("node_type"); ok {
		request["NodeType"] = v
	}

	request["Password"] = d.Get("password")
	if request["Password"].(string) == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			plaintext, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request["Password"] = plaintext
		}
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("port"); ok {
		request["Port"] = v
	}
	if v, ok := d.GetOk("private_ip"); ok {
		request["PrivateIpAddress"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("restore_time"); ok {
		request["RestoreTime"] = v
	}
	if v, ok := d.GetOk("secondary_zone_id"); ok {
		request["SecondaryZoneId"] = v
	}
	if v, ok := d.GetOk("srcdb_instance_id"); ok {
		request["SrcDBInstanceId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request["NetworkType"] = "VPC"
		request["VpcId"] = vsw["VpcId"]
		request["VSwitchId"] = vswitchId
		if v, ok := request["ZoneId"].(string); !ok || v == "" {
			request["ZoneId"] = vsw["ZoneId"]
		}
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kvstore_instance", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["InstanceId"]))
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 300*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudKvstoreInstanceUpdate(d, meta)
}
func resourceAlicloudKvstoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	object, err := r_kvstoreService.DescribeKvstoreInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kvstore_instance r_kvstoreService.DescribeKvstoreInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	net, _ := r_kvstoreService.DescribeKvstoreConnection(d.Id())
	for _, instanceNetInfo := range net {
		if instanceNetInfo.DBInstanceNetType == "2" {
			d.Set("private_connection_port", instanceNetInfo.Port)
		}
	}
	d.Set("capacity", fmt.Sprint(formatInt(object["Capacity"])))
	d.Set("config", object["Config"])
	d.Set("connection_domain", object["ConnectionDomain"])
	d.Set("db_instance_name", object["InstanceName"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("instance_class", object["InstanceClass"])
	d.Set("instance_release_protection", object["InstanceReleaseProtection"])
	d.Set("instance_type", object["InstanceType"])
	d.Set("maintain_end_time", object["MaintainEndTime"])
	d.Set("maintain_start_time", object["MaintainStartTime"])
	d.Set("payment_type", object["ChargeType"])
	d.Set("private_ip", object["PrivateIp"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["InstanceStatus"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("vpc_auth_mode", object["VpcAuthMode"])
	d.Set("zone_id", object["ZoneId"])
	describeBackupPolicyObject, err := r_kvstoreService.DescribeBackupPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("backup_period", strings.Split(describeBackupPolicyObject["PreferredBackupPeriod"].(string), ","))
	d.Set("backup_time", describeBackupPolicyObject["PreferredBackupTime"])
	if object["ChargeType"] == string(PrePaid) {
		describeInstanceAutoRenewalAttributeObject, err := r_kvstoreService.DescribeInstanceAutoRenewalAttribute(d.Id())
		if err != nil {
			return WrapError(err)
		}
		d.Set("auto_renew", describeInstanceAutoRenewalAttributeObject["AutoRenew"])
		d.Set("auto_renew_period", fmt.Sprint(formatInt(describeInstanceAutoRenewalAttributeObject["Duration"])))
	}
	if _, ok := d.GetOk("ssl_enable"); ok {
		describeInstanceSSLObject, err := r_kvstoreService.DescribeInstanceSSL(d.Id())
		if err != nil {
			return WrapError(err)
		}
		d.Set("ssl_enable", describeInstanceSSLObject["SSLEnabled"])
	}
	if _, ok := d.GetOk("security_group_id"); ok {
		describeSecurityGroupConfigurationObject, err := r_kvstoreService.DescribeSecurityGroupConfiguration(d.Id())
		if err != nil {
			return WrapError(err)
		}
		sgs := make([]string, 0)
		for _, sg := range describeSecurityGroupConfigurationObject {
			sg_data := sg.(map[string]interface{})
			sgs = append(sgs, sg_data["SecurityGroupId"].(string))
		}
		d.Set("security_group_id", strings.Join(sgs, ","))
	}
	if _, ok := d.GetOk("security_ips"); ok {
		describeSecurityIpsObject, err := r_kvstoreService.DescribeSecurityIps(d.Id())
		if err != nil {
			return WrapError(err)
		}
		d.Set("security_ip_group_attribute", describeSecurityIpsObject["SecurityIpGroupAttribute"])
		d.Set("security_ip_group_name", describeSecurityIpsObject["SecurityIpGroupName"])
		d.Set("security_ips", strings.Split(describeSecurityIpsObject["SecurityIpList"].(string), ","))
	}

	return nil
}
func resourceAlicloudKvstoreInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	var response map[string]interface{}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := r_kvstoreService.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("config") {
		update = true
	}
	respJson, err := convertMaptoJsonString(d.Get("config").(map[string]interface{}))
	request["Config"] = respJson

	if err != nil {
		return WrapError(err)
	}

	if update {
		action := "ModifyInstanceConfig"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("config")
	}
	update = false
	modifyInstanceSSLReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("ssl_enable") {
		update = true
	}
	modifyInstanceSSLReq["SSLEnabled"] = d.Get("ssl_enable")
	if update {
		action := "ModifyInstanceSSL"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyInstanceSSLReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyInstanceSSLReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("ssl_enable")
	}
	update = false
	modifyInstanceVpcAuthModeReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("vpc_auth_mode") {
		update = true
	}
	modifyInstanceVpcAuthModeReq["VpcAuthMode"] = d.Get("vpc_auth_mode").(string)
	if update {
		action := "ModifyInstanceVpcAuthMode"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyInstanceVpcAuthModeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyInstanceVpcAuthModeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("vpc_auth_mode")
	}
	update = false
	modifyResourceGroupReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	modifyResourceGroupReq["ResourceGroupId"] = d.Get("resource_group_id")
	if update {
		action := "ModifyResourceGroup"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ModifyResourceGroup")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyResourceGroupReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyResourceGroupReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
	}
	update = false
	modifySecurityGroupConfigurationReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("security_group_id") {
		update = true
	}
	modifySecurityGroupConfigurationReq["SecurityGroupId"] = d.Get("security_group_id")
	if update {
		action := "ModifySecurityGroupConfiguration"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifySecurityGroupConfigurationReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifySecurityGroupConfigurationReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_group_id")
	}
	update = false
	modifyAuditLogConfigReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("db_audit") {
		update = true
		modifyAuditLogConfigReq["DbAudit"] = d.Get("db_audit")
	}
	if d.HasChange("db_audit_retention") {
		update = true
		modifyAuditLogConfigReq["Retention"] = d.Get("db_audit_retention")
	}
	if update {
		action := "ModifyAuditLogConfig"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyAuditLogConfigReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyAuditLogConfigReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("db_audit")
		d.SetPartial("db_audit_retention")
	}
	update = false
	modifyInstanceAutoRenewalAttributeReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("auto_renew") {
		update = true
		modifyInstanceAutoRenewalAttributeReq["AutoRenew"] = d.Get("auto_renew")
	}
	if !d.IsNewResource() && d.HasChange("auto_renew_period") {
		update = true
		modifyInstanceAutoRenewalAttributeReq["Duration"] = d.Get("auto_renew_period")
	}
	if update {
		action := "ModifyInstanceAutoRenewalAttribute"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyInstanceAutoRenewalAttributeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyInstanceAutoRenewalAttributeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("auto_renew")
		d.SetPartial("auto_renew_period")
	}
	update = false
	modifyInstanceMaintainTimeReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("maintain_end_time") {
		update = true
	}
	modifyInstanceMaintainTimeReq["MaintainEndTime"] = d.Get("maintain_end_time")
	if d.HasChange("maintain_start_time") {
		update = true
	}
	modifyInstanceMaintainTimeReq["MaintainStartTime"] = d.Get("maintain_start_time")
	if update {
		action := "ModifyInstanceMaintainTime"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyInstanceMaintainTimeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyInstanceMaintainTimeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("maintain_end_time")
		d.SetPartial("maintain_start_time")
	}
	update = false
	modifyInstanceMajorVersionReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("engine_version") {
		update = true
	}
	modifyInstanceMajorVersionReq["MajorVersion"] = d.Get("engine_version")
	if update {
		if v, ok := d.GetOk("effective_time"); ok {
			modifyInstanceMajorVersionReq["EffectiveTime"] = v
		}
		action := "ModifyInstanceMajorVersion"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyInstanceMajorVersionReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyInstanceMajorVersionReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 300*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("engine_version")
	}
	update = false
	resetAccountPasswordReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && (d.HasChange("password") || d.HasChange("kms_encrypted_password")) {
		update = true
		password := d.Get("password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password != "" {
			resetAccountPasswordReq["AccountPassword"] = password
		} else {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			plaintext, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			resetAccountPasswordReq["AccountPassword"] = plaintext
		}

		if password == "" && kmsPassword == "" {
			return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
		}
	}

	if update {
		if v, ok := d.GetOk("account_name"); ok {
			resetAccountPasswordReq["AccountName"] = v
		}
		action := "ResetAccountPassword"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, resetAccountPasswordReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, resetAccountPasswordReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("kms_encrypted_password")
		d.SetPartial("kms_encryption_context")
		d.SetPartial("password")
	}
	update = false
	modifyBackupPolicyReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("backup_period") {
		update = true
	}
	modifyBackupPolicyReq["PreferredBackupPeriod"] = convertListToCommaSeparate(d.Get("backup_period").(*schema.Set).List())
	if d.HasChange("backup_time") {
		update = true
	}
	modifyBackupPolicyReq["PreferredBackupTime"] = d.Get("backup_time")
	if update {
		if v, ok := d.GetOk("enable_backup_log"); ok {
			modifyBackupPolicyReq["EnableBackupLog"] = v
		}
		action := "ModifyBackupPolicy"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyBackupPolicyReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyBackupPolicyReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("backup_period")
		d.SetPartial("backup_time")
	}
	update = false
	modifyDBInstanceConnectionStringReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("connection_domain") {
		update = true
	}
	modifyDBInstanceConnectionStringReq["CurrentConnectionString"] = d.Get("connection_domain")
	modifyDBInstanceConnectionStringReq["IPType"] = "Private"
	if d.HasChange("private_connection_prefix") {
		update = true
		modifyDBInstanceConnectionStringReq["NewConnectionString"] = d.Get("private_connection_prefix")
	}
	if update {
		action := "ModifyDBInstanceConnectionString"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyDBInstanceConnectionStringReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceConnectionStringReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("connection_domain")
		d.SetPartial("private_connection_prefix")
	}
	update = false
	modifyInstanceAttributeReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && (d.HasChange("db_instance_name") || d.HasChange("instance_name")) {
		update = true
		if _, ok := d.GetOk("db_instance_name"); ok {
			modifyInstanceAttributeReq["InstanceName"] = d.Get("db_instance_name")
		} else {
			modifyInstanceAttributeReq["InstanceName"] = d.Get("instance_name")
		}
	}
	if d.HasChange("instance_release_protection") || d.IsNewResource() {
		update = true
		modifyInstanceAttributeReq["InstanceReleaseProtection"] = d.Get("instance_release_protection")
	}
	if !d.IsNewResource() && (d.HasChange("password") || d.HasChange("kms_encrypted_password")) {
		update = true
		password := d.Get("password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password != "" {
			modifyInstanceAttributeReq["NewPassword"] = password
		} else {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			plaintext, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			modifyInstanceAttributeReq["NewPassword"] = plaintext
		}

		if password == "" && kmsPassword == "" {
			return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
		}
	}
	if update {
		action := "ModifyInstanceAttribute"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyInstanceAttributeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyInstanceAttributeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_name")
		d.SetPartial("db_instance_name")
		d.SetPartial("instance_release_protection")
		d.SetPartial("kms_encrypted_password")
		d.SetPartial("kms_encryption_context")
		d.SetPartial("password")
	}
	update = false
	migrateToOtherZoneReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && (d.HasChange("zone_id") || d.HasChange("availability_zone")) {
		update = true
	}
	if _, ok := d.GetOk("zone_id"); ok {
		migrateToOtherZoneReq["ZoneId"] = d.Get("zone_id")
	} else {
		migrateToOtherZoneReq["ZoneId"] = d.Get("availability_zone")
	}
	if !d.IsNewResource() && d.HasChange("vswitch_id") {
		update = true
		migrateToOtherZoneReq["VSwitchId"] = d.Get("vswitch_id")
	}
	if update {
		if v, ok := d.GetOk("effective_time"); ok {
			migrateToOtherZoneReq["EffectiveTime"] = v
		}
		if v, ok := d.GetOk("secondary_zone_id"); ok {
			migrateToOtherZoneReq["SecondaryZoneId"] = v
		}
		action := "MigrateToOtherZone"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, migrateToOtherZoneReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, migrateToOtherZoneReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 600*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("availability_zone")
		d.SetPartial("zone_id")
		d.SetPartial("vswitch_id")
	}
	update = false
	modifySecurityIpsReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("security_ips") {
		update = true
	}
	modifySecurityIpsReq["SecurityIps"] = convertListToCommaSeparate(d.Get("security_ips").(*schema.Set).List())
	if d.HasChange("security_ip_group_attribute") {
		update = true
		modifySecurityIpsReq["SecurityIpGroupAttribute"] = d.Get("security_ip_group_attribute")
	}
	if d.HasChange("security_ip_group_name") {
		update = true
		modifySecurityIpsReq["SecurityIpGroupName"] = d.Get("security_ip_group_name")
	}
	if update {
		if v, ok := d.GetOk("modify_mode"); ok {
			modifySecurityIpsReq["ModifyMode"] = v
		}
		action := "ModifySecurityIps"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifySecurityIpsReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifySecurityIpsReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_ips")
		d.SetPartial("security_ip_group_attribute")
		d.SetPartial("security_ip_group_name")
	}
	update = false
	transformInstanceChargeTypeReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && (d.HasChange("payment_type") || d.HasChange("instance_charge_type")) {
		update = true
	}
	if _, ok := d.GetOk("payment_type"); ok {
		transformInstanceChargeTypeReq["ChargeType"] = d.Get("payment_type")
	} else {
		transformInstanceChargeTypeReq["ChargeType"] = d.Get("instance_charge_type")
	}
	transformInstanceChargeTypeReq["AutoPay"] = true
	if !d.IsNewResource() && d.HasChange("period") {
		update = true
		transformInstanceChargeTypeReq["Period"] = d.Get("period")
	}
	if update {
		action := "TransformInstanceChargeType"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, transformInstanceChargeTypeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, transformInstanceChargeTypeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_charge_type")
		d.SetPartial("payment_type")
		d.SetPartial("period")
	}
	update = false
	modifyInstanceSpecReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	modifyInstanceSpecReq["AutoPay"] = true
	modifyInstanceSpecReq["EffectiveTime"] = "0"

	if !d.IsNewResource() && d.HasChange("engine_version") {

		modifyInstanceSpecReq["MajorVersion"] = d.Get("engine_version")

	}
	if !d.IsNewResource() && d.HasChange("instance_class") {
		update = true
		modifyInstanceSpecReq["InstanceClass"] = d.Get("instance_class")
	}

	if update {
		if v, ok := d.GetOk("business_info"); ok {
			modifyInstanceSpecReq["BusinessInfo"] = v
		}
		if v, ok := d.GetOk("coupon_no"); ok {
			modifyInstanceSpecReq["CouponNo"] = v
		}
		if v, ok := d.GetOk("effective_time"); ok {
			modifyInstanceSpecReq["EffectiveTime"] = v
		}
		if v, ok := d.GetOkExists("force_upgrade"); ok {
			modifyInstanceSpecReq["ForceUpgrade"] = v
		}
		if v, ok := d.GetOk("order_type"); ok {
			modifyInstanceSpecReq["OrderType"] = v
		}
		if v, ok := d.GetOk("source_biz"); ok {
			modifyInstanceSpecReq["SourceBiz"] = v
		}
		action := "ModifyInstanceSpec"
		conn, err := client.NewRedisaClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ModifyInstanceSpec")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, modifyInstanceSpecReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"MissingRedisUsedmemoryUnsupportPerfItem"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyInstanceSpecReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 360*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("engine_version")
		d.SetPartial("instance_class")
	}
	if d.HasChange("parameters") {
		request := r_kvstore.CreateModifyInstanceConfigRequest()
		request.InstanceId = d.Id()
		config := make(map[string]interface{})
		documented := d.Get("parameters").(*schema.Set).List()
		if len(documented) > 0 {
			for _, i := range documented {
				key := i.(map[string]interface{})["name"].(string)
				value := i.(map[string]interface{})["value"]
				config[key] = value
			}
			cfg, _ := convertMaptoJsonString(config)
			request.Config = cfg

			raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
				return r_kvstoreClient.ModifyInstanceConfig(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		d.SetPartial("parameters")
	}
	if d.HasChange("enable_public") {
		prefix := d.Get("connection_string_prefix").(string)
		port := fmt.Sprint(d.Get("port").(int))
		target := d.Get("enable_public").(bool)

		if target {
			request := r_kvstore.CreateAllocateInstancePublicConnectionRequest()
			request.InstanceId = d.Id()
			request.ConnectionStringPrefix = prefix
			request.Port = port

			raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
				return client.AllocateInstancePublicConnection(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if !target && d.Get("connection_string") != "" {
			request := r_kvstore.CreateReleaseInstancePublicConnectionRequest()
			request.InstanceId = d.Id()
			request.CurrentConnectionString = d.Get("connection_string").(string)

			raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
				return client.ReleaseInstancePublicConnection(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("enable_public")
	}
	if d.HasChange("payment_type") {
		object, err := r_kvstoreService.DescribeKvstoreInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("payment_type").(string)
		if object["ChargeType"].(string) != target {
			if target == "PrePaid" {
				request := map[string]interface{}{
					"InstanceId": d.Id(),
				}
				if v, ok := d.GetOk("period"); ok {
					request["Period"] = v
				}
				if v, ok := d.GetOk("auto_renew"); ok {
					request["AutoPay"] = v
				}
				action := "TransformToPrePaid"
				conn, err := client.NewRedisaClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, r_kvstoreService.KvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("payment_type")
		}
	}
	d.Partial(false)
	return resourceAlicloudKvstoreInstanceRead(d, meta)
}
func resourceAlicloudKvstoreInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteInstance"
	var response map[string]interface{}
	conn, err := client.NewRedisaClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	if v, ok := d.GetOk("global_instance_id"); ok {
		request["GlobalInstanceId"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertModifyModeRequest(input int) string {
	switch input {
	case 0:
		return "Cover"
	case 1:
		return "Append"
	case 2:
		return "Delete"
	}
	return ""
}
