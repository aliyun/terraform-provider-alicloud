package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudNasAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasAccessRuleCreate,
		Read:   resourceAlicloudNasAccessRuleRead,
		Update: resourceAlicloudNasAccessRuleUpdate,
		Delete: resourceAlicloudNasAccessRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_cidr_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rw_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"RDWR", "RDONLY"}),
			},
			"user_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"no_squash", "root_squash", "all_squash"}),
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateIntegerInRange(1, 100),
			},
			"access_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudNasAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateCreateAccessRuleRequest()
	request.RegionId = string(client.Region)
	request.AccessGroupName = d.Get("access_group_name").(string)
	request.SourceCidrIp = d.Get("source_cidr_ip").(string)
	if v, ok := d.GetOk("rw_access_type"); ok && v.(string) != "" {
		request.RWAccessType = v.(string)
	}
	if v, ok := d.GetOk("user_access_type"); ok && v.(string) != "" {
		request.UserAccessType = v.(string)
	}
	request.Priority = requests.NewInteger(d.Get("priority").(int))
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateAccessRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_access_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*nas.CreateAccessRuleResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(fmt.Sprintf("%s%s%s", d.Get("access_group_name").(string), COLON_SEPARATED, response.AccessRuleId))
	return resourceAlicloudNasAccessRuleRead(d, meta)
}

func resourceAlicloudNasAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		err = WrapError(err)
		return err
	}
	request := nas.CreateModifyAccessRuleRequest()
	request.RegionId = client.RegionId
	request.AccessGroupName = d.Get("access_group_name").(string)
	request.AccessRuleId = parts[1]
	request.SourceCidrIp = d.Get("source_cidr_ip").(string)
	request.RWAccessType = d.Get("rw_access_type").(string)
	request.UserAccessType = d.Get("user_access_type").(string)
	request.Priority = requests.NewInteger(d.Get("priority").(int))
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.ModifyAccessRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAlicloudNasAccessRuleRead(d, meta)
}

func resourceAlicloudNasAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasAccessRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("access_rule_id", object.AccessRuleId)
	d.Set("source_cidr_ip", object.SourceCidrIp)
	d.Set("access_group_name", parts[0])
	d.Set("priority", object.Priority)
	d.Set("rw_access_type", object.RWAccess)
	d.Set("user_access_type", object.UserAccess)

	return nil
}

func resourceAlicloudNasAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	request := nas.CreateDeleteAccessRuleRequest()
	request.RegionId = client.RegionId
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		err = WrapError(err)
		return err
	}
	request.AccessRuleId = parts[1]
	request.AccessGroupName = parts[0]

	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.DeleteAccessRule(request)
	})

	if err != nil {
		if IsExceptedError(err, ForbiddenNasNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(nasService.WaitForNasAccessRule(d.Id(), Deleted, DefaultTimeoutMedium))

}
