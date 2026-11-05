package ddoscoo

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDdoscooDomainResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDdosCooDomainResourceCreate,
		Read:   resourceAliCloudDdosCooDomainResourceRead,
		Update: resourceAliCloudDdosCooDomainResourceUpdate,
		Delete: resourceAliCloudDdosCooDomainResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ai_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"watch", "defense"}, false),
			},
			"ai_template": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"level30", "level60", "level90"}, false),
			},
			"black_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bw_list_enable": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.IntInSlice([]int{0, 1}),
			},
			"cc_global_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"close", "open"}, false),
			},
			"cert": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"cert_identifier": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"cert_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cert_region": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_ciphers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"custom_headers": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"https_ext": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ocsp_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"proxy_types": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_ports": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"proxy_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"real_servers": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rs_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ssl13_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_ciphers": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"default", "all", "strong", "improved"}, false),
			},
			"ssl_protocols": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"tls1.0", "tls1.1", "tls1.2"}, false),
			},
			"tls13_custom_ciphers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"white_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudDdosCooDomainResourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDomainResource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsMapsArray := convert.ConvertToInterfaceArray(v)

		request["InstanceIds"] = instanceIdsMapsArray
	}

	if v, ok := d.GetOk("proxy_types"); ok {
		proxyTypesMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convert.ConvertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["ProxyPorts"] = convert.ConvertToInterfaceArray(dataLoop1Tmp["proxy_ports"])
			dataLoop1Map["ProxyType"] = dataLoop1Tmp["proxy_type"]
			proxyTypesMapsArray = append(proxyTypesMapsArray, dataLoop1Map)
		}
		request["ProxyTypes"] = proxyTypesMapsArray
	}

	request["RsType"] = d.Get("rs_type")
	if v, ok := d.GetOk("real_servers"); ok {
		realServersMapsArray := convert.ConvertToInterfaceArray(v)

		request["RealServers"] = realServersMapsArray
	}

	if v, ok := d.GetOk("https_ext"); ok {
		request["HttpsExt"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ddoscoo_domain_resource", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Domain"]))

	return resourceAliCloudDdosCooDomainResourceUpdate(d, meta)
}

func resourceAliCloudDdosCooDomainResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosCooServiceV2 := NewService(client)

	objectRaw, err := ddosCooServiceV2.DescribeDdosCooDomainResource(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddoscoo_domain_resource DescribeDdosCooDomainResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("cname", objectRaw["Cname"])
	d.Set("https_ext", objectRaw["HttpsExt"])
	ocsp_enabledVal, formatErr := convert.FormatBool(convertDdosCooDomainResourceWebRulesOcspEnabledResponse(objectRaw["OcspEnabled"]))
	if formatErr != nil {
		return formatErr
	}
	d.Set("ocsp_enabled", ocsp_enabledVal)
	d.Set("rs_type", objectRaw["RsType"])
	d.Set("ssl13_enabled", objectRaw["Ssl13Enabled"])
	d.Set("ssl_ciphers", objectRaw["SslCiphers"])
	d.Set("ssl_protocols", objectRaw["SslProtocols"])
	d.Set("domain", objectRaw["Domain"])

	blackListRaw := make([]interface{}, 0)
	if objectRaw["BlackList"] != nil {
		blackListRaw = convert.ConvertToInterfaceArray(objectRaw["BlackList"])
	}

	d.Set("black_list", blackListRaw)
	customCiphersRaw := make([]interface{}, 0)
	if objectRaw["CustomCiphers"] != nil {
		customCiphersRaw = convert.ConvertToInterfaceArray(objectRaw["CustomCiphers"])
	}

	d.Set("custom_ciphers", customCiphersRaw)
	instanceIdsRaw := make([]interface{}, 0)
	if objectRaw["InstanceIds"] != nil {
		instanceIdsRaw = convert.ConvertToInterfaceArray(objectRaw["InstanceIds"])
	}

	d.Set("instance_ids", instanceIdsRaw)
	proxyTypesRaw := objectRaw["ProxyTypes"]
	proxyTypesMaps := make([]map[string]interface{}, 0)
	if proxyTypesRaw != nil {
		for _, proxyTypesChildRaw := range convert.ConvertToInterfaceArray(proxyTypesRaw) {
			proxyTypesMap := make(map[string]interface{})
			proxyTypesChildRaw := proxyTypesChildRaw.(map[string]interface{})
			proxyTypesMap["proxy_type"] = proxyTypesChildRaw["ProxyType"]

			proxyPortsRaw := make([]interface{}, 0)
			if proxyTypesChildRaw["ProxyPorts"] != nil {
				proxyPortsRaw = convert.ConvertToInterfaceArray(proxyTypesChildRaw["ProxyPorts"])
			}

			proxyTypesMap["proxy_ports"] = proxyPortsRaw
			proxyTypesMaps = append(proxyTypesMaps, proxyTypesMap)
		}
	}
	if err := d.Set("proxy_types", proxyTypesMaps); err != nil {
		return err
	}
	realServersRaw := make([]interface{}, 0)
	if objectRaw["RealServers"] != nil {
		realServersRaw = convert.ConvertToInterfaceArray(objectRaw["RealServers"])
	}

	d.Set("real_servers", realServersRaw)
	whiteListRaw := make([]interface{}, 0)
	if objectRaw["WhiteList"] != nil {
		whiteListRaw = convert.ConvertToInterfaceArray(objectRaw["WhiteList"])
	}

	d.Set("white_list", whiteListRaw)

	objectRaw, err = ddosCooServiceV2.DescribeDomainResourceDescribeWebRules(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	d.Set("cert_name", objectRaw["UserCertName"])
	d.Set("domain", objectRaw["Domain"])

	tls13CustomCiphersRaw := make([]interface{}, 0)
	if objectRaw["Tls13CustomCiphers"] != nil {
		tls13CustomCiphersRaw = convert.ConvertToInterfaceArray(objectRaw["Tls13CustomCiphers"])
	}

	d.Set("tls13_custom_ciphers", tls13CustomCiphersRaw)

	objectRaw, err = ddosCooServiceV2.DescribeDomainResourceDescribeDomainCcProtectSwitch(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	d.Set("ai_mode", objectRaw["AiMode"])
	d.Set("ai_template", objectRaw["AiTemplate"])
	d.Set("bw_list_enable", objectRaw["BlackWhiteListEnable"])
	d.Set("cc_global_switch", objectRaw["CcGlobalSwitch"])

	objectRaw, err = ddosCooServiceV2.DescribeDomainResourceDescribeHeaders(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	d.Set("custom_headers", convert.ConvertObjectToJsonString(objectRaw["Headers"]))

	return nil
}

func resourceAliCloudDdosCooDomainResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyDomainResource"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()
	if !d.IsNewResource() && d.HasChange("instance_ids") {
		update = true
	}
	if v, ok := d.GetOk("instance_ids"); ok || d.HasChange("instance_ids") {
		instanceIdsMapsArray := convert.ConvertToInterfaceArray(v)

		request["InstanceIds"] = instanceIdsMapsArray
	}

	if !d.IsNewResource() && d.HasChange("proxy_types") {
		update = true
	}
	if v, ok := d.GetOk("proxy_types"); ok || d.HasChange("proxy_types") {
		proxyTypesMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convert.ConvertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["ProxyPorts"] = convert.ConvertToInterfaceArray(dataLoop1Tmp["proxy_ports"])
			dataLoop1Map["ProxyType"] = dataLoop1Tmp["proxy_type"]
			proxyTypesMapsArray = append(proxyTypesMapsArray, dataLoop1Map)
		}
		request["ProxyTypes"] = proxyTypesMapsArray
	}

	if !d.IsNewResource() && d.HasChange("real_servers") {
		update = true
	}
	if v, ok := d.GetOk("real_servers"); ok || d.HasChange("real_servers") {
		realServersMapsArray := convert.ConvertToInterfaceArray(v)

		request["RealServers"] = realServersMapsArray
	}

	if !d.IsNewResource() && d.HasChange("rs_type") {
		update = true
	}
	request["RsType"] = d.Get("rs_type")
	if !d.IsNewResource() && d.HasChange("https_ext") {
		update = true
		request["HttpsExt"] = d.Get("https_ext")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifyOcspStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()

	if d.HasChange("ocsp_enabled") {
		update = true
	}
	request["Enable"] = convertDdosCooDomainResourceEnableRequest(d.Get("ocsp_enabled").(bool))
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "AssociateWebCert"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()
	if d.HasChange("cert_region") {
		update = true
		request["CertRegion"] = d.Get("cert_region")
	}

	if d.HasChange("cert_identifier") {
		update = true
		request["CertIdentifier"] = d.Get("cert_identifier")
	}

	if d.HasChange("key") {
		update = true
		request["Key"] = d.Get("key")
	}

	if d.HasChange("cert") {
		update = true
		request["Cert"] = d.Get("cert")
	}

	if d.HasChange("cert_name") {
		update = true
		request["CertName"] = d.Get("cert_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifyHeaders"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("custom_headers") {
		update = true
	}
	request["CustomHeaders"] = d.Get("custom_headers")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifyWebAIProtectMode"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()
	request["RegionId"] = client.RegionId
	config := make(map[string]interface{})

	if d.HasChange("ai_template") {
		update = true
	}
	if v, ok := d.GetOk("ai_template"); ok {
		config["AiTemplate"] = v
	}

	if d.HasChange("ai_mode") {
		update = true
	}
	if v, ok := d.GetOk("ai_mode"); ok {
		config["AiMode"] = v
	}

	request["Config"] = convert.ConvertObjectToJsonString(config)

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ConfigDomainSecurityProfile"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()

	ddosCooServiceV2 := NewService(client)

	objectRaw, err := ddosCooServiceV2.DescribeDomainResourceDescribeDomainCcProtectSwitch(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	v, ok := d.GetOkExists("bw_list_enable")
	if fmt.Sprint(objectRaw["BlackWhiteListEnable"]) != fmt.Sprint(d.Get("bw_list_enable")) && ok {
		update = true

		config = make(map[string]interface{})
		config["bwlist_enable"] = v
		request["Config"] = convert.ConvertObjectToJsonString(config)
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ConfigWebIpSet"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("black_list") {
		update = true
	}
	if v, ok := d.GetOk("black_list"); ok || d.HasChange("black_list") {
		blackListMapsArray := convert.ConvertToInterfaceArray(v)

		request["BlackList"] = blackListMapsArray
	}

	if d.HasChange("white_list") {
		update = true
	}
	if v, ok := d.GetOk("white_list"); ok || d.HasChange("white_list") {
		whiteListMapsArray := convert.ConvertToInterfaceArray(v)

		request["WhiteList"] = whiteListMapsArray
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifyWebCCGlobalSwitch"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("cc_global_switch") {
		update = true
	}
	request["CcGlobalSwitch"] = d.Get("cc_global_switch")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifyTlsConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = d.Id()
	request["RegionId"] = client.RegionId
	config = make(map[string]interface{})

	if d.HasChange("custom_ciphers") {
		update = true
	}
	if v, ok := d.GetOk("custom_ciphers"); ok {
		config["custom_ciphers"] = convert.ConvertToInterfaceArray(v)
	}

	if d.HasChange("ssl13_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("ssl13_enabled"); ok {
		config["ssl_13_enabled"] = v
	}

	if d.HasChange("tls13_custom_ciphers") {
		update = true
	}
	if v, ok := d.GetOk("tls13_custom_ciphers"); ok {
		config["tls_13_custom_ciphers"] = convert.ConvertToInterfaceArray(v)
	}

	if d.HasChange("ssl_protocols") {
		update = true
	}
	if v, ok := d.GetOk("ssl_protocols"); ok {
		config["ssl_protocols"] = v
	}

	if d.HasChange("ssl_ciphers") {
		update = true
	}
	if v, ok := d.GetOk("ssl_ciphers"); ok {
		config["ssl_ciphers"] = v
	}

	request["Config"] = convert.ConvertObjectToJsonString(config)

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAliCloudDdosCooDomainResourceRead(d, meta)
}

func resourceAliCloudDdosCooDomainResourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDomainResource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Domain"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertDdosCooDomainResourceWebRulesOcspEnabledResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	}
	return source
}

func convertDdosCooDomainResourceEnableRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "true":
		return "1"
	case "false":
		return "0"
	}
	return source
}
