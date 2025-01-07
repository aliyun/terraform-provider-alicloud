package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudBastionhostInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBastionhostInstanceCreate,
		Read:   resourceAlicloudBastionhostInstanceRead,
		Update: resourceAlicloudBastionhostInstanceUpdate,
		Delete: resourceAlicloudBastionhostInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"license_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plan_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloudbastion", "cloudbastion_ha"}, false),
			},
			"storage": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_white_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 6, 12, 24, 36}),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_public_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ad_auth_server": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					d1, d2 := d.GetChange("ad_auth_server")
					if len(d1.(*schema.Set).List()) == 0 || len(d2.(*schema.Set).List()) == 0 {
						return false
					}
					return compareMapWithIgnoreEquivalent(d1.(*schema.Set).List()[0].(map[string]interface{}), d2.(*schema.Set).List()[0].(map[string]interface{}), []string{"password"})
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account": {
							Type:     schema.TypeString,
							Required: true,
						},
						"base_dn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"domain": {
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
						"is_ssl": {
							Type:     schema.TypeBool,
							Required: true,
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
							Optional:  true,
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
				},
			},
			"ldap_auth_server": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					d1, d2 := d.GetChange("ldap_auth_server")
					if len(d1.(*schema.Set).List()) == 0 || len(d2.(*schema.Set).List()) == 0 {
						return false
					}
					return compareMapWithIgnoreEquivalent(d1.(*schema.Set).List()[0].(map[string]interface{}), d2.(*schema.Set).List()[0].(map[string]interface{}), []string{"password"})
				},
				Elem: &schema.Resource{
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
							Optional:  true,
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
				},
			},
			"renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.Any(validation.IntBetween(1, 9), validation.IntInSlice([]int{12, 24, 36})),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
						return false
					}
					return true
				},
			},
			"renewal_period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"M", "Y"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
						return false
					}
					return true
				},
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AutoRenewal", "ManualRenewal", "NotRenewal"}, false),
			},
		},
	}
}

func resourceAlicloudBastionhostInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	var endpoint string
	action := "CreateInstance"
	request := make(map[string]interface{})
	parameterMapList := make([]map[string]interface{}, 0)
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "NetworkType",
		"Value": "vpc",
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "LicenseCode",
		"Value": d.Get("license_code").(string),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "PlanCode",
		"Value": d.Get("plan_code").(string),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Storage",
		"Value": d.Get("storage").(string),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Bandwidth",
		"Value": d.Get("bandwidth").(string),
	})
	request["SubscriptionType"] = "Subscription"
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}

	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renew_period", "renewal_status", d.Get("renewal_status")))
	}
	request["ProductCode"] = "bastionhost"
	request["ProductType"] = "bastionhost"
	if client.IsInternationalAccount() {
		request["ProductType"] = "bastionhost_std_public_intl"
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "RegionId",
		"Value": client.RegionId,
	})
	request["Parameter"] = parameterMapList
	request["ClientToken"] = buildClientToken("CreateInstance")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "bastionhost_std_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_instance", action, AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["InstanceId"]))

	bastionhostService := YundunBastionhostService{client}

	// check RAM policy
	if err := bastionhostService.ProcessRolePolicy(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	// wait for order complete
	stateConf := BuildStateConf([]string{}, []string{"PENDING"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"UPGRADING", "UPGRADE_FAILED", "CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	rawSecurityGroupIds := d.Get("security_group_ids").(*schema.Set).List()
	securityGroupIds := make([]string, len(rawSecurityGroupIds))
	for index, rawSecurityGroupId := range rawSecurityGroupIds {
		securityGroupIds[index] = rawSecurityGroupId.(string)
	}
	// start instance
	if err := bastionhostService.StartBastionhostInstance(d.Id(), d.Get("vswitch_id").(string), securityGroupIds); err != nil {
		return WrapError(err)
	}
	// wait for pending
	stateConf = BuildStateConf([]string{"PENDING", "CREATING"}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 600*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"UPGRADING", "UPGRADE_FAILED", "CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudBastionhostInstanceUpdate(d, meta)
}

func resourceAlicloudBastionhostInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	BastionhostService := YundunBastionhostService{client}
	instance, err := BastionhostService.DescribeBastionhostInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", instance["Description"])
	d.Set("license_code", instance["LicenseCode"])
	d.Set("vswitch_id", instance["VswitchId"])
	d.Set("security_group_ids", instance["AuthorizedSecurityGroups"])
	d.Set("enable_public_access", instance["PublicNetworkAccess"])
	d.Set("resource_group_id", instance["ResourceGroupId"])
	// instance["Storage"] is in byte, and it is larger than request param
	if v, err := strconv.ParseInt(instance["Storage"].(json.Number).String(), 10, 64); err != nil {
		return WrapError(err)
	} else {
		d.Set("storage", fmt.Sprint(bytesToTB(v)-1))
	}

	d.Set("bandwidth", instance["BandwidthPackage"])

	if fmt.Sprint(instance["PublicNetworkAccess"]) == "true" {
		d.Set("public_white_list", instance["PublicWhiteList"])
	}

	instance, err = BastionhostService.DescribeBastionhostInstances(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("plan_code", instance["PlanCode"])

	tags, err := BastionhostService.ListTagResources(d.Id(), "INSTANCE")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(tags))

	adAuthServer, err := BastionhostService.DescribeBastionhostAdAuthServer(d.Id())
	if err != nil {
		return WrapError(err)
	}
	adAuthServerMap := map[string]interface{}{
		"account":        adAuthServer["Account"],
		"base_dn":        adAuthServer["BaseDN"],
		"domain":         adAuthServer["Domain"],
		"email_mapping":  adAuthServer["EmailMapping"],
		"filter":         adAuthServer["Filter"],
		"is_ssl":         adAuthServer["IsSSL"],
		"mobile_mapping": adAuthServer["MobileMapping"],
		"name_mapping":   adAuthServer["NameMapping"],
		"port":           formatInt(adAuthServer["Port"]),
		"server":         adAuthServer["Server"],
		"standby_server": adAuthServer["StandbyServer"],
		"has_password":   adAuthServer["HasPassword"],
	}
	d.Set("ad_auth_server", []map[string]interface{}{adAuthServerMap})

	ldapAuthServer, err := BastionhostService.DescribeBastionhostLdapAuthServer(d.Id())
	if err != nil {
		return WrapError(err)
	}
	ldapAuthServerMap := map[string]interface{}{
		"account":            ldapAuthServer["Account"],
		"base_dn":            ldapAuthServer["BaseDN"],
		"email_mapping":      ldapAuthServer["EmailMapping"],
		"filter":             ldapAuthServer["Filter"],
		"is_ssl":             ldapAuthServer["IsSSL"],
		"login_name_mapping": ldapAuthServer["LoginNameMapping"],
		"mobile_mapping":     ldapAuthServer["MobileMapping"],
		"name_mapping":       ldapAuthServer["NameMapping"],
		"port":               ldapAuthServer["Port"],
		"server":             ldapAuthServer["Server"],
		"standby_server":     ldapAuthServer["StandbyServer"],
		"has_password":       adAuthServer["HasPassword"],
	}
	d.Set("ldap_auth_server", []map[string]interface{}{ldapAuthServerMap})

	bssOpenApiService := BssOpenApiService{client}
	// can not set region when invoking QueryAvailableInstances for bastionhost instance
	getQueryInstanceObject, err := bssOpenApiService.QueryAvailableInstances(d.Id(), "", "bastionhost", "bastionhost", "bastionhost", "bastionhost_std_public_intl")
	if err != nil {
		return WrapError(err)
	}

	d.Set("renewal_status", getQueryInstanceObject["RenewStatus"])
	if v, ok := getQueryInstanceObject["RenewalDuration"]; ok && v != nil {
		d.Set("renew_period", getQueryInstanceObject["RenewalDuration"])
	}
	d.Set("renewal_period_unit", getQueryInstanceObject["RenewalDurationUnit"])

	return nil
}

func resourceAlicloudBastionhostInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bastionhostService := YundunBastionhostService{client}
	var err error

	d.Partial(true)

	if d.HasChange("tags") {
		if err := bastionhostService.setInstanceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if d.HasChange("description") {
		if err := bastionhostService.UpdateBastionhostInstanceDescription(d.Id(), d.Get("description").(string)); err != nil {
			return WrapError(err)
		}
		d.SetPartial("description")
	}

	if d.HasChange("resource_group_id") {
		if err := bastionhostService.UpdateResourceGroup(d.Id(), "INSTANCE", d.Get("resource_group_id").(string)); err != nil {
			return WrapError(err)
		}
		d.SetPartial("resource_group_id")
	}

	if !d.IsNewResource() && d.HasChange("license_code") {
		params := map[string]string{
			"LicenseCode": "license_code",
		}
		if err := bastionhostService.UpdateInstanceSpec(params, d, meta); err != nil {
			return WrapError(err)
		}
		stateConf := BuildStateConf([]string{"UPGRADING"}, []string{"PENDING", "RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"CREATING", "UPGRADE_FAILED", "CREATE_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("license_code")
	}

	if !d.IsNewResource() && d.HasChange("security_group_ids") {
		securityGroupIds := d.Get("security_group_ids").(*schema.Set).List()
		sgs := make([]string, 0, len(securityGroupIds))
		for _, rawSecurityGroupId := range securityGroupIds {
			sgs = append(sgs, rawSecurityGroupId.(string))
		}
		if err := bastionhostService.UpdateBastionhostSecurityGroups(d.Id(), sgs); err != nil {
			return WrapError(err)
		}
		stateConf := BuildStateConf([]string{"UPGRADING"}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"CREATING", "UPGRADE_FAILED", "CREATE_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_group_ids")
	}

	_, ok := d.GetOkExists("enable_public_access")
	if d.HasChange("enable_public_access") || (d.IsNewResource() && ok) {
		client := meta.(*connectivity.AliyunClient)
		BastionhostService := YundunBastionhostService{client}
		instance, err := BastionhostService.DescribeBastionhostInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := strconv.FormatBool(d.Get("enable_public_access").(bool))
		if strconv.FormatBool(instance["PublicNetworkAccess"].(bool)) != target {
			if target == "false" {
				err := BastionhostService.DisableInstancePublicAccess(d.Id())
				if err != nil {
					return WrapError(err)
				}
			} else {
				err := BastionhostService.EnableInstancePublicAccess(d.Id())
				if err != nil {
					return WrapError(err)
				}
			}
		}

		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("enable_public_access")
	}

	if d.HasChange("ad_auth_server") {
		if v, ok := d.GetOk("ad_auth_server"); ok && len(v.(*schema.Set).List()) > 0 {
			var response map[string]interface{}
			modifyAdRequest := map[string]interface{}{
				"InstanceId": d.Id(),
				"RegionId":   client.RegionId,
			}
			adAuthServer := v.(*schema.Set).List()[0].(map[string]interface{})
			modifyAdRequest["Account"] = adAuthServer["account"]
			modifyAdRequest["BaseDN"] = adAuthServer["base_dn"]
			modifyAdRequest["Domain"] = adAuthServer["domain"]
			modifyAdRequest["IsSSL"] = adAuthServer["is_ssl"]
			modifyAdRequest["Port"] = adAuthServer["port"]
			modifyAdRequest["Server"] = adAuthServer["server"]
			modifyAdRequest["EmailMapping"] = adAuthServer["email_mapping"]
			modifyAdRequest["Filter"] = adAuthServer["filter"]
			modifyAdRequest["MobileMapping"] = adAuthServer["mobile_mapping"]
			modifyAdRequest["NameMapping"] = adAuthServer["name_mapping"]
			modifyAdRequest["Password"] = adAuthServer["password"]
			modifyAdRequest["StandbyServer"] = adAuthServer["standby_server"]

			action := "ModifyInstanceADAuthServer"
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, modifyAdRequest, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, modifyAdRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			d.SetPartial("ad_auth_server")
		}
	}

	if d.HasChange("ldap_auth_server") {
		if v, ok := d.GetOk("ldap_auth_server"); ok && len(v.(*schema.Set).List()) > 0 {
			var response map[string]interface{}
			modifyLdapRequest := map[string]interface{}{
				"InstanceId": d.Id(),
				"RegionId":   client.RegionId,
			}

			adAuthServer := v.(*schema.Set).List()[0].(map[string]interface{})
			modifyLdapRequest["Account"] = adAuthServer["account"]
			modifyLdapRequest["BaseDN"] = adAuthServer["base_dn"]
			modifyLdapRequest["Port"] = adAuthServer["port"]
			modifyLdapRequest["Server"] = adAuthServer["server"]
			modifyLdapRequest["Password"] = adAuthServer["password"]
			modifyLdapRequest["IsSSL"] = adAuthServer["is_ssl"]
			modifyLdapRequest["LoginNameMapping"] = adAuthServer["login_name_mapping"]
			modifyLdapRequest["EmailMapping"] = adAuthServer["email_mapping"]
			modifyLdapRequest["Filter"] = adAuthServer["filter"]
			modifyLdapRequest["MobileMapping"] = adAuthServer["mobile_mapping"]
			modifyLdapRequest["NameMapping"] = adAuthServer["name_mapping"]
			modifyLdapRequest["StandbyServer"] = adAuthServer["standby_server"]

			action := "ModifyInstanceLDAPAuthServer"
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, modifyLdapRequest, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, modifyLdapRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			d.SetPartial("ldap_auth_server")
		}
	}

	var setRenewalResponse map[string]interface{}
	var endpoint string
	update := false
	setRenewalReq := map[string]interface{}{
		"InstanceIDs":      d.Id(),
		"ProductCode":      "bastionhost",
		"ProductType":      "bastionhost",
		"SubscriptionType": "Subscription",
	}
	if client.IsInternationalAccount() {
		setRenewalReq["ProductType"] = "bastionhost_std_public_intl"
	}

	if !d.IsNewResource() && d.HasChange("renewal_status") {
		update = true
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		setRenewalReq["RenewalStatus"] = v
	}

	if !d.IsNewResource() && d.HasChange("renew_period") {
		update = true
		if v, ok := d.GetOk("renew_period"); ok {
			setRenewalReq["RenewalPeriod"] = v
		}
	}

	if d.HasChange("renewal_period_unit") {
		update = true
	}
	if v, ok := d.GetOk("renewal_period_unit"); ok {
		setRenewalReq["RenewalPeriodUnit"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renewal_period_unit", "renewal_status", d.Get("renewal_status")))
	}

	if update {
		action := "SetRenewal"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			setRenewalResponse, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, setRenewalReq, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					setRenewalReq["ProductType"] = "bastionhost_std_public_intl"
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, setRenewalResponse, setRenewalReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(setRenewalResponse["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, setRenewalResponse))
		}

		d.SetPartial("renewal_status")
		d.SetPartial("renew_period")
		d.SetPartial("renewal_period_unit")
	}

	update = false
	configInstanceWhiteListReq := map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	if d.HasChange("public_white_list") {
		update = true
	}
	if v, ok := d.GetOk("public_white_list"); ok {
		configInstanceWhiteListReq["WhiteList"] = v
	}

	if update {
		action := "ConfigInstanceWhiteList"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, configInstanceWhiteListReq, false)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, configInstanceWhiteListReq)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("public_white_list")
	}

	d.Partial(false)
	// wait for order complete
	return resourceAlicloudBastionhostInstanceRead(d, meta)
}

func resourceAlicloudBastionhostInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceBastionhostInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
