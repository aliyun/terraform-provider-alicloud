package alb

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceListenerAclAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerAclAttachmentCreate,
		Read:   resourceListenerAclAttachmentRead,
		Delete: resourceListenerAclAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"acl_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"White", "Black"}, false),
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceListenerAclAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	albService := NewService(client)

	action := "AssociateAclsWithListener"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["ListenerId"] = d.Get("listener_id")
	request["AclIds.1"] = d.Get("acl_id")

	request["ClientToken"] = helper.BuildClientToken(action)

	request["AclType"] = d.Get("acl_type")
	// The acl can still be Creating or Configuring when this resource is created right after
	// alicloud_alb_acl or alicloud_alb_acl_entry_attachment; AssociateAclsWithListener rejects
	// such an acl with IncorrectStatus.Acl, so wait for it to become Available first.
	stateConf := helper.BuildStateConf([]string{"Creating", "Configuring"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AclStateRefreshFunc(fmt.Sprint(d.Get("acl_id")), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, fmt.Sprint(d.Get("acl_id")))
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ResourceInConfiguring.Listener", "IncorrectStatus.Listener", "IncorrectStatus.Acl", "Conflict.Acl", "LockFailed", "AclStatusNotStable"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alb_listener_acl_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ListenerId"], request["AclIds.1"]))

	stateConf = helper.BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.ListenerAclAttachmentStateRefreshFunc(d.Id(), "$.AclConfig.AclRelations[0].Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceListenerAclAttachmentRead(d, meta)
}

func resourceListenerAclAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := NewService(client)

	objectRaw, err := albService.DescribeListenerAclAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_listener_acl_attachment DescribeListenerAclAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("listener_id", objectRaw["ListenerId"])
	aclConfig1RawObj, _ := jsonpath.Get("$.AclConfig", objectRaw)
	aclConfig1Raw := make(map[string]interface{})
	if aclConfig1RawObj != nil {
		aclConfig1Raw = aclConfig1RawObj.(map[string]interface{})
	}
	d.Set("acl_type", aclConfig1Raw["AclType"])
	aclRelations1RawObj, _ := jsonpath.Get("$.AclConfig.AclRelations[*]", objectRaw)
	aclRelations1Raw := make([]interface{}, 0)
	if aclRelations1RawObj != nil {
		aclRelations1Raw = aclRelations1RawObj.([]interface{})
	}
	aclRelationsChild1Raw := aclRelations1Raw[0].(map[string]interface{})
	parts := strings.Split(d.Id(), ":")
	for _, vv := range aclRelations1Raw {
		if vv.(map[string]interface{})["AclId"] == parts[1] {
			aclRelationsChild1Raw = vv.(map[string]interface{})
			break
		}
	}
	d.Set("status", aclRelationsChild1Raw["Status"])
	d.Set("acl_id", aclRelationsChild1Raw["AclId"])

	return nil
}

func resourceListenerAclAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	albService := NewService(client)
	parts := strings.Split(d.Id(), ":")
	action := "DissociateAclsFromListener"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["ListenerId"] = parts[0]
	request["AclIds.1"] = parts[1]

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"LockFailed", "ResourceInConfiguring.Listener", "IncorrectStatus.Listener", "IncorrectStatus.Acl", "AclStatusNotStable"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albService.ListenerAclAttachmentStateRefreshFunc(d.Id(), "$.AclConfig.AclRelations[0].Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
