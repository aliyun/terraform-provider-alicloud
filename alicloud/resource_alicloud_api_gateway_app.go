package alicloud

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
		},
	}
}

func resourceAliyunApigatewayAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cloudapi.CreateCreateAppRequest()
	request.AppName = d.Get("name").(string)
	if v, exist := d.GetOk("description"); exist {
		request.Description = v.(string)
	}
	request.Description = d.Get("description").(string)

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateApp(request)
		})
		if err != nil {
			if IsExceptedError(err, RepeatedCommit) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cloudapi.CreateAppResponse)
		d.SetId(strconv.FormatInt(response.AppId, 10))
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apigateway_app", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAliyunApigatewayAppRead(d, meta)
}

func resourceAliyunApigatewayAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	request := cloudapi.CreateDescribeAppRequest()
	request.AppId = requests.Integer(d.Id())

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
		return nil
	}); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAliyunApigatewayAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("name") || d.HasChange("description") {
		request := cloudapi.CreateModifyAppRequest()
		request.AppId = requests.Integer(d.Id())
		request.AppName = d.Get("name").(string)
		if v, exist := d.GetOk("description"); exist {
			request.Description = v.(string)
		}

		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApp(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}
	time.Sleep(3 * time.Second)
	return resourceAliyunApigatewayAppRead(d, meta)
}

func resourceAliyunApigatewayAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	request := cloudapi.CreateDeleteAppRequest()
	request.AppId = requests.Integer(d.Id())

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApp(request)
	})
	if err != nil {
		if IsExceptedError(err, NotFoundApp) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return WrapError(cloudApiService.WaitForApiGatewayApp(d.Id(), Deleted, DefaultTimeout))
}
