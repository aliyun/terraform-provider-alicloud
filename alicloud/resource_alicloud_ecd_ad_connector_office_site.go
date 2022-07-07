package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEcdAdConnectorOfficeSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdAdConnectorOfficeSiteCreate,
		Read:   resourceAlicloudEcdAdConnectorOfficeSiteRead,
		Update: resourceAlicloudEcdAdConnectorOfficeSiteUpdate,
		Delete: resourceAlicloudEcdAdConnectorOfficeSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ad_connector_office_site_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_:-]{1,254}$`), "It must be `2` to `255` characters in length, The name must start with a letter, and can contain letters, digits, colons (:), underscores (_), and hyphens (-).")),
			},
			"ad_hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 200),
			},
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cen_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"desktop_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ANY", "INTERNET", "VPC"}, false),
			},
			"dns_address": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"domain_user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_admin_access": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enable_internet_access": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"mfa_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ASP", "HDX"}, false),
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
			},
			"sub_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"verify_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudEcdAdConnectorOfficeSiteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateADConnectorOfficeSite"
	request := make(map[string]interface{})
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request["OfficeSiteName"] = d.Get("ad_connector_office_site_name")

	if v, ok := d.GetOk("ad_hostname"); ok {
		request["AdHostname"] = v
	}
	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	request["CenId"] = d.Get("cen_id")
	if v, ok := d.GetOk("cen_owner_id"); ok {
		request["CenOwnerId"] = v
	}
	request["CidrBlock"] = d.Get("cidr_block")
	if v, ok := d.GetOk("desktop_access_type"); ok {
		request["DesktopAccessType"] = v
	}
	request["DnsAddress"] = d.Get("dns_address")
	request["DomainName"] = d.Get("domain_name")
	if v, ok := d.GetOk("domain_password"); ok {
		request["DomainPassword"] = v
	}
	if v, ok := d.GetOk("domain_user_name"); ok {
		request["DomainUserName"] = v
	}
	if v, ok := d.GetOkExists("enable_admin_access"); ok {
		request["EnableAdminAccess"] = v
	}
	if v, ok := d.GetOkExists("enable_internet_access"); ok {
		request["EnableInternetAccess"] = v
	}
	if v, ok := d.GetOkExists("mfa_enabled"); ok {
		request["MfaEnabled"] = v
	}
	if v, ok := d.GetOk("protocol_type"); ok {
		request["ProtocolType"] = v
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
	if v, ok := d.GetOk("verify_code"); ok {
		request["VerifyCode"] = v
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_ad_connector_office_site", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["OfficeSiteId"]))
	ecdService := EcdService{client}
	stateConf := BuildStateConf([]string{"REGISTERING", "CONFIGTRUSTING"}, []string{"REGISTERED", "ERROR", "NEEDCONFIGTRUST", "NEEDCONFIGUSER", "CONFIGTRUSTFAILED"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecdService.EcdAdConnectorOfficeSiteStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcdAdConnectorOfficeSiteRead(d, meta)
}
func resourceAlicloudEcdAdConnectorOfficeSiteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	object, err := ecdService.DescribeEcdAdConnectorOfficeSite(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_ad_connector_office_site ecdService.DescribeEcdAdConnectorOfficeSite Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ad_connector_office_site_name", object["Name"])
	if v, ok := object["Bandwidth"]; ok && fmt.Sprint(v) != "0" {
		d.Set("bandwidth", formatInt(v))
	}
	d.Set("cen_id", object["CenId"])
	d.Set("cidr_block", object["CidrBlock"])
	d.Set("desktop_access_type", object["DesktopAccessType"])
	d.Set("dns_address", object["DnsAddress"])
	d.Set("domain_name", object["DomainName"])
	d.Set("domain_user_name", object["DomainUserName"])
	d.Set("enable_admin_access", object["EnableAdminAccess"])
	d.Set("enable_internet_access", object["EnableInternetAccess"])
	d.Set("mfa_enabled", object["MfaEnabled"])
	d.Set("status", object["Status"])
	d.Set("sub_domain_dns_address", object["SubDnsAddress"])
	d.Set("sub_domain_name", object["SubDomainName"])
	return nil
}
func resourceAlicloudEcdAdConnectorOfficeSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudEcdAdConnectorOfficeSiteRead(d, meta)
}
func resourceAlicloudEcdAdConnectorOfficeSiteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteOfficeSites"
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"OfficeSiteId": []string{d.Id()},
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
	ecdService := EcdService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecdService.EcdAdConnectorOfficeSiteStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
