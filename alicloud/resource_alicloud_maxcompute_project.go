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
	if v, ok := d.GetOk("project_name"); ok {
		request["name"] = v
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("properties"); !IsNil(v) {
		sqlMeteringMax1, _ := jsonpath.Get("$[0].sql_metering_max", v)
		if sqlMeteringMax1 != nil && sqlMeteringMax1 != "" {
			objectDataLocalMap["sqlMeteringMax"] = sqlMeteringMax1
		}
		typeSystem1, _ := jsonpath.Get("$[0].type_system", v)
		if typeSystem1 != nil && typeSystem1 != "" {
			objectDataLocalMap["typeSystem"] = typeSystem1
		}
		encryption := make(map[string]interface{})
		enable1, _ := jsonpath.Get("$[0].encryption[0].enable", v)
		if enable1 != nil && enable1 != "" {
			encryption["enable"] = enable1
		}
		algorithm1, _ := jsonpath.Get("$[0].encryption[0].algorithm", v)
		if algorithm1 != nil && algorithm1 != "" {
			encryption["algorithm"] = algorithm1
		}
		key1, _ := jsonpath.Get("$[0].encryption[0].key", v)
		if key1 != nil && key1 != "" {
			encryption["key"] = key1
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

	if objectRaw["comment"] != nil {
		d.Set("comment", objectRaw["comment"])
	}
	if objectRaw["createdTime"] != nil {
		d.Set("create_time", objectRaw["createdTime"])
	}
	if objectRaw["defaultQuota"] != nil {
		d.Set("default_quota", objectRaw["defaultQuota"])
	}
	if objectRaw["owner"] != nil {
		d.Set("owner", objectRaw["owner"])
	}
	if objectRaw["regionId"] != nil {
		d.Set("region_id", objectRaw["regionId"])
	}
	if objectRaw["status"] != nil {
		d.Set("status", objectRaw["status"])
	}
	if objectRaw["type"] != nil {
		d.Set("type", objectRaw["type"])
	}
	if objectRaw["name"] != nil {
		d.Set("project_name", objectRaw["name"])
	}

	ipWhiteListMaps := make([]map[string]interface{}, 0)
	ipWhiteListMap := make(map[string]interface{})
	ipWhiteList1Raw := make(map[string]interface{})
	if objectRaw["ipWhiteList"] != nil {
		ipWhiteList1Raw = objectRaw["ipWhiteList"].(map[string]interface{})
	}
	if len(ipWhiteList1Raw) > 0 {
		ipWhiteListMap["ip_list"] = ipWhiteList1Raw["ipList"]
		ipWhiteListMap["vpc_ip_list"] = ipWhiteList1Raw["vpcIpList"]

		ipWhiteListMaps = append(ipWhiteListMaps, ipWhiteListMap)
	}
	if objectRaw["ipWhiteList"] != nil {
		if err := d.Set("ip_white_list", ipWhiteListMaps); err != nil {
			return err
		}
	}
	propertiesMaps := make([]map[string]interface{}, 0)
	propertiesMap := make(map[string]interface{})
	properties1Raw := make(map[string]interface{})
	if objectRaw["properties"] != nil {
		properties1Raw = objectRaw["properties"].(map[string]interface{})
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
		encryption1Raw := make(map[string]interface{})
		if properties1Raw["encryption"] != nil {
			encryption1Raw = properties1Raw["encryption"].(map[string]interface{})
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
		tableLifecycle1Raw := make(map[string]interface{})
		if properties1Raw["tableLifecycle"] != nil {
			tableLifecycle1Raw = properties1Raw["tableLifecycle"].(map[string]interface{})
		}
		if len(tableLifecycle1Raw) > 0 {
			tableLifecycleMap["type"] = tableLifecycle1Raw["type"]
			tableLifecycleMap["value"] = tableLifecycle1Raw["value"]

			tableLifecycleMaps = append(tableLifecycleMaps, tableLifecycleMap)
		}
		propertiesMap["table_lifecycle"] = tableLifecycleMaps
		propertiesMaps = append(propertiesMaps, propertiesMap)
	}
	if objectRaw["properties"] != nil {
		if err := d.Set("properties", propertiesMaps); err != nil {
			return err
		}
	}
	securityPropertiesMaps := make([]map[string]interface{}, 0)
	securityPropertiesMap := make(map[string]interface{})
	securityProperties1Raw := make(map[string]interface{})
	if objectRaw["securityProperties"] != nil {
		securityProperties1Raw = objectRaw["securityProperties"].(map[string]interface{})
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
		projectProtection1Raw := make(map[string]interface{})
		if securityProperties1Raw["projectProtection"] != nil {
			projectProtection1Raw = securityProperties1Raw["projectProtection"].(map[string]interface{})
		}
		if len(projectProtection1Raw) > 0 {
			projectProtectionMap["exception_policy"] = projectProtection1Raw["exceptionPolicy"]
			projectProtectionMap["protected"] = projectProtection1Raw["protected"]

			projectProtectionMaps = append(projectProtectionMaps, projectProtectionMap)
		}
		securityPropertiesMap["project_protection"] = projectProtectionMaps
		securityPropertiesMaps = append(securityPropertiesMaps, securityPropertiesMap)
	}
	if objectRaw["securityProperties"] != nil {
		if err := d.Set("security_properties", securityPropertiesMaps); err != nil {
			return err
		}
	}

	objectRaw, err = maxComputeServiceV2.DescribeProjectListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
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

	if v := d.Get("properties"); v != nil {
		timezone1, _ := jsonpath.Get("$[0].timezone", v)
		if timezone1 != nil && (d.HasChange("properties.0.timezone") || timezone1 != "") {
			objectDataLocalMap["timezone"] = timezone1
		}
		retentionDays1, _ := jsonpath.Get("$[0].retention_days", v)
		if retentionDays1 != nil && (d.HasChange("properties.0.retention_days") || retentionDays1 != "") {
			objectDataLocalMap["retentionDays"] = retentionDays1
		}
		allowFullScan1, _ := jsonpath.Get("$[0].allow_full_scan", v)
		if allowFullScan1 != nil && (d.HasChange("properties.0.allow_full_scan") || allowFullScan1 != "") {
			objectDataLocalMap["allowFullScan"] = allowFullScan1
		}
		enableDecimal21, _ := jsonpath.Get("$[0].enable_decimal2", v)
		if enableDecimal21 != nil && (d.HasChange("properties.0.enable_decimal2") || enableDecimal21 != "") {
			objectDataLocalMap["enableDecimal2"] = enableDecimal21
		}
		tableLifecycle := make(map[string]interface{})
		value1, _ := jsonpath.Get("$[0].table_lifecycle[0].value", v)
		if value1 != nil && (d.HasChange("properties.0.table_lifecycle.0.value") || value1 != "") {
			tableLifecycle["value"] = value1
		}
		type1, _ := jsonpath.Get("$[0].table_lifecycle[0].type", v)
		if type1 != nil && (d.HasChange("properties.0.table_lifecycle.0.type") || type1 != "") {
			tableLifecycle["type"] = type1
		}

		objectDataLocalMap["tableLifecycle"] = tableLifecycle
		sqlMeteringMax1, _ := jsonpath.Get("$[0].sql_metering_max", v)
		if sqlMeteringMax1 != nil && (d.HasChange("properties.0.sql_metering_max") || sqlMeteringMax1 != "") {
			objectDataLocalMap["sqlMeteringMax"] = sqlMeteringMax1
		}
		typeSystem1, _ := jsonpath.Get("$[0].type_system", v)
		if typeSystem1 != nil && (d.HasChange("properties.0.type_system") || typeSystem1 != "") {
			objectDataLocalMap["typeSystem"] = typeSystem1
		}
		encryption := make(map[string]interface{})
		enable1, _ := jsonpath.Get("$[0].encryption[0].enable", v)
		if enable1 != nil && (d.HasChange("properties.0.encryption.0.enable") || enable1 != "") {
			encryption["enable"] = enable1
		}
		algorithm1, _ := jsonpath.Get("$[0].encryption[0].algorithm", v)
		if algorithm1 != nil && (d.HasChange("properties.0.encryption.0.algorithm") || algorithm1 != "") {
			encryption["algorithm"] = algorithm1
		}
		key1, _ := jsonpath.Get("$[0].encryption[0].key", v)
		if key1 != nil && (d.HasChange("properties.0.encryption.0.key") || key1 != "") {
			encryption["key"] = key1
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	projectName = d.Id()
	action = fmt.Sprintf("/api/v1/projects/%s/quota", projectName)
	conn, err = client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["projectName"] = d.Id()

	if !d.IsNewResource() && d.HasChange("default_quota") {
		update = true
	}
	if v, ok := d.GetOk("default_quota"); ok || d.HasChange("default_quota") {
		request["quota"] = v
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

	if v := d.Get("ip_white_list"); v != nil {
		ipList1, _ := jsonpath.Get("$[0].ip_list", v)
		if ipList1 != nil && (d.HasChange("ip_white_list.0.ip_list") || ipList1 != "") {
			objectDataLocalMap["ipList"] = ipList1
		}
		vpcIpList1, _ := jsonpath.Get("$[0].vpc_ip_list", v)
		if vpcIpList1 != nil && (d.HasChange("ip_white_list.0.vpc_ip_list") || vpcIpList1 != "") {
			objectDataLocalMap["vpcIpList"] = vpcIpList1
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

	if v := d.Get("security_properties"); v != nil {
		protected1, _ := jsonpath.Get("$[0].project_protection[0].protected", v)
		if protected1 != nil && (d.HasChange("security_properties.0.project_protection.0.protected") || protected1 != "") {
			objectDataLocalMap["protected"] = protected1
		}
		exceptionPolicy1, _ := jsonpath.Get("$[0].project_protection[0].exception_policy", v)
		if exceptionPolicy1 != nil && (d.HasChange("security_properties.0.project_protection.0.exception_policy") || exceptionPolicy1 != "") {
			objectDataLocalMap["exceptionPolicy"] = exceptionPolicy1
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
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
