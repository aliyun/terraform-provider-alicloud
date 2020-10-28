package alicloud

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				ValidateFunc: validation.StringLenBetween(5, 67),
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"function_args": {
				Type:     schema.TypeSet,
				Set:      expirationCdnDomainConfigHash,
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
	cdnService := &CdnService{client: client}

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
	bytconfig, _ := json.Marshal(config)

	request := cdn.CreateBatchSetCdnDomainConfigRequest()
	request.RegionId = client.RegionId
	request.DomainNames = d.Get("domain_name").(string)
	request.Functions = string(bytconfig)
	raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.BatchSetCdnDomainConfig(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "cdn_domain_config", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(fmt.Sprintf("%s:%s", request.DomainNames, d.Get("function_name").(string)))

	err = cdnService.WaitForCdnDomain(d.Get("domain_name").(string), Online, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudCdnDomainConfigRead(d, meta)
}

func resourceAlicloudCdnDomainConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnService := &CdnService{client: client}

	config, err := cdnService.DescribeCdnDomainConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	var funArgs []map[string]string

	for _, args := range config.FunctionArgs.FunctionArg {
		// This two function args is extra, filter them to pass test check.
		if args.ArgName == "aliyun_id" || args.ArgName == "scheme_origin_port" {
			continue
		}
		// private_oss_tbl always is changed and used to enable Alibaba Cloud OSS Private Bucket Back to Source Authorization
		if args.ArgName == "private_oss_tbl" {
			continue
		}
		funArgs = append(funArgs, map[string]string{
			"arg_name":  args.ArgName,
			"arg_value": args.ArgValue,
		})
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("domain_name", parts[0])
	d.Set("function_name", parts[1])
	d.Set("function_args", funArgs)

	return nil
}

func resourceAlicloudCdnDomainConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	cdnService := &CdnService{client: client}
	request := cdn.CreateDeleteSpecificConfigRequest()
	config, err := cdnService.DescribeCdnDomainConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.ConfigId = config.ConfigId
	request.DomainName = parts[0]
	raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DeleteSpecificConfig(request)
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(cdnService.WaitForCdnDomain(d.Id(), Deleted, DefaultTimeout))
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
	return hashcode.String(buf.String())
}
