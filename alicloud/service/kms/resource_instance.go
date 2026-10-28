package kms

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"end_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"force_delete_without_backup": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"true", "false"}, false),
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(100, 100000),
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
				ValidateFunc: validators.StringInSlice([]string{"0", "1"}, false),
			},
			"log_storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.IntBetween(0, 500000),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validators.IntInSlice([]int{0, 1, 2, 3, 6, 12, 24, 36}),
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
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"3", "5"}, false),
			},
			"renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 36),
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
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"renewal_period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"M", "Y"}, false),
			},
			"secret_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 100000),
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
				ValidateFunc: validators.IntInSlice([]int{0, 1000, 2000, 4000, 200}),
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
			"tags": tags.TagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 10000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
						return true
					}
					return false
				},
			},
			"vswitch_ids": {
				Type:             schema.TypeList,
				Required:         true,
				ForceNew:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: kmsInstanceIdsOrderDiffSuppressFunc,
			},
			"zone_ids": {
				Type:             schema.TypeList,
				Required:         true,
				ForceNew:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: kmsInstanceIdsOrderDiffSuppressFunc,
			},
		},
	}
}

func resourceAliCloudKmsInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = helper.BuildClientToken(action)

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

	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOkExists("renew_period"); ok && v.(int) > 0 {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOk("renew_status"); ok {
		request["RenewalStatus"] = v
	}
	request["SubscriptionType"] = "Subscription"
	if v, ok := d.GetOk("payment_type"); ok {
		request["SubscriptionType"] = v.(string)
	}
	var endpoint string
	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "kms"
		request["ProductType"] = "kms_ppi_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "kms"
		request["ProductType"] = "kms_ddi_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "kms"
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}
	wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalError"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "kms"
				request["ProductType"] = "kms_ddi_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductCode"] = "kms"
					request["ProductType"] = "kms_ppi_public_intl"
				}
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_kms_instance", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "null" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	if fmt.Sprint(response["Success"]) != "true" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	kmsService := KmsServiceV2{client}
	stateConfig := helper.BuildStateConf([]string{}, []string{fmt.Sprint(id)}, d.Timeout(schema.TimeoutCreate), 10*time.Second, kmsService.KmsInstanceStateRefreshFunc(d.Id(), "InstanceId", []string{}))
	if _, err := stateConfig.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	action = "ConnectKmsInstance"
	request = make(map[string]interface{})
	request["VpcId"] = d.Get("vpc_id")
	request["KmsInstanceId"] = d.Id()
	request["KMProvider"] = "Aliyun"
	zoneIdsJsonPath, err := jsonpath.Get("$", d.Get("zone_ids"))
	if err == nil {
		request["ZoneIds"] = convert.ConvertListToCommaSeparate(convert.ConvertToInterfaceArray(zoneIdsJsonPath))
	}

	vswitchIdsJsonPath, err := jsonpath.Get("$", d.Get("vswitch_ids"))
	if err == nil {
		request["VSwitchIds"] = convert.ConvertListToCommaSeparate(convert.ConvertToInterfaceArray(vswitchIdsJsonPath))
	}

	wait = helper.IncrementalWait(5*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Forbidden.RamRoleNotFound"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_kms_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["KmsInstanceId"]))

	kmsServiceV2 := KmsServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Connected"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "Status", []string{"Error"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudKmsInstanceUpdate(d, meta)
}

func resourceAliCloudKmsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsServiceV2 := KmsServiceV2{client}

	objectRaw, err := kmsServiceV2.DescribeKmsInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_instance DescribeKmsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("ca_certificate_chain_pem", objectRaw["CaCertificateChainPem"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("deletion_protection", objectRaw["DeletionProtection"])
	d.Set("end_date", objectRaw["EndDate"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("key_num", objectRaw["KeyNum"])
	d.Set("log", objectRaw["Log"])
	d.Set("log_storage", objectRaw["LogStorage"])
	d.Set("payment_type", convertKmsInstanceKmsInstanceChargeTypeResponse(objectRaw["ChargeType"]))
	d.Set("product_version", objectRaw["ProductVersion"])
	secret_numVal, formatErr := convert.FormatInt(objectRaw["SecretNum"])
	if formatErr != nil {
		return formatErr
	}
	d.Set("secret_num", secret_numVal)
	d.Set("spec", objectRaw["Spec"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("vpc_num", objectRaw["VpcNum"])

	bindVpcRaw, _ := jsonpath.Get("$.BindVpcs.BindVpc", objectRaw)
	bindVpcsMaps := make([]map[string]interface{}, 0)
	if bindVpcRaw != nil {
		for _, bindVpcChildRaw := range convert.ConvertToInterfaceArray(bindVpcRaw) {
			bindVpcsMap := make(map[string]interface{})
			bindVpcChildRaw := bindVpcChildRaw.(map[string]interface{})
			bindVpcsMap["region_id"] = bindVpcChildRaw["RegionId"]
			bindVpcsMap["vswitch_id"] = bindVpcChildRaw["VSwitchId"]
			bindVpcsMap["vpc_id"] = bindVpcChildRaw["VpcId"]
			if v, ok := bindVpcChildRaw["VpcOwnerId"]; ok {
				bindVpcsMap["vpc_owner_id"] = v
			}

			bindVpcsMaps = append(bindVpcsMaps, bindVpcsMap)
		}
	}
	if err := d.Set("bind_vpcs", bindVpcsMaps); err != nil {
		return err
	}

	vswitchIdsRaw := make([]interface{}, 0)
	if objectRaw["VswitchIds"] != nil {
		vswitchIdsRaw = convert.ConvertToInterfaceArray(objectRaw["VswitchIds"])
	}

	d.Set("vswitch_ids", vswitchIdsRaw)
	zoneIdsRaw := make([]interface{}, 0)
	if objectRaw["ZoneIds"] != nil {
		zoneIdsRaw = convert.ConvertToInterfaceArray(objectRaw["ZoneIds"])
	}

	d.Set("zone_ids", zoneIdsRaw)

	objectRaw, err = kmsServiceV2.DescribeInstanceQueryAvailableInstances(d)
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	d.Set("renew_period", objectRaw["RenewalDuration"])
	d.Set("renew_status", objectRaw["RenewStatus"])

	objectRaw, err = kmsServiceV2.DescribeInstanceListTagResources(d.Id())
	if err != nil && !errs.NotFoundError(err) && !errs.IsExpectedErrors(err, []string{"Rejected.Uninitialized"}) {
		return errs.WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tags.TagsToMap(tagsMaps))
	return nil
}

func resourceAliCloudKmsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)

	parameterMapList := make([]map[string]interface{}, 0)

	if !d.IsNewResource() && d.HasChange("vpc_num") {
		update = true

		if v, ok := d.GetOk("vpc_num"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "VpcNum",
				"Value": fmt.Sprint(v),
			})
		}
	}

	if !d.IsNewResource() && d.HasChange("spec") {
		update = true

		if v, ok := d.GetOk("spec"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "Spec",
				"Value": fmt.Sprint(v),
			})
		}
	}

	if !d.IsNewResource() && d.HasChange("secret_num") {
		update = true

		if v, ok := d.GetOk("secret_num"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "SecretNum",
				"Value": fmt.Sprint(v),
			})
		}
	}

	if !d.IsNewResource() && d.HasChange("key_num") {
		update = true

		if v, ok := d.GetOk("key_num"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "KeyNum",
				"Value": fmt.Sprint(v),
			})
		}
	}

	if !d.IsNewResource() && d.HasChange("product_version") {
		update = true

		if v, ok := d.GetOk("product_version"); ok {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "ProductVersion",
				"Value": v,
			})
		}
	}

	if !d.IsNewResource() && d.HasChanges("log", "log_storage") {
		update = true
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

	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	request["SubscriptionType"] = d.Get("payment_type")
	request["ModifyType"] = "Upgrade"
	var endpoint string
	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "kms"
		request["ProductType"] = "kms_ppi_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "kms"
		request["ProductType"] = "kms_ddi_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "kms"
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}
	if update {
		wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"InternalError"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductCode"] = "kms"
					request["ProductType"] = "kms_ddi_public_intl"
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductCode"] = "kms"
						request["ProductType"] = "kms_ppi_public_intl"
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		kmsServiceV2 := KmsServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "#KeyNum", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
	if v, ok := d.GetOk("bind_vpcs"); ok || d.HasChange("bind_vpcs") {
		bindVpcsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convert.ConvertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["VpcId"] = dataLoopTmp["vpc_id"]

			if vpcOwnerIdRaw, ok := dataLoopTmp["vpc_owner_id"]; ok && vpcOwnerIdRaw != "" {
				vpcOwnerId1, _ := strconv.ParseInt(vpcOwnerIdRaw.(string), 10, 64)
				dataLoopMap["VpcOwnerId"] = vpcOwnerId1

			}
			dataLoopMap["RegionId"] = dataLoopTmp["region_id"]
			dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
			bindVpcsMapsArray = append(bindVpcsMapsArray, dataLoopMap)
		}
		bindVpcsMapsJson, err := json.Marshal(bindVpcsMapsArray)
		if err != nil {
			return errs.WrapError(err)
		}
		query["BindVpcs"] = string(bindVpcsMapsJson)
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcGet("Kms", "2016-01-20", action, query, request)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Forbidden.KMSInstanceInUpgrade"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "UpdateKmsInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KmsInstanceId"] = d.Id()

	if d.HasChange("instance_name") {
		update = true
		request["InstanceName"] = d.Get("instance_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "SetRenewal"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()

	if !d.IsNewResource() && d.HasChange("renew_status") {
		update = true
	}
	request["RenewalStatus"] = d.Get("renew_status")
	if !d.IsNewResource() && d.HasChange("renew_period") {
		update = true
		request["RenewalPeriod"] = d.Get("renew_period")
	}

	if v, ok := d.GetOk("renewal_period_unit"); ok {
		request["RenewalPeriodUnit"] = v
	}
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		request["SubscriptionType"] = d.Get("payment_type")
	}

	request["ProductCode"] = "kms"
	request["ProductType"] = "kms_ddi_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "kms"
		request["ProductType"] = "kms_ppi_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "kms"
		request["ProductType"] = "kms_ddi_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "kms"
			request["ProductType"] = "kms_ppi_public_intl"
		}
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductCode"] = "kms"
					request["ProductType"] = "kms_ddi_public_intl"
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductCode"] = "kms"
						request["ProductType"] = "kms_ppi_public_intl"
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "SetDeletionProtection"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KmsInstanceId"] = d.Id()

	if d.HasChange("deletion_protection") {
		update = true

		if v, ok := d.GetOkExists("deletion_protection"); ok {
			request["EnableDeletionProtection"] = v
		}
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		kmsServiceV2 := KmsServiceV2{client}
		if err := kmsServiceV2.SetResourceTags(d, "instance"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudKmsInstanceRead(d, meta)
}

func resourceAliCloudKmsInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	enableDelete := false
	if v, ok := d.GetOkExists("payment_type"); ok {
		if helper.InArray(fmt.Sprint(v), []string{"Subscription"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		action := "RefundInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["InstanceId"] = d.Id()

		request["ClientToken"] = helper.BuildClientToken(action)

		request["ImmediatelyRelease"] = "1"
		var endpoint string
		request["ProductCode"] = "kms"
		request["ProductType"] = "kms_ddi_public_cn"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "kms"
			request["ProductType"] = "kms_ppi_public_cn"
		}
		if client.IsInternationalAccount() {
			request["ProductCode"] = "kms"
			request["ProductType"] = "kms_ddi_public_intl"
			if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
				request["ProductCode"] = "kms"
				request["ProductType"] = "kms_ppi_public_intl"
			}
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductCode"] = "kms"
					request["ProductType"] = "kms_ddi_public_intl"
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductCode"] = "kms"
						request["ProductType"] = "kms_ppi_public_intl"
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ResourceNotExists"}) || errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		kmsServiceV2 := KmsServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 30*time.Second, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "$.InstanceId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	enableDelete = false
	if v, ok := d.GetOkExists("payment_type"); ok {
		if helper.InArray(fmt.Sprint(v), []string{"PayAsYouGo"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		action := "ReleaseKmsInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["KmsInstanceId"] = d.Id()

		if v, ok := d.GetOk("force_delete_without_backup"); ok {
			request["ForceDeleteWithoutBackup"] = v
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

	}
	return nil
}

// kmsInstanceIdsOrderDiffSuppressFunc suppresses order-only changes to vswitch_ids and
// zone_ids. The API does not report their order back, and both fields are ForceNew, so
// without this an order-only difference would destroy and recreate a live instance.
func kmsInstanceIdsOrderDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	oldValue, newValue := d.GetChange(strings.SplitN(k, ".", 2)[0])
	return kmsInstanceEqualIgnoringOrder(oldValue, newValue)
}

func kmsInstanceEqualIgnoringOrder(old, new interface{}) bool {
	oldList := convert.ConvertToInterfaceArray(old)
	newList := convert.ConvertToInterfaceArray(new)
	if len(oldList) != len(newList) {
		return false
	}

	counter := make(map[string]int, len(oldList))
	for _, item := range oldList {
		counter[fmt.Sprint(item)]++
	}
	for _, item := range newList {
		key := fmt.Sprint(item)
		counter[key]--
		if counter[key] < 0 {
			return false
		}
	}
	return true
}

func convertKmsInstanceKmsInstanceChargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PREPAY":
		return "Subscription"
	case "POSTPAY":
		return "PayAsYouGo"
	}
	return source
}
