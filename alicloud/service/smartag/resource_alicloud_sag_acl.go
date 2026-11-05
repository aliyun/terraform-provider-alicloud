package smartag

import (
	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudSagAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagAclCreate,
		Read:   resourceAlicloudSagAclRead,
		Update: resourceAlicloudSagAclUpdate,
		Delete: resourceAlicloudSagAclDelete,
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

func resourceAlicloudSagAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartagsdk.CreateCreateACLRequest()

	request.Name = d.Get("name").(string)

	raw, err := client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
		return sagClient.CreateACL(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_sag_acl", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartagsdk.CreateACLResponse)
	d.SetId(response.AclId)
	return resourceAlicloudSagAclRead(d, meta)
}

func resourceAlicloudSagAclRead(d *schema.ResourceData, meta interface{}) error {
	sagService := Service{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagAcl(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("name", object.Name)
	return nil
}

func resourceAlicloudSagAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("name") {
		request := smartagsdk.CreateModifyACLRequest()
		request.AclId = d.Id()
		request.Name = d.Get("name").(string)

		raw, err := client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.ModifyACL(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudSagAclRead(d, meta)
}

func resourceAlicloudSagAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := Service{client}
	request := smartagsdk.CreateDeleteACLRequest()
	request.AclId = d.Id()

	raw, err := client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
		return sagClient.DeleteACL(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ParameterSagACLId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errs.WrapError(sagService.WaitForSagAcl(d.Id(), Deleted, DefaultTimeoutMedium))
}
