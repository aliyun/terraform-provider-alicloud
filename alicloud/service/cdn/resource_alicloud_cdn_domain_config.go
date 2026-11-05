package cdn

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"encoding/json"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAliCloudCdnDomainConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCdnDomainConfigCreate,
		Read:   resourceAliCloudCdnDomainConfigRead,
		Update: resourceAliCloudCdnDomainConfigUpdate,
		Delete: resourceAliCloudCdnDomainConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(5, 67),
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"function_args": {
				Type:     schema.TypeSet,
				Set:      expirationCdnDomainConfigHash,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arg_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"arg_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"config_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCdnDomainConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnService := &CdnService{client: client}
	var response map[string]interface{}
	var err error
	action := "BatchSetCdnDomainConfig"

	config := make([]map[string]interface{}, 1)
	functionArgs := d.Get("function_args").(*schema.Set).List()
	args := make([]map[string]interface{}, len(functionArgs))
	for key, value := range functionArgs {
		arg := value.(map[string]interface{})
		args[key] = map[string]interface{}{
			"argName":  arg["arg_name"],
			"argValue": arg["arg_value"],
		}
	}
	config[0] = map[string]interface{}{
		"functionArgs": args,
		"functionName": d.Get("function_name").(string),
	}
	if v, ok := d.GetOk("parent_id"); ok {
		config[0]["parentId"] = v.(string)
	}
	bytconfig, _ := json.Marshal(config)

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"DomainNames": d.Get("domain_name").(string),
		"Functions":   string(bytconfig),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Cdn", "2018-05-10", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ServiceBusy", "FlowControlError"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cdn_domain_config", action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DomainConfigList.DomainConfigModel", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, "request.DomainNames", "$.DomainConfigList.DomainConfigModel", err)
	}

	if len(v.([]interface{})) != 1 {
		return errs.WrapErrorf(errs.NotFoundErr("CDN", d.Get("function_name").(string)), errs.DefaultErrorMsg, err)
	}

	val := v.([]interface{})[0].(map[string]interface{})
	d.SetId(fmt.Sprintf("%s:%s:%s", val["DomainName"], val["FunctionName"], val["ConfigId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, cdnService.CdnDomainConfigRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	err = cdnService.WaitForCdnDomain(d.Get("domain_name").(string), Online, DefaultTimeoutMedium)
	if err != nil {
		return errs.WrapError(err)
	}

	return resourceAliCloudCdnDomainConfigRead(d, meta)
}

func resourceAliCloudCdnDomainConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnService := &CdnService{client: client}
	var response map[string]interface{}
	var err error

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"DomainNames": parts[0],
	}

	update := false
	if d.HasChange("function_args") || d.HasChange("parent_id") {
		update = true
		config := make([]map[string]interface{}, 1)
		functionArgs := d.Get("function_args").(*schema.Set).List()
		args := make([]map[string]interface{}, len(functionArgs))
		for key, value := range functionArgs {
			arg := value.(map[string]interface{})
			args[key] = map[string]interface{}{
				"argName":  arg["arg_name"],
				"argValue": arg["arg_value"],
			}
		}

		config[0] = map[string]interface{}{
			"functionArgs": args,
			"functionName": parts[1],
			"configId":     parts[2],
		}

		if v, ok := d.GetOk("parent_id"); ok {
			config[0]["parentId"] = v
		}

		bytconfig, _ := json.Marshal(config)
		request["Functions"] = string(bytconfig)
	}

	if update {
		action := "BatchSetCdnDomainConfig"

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Cdn", "2018-05-10", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"ServiceBusy", "FlowControlError"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cdn_domain_config", action, errs.AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$.DomainConfigList.DomainConfigModel", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, "request.DomainNames", "$.DomainConfigList.DomainConfigModel", err)
		}

		if len(v.([]interface{})) > 1 {
			return errs.WrapErrorf(errs.NotFoundErr("CDN", d.Get("function_name").(string)), errs.DefaultErrorMsg, err)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, cdnService.CdnDomainConfigRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	err = cdnService.WaitForCdnDomain(d.Get("domain_name").(string), Online, DefaultTimeoutMedium)
	if err != nil {
		return errs.WrapError(err)
	}

	return resourceAliCloudCdnDomainConfigRead(d, meta)
}

func resourceAliCloudCdnDomainConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnService := &CdnService{client: client}

	v, err := cdnService.DescribeCdnDomainConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	val := v.(map[string]interface{})
	args := make([]interface{}, 0)
	if v, ok := val["FunctionArgs"].(map[string]interface{}); ok {
		if tmp, exist := v["FunctionArg"].([]interface{}); exist {
			args = tmp
		}
	}

	var funArgs []map[string]string
	ignoreFunctionArg := []string{"aliyun_id", "scheme_origin_port", "dsl", "session_timeout",
		"oss_pri_buckets", "private_oss_tbl", "private_oss_ram_unauthorized", "cert", "cert_name",
		"consistent_hash", "cert_type", "https", "dkey", "pkey"}

	for _, k := range args {

		arg := k.(map[string]interface{})
		// This two function args is extra, filter them to pass test check.
		if helper.InArray(fmt.Sprint(arg["ArgName"]), ignoreFunctionArg) {
			continue
		}

		funArgs = append(funArgs, map[string]string{
			"arg_name":  arg["ArgName"].(string),
			"arg_value": arg["ArgValue"].(string),
		})
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("domain_name", parts[0])
	d.Set("function_name", parts[1])
	d.Set("config_id", val["ConfigId"])
	if len(parts) == 2 {
		d.SetId(fmt.Sprintf("%s:%s:%s", parts[0], parts[1], val["ConfigId"]))
	}

	d.Set("status", val["Status"])
	d.Set("function_args", funArgs)
	d.Set("parent_id", val["ParentId"])

	return nil
}

func resourceAliCloudCdnDomainConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnService := &CdnService{client: client}

	_, err := cdnService.DescribeCdnDomainConfig(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapError(err)
	}

	action := "DeleteSpecificConfig"
	var response map[string]interface{}

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"DomainName": parts[0],
		"ConfigId":   parts[2],
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Cdn", "2018-05-10", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ServiceBusy", "FlowControlError"}) || helper.NeedRetry(err) {
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

	return nil
}

func expirationCdnDomainConfigHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["arg_name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["arg_value"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return helper.Hashcode(buf.String())
}
