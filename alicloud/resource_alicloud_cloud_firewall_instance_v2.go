package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallInstanceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallInstanceV2Create,
		Read:   resourceAliCloudCloudFirewallInstanceV2Read,
		Update: resourceAliCloudCloudFirewallInstanceV2Update,
		Delete: resourceAliCloudCloudFirewallInstanceV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cfw_log": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Upgrade", "Downgrade"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"product_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"cfw"}, false),
			},
			"product_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"cfw_elasticity_public_cn", "cfw_elasticity_public_intl", "cfw_sub_public_cn", "cfw_sub_public_intl"}, false),
			},
			"release_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renewal_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"renewal_duration_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sdl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"spec": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"payg_version", "premium_version", "enterprise_version", "ultimate_version"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallInstanceV2Create(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOkExists("sdl"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "cfw_ndlp_enable",
			"Value": fmt.Sprint(v),
		})
	}
	if v, ok := d.GetOkExists("cfw_log"); ok {
		if d.Get("payment_type").(string) == "Subscription" {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "cfw_log",
				"Value": fmt.Sprint(v),
			})
		} else {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "CfwLog",
				"Value": fmt.Sprint(v),
			})
		}
	}
	if v, ok := d.GetOk("spec"); ok {
		if d.Get("payment_type").(string) == "PayAsYouGo" {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "spec",
				"Value": convertCloudFirewallInstanceSpecRequest(v),
			})
		} else {
			parameterMapList = append(parameterMapList, map[string]interface{}{
				"Code":  "cfw_spec",
				"Value": convertCloudFirewallInstanceSpecRequest(v),
			})
		}
	}
	request["Parameter"] = parameterMapList

	request["ProductCode"] = d.Get("product_code")
	request["SubscriptionType"] = d.Get("payment_type")
	request["ProductType"] = d.Get("product_type")

	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOkExists("renewal_duration"); ok {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	var endpoint string
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable", NotFoundArticle}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_instance_v2", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudFirewallServiceV2.CloudFirewallInstanceStateRefreshFuncWithApi(d.Id(), "$.InstanceStatus", []string{}, cloudFirewallServiceV2.DescribeInstanceDescribeUserBuyVersion))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCloudFirewallInstanceV2Update(d, meta)
}

func resourceAliCloudCloudFirewallInstanceV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_instance_v2 DescribeCloudFirewallInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("payment_type", objectRaw["SubscriptionType"])
	d.Set("product_code", objectRaw["ProductCode"])
	d.Set("product_type", objectRaw["ProductType"])
	d.Set("release_time", objectRaw["ReleaseTime"])
	d.Set("renewal_duration", objectRaw["RenewalDuration"])
	d.Set("renewal_duration_unit", objectRaw["RenewalDurationUnit"])
	d.Set("renewal_status", objectRaw["RenewStatus"])

	objectRaw, err = cloudFirewallServiceV2.DescribeInstanceDescribeUserBuyVersion(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("cfw_log", objectRaw["LogStatus"])
	d.Set("sdl", convertCloudFirewallInstanceSdlResponse(objectRaw["Sdl"]))
	d.Set("spec", convertCloudFirewallInstanceVersionResponse(objectRaw["Version"]))
	d.Set("status", objectRaw["InstanceStatus"])
	d.Set("user_status", fmt.Sprint(objectRaw["UserStatus"]))

	return nil
}

func resourceAliCloudCloudFirewallInstanceV2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyCfwInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	updateList := make([]map[string]interface{}, 0)

	objectRaw, err := cloudFirewallServiceV2.DescribeInstanceDescribeUserBuyVersion(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if fmt.Sprint(convertCloudFirewallInstanceSdlResponse(objectRaw["Sdl"])) != fmt.Sprint(d.Get("sdl")) {
		update = true
	}
	if v, ok := d.GetOkExists("sdl"); ok {
		updateList = append(updateList, map[string]interface{}{
			"Code":  "Sdl",
			"Value": convertCloudFirewallInstanceSdlRequest(v),
		})
	}

	request["UpdateList"] = updateList

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
	update = false
	action = "ModifyInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	var endpoint string

	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	request["ProductType"] = d.Get("product_type")
	request["ProductCode"] = d.Get("product_code")
	request["SubscriptionType"] = d.Get("payment_type")

	if v, ok := d.GetOk("modify_type"); ok {
		request["ModifyType"] = v
	}

	parameterMapList := make([]map[string]interface{}, 0)

	if !d.IsNewResource() && d.HasChange("cfw_log") {
		update = true
	}
	if v, ok := d.GetOkExists("cfw_log"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwLog",
			"Value": v,
		})
	}

	request["Parameter"] = parameterMapList

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			if fmt.Sprint(response["Code"]) == "SYSTEM.CONCURRENT_OPERATE" {
				wait()
				return resource.RetryableError(fmt.Errorf("%s", response))
			}

			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
		if d.HasChange("cfw_log") {
			stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("cfw_log"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cloudFirewallServiceV2.CloudFirewallInstanceStateRefreshFuncWithApi(d.Id(), "LogStatus", []string{}, cloudFirewallServiceV2.DescribeInstanceDescribeUserBuyVersion))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		bssOpenApiService := BssOpenApiService{client}
		stateConf := BuildStateConf([]string{}, []string{"Paid"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, bssOpenApiService.CloudFirewallInstanceOrderDetailStateRefreshFunc(fmt.Sprint(response["Data"].(map[string]interface{})["OrderId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "SetRenewal"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()

	request["ProductCode"] = d.Get("product_code")
	request["ProductType"] = d.Get("product_type")
	request["SubscriptionType"] = d.Get("payment_type")

	if d.HasChange("renewal_duration_unit") {
		update = true
	}
	if v, ok := d.GetOk("renewal_duration_unit"); ok {
		request["RenewalPeriodUnit"] = v
	}

	if !d.IsNewResource() && d.HasChange("renewal_status") {
		update = true
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}

	if !d.IsNewResource() && d.HasChange("renewal_duration") {
		update = true
	}
	if v, ok := d.GetOk("renewal_duration"); ok {
		request["RenewalPeriod"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
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
	}

	d.Partial(false)
	return resourceAliCloudCloudFirewallInstanceV2Read(d, meta)
}

func resourceAliCloudCloudFirewallInstanceV2Delete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type").(string) == "Subscription" {
		log.Printf("[WARN] Cannot destroy resourceAliCloudCloudFirewallInstance. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	enableDelete := false
	if v, ok := d.GetOkExists("payment_type"); ok {
		if InArray(fmt.Sprint(v), []string{"PayAsYouGo"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		action := "ReleasePostInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		var endpoint string
		request = make(map[string]interface{})
		request["InstanceId"] = d.Id()

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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

		cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudFirewallServiceV2.CloudFirewallInstanceStateRefreshFunc(d.Id(), "#$.Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}
	return nil
}

func convertCloudFirewallInstanceVersionResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "2":
		return "premium_version"
	case "3":
		return "enterprise_version"
	case "4":
		return "ultimate_version"
	case "10":
		return "payg_version"
	}

	return source
}

func convertCloudFirewallInstanceSdlRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "true":
		return 1
	case "false":
		return 0
	}

	return source
}

func convertCloudFirewallInstanceSdlResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "0":
		return false
	case "1":
		return true
	}

	return source
}
