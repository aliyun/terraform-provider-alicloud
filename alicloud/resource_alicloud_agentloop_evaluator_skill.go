// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAgentloopEvaluatorSkill() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopEvaluatorSkillCreate,
		Read:   resourceAliCloudAgentloopEvaluatorSkillRead,
		Update: resourceAliCloudAgentloopEvaluatorSkillUpdate,
		Delete: resourceAliCloudAgentloopEvaluatorSkillDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agent_space": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"files": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"skill_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAgentloopEvaluatorSkillCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	name := d.Get("name")
	action := fmt.Sprintf("/api/v1/evaluator/%s/skill", name)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["agentSpace"] = StringPointer(d.Get("agent_space").(string))
	if v, ok := d.GetOk("skill_name"); ok {
		request["skillName"] = v
	}

	if v, ok := d.GetOk("files"); ok {
		filesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["name"] = dataLoopTmp["name"]
			dataLoopMap["content"] = dataLoopTmp["content"]
			dataLoopMap["remark"] = dataLoopTmp["remark"]
			filesMapsArray = append(filesMapsArray, dataLoopMap)
		}
		request["files"] = filesMapsArray
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("display_name"); ok {
		request["displayName"] = v
	}
	if v, ok := d.GetOkExists("enable"); ok {
		request["enable"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_evaluator_skill", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", *query["agentSpace"], name, request["skillName"]))

	return resourceAliCloudAgentloopEvaluatorSkillRead(d, meta)
}

func resourceAliCloudAgentloopEvaluatorSkillRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopServiceV2 := AgentloopServiceV2{client}

	objectRaw, err := agentloopServiceV2.DescribeAgentloopEvaluatorSkill(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_evaluator_skill DescribeAgentloopEvaluatorSkill Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	skillRawObj, _ := jsonpath.Get("$.skill", objectRaw)
	skillRaw := make(map[string]interface{})
	if skillRawMap, ok := skillRawObj.(map[string]interface{}); ok {
		skillRaw = skillRawMap
	}
	d.Set("created_at", skillRaw["createdAt"])
	d.Set("description", skillRaw["description"])
	d.Set("display_name", skillRaw["displayName"])
	d.Set("enable", skillRaw["enable"])
	d.Set("updated_at", skillRaw["updatedAt"])
	d.Set("skill_name", skillRaw["skillName"])

	// GetEvaluatorSkill does not return the per-file remark. Preserve the
	// configured remark (matched by file name) so the read-back does not
	// produce a perpetual diff.
	configuredRemarks := make(map[string]interface{})
	for _, f := range convertToInterfaceArray(d.Get("files")) {
		if fm, ok := f.(map[string]interface{}); ok {
			configuredRemarks[fmt.Sprint(fm["name"])] = fm["remark"]
		}
	}
	filesRaw, _ := jsonpath.Get("$.skill.files", objectRaw)
	filesMaps := make([]map[string]interface{}, 0)
	if filesRaw != nil {
		for _, filesChildRawObj := range convertToInterfaceArray(filesRaw) {
			if filesChildRaw, ok := filesChildRawObj.(map[string]interface{}); ok {
				filesMap := make(map[string]interface{})
				filesMap["content"] = filesChildRaw["content"]
				filesMap["name"] = filesChildRaw["name"]
				filesMap["remark"] = filesChildRaw["remark"]
				if filesMap["remark"] == nil || filesMap["remark"] == "" {
					filesMap["remark"] = configuredRemarks[fmt.Sprint(filesChildRaw["name"])]
				}

				filesMaps = append(filesMaps, filesMap)
			}
		}
	}
	// The server does not guarantee the file order, while the schema keeps
	// files as an order-sensitive TypeList. Rebuild the state in the order of
	// the current configuration (matched by file name) so a refresh never
	// produces a perpetual diff; files returned by the API but absent from
	// the configuration are appended in name order. On import there is no
	// configuration to match, so fall back to a stable name ordering.
	filesByName := make(map[string]map[string]interface{})
	for _, filesMap := range filesMaps {
		filesByName[fmt.Sprint(filesMap["name"])] = filesMap
	}
	orderedFiles := make([]map[string]interface{}, 0, len(filesMaps))
	consumed := make(map[string]bool)
	for _, f := range convertToInterfaceArray(d.Get("files")) {
		if fm, ok := f.(map[string]interface{}); ok {
			name := fmt.Sprint(fm["name"])
			if filesMap, exists := filesByName[name]; exists && !consumed[name] {
				orderedFiles = append(orderedFiles, filesMap)
				consumed[name] = true
			}
		}
	}
	remaining := make([]map[string]interface{}, 0)
	for _, filesMap := range filesMaps {
		if !consumed[fmt.Sprint(filesMap["name"])] {
			remaining = append(remaining, filesMap)
		}
	}
	sort.Slice(remaining, func(i, j int) bool {
		return fmt.Sprint(remaining[i]["name"]) < fmt.Sprint(remaining[j]["name"])
	})
	orderedFiles = append(orderedFiles, remaining...)
	if err := d.Set("files", orderedFiles); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("agent_space", parts[0])
	d.Set("name", parts[1])

	return nil
}

func resourceAliCloudAgentloopEvaluatorSkillUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	name := parts[1]
	skillName := parts[2]
	action := fmt.Sprintf("/api/v1/evaluator/%s/skill/%s", name, skillName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["agentSpace"] = StringPointer(parts[0])

	if d.HasChange("files") {
		update = true
		if v, ok := d.GetOk("files"); ok {
			filesMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["name"] = dataLoopTmp["name"]
				dataLoopMap["content"] = dataLoopTmp["content"]
				dataLoopMap["remark"] = dataLoopTmp["remark"]
				filesMapsArray = append(filesMapsArray, dataLoopMap)
			}
			request["files"] = filesMapsArray
		}
	}

	// The backend rejects an update without enable ("enable is mandatory for
	// this action") although the spec marks it optional, and the PUT applies
	// full-replacement semantics that blank out description/displayName when
	// they are absent, so always send the current values of all three.
	request["enable"] = d.Get("enable")
	request["displayName"] = d.Get("display_name")
	request["description"] = d.Get("description")

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
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

	d.Partial(false)
	return resourceAliCloudAgentloopEvaluatorSkillRead(d, meta)
}

func resourceAliCloudAgentloopEvaluatorSkillDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	name := parts[1]
	skillName := parts[2]
	action := fmt.Sprintf("/api/v1/evaluator/%s/skill/%s", name, skillName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["agentSpace"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EvaluatorSkillNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
