// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsEtl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsEtlCreate,
		Read:   resourceAliCloudSlsEtlRead,
		Update: resourceAliCloudSlsEtlUpdate,
		Delete: resourceAliCloudSlsEtlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"script": {
							Type:     schema.TypeString,
							Required: true,
						},
						"to_time": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"parameters": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"sink": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datasets": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"project": {
										Type:     schema.TypeString,
										Required: true,
									},
									"endpoint": {
										Type:     schema.TypeString,
										Required: true,
									},
									"logstore": {
										Type:     schema.TypeString,
										Required: true,
									},
									"role_arn": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"logstore": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"lang": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"SPL"}, false),
						},
						"from_time": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudSlsEtlCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/etls")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project").(string))
	request["name"] = d.Get("job_name")
	request["displayName"] = d.Get("display_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("configuration"); v != nil {
		script1, _ := jsonpath.Get("$[0].script", v)
		if script1 != nil && script1 != "" {
			objectDataLocalMap["script"] = script1
		}
		lang1, _ := jsonpath.Get("$[0].lang", v)
		if lang1 != nil && lang1 != "" {
			objectDataLocalMap["lang"] = lang1
		}
		fromTime1, _ := jsonpath.Get("$[0].from_time", v)
		if fromTime1 != nil && fromTime1 != "" {
			objectDataLocalMap["fromTime"] = fromTime1
		}
		toTime1, _ := jsonpath.Get("$[0].to_time", v)
		if toTime1 != nil && toTime1 != "" {
			objectDataLocalMap["toTime"] = toTime1
		}
		logstore1, _ := jsonpath.Get("$[0].logstore", v)
		if logstore1 != nil && logstore1 != "" {
			objectDataLocalMap["logstore"] = logstore1
		}
		roleArn1, _ := jsonpath.Get("$[0].role_arn", v)
		if roleArn1 != nil && roleArn1 != "" {
			objectDataLocalMap["roleArn"] = roleArn1
		}
		parameters1, _ := jsonpath.Get("$[0].parameters", v)
		if parameters1 != nil && parameters1 != "" {
			objectDataLocalMap["parameters"] = parameters1
		}
		if v, ok := d.GetOk("configuration"); ok {
			localData, err := jsonpath.Get("$[0].sink", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["name"] = dataLoopTmp["name"]
				dataLoopMap["endpoint"] = dataLoopTmp["endpoint"]
				dataLoopMap["project"] = dataLoopTmp["project"]
				dataLoopMap["logstore"] = dataLoopTmp["logstore"]
				dataLoopMap["datasets"] = dataLoopTmp["datasets"]
				dataLoopMap["roleArn"] = dataLoopTmp["role_arn"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["sinks"] = localMaps
		}

		request["configuration"] = objectDataLocalMap
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateETL", action), query, body, nil, hostMap, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"403"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_etl", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["name"]))

	slsServiceV2 := SlsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, slsServiceV2.SlsEtlStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudSlsEtlRead(d, meta)
}

func resourceAliCloudSlsEtlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsEtl(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_etl DescribeSlsEtl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("display_name", objectRaw["displayName"])
	d.Set("status", objectRaw["status"])

	configurationMaps := make([]map[string]interface{}, 0)
	configurationMap := make(map[string]interface{})
	configurationRaw := make(map[string]interface{})
	if objectRaw["configuration"] != nil {
		configurationRaw = objectRaw["configuration"].(map[string]interface{})
	}
	if len(configurationRaw) > 0 {
		configurationMap["from_time"] = configurationRaw["fromTime"]
		configurationMap["lang"] = configurationRaw["lang"]
		configurationMap["logstore"] = configurationRaw["logstore"]
		configurationMap["parameters"] = configurationRaw["parameters"]
		configurationMap["role_arn"] = configurationRaw["roleArn"]
		configurationMap["script"] = configurationRaw["script"]
		configurationMap["to_time"] = configurationRaw["toTime"]

		sinksRaw := configurationRaw["sinks"]
		sinkMaps := make([]map[string]interface{}, 0)
		if sinksRaw != nil {
			for _, sinksChildRaw := range sinksRaw.([]interface{}) {
				sinkMap := make(map[string]interface{})
				sinksChildRaw := sinksChildRaw.(map[string]interface{})
				sinkMap["endpoint"] = sinksChildRaw["endpoint"]
				sinkMap["logstore"] = sinksChildRaw["logstore"]
				sinkMap["name"] = sinksChildRaw["name"]
				sinkMap["project"] = sinksChildRaw["project"]
				sinkMap["role_arn"] = sinksChildRaw["roleArn"]

				datasetsRaw := make([]interface{}, 0)
				if sinksChildRaw["datasets"] != nil {
					datasetsRaw = sinksChildRaw["datasets"].([]interface{})
				}

				sinkMap["datasets"] = datasetsRaw
				sinkMaps = append(sinkMaps, sinkMap)
			}
		}
		configurationMap["sink"] = sinkMaps
		configurationMaps = append(configurationMaps, configurationMap)
	}
	if err := d.Set("configuration", configurationMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("project", parts[0])
	d.Set("job_name", parts[1])

	return nil
}

func resourceAliCloudSlsEtlUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	etlName := parts[1]
	action := fmt.Sprintf("/etls/%s", etlName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])

	if d.HasChange("display_name") {
		update = true
	}
	request["displayName"] = d.Get("display_name")
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if d.HasChange("configuration") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("configuration"); v != nil {
		script1, _ := jsonpath.Get("$[0].script", v)
		if script1 != nil && (d.HasChange("configuration.0.script") || script1 != "") {
			objectDataLocalMap["script"] = script1
		}
		fromTime1, _ := jsonpath.Get("$[0].from_time", v)
		if fromTime1 != nil && (d.HasChange("configuration.0.from_time") || fromTime1 != "") {
			objectDataLocalMap["fromTime"] = fromTime1
		}
		toTime1, _ := jsonpath.Get("$[0].to_time", v)
		if toTime1 != nil && (d.HasChange("configuration.0.to_time") || toTime1 != "") {
			objectDataLocalMap["toTime"] = toTime1
		}
		logstore1, _ := jsonpath.Get("$[0].logstore", v)
		if logstore1 != nil && (d.HasChange("configuration.0.logstore") || logstore1 != "") {
			objectDataLocalMap["logstore"] = logstore1
		}
		roleArn1, _ := jsonpath.Get("$[0].role_arn", v)
		if roleArn1 != nil && (d.HasChange("configuration.0.role_arn") || roleArn1 != "") {
			objectDataLocalMap["roleArn"] = roleArn1
		}
		parameters1, _ := jsonpath.Get("$[0].parameters", v)
		if parameters1 != nil && (d.HasChange("configuration.0.parameters") || parameters1 != "") {
			objectDataLocalMap["parameters"] = parameters1
		}
		if v, ok := d.GetOk("configuration"); ok {
			localData, err := jsonpath.Get("$[0].sink", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["name"] = dataLoopTmp["name"]
				dataLoopMap["project"] = dataLoopTmp["project"]
				dataLoopMap["logstore"] = dataLoopTmp["logstore"]
				dataLoopMap["datasets"] = dataLoopTmp["datasets"]
				dataLoopMap["roleArn"] = dataLoopTmp["role_arn"]
				dataLoopMap["endpoint"] = dataLoopTmp["endpoint"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["sinks"] = localMaps
		}

		lang1, _ := jsonpath.Get("$[0].lang", v)
		if lang1 != nil && (d.HasChange("configuration.0.lang") || lang1 != "") {
			objectDataLocalMap["lang"] = lang1
		}

		request["configuration"] = objectDataLocalMap
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateETL", action), query, body, nil, hostMap, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"403"}) || NeedRetry(err) {
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
		slsServiceV2 := SlsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, slsServiceV2.SlsEtlStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudSlsEtlRead(d, meta)
}

func resourceAliCloudSlsEtlDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	etlName := parts[1]
	action := fmt.Sprintf("/etls/%s", etlName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteETL", action), query, nil, nil, hostMap, false)

		if err != nil {
			if IsExpectedErrors(err, []string{"403"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
