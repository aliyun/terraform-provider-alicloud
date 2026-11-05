package smartag

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudCloudConnectNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudConnectNetworkCreate,
		Read:   resourceAlicloudCloudConnectNetworkRead,
		Update: resourceAlicloudCloudConnectNetworkUpdate,
		Delete: resourceAlicloudCloudConnectNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"is_default": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.ValidateVpnCIDRNetworkAddress,
			},
		},
	}
}

func resourceAlicloudCloudConnectNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartagsdk.CreateCreateCloudConnectNetworkRequest()

	request.IsDefault = requests.NewBoolean(d.Get("is_default").(bool))
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("cidr_block"); ok && v.(string) != "" {
		request.CidrBlock = v.(string)
	}

	var err error
	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.CreateCloudConnectNetwork(request)
		})
		if err != nil {
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			if errs.IsExpectedErrors(err, []string{"ResourceInOperating", "InternalError"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_connect_network", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*smartagsdk.CreateCloudConnectNetworkResponse)
	d.SetId(response.CcnId)

	return resourceAlicloudCloudConnectNetworkRead(d, meta)
}

func resourceAlicloudCloudConnectNetworkRead(d *schema.ResourceData, meta interface{}) error {
	ccnService := Service{meta.(*connectivity.AliyunClient)}
	object, err := ccnService.DescribeCloudConnectNetwork(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("cidr_block", object.CidrBlock)
	d.Set("is_default", object.IsDefault)

	return nil
}

func resourceAlicloudCloudConnectNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	request := smartagsdk.CreateModifyCloudConnectNetworkRequest()
	request.CcnId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
	}
	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if d.HasChange("cidr_block") {
		request.CidrBlock = d.Get("cidr_block").(string)
		update = true
	}
	if update {
		raw, err := client.WithSagClient(func(ccnClient *smartagsdk.Client) (interface{}, error) {
			return ccnClient.ModifyCloudConnectNetwork(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudCloudConnectNetworkRead(d, meta)
}

func resourceAlicloudCloudConnectNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := Service{client}
	request := smartagsdk.CreateDeleteCloudConnectNetworkRequest()
	request.CcnId = d.Id()

	raw, err := client.WithSagClient(func(ccnClient *smartagsdk.Client) (interface{}, error) {
		return ccnClient.DeleteCloudConnectNetwork(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ParameterCcnInstanceId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errs.WrapError(sagService.WaitForCloudConnectNetwork(d.Id(), Deleted, DefaultTimeoutMedium))
}
