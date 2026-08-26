// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudIaCServiceModule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudIaCServiceModuleCreate,
		Read:   resourceAliCloudIaCServiceModuleRead,
		Update: resourceAliCloudIaCServiceModuleUpdate,
		Delete: resourceAliCloudIaCServiceModuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"module_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"version_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudIaCServiceModuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/modules")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["name"] = d.Get("module_name")
	request["source"] = d.Get("source")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("source_path"); ok {
		request["sourcePath"] = v
	}
	if v, ok := d.GetOk("state_path"); ok {
		request["statePath"] = v
	}
	if v, ok := d.GetOk("version_strategy"); ok {
		request["versionStrategy"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["tagKey"] = dataLoopTmp["tag_key"]
			dataLoopMap["tagValue"] = dataLoopTmp["tag_value"]
			tagsMapsArray = append(tagsMapsArray, dataLoopMap)
		}
		request["tags"] = tagsMapsArray
	}

	request["clientToken"] = buildClientToken(action)

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("IaCService", "2021-08-06", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ia_c_service_module", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.moduleId", response)
	d.SetId(fmt.Sprint(id))

	iaCServiceServiceV2 := IaCServiceServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, iaCServiceServiceV2.IaCServiceModuleStateRefreshFunc(d.Id(), "status", []string{"Errored"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudIaCServiceModuleRead(d, meta)
}

func resourceAliCloudIaCServiceModuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	iaCServiceServiceV2 := IaCServiceServiceV2{client}

	objectRaw, err := iaCServiceServiceV2.DescribeIaCServiceModule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ia_c_service_module DescribeIaCServiceModule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("latest_version", objectRaw["latestVersion"])
	d.Set("module_name", objectRaw["name"])
	d.Set("output_path", objectRaw["outputPath"])
	d.Set("source", objectRaw["source"])
	d.Set("source_path", objectRaw["sourcePath"])
	d.Set("state_path", objectRaw["statePath"])
	d.Set("status", objectRaw["status"])
	d.Set("version_strategy", objectRaw["versionStrategy"])

	groupInfoRaw := objectRaw["groupInfo"]
	groupInfoMaps := make([]map[string]interface{}, 0)
	if groupInfoRaw != nil {
		groupInfoMap := make(map[string]interface{})
		groupInfoChildRaw := groupInfoRaw.(map[string]interface{})
		groupInfoMap["group_id"] = groupInfoChildRaw["groupId"]
		groupInfoMap["group_name"] = groupInfoChildRaw["groupName"]
		groupInfoMap["project_id"] = groupInfoChildRaw["projectId"]
		groupInfoMap["project_name"] = groupInfoChildRaw["projectName"]

		groupInfoMaps = append(groupInfoMaps, groupInfoMap)
	}
	if err := d.Set("group_info", groupInfoMaps); err != nil {
		return err
	}
	tagsRaw := objectRaw["tags"]
	tagsMaps := make([]map[string]interface{}, 0)
	if tagsRaw != nil {
		for _, tagsChildRaw := range convertToInterfaceArray(tagsRaw) {
			tagsMap := make(map[string]interface{})
			tagsChildRaw := tagsChildRaw.(map[string]interface{})
			tagsMap["tag_key"] = tagsChildRaw["tagKey"]
			tagsMap["tag_value"] = tagsChildRaw["tagValue"]

			tagsMaps = append(tagsMaps, tagsMap)
		}
	}
	if err := d.Set("tags", tagsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudIaCServiceModuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	moduleId := d.Id()
	action := fmt.Sprintf("/modules/%s", moduleId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("module_name") {
		update = true
	}
	request["name"] = d.Get("module_name")
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOkExists("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("source_path") {
		update = true
	}
	if v, ok := d.GetOkExists("source_path"); ok || d.HasChange("source_path") {
		request["sourcePath"] = v
	}
	if d.HasChange("state_path") {
		update = true
	}
	if v, ok := d.GetOkExists("state_path"); ok || d.HasChange("state_path") {
		request["statePath"] = v
	}
	if d.HasChange("version_strategy") {
		update = true
	}
	if v, ok := d.GetOkExists("version_strategy"); ok || d.HasChange("version_strategy") {
		request["versionStrategy"] = v
	}
	if d.HasChange("tags") {
		update = true
		tagsMapsArray := make([]interface{}, 0)
		if v, ok := d.GetOk("tags"); ok {
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["tagKey"] = dataLoopTmp["tag_key"]
				dataLoopMap["tagValue"] = dataLoopTmp["tag_value"]
				tagsMapsArray = append(tagsMapsArray, dataLoopMap)
			}
		}
		request["tags"] = tagsMapsArray
	}

	request["clientToken"] = buildClientToken(action)

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("IaCService", "2021-08-06", action, query, nil, body, true)
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

		iaCServiceServiceV2 := IaCServiceServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, iaCServiceServiceV2.IaCServiceModuleStateRefreshFunc(d.Id(), "status", []string{"Errored"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudIaCServiceModuleRead(d, meta)
}

func resourceAliCloudIaCServiceModuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	moduleId := d.Id()
	action := fmt.Sprintf("/modules/%s", moduleId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("IaCService", "2021-08-06", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"InvalidModule.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
