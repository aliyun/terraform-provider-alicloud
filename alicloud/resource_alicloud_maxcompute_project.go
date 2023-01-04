package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMaxcomputeProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMaxcomputeProjectCreate,
		Read:   resourceAlicloudMaxcomputeProjectRead,
		Update: resourceAlicloudMaxcomputeProjectUpdate,
		Delete: resourceAlicloudMaxcomputeProjectDelete,
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
				Removed:      "Field 'order_type' has been removed from provider version 1.196.0.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(3, 27),
				Removed:      "Field 'name' has been removed from provider version 1.196.0.",
			},
			"specification_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"OdpsStandard"}, false),
				Removed:      "Field 'specification_type' has been removed from provider version 1.196.0.",
			},
			"comment": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"default_quota": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"ip_white_list": {
				Optional: true,
				MaxItems: 1,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_list": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"vpc_ip_list": {
							Optional: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"owner": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"product_type": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription", "Dev"}, false),
			},
			"project_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"properties": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeSet,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_full_scan": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeBool,
						},
						"enable_decimal2": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeBool,
						},
						"encryption": {
							Optional: true,
							ForceNew: true,
							Computed: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"algorithm": {
										Optional: true,
										ForceNew: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"enable": {
										Optional: true,
										ForceNew: true,
										Computed: true,
										Type:     schema.TypeBool,
									},
									"key": {
										Optional: true,
										ForceNew: true,
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"retention_days": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeInt,
						},
						"sql_metering_max": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeString,
						},
						"table_lifecycle": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"value": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"timezone": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeString,
						},
						"type_system": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"security_properties": {
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_download_privilege": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeBool,
						},
						"label_security": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeBool,
						},
						"object_creator_has_access_permission": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeBool,
						},
						"object_creator_has_grant_permission": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeBool,
						},
						"project_protection": {
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Type:     schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exception_policy": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"protected": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeBool,
									},
								},
							},
						},
						"using_acl": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeBool,
						},
						"using_policy": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeBool,
						},
					},
				},
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"type": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudMaxcomputeProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("project_name"); ok {
		request["name"] = v
	}
	if v, ok := d.GetOk("comment"); ok {
		request["comment"] = v
	}
	if v, ok := d.GetOk("product_type"); ok {
		request["productType"] = convertMaxComputeProjectProductTypeRequest(v)
	}
	if v, ok := d.GetOk("default_quota"); ok {
		request["defaultQuota"] = v
	}
	if v, ok := d.GetOk("properties"); ok {
		properties := make(map[string]interface{})
		propertiesOrigin := v.(*schema.Set).List()
		propertiesRaw := make(map[string]interface{})
		if len(propertiesOrigin) > 0 {
			propertiesRaw = propertiesOrigin[0].(map[string]interface{})
			if sqlMeteringMax, ok := propertiesRaw["sql_metering_max"]; ok {
				properties["sqlMeteringMax"] = sqlMeteringMax
			}
			if typeSystem, ok := propertiesRaw["type_system"]; ok {
				properties["typeSystem"] = typeSystem
			}
			encryptionOrigin := propertiesRaw["encryption"].(*schema.Set).List()
			if len(encryptionOrigin) > 0 {
				encryptionRaw, encryptionRawOk := encryptionOrigin[0].(map[string]interface{})
				if encryptionRawOk {
					encryption := make(map[string]interface{})
					if enable, ok := encryptionRaw["enable"]; ok {
						encryption["enable"] = enable
					}
					if key, ok := encryptionRaw["key"]; ok {
						encryption["key"] = key
					}
					if algorithm, ok := encryptionRaw["algorithm"]; ok {
						encryption["algorithm"] = algorithm
					}
					properties["encryption"] = encryption
				}
			}
			request["properties"] = properties
		}
	}

	var response map[string]interface{}
	action := "/api/v1/projects"
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), nil, nil, request, &util.RuntimeOptions{})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_maxcompute_project", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["name"]))

	return resourceAlicloudMaxcomputeProjectUpdate(d, meta)
}

func resourceAlicloudMaxcomputeProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxComputeService := MaxComputeService{client}

	object, err := maxComputeService.DescribeMaxcomputeProject(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_maxcompute_project maxComputeService.DescribeMaxcomputeProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("project_name", object["name"])
	d.Set("owner", object["owner"])
	d.Set("type", object["type"])
	d.Set("comment", object["comment"])
	d.Set("default_quota", object["defaultQuota"])
	d.Set("status", object["status"])
	ipWhiteList, err := jsonpath.Get("$.ipWhiteList", object)
	if err == nil {
		ipWhiteList60Maps := make([]map[string]interface{}, 0)
		ipWhiteList60Map := make(map[string]interface{})
		ipWhiteList60Raw := ipWhiteList.(map[string]interface{})
		ipWhiteList60Map["ip_list"] = ipWhiteList60Raw["ipList"]
		ipWhiteList60Map["vpc_ip_list"] = ipWhiteList60Raw["vpcIpList"]
		ipWhiteList60Maps = append(ipWhiteList60Maps, ipWhiteList60Map)
		d.Set("ip_white_list", ipWhiteList60Maps)
	}
	securityPropertiesMaps := make([]map[string]interface{}, 0)
	securityPropertiesMap := make(map[string]interface{})
	securityPropertiesRaw := object["securityProperties"].(map[string]interface{})
	securityPropertiesMap["enable_download_privilege"] = securityPropertiesRaw["enableDownloadPrivilege"]
	securityPropertiesMap["using_acl"] = securityPropertiesRaw["usingAcl"]
	securityPropertiesMap["using_policy"] = securityPropertiesRaw["usingPolicy"]
	securityPropertiesMap["label_security"] = securityPropertiesRaw["labelSecurity"]
	securityPropertiesMap["object_creator_has_access_permission"] = securityPropertiesRaw["objectCreatorHasAccessPermission"]
	securityPropertiesMap["object_creator_has_grant_permission"] = securityPropertiesRaw["objectCreatorHasGrantPermission"]
	projectProtectionMaps := make([]map[string]interface{}, 0)
	projectProtectionMap := make(map[string]interface{})
	projectProtectionRaw := securityPropertiesRaw["projectProtection"].(map[string]interface{})
	projectProtectionMap["protected"] = projectProtectionRaw["protected"]
	exceptionPolicy, err := jsonpath.Get("$.exceptionPolicy", projectProtectionRaw)
	if err == nil {
		projectProtectionMap["exceptionPolicy"] = exceptionPolicy
	}
	projectProtectionMaps = append(projectProtectionMaps, projectProtectionMap)
	securityPropertiesMap["project_protection"] = projectProtectionMaps
	securityPropertiesMaps = append(securityPropertiesMaps, securityPropertiesMap)
	err = d.Set("security_properties", securityPropertiesMaps)
	if err != nil {
		return WrapError(err)
	}
	propertiesMaps := make([]map[string]interface{}, 0)
	propertiesMap := make(map[string]interface{})
	propertiesRaw := object["properties"].(map[string]interface{})
	propertiesMap["timezone"] = propertiesRaw["timezone"]
	propertiesMap["allow_full_scan"] = propertiesRaw["allowFullScan"]
	propertiesMap["enable_decimal2"] = propertiesRaw["enableDecimal2"]
	propertiesMap["retention_days"] = formatInt(propertiesRaw["retentionDays"])
	propertiesMap["sql_metering_max"] = propertiesRaw["sqlMeteringMax"]
	propertiesMap["type_system"] = propertiesRaw["typeSystem"]
	tableLifecycleMaps := make([]map[string]interface{}, 0)
	tableLifecycleMap := make(map[string]interface{})
	tableLifecycleRaw := propertiesRaw["tableLifecycle"].(map[string]interface{})
	tableLifecycleMap["type"] = tableLifecycleRaw["type"]
	tableLifecycleMap["value"] = tableLifecycleRaw["value"]
	tableLifecycleMaps = append(tableLifecycleMaps, tableLifecycleMap)
	propertiesMap["table_lifecycle"] = tableLifecycleMaps
	encryptionMaps := make([]map[string]interface{}, 0)
	encryptionMap := make(map[string]interface{})
	encryptionRaw := propertiesRaw["encryption"].(map[string]interface{})
	encryptionMap["enable"] = encryptionRaw["enable"]
	encryptionMap["algorithm"] = encryptionRaw["algorithm"]
	encryptionMap["key"] = encryptionRaw["key"]
	encryptionMaps = append(encryptionMaps, encryptionMap)
	propertiesMap["encryption"] = encryptionMaps
	propertiesMaps = append(propertiesMaps, propertiesMap)
	d.Set("properties", propertiesMaps)

	return nil
}

func resourceAlicloudMaxcomputeProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewOdpsClient()
	update := false
	body := map[string]interface{}{}
	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			body["comment"] = v
			update = true
		}
	}
	if d.HasChange("properties") {
		if v, ok := d.GetOk("properties"); ok {
			update = true
			properties := make(map[string]interface{})
			propertiesOrigin := v.(*schema.Set).List()
			propertiesRaw := make(map[string]interface{})
			if len(propertiesOrigin) > 0 {
				propertiesRaw = propertiesOrigin[0].(map[string]interface{})
				if sqlMeteringMax, ok := propertiesRaw["sql_metering_max"]; ok {
					properties["sqlMeteringMax"] = sqlMeteringMax
				}
				if typeSystem, ok := propertiesRaw["type_system"]; ok {
					properties["typeSystem"] = typeSystem
				}
				if allowFullScan, ok := propertiesRaw["allow_full_scan"]; ok {
					properties["allowFullScan"] = allowFullScan
				}
				if enableDecimal2, ok := propertiesRaw["enable_decimal2"]; ok {
					properties["enableDecimal2"] = enableDecimal2
				}
				if retentionDays, ok := propertiesRaw["retention_days"]; ok {
					properties["retentionDays"] = retentionDays
				}
				if timezone, ok := propertiesRaw["timezone"]; ok {
					properties["timezone"] = timezone
				}
				tableLifecycleOrigin := propertiesRaw["table_lifecycle"].(*schema.Set).List()
				if len(tableLifecycleOrigin) > 0 {
					tableLifecycleRaw, tableLifecycleRawOk := tableLifecycleOrigin[0].(map[string]interface{})
					if tableLifecycleRawOk {
						tableLifecycle := make(map[string]interface{})
						tableLifecycle["type"] = tableLifecycleRaw["type"]
						tableLifecycle["value"] = tableLifecycleRaw["value"]
						properties["tableLifecycle"] = tableLifecycle
					}
				}
				body["properties"] = properties
			}
		}
	}
	if update {
		action := "/api/v1/projects/" + d.Id() + "/meta"
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), nil, nil, body, &util.RuntimeOptions{})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			response := resp
			addDebug(action, response, body)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_maxcompute_project", action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	body = map[string]interface{}{}
	if d.HasChange("default_quota") || d.IsNewResource() {
		if v, ok := d.GetOk("default_quota"); ok {
			body["quota"] = v
			update = true
		}
	}
	if update {
		action := "/api/v1/projects/" + d.Id() + "/quota"
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), nil, nil, body, &util.RuntimeOptions{})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			response := resp
			addDebug(action, response, body)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_maxcompute_project", action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	body = map[string]interface{}{}
	if d.HasChange("ip_white_list") || d.IsNewResource() {
		if v, ok := d.GetOk("ip_white_list"); ok {
			update = true
			ipWhiteListRaw := v.(*schema.Set).List()
			if len(ipWhiteListRaw) > 0 {
				ipWhiteListRawFirst := ipWhiteListRaw[0].(map[string]interface{})
				ipWhiteListRawMap := make(map[string]interface{})
				ipWhiteListRawMap["ipList"] = ipWhiteListRawFirst["ip_list"]
				ipWhiteListRawMap["vpcIpList"] = ipWhiteListRawFirst["vpc_ip_list"]
				body["ipWhiteList"] = ipWhiteListRawMap
			}
		}
	}
	if update {
		action := "/api/v1/projects/" + d.Id() + "/ipWhiteList"
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), nil, nil, body, &util.RuntimeOptions{})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			response := resp
			addDebug(action, response, body)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_maxcompute_project", action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	body = map[string]interface{}{}
	if d.HasChange("security_properties") || d.IsNewResource() {
		if v, ok := d.GetOk("security_properties"); ok {
			update = true
			securityPropertiesRaw := v.(*schema.Set).List()
			if len(securityPropertiesRaw) > 0 {
				securityPropertiesRawFirst := securityPropertiesRaw[0].(map[string]interface{})
				securityPropertiesMap := make(map[string]interface{})
				securityPropertiesMap["usingAcl"] = securityPropertiesRawFirst["using_acl"]
				securityPropertiesMap["usingPolicy"] = securityPropertiesRawFirst["using_policy"]
				securityPropertiesMap["objectCreatorHasAccessPermission"] = securityPropertiesRawFirst["object_creator_has_access_permission"]
				securityPropertiesMap["objectCreatorHasGrantPermission"] = securityPropertiesRawFirst["object_creator_has_grant_permission"]
				securityPropertiesMap["labelSecurity"] = securityPropertiesRawFirst["label_security"]
				securityPropertiesMap["enableDownloadPrivilege"] = securityPropertiesRawFirst["enable_download_privilege"]
				projectProtectionRaw := securityPropertiesRawFirst["project_protection"].(*schema.Set).List()
				if len(projectProtectionRaw) > 0 {
					projectProtection := make(map[string]interface{})
					projectProtectionRawFirst := projectProtectionRaw[0].(map[string]interface{})
					projectProtection["protected"] = projectProtectionRawFirst["protected"]
					projectProtection["exceptionPolicy"] = projectProtectionRawFirst["exception_policy"]
					securityPropertiesMap["projectProtection"] = projectProtection
				}
				body = securityPropertiesMap
			}
		}
	}
	if update {
		action := "/api/v1/projects/" + d.Id() + "/security"
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), nil, nil, body, &util.RuntimeOptions{})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			response := resp
			addDebug(action, response, body)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_maxcompute_project", action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudMaxcomputeProjectRead(d, meta)
}

func resourceAlicloudMaxcomputeProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxComputeService := MaxComputeService{client}
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	query := make(map[string]*string)
	query["isLogical"] = StringPointer("false")

	action := "/api/v1/projects/" + d.Id()
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"OBJECT_NOT_EXIST"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, maxComputeService.MaxcomputeProjectStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertMaxComputeProjectProductTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PAYASYOUGO"
	case "Subscription":
		return "SUBSCRIPTION"
	case "Dev":
		return "DEV"
	}
	return source
}
