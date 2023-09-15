// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
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
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
							Type:     schema.TypeInt,
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
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(1000, 100000),
			},
			"product_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"3"}, false),
			},
			"renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1, 36),
			},
			"renew_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
			},
			"secret_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(0, 100000),
			},
			"spec": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{1000, 2000, 4000}),
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
				Required:     true,
				ValidateFunc: IntBetween(1, 10000),
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone_ids": {
				Type:     schema.TypeList,
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
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
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
	request["Parameter"] = parameterMapList

	request["Period"] = "1"
	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOk("renew_status"); ok {
		request["RenewalStatus"] = v
	}
	request["ProductType"] = "kms_ddi_public_cn"
	request["SubscriptionType"] = "Subscription"
	request["ProductCode"] = "kms"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "kms_ddi_public_intl"
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	kmsServiceV2 := KmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(id)}, d.Timeout(schema.TimeoutCreate), 10*time.Second, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "InstanceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	action = "ConnectKmsInstance"
	conn, err = client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["KmsInstanceId"] = d.Id()

	request["VpcId"] = d.Get("vpc_id")
	request["KMProvider"] = "Aliyun"
	jsonPathResult1, err := jsonpath.Get("$", d.Get("zone_ids"))
	if err != nil {
		return WrapError(err)
	}
	request["ZoneIds"] = convertArrayToString(jsonPathResult1, ",")

	jsonPathResult2, err := jsonpath.Get("$", d.Get("vswitch_ids"))
	if err != nil {
		return WrapError(err)
	}
	request["VSwitchIds"] = convertArrayToString(jsonPathResult2, ",")

	wait = incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["KmsInstanceId"]))

	kmsServiceV2 = KmsServiceV2{client}
	stateConf = BuildStateConf([]string{}, []string{"Connected"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
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

	d.Set("ca_certificate_chain_pem", objectRaw["CaCertificateChainPem"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("key_num", objectRaw["KeyNum"])
	d.Set("secret_num", objectRaw["SecretNum"])
	d.Set("spec", objectRaw["Spec"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("vpc_num", objectRaw["VpcNum"])
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
	d.Set("bind_vpcs", bindVpcsMaps)

	d.Set("vswitch_ids", objectRaw["VswitchIds"])

	d.Set("zone_ids", objectRaw["ZoneIds"])

	return nil
}

func resourceAliCloudKmsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyInstance"
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChanges("vpc_num", "spec", "secret_num", "key_num") {
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
	request["Parameter"] = parameterMapList

	request["ProductType"] = "kms_ddi_public_cn"
	request["ProductCode"] = "kms"
	request["ModifyType"] = "Upgrade"
	request["SubscriptionType"] = "Subscription"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"NotApplicable"}) {
					conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
					request["ProductType"] = "kms_ddi_public_intl"
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		kmsServiceV2 := KmsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("key_num"))}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "KeyNum", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateKmsInstanceBindVpc"
	conn, err = client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["KmsInstanceId"] = d.Id()
	if d.HasChange("bind_vpcs") {
		update = true
		request["BindVpcs"] = "[]"
		if v, ok := d.GetOk("bind_vpcs"); ok {
			bindVpcsMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range v.(*schema.Set).List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["VpcId"] = dataLoopTmp["vpc_id"]
				dataLoopMap["VpcOwnerId"] = dataLoopTmp["vpc_owner_id"]
				dataLoopMap["RegionId"] = dataLoopTmp["region_id"]
				dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
				bindVpcsMaps = append(bindVpcsMaps, dataLoopMap)
			}
			bindVpcsMapsJson, err := json.Marshal(bindVpcsMaps)
			if err != nil {
				return WrapError(err)
			}
			request["BindVpcs"] = string(bindVpcsMapsJson)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2016-01-20"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAliCloudKmsInstanceRead(d, meta)
}

func resourceAliCloudKmsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Instance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
