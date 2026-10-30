package smartag

import (
	"fmt"
	"time"

	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudCloudConnectNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudConnectNetworkAttachmentCreate,
		Read:   resourceAlicloudCloudConnectNetworkAttachmentRead,
		Delete: resourceAlicloudCloudConnectNetworkAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sag_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCloudConnectNetworkAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartagsdk.CreateBindSmartAccessGatewayRequest()

	request.RegionId = client.RegionId
	request.CcnId = d.Get("ccn_id").(string)
	request.SmartAGId = d.Get("sag_id").(string)
	var err error
	var raw interface{}
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.BindSmartAccessGateway(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "UnknownError"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_connect_network_attachment", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(fmt.Sprintf("%s%s%s", request.CcnId, ":", request.SmartAGId))

	return resourceAlicloudCloudConnectNetworkAttachmentRead(d, meta)
}

func resourceAlicloudCloudConnectNetworkAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	sagService := Service{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeCloudConnectNetworkAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("ccn_id", object.AssociatedCcnId)
	d.Set("sag_id", object.SmartAGId)

	return nil
}

func resourceAlicloudCloudConnectNetworkAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := Service{client}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	ccnId := parts[0]
	sagId := parts[1]

	request := smartagsdk.CreateUnbindSmartAccessGatewayRequest()
	request.RegionId = client.RegionId
	request.CcnId = ccnId
	request.SmartAGId = sagId

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.UnbindSmartAccessGateway(request)
		})

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectStatus", "TaskConflict"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errs.WrapError(sagService.WaitForCloudConnectNetworkAttachment(d.Id(), Deleted, DefaultTimeoutMedium))
}
