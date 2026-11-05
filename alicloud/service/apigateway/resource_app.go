package apigateway

import (
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceApp() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAppCreate,
		Read:   ResourceAppRead,
		Update: ResourceAppUpdate,
		Delete: ResourceAppDelete,
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
			"tags": tags.TagsSchema(),
		},
	}
}

func ResourceAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cloudapi.CreateCreateAppRequest()
	request.RegionId = client.RegionId
	request.AppName = d.Get("name").(string)
	if v, exist := d.GetOk("description"); exist {
		request.Description = v.(string)
	}
	request.Description = d.Get("description").(string)

	if err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateApp(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"RepeatedCommit"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.CreateAppResponse)
		d.SetId(strconv.FormatInt(response.AppId, 10))
		return nil
	}); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_apigateway_app", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return ResourceAppUpdate(d, meta)
}

func ResourceAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	if err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		tags, err := cloudApiService.DescribeTags(d.Id(), nil, tags.TagResourceApp)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotFoundResourceId"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		d.Set("tags", cloudApiService.tagsToMap(tags))
		return nil
	}); err != nil {
		return errs.WrapError(err)
	}

	if err := retry.Retry(3*time.Second, func() *retry.RetryError {
		object, err := cloudApiService.DescribeApiGatewayApp(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		d.Set("name", object.AppName)
		d.Set("description", object.Description)
		return nil
	}); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func ResourceAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	if err := cloudApiService.setInstanceTags(d, tags.TagResourceApp); err != nil {
		return errs.WrapError(err)
	}
	if d.IsNewResource() {
		return ResourceAppRead(d, meta)
	}
	if d.HasChange("name") || d.HasChange("description") {
		request := cloudapi.CreateModifyAppRequest()
		request.RegionId = client.RegionId
		request.AppId = requests.Integer(d.Id())
		request.AppName = d.Get("name").(string)
		if v, exist := d.GetOk("description"); exist {
			request.Description = v.(string)
		}

		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApp(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	time.Sleep(3 * time.Second)
	return ResourceAppRead(d, meta)
}

func ResourceAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	request := cloudapi.CreateDeleteAppRequest()
	request.RegionId = client.RegionId
	request.AppId = requests.Integer(d.Id())

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApp(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFoundApp"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errs.WrapError(cloudApiService.WaitForApiGatewayApp(d.Id(), Deleted, DefaultTimeout))
}
