package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

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

			"app_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"api_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"stage_name": &schema.Schema{
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

	authorizationReq := cloudapi.CreateSetAppsAuthoritiesRequest()
	authorizationReq.GroupId = groupId
	authorizationReq.ApiId = apiId
	authorizationReq.AppIds = appId
	authorizationReq.StageName = stageName

	_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.SetAppsAuthorities(authorizationReq)
	})
	if err != nil {
		return fmt.Errorf("APP Attachment: Authorizing api to app got an error: %#v.", err)
	}

	id := fmt.Sprintf("%s%s%s%s%s%s%s", groupId, COLON_SEPARATED, apiId, COLON_SEPARATED, appId, COLON_SEPARATED, stageName)

	err = cloudApiService.WaitForAppAttachmentAuthorization(id, 10)
	if err != nil {
		return fmt.Errorf("APP Attachment: Authorizing api to app got an error: %#v.", err)
	}

	d.SetId(id)
	return resourceAliyunApigatewayAppAttachmentRead(d, meta)
}

func resourceAliyunApigatewayAppAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	_, err := cloudApiService.DescribeAuthorization(d.Id())
	if err != nil {
		return fmt.Errorf("Describe apigatway App error: %#v", err)
	}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	d.Set("group_id", split[0])
	d.Set("api_id", split[1])
	d.Set("app_id", split[2])
	d.Set("stage_name", split[3])

	return nil
}

func resourceAliyunApigatewayAppAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	reqRemoveAuth := cloudapi.CreateRemoveAppsAuthoritiesRequest()
	split := strings.Split(d.Id(), COLON_SEPARATED)
	reqRemoveAuth.GroupId = split[0]
	reqRemoveAuth.ApiId = split[1]
	reqRemoveAuth.AppIds = split[2]
	reqRemoveAuth.StageName = split[3]

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.RemoveAppsAuthorities(reqRemoveAuth)
		})
		if err != nil {
			if IsExceptedError(err, NotFoundAuthorization) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting authorization failed: %#v", err))
		}

		if _, err := cloudApiService.DescribeAuthorization(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error describing authorization failed when deleting: %#v", err))
		}
		return nil
	})
}
