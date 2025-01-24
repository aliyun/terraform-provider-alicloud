package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudKmsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudKmsInstanceCreate,
		Read:   resourceAliCloudKmsInstanceRead,
		Update: resourceAliCloudKmsInstanceUpdate,
		Delete: resourceAliCloudKmsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bind_vpcs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpc_owner_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ca_certificate_chain_pem": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"force_delete_without_backup": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(100, 100000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"log": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"0", "1"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"log_storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 500000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: IntInSlice([]int{0, 1, 2, 3, 6, 12, 24, 36}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"product_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"3", "5"}, false),
			},
			"renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 36),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"renew_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"secret_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 100000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"spec": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1000, 2000, 4000, 200}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 10000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"vswitch_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudKmsInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	var endpoint string
	var err error
	query := make(map[string]interface{})
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	parameterMapList := make([]map[string]interface{}, 0)
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Region",
		"Value": client.RegionId,
	})
	if v, ok := d.GetOk("vpc_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "VpcNum",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("spec"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Spec",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("secret_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "SecretNum",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("key_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "KeyNum",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("product_version"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "ProductVersion",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("log"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Log",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("log_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "LogStorage",
			"Value": fmt.Sprint(v),
		})
	}
	request["Parameter"] = parameterMapList

	request["Period"] = d.Get("period")
	if v, ok := d.GetOkExists("renew_period"); ok && v.(int) > 0 && request["SubscriptionType"] == "Subscription" {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOk("renew_status"); ok && request["SubscriptionType"] == "Subscription" {
		request["RenewalStatus"] = v
	}

	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if client.IsInternationalAccount() {
		request["ProductType"] = "kms_ddi_public_intl"
	}
	request["SubscriptionType"] = "Subscription"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductType"] = "kms_ppi_public_cn"
		request["SubscriptionType"] = v
		if client.IsInternationalAccount() {
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "kms_ddi_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductType"] = "kms_ppi_public_intl"
				}
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	kmsServiceV2 := KmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(id)}, d.Timeout(schema.TimeoutCreate), 10*time.Second, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "InstanceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	action = "ConnectKmsInstance"
	request = make(map[string]interface{})
	request["KmsInstanceId"] = d.Id()

	request["VpcId"] = d.Get("vpc_id")
	request["KMProvider"] = "Aliyun"
	jsonPathResult1, err := jsonpath.Get("$", d.Get("zone_ids"))
	if err == nil {
		request["ZoneIds"] = convertListToCommaSeparate(jsonPathResult1.(*schema.Set).List())
	}

	jsonPathResult2, err := jsonpath.Get("$", d.Get("vswitch_ids"))
	if err == nil {
		request["VSwitchIds"] = convertListToCommaSeparate(jsonPathResult2.(*schema.Set).List())
	}

	wait = incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"Forbidden.RamRoleNotFound"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["KmsInstanceId"]))

	kmsServiceV2 = KmsServiceV2{client}
	stateConf = BuildStateConf([]string{}, []string{"Connected"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "Status", []string{"Error"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudKmsInstanceUpdate(d, meta)
}

func resourceAliCloudKmsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsServiceV2 := KmsServiceV2{client}

	objectRaw, err := kmsServiceV2.DescribeKmsInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_instance DescribeKmsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CaCertificateChainPem"] != nil {
		d.Set("ca_certificate_chain_pem", objectRaw["CaCertificateChainPem"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["EndDate"] != nil {
		d.Set("end_date", objectRaw["EndDate"])
	}
	if objectRaw["InstanceName"] != nil {
		d.Set("instance_name", objectRaw["InstanceName"])
	}
	if objectRaw["KeyNum"] != nil {
		d.Set("key_num", objectRaw["KeyNum"])
	}
	if objectRaw["Log"] != nil {
		d.Set("log", objectRaw["Log"])
	}
	if objectRaw["LogStorage"] != nil {
		d.Set("log_storage", objectRaw["LogStorage"])
	}
	if objectRaw["SecretNum"] != nil {
		d.Set("secret_num", formatInt(objectRaw["SecretNum"]))
	}
	if objectRaw["Spec"] != nil {
		d.Set("spec", objectRaw["Spec"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["VpcId"] != nil {
		d.Set("vpc_id", objectRaw["VpcId"])
	}
	if objectRaw["VpcNum"] != nil {
		d.Set("vpc_num", objectRaw["VpcNum"])
	}

	bindVpc1Raw, _ := jsonpath.Get("$.BindVpcs.BindVpc", objectRaw)
	bindVpcsMaps := make([]map[string]interface{}, 0)
	if bindVpc1Raw != nil {
		for _, bindVpcChild1Raw := range bindVpc1Raw.([]interface{}) {
			bindVpcsMap := make(map[string]interface{})
			bindVpcChild1Raw := bindVpcChild1Raw.(map[string]interface{})
			bindVpcsMap["region_id"] = bindVpcChild1Raw["RegionId"]
			bindVpcsMap["vswitch_id"] = bindVpcChild1Raw["VSwitchId"]
			bindVpcsMap["vpc_id"] = bindVpcChild1Raw["VpcId"]
			bindVpcsMap["vpc_owner_id"] = bindVpcChild1Raw["VpcOwnerId"]

			bindVpcsMaps = append(bindVpcsMaps, bindVpcsMap)
		}
	}
	if bindVpc1Raw != nil {
		if err := d.Set("bind_vpcs", bindVpcsMaps); err != nil {
			return err
		}
	}

	vswitchIds := make([]interface{}, 0)
	if objectRaw["VswitchIds"] != nil {
		vswitchIds = objectRaw["VswitchIds"].([]interface{})
	}

	d.Set("vswitch_ids", vswitchIds)
	zoneIds := make([]interface{}, 0)
	if objectRaw["ZoneIds"] != nil {
		zoneIds = objectRaw["ZoneIds"].([]interface{})
	}
	d.Set("zone_ids", zoneIds)

	return nil
}

func resourceAliCloudKmsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChanges("vpc_num", "spec", "secret_num", "key_num", "log", "log_storage") {
		update = true
	}
	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("vpc_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "VpcNum",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("spec"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Spec",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("secret_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "SecretNum",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("key_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "KeyNum",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOk("log"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Log",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("log_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "LogStorage",
			"Value": fmt.Sprint(v),
		})
	}
	request["Parameter"] = parameterMapList

	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if client.IsInternationalAccount() {
		request["ProductType"] = "kms_ddi_public_intl"
	}
	request["SubscriptionType"] = "Subscription"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductType"] = "kms_ppi_public_cn"
		request["SubscriptionType"] = v
		if client.IsInternationalAccount() {
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}
	request["ModifyType"] = "Upgrade"
	if update && request["SubscriptionType"] == "Subscription" {
		var endpoint string
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductType"] = "kms_ddi_public_intl"
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductType"] = "kms_ppi_public_intl"
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
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

		kmsServiceV2 := KmsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("key_num"))}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "KeyNum", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateKmsInstanceBindVpc"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["KmsInstanceId"] = d.Id()

	if d.HasChange("bind_vpcs") {
		update = true
	}
	query["BindVpcs"] = "[]"
	if v, ok := d.GetOk("bind_vpcs"); ok || d.HasChange("bind_vpcs") {
		bindVpcsMaps := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["VpcId"] = dataLoopTmp["vpc_id"]

			if vpcOwnerId, ok := dataLoopTmp["vpc_owner_id"]; ok {
				vpcOwnerIdNumber, err := strconv.ParseInt(fmt.Sprint(vpcOwnerId), 10, 64)
				if err != nil {
					return WrapError(fmt.Errorf("convert vpc_owner_id to int64 failed, value: %v", vpcOwnerId))
				}
				dataLoopMap["VpcOwnerId"] = vpcOwnerIdNumber
			}

			dataLoopMap["RegionId"] = dataLoopTmp["region_id"]
			dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
			bindVpcsMaps = append(bindVpcsMaps, dataLoopMap)
		}
		bindVpcsMapsJson, err := json.Marshal(bindVpcsMaps)
		if err != nil {
			return WrapError(err)
		}
		query["BindVpcs"] = string(bindVpcsMapsJson)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, false)
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
	}

	d.Partial(false)
	return resourceAliCloudKmsInstanceRead(d, meta)
}

func resourceAliCloudKmsInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); !ok || v.(string) == "Subscription" {
		client := meta.(*connectivity.AliyunClient)
		action := "RefundInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		var err error
		query := make(map[string]interface{})
		request = make(map[string]interface{})
		request["InstanceId"] = d.Id()

		request["ClientToken"] = buildClientToken(action)
		request["ProductCode"] = "kms"
		request["ProductType"] = "kms_ddi_public_cn"
		if client.IsInternationalAccount() {
			request["ProductType"] = "kms_ddi_public_intl"
		}
		request["ImmediatelyRelease"] = "1"
		var endpoint string
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductType"] = "kms_ddi_public_intl"
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceNotExists"}) || NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		kmsServiceV2 := KmsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 30*time.Second, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "InstanceId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		return nil
	}

	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
		client := meta.(*connectivity.AliyunClient)
		action := "ReleaseKmsInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		var err error
		query := make(map[string]interface{})
		request = make(map[string]interface{})
		request["KmsInstanceId"] = d.Id()

		if v, ok := d.GetOk("force_delete_without_backup"); ok {
			request["ForceDeleteWithoutBackup"] = v
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, false)

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
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		return nil
	}
	return nil
}
