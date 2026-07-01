package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudResourceManagerHandshakeAcceptance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerHandshakeAcceptanceCreate,
		Read:   resourceAliCloudResourceManagerHandshakeAcceptanceRead,
		Delete: resourceAliCloudResourceManagerHandshakeAcceptanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"handshake_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"note": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_directory_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_entity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudResourceManagerHandshakeAcceptanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AcceptHandshake"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["HandshakeId"] = d.Get("handshake_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_handshake_acceptance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(d.Get("handshake_id")))

	return resourceAliCloudResourceManagerHandshakeAcceptanceRead(d, meta)
}

func resourceAliCloudResourceManagerHandshakeAcceptanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerHandshake(d.Id())
	if err != nil {
		// Once accepted, the invited account can no longer read the handshake: GetHandshake returns
		// HandshakeStatusMismatch. Treat that as a successfully accepted invitation.
		if IsExpectedErrors(err, []string{"HandshakeStatusMismatch"}) {
			d.Set("handshake_id", d.Id())
			d.Set("status", "Accepted")
			return nil
		}
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_handshake_acceptance DescribeResourceManagerHandshake Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("handshake_id", d.Id())
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("expire_time", objectRaw["ExpireTime"])
	d.Set("master_account_id", objectRaw["MasterAccountId"])
	d.Set("master_account_name", objectRaw["MasterAccountName"])
	d.Set("modify_time", objectRaw["ModifyTime"])
	d.Set("note", objectRaw["Note"])
	d.Set("resource_directory_id", objectRaw["ResourceDirectoryId"])
	d.Set("status", objectRaw["Status"])
	d.Set("target_entity", objectRaw["TargetEntity"])
	d.Set("target_type", objectRaw["TargetType"])

	return nil
}

func resourceAliCloudResourceManagerHandshakeAcceptanceDelete(d *schema.ResourceData, meta interface{}) error {
	// Accepting a handshake is irreversible: once the invited account joins the resource directory,
	// there is no API to revoke the acceptance. Removing this resource only drops it from state.
	log.Printf("[WARN] Cannot destroy alicloud_resource_manager_handshake_acceptance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
