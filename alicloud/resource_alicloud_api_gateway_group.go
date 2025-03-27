package alicloud

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunApigatewayGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayGroupCreate,
		Read:   resourceAliyunApigatewayGroupRead,
		Update: resourceAliyunApigatewayGroupUpdate,
		Delete: resourceAliyunApigatewayGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_intranet_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"base_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_log_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"request_body": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"response_body": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"query_string": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"request_headers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"response_headers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"jwt_claims": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunApigatewayGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cloudapi.CreateCreateApiGroupRequest()
	request.RegionId = client.RegionId
	request.GroupName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.BasePath = d.Get("base_path").(string)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = v.(string)
	}
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateApiGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RepeatedCommit"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.CreateApiGroupResponse)
		d.SetId(response.GroupId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAliyunApigatewayGroupUpdate(d, meta)
}

func resourceAliyunApigatewayGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	apiGroup, err := cloudApiService.DescribeApiGatewayGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", apiGroup.GroupName)
	d.Set("description", apiGroup.Description)
	d.Set("sub_domain", apiGroup.SubDomain)
	d.Set("vpc_domain", apiGroup.VpcDomain)
	if apiGroup.VpcDomain != "" {
		d.Set("vpc_intranet_enable", true)
	} else {
		d.Set("vpc_intranet_enable", false)
	}
	d.Set("instance_id", apiGroup.InstanceId)
	d.Set("base_path", apiGroup.BasePath)
	if apiGroup.UserLogConfig != "" {
		var logConfig ApiGatewayUserLogConfig
		err := json.Unmarshal([]byte(apiGroup.UserLogConfig), &logConfig)
		if err != nil {
			return WrapError(err)
		}
		userLogConfig := map[string]interface{}{}
		userLogConfig["request_body"] = logConfig.RequestBody
		userLogConfig["response_body"] = logConfig.ResponseBody
		userLogConfig["query_string"] = logConfig.QueryString
		userLogConfig["request_headers"] = logConfig.RequestHeaders
		userLogConfig["response_headers"] = logConfig.ResponseHeaders
		userLogConfig["jwt_claims"] = logConfig.JwtClaims
		if err := d.Set("user_log_config", []map[string]interface{}{userLogConfig}); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func resourceAliyunApigatewayGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	request := cloudapi.CreateModifyApiGroupRequest()
	request.RegionId = client.RegionId
	request.GroupId = d.Id()

	if !d.IsNewResource() && d.HasChanges("name", "description", "base_path") {
		request.Description = d.Get("description").(string)
		request.GroupName = d.Get("name").(string)
		request.BasePath = d.Get("base_path").(string)
		update = true
	}
	if d.HasChanges("user_log_config") {
		logConfig, err := userLogConfigToJsonStr(d)
		if err != nil {
			return WrapError(err)
		}
		request.UserLogConfig = logConfig
		update = true
	}

	if update {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApiGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	update = false
	request2 := cloudapi.CreateModifyIntranetDomainPolicyRequest()
	request2.RegionId = client.RegionId
	request2.GroupId = d.Id()
	if d.HasChanges("vpc_intranet_enable") {
		request2.VpcIntranetEnable = requests.NewBoolean(d.Get("vpc_intranet_enable").(bool))
		update = true
	}
	if update {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyIntranetDomainPolicy(request2)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request2.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request2.GetActionName(), raw, request2.RpcRequest, request2)
	}

	return resourceAliyunApigatewayGroupRead(d, meta)
}

func resourceAliyunApigatewayGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	request := cloudapi.CreateDeleteApiGroupRequest()
	request.RegionId = client.RegionId
	request.GroupId = d.Id()

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApiGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(cloudApiService.WaitForApiGatewayGroup(d.Id(), Deleted, DefaultTimeout))

}

func userLogConfigToJsonStr(d *schema.ResourceData) (string, error) {
	var logConfig ApiGatewayUserLogConfig
	var l []interface{}
	v, ok := d.GetOk("user_log_config")
	if ok {
		l = v.([]interface{})
		config := l[0].(map[string]interface{})
		if val, exist := config["request_body"]; exist {
			logConfig.RequestBody = val.(bool)
		}
		if val, exist := config["response_body"]; exist {
			logConfig.ResponseBody = val.(bool)
		}
		if val, exist := config["query_string"]; exist {
			logConfig.QueryString = val.(string)
		}
		if val, exist := config["request_headers"]; exist {
			logConfig.RequestHeaders = val.(string)
		}
		if val, exist := config["response_headers"]; exist {
			logConfig.ResponseHeaders = val.(string)
		}
		if val, exist := config["jwt_claims"]; exist {
			logConfig.JwtClaims = val.(string)
		}
	}

	configStr, err := json.Marshal(logConfig)
	if err != nil {
		return "", WrapError(err)
	}
	return string(configStr), nil
}
