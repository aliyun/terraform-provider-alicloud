package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"time"
)

func resourceAlicloudNAS_AccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNAS_AccessRuleCreate,
		Read:   resourceAlicloudNAS_AccessRuleRead,
		Update: resourceAlicloudNAS_AccessRuleUpdate,
		Delete: resourceAlicloudNAS_AccessRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"accessgroup_name": &schema.Schema{
				Type:         schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sourcecidr_ip": &schema.Schema{
				Type:         schema.TypeString,
				Required: true,
			},
			"rwaccess_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"RDWR",
					"RDONLY",
				}),
			},
			"useraccess_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"no_squash",
					"root_squash",
					"all_squash",
				}),
			},
			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:          60,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
		},
	}
}

func resourceAlicloudNAS_AccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateCreateAccessRuleRequest()
	request.RegionId = string(client.Region)
	request.AccessGroupName = d.Get("accessgroup_name").(string)
	request.SourceCidrIp = d.Get("sourcecidr_ip").(string)
	if d.Get("rwaccess_type").(string) != "" {
		request.RWAccessType = d.Get("rwaccess_type").(string)
	}
	if  d.Get("useraccess_type").(string) != "" {
		request.UserAccessType = d.Get("useraccess_type").(string)
	}
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateAccessRule(request)
	})
	ar, _ := raw.(*nas.CreateAccessRuleResponse)
	if err != nil {
		return fmt.Errorf("Error Waitting for NAS available: %#v", err)
	}
	d.SetId(d.Get("accessgroup_name").(string) + ":" + ar.AccessRuleId)
	return resourceAlicloudNAS_AccessRuleUpdate(d, meta)
}

func resourceAlicloudNAS_AccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	split := strings.Split(d.Id(), ":")
	attributeUpdate := false
	request := nas.CreateModifyAccessRuleRequest()
	request.AccessGroupName = split[0]
	request.AccessRuleId = split[1]
	if d.HasChange("sourcecidr_ip") {
		attributeUpdate = true
		d.SetPartial("sourcecidr_ip")
		request.SourceCidrIp = d.Get("sourcecidr_ip").(string)
	}
	if d.HasChange("rwaccess_type") {
		attributeUpdate = true
		d.SetPartial("rwaccess_type")
		request.RWAccessType = d.Get("rwaccess_type").(string)
	}
	if d.HasChange("useraccess_type") {
		attributeUpdate = true
		d.SetPartial("useraccess_type")
		request.UserAccessType = d.Get("useraccess_type").(string)
	}
	if d.HasChange("priority") {
		attributeUpdate = true
		d.SetPartial("priority")
		request.Priority = requests.Integer(d.Get("priority").(int))
	}
	if attributeUpdate {
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.ModifyAccessRule(request)
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAlicloudNAS_AccessRuleRead(d, meta)
}

func resourceAlicloudNAS_AccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	split := strings.Split(d.Id(), ":")
	resp, err := nasService.DescribeAccessRules(split[0])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("sourcecidr_ip", resp.SourceCidrIp)
	d.Set("accessrule_id", resp.AccessRuleId)
	d.Set("priority", resp.Priority)
	if resp.RWAccess != "" {
		d.Set("rwaccess_type", resp.RWAccess)
	}
	if resp.UserAccess != "" {
		d.Set("useraccess_type", resp.UserAccess)
	}


	return nil
}

func resourceAlicloudNAS_AccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	split := strings.Split(d.Id(), ":")
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
			return resource.RetryableError(fmt.Errorf("Delete NAS timeout and got an error: %#v.", err))
		}

		if _, err := nasService.DescribeAccessRules(split[0]); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

