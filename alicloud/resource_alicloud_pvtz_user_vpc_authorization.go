package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudPvtzUserVpcAuthorization() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzUserVpcAuthorizationCreate,
		Read:   resourceAlicloudPvtzUserVpcAuthorizationRead,
		Update: resourceAlicloudPvtzUserVpcAuthorizationUpdate,
		Delete: resourceAlicloudPvtzUserVpcAuthorizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auth_channel": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RESOURCE_DIRECTORY"}, false),
			},
			"auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"NORMAL", "CLOUD_PRODUCT"}, false),
			},
			"authorized_user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPvtzUserVpcAuthorizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddUserVpcAuthorization"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("auth_channel"); ok {
		request["AuthChannel"] = v
	}
	if v, ok := d.GetOk("auth_type"); ok {
		request["AuthType"] = v
	}
	request["AuthorizedUserId"] = d.Get("authorized_user_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"System.Busy"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_user_vpc_authorization", action, AlibabaCloudSdkGoERROR)
	}

	// Default auth_type to NORMAL when omitted, matching the API default for
	// DeleteUserVpcAuthorization. Reading request["AuthType"] would yield nil
	// and produce a "<nil>" literal in the resource ID when the optional field
	// is unset, breaking subsequent Read and Delete operations.
	authType := d.Get("auth_type").(string)
	if authType == "" {
		authType = "NORMAL"
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("authorized_user_id"), authType))

	return resourceAlicloudPvtzUserVpcAuthorizationRead(d, meta)
}
func resourceAlicloudPvtzUserVpcAuthorizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	// Migrate legacy state IDs produced by a previous bug where unset auth_type
	// resulted in ID = "<authorized_user_id>:<nil>". Treat "<nil>" as "NORMAL"
	// so Read/Delete send the correct AuthType to the API.
	if parts[1] == "<nil>" {
		parts[1] = "NORMAL"
		d.SetId(fmt.Sprintf("%s:%s", parts[0], parts[1]))
	}
	_, err = pvtzService.DescribePvtzUserVpcAuthorization(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pvtz_user_vpc_authorization pvtzService.DescribePvtzUserVpcAuthorization Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("auth_type", parts[1])
	d.Set("authorized_user_id", parts[0])
	return nil
}
func resourceAlicloudPvtzUserVpcAuthorizationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudPvtzUserVpcAuthorizationRead(d, meta)
}
func resourceAlicloudPvtzUserVpcAuthorizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	// Migrate legacy state IDs where auth_type was "<nil>" (see Read for details).
	if parts[1] == "<nil>" {
		parts[1] = "NORMAL"
	}
	action := "DeleteUserVpcAuthorization"
	var response map[string]interface{}
	request := map[string]interface{}{
		"AuthType":         parts[1],
		"AuthorizedUserId": parts[0],
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"System.Busy"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ZoneVpc.Auth.NotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
