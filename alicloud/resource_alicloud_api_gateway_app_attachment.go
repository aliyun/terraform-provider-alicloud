package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunApigatewayAppAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayAppAttachmentCreate,
		Read:   resourceAliyunApigatewayAppAttachmentRead,
		Delete: resourceAliyunApigatewayAppAttachmentDelete,

		Schema: map[string]*schema.Schema{

			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"stage_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(StageNamePre), string(StageNameRelease), string(StageNameTest)}),
			},
		},
	}
}

func resourceAliyunApigatewayAppAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	apiId := d.Get("api_id").(string)
	groupId := d.Get("group_id").(string)
	stageName := d.Get("stage_name").(string)
	appId := d.Get("app_id").(string)

	request := cloudapi.CreateSetAppsAuthoritiesRequest()
	request.GroupId = groupId
	request.ApiId = apiId
	request.AppIds = appId
	request.StageName = stageName

	_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.SetAppsAuthorities(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apigateway_app_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	id := fmt.Sprintf("%s%s%s%s%s%s%s", groupId, COLON_SEPARATED, apiId, COLON_SEPARATED, appId, COLON_SEPARATED, stageName)

	err = cloudApiService.WaitForApiGatewayAppAttachment(id, Normal, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}

	d.SetId(id)
	return resourceAliyunApigatewayAppAttachmentRead(d, meta)
}

func resourceAliyunApigatewayAppAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	_, err := cloudApiService.DescribeApiGatewayAppAttachment(d.Id())
	if err != nil {
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("group_id", parts[0])
	d.Set("api_id", parts[1])
	d.Set("app_id", parts[2])
	d.Set("stage_name", parts[3])

	return nil
}

func resourceAliyunApigatewayAppAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	request := cloudapi.CreateRemoveAppsAuthoritiesRequest()
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	request.GroupId = parts[0]
	request.ApiId = parts[1]
	request.AppIds = parts[2]
	request.StageName = parts[3]

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.RemoveAppsAuthorities(request)
	})
	if err != nil {
		if IsExceptedError(err, NotFoundAuthorization) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	return WrapError(cloudApiService.WaitForApiGatewayAppAttachment(d.Id(), Deleted, DefaultLongTimeout))
}
