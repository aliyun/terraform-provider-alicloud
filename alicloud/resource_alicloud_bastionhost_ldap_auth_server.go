package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudBastionhostLdapAuthServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBastionhostLdapAuthServerCreate,
		Read:   resourceAlicloudBastionhostLdapAuthServerRead,
		Update: resourceAlicloudBastionhostLdapAuthServerUpdate,
		Delete: resourceAlicloudBastionhostLdapAuthServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"login_name_mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mobile_mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"server": {
				Type:     schema.TypeString,
				Required: true,
			},
			"standby_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudBastionhostLdapAuthServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyInstanceLDAPAuthServer"
	request := make(map[string]interface{})
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("account"); ok {
		request["Account"] = v
	}
	if v, ok := d.GetOk("base_dn"); ok {
		request["BaseDN"] = v
	}
	if v, ok := d.GetOk("email_mapping"); ok {
		request["EmailMapping"] = v
	}
	if v, ok := d.GetOk("filter"); ok {
		request["Filter"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOkExists("is_ssl"); ok {
		request["IsSSL"] = v
	}
	if v, ok := d.GetOk("login_name_mapping"); ok {
		request["LoginNameMapping"] = v
	}
	if v, ok := d.GetOk("mobile_mapping"); ok {
		request["MobileMapping"] = v
	}
	if v, ok := d.GetOk("name_mapping"); ok {
		request["NameMapping"] = v
	}
	request["Password"] = d.Get("password")
	if v, ok := d.GetOk("port"); ok {
		request["Port"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("server"); ok {
		request["Server"] = v
	}
	if v, ok := d.GetOk("standby_server"); ok {
		request["StandbyServer"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_ldap_auth_server", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"]))

	return resourceAlicloudBastionhostLdapAuthServerRead(d, meta)
}
func resourceAlicloudBastionhostLdapAuthServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := YundunBastionhostService{client}
	object, err := yundunBastionhostService.DescribeBastionhostLdapAuthServer(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_ldap_auth_server yundunBastionhostService.DescribeBastionhostLdapAuthServer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	d.Set("account", object["Account"])
	d.Set("base_dn", object["BaseDN"])
	d.Set("email_mapping", object["EmailMapping"])
	d.Set("filter", object["Filter"])
	d.Set("is_ssl", object["IsSSL"])
	d.Set("login_name_mapping", object["LoginNameMapping"])
	d.Set("mobile_mapping", object["MobileMapping"])
	d.Set("name_mapping", object["NameMapping"])
	d.Set("port", formatInt(object["Port"]))
	d.Set("server", object["Server"])
	d.Set("standby_server", object["StandbyServer"])
	return nil
}
func resourceAlicloudBastionhostLdapAuthServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if v, ok := d.GetOk("account"); ok {
		request["Account"] = v
	}
	if d.HasChange("account") {
		update = true
	}
	if v, ok := d.GetOk("base_dn"); ok {
		request["BaseDN"] = v
	}
	if d.HasChange("base_dn") {
		update = true
	}
	if v, ok := d.GetOk("email_mapping"); ok {
		request["EmailMapping"] = v
	}
	if d.HasChange("email_mapping") {
		update = true
	}
	if v, ok := d.GetOk("filter"); ok {
		request["Filter"] = v
	}
	if d.HasChange("filter") {
		update = true
	}
	if v, ok := d.GetOkExists("is_ssl"); ok {
		request["IsSSL"] = v
	}
	if d.HasChange("is_ssl") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOk("login_name_mapping"); ok {
		request["LoginNameMapping"] = v
	}
	if d.HasChange("login_name_mapping") {
		update = true
	}
	if v, ok := d.GetOk("mobile_mapping"); ok {
		request["MobileMapping"] = v
	}
	if d.HasChange("mobile_mapping") {
		update = true
	}
	if v, ok := d.GetOk("name_mapping"); ok {
		request["NameMapping"] = v
	}
	if d.HasChange("name_mapping") {
		update = true
	}
	request["Password"] = d.Get("password")

	if d.HasChange("password") {
		update = true
	}
	if v, ok := d.GetOk("port"); ok {
		request["Port"] = v
	}
	if d.HasChange("port") {
		update = true
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("server"); ok {
		request["Server"] = v
	}
	if d.HasChange("server") {
		update = true
	}
	if v, ok := d.GetOk("standby_server"); ok {
		request["StandbyServer"] = v
	}
	if d.HasChange("standby_server") {
		update = true
	}
	if update {
		action := "ModifyInstanceLDAPAuthServer"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudBastionhostLdapAuthServerRead(d, meta)
}
func resourceAlicloudBastionhostLdapAuthServerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudBastionhostLdapAuthServer. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
