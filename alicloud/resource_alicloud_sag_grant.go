package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudSagGrant() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagGrantCreate,
		Read:   resourceAlicloudSagGrantRead,
		Delete: resourceAlicloudSagGrantDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sag_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ccn_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ccn_uid": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudSagGrantCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateGrantSagInstanceToCcnRequest()

	request.RegionId = client.RegionId
	request.SmartAGId = d.Get("sag_id").(string)
	request.CcnInstanceId = d.Get("ccn_id").(string)
	request.CcnUid = requests.NewInteger(d.Get("ccn_uid").(int))
	var err error
	var raw interface{}
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.GrantSagInstanceToCcn(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{OperationBlocking, UnknownError}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sag_grant", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(fmt.Sprintf("%s%s%s", request.SmartAGId, COLON_SEPARATED, request.CcnInstanceId))

	return resourceAlicloudSagGrantRead(d, meta)
}

func resourceAlicloudSagGrantRead(d *schema.ResourceData, meta interface{}) error {
	sagService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagGrant(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("sag_id", object.SmartAGId)
	d.Set("ccn_id", object.CcnInstanceId)
	d.Set("ccn_uid", object.CcnUid)

	return nil
}

func resourceAlicloudSagGrantDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	sagId := parts[0]
	ccnId := parts[1]
	request := smartag.CreateRevokeSagInstanceFromCcnRequest()
	request.RegionId = client.RegionId
	request.SmartAGId = sagId
	request.CcnInstanceId = ccnId

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.RevokeSagInstanceFromCcn(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{IncorrectStatus, TaskConflict}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if IsExceptedError(err, InvalidInstanceIdNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(sagService.WaitForSagGrant(d.Id(), Deleted, DefaultTimeoutMedium))
}
