// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_quota": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"product_type": {
				Type:     schema.TypeString,
				Optional: true,
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
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sql_metering_max": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type_system": {
							Type:     schema.TypeString,
							Optional: true,
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
						"retention_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"encryption": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"algorithm": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"allow_full_scan": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_decimal2": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
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
										Type:     schema.TypeString,
										Optional: true,
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
			},
			"tags": tagsSchema(),
			"type": {
				Type:     schema.TypeString,
				Computed: true,
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
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["name"] = d.Get("project_name")

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("properties"); !IsNil(v) {
		nodeNative, _ := jsonpath.Get("$[0].sql_metering_max", d.Get("properties"))
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["sqlMeteringMax"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].type_system", d.Get("properties"))
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["typeSystem"] = nodeNative1
		}
		encryption := make(map[string]interface{})
		nodeNative2, _ := jsonpath.Get("$[0].encryption[0].enable", d.Get("properties"))
		if nodeNative2 != nil && nodeNative2 != "" {
			encryption["enable"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].encryption[0].algorithm", d.Get("properties"))
		if nodeNative3 != nil && nodeNative3 != "" {
			encryption["algorithm"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].encryption[0].key", d.Get("properties"))
		if nodeNative4 != nil && nodeNative4 != "" {
			encryption["key"] = nodeNative4
		}

		objectDataLocalMap["encryption"] = encryption

		request["properties"] = objectDataLocalMap
	}

	if v, ok := d.GetOk("default_quota"); ok {
		request["defaultQuota"] = v
	}
	if v, ok := d.GetOk("comment"); ok {
		request["comment"] = v
	}
	if v, ok := d.GetOk("product_type"); ok {
		request["productType"] = v
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

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

	data1RawObj, _ := jsonpath.Get("$.data", objectRaw)
	data1Raw := make(map[string]interface{})
	if data1RawObj != nil {
		data1Raw = data1RawObj.(map[string]interface{})
	}
	if data1Raw["comment"] != nil {
		d.Set("comment", data1Raw["comment"])
	}
	if data1Raw["createdTime"] != nil {
		d.Set("create_time", data1Raw["createdTime"])
	}
	if data1Raw["defaultQuota"] != nil {
		d.Set("default_quota", data1Raw["defaultQuota"])
	}
	if data1Raw["owner"] != nil {
		d.Set("owner", data1Raw["owner"])
	}
	if data1Raw["status"] != nil {
		d.Set("status", data1Raw["status"])
	}
	if data1Raw["type"] != nil {
		d.Set("type", data1Raw["type"])
	}
	if data1Raw["name"] != nil {
		d.Set("project_name", data1Raw["name"])
	}

	ipWhiteListMaps := make([]map[string]interface{}, 0)
	ipWhiteListMap := make(map[string]interface{})
	ipWhiteList1RawObj, _ := jsonpath.Get("$.data.ipWhiteList", objectRaw)
	ipWhiteList1Raw := make(map[string]interface{})
	if ipWhiteList1RawObj != nil {
		ipWhiteList1Raw = ipWhiteList1RawObj.(map[string]interface{})
	}
	if len(ipWhiteList1Raw) > 0 {
		ipWhiteListMap["ip_list"] = ipWhiteList1Raw["ipList"]
		ipWhiteListMap["vpc_ip_list"] = ipWhiteList1Raw["vpcIpList"]

		ipWhiteListMaps = append(ipWhiteListMaps, ipWhiteListMap)
	}
	if ipWhiteList1RawObj != nil {
		d.Set("ip_white_list", ipWhiteListMaps)
	}
	propertiesMaps := make([]map[string]interface{}, 0)
	propertiesMap := make(map[string]interface{})
	properties1RawObj, _ := jsonpath.Get("$.data.properties", objectRaw)
	properties1Raw := make(map[string]interface{})
	if properties1RawObj != nil {
		properties1Raw = properties1RawObj.(map[string]interface{})
	}
	if len(properties1Raw) > 0 {
		propertiesMap["allow_full_scan"] = properties1Raw["allowFullScan"]
		propertiesMap["enable_decimal2"] = properties1Raw["enableDecimal2"]
		propertiesMap["retention_days"] = properties1Raw["retentionDays"]
		propertiesMap["sql_metering_max"] = properties1Raw["sqlMeteringMax"]
		propertiesMap["timezone"] = properties1Raw["timezone"]
		propertiesMap["type_system"] = properties1Raw["typeSystem"]

		encryptionMaps := make([]map[string]interface{}, 0)
		encryptionMap := make(map[string]interface{})
		encryption1RawObj, _ := jsonpath.Get("$.data.properties.encryption", objectRaw)
		encryption1Raw := make(map[string]interface{})
		if encryption1RawObj != nil {
			encryption1Raw = encryption1RawObj.(map[string]interface{})
		}
		if len(encryption1Raw) > 0 {
			encryptionMap["algorithm"] = encryption1Raw["algorithm"]
			encryptionMap["enable"] = encryption1Raw["enable"]
			encryptionMap["key"] = encryption1Raw["key"]

			encryptionMaps = append(encryptionMaps, encryptionMap)
		}
		propertiesMap["encryption"] = encryptionMaps
		tableLifecycleMaps := make([]map[string]interface{}, 0)
		tableLifecycleMap := make(map[string]interface{})
		tableLifecycle1RawObj, _ := jsonpath.Get("$.data.properties.tableLifecycle", objectRaw)
		tableLifecycle1Raw := make(map[string]interface{})
		if tableLifecycle1RawObj != nil {
			tableLifecycle1Raw = tableLifecycle1RawObj.(map[string]interface{})
		}
		if len(tableLifecycle1Raw) > 0 {
			tableLifecycleMap["type"] = tableLifecycle1Raw["type"]
			tableLifecycleMap["value"] = tableLifecycle1Raw["value"]

			tableLifecycleMaps = append(tableLifecycleMaps, tableLifecycleMap)
		}
		propertiesMap["table_lifecycle"] = tableLifecycleMaps
		propertiesMaps = append(propertiesMaps, propertiesMap)
	}
	if properties1RawObj != nil {
		d.Set("properties", propertiesMaps)
	}
	securityPropertiesMaps := make([]map[string]interface{}, 0)
	securityPropertiesMap := make(map[string]interface{})
	securityProperties1RawObj, _ := jsonpath.Get("$.data.securityProperties", objectRaw)
	securityProperties1Raw := make(map[string]interface{})
	if securityProperties1RawObj != nil {
		securityProperties1Raw = securityProperties1RawObj.(map[string]interface{})
	}
	if len(securityProperties1Raw) > 0 {
		securityPropertiesMap["enable_download_privilege"] = securityProperties1Raw["enableDownloadPrivilege"]
		securityPropertiesMap["label_security"] = securityProperties1Raw["labelSecurity"]
		securityPropertiesMap["object_creator_has_access_permission"] = securityProperties1Raw["objectCreatorHasAccessPermission"]
		securityPropertiesMap["object_creator_has_grant_permission"] = securityProperties1Raw["objectCreatorHasGrantPermission"]
		securityPropertiesMap["using_acl"] = securityProperties1Raw["usingAcl"]
		securityPropertiesMap["using_policy"] = securityProperties1Raw["usingPolicy"]

		projectProtectionMaps := make([]map[string]interface{}, 0)
		projectProtectionMap := make(map[string]interface{})
		projectProtection1RawObj, _ := jsonpath.Get("$.data.securityProperties.projectProtection", objectRaw)
		projectProtection1Raw := make(map[string]interface{})
		if projectProtection1RawObj != nil {
			projectProtection1Raw = projectProtection1RawObj.(map[string]interface{})
		}
		if len(projectProtection1Raw) > 0 {
			projectProtectionMap["exception_policy"] = projectProtection1Raw["exceptionPolicy"]
			projectProtectionMap["protected"] = projectProtection1Raw["protected"]

			projectProtectionMaps = append(projectProtectionMaps, projectProtectionMap)
		}
		securityPropertiesMap["project_protection"] = projectProtectionMaps
		securityPropertiesMaps = append(securityPropertiesMaps, securityPropertiesMap)
	}
	if securityProperties1RawObj != nil {
		d.Set("security_properties", securityPropertiesMaps)
	}

	objectRaw, err = maxComputeServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

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
	d.Partial(true)
	projectName := d.Id()
	action := fmt.Sprintf("/api/v1/projects/%s/meta", projectName)
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["projectName"] = d.Id()

	if d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok || d.HasChange("status") {
		request["status"] = v
	}
	if d.HasChange("properties") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("properties"); !IsNil(v) || d.HasChange("properties") {
		nodeNative, _ := jsonpath.Get("$[0].timezone", v)
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["timezone"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].retention_days", v)
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["retentionDays"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].allow_full_scan", v)
		if nodeNative2 != nil && nodeNative2 != "" {
			objectDataLocalMap["allowFullScan"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].enable_decimal2", v)
		if nodeNative3 != nil && nodeNative3 != "" {
			objectDataLocalMap["enableDecimal2"] = nodeNative3
		}
		tableLifecycle := make(map[string]interface{})
		nodeNative4, _ := jsonpath.Get("$[0].table_lifecycle[0].value", v)
		if nodeNative4 != nil && nodeNative4 != "" {
			tableLifecycle["value"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].table_lifecycle[0].type", v)
		if nodeNative5 != nil && nodeNative5 != "" {
			tableLifecycle["type"] = nodeNative5
		}

		objectDataLocalMap["tableLifecycle"] = tableLifecycle
		nodeNative6, _ := jsonpath.Get("$[0].sql_metering_max", v)
		if nodeNative6 != nil && nodeNative6 != "" {
			objectDataLocalMap["sqlMeteringMax"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].type_system", v)
		if nodeNative7 != nil && nodeNative7 != "" {
			objectDataLocalMap["typeSystem"] = nodeNative7
		}
		encryption := make(map[string]interface{})
		nodeNative8, _ := jsonpath.Get("$[0].encryption[0].enable", v)
		if nodeNative8 != nil && nodeNative8 != "" {
			encryption["enable"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].encryption[0].algorithm", v)
		if nodeNative9 != nil {
			encryption["algorithm"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].encryption[0].key", v)
		if nodeNative10 != nil {
			encryption["key"] = nodeNative10
		}

		if encryption["enable"] == "true" {
			objectDataLocalMap["encryption"] = encryption
		}

		request["properties"] = objectDataLocalMap
	}

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	projectName = d.Id()
	action = fmt.Sprintf("/api/v1/projects/%s/ipWhiteList", projectName)
	conn, err = client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["projectName"] = d.Id()

	if d.HasChange("ip_white_list") {
		update = true
	}
	objectDataLocalMap = make(map[string]interface{})

	if v := d.Get("ip_white_list"); !IsNil(v) || d.HasChange("ip_white_list") {
		nodeNative, _ := jsonpath.Get("$[0].ip_list", v)
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["ipList"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].vpc_ip_list", v)
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["vpcIpList"] = nodeNative1
		}

		request["ipWhiteList"] = objectDataLocalMap
	}

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	projectName = d.Id()
	action = fmt.Sprintf("/api/v1/projects/%s/security", projectName)
	conn, err = client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["projectName"] = d.Id()

	if d.HasChange("security_properties.0.using_acl") {
		update = true
	}
	if v, ok := d.GetOk("security_properties"); ok || d.HasChange("security_properties") {
		jsonPathResult, err := jsonpath.Get("$[0].using_acl", v)
		if err == nil && jsonPathResult != "" {
			request["usingAcl"] = jsonPathResult
		}
	}
	if d.HasChange("security_properties.0.using_policy") {
		update = true
	}
	if v, ok := d.GetOk("security_properties"); ok || d.HasChange("security_properties") {
		jsonPathResult1, err := jsonpath.Get("$[0].using_policy", v)
		if err == nil && jsonPathResult1 != "" {
			request["usingPolicy"] = jsonPathResult1
		}
	}
	if d.HasChange("security_properties.0.object_creator_has_access_permission") {
		update = true
	}
	if v, ok := d.GetOk("security_properties"); ok || d.HasChange("security_properties") {
		jsonPathResult2, err := jsonpath.Get("$[0].object_creator_has_access_permission", v)
		if err == nil && jsonPathResult2 != "" {
			request["objectCreatorHasAccessPermission"] = jsonPathResult2
		}
	}
	if d.HasChange("security_properties.0.object_creator_has_grant_permission") {
		update = true
	}
	if v, ok := d.GetOk("security_properties"); ok || d.HasChange("security_properties") {
		jsonPathResult3, err := jsonpath.Get("$[0].object_creator_has_grant_permission", v)
		if err == nil && jsonPathResult3 != "" {
			request["objectCreatorHasGrantPermission"] = jsonPathResult3
		}
	}
	if d.HasChange("security_properties.0.label_security") {
		update = true
	}
	if v, ok := d.GetOk("security_properties"); ok || d.HasChange("security_properties") {
		jsonPathResult4, err := jsonpath.Get("$[0].label_security", v)
		if err == nil && jsonPathResult4 != "" {
			request["labelSecurity"] = jsonPathResult4
		}
	}
	if d.HasChange("security_properties.0.enable_download_privilege") {
		update = true
	}
	if v, ok := d.GetOk("security_properties"); ok || d.HasChange("security_properties") {
		jsonPathResult5, err := jsonpath.Get("$[0].enable_download_privilege", v)
		if err == nil && jsonPathResult5 != "" {
			request["enableDownloadPrivilege"] = jsonPathResult5
		}
	}
	if d.HasChange("security_properties") {
		update = true
	}
	objectDataLocalMap = make(map[string]interface{})

	if v := d.Get("security_properties"); !IsNil(v) || d.HasChange("security_properties") {
		nodeNative, _ := jsonpath.Get("$[0].project_protection[0].protected", v)
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["protected"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].project_protection[0].exception_policy", v)
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["exceptionPolicy"] = nodeNative1
		}

		request["projectProtection"] = objectDataLocalMap
	}

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	projectName = d.Id()
	action = fmt.Sprintf("/api/v1/projects/%s/status", projectName)
	conn, err = client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["projectName"] = d.Id()

	if d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok || d.HasChange("status") {
		request["status"] = v
	}
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
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
	d.Partial(false)
	return resourceAliCloudMaxComputeProjectRead(d, meta)
}

func resourceAliCloudMaxComputeProjectDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	projectName := d.Id()
	action := fmt.Sprintf("/api/v1/projects/%s", projectName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["projectName"] = d.Id()

	query["isLogical"] = StringPointer("false")
	if v, ok := d.GetOk("is_logical"); ok {
		query["isLogical"] = StringPointer(v.(string))
	}

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	maxComputeServiceV2 := MaxComputeServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, maxComputeServiceV2.MaxComputeProjectStateRefreshFunc(d.Id(), "$.data.status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
