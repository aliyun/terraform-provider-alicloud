package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// apiGatewaySlsLogStorePrefix is the prefix the CloudAPI service prepends to a
// SlsLogStore value that does not already start with it. CreateLogConfig accepts
// the bare name while DescribeLogConfig returns the canonical (prefixed) value,
// causing an input/output mismatch and a perpetual diff. Read strips the prefix
// so state matches the configured value, and the diff suppress func keeps plans
// stable for users who configure the value with the prefix already set.
const apiGatewaySlsLogStorePrefix = "acs_apigateway-"

// normalizeApiGatewaySlsLogStore removes the service-side prefix from the
// SlsLogStore value returned by DescribeLogConfig. Non-string values are
// treated as empty so an absent field does not produce a "<nil>" placeholder.
func normalizeApiGatewaySlsLogStore(v interface{}) string {
	s, _ := v.(string)
	return strings.TrimPrefix(s, apiGatewaySlsLogStorePrefix)
}

// apiGatewayLogConfigSlsLogStoreDiffSuppressFunc treats the bare and prefixed
// forms of sls_log_store as equivalent so plans stay stable regardless of
// whether the user configured the value with or without the service prefix.
func apiGatewayLogConfigSlsLogStoreDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return normalizeApiGatewaySlsLogStore(old) == normalizeApiGatewaySlsLogStore(new)
}

func resourceAlicloudApiGatewayLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudApiGatewayLogConfigCreate,
		Read:   resourceAlicloudApiGatewayLogConfigRead,
		Update: resourceAlicloudApiGatewayLogConfigUpdate,
		Delete: resourceAlicloudApiGatewayLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"sls_project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sls_log_store": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: apiGatewayLogConfigSlsLogStoreDiffSuppressFunc,
			},
			"log_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PROVIDER"}, false),
			},
		},
	}
}

func resourceAlicloudApiGatewayLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateLogConfig"
	request := make(map[string]interface{})
	var err error

	request["SlsProject"] = d.Get("sls_project")
	request["SlsLogStore"] = d.Get("sls_log_store")
	request["LogType"] = d.Get("log_type")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_log_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["LogType"]))

	return resourceAlicloudApiGatewayLogConfigRead(d, meta)
}

func resourceAlicloudApiGatewayLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	object, err := cloudApiService.DescribeApiGatewayLogConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("sls_project", object["SlsProject"])
	d.Set("sls_log_store", normalizeApiGatewaySlsLogStore(object["SlsLogStore"]))
	d.Set("log_type", object["LogType"])

	return nil
}

func resourceAlicloudApiGatewayLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"LogType": d.Id(),
	}

	if d.HasChange("sls_project") || d.HasChange("sls_log_store") {
		update = true
	}
	request["SlsProject"] = d.Get("sls_project")
	request["SlsLogStore"] = d.Get("sls_log_store")

	if update {
		action := "ModifyLogConfig"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudApiGatewayLogConfigRead(d, meta)
}

func resourceAlicloudApiGatewayLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLogConfig"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"LogType": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
