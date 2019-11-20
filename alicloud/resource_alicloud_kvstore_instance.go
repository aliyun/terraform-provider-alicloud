package alicloud

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"

	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudKVStoreInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKVStoreInstanceCreate,
		Read:   resourceAlicloudKVStoreInstanceRead,
		Update: resourceAlicloudKVStoreInstanceUpdate,
		Delete: resourceAlicloudKVStoreInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRKVInstanceName,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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
			"instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validateInstanceChargeType,
				Optional:     true,
				Default:      PostPaid,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: rkvPostPaidDiffSuppressFunc,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: rkvPostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 12),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: rkvPostPaidDiffSuppressFunc,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(KVStoreRedis),
				ValidateFunc: validateAllowedStringValue([]string{
					string(KVStoreMemcache),
					string(KVStoreRedis),
				}),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  KVStore2Dot8,
			},
			"connection_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"security_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},

			"vpc_auth_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{"Open", "Close"}),
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
				Optional: true,
				Computed: true,
			},

			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudKVStoreInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	request, err := buildKVStoreCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.CreateInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kvstore_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.CreateInstanceResponse)
	d.SetId(response.InstanceId)

	// wait instance status change from Creating to Normal
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudKVStoreInstanceUpdate(d, meta)
}

func resourceAlicloudKVStoreInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	d.Partial(true)
	stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging", "Changing"}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

	if err := kvstoreService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}
	if d.HasChange("parameters") {
		config := make(map[string]interface{})
		documented := d.Get("parameters").(*schema.Set).List()
		if len(documented) > 0 {
			for _, i := range documented {
				key := i.(map[string]interface{})["name"].(string)
				value := i.(map[string]interface{})["value"]
				config[key] = value
			}
			cfg, _ := json.Marshal(config)
			if err := kvstoreService.ModifyInstanceConfig(d.Id(), string(cfg)); err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("parameters")
	}

	if d.HasChange("security_ips") {
		// wait instance status is Normal before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		request := r_kvstore.CreateModifySecurityIpsRequest()
		request.RegionId = client.RegionId
		request.SecurityIpGroupName = "default"
		request.InstanceId = d.Id()
		if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
			request.SecurityIps = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
		} else {
			return WrapError(Error("Security ips cannot be empty"))
		}
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifySecurityIps(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("security_ips")
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("vpc_auth_mode") {
		if vswitchId, ok := d.GetOk("vswitch_id"); ok && vswitchId.(string) != "" {
			// vpc_auth_mode works only if the network type is VPC
			instanceType := d.Get("instance_type").(string)
			if string(KVStoreRedis) == instanceType {
				// wait instance status is Normal before modifying
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapError(err)
				}

				request := r_kvstore.CreateModifyInstanceVpcAuthModeRequest()
				request.RegionId = client.RegionId
				request.InstanceId = d.Id()
				request.VpcAuthMode = d.Get("vpc_auth_mode").(string)

				raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
					return rkvClient.ModifyInstanceVpcAuthMode(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				d.SetPartial("vpc_auth_mode")

				// The auth mode take some time to be effective, so wait to ensure the state !
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapError(err)
				}
			}
		}
	}

	if !d.IsNewResource() && (d.HasChange("instance_charge_type") || d.HasChange("period")) {
		prePaidRequest := r_kvstore.CreateTransformToPrePaidRequest()
		prePaidRequest.RegionId = client.RegionId
		prePaidRequest.InstanceId = d.Id()
		prePaidRequest.Period = requests.Integer(strconv.Itoa(d.Get("period").(int)))

		// for now we just support charge change from PostPaid to PrePaid
		configPayType := PayType(d.Get("instance_charge_type").(string))
		if configPayType == PrePaid {
			raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
				return rkvClient.TransformToPrePaid(prePaidRequest)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), prePaidRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(prePaidRequest.GetActionName(), raw, prePaidRequest.RpcRequest, prePaidRequest)
			// wait instance status is Normal after modifying
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}
			d.SetPartial("instance_charge_type")
			d.SetPartial("period")
		}
	}

	if d.HasChange("auto_renew") || d.HasChange("auto_renew_period") {
		request := r_kvstore.CreateModifyInstanceAutoRenewalAttributeRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()

		auto_renew := d.Get("auto_renew").(bool)
		if auto_renew {
			request.AutoRenew = "True"
		} else {
			request.AutoRenew = "False"
		}
		request.Duration = strconv.Itoa(d.Get("auto_renew_period").(int))

		raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
			return client.ModifyInstanceAutoRenewalAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("auto_renew")
		d.SetPartial("auto_renew_period")
	}

	if d.HasChange("maintain_start_time") || d.HasChange("maintain_end_time") {
		request := r_kvstore.CreateModifyInstanceMaintainTimeRequest()
		request.RegionId = client.RegionId
		request.InstanceId = d.Id()
		request.MaintainStartTime = d.Get("maintain_start_time").(string)
		request.MaintainEndTime = d.Get("maintain_end_time").(string)

		raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
			return client.ModifyInstanceMaintainTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("maintain_start_time")
		d.SetPartial("maintain_end_time")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudKVStoreInstanceRead(d, meta)
	}

	if d.HasChange("instance_class") {
		// wait instance status is Normal before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}

		request := r_kvstore.CreateModifyInstanceSpecRequest()
		request.RegionId = client.RegionId
		request.InstanceId = d.Id()
		request.InstanceClass = d.Get("instance_class").(string)
		request.EffectiveTime = "Immediately"

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
				return rkvClient.ModifyInstanceSpec(request)
			})
			if err != nil {
				if IsExceptedError(err, "MissingRedisUsedmemoryUnsupportPerfItem") {
					time.Sleep(time.Duration(5) * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		// There needs more time to sync instance class update
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			object, err := kvstoreService.DescribeKVstoreInstance(d.Id())
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if object.InstanceClass != request.InstanceClass {
				return resource.RetryableError(Error("Waitting for instance class is changed timeout. Expect instance class %s, got %s.",
					object.InstanceClass, request.InstanceClass))
			}
			return nil
		})
		if err != nil {
			return WrapError(err)
		}

		d.SetPartial("instance_class")
	}

	request := r_kvstore.CreateModifyInstanceAttributeRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	update := false
	if d.HasChange("instance_name") {
		request.InstanceName = d.Get("instance_name").(string)
		update = true
	}

	if d.HasChange("password") || d.HasChange("kms_encrypted_password") {
		if v := d.Get("password").(string); v != "" {
			d.SetPartial("password")
			request.NewPassword = v
			update = true
		}
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request.NewPassword = decryptResp.Plaintext
			d.SetPartial("kms_encrypted_password")
			d.SetPartial("kms_encryption_context")
			update = true
		}
	}

	if update {
		// wait instance status is Normal before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifyInstanceAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		d.SetPartial("instance_name")
		d.SetPartial("password")
	}

	d.Partial(false)
	return resourceAlicloudKVStoreInstanceRead(d, meta)
}

func resourceAlicloudKVStoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	object, err := kvstoreService.DescribeKVstoreInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_name", object.InstanceName)
	d.Set("instance_class", object.InstanceClass)
	d.Set("availability_zone", object.ZoneId)
	d.Set("instance_charge_type", object.ChargeType)
	d.Set("instance_type", object.InstanceType)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("engine_version", object.EngineVersion)
	d.Set("connection_domain", object.ConnectionDomain)
	d.Set("private_ip", object.PrivateIp)
	d.Set("security_ips", strings.Split(object.SecurityIPList, COMMA_SEPARATED))
	d.Set("vpc_auth_mode", object.VpcAuthMode)
	d.Set("maintain_start_time", object.MaintainStartTime)
	d.Set("maintain_end_time", object.MaintainEndTime)
	tags, err := kvstoreService.DescribeTags(d.Id(), TagResourceInstance)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", kvstoreService.tagsToMap(tags))

	if object.ChargeType == string(Prepaid) {
		request := r_kvstore.CreateDescribeInstanceAutoRenewalAttributeRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = d.Id()

		raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
			return client.DescribeInstanceAutoRenewalAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*r_kvstore.DescribeInstanceAutoRenewalAttributeResponse)
		if len(response.Items.Item) > 0 {
			renew := response.Items.Item[0]
			auto_renew := bool(renew.AutoRenew == "True")

			d.Set("auto_renew", auto_renew)
			d.Set("auto_renew_period", renew.Duration)
		}
	}
	//refresh parameters
	if err = refreshParameters(d, meta); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudKVStoreInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	object, err := kvstoreService.DescribeKVstoreInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if PayType(object.ChargeType) == Prepaid {
		return WrapError(Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically"))
	}
	request := r_kvstore.CreateDeleteInstanceRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DeleteInstance(request)
	})

	if err != nil {
		if !IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	stateConf := BuildStateConf([]string{"Creating", "Active", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return WrapError(err)
}

func buildKVStoreCreateRequest(d *schema.ResourceData, meta interface{}) (*r_kvstore.CreateInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := r_kvstore.CreateCreateInstanceRequest()
	request.InstanceName = Trim(d.Get("instance_name").(string))
	request.RegionId = client.RegionId
	request.InstanceType = Trim(d.Get("instance_type").(string))
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.InstanceClass = Trim(d.Get("instance_class").(string))
	request.ChargeType = Trim(d.Get("instance_charge_type").(string))

	request.Password = d.Get("password").(string)
	if request.Password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return request, WrapError(err)
			}
			request.Password = decryptResp.Plaintext
		}
	}

	request.BackupId = Trim(d.Get("backup_id").(string))

	if PayType(request.ChargeType) == PrePaid {
		request.Period = strconv.Itoa(d.Get("period").(int))
	}

	if zone, ok := d.GetOk("availability_zone"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	request.NetworkType = strings.ToUpper(string(Classic))
	if vswitchId, ok := d.GetOk("vswitch_id"); ok && vswitchId.(string) != "" {
		request.VSwitchId = vswitchId.(string)
		request.NetworkType = strings.ToUpper(string(Vpc))
		request.PrivateIpAddress = Trim(d.Get("private_ip").(string))

		// check vswitchId in zone
		object, err := vpcService.DescribeVSwitch(vswitchId.(string))
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = object.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(object.ZoneId)[len(object.ZoneId)-1])) {
				return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %s", object.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != object.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s", object.VSwitchId, request.ZoneId))
		}

		request.VpcId = object.VpcId
	}

	request.Token = buildClientToken(request.GetActionName())

	return request, nil
}

func refreshParameters(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	var param []map[string]interface{}
	documented, ok := d.GetOk("parameters")
	if !ok {
		d.Set("parameters", param)
		return nil
	}
	object, err := kvstoreService.DescribeParameters(d.Id())
	if err != nil {
		return WrapError(err)
	}

	var parameters = make(map[string]interface{})
	for _, i := range object.RunningParameters.Parameter {
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, i := range object.ConfigParameters.Parameter {
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, parameter := range documented.(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, value := range parameters {
			if value.(map[string]interface{})["name"] == name {
				param = append(param, value.(map[string]interface{}))
				break
			}
		}
	}

	d.Set("parameters", param)
	return nil
}
