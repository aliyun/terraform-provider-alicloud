// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudMaxComputeProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMaxComputeProjectCreate,
		Read:   resourceAliCloudMaxComputeProjectRead,
		Update: resourceAliCloudMaxComputeProjectUpdate,
		Delete: resourceAliCloudMaxComputeProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cost_storage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_quota": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_white_list": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_ip_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"is_logical": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"properties": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sql_metering_max": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tunnel_quota": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"table_lifecycle": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"type_system": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enable_data_masking": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_tunnel_quota_route": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"retention_days": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"encryption": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"algorithm": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"allow_full_scan": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_dr": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_decimal2": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_properties": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"using_policy": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"label_security": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"object_creator_has_grant_permission": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"object_creator_has_access_permission": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"using_acl": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_download_privilege": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"project_protection": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protected": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"exception_policy": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsJSON,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											equal, _ := compareJsonTemplateAreEquivalent(old, new)
											return equal
										},
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"AVAILABLE", "READONLY", "DELETING", "FROZEN"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "FROZEN" && strings.HasPrefix(old, "FROZEN")
				},
			},
			"tags": tagsSchema(),
			"three_tier_model": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"trusted_projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudMaxComputeProjectCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/projects")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("project_name"); ok {
		request["name"] = v
	}

	if v, ok := d.GetOk("comment"); ok {
		request["comment"] = v
	}
	if v, ok := d.GetOk("default_quota"); ok {
		request["defaultQuota"] = v
	}
	properties := make(map[string]interface{})

	if v := d.Get("properties"); !IsNil(v) {
		encryption := make(map[string]interface{})
		algorithm1, _ := jsonpath.Get("$[0].encryption[0].algorithm", d.Get("properties"))
		if algorithm1 != nil && algorithm1 != "" {
			encryption["algorithm"] = algorithm1
		}
		enable1, _ := jsonpath.Get("$[0].encryption[0].enable", d.Get("properties"))
		if enable1 != nil && enable1 != "" {
			encryption["enable"] = enable1
		}
		key1, _ := jsonpath.Get("$[0].encryption[0].key", d.Get("properties"))
		if key1 != nil && key1 != "" {
			encryption["key"] = key1
		}

		if len(encryption) > 0 {
			properties["encryption"] = encryption
		}
		sqlMeteringMax1, _ := jsonpath.Get("$[0].sql_metering_max", v)
		if sqlMeteringMax1 != nil && sqlMeteringMax1 != "" {
			properties["sqlMeteringMax"] = sqlMeteringMax1
		}
		typeSystem1, _ := jsonpath.Get("$[0].type_system", v)
		if typeSystem1 != nil && typeSystem1 != "" {
			properties["typeSystem"] = typeSystem1
		}

		request["properties"] = properties
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("MaxCompute", "2022-01-04", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_maxcompute_project", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["name"]))

	return resourceAliCloudMaxComputeProjectUpdate(d, meta)
}

func resourceAliCloudMaxComputeProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxComputeServiceV2 := MaxComputeServiceV2{client}

	objectRaw, err := maxComputeServiceV2.DescribeMaxComputeProject(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_maxcompute_project DescribeMaxComputeProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("comment", objectRaw["comment"])
	d.Set("cost_storage", objectRaw["costStorage"])
	d.Set("create_time", objectRaw["createdTime"])
	d.Set("default_quota", objectRaw["defaultQuota"])
	d.Set("owner", objectRaw["owner"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("status", objectRaw["status"])
	d.Set("three_tier_model", objectRaw["threeTierModel"])
	d.Set("type", objectRaw["type"])
	d.Set("project_name", objectRaw["name"])

	ipWhiteListMaps := make([]map[string]interface{}, 0)
	ipWhiteListMap := make(map[string]interface{})
	ipWhiteListRaw := make(map[string]interface{})
	if objectRaw["ipWhiteList"] != nil {
		ipWhiteListRaw = objectRaw["ipWhiteList"].(map[string]interface{})
	}
	if len(ipWhiteListRaw) > 0 {
		ipWhiteListMap["ip_list"] = ipWhiteListRaw["ipList"]
		ipWhiteListMap["vpc_ip_list"] = ipWhiteListRaw["vpcIpList"]

		ipWhiteListMaps = append(ipWhiteListMaps, ipWhiteListMap)
	}
	if err := d.Set("ip_white_list", ipWhiteListMaps); err != nil {
		return err
	}
	propertiesMaps := make([]map[string]interface{}, 0)
	propertiesMap := make(map[string]interface{})
	propertiesRaw := make(map[string]interface{})
	if objectRaw["properties"] != nil {
		propertiesRaw = objectRaw["properties"].(map[string]interface{})
	}
	if len(propertiesRaw) > 0 {
		propertiesMap["allow_full_scan"] = propertiesRaw["allowFullScan"]
		propertiesMap["enable_data_masking"] = propertiesRaw["enableDataMasking"]
		propertiesMap["enable_decimal2"] = propertiesRaw["enableDecimal2"]
		propertiesMap["enable_dr"] = propertiesRaw["enableDr"]
		propertiesMap["enable_tunnel_quota_route"] = propertiesRaw["enableTunnelQuotaRoute"]
		propertiesMap["retention_days"] = propertiesRaw["retentionDays"]
		propertiesMap["sql_metering_max"] = propertiesRaw["sqlMeteringMax"]
		propertiesMap["timezone"] = propertiesRaw["timezone"]
		propertiesMap["tunnel_quota"] = propertiesRaw["tunnelQuota"]
		propertiesMap["type_system"] = propertiesRaw["typeSystem"]

		encryptionMaps := make([]map[string]interface{}, 0)
		encryptionMap := make(map[string]interface{})
		encryptionRaw := make(map[string]interface{})
		if propertiesRaw["encryption"] != nil {
			encryptionRaw = propertiesRaw["encryption"].(map[string]interface{})
		}
		if len(encryptionRaw) > 0 {
			encryptionMap["algorithm"] = encryptionRaw["algorithm"]
			encryptionMap["enable"] = encryptionRaw["enable"]
			encryptionMap["key"] = encryptionRaw["key"]

			encryptionMaps = append(encryptionMaps, encryptionMap)
		}
		propertiesMap["encryption"] = encryptionMaps
		tableLifecycleMaps := make([]map[string]interface{}, 0)
		tableLifecycleMap := make(map[string]interface{})
		tableLifecycleRaw := make(map[string]interface{})
		if propertiesRaw["tableLifecycle"] != nil {
			tableLifecycleRaw = propertiesRaw["tableLifecycle"].(map[string]interface{})
		}
		if len(tableLifecycleRaw) > 0 {
			tableLifecycleMap["type"] = tableLifecycleRaw["type"]
			tableLifecycleMap["value"] = tableLifecycleRaw["value"]

			tableLifecycleMaps = append(tableLifecycleMaps, tableLifecycleMap)
		}
		propertiesMap["table_lifecycle"] = tableLifecycleMaps
		propertiesMaps = append(propertiesMaps, propertiesMap)
	}
	if err := d.Set("properties", propertiesMaps); err != nil {
		return err
	}
	securityPropertiesMaps := make([]map[string]interface{}, 0)
	securityPropertiesMap := make(map[string]interface{})
	securityPropertiesRaw := make(map[string]interface{})
	if objectRaw["securityProperties"] != nil {
		securityPropertiesRaw = objectRaw["securityProperties"].(map[string]interface{})
	}
	if len(securityPropertiesRaw) > 0 {
		securityPropertiesMap["enable_download_privilege"] = securityPropertiesRaw["enableDownloadPrivilege"]
		securityPropertiesMap["label_security"] = securityPropertiesRaw["labelSecurity"]
		securityPropertiesMap["object_creator_has_access_permission"] = securityPropertiesRaw["objectCreatorHasAccessPermission"]
		securityPropertiesMap["object_creator_has_grant_permission"] = securityPropertiesRaw["objectCreatorHasGrantPermission"]
		securityPropertiesMap["using_acl"] = securityPropertiesRaw["usingAcl"]
		securityPropertiesMap["using_policy"] = securityPropertiesRaw["usingPolicy"]

		projectProtectionMaps := make([]map[string]interface{}, 0)
		projectProtectionMap := make(map[string]interface{})
		projectProtectionRaw := make(map[string]interface{})
		if securityPropertiesRaw["projectProtection"] != nil {
			projectProtectionRaw = securityPropertiesRaw["projectProtection"].(map[string]interface{})
		}
		if len(projectProtectionRaw) > 0 {
			projectProtectionMap["exception_policy"] = projectProtectionRaw["exceptionPolicy"]
			projectProtectionMap["protected"] = projectProtectionRaw["protected"]

			projectProtectionMaps = append(projectProtectionMaps, projectProtectionMap)
		}
		securityPropertiesMap["project_protection"] = projectProtectionMaps
		securityPropertiesMaps = append(securityPropertiesMaps, securityPropertiesMap)
	}
	if err := d.Set("security_properties", securityPropertiesMaps); err != nil {
		return err
	}

	objectRaw, err = maxComputeServiceV2.DescribeProjectListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagResourcesRaw := make(map[string]interface{})
	if objectRaw["TagResources"] != nil {
		tagResourcesRaw = objectRaw["TagResources"].(map[string]interface{})
	}
	if len(tagResourcesRaw) > 0 {

		tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
		d.Set("tags", tagsToMap(tagsMaps))

		objectRaw, err = maxComputeServiceV2.DescribeProjectGetTrustedProjects(d.Id())
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		dataRaw := make([]interface{}, 0)
		if objectRaw["data"] != nil {
			dataRaw = convertToInterfaceArray(objectRaw["data"])
		}

		d.Set("trusted_projects", dataRaw)
	}

	d.Set("project_name", d.Id())

	return nil
}

func resourceAliCloudMaxComputeProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	maxComputeServiceV2 := MaxComputeServiceV2{client}
	objectRaw, _ := maxComputeServiceV2.DescribeMaxComputeProject(d.Id())

	initedThreeTierModel := false
	if _, ok := d.GetOkExists("three_tier_model"); ok && d.IsNewResource() {
		initedThreeTierModel = true
	}
	if initedThreeTierModel || d.HasChange("three_tier_model") {
		var err error
		target := d.Get("three_tier_model").(bool)

		currentStatus, err := jsonpath.Get("threeTierModel", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "threeTierModel", objectRaw)
		}
		if formatBool(currentStatus) != target {
			if target == true {
				projectName := d.Id()
				action := fmt.Sprintf("/api/v1/projects/%s/modelTier", projectName)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
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
		}
	}

	var err error
	projectName := d.Id()
	action := fmt.Sprintf("/api/v1/projects/%s/quota", projectName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("default_quota") {
		update = true
	}
	if v, ok := d.GetOk("default_quota"); ok || d.HasChange("default_quota") {
		request["quota"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
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
	projectName = d.Id()
	action = fmt.Sprintf("/api/v1/projects/%s/ipWhiteList", projectName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("ip_white_list") {
		update = true
	}
	ipWhiteList := make(map[string]interface{})

	if v := d.Get("ip_white_list"); !IsNil(v) || d.HasChange("ip_white_list") {
		ipList1, _ := jsonpath.Get("$[0].ip_list", v)
		if ipList1 != nil && ipList1 != "" {
			ipWhiteList["ipList"] = ipList1
		}
		vpcIpList1, _ := jsonpath.Get("$[0].vpc_ip_list", v)
		if vpcIpList1 != nil && vpcIpList1 != "" {
			ipWhiteList["vpcIpList"] = vpcIpList1
		}

		request["ipWhiteList"] = ipWhiteList
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
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
	projectName = d.Id()
	action = fmt.Sprintf("/api/v1/projects/%s/security", projectName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("security_properties.0.enable_download_privilege") {
		update = true
	}
	if v, ok := d.GetOkExists("security_properties"); ok || d.HasChange("security_properties") {
		securityPropertiesEnableDownloadPrivilegeJsonPath, err := jsonpath.Get("$[0].enable_download_privilege", v)
		if err == nil && securityPropertiesEnableDownloadPrivilegeJsonPath != "" {
			request["enableDownloadPrivilege"] = securityPropertiesEnableDownloadPrivilegeJsonPath
		}
	}
	if d.HasChange("security_properties.0.label_security") {
		update = true
	}
	if v, ok := d.GetOkExists("security_properties"); ok || d.HasChange("security_properties") {
		securityPropertiesLabelSecurityJsonPath, err := jsonpath.Get("$[0].label_security", v)
		if err == nil && securityPropertiesLabelSecurityJsonPath != "" {
			request["labelSecurity"] = securityPropertiesLabelSecurityJsonPath
		}
	}
	if d.HasChange("security_properties.0.object_creator_has_access_permission") {
		update = true
	}
	if v, ok := d.GetOkExists("security_properties"); ok || d.HasChange("security_properties") {
		securityPropertiesObjectCreatorHasAccessPermissionJsonPath, err := jsonpath.Get("$[0].object_creator_has_access_permission", v)
		if err == nil && securityPropertiesObjectCreatorHasAccessPermissionJsonPath != "" {
			request["objectCreatorHasAccessPermission"] = securityPropertiesObjectCreatorHasAccessPermissionJsonPath
		}
	}
	if d.HasChange("security_properties.0.object_creator_has_grant_permission") {
		update = true
	}
	if v, ok := d.GetOkExists("security_properties"); ok || d.HasChange("security_properties") {
		securityPropertiesObjectCreatorHasGrantPermissionJsonPath, err := jsonpath.Get("$[0].object_creator_has_grant_permission", v)
		if err == nil && securityPropertiesObjectCreatorHasGrantPermissionJsonPath != "" {
			request["objectCreatorHasGrantPermission"] = securityPropertiesObjectCreatorHasGrantPermissionJsonPath
		}
	}
	if d.HasChange("security_properties") {
		update = true
	}
	projectProtection := make(map[string]interface{})

	if v := d.Get("security_properties"); !IsNil(v) || d.HasChange("security_properties") {
		exceptionPolicy1, _ := jsonpath.Get("$[0].project_protection[0].exception_policy", v)
		if exceptionPolicy1 != nil && exceptionPolicy1 != "" {
			projectProtection["exceptionPolicy"] = exceptionPolicy1
		}
		protected1, _ := jsonpath.Get("$[0].project_protection[0].protected", v)
		if protected1 != nil && protected1 != "" {
			projectProtection["protected"] = protected1
		}

		request["projectProtection"] = projectProtection
	}

	if d.HasChange("security_properties.0.using_acl") {
		update = true
	}
	if v, ok := d.GetOkExists("security_properties"); ok || d.HasChange("security_properties") {
		securityPropertiesUsingAclJsonPath, err := jsonpath.Get("$[0].using_acl", v)
		if err == nil && securityPropertiesUsingAclJsonPath != "" {
			request["usingAcl"] = securityPropertiesUsingAclJsonPath
		}
	}
	if d.HasChange("security_properties.0.using_policy") {
		update = true
	}
	if v, ok := d.GetOkExists("security_properties"); ok || d.HasChange("security_properties") {
		securityPropertiesUsingPolicyJsonPath, err := jsonpath.Get("$[0].using_policy", v)
		if err == nil && securityPropertiesUsingPolicyJsonPath != "" {
			request["usingPolicy"] = securityPropertiesUsingPolicyJsonPath
		}
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
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
	projectName = d.Id()
	action = fmt.Sprintf("/api/v1/projects/%s/status", projectName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok || d.HasChange("status") {
		request["status"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
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
	projectName = d.Id()
	action = fmt.Sprintf("/api/v1/projects/%s/meta", projectName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("comment") {
		update = true
	}
	if v, ok := d.GetOk("comment"); ok || d.HasChange("comment") {
		request["comment"] = v
	}
	if d.HasChange("properties") {
		update = true
	}
	properties := make(map[string]interface{})

	if v := d.Get("properties"); !IsNil(v) || d.HasChange("properties") {
		allowFullScan1, _ := jsonpath.Get("$[0].allow_full_scan", v)
		if allowFullScan1 != nil && allowFullScan1 != "" {
			properties["allowFullScan"] = allowFullScan1
		}
		enableDataMasking1, _ := jsonpath.Get("$[0].enable_data_masking", v)
		if enableDataMasking1 != nil && enableDataMasking1 != "" {
			properties["enableDataMasking"] = enableDataMasking1
		}
		enableDecimal21, _ := jsonpath.Get("$[0].enable_decimal2", v)
		if enableDecimal21 != nil && enableDecimal21 != "" {
			properties["enableDecimal2"] = enableDecimal21
		}
		enableDr1, _ := jsonpath.Get("$[0].enable_dr", v)
		if enableDr1 != nil && enableDr1 != "" && d.HasChange("properties.0.enable_dr") {
			properties["enableDr"] = enableDr1
		}
		enableTunnelQuotaRoute1, _ := jsonpath.Get("$[0].enable_tunnel_quota_route", v)
		if enableTunnelQuotaRoute1 != nil && enableTunnelQuotaRoute1 != "" {
			properties["enableTunnelQuotaRoute"] = enableTunnelQuotaRoute1
		}
		encryption := make(map[string]interface{})
		algorithm1, _ := jsonpath.Get("$[0].encryption[0].algorithm", d.Get("properties"))
		if algorithm1 != nil && algorithm1 != "" {
			encryption["algorithm"] = algorithm1
		}
		enable1, _ := jsonpath.Get("$[0].encryption[0].enable", d.Get("properties"))
		if enable1 != nil && enable1 != "" {
			encryption["enable"] = enable1
		}
		key1, _ := jsonpath.Get("$[0].encryption[0].key", d.Get("properties"))
		if key1 != nil && key1 != "" {
			encryption["key"] = key1
		}

		if len(encryption) > 0 && d.HasChange("properties.0.encryption") {
			properties["encryption"] = encryption
		}
		retentionDays1, _ := jsonpath.Get("$[0].retention_days", v)
		if retentionDays1 != nil && retentionDays1 != "" {
			properties["retentionDays"] = retentionDays1
		}
		sqlMeteringMax1, _ := jsonpath.Get("$[0].sql_metering_max", v)
		if sqlMeteringMax1 != nil && sqlMeteringMax1 != "" {
			properties["sqlMeteringMax"] = sqlMeteringMax1
		}
		tableLifecycle := make(map[string]interface{})
		type1, _ := jsonpath.Get("$[0].table_lifecycle[0].type", d.Get("properties"))
		if type1 != nil && type1 != "" {
			tableLifecycle["type"] = type1
		}
		value1, _ := jsonpath.Get("$[0].table_lifecycle[0].value", d.Get("properties"))
		if value1 != nil && value1 != "" {
			tableLifecycle["value"] = value1
		}

		if len(tableLifecycle) > 0 {
			properties["tableLifecycle"] = tableLifecycle
		}
		timezone1, _ := jsonpath.Get("$[0].timezone", v)
		if timezone1 != nil && timezone1 != "" {
			properties["timezone"] = timezone1
		}
		tunnelQuota1, _ := jsonpath.Get("$[0].tunnel_quota", v)
		if tunnelQuota1 != nil && tunnelQuota1 != "" {
			properties["tunnelQuota"] = tunnelQuota1
		}
		typeSystem1, _ := jsonpath.Get("$[0].type_system", v)
		if typeSystem1 != nil && typeSystem1 != "" {
			properties["typeSystem"] = typeSystem1
		}

		request["properties"] = properties
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
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
		maxComputeServiceV2 := MaxComputeServiceV2{client}
		if err := maxComputeServiceV2.SetResourceTags(d, "project"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudMaxComputeProjectRead(d, meta)
}

func resourceAliCloudMaxComputeProjectDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	projectName := d.Id()
	action := fmt.Sprintf("/api/v1/projects/%s", projectName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	query["isLogical"] = StringPointer("false")
	if v, ok := d.GetOk("is_logical"); ok {
		query["isLogical"] = StringPointer(v.(string))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("MaxCompute", "2022-01-04", action, query, nil, nil, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	maxComputeServiceV2 := MaxComputeServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, maxComputeServiceV2.MaxComputeProjectStateRefreshFunc(d.Id(), "status", []string{}))
	stateConf.ContinuousTargetOccurence = 5
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
