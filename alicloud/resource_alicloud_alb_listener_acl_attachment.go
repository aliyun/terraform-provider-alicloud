package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAlbListenerAclAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlbListenerAclAttachmentCreate,
		Read:   resourceAliCloudAlbListenerAclAttachmentRead,
		Delete: resourceAliCloudAlbListenerAclAttachmentDelete,
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
				ValidateFunc: StringInSlice([]string{"White", "Black"}, false),
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

func resourceAliCloudAlbListenerAclAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AssociateAclsWithListener"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["ListenerId"] = d.Get("listener_id")
	request["AclIds.1"] = d.Get("acl_id")

	request["ClientToken"] = buildClientToken(action)

	request["AclType"] = d.Get("acl_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceInConfiguring.Listener", "IncorrectStatus.Listener", "Conflict.Acl"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_listener_acl_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ListenerId"], request["AclIds.1"]))

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albServiceV2.AlbListenerAclAttachmentStateRefreshFunc(d.Id(), "$.AclConfig.AclRelations[0].Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAlbListenerAclAttachmentRead(d, meta)
}

func resourceAliCloudAlbListenerAclAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}

	objectRaw, err := albServiceV2.DescribeAlbListenerAclAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_listener_acl_attachment DescribeAlbListenerAclAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAliCloudAlbListenerAclAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DissociateAclsFromListener"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["ListenerId"] = parts[0]
	request["AclIds.1"] = parts[1]

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"LockFailed", "ResourceInConfiguring.Listener", "IncorrectStatus.Listener", "IncorrectStatus.Acl"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albServiceV2.AlbListenerAclAttachmentStateRefreshFunc(d.Id(), "$.AclConfig.AclRelations[0].Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
