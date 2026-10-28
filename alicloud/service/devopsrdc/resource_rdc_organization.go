package devopsrdc

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRdcOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceRdcOrganizationCreate,
		Read:   resourceRdcOrganizationRead,
		Update: resourceRdcOrganizationUpdate,
		Delete: resourceRdcOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.238.0 and will be removed in the future.",
		Schema: map[string]*schema.Schema{
			"desired_member_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"organization_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"real_pk": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRdcOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDevopsOrganization"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("desired_member_count"); ok {
		request["DesiredMemberCount"] = v
	}
	request["OrgName"] = d.Get("organization_name")
	if v, ok := d.GetOk("real_pk"); ok {
		request["RealPk"] = v
	}
	request["Source"] = d.Get("source")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("devops-rdc", "2020-03-03", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_rdc_organization", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Object"]))

	return resourceRdcOrganizationRead(d, meta)
}

func resourceRdcOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)
	object, err := service.DescribeRdcOrganization(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rdc_organization service.DescribeRdcOrganization Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("organization_name", object["Name"])
	return nil
}

func resourceRdcOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceRdcOrganizationRead(d, meta)
}

func resourceRdcOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDevopsOrganization"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"OrgId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("devops-rdc", "2020-03-03", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidOrganization.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
