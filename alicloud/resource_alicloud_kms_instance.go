// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_num": {
				Type:     schema.TypeString,
				Required: true,
			},
			"payment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"product_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"kms_ddi_public_cn", "kms_ddi_public_intl"}, false),
			},
			"product_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "3",
			},
			"renew_period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(1, 36),
			},
			"renew_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "AutoRenewal",
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
			},
			"secret_num": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringLenBetween(0, 100000),
			},
			"spec": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"1000", "2000", "4000"}, false),
			},
			"vpc_num": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringLenBetween(1, 10000),
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

	request["Parameter"] = []map[string]string{
		{
			"Code":  "Region",
			"Value": client.RegionId,
		},
		{
			"Code":  "VpcNum",
			"Value": d.Get("vpc_num").(string),
		},
		{
			"Code":  "Spec",
			"Value": d.Get("spec").(string),
		},
		{
			"Code":  "SecretNum",
			"Value": d.Get("secret_num").(string),
		},
		{
			"Code":  "KeyNum",
			"Value": d.Get("key_num").(string),
		},
		{
			"Code":  "ProductVersion",
			"Value": d.Get("product_version").(string),
		},
	}

	request["Period"] = "1"
	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOk("renew_status"); ok {
		request["RenewalStatus"] = v
	}
	request["ProductType"] = d.Get("product_type")
	request["ProductCode"] = "kms"
	request["SubscriptionType"] = "Subscription"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
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

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	kmsServiceV2 := KmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutCreate), 10*time.Second, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "KmsInstanceId", []string{}))
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

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("key_num", objectRaw["KeyLimits"])
	d.Set("spec", objectRaw["QpsLimits"])
	d.Set("vpc_num", objectRaw["BindVpcLimits"])
	secretManagerParameters1RawObj, _ := jsonpath.Get("$.SecretManagerParameters", objectRaw)
	secretManagerParameters1Raw := make(map[string]interface{})
	if secretManagerParameters1RawObj != nil {
		secretManagerParameters1Raw = secretManagerParameters1RawObj.(map[string]interface{})
	}
	d.Set("secret_num", secretManagerParameters1Raw["SecretLimits"])

	return nil
}

func resourceAliCloudKmsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "ModifyInstance"
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("") {
		update = true
	}
	if d.HasChange("") {
		update = true
	}
	if d.HasChange("") {
		update = true
	}
	if d.HasChange("") {
		update = true
	}
	request["Parameter"] = []map[string]string{
		{
			"Code":  "VpcNum",
			"Value": d.Get("vpc_num").(string),
		},
		{
			"Code":  "Spec",
			"Value": d.Get("spec").(string),
		},
		{
			"Code":  "SecretNum",
			"Value": d.Get("secret_num").(string),
		},
		{
			"Code":  "KeyNum",
			"Value": d.Get("key_num").(string),
		},
	}

	request["ProductType"] = d.Get("product_type")
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
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		kmsServiceV2 := KmsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, kmsServiceV2.KmsInstanceStateRefreshFunc(d.Id(), "KeyLimits", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudKmsInstanceRead(d, meta)
}

func resourceAliCloudKmsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Instance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
