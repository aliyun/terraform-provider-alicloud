package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaSiteCreate,
		Read:   resourceAliCloudEsaSiteRead,
		Update: resourceAliCloudEsaSiteUpdate,
		Delete: resourceAliCloudEsaSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"add_client_geolocation_header": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"add_real_client_ip_header": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ai_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ai_template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"automatic_frequency_control_action_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"automatic_frequency_control_enable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"automatic_frequency_control_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cache_architecture_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cache_reserve_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cache_reserve_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"case_insensitive": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"coverage": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cross_border_optimization": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"development_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flatten_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"global_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"paused": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"performance_data_collection_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"real_client_ip_header_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"seo_bypass": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"site_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"site_name_exclusive": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"site_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"site_waf_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"add_bot_protection_headers": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"client_ip_identifier": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"headers": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"add_security_headers": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"security_level": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"bandwidth_abuse_protection": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"bot_management": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definite_bots": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"likely_bots": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"js_detection": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"verified_bots": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"effect_on_static": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"version_management": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEsaSiteCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSite"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Coverage"] = d.Get("coverage")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	request["SiteName"] = d.Get("site_name")
	request["AccessType"] = d.Get("access_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_site", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SiteId"]))

	esaServiceV2 := EsaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"pending"}, d.Timeout(schema.TimeoutCreate), 3*time.Minute, esaServiceV2.EsaSiteStateRefreshFuncWithApi(d.Id(), "Status", []string{}, esaServiceV2.DescribeEsaSite))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEsaSiteUpdate(d, meta)
}

func resourceAliCloudEsaSiteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaSite(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_site DescribeEsaSite Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_type", objectRaw["AccessType"])
	d.Set("coverage", objectRaw["Coverage"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("site_name", objectRaw["SiteName"])
	d.Set("status", objectRaw["Status"])
	d.Set("version_management", objectRaw["VersionManagement"])

	objectRaw, err = esaServiceV2.DescribeSiteListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = esaServiceV2.DescribeSiteGetManagedTransform(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("add_client_geolocation_header", objectRaw["AddClientGeolocationHeader"])
	d.Set("add_real_client_ip_header", objectRaw["AddRealClientIpHeader"])
	d.Set("real_client_ip_header_name", objectRaw["RealClientIpHeaderName"])
	d.Set("site_version", objectRaw["SiteVersion"])

	objectRaw, err = esaServiceV2.DescribeSiteGetCacheTag(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("case_insensitive", objectRaw["CaseInsensitive"])
	d.Set("site_version", objectRaw["SiteVersion"])
	d.Set("tag_name", objectRaw["TagName"])

	objectRaw, err = esaServiceV2.DescribeSiteGetIPv6(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("ipv6_enable", objectRaw["Enable"])
	d.Set("ipv6_region", objectRaw["Region"])

	objectRaw, err = esaServiceV2.DescribeSiteGetCacheReserve(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("cache_reserve_enable", objectRaw["Enable"])
	d.Set("cache_reserve_instance_id", objectRaw["CacheReserveInstanceId"])

	objectRaw, err = esaServiceV2.DescribeSiteDescribeHttpDDoSAttackIntelligentProtection(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("ai_mode", objectRaw["AiMode"])
	d.Set("ai_template", objectRaw["AiTemplate"])

	objectRaw, err = esaServiceV2.DescribeSiteGetTieredCache(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("cache_architecture_mode", objectRaw["CacheArchitectureMode"])

	objectRaw, err = esaServiceV2.DescribeSiteGetCrossBorderOptimization(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("cross_border_optimization", objectRaw["Enable"])

	objectRaw, err = esaServiceV2.DescribeSiteGetSiteNameExclusive(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("site_name_exclusive", objectRaw["Enable"])

	objectRaw, err = esaServiceV2.DescribeSiteGetCnameFlattening(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("flatten_mode", objectRaw["FlattenMode"])

	objectRaw, err = esaServiceV2.DescribeSiteGetSeoBypass(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("seo_bypass", objectRaw["Enable"])

	objectRaw, err = esaServiceV2.DescribeSiteGetDevelopmentMode(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("development_mode", objectRaw["Enable"])

	objectRaw, err = esaServiceV2.DescribeSiteGetSitePause(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("paused", objectRaw["Paused"])

	objectRaw, err = esaServiceV2.DescribeSiteDescribeHttpDDoSAttackProtection(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("global_mode", objectRaw["GlobalMode"])

	objectRaw, err = esaServiceV2.DescribeSiteGetAutomaticFrequencyControlConfig(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("automatic_frequency_control_action_type", objectRaw["ActionType"])
	d.Set("automatic_frequency_control_enable", objectRaw["Enable"])
	d.Set("automatic_frequency_control_level", objectRaw["Level"])

	objectRaw, err = esaServiceV2.DescribeSiteGetPerformanceDataCollection(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("performance_data_collection_enable", objectRaw["Enable"])

	objectRaw, err = esaServiceV2.DescribeSiteGetSiteWafSettings(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	siteWafSettingsMaps := make([]map[string]interface{}, 0)
	siteWafSettingsMap := make(map[string]interface{})

	addBotProtectionHeadersMaps := make([]map[string]interface{}, 0)
	addBotProtectionHeadersMap := make(map[string]interface{})
	addBotProtectionHeadersRawObj, _ := jsonpath.Get("$.AddBotProtectionHeaders", objectRaw)
	addBotProtectionHeadersRaw := make(map[string]interface{})
	if addBotProtectionHeadersRawObj != nil {
		addBotProtectionHeadersRaw = addBotProtectionHeadersRawObj.(map[string]interface{})
	}
	if len(addBotProtectionHeadersRaw) > 0 {
		addBotProtectionHeadersMap["enable"] = addBotProtectionHeadersRaw["Enable"]

		addBotProtectionHeadersMaps = append(addBotProtectionHeadersMaps, addBotProtectionHeadersMap)
	}
	siteWafSettingsMap["add_bot_protection_headers"] = addBotProtectionHeadersMaps
	addSecurityHeadersMaps := make([]map[string]interface{}, 0)
	addSecurityHeadersMap := make(map[string]interface{})
	addSecurityHeadersRawObj, _ := jsonpath.Get("$.AddSecurityHeaders", objectRaw)
	addSecurityHeadersRaw := make(map[string]interface{})
	if addSecurityHeadersRawObj != nil {
		addSecurityHeadersRaw = addSecurityHeadersRawObj.(map[string]interface{})
	}
	if len(addSecurityHeadersRaw) > 0 {
		addSecurityHeadersMap["enable"] = addSecurityHeadersRaw["Enable"]

		addSecurityHeadersMaps = append(addSecurityHeadersMaps, addSecurityHeadersMap)
	}
	siteWafSettingsMap["add_security_headers"] = addSecurityHeadersMaps
	bandwidthAbuseProtectionMaps := make([]map[string]interface{}, 0)
	bandwidthAbuseProtectionMap := make(map[string]interface{})
	bandwidthAbuseProtectionRawObj, _ := jsonpath.Get("$.BandwidthAbuseProtection", objectRaw)
	bandwidthAbuseProtectionRaw := make(map[string]interface{})
	if bandwidthAbuseProtectionRawObj != nil {
		bandwidthAbuseProtectionRaw = bandwidthAbuseProtectionRawObj.(map[string]interface{})
	}
	if len(bandwidthAbuseProtectionRaw) > 0 {
		bandwidthAbuseProtectionMap["action"] = bandwidthAbuseProtectionRaw["Action"]
		bandwidthAbuseProtectionMap["id"] = bandwidthAbuseProtectionRaw["Id"]
		bandwidthAbuseProtectionMap["status"] = bandwidthAbuseProtectionRaw["Status"]

		bandwidthAbuseProtectionMaps = append(bandwidthAbuseProtectionMaps, bandwidthAbuseProtectionMap)
	}
	siteWafSettingsMap["bandwidth_abuse_protection"] = bandwidthAbuseProtectionMaps
	botManagementMaps := make([]map[string]interface{}, 0)
	botManagementMap := make(map[string]interface{})
	botManagementRawObj, _ := jsonpath.Get("$.BotManagement", objectRaw)
	botManagementRaw := make(map[string]interface{})
	if botManagementRawObj != nil {
		botManagementRaw = botManagementRawObj.(map[string]interface{})
	}
	if len(botManagementRaw) > 0 {

		definiteBotsMaps := make([]map[string]interface{}, 0)
		definiteBotsMap := make(map[string]interface{})
		definiteBotsRawObj, _ := jsonpath.Get("$.BotManagement.DefiniteBots", objectRaw)
		definiteBotsRaw := make(map[string]interface{})
		if definiteBotsRawObj != nil {
			definiteBotsRaw = definiteBotsRawObj.(map[string]interface{})
		}
		if len(definiteBotsRaw) > 0 {
			definiteBotsMap["action"] = definiteBotsRaw["Action"]
			definiteBotsMap["id"] = definiteBotsRaw["Id"]

			definiteBotsMaps = append(definiteBotsMaps, definiteBotsMap)
		}
		botManagementMap["definite_bots"] = definiteBotsMaps
		effectOnStaticMaps := make([]map[string]interface{}, 0)
		effectOnStaticMap := make(map[string]interface{})
		effectOnStaticRawObj, _ := jsonpath.Get("$.BotManagement.EffectOnStatic", objectRaw)
		effectOnStaticRaw := make(map[string]interface{})
		if effectOnStaticRawObj != nil {
			effectOnStaticRaw = effectOnStaticRawObj.(map[string]interface{})
		}
		if len(effectOnStaticRaw) > 0 {
			effectOnStaticMap["enable"] = effectOnStaticRaw["Enable"]

			effectOnStaticMaps = append(effectOnStaticMaps, effectOnStaticMap)
		}
		botManagementMap["effect_on_static"] = effectOnStaticMaps
		jSDetectionMaps := make([]map[string]interface{}, 0)
		jSDetectionMap := make(map[string]interface{})
		jSDetectionRawObj, _ := jsonpath.Get("$.BotManagement.JSDetection", objectRaw)
		jSDetectionRaw := make(map[string]interface{})
		if jSDetectionRawObj != nil {
			jSDetectionRaw = jSDetectionRawObj.(map[string]interface{})
		}
		if len(jSDetectionRaw) > 0 {
			jSDetectionMap["enable"] = jSDetectionRaw["Enable"]

			jSDetectionMaps = append(jSDetectionMaps, jSDetectionMap)
		}
		botManagementMap["js_detection"] = jSDetectionMaps
		likelyBotsMaps := make([]map[string]interface{}, 0)
		likelyBotsMap := make(map[string]interface{})
		likelyBotsRawObj, _ := jsonpath.Get("$.BotManagement.LikelyBots", objectRaw)
		likelyBotsRaw := make(map[string]interface{})
		if likelyBotsRawObj != nil {
			likelyBotsRaw = likelyBotsRawObj.(map[string]interface{})
		}
		if len(likelyBotsRaw) > 0 {
			likelyBotsMap["action"] = likelyBotsRaw["Action"]
			likelyBotsMap["id"] = likelyBotsRaw["Id"]

			likelyBotsMaps = append(likelyBotsMaps, likelyBotsMap)
		}
		botManagementMap["likely_bots"] = likelyBotsMaps
		verifiedBotsMaps := make([]map[string]interface{}, 0)
		verifiedBotsMap := make(map[string]interface{})
		verifiedBotsRawObj, _ := jsonpath.Get("$.BotManagement.VerifiedBots", objectRaw)
		verifiedBotsRaw := make(map[string]interface{})
		if verifiedBotsRawObj != nil {
			verifiedBotsRaw = verifiedBotsRawObj.(map[string]interface{})
		}
		if len(verifiedBotsRaw) > 0 {
			verifiedBotsMap["action"] = verifiedBotsRaw["Action"]
			verifiedBotsMap["id"] = verifiedBotsRaw["Id"]

			verifiedBotsMaps = append(verifiedBotsMaps, verifiedBotsMap)
		}
		botManagementMap["verified_bots"] = verifiedBotsMaps
		botManagementMaps = append(botManagementMaps, botManagementMap)
	}
	siteWafSettingsMap["bot_management"] = botManagementMaps
	clientIpIdentifierMaps := make([]map[string]interface{}, 0)
	clientIpIdentifierMap := make(map[string]interface{})
	clientIpIdentifierRawObj, _ := jsonpath.Get("$.ClientIpIdentifier", objectRaw)
	clientIpIdentifierRaw := make(map[string]interface{})
	if clientIpIdentifierRawObj != nil {
		clientIpIdentifierRaw = clientIpIdentifierRawObj.(map[string]interface{})
	}
	if len(clientIpIdentifierRaw) > 0 {
		clientIpIdentifierMap["mode"] = clientIpIdentifierRaw["Mode"]

		headersRaw, _ := jsonpath.Get("$.ClientIpIdentifier.Headers", objectRaw)
		clientIpIdentifierMap["headers"] = headersRaw
		clientIpIdentifierMaps = append(clientIpIdentifierMaps, clientIpIdentifierMap)
	}
	siteWafSettingsMap["client_ip_identifier"] = clientIpIdentifierMaps
	securityLevelMaps := make([]map[string]interface{}, 0)
	securityLevelMap := make(map[string]interface{})
	securityLevelRawObj, _ := jsonpath.Get("$.SecurityLevel", objectRaw)
	securityLevelRaw := make(map[string]interface{})
	if securityLevelRawObj != nil {
		securityLevelRaw = securityLevelRawObj.(map[string]interface{})
	}
	if len(securityLevelRaw) > 0 {
		securityLevelMap["value"] = securityLevelRaw["Value"]

		securityLevelMaps = append(securityLevelMaps, securityLevelMap)
	}
	siteWafSettingsMap["security_level"] = securityLevelMaps
	siteWafSettingsMaps = append(siteWafSettingsMaps, siteWafSettingsMap)
	if err := d.Set("site_waf_settings", siteWafSettingsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudEsaSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	esaServiceV2 := EsaServiceV2{client}
	objectRaw, _ := esaServiceV2.DescribeEsaSite(d.Id())

	if d.HasChange("version_management") {
		var err error
		target := d.Get("version_management").(bool)

		currentStatus, err := jsonpath.Get("VersionManagement", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "VersionManagement", objectRaw)
		}
		if formatBool(currentStatus) != target {
			if target == true {
				action := "ActivateVersionManagement"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["SiteId"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
			if target == false {
				action := "DeactivateVersionManagement"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["SiteId"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
		}
	}

	var err error
	action := "UpdateSiteCoverage"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("coverage") {
		update = true
	}
	request["Coverage"] = d.Get("coverage")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy"}) || NeedRetry(err) {
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
		esaServiceV2 := EsaServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"pending"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, esaServiceV2.EsaSiteStateRefreshFuncWithApi(d.Id(), "Status", []string{}, esaServiceV2.DescribeEsaSite))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateIPv6"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("ipv6_region") {
		update = true
		request["Region"] = d.Get("ipv6_region")
	}

	if d.HasChange("ipv6_enable") {
		update = true
	}
	request["Enable"] = d.Get("ipv6_enable")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "UpdateTieredCache"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("cache_architecture_mode") {
		update = true
	}
	request["CacheArchitectureMode"] = d.Get("cache_architecture_mode")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "UpdateManagedTransform"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("add_real_client_ip_header") {
		update = true
	}
	if v, ok := d.GetOk("add_real_client_ip_header"); ok {
		request["AddRealClientIpHeader"] = v
	}

	if d.HasChange("real_client_ip_header_name") {
		update = true
	}
	if v, ok := d.GetOk("real_client_ip_header_name"); ok {
		request["RealClientIpHeaderName"] = v
	}

	if d.HasChange("site_version") {
		update = true
		request["SiteVersion"] = d.Get("site_version")
	}

	if d.HasChange("add_client_geolocation_header") {
		update = true
	}
	if v, ok := d.GetOk("add_client_geolocation_header"); ok {
		request["AddClientGeolocationHeader"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "UpdateCrossBorderOptimization"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("cross_border_optimization") {
		update = true
	}
	request["Enable"] = d.Get("cross_border_optimization")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "UpdateSiteNameExclusive"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("site_name_exclusive") {
		update = true
	}
	request["Enable"] = d.Get("site_name_exclusive")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "UpdateCnameFlattening"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("flatten_mode") {
		update = true
	}
	request["FlattenMode"] = d.Get("flatten_mode")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "UpdateSeoBypass"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("seo_bypass") {
		update = true
	}
	request["Enable"] = d.Get("seo_bypass")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "UpdateCacheTag"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("case_insensitive") {
		update = true
		request["CaseInsensitive"] = d.Get("case_insensitive")
	}

	if d.HasChange("site_version") {
		update = true
		request["SiteVersion"] = d.Get("site_version")
	}

	if d.HasChange("tag_name") {
		update = true
		request["TagName"] = d.Get("tag_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "UpdateDevelopmentMode"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("development_mode") {
		update = true
	}
	request["Enable"] = d.Get("development_mode")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "UpdateSitePause"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("paused") {
		update = true
	}
	request["Paused"] = d.Get("paused")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "UpdateCacheReserve"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("cache_reserve_instance_id") {
		update = true
		request["CacheReserveInstanceId"] = d.Get("cache_reserve_instance_id")
	}

	if d.HasChange("cache_reserve_enable") {
		update = true
		request["Enable"] = d.Get("cache_reserve_enable")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "SetHttpDDoSAttackProtection"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("global_mode") {
		update = true
	}
	request["GlobalMode"] = d.Get("global_mode")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "SetHttpDDoSAttackIntelligentProtection"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("ai_template") {
		update = true
	}
	request["AiTemplate"] = d.Get("ai_template")
	if d.HasChange("ai_mode") {
		update = true
	}
	request["AiMode"] = d.Get("ai_mode")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	update = false
	action = "EditSiteWafSettings"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("site_version") {
		update = true
		request["SiteVersion"] = d.Get("site_version")
	}

	if d.HasChange("site_waf_settings") {
		update = true
		settings := make(map[string]interface{})

		if v := d.Get("site_waf_settings"); v != nil {
			clientIpIdentifier := make(map[string]interface{})
			headers1, _ := jsonpath.Get("$[0].client_ip_identifier[0].headers", d.Get("site_waf_settings"))
			if headers1 != nil && headers1 != "" {
				clientIpIdentifier["Headers"] = headers1
			}
			mode1, _ := jsonpath.Get("$[0].client_ip_identifier[0].mode", d.Get("site_waf_settings"))
			if mode1 != nil && mode1 != "" {
				clientIpIdentifier["Mode"] = mode1
			}

			if len(clientIpIdentifier) > 0 {
				settings["ClientIpIdentifier"] = clientIpIdentifier
			}
			securityLevel := make(map[string]interface{})
			value1, _ := jsonpath.Get("$[0].security_level[0].value", d.Get("site_waf_settings"))
			if value1 != nil && value1 != "" {
				securityLevel["Value"] = value1
			}

			if len(securityLevel) > 0 {
				settings["SecurityLevel"] = securityLevel
			}
			bandwidthAbuseProtection := make(map[string]interface{})
			id1, _ := jsonpath.Get("$[0].bandwidth_abuse_protection[0].id", d.Get("site_waf_settings"))
			if id1 != nil && id1 != "" {
				bandwidthAbuseProtection["Id"] = id1
			}
			status1, _ := jsonpath.Get("$[0].bandwidth_abuse_protection[0].status", d.Get("site_waf_settings"))
			if status1 != nil && status1 != "" {
				bandwidthAbuseProtection["Status"] = status1
			}
			action1, _ := jsonpath.Get("$[0].bandwidth_abuse_protection[0].action", d.Get("site_waf_settings"))
			if action1 != nil && action1 != "" {
				bandwidthAbuseProtection["Action"] = action1
			}

			if len(bandwidthAbuseProtection) > 0 {
				settings["BandwidthAbuseProtection"] = bandwidthAbuseProtection
			}
			addSecurityHeaders := make(map[string]interface{})
			enable1, _ := jsonpath.Get("$[0].add_security_headers[0].enable", d.Get("site_waf_settings"))
			if enable1 != nil && enable1 != "" {
				addSecurityHeaders["Enable"] = enable1
			}

			if len(addSecurityHeaders) > 0 {
				settings["AddSecurityHeaders"] = addSecurityHeaders
			}
			addBotProtectionHeaders := make(map[string]interface{})
			enable3, _ := jsonpath.Get("$[0].add_bot_protection_headers[0].enable", d.Get("site_waf_settings"))
			if enable3 != nil && enable3 != "" {
				addBotProtectionHeaders["Enable"] = enable3
			}

			if len(addBotProtectionHeaders) > 0 {
				settings["AddBotProtectionHeaders"] = addBotProtectionHeaders
			}
			botManagement := make(map[string]interface{})
			likelyBots := make(map[string]interface{})
			action3, _ := jsonpath.Get("$[0].bot_management[0].likely_bots[0].action", d.Get("site_waf_settings"))
			if action3 != nil && action3 != "" {
				likelyBots["Action"] = action3
			}
			id3, _ := jsonpath.Get("$[0].bot_management[0].likely_bots[0].id", d.Get("site_waf_settings"))
			if id3 != nil && id3 != "" {
				likelyBots["Id"] = id3
			}

			if len(likelyBots) > 0 {
				botManagement["LikelyBots"] = likelyBots
			}
			effectOnStatic := make(map[string]interface{})
			enable5, _ := jsonpath.Get("$[0].bot_management[0].effect_on_static[0].enable", d.Get("site_waf_settings"))
			if enable5 != nil && enable5 != "" {
				effectOnStatic["Enable"] = enable5
			}

			if len(effectOnStatic) > 0 {
				botManagement["EffectOnStatic"] = effectOnStatic
			}
			definiteBots := make(map[string]interface{})
			id5, _ := jsonpath.Get("$[0].bot_management[0].definite_bots[0].id", d.Get("site_waf_settings"))
			if id5 != nil && id5 != "" {
				definiteBots["Id"] = id5
			}
			action5, _ := jsonpath.Get("$[0].bot_management[0].definite_bots[0].action", d.Get("site_waf_settings"))
			if action5 != nil && action5 != "" {
				definiteBots["Action"] = action5
			}

			if len(definiteBots) > 0 {
				botManagement["DefiniteBots"] = definiteBots
			}
			verifiedBots := make(map[string]interface{})
			action7, _ := jsonpath.Get("$[0].bot_management[0].verified_bots[0].action", d.Get("site_waf_settings"))
			if action7 != nil && action7 != "" {
				verifiedBots["Action"] = action7
			}
			id7, _ := jsonpath.Get("$[0].bot_management[0].verified_bots[0].id", d.Get("site_waf_settings"))
			if id7 != nil && id7 != "" {
				verifiedBots["Id"] = id7
			}

			if len(verifiedBots) > 0 {
				botManagement["VerifiedBots"] = verifiedBots
			}
			jSDetection := make(map[string]interface{})
			enable7, _ := jsonpath.Get("$[0].bot_management[0].js_detection[0].enable", d.Get("site_waf_settings"))
			if enable7 != nil && enable7 != "" {
				jSDetection["Enable"] = enable7
			}

			if len(jSDetection) > 0 {
				botManagement["JSDetection"] = jSDetection
			}

			if len(botManagement) > 0 {
				settings["BotManagement"] = botManagement
			}

			settingsJson, err := json.Marshal(settings)
			if err != nil {
				return WrapError(err)
			}
			request["Settings"] = string(settingsJson)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
	update = false
	action = "SetAutomaticFrequencyControlConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("site_version") {
		update = true
		request["SiteVersion"] = d.Get("site_version")
	}

	if d.HasChange("automatic_frequency_control_action_type") {
		update = true
	}
	request["ActionType"] = d.Get("automatic_frequency_control_action_type")
	if d.HasChange("automatic_frequency_control_level") {
		update = true
	}
	request["Level"] = d.Get("automatic_frequency_control_level")
	if d.HasChange("automatic_frequency_control_enable") {
		update = true
	}
	request["Enable"] = d.Get("automatic_frequency_control_enable")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
	update = false
	action = "UpdatePerformanceDataCollection"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()

	if d.HasChange("performance_data_collection_enable") {
		update = true
	}
	request["Enable"] = d.Get("performance_data_collection_enable")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
	if d.HasChange("tags") {
		esaServiceV2 := EsaServiceV2{client}
		if err := esaServiceV2.SetResourceTags(d, "Site"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudEsaSiteRead(d, meta)
}

func resourceAliCloudEsaSiteDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSite"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SiteId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	esaServiceV2 := EsaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 10*time.Second, esaServiceV2.DescribeAsyncEsaSiteStateRefreshFunc(d, response, "$.SiteModel.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
