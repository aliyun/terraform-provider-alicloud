package ecd

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAdConnectorDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAdConnectorDirectoryCreate,
		Read:   resourceAdConnectorDirectoryRead,
		Update: resourceAdConnectorDirectoryUpdate,
		Delete: resourceAdConnectorDirectoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			// Create waits for the directory to leave REGISTERING and settle on
			// REGISTERED, ERROR, or NEEDCONFIGCONNECTOR. The AD connector
			// directory's registration state machine retries the unreachable AD
			// longer than the sibling office site's fast-failing trust config
			// (office_site uses 5m), so the create wait needs a longer window to
			// absorb a slow transition out of REGISTERING before timing out.
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"directory_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_:-]{1,254}$`), "It must be `2` to `255` characters in length, The name must start with a letter, and can contain letters, digits, colons (:), underscores (_), and hyphens (-).")),
			},
			"desktop_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC", "INTERNET", "ANY"}, false),
			},
			"dns_address": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"domain_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"domain_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"domain_user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enable_admin_access": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"mfa_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"specification": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_domain_dns_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"sub_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
		},
	}
}

func resourceAdConnectorDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateADConnectorDirectory"
	request := make(map[string]interface{})
	var err error
	request["DirectoryName"] = d.Get("directory_name")
	if v, ok := d.GetOk("desktop_access_type"); ok {
		request["DesktopAccessType"] = v
	}
	request["DnsAddress"] = d.Get("dns_address")
	request["DomainName"] = d.Get("domain_name")
	request["DomainPassword"] = d.Get("domain_password")
	request["DomainUserName"] = d.Get("domain_user_name")
	if v, ok := d.GetOkExists("enable_admin_access"); ok {
		request["EnableAdminAccess"] = v
	}
	if v, ok := d.GetOkExists("mfa_enabled"); ok {
		request["MfaEnabled"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("specification"); ok {
		request["Specification"] = v
	}
	if v, ok := d.GetOk("sub_domain_dns_address"); ok {
		request["SubDomainDnsAddress"] = v
	}
	if v, ok := d.GetOk("sub_domain_name"); ok {
		request["SubDomainName"] = v
	}
	request["VSwitchId"] = d.Get("vswitch_ids")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecd_ad_connector_directory", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DirectoryId"]))

	// Wait for the directory to settle before reading, so the dependent
	// VSwitch/VPC are not touched while the directory is still REGISTERING.
	// Previously the resource returned immediately after SetId, which left the
	// directory mid-registration and caused OfficeSiteNameUsed on the next run.
	// The terminal set mirrors alicloud_ecd_ad_connector_office_site's
	// EcdAdConnectorOfficeSiteStateRefreshFunc: REGISTERED is the success, ERROR a
	// failure, NEEDCONFIGCONNECTOR a connector-config-needed state the operator
	// must resolve out-of-band. All three are stable, so the resource only
	// returns once the directory is no longer REGISTERING.
	ecdService := NewService(client)
	stateConf := helper.BuildStateConf([]string{"REGISTERING"}, []string{"REGISTERED", "ERROR", "NEEDCONFIGCONNECTOR"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecdService.EcdAdConnectorDirectoryStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAdConnectorDirectoryRead(d, meta)
}

func resourceAdConnectorDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	object, err := ecdService.DescribeEcdAdConnectorDirectory(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_ad_connector_directory ecdService.DescribeEcdAdConnectorDirectory Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("directory_name", object["Name"])
	d.Set("desktop_access_type", object["DesktopAccessType"])
	d.Set("dns_address", object["DnsAddress"])
	d.Set("domain_name", object["DomainName"])
	d.Set("domain_user_name", object["DomainUserName"])
	d.Set("enable_admin_access", object["EnableAdminAccess"])
	d.Set("mfa_enabled", object["MfaEnabled"])
	d.Set("status", object["Status"])
	d.Set("sub_domain_dns_address", object["SubDnsAddress"])
	d.Set("sub_domain_name", object["SubDomainName"])
	d.Set("vswitch_ids", object["VSwitchIds"])
	return nil
}

func resourceAdConnectorDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAdConnectorDirectoryRead(d, meta)
}

func resourceAdConnectorDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDirectories"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DirectoryId": []string{d.Id()},
	}
	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"InvalidDirectoryStatus"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	// Wait for the directory to disappear before returning, so dependent
	// VSwitch/VPC are only destroyed once the directory is gone. Previously the
	// resource returned immediately after the delete RPC, causing the VSwitch to
	// be destroyed while the directory still referenced it (DependencyViolation).
	ecdService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, ecdService.EcdAdConnectorDirectoryStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
