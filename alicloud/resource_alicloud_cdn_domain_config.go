package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCdnDomainConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCdnDomainConfigCreate,
		Read:   resourceAlicloudCdnDomainConfigRead,
		Delete: resourceAlicloudCdnDomainConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDomainName,
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"function_args": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
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
		},
	}
}

func resourceAlicloudCdnDomainConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	config := make([]map[string]interface{}, 1)
	functionArgs := d.Get("function_args").([]interface{})
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
	bytconfig, _ := json.Marshal(config)

	request := cdn.CreateBatchSetCdnDomainConfigRequest()
	request.DomainNames = d.Get("domain_name").(string)
	request.Functions = string(bytconfig)
	raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.BatchSetCdnDomainConfig(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "cdn_domain_config", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	d.SetId(fmt.Sprintf("%s:%s", request.DomainNames, d.Get("function_name").(string)))

	cdnservice := &CdnService{client: client}
	err = cdnservice.WaitForCdnDomain(d.Get("domain_name").(string), Online, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudCdnDomainConfigRead(d, meta)
}

func resourceAlicloudCdnDomainConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	cdnservice := &CdnService{client: client}
	config, err := cdnservice.DescribeDomainConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	funArgs := make([]map[string]string, len(config.FunctionArgs.FunctionArg))
	for key, args := range config.FunctionArgs.FunctionArg {
		funArgs[key] = map[string]string{
			"arg_name":  args.ArgName,
			"arg_value": args.ArgValue,
		}
	}
	d.Set("function_name", config.FunctionName)
	d.Set("function_args", funArgs)

	return nil
}

func resourceAlicloudCdnDomainConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	cdnservice := &CdnService{client: client}
	request := cdn.CreateDeleteSpecificConfigRequest()
	config, err := cdnservice.DescribeDomainConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	request.ConfigId = config.ConfigId
	request.DomainName = d.Get("domain_name").(string)
	raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DeleteSpecificConfig(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return nil
}
