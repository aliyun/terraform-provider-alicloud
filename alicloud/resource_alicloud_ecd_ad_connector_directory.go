package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcdAdConnectorDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdAdConnectorDirectoryCreate,
		Read:   resourceAlicloudEcdAdConnectorDirectoryRead,
		Update: resourceAlicloudEcdAdConnectorDirectoryUpdate,
		Delete: resourceAlicloudEcdAdConnectorDirectoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
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

func resourceAlicloudEcdAdConnectorDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateADConnectorDirectory"
	request := make(map[string]interface{})
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_ad_connector_directory", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DirectoryId"]))

	return resourceAlicloudEcdAdConnectorDirectoryRead(d, meta)
}
func resourceAlicloudEcdAdConnectorDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	object, err := ecdService.DescribeEcdAdConnectorDirectory(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_ad_connector_directory ecdService.DescribeEcdAdConnectorDirectory Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
func resourceAlicloudEcdAdConnectorDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudEcdAdConnectorDirectoryRead(d, meta)
}
func resourceAlicloudEcdAdConnectorDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDirectories"
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DirectoryId": []string{d.Id()},
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidDirectoryStatus"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
