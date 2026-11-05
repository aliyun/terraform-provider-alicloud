package apigateway

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceVpcAccess() *schema.Resource {
	return &schema.Resource{
		Create: ResourceVpcAccessCreate,
		Read:   ResourceVpcAccessRead,
		Delete: ResourceVpcAccessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"vpc_target_host_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_access_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceVpcAccessCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cloudapi.CreateSetVpcAccessRequest()
	request.RegionId = client.RegionId
	request.Name = d.Get("name").(string)
	request.VpcId = d.Get("vpc_id").(string)
	request.InstanceId = d.Get("instance_id").(string)
	request.Port = requests.NewInteger(d.Get("port").(int))
	request.VpcTargetHostName = d.Get("vpc_target_host_name").(string)

	var raw interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		raw, err = client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.SetVpcAccess(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_api_gateway_vpc_access", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s:%s", request.Name, request.VpcId, request.InstanceId, request.Port))

	return ResourceVpcAccessRead(d, meta)
}

func ResourceVpcAccessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	object, err := cloudApiService.DescribeApiGatewayVpcAccess(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("name", object["Name"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("instance_id", object["InstanceId"])
	d.Set("port", object["Port"])
	d.Set("vpc_target_host_name", object["VpcTargetHostName"])
	if v, ok := object["VpcAccessId"]; ok {
		d.Set("vpc_access_id", v)
	}

	return nil
}

func ResourceVpcAccessDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cloudapi.CreateRemoveVpcAccessRequest()
	request.RegionId = client.RegionId
	request.VpcId = d.Get("vpc_id").(string)
	request.InstanceId = d.Get("instance_id").(string)
	request.Port = requests.NewInteger(d.Get("port").(int))

	var raw interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		raw, err = client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.RemoveVpcAccess(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return nil

}
