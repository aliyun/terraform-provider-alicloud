package alicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	dms_enterprise "github.com/aliyun/alibaba-cloud-sdk-go/services/dms-enterprise"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDmsEnterpriseUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDmsEnterpriseUserCreate,
		Read:   resourceAlicloudDmsEnterpriseUserRead,
		Update: resourceAlicloudDmsEnterpriseUserUpdate,
		Delete: resourceAlicloudDmsEnterpriseUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"max_execute_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_result_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"mobile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLE", "NORMAL"}, false),
				Default:      "NORMAL",
			},
			"tid": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"nick_name"},
			},
			"nick_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'nick_name' has been deprecated from version 1.100.0. Use 'user_name' instead.",
				ConflictsWith: []string{"user_name"},
			},
		},
	}
}

func resourceAlicloudDmsEnterpriseUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := dms_enterprise.CreateRegisterUserRequest()
	if v, ok := d.GetOk("mobile"); ok {
		request.Mobile = v.(string)
	}

	if v, ok := d.GetOk("role_names"); ok {
		request.RoleNames = convertListToCommaSeparate(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("tid"); ok {
		request.Tid = requests.NewInteger(v.(int))
	}

	if v, err := strconv.Atoi(d.Get("uid").(string)); err == nil {
		request.Uid = requests.NewInteger(v)
	} else {
		return WrapError(err)
	}
	if v, ok := d.GetOk("user_name"); ok {
		request.UserNick = v.(string)
	} else if v, ok := d.GetOk("nick_name"); ok {
		request.UserNick = v.(string)
	} else {
		return WrapError(Error(`[ERROR] Argument "nick_name" or "user_name" must be set one!`))
	}

	raw, err := client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
		return dms_enterpriseClient.RegisterUser(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dms_enterprise_user", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	d.SetId(fmt.Sprintf("%v", request.Uid))

	return resourceAlicloudDmsEnterpriseUserUpdate(d, meta)
}
func resourceAlicloudDmsEnterpriseUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dms_enterpriseService := Dms_enterpriseService{client}
	object, err := dms_enterpriseService.DescribeDmsEnterpriseUser(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_enterprise_user dms_enterpriseService.DescribeDmsEnterpriseUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("uid", d.Id())
	d.Set("mobile", object.Mobile)
	d.Set("role_names", object.RoleNameList.RoleNames)
	d.Set("status", object.State)
	d.Set("user_name", object.NickName)
	d.Set("nick_name", object.NickName)
	return nil
}
func resourceAlicloudDmsEnterpriseUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dms_enterpriseService := Dms_enterpriseService{client}
	d.Partial(true)

	update := false
	request := dms_enterprise.CreateUpdateUserRequest()
	if v, err := strconv.Atoi(d.Id()); err == nil {
		request.Uid = requests.NewInteger(v)
	} else {
		return WrapError(err)
	}
	if !d.IsNewResource() && d.HasChange("mobile") {
		update = true
		request.Mobile = d.Get("mobile").(string)
	}
	if !d.IsNewResource() && d.HasChange("role_names") {
		update = true
		request.RoleNames = convertListToCommaSeparate(d.Get("role_names").(*schema.Set).List())
	}
	if !d.IsNewResource() && d.HasChange("user_name") {
		update = true
		request.UserNick = d.Get("user_name").(string)
	}
	if !d.IsNewResource() && d.HasChange("nick_name") {
		update = true
		request.UserNick = d.Get("nick_name").(string)
	}
	if update {
		if _, ok := d.GetOk("max_execute_count"); ok {
			request.MaxExecuteCount = requests.NewInteger(d.Get("max_execute_count").(int))
		}
		if _, ok := d.GetOk("max_result_count"); ok {
			request.MaxResultCount = requests.NewInteger(d.Get("max_result_count").(int))
		}
		if _, ok := d.GetOk("tid"); ok {
			request.Tid = requests.NewInteger(d.Get("tid").(int))
		}
		raw, err := client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
			return dms_enterpriseClient.UpdateUser(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("mobile")
		d.SetPartial("role_names")
		d.SetPartial("nick_name")
		d.SetPartial("user_name")
	}
	if d.HasChange("status") {
		object, err := dms_enterpriseService.DescribeDmsEnterpriseUser(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object.State != target {
			if target == "DISABLE" {
				request := dms_enterprise.CreateDisableUserRequest()
				if v, err := strconv.Atoi(d.Id()); err == nil {
					request.Uid = requests.NewInteger(v)
				} else {
					return WrapError(err)
				}
				if v, ok := d.GetOk("tid"); ok {
					request.Tid = requests.NewInteger(v.(int))
				}
				raw, err := client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
					return dms_enterpriseClient.DisableUser(request)
				})
				addDebug(request.GetActionName(), raw)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
			if target == "NORMAL" {
				request := dms_enterprise.CreateEnableUserRequest()
				if v, err := strconv.Atoi(d.Id()); err == nil {
					request.Uid = requests.NewInteger(v)
				} else {
					return WrapError(err)
				}
				if v, ok := d.GetOk("tid"); ok {
					request.Tid = requests.NewInteger(v.(int))
				}
				raw, err := client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
					return dms_enterpriseClient.EnableUser(request)
				})
				addDebug(request.GetActionName(), raw)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudDmsEnterpriseUserRead(d, meta)
}
func resourceAlicloudDmsEnterpriseUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := dms_enterprise.CreateDeleteUserRequest()
	if v, err := strconv.Atoi(d.Id()); err == nil {
		request.Uid = requests.NewInteger(v)
	} else {
		return WrapError(err)
	}
	if v, ok := d.GetOk("tid"); ok {
		request.Tid = requests.NewInteger(v.(int))
	}
	raw, err := client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
		return dms_enterpriseClient.DeleteUser(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
