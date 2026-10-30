package smartag

import (
	"fmt"

	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudSagSnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagSnatCreate,
		Read:   resourceAlicloudSagSnatRead,
		Delete: resourceAlicloudSagSnatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sag_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.ValidateCIDRNetworkAddress,
			},
			"snat_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
		},
	}
}

func resourceAlicloudSagSnatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartagsdk.CreateAddSnatEntryRequest()

	request.SmartAGId = d.Get("sag_id").(string)
	request.CidrBlock = d.Get("cidr_block").(string)
	request.SnatIp = d.Get("snat_ip").(string)

	raw, err := client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
		return sagClient.AddSnatEntry(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_sag_snat_entry", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartagsdk.AddSnatEntryResponse)
	d.SetId(fmt.Sprintf("%s%s%s", request.SmartAGId, ":", response.InstanceId))

	return resourceAlicloudSagSnatRead(d, meta)
}

func resourceAlicloudSagSnatRead(d *schema.ResourceData, meta interface{}) error {
	sagService := Service{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagSnatEntry(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, _ := helper.ParseResourceId(d.Id(), 2)
	d.Set("sag_id", parts[0])
	d.Set("cidr_block", object.CidrBlock)
	d.Set("snat_ip", object.SnatIp)

	return nil
}

func resourceAlicloudSagSnatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := Service{client}
	request := smartagsdk.CreateDeleteSnatEntryRequest()
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request.SmartAGId = parts[0]
	request.InstanceId = parts[1]

	raw, err := client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
		return sagClient.DeleteSnatEntry(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ParameterSagSnatId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errs.WrapError(sagService.WaitForSagSnatEntry(d.Id(), Deleted, DefaultTimeoutMedium))
}
