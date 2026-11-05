package smartag

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudCloudConnectNetworkGrant() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudConnectNetworkGrantCreate,
		Read:   resourceAlicloudCloudConnectNetworkGrantRead,
		Delete: resourceAlicloudCloudConnectNetworkGrantDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cen_uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCloudConnectNetworkGrantCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartagsdk.CreateGrantInstanceToCbnRequest()

	request.RegionId = client.RegionId
	request.CcnInstanceId = d.Get("ccn_id").(string)
	request.CenInstanceId = d.Get("cen_id").(string)
	request.CenUid = requests.Integer(d.Get("cen_uid").(string))

	var err error
	var raw interface{}
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.GrantInstanceToCbn(request)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_connect_network_grant", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(fmt.Sprintf("%s%s%s", request.CcnInstanceId, ":", request.CenInstanceId))

	return resourceAlicloudCloudConnectNetworkGrantRead(d, meta)
}

func resourceAlicloudCloudConnectNetworkGrantRead(d *schema.ResourceData, meta interface{}) error {
	sagService := Service{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeCloudConnectNetworkGrant(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("ccn_id", object.CcnInstanceId)
	d.Set("cen_id", object.CenInstanceId)
	d.Set("cen_uid", fmt.Sprint(object.CenUid))

	return nil
}

func resourceAlicloudCloudConnectNetworkGrantDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := Service{client}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := smartagsdk.CreateRevokeInstanceFromCbnRequest()
	request.RegionId = client.RegionId
	request.CcnInstanceId = parts[0]
	request.CenInstanceId = parts[1]

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.RevokeInstanceFromCbn(request)
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
	return errs.WrapError(sagService.WaitForCloudConnectNetworkGrant(d.Id(), Deleted, DefaultTimeoutMedium))
}
