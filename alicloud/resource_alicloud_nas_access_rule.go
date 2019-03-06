package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/resource"
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
			"access_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_cidr_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"rw_access_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"RDWR", "RDONLY"}),
			},
			"user_access_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"no_squash", "root_squash", "all_squash"}),
			},
			"priority": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 100),
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
	if d.Get("rw_access_type").(string) != "" {
		request.RWAccessType = d.Get("rw_access_type").(string)
	}
	if d.Get("user_access_type").(string) != "" {
		request.UserAccessType = d.Get("user_access_type").(string)
	}
	request.Priority = requests.NewInteger(d.Get("priority").(int))
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateAccessRule(request)
	})
	ar, _ := raw.(*nas.CreateAccessRuleResponse)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "nas_access_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%s%s%s", d.Get("access_group_name").(string), COLON_SEPARATED, ar.AccessRuleId))
	return resourceAlicloudNasAccessRuleRead(d, meta)
}

func resourceAlicloudNasAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	split := strings.Split(d.Id(), COLON_SEPARATED)

	request := nas.CreateModifyAccessRuleRequest()
	request.AccessGroupName = d.Get("access_group_name").(string)
	request.AccessRuleId = split[1]
	request.SourceCidrIp = d.Get("source_cidr_ip").(string)
	request.RWAccessType = d.Get("rw_access_type").(string)
	request.UserAccessType = d.Get("user_access_type").(string)
	request.Priority = requests.NewInteger(d.Get("priority").(int))
	_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.ModifyAccessRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return resourceAlicloudNasAccessRuleRead(d, meta)
}

func resourceAlicloudNasAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)
	resp, err := nasService.DescribeNasAccessRule(split[0])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("source_cidr_ip", resp.SourceCidrIp)
	d.Set("access_group_name", split[0])
	d.Set("priority", resp.Priority)
	d.Set("rw_access_type", resp.RWAccess)
	d.Set("user_access_type", resp.UserAccess)

	return nil
}

func resourceAlicloudNasAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)
	request := nas.CreateDeleteAccessRuleRequest()

	request.AccessRuleId = split[1]
	request.AccessGroupName = split[0]
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DeleteAccessRule(request)
		})

		if err != nil {
			if IsExceptedError(err, ForbiddenNasNotFound) {
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		if _, err := nasService.DescribeNasAccessRule(split[0]); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))

	})
}
