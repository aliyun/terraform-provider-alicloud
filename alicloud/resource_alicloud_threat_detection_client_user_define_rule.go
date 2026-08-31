// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudThreatDetectionClientUserDefineRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionClientUserDefineRuleCreate,
		Read:   resourceAliCloudThreatDetectionClientUserDefineRuleRead,
		Update: resourceAliCloudThreatDetectionClientUserDefineRuleUpdate,
		Delete: resourceAliCloudThreatDetectionClientUserDefineRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{0, 1}),
			},
			"client_user_define_rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cmdline": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"new_file_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_cmdline": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_proc_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"windows", "linux", "all"}, false),
			},
			"port_str": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "The port number. When the value of the Type attribute is 3, the PortStr attribute is required. Value range: **1-65535**."),
			},
			"proc_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"registry_content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"registry_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 7}),
			},
		},
	}
}

func resourceAliCloudThreatDetectionClientUserDefineRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddClientUserDefineRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Type"] = d.Get("type")
	request["ActionType"] = d.Get("action_type")
	if v, ok := d.GetOk("proc_path"); ok {
		request["ProcPath"] = v
	}
	if v, ok := d.GetOk("cmdline"); ok {
		request["Cmdline"] = v
	}
	request["Platform"] = d.Get("platform")
	if v, ok := d.GetOk("file_path"); ok {
		request["FilePath"] = v
	}
	if v, ok := d.GetOk("registry_key"); ok {
		request["RegistryKey"] = v
	}
	if v, ok := d.GetOk("registry_content"); ok {
		request["RegistryContent"] = v
	}
	if v, ok := d.GetOk("new_file_path"); ok {
		request["NewFilePath"] = v
	}
	if v, ok := d.GetOk("parent_proc_path"); ok {
		request["ParentProcPath"] = v
	}
	if v, ok := d.GetOk("parent_cmdline"); ok {
		request["ParentCmdline"] = v
	}
	if v, ok := d.GetOk("port_str"); ok {
		request["PortStr"] = v
	}
	request["Name"] = d.Get("client_user_define_rule_name")
	if v, ok := d.GetOk("hash"); ok {
		request["Md5List"] = v
	}
	if v, ok := d.GetOk("ip"); ok {
		request["IP"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_client_user_define_rule", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.UserDefineRuleAddResult.Id", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudThreatDetectionClientUserDefineRuleRead(d, meta)
}

func resourceAliCloudThreatDetectionClientUserDefineRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionClientUserDefineRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_client_user_define_rule DescribeThreatDetectionClientUserDefineRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("action_type", objectRaw["ActionType"])
	d.Set("client_user_define_rule_name", objectRaw["Name"])
	d.Set("cmdline", objectRaw["Cmdline"])
	d.Set("create_time", objectRaw["GmtCreate"])
	d.Set("file_path", objectRaw["FilePath"])
	d.Set("hash", objectRaw["Md5List"])
	d.Set("ip", objectRaw["IP"])
	d.Set("new_file_path", objectRaw["NewFilePath"])
	d.Set("parent_cmdline", objectRaw["ParentCmdline"])
	d.Set("parent_proc_path", objectRaw["ParentProcPath"])
	d.Set("platform", objectRaw["Platform"])
	d.Set("port_str", objectRaw["PortStr"])
	d.Set("proc_path", objectRaw["ProcPath"])
	d.Set("registry_content", objectRaw["RegistryContent"])
	d.Set("registry_key", objectRaw["RegistryKey"])
	d.Set("type", objectRaw["Type"])

	return nil
}

func resourceAliCloudThreatDetectionClientUserDefineRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyClientUserDefineRule"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = d.Id()
	if d.HasChange("action_type") {
		update = true
	}
	request["ActionType"] = d.Get("action_type")
	if d.HasChange("proc_path") {
		update = true
		request["ProcPath"] = d.Get("proc_path")
	}
	if v, ok := d.GetOk("proc_path"); ok {
		request["ProcPath"] = v
	}

	if d.HasChange("cmdline") {
		update = true
		request["Cmdline"] = d.Get("cmdline")
	}
	if v, ok := d.GetOk("cmdline"); ok {
		request["Cmdline"] = v
	}

	if d.HasChange("platform") {
		request["Platform"] = d.Get("platform")
		update = true
	}
	if v, ok := d.GetOk("platform"); ok {
		request["Platform"] = v
	}

	if d.HasChange("ip") {
		update = true
		request["IP"] = d.Get("ip")
	}
	if v, ok := d.GetOk("ip"); ok {
		request["IP"] = v
	}
	request["Platform"] = d.Get("platform")
	if d.HasChange("file_path") {
		update = true
		request["FilePath"] = d.Get("file_path")
	}
	if v, ok := d.GetOk("file_path"); ok {
		request["FilePath"] = v
	}

	if d.HasChange("registry_key") {
		update = true
		request["RegistryKey"] = d.Get("registry_key")
	}
	if v, ok := d.GetOk("registry_key"); ok {
		request["RegistryKey"] = v
	}

	if d.HasChange("registry_content") {
		update = true
		request["RegistryContent"] = d.Get("registry_content")
	}
	if v, ok := d.GetOk("registry_content"); ok {
		request["RegistryContent"] = v
	}

	if d.HasChange("new_file_path") {
		update = true
		request["NewFilePath"] = d.Get("new_file_path")
	}
	if v, ok := d.GetOk("new_file_path"); ok {
		request["NewFilePath"] = v
	}

	if d.HasChange("parent_proc_path") {
		update = true
		request["ParentProcPath"] = d.Get("parent_proc_path")
	}
	if v, ok := d.GetOk("parent_proc_path"); ok {
		request["ParentProcPath"] = v
	}

	if d.HasChange("parent_cmdline") {
		update = true
		request["ParentCmdline"] = d.Get("parent_cmdline")
	}
	if v, ok := d.GetOk("parent_cmdline"); ok {
		request["ParentCmdline"] = v
	}

	if d.HasChange("port_str") {
		update = true
		request["PortStr"] = d.Get("port_str")
	}
	if v, ok := d.GetOk("port_str"); ok {
		request["PortStr"] = v
	}

	if d.HasChange("hash") {
		update = true
		request["Md5List"] = d.Get("hash")
	}
	if v, ok := d.GetOk("hash"); ok {
		request["Md5List"] = v
	}

	if d.HasChange("client_user_define_rule_name") {
		update = true
	}
	request["Name"] = d.Get("client_user_define_rule_name")
	if d.HasChange("type") {
		update = true
	}
	request["Type"] = d.Get("type")
	if d.HasChange("ip") {
		request["IP"] = d.Get("ip")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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

	return resourceAliCloudThreatDetectionClientUserDefineRuleRead(d, meta)
}

func resourceAliCloudThreatDetectionClientUserDefineRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteClientUserDefineRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["IdList.1"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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

	return nil
}
