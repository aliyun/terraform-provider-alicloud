package cen

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudCenInstanceGrant() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenInstanceGrantCreate,
		Read:   resourceAlicloudCenInstanceGrantRead,
		Delete: resourceAlicloudCenInstanceGrantDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"child_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cen_owner_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenInstanceGrantCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenId := d.Get("cen_id").(string)
	ownerId := d.Get("cen_owner_id").(string)
	regionId := client.RegionId
	instanceId := d.Get("child_instance_id").(string)
	instanceType, err := GetCenChildInstanceType(instanceId)
	if err != nil {
		return errs.WrapError(err)
	}

	request := vpc.CreateGrantInstanceToCenRequest()
	request.RegionId = regionId
	request.CenId = cenId
	request.CenOwnerId = requests.Integer(ownerId)
	request.InstanceId = instanceId
	request.InstanceType = instanceType
	request.ClientToken = helper.BuildClientToken(request.GetActionName())

	var raw interface{}
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.GrantInstanceToCen(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "UnknownError", "TaskConflict"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_instance_grant", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(cenId + COLON_SEPARATED + instanceId + COLON_SEPARATED + string(ownerId))

	return resourceAlicloudCenInstanceGrantRead(d, meta)
}

func resourceAlicloudCenInstanceGrantRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	instanceId := parts[1]

	object, err := cbnService.DescribeCenInstanceGrant(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("cen_id", object.CenInstanceId)
	d.Set("cen_owner_id", fmt.Sprintf("%d", object.CenOwnerId))
	d.Set("child_instance_id", instanceId)

	return nil
}

func resourceAlicloudCenInstanceGrantDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	cenId := parts[0]
	instanceId := parts[1]
	ownerId := parts[2]
	instanceType, err := GetCenChildInstanceType(instanceId)
	if err != nil {
		return errs.WrapError(err)
	}

	request := vpc.CreateRevokeInstanceFromCenRequest()
	request.RegionId = client.RegionId
	request.InstanceId = instanceId
	request.InstanceType = instanceType
	request.CenId = cenId
	request.CenOwnerId = requests.Integer(ownerId)

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.RevokeInstanceFromCen(request)
		})

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
				return nil
			} else if errs.IsExpectedErrors(err, []string{"IncorrectStatus", "TaskConflict"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return errs.WrapError(cbnService.WaitForCenInstanceGrant(d.Id(), Deleted, DefaultCenTimeout))
}
