package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudLogAudit() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogAuditCreate,
		Read:   resourceAlicloudLogAuditRead,
		Update: resourceAlicloudLogAuditUpdate,
		Delete: resourceAlicloudLogAuditDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aliuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"variable_map": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"multi_account": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceAlicloudLogAuditCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("display_name").(string))
	return resourceAlicloudLogAuditUpdate(d, meta)
}

func resourceAlicloudLogAuditUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := slsPop.CreateAnalyzeAppLogRequest()
	request.AppType = "audit"
	request.DisplayName = d.Id()

	var variableMap = map[string]interface{}{}
	mutiAccount := expandStringList(d.Get("multi_account").(*schema.Set).List())
	if len(mutiAccount) > 0 {
		mutiAccountMap := map[string]string{}
		mutiAccountList := []map[string]string{}
		for _, v := range mutiAccount {
			mutiAccountMap["uid"] = v
			mutiAccountList = append(mutiAccountList, mutiAccountMap)
		}
		data, err := json.Marshal(mutiAccountList)
		if err != nil {
			return WrapError(err)
		}
		resultMutiAccount := string(data)
		variableMap["multi_account"] = resultMutiAccount
	}
	variableMap["region"] = client.RegionId
	variableMap["aliuid"] = d.Get("aliuid").(string)
	variableMap["project"] = fmt.Sprintf("slsaudit-center-%s-%s", variableMap["aliuid"], variableMap["region"])
	variableMap["logstore"] = "xx"

	if tempMap, ok := d.GetOk("variable_map"); ok {
		for k, v := range tempMap.(map[string]interface{}) {
			variableMap[k] = v
		}
	}

	b, err := json.Marshal(variableMap)
	if err != nil {
		return WrapError(err)
	} else {
		request.VariableMap = string(b[:])
	}
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		rep, err := client.WithLogPopClient(func(client *slsPop.Client) (interface{}, error) {
			return client.AnalyzeAppLog(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), rep, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_audit", request.GetActionName(), AliyunLogGoSdkERROR)
	}
	return resourceAlicloudLogAuditRead(d, meta)
}

func resourceAlicloudLogAuditRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	response, err := logService.DescribeLogAudit(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	displayName, initMap, err := getInitParameter(response.GetHttpContentString())
	if err != nil {
		return WrapError(err)
	}
	d.Set("display_name", displayName)
	d.Set("aliuid", initMap["aliuid"].(string))
	delete(initMap, "region")
	delete(initMap, "aliuid")
	delete(initMap, "project")
	delete(initMap, "logstore")
	delete(initMap, "multi_account")
	d.Set("variable_map", initMap)
	return nil
}

func resourceAlicloudLogAuditDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudLogAuditInstance.")
	return nil
}

func getInitParameter(rep string) (displayName string, initMap map[string]interface{}, err error) {
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(rep), &m)
	if _, ok := m["AppModel"].(map[string]interface{}); ok {
		model := make(map[string]interface{})
		err = json.Unmarshal([]byte(rep), &model)
		if d, ok := model["AppModel"].(map[string]interface{}); ok {
			displayName = d["DisplayName"].(string)
			configNew := d["Config"]
			m := make(map[string]interface{})
			err = json.Unmarshal([]byte(configNew.(string)), &m)
			initMap = m["initParam"].(map[string]interface{})
		}
	}
	return displayName, initMap, err
}
