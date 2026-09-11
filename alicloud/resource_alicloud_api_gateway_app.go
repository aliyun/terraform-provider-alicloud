package alicloud

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunApigatewayApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayAppCreate,
		Read:   resourceAliyunApigatewayAppRead,
		Update: resourceAliyunApigatewayAppUpdate,
		Delete: resourceAliyunApigatewayAppDelete,
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
			"tags": tagsSchema(),

			// app_code: settable at create; modifiable post-create via ResetAppCode;
			// refreshed from DescribeAppSecurity. E-TF-0006: regular modifiable attribute.
			"app_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// extend: settable at create and update (cspec @rac=["create","update"]);
			// refreshed from DescribeApp.
			"extend": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// disabled: modifiable via ModifyApp (cspec @rac=["update"]);
			// CreateApp does not accept Disabled, so it is applied via ModifyApp right after create.
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			// app_key: server-managed AppKey surfaced from DescribeAppSecurity;
			// regenerated server-side by ResetAppSecret (triggered via
			// app_secret_reset). Exported attribute (Computed-only), not settable
			// or modifiable via config.
			"app_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// app_secret: server-managed AppSecret surfaced from DescribeAppSecurity;
			// regenerated server-side by ResetAppSecret (triggered via
			// app_secret_reset). Exported attribute (Computed-only), not settable
			// or modifiable via config.
			"app_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			// app_secret_reset: trigger field. Any non-empty change invokes the
			// ResetAppSecret operation; the new AppKey/AppSecret are refreshed into state.
			"app_secret_reset": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliyunApigatewayAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cloudapi.CreateCreateAppRequest()
	request.RegionId = client.RegionId
	request.AppName = d.Get("name").(string)
	if v, exist := d.GetOk("description"); exist {
		request.Description = v.(string)
	}
	if v, exist := d.GetOk("app_code"); exist {
		request.AppCode = v.(string)
	}
	if v, exist := d.GetOk("extend"); exist {
		request.Extend = v.(string)
	}
	// app_key/app_secret are exported (Computed-only) credentials surfaced from
	// DescribeAppSecurity and regenerated server-side via app_secret_reset; they
	// are not forwarded to CreateApp.
	// Disabled is not accepted by CreateApp (cspec @rac=["update"]); applied via ModifyApp below.

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateApp(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RepeatedCommit"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.CreateAppResponse)
		d.SetId(strconv.FormatInt(response.AppId, 10))
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_app", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return resourceAliyunApigatewayAppUpdate(d, meta)
}

func resourceAliyunApigatewayAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		tags, err := cloudApiService.DescribeTags(d.Id(), nil, TagResourceApp)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFoundResourceId"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("tags", cloudApiService.tagsToMap(tags))
		return nil
	}); err != nil {
		return WrapError(err)
	}

	if err := resource.Retry(3*time.Second, func() *resource.RetryError {
		object, err := cloudApiService.DescribeApiGatewayApp(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		d.Set("name", object.AppName)
		d.Set("description", object.Description)
		d.Set("extend", object.Extend)
		// Disabled is returned by DescribeApp but is absent from the SDK
		// DescribeAppResponse struct, so it is parsed from the raw HTTP
		// response body. This mirrors the ModifyApp write-side path, which
		// injects Disabled via QueryParams because ModifyAppRequest also
		// lacks the field. The raw body is a JSON object; Unmarshal into a
		// map and type-assert the bool field. jsonpath.Get cannot take a raw
		// string (it expects a parsed map/interface{}) and silently swallows
		// the "unsupported value type string" error, so d.Set never runs.
		var raw map[string]interface{}
		if e := json.Unmarshal([]byte(object.GetHttpContentString()), &raw); e == nil {
			if v, ok := raw["Disabled"].(bool); ok {
				d.Set("disabled", v)
			}
		}
		return nil
	}); err != nil {
		return WrapError(err)
	}

	// Read app_key/app_secret/app_code via DescribeAppSecurity. The SDK response
	// struct carries AppCode/AppKey/AppSecret/ModifiedTime/CreatedTime, matching
	// the cspec App_operation_DescribeAppSecurity_mapping.
	if err := resource.Retry(3*time.Second, func() *resource.RetryError {
		object, err := cloudApiService.DescribeApiGatewayAppSecurity(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("app_key", object.AppKey)
		d.Set("app_secret", object.AppSecret)
		d.Set("app_code", object.AppCode)
		return nil
	}); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAliyunApigatewayAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	if err := cloudApiService.setInstanceTags(d, TagResourceApp); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		// CreateApp does not accept Disabled; apply via ModifyApp when the user set it.
		if _, ok := d.GetOkExists("disabled"); ok {
			if err := resourceAliyunApigatewayAppModifyApp(d, meta); err != nil {
				return err
			}
		}
		d.Partial(false)
		return resourceAliyunApigatewayAppRead(d, meta)
	}

	// ResetAppSecret trigger: any non-empty change to app_secret_reset invokes
	// ResetAppSecret, after which Read refreshes AppKey/AppSecret.
	if d.HasChange("app_secret_reset") {
		_, n := d.GetChange("app_secret_reset")
		if v := n.(string); v != "" {
			if appKey, ok := d.GetOk("app_key"); ok && appKey.(string) != "" {
				request := cloudapi.CreateResetAppSecretRequest()
				request.RegionId = client.RegionId
				request.AppKey = appKey.(string)
				if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
						return cloudApiClient.ResetAppSecret(request)
					})
					if err != nil {
						if IsExpectedErrors(err, []string{"NotFoundApp"}) {
							return resource.NonRetryableError(err)
						}
						return resource.RetryableError(err)
					}
					addDebug(request.GetActionName(), raw, request.RpcRequest, request)
					return nil
				}); err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	// ResetAppCode: when app_code HasChange post-create, call ResetAppCode with
	// old/new values. ResetAppCode requires the existing AppCode as identifier;
	// skip when old value is empty (no AppCode to reset from).
	if d.HasChange("app_code") {
		o, n := d.GetChange("app_code")
		oldVal := o.(string)
		newVal := n.(string)
		if oldVal != "" {
			request := cloudapi.CreateResetAppCodeRequest()
			request.RegionId = client.RegionId
			request.AppCode = oldVal
			if newVal != "" {
				request.NewAppCode = newVal
			}
			if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
					return cloudApiClient.ResetAppCode(request)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"NotFoundApp"}) {
						return resource.NonRetryableError(err)
					}
					return resource.RetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}
	}

	// ModifyApp: name/description/extend/disabled changes. Drift guards for
	// Optional+Computed fields: extend only when user set non-empty value,
	// disabled only when user explicitly set it.
	modifyNeeded := d.HasChange("name") || d.HasChange("description")
	if v, exist := d.GetOk("extend"); exist && v.(string) != "" {
		modifyNeeded = true
	}
	if _, ok := d.GetOkExists("disabled"); ok {
		modifyNeeded = true
	}
	if modifyNeeded {
		if err := resourceAliyunApigatewayAppModifyApp(d, meta); err != nil {
			return err
		}
	}
	time.Sleep(3 * time.Second)
	return resourceAliyunApigatewayAppRead(d, meta)
}

// resourceAliyunApigatewayAppModifyApp issues a ModifyApp call. The SDK
// ModifyAppRequest does not carry a Disabled field (cspec M-RT-0040 added
// Disabled to the ModifyApp input mapping), so Disabled is injected via
// QueryParams, mirroring the pattern used elsewhere for SDK-lagging fields.
func resourceAliyunApigatewayAppModifyApp(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cloudapi.CreateModifyAppRequest()
	request.RegionId = client.RegionId
	request.AppId = requests.Integer(d.Id())
	request.AppName = d.Get("name").(string)
	if v, exist := d.GetOk("description"); exist {
		request.Description = v.(string)
	}
	if v, exist := d.GetOk("extend"); exist && v.(string) != "" {
		request.Extend = v.(string)
	}
	if v, ok := d.GetOkExists("disabled"); ok {
		if v.(bool) {
			request.QueryParams["Disabled"] = "true"
		} else {
			request.QueryParams["Disabled"] = "false"
		}
	}
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.ModifyApp(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func resourceAliyunApigatewayAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	request := cloudapi.CreateDeleteAppRequest()
	request.RegionId = client.RegionId
	request.AppId = requests.Integer(d.Id())

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApp(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundApp"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(cloudApiService.WaitForApiGatewayApp(d.Id(), Deleted, DefaultTimeout))
}
