package smartag

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudSagClientUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagClientUserCreate,
		Read:   resourceAlicloudSagClientUserRead,
		Update: resourceAlicloudSagClientUserUpdate,
		Delete: resourceAlicloudSagClientUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sag_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"user_mail": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"password": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: sagClientUserPasswordSuppressFunc,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
		},
	}
}

func resourceAlicloudSagClientUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartagsdk.CreateCreateSmartAccessGatewayClientUserRequest()

	request.SmartAGId = d.Get("sag_id").(string)
	request.UserMail = d.Get("user_mail").(string)
	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))

	if v, ok := d.GetOk("client_ip"); ok && v.(string) != "" {
		request.ClientIp = v.(string)
	}
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)

		password := d.Get("password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)
		if password == "" && kmsPassword == "" {
			return errs.WrapError(errs.Error("One of the 'password' and 'kms_encrypted_password' should be set when 'user_name' was set."))
		}
		if password != "" {
			request.Password = password
		} else {
			sagService := Service{client}
			decryptResp, err := sagService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return errs.WrapError(err)
			}
			request.Password = decryptResp
		}
	}

	raw, err := client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
		return sagClient.CreateSmartAccessGatewayClientUser(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_sag_client_user", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartagsdk.CreateSmartAccessGatewayClientUserResponse)
	d.SetId(fmt.Sprintf("%s%s%s", d.Get("sag_id").(string), ":", response.UserName))

	return resourceAlicloudSagClientUserRead(d, meta)
}

func resourceAlicloudSagClientUserRead(d *schema.ResourceData, meta interface{}) error {
	sagService := Service{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagClientUser(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, _ := helper.ParseResourceId(d.Id(), 2)

	d.Set("sag_id", parts[0])
	d.Set("bandwidth", object.Bandwidth)
	d.Set("user_mail", object.UserMail)
	d.Set("client_ip", object.ClientIp)
	d.Set("user_name", object.UserName)

	return nil
}

func resourceAlicloudSagClientUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	if d.HasChange("bandwidth") {
		request := smartagsdk.CreateModifySmartAccessGatewayClientUserRequest()
		request.SmartAGId = parts[0]
		request.UserName = parts[1]
		request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))
		raw, err := client.WithSagClient(func(ccnClient *smartagsdk.Client) (interface{}, error) {
			return ccnClient.ModifySmartAccessGatewayClientUser(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudSagClientUserRead(d, meta)
}

func resourceAlicloudSagClientUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := Service{client}
	request := smartagsdk.CreateDeleteSmartAccessGatewayClientUserRequest()
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request.SmartAGId = parts[0]
	request.UserName = parts[1]

	raw, err := client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
		return sagClient.DeleteSmartAccessGatewayClientUser(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ParameterSagClientId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errs.WrapError(sagService.WaitForSagClientUser(d.Id(), Deleted, DefaultTimeoutMedium))
}
