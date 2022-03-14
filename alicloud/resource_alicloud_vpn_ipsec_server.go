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

func resourceAlicloudVpnIpsecServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpnIpsecServerCreate,
		Read:   resourceAlicloudVpnIpsecServerRead,
		Update: resourceAlicloudVpnIpsecServerUpdate,
		Delete: resourceAlicloudVpnIpsecServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"client_ip_pool": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"effect_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipsec_server_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{1,127}`), "The name must be `2` to `128` characters in length and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter."),
			},
			"local_subnet": {
				Type:     schema.TypeString,
				Required: true,
			},
			"psk": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"psk_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ike_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike_version": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ikev1", "ikev2"}, false),
						},
						"ike_mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_enc_alg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_auth_alg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_pfs": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ike_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
						"local_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"remote_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ipsec_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipsec_enc_alg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipsec_auth_alg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipsec_pfs": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipsec_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudVpnIpsecServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateIpsecServer"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request["ClientIpPool"] = d.Get("client_ip_pool")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("effect_immediately"); ok {
		request["EffectImmediately"] = v
	}
	if v, ok := d.GetOk("ipsec_server_name"); ok {
		request["IpSecServerName"] = v
	}
	request["LocalSubnet"] = d.Get("local_subnet")
	if v, ok := d.GetOk("psk"); ok {
		request["Psk"] = v
	}
	if v, ok := d.GetOkExists("psk_enabled"); ok {
		request["PskEnabled"] = v
	}

	if v, ok := d.GetOk("ike_config"); ok {
		ikeConfig := make(map[string]interface{})
		for _, ikeConfigArgs := range v.(*schema.Set).List() {
			ikeConfigArg := ikeConfigArgs.(map[string]interface{})
			ikeConfig["IkeVersion"] = ikeConfigArg["ike_version"]
			ikeConfig["IkeMode"] = ikeConfigArg["ike_mode"]
			ikeConfig["IkeEncAlg"] = ikeConfigArg["ike_enc_alg"]
			ikeConfig["IkeAuthAlg"] = ikeConfigArg["ike_auth_alg"]
			ikeConfig["IkePfs"] = ikeConfigArg["ike_pfs"]
			ikeConfig["IkeLifetime"] = ikeConfigArg["ike_lifetime"]
			ikeConfig["LocalId"] = ikeConfigArg["local_id"]
			ikeConfig["RemoteId"] = ikeConfigArg["remote_id"]
		}

		if v, err := convertMaptoJsonString(ikeConfig); err != nil {
			return WrapError(err)
		} else {
			request["IkeConfig"] = v
		}
	}

	if v, ok := d.GetOk("ipsec_config"); ok {
		ipsecConfig := make(map[string]interface{})
		for _, ipsecConfigArgs := range v.(*schema.Set).List() {
			ipsecConfigArg := ipsecConfigArgs.(map[string]interface{})
			ipsecConfig["IpsecEncAlg"] = ipsecConfigArg["IpsecEncAlg"]
			ipsecConfig["IpsecAuthAlg"] = ipsecConfigArg["IpsecAuthAlg"]
			ipsecConfig["IpsecPfs"] = ipsecConfigArg["IpsecPfs"]
			ipsecConfig["IpsecLifetime"] = ipsecConfigArg["IpsecLifetime"]
		}

		if v, err := convertMaptoJsonString(ipsecConfig); err != nil {
			return WrapError(err)
		} else {
			request["IpsecConfig"] = v
		}
	}

	request["RegionId"] = client.RegionId
	request["VpnGatewayId"] = d.Get("vpn_gateway_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateIpsecServer")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_ipsec_server", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["IpsecServerId"]))

	return resourceAlicloudVpnIpsecServerRead(d, meta)
}
func resourceAlicloudVpnIpsecServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpnIpsecServer(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_ipsec_server vpcService.DescribeVpnIpsecServer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("client_ip_pool", object["ClientIpPool"])
	d.Set("effect_immediately", object["EffectImmediately"])
	d.Set("ipsec_server_name", object["IpsecServerName"])
	d.Set("local_subnet", object["LocalSubnet"])
	d.Set("psk", object["Psk"])
	d.Set("psk_enabled", object["PskEnabled"])
	d.Set("vpn_gateway_id", object["VpnGatewayId"])
	if v, ok := object["IkeConfig"]; ok {
		ikeConfigSli := make([]map[string]interface{}, 0)
		if len(v.(map[string]interface{})) > 0 {
			if ikeConfigArg, ok := object["IkeConfig"].(map[string]interface{}); ok {
				ikeConfigMap := make(map[string]interface{})
				ikeConfigMap["ike_auth_alg"] = ikeConfigArg["IkeAuthAlg"]
				ikeConfigMap["ike_enc_alg"] = ikeConfigArg["IkeEncAlg"]
				ikeConfigMap["ike_lifetime"] = formatInt(ikeConfigArg["IkeLifetime"])
				ikeConfigMap["ike_mode"] = ikeConfigArg["IkeMode"]
				ikeConfigMap["ike_pfs"] = ikeConfigArg["IkePfs"]
				ikeConfigMap["ike_version"] = ikeConfigArg["IkeVersion"]
				ikeConfigMap["local_id"] = ikeConfigArg["LocalId"]
				ikeConfigMap["remote_id"] = ikeConfigArg["RemoteId"]
				ikeConfigSli = append(ikeConfigSli, ikeConfigMap)
			}
		}
		d.Set("ike_config", ikeConfigSli)
	}

	if v, ok := object["IpsecConfig"]; ok {
		ipsecConfigSli := make([]map[string]interface{}, 0)
		if len(v.(map[string]interface{})) > 0 {
			ipsecConfig := object["IpsecConfig"]
			if ipsecConfigArg, ok := ipsecConfig.(map[string]interface{}); ok {
				ipsecConfigMap := make(map[string]interface{})
				ipsecConfigMap["ipsec_auth_alg"] = ipsecConfigArg["IpsecAuthAlg"]
				ipsecConfigMap["ipsec_enc_alg"] = ipsecConfigArg["IpsecEncAlg"]
				ipsecConfigMap["ipsec_lifetime"] = formatInt(ipsecConfigArg["IpsecLifetime"])
				ipsecConfigMap["ipsec_pfs"] = ipsecConfigArg["IpsecPfs"]
				ipsecConfigSli = append(ipsecConfigSli, ipsecConfigMap)
			}
		}
		d.Set("ipsec_config", ipsecConfigSli)
	}

	return nil
}
func resourceAlicloudVpnIpsecServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"IpsecServerId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("client_ip_pool") {
		update = true
		request["ClientIpPool"] = d.Get("client_ip_pool")
	}
	if d.HasChange("effect_immediately") {
		update = true
		if v, ok := d.GetOkExists("effect_immediately"); ok {
			request["EffectImmediately"] = v
		}
	}
	if d.HasChange("ipsec_server_name") {
		update = true
		if v, ok := d.GetOk("ipsec_server_name"); ok {
			request["IpsecServerName"] = v
		}
	}
	if d.HasChange("local_subnet") {
		update = true
		request["LocalSubnet"] = d.Get("local_subnet")
	}
	if d.HasChange("psk") {
		update = true
		if v, ok := d.GetOk("psk"); ok {
			request["Psk"] = v
		}
	}
	if d.HasChange("ike_config") {
		update = true
		if v, ok := d.GetOk("ike_config"); ok {
			ikeConfig := make(map[string]interface{})
			for _, ikeConfigArgs := range v.(*schema.Set).List() {
				ikeConfigArg := ikeConfigArgs.(map[string]interface{})
				ikeConfig["IkeVersion"] = ikeConfigArg["ike_version"]
				ikeConfig["IkeMode"] = ikeConfigArg["ike_mode"]
				ikeConfig["IkeEncAlg"] = ikeConfigArg["ike_enc_alg"]
				ikeConfig["IkeAuthAlg"] = ikeConfigArg["ike_auth_alg"]
				ikeConfig["IkePfs"] = ikeConfigArg["ike_pfs"]
				ikeConfig["IkeLifetime"] = ikeConfigArg["ike_lifetime"]
				ikeConfig["LocalId"] = ikeConfigArg["local_id"]
				ikeConfig["RemoteId"] = ikeConfigArg["remote_id"]
			}

			if v, err := convertMaptoJsonString(ikeConfig); err != nil {
				return WrapError(err)
			} else {
				request["IkeConfig"] = v
			}
		}
	}

	if d.HasChange("ipsec_config") {
		update = true
		if v, ok := d.GetOk("ipsec_config"); ok {
			ipsecConfig := make(map[string]interface{})
			for _, ipsecConfigArgs := range v.(*schema.Set).List() {
				ipsecConfigArg := ipsecConfigArgs.(map[string]interface{})
				ipsecConfig["IpsecEncAlg"] = ipsecConfigArg["ipsec_enc_alg"]
				ipsecConfig["IpsecAuthAlg"] = ipsecConfigArg["ipsec_auth_alg"]
				ipsecConfig["IpsecPfs"] = ipsecConfigArg["ipsec_pfs"]
				ipsecConfig["IpsecLifetime"] = ipsecConfigArg["ipsec_lifetime"]
			}
			if v, err := convertMaptoJsonString(ipsecConfig); err != nil {
				return WrapError(err)
			} else {
				request["IpsecConfig"] = v
			}
		}
	}

	if d.HasChange("psk_enabled") {
		update = true
		if v, ok := d.GetOkExists("psk_enabled"); ok {
			request["PskEnabled"] = v
		}
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateIpsecServer"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateIpsecServer")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
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
	return resourceAlicloudVpnIpsecServerRead(d, meta)
}
func resourceAlicloudVpnIpsecServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteIpsecServer"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"IpsecServerId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteIpsecServer")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
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
