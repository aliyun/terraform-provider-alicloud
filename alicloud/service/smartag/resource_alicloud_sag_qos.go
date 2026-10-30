package smartag

import (
	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudSagQos() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagQosCreate,
		Read:   resourceAlicloudSagQosRead,
		Update: resourceAlicloudSagQosUpdate,
		Delete: resourceAlicloudSagQosDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
		},
	}
}

func resourceAlicloudSagQosCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartagsdk.CreateCreateQosRequest()

	request.QosName = d.Get("name").(string)

	raw, err := client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
		return sagClient.CreateQos(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_sag_qos", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartagsdk.CreateQosResponse)
	d.SetId(response.QosId)
	return resourceAlicloudSagQosRead(d, meta)
}

func resourceAlicloudSagQosRead(d *schema.ResourceData, meta interface{}) error {
	sagService := Service{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagQos(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("name", object.QosName)

	return nil
}

func resourceAlicloudSagQosUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("name") {
		request := smartagsdk.CreateModifyQosRequest()
		request.QosId = d.Id()
		request.QosName = d.Get("name").(string)
		raw, err := client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.ModifyQos(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudSagQosRead(d, meta)
}

func resourceAlicloudSagQosDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := Service{client}
	request := smartagsdk.CreateDeleteQosRequest()
	request.QosId = d.Id()

	raw, err := client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
		return sagClient.DeleteQos(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ParameterSagQosId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errs.WrapError(sagService.WaitForSagQos(d.Id(), Deleted, DefaultTimeoutMedium))
}
