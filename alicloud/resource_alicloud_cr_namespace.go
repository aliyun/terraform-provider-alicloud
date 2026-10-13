// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
//
// NOTE (RD 补全): DefaultRepoConfiguration block 由 RD 依据 CloudSpec pre 模型
// (alicloud.CR.cr.v20181201/Namespace.cspec, AMP pre QA-PASS) 补齐——ACube 生成器与
// terraform-cspec terraform_pre 均不处理 struct 引用属性 (D2 同类缺陷)，故生成产物不含该块。
// create/update 以 JSON 对象发送 DefaultRepoConfiguration (OpenAPI style:json,
// 与 rds_custom/esa_waf_rule 等的 RpcPost json.Marshal 模式一致)；read 从 GetNamespace
// 响应的 DefaultRepoConfiguration 嵌套对象 flatten。
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudCrNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrNamespaceCreate,
		Read:   resourceAliCloudCrNamespaceRead,
		Update: resourceAliCloudCrNamespaceUpdate,
		Delete: resourceAliCloudCrNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_create_repo": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_repo_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "PUBLIC"}, false),
						},
						"tag_immutability": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"artifact_build_rule_parameters": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_index_only": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"priority_file": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"default_repo_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCrNamespaceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateNamespace"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("namespace_name"); ok {
		request["NamespaceName"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("default_repo_type"); ok {
		request["DefaultRepoType"] = v
	}
	if v, ok := d.GetOkExists("auto_create_repo"); ok {
		request["AutoCreateRepo"] = v
	}
	drc, err := buildCrNamespaceDefaultRepoConfiguration(d)
	if err != nil {
		return err
	}
	if drc != "" {
		request["DefaultRepoConfiguration"] = drc
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_namespace", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], request["NamespaceName"]))

	return resourceAliCloudCrNamespaceRead(d, meta)
}

func resourceAliCloudCrNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crServiceV2 := CrServiceV2{client}

	objectRaw, err := crServiceV2.DescribeCrNamespace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_namespace DescribeCrNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auto_create_repo", objectRaw["AutoCreateRepo"])
	d.Set("default_repo_type", objectRaw["DefaultRepoType"])
	d.Set("namespace_id", objectRaw["NamespaceId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["NamespaceStatus"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("namespace_name", objectRaw["NamespaceName"])
	if v, ok := objectRaw["DefaultRepoConfiguration"]; ok && v != nil {
		if err := d.Set("default_repo_configuration", flattenCrNamespaceDefaultRepoConfiguration(v)); err != nil {
			return WrapError(err)
		}
	} else {
		d.Set("default_repo_configuration", []interface{}{})
	}

	return nil
}

func resourceAliCloudCrNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateNamespace"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["NamespaceName"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("default_repo_type") {
		update = true
		request["DefaultRepoType"] = d.Get("default_repo_type")
	}

	if d.HasChange("auto_create_repo") {
		update = true
		request["AutoCreateRepo"] = d.Get("auto_create_repo")
	}

	if d.HasChange("default_repo_configuration") {
		update = true
		drc, drcErr := buildCrNamespaceDefaultRepoConfiguration(d)
		if drcErr != nil {
			return drcErr
		}
		request["DefaultRepoConfiguration"] = drc
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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

	return resourceAliCloudCrNamespaceRead(d, meta)
}

func resourceAliCloudCrNamespaceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteNamespace"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["NamespaceName"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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

	crServiceV2 := CrServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"DELETING"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, crServiceV2.CrNamespaceStateRefreshFunc(d.Id(), "NamespaceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

// buildCrNamespaceDefaultRepoConfiguration serializes the default_repo_configuration
// block into a JSON string for the DefaultRepoConfiguration query param (OpenAPI style:json).
// 补全依据: CloudSpec pre 模型 Namespace.cspec (create/update paramSerializer,
// Key=DefaultRepoConfiguration, Value=RepoConfiguration 对象 schema)。
func buildCrNamespaceDefaultRepoConfiguration(d *schema.ResourceData) (string, error) {
	configs := d.Get("default_repo_configuration").([]interface{})
	if len(configs) == 0 || configs[0] == nil {
		return "", nil
	}
	config := configs[0].(map[string]interface{})
	drc := map[string]interface{}{}
	if v, ok := config["repo_type"].(string); ok && v != "" {
		drc["RepoType"] = v
	}
	if v, ok := config["tag_immutability"].(bool); ok {
		drc["TagImmutability"] = v
	}
	if params, ok := config["artifact_build_rule_parameters"].([]interface{}); ok && len(params) > 0 && params[0] != nil {
		param := params[0].(map[string]interface{})
		abrp := map[string]interface{}{}
		if v, ok := param["image_index_only"].(bool); ok {
			abrp["ImageIndexOnly"] = v
		}
		if v, ok := param["priority_file"].(string); ok && v != "" {
			abrp["PriorityFile"] = v
		}
		if len(abrp) > 0 {
			drc["ArtifactBuildRuleParameters"] = abrp
		}
	}
	if len(drc) == 0 {
		return "", nil
	}
	configJSON, err := json.Marshal(drc)
	if err != nil {
		return "", WrapError(err)
	}
	return string(configJSON), nil
}

// flattenCrNamespaceDefaultRepoConfiguration flattens the GetNamespace response
// DefaultRepoConfiguration nested object into the schema. 补全依据: CloudSpec pre 模型
// Namespace.cspec (get responsePath $.DefaultRepoConfiguration.*；线上 GetNamespace 响应
// 200 含 DefaultRepoConfiguration $ref RepoConfiguration 已核验)。
func flattenCrNamespaceDefaultRepoConfiguration(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	config, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	result := map[string]interface{}{
		"repo_type":        config["RepoType"],
		"tag_immutability": config["TagImmutability"],
	}
	if abrp, ok := config["ArtifactBuildRuleParameters"].(map[string]interface{}); ok {
		result["artifact_build_rule_parameters"] = []map[string]interface{}{
			{
				"image_index_only": abrp["ImageIndexOnly"],
				"priority_file":    abrp["PriorityFile"],
			},
		}
	}
	return []map[string]interface{}{result}
}
