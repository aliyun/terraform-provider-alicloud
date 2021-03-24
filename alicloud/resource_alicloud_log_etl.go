package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudLogETL() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogETLCreate,
		Read:   resourceAlicloudLogETLRead,
		Update: resourceAlicloudLogETLUpdate,
		Delete: resourceAlicloudLogETLDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"etl_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Resident",
			},

			"etl_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  sls.ETLType,
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"STARTING", "RUNNING", "STOPPING", "STOPPED"}, false),
			},

			"create_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"last_modified_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"access_key_id": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"access_key_secret": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"kms_encrypted_access_key_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},

			"kms_encryption_access_key_id_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_access_key_id").(string) == ""
				},
				Elem: schema.TypeString,
			},

			"kms_encrypted_access_key_secret": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},

			"kms_encryption_access_key_secret_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_access_key_secret").(string) == ""
				},
				Elem: schema.TypeString,
			},

			"from_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"to_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"script": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  sls.ETLVersion,
			},
			"logstore": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"etl_sinks": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							Sensitive: true,
						},
						"kms_encrypted_access_key_id": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: kmsDiffSuppressFunc,
						},
						"access_key_secret": {
							Type:     schema.TypeString,
							Optional: true,
							Sensitive: true,
						},
						"kms_encrypted_access_key_secret": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: kmsDiffSuppressFunc,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"project": {
							Type:     schema.TypeString,
							Required: true,
						},
						"logstore": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  sls.ETLSinksType,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudLogETLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestinfo *sls.Client
	logService := LogService{client}
	etlJob, err := getETLJob(d, meta)
	if err != nil {
		return err
	}

	project := d.Get("project").(string)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestinfo = slsClient
			return nil, slsClient.CreateETL(project, etlJob)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("CreateETL", raw, requestinfo, map[string]interface{}{
				"project":  project,
				"logstore": d.Get("logstore").(string),
			})
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_etl", "CreateETL", AliyunLogGoSdkERROR)
	}
	d.SetId(fmt.Sprintf("%s%s%s", project, COLON_SEPARATED, d.Get("etl_name").(string)))
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, logService.LogOssShipperStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudLogETLRead(d, meta)
}

func resourceAlicloudLogETLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	etl, err := logService.DescribeLogEtl(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_log_etl LogService.DescribeLogEtl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("etl_name", parts[1])
	d.Set("project", parts[0])
	d.Set("display_name", etl.DisplayName)
	d.Set("description", etl.Description)
	d.Set("schedule", etl.Schedule.Type)
	d.Set("etl_type", etl.Type)
	d.Set("status", etl.Status)
	d.Set("create_time", etl.CreateTime)
	d.Set("last_modified_time", etl.LastModifiedTime)
	d.Set("from_time", etl.Configuration.FromTime)
	d.Set("to_time", etl.Configuration.ToTime)
	d.Set("script", etl.Configuration.Script)
	d.Set("version", etl.Configuration.Version)
	d.Set("logstore", etl.Configuration.Logstore)
	d.Set("parameters", etl.Configuration.Parameters)
	d.Set("role_arn", etl.Configuration.RoleArn)

	var etl_sinks []map[string]interface{}
	for _, etl_sink := range etl.Configuration.ETLSinks {
		temp := map[string]interface{}{
			"access_key_id":     etl_sink.AccessKeyId,
			"access_key_secret": etl_sink.AccessKeySecret,
			"endpoint":          etl_sink.Endpoint,
			"name":              etl_sink.Name,
			"project":           etl_sink.Project,
			"logstore":          etl_sink.Logstore,
			"role_arn":          etl_sink.RoleArn,
			"type":              etl_sink.Type,
		}
		etl_sinks = append(etl_sinks, temp)
	}
	d.Set("etl_sinks", etl_sinks)
	return nil
}

func resourceAlicloudLogETLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	if d.HasChange("status") {
		logService := LogService{client}
		status := d.Get("status").(string)
		if status == "STARTING" || status == "RUNNING" {
			if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
					return nil, slsClient.StartETL(parts[0], d.Get("etl_name").(string))
				})
				if err != nil {
					if IsExpectedErrors(err, []string{LogClientTimeout}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "StartLogETL", AliyunLogGoSdkERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, logService.LogOssShipperStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		} else if status == "STOPPING" || status == "STOPPED" {
			if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
					return nil, slsClient.StopETL(parts[0], d.Get("etl_name").(string))
				})
				if err != nil {
					if IsExpectedErrors(err, []string{LogClientTimeout}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "StopLogETL", AliyunLogGoSdkERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"STOPPED"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, logService.LogOssShipperStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
		return resourceAlicloudLogETLRead(d, meta)
	}

	if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			etl, err := getETLJob(d, meta)
			if err != nil {
				return nil, err
			}
			return nil, slsClient.UpdateETL(parts[0], etl)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateLogETL", AliyunLogGoSdkERROR)
	}
	return resourceAlicloudLogETLRead(d, meta)
}

func resourceAlicloudLogETLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *sls.Client
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteETL(parts[0], parts[1])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("DeleteLogETL", raw, requestInfo, map[string]interface{}{
				"project_name": parts[0],
				"elt_name":     parts[1],
			})
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_etl", "DeleteLogETL", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogETL(d.Id(), Deleted, DefaultTimeout))
}

func getETLJob(d *schema.ResourceData, meta interface{}) (sls.ETL, error) {
	client := meta.(*connectivity.AliyunClient)
	var etlSinks = []sls.ETLSink{}
	var config = sls.ETLConfiguration{}
	var etlJob = sls.ETL{}
	schedule := sls.ETLSchedule{
		Type: d.Get("schedule").(string),
	}
	parms := map[string]string{}
	if temp, ok := d.GetOk("parameters"); ok {
		for k, v := range temp.(map[string]interface{}) {
			parms[k] = v.(string)
		}
	}
	config = sls.ETLConfiguration{
		FromTime:   int64(d.Get("from_time").(int)),
		Logstore:   d.Get("logstore").(string),
		Parameters: parms,
		RoleArn:    d.Get("role_arn").(string),
		Script:     d.Get("script").(string),
		ToTime:     int32(d.Get("to_time").(int)),
		Version:    int8(d.Get("version").(int)),
	}
	for _, f := range d.Get("etl_sinks").(*schema.Set).List() {
		v := f.(map[string]interface{})
		sink := sls.ETLSink{
			AccessKeyId:     v["access_key_id"].(string),
			AccessKeySecret: v["access_key_secret"].(string),
			Endpoint:        v["endpoint"].(string),
			Name:            v["name"].(string),
			Project:         v["project"].(string),
			Type:            v["type"].(string),
			RoleArn:         v["role_arn"].(string),
			Logstore:        v["logstore"].(string),
		}
		akId := v["access_key_id"].(string)
		kmsAkId := v["kms_encrypted_access_key_id"].(string)
		if akId == "" && kmsAkId == "" {
			return etlJob, WrapError(Error("One of the 'access_key_id' and 'kms_encrypted_access_key_id' should be set."))
		}
		if akId != "" {
			sink.AccessKeyId = akId
		} else {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(kmsAkId, d.Get("kms_encryption_access_key_id_context").(map[string]interface{}))
			if err != nil {
				return etlJob, WrapError(err)
			}
			sink.AccessKeyId = decryptResp.Plaintext
		}
		akSecret := v["access_key_secret"].(string)
		kmsAkSecret := v["kms_encrypted_access_key_secret"].(string)
		if akSecret == "" && kmsAkSecret == "" {
			return etlJob, WrapError(Error("One of the 'access_key_secret' and 'kms_encrypted_access_key_secret' should be set."))
		}
		if akSecret != "" {
			sink.AccessKeySecret = akSecret
		} else {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(kmsAkSecret, d.Get("kms_encryption_access_key_secret_context").(map[string]interface{}))
			if err != nil {
				return etlJob, WrapError(err)
			}
			sink.AccessKeySecret = decryptResp.Plaintext
		}
		etlSinks = append(etlSinks, sink)
	}
	config.ETLSinks = etlSinks

	etlJob = sls.ETL{
		Configuration:    config,
		DisplayName:      d.Get("display_name").(string),
		Description:      d.Get("description").(string),
		Name:             d.Get("etl_name").(string),
		Schedule:         schedule,
		Type:             d.Get("etl_type").(string),
		Status:           d.Get("status").(string),
		CreateTime:       int32(d.Get("create_time").(int)),
		LastModifiedTime: int32(d.Get("last_modified_time").(int)),
	}

	akId := d.Get("access_key_id").(string)
	kmsAkId := d.Get("kms_encrypted_access_key_id").(string)
	if akId == "" && kmsAkId == "" {
		return etlJob, WrapError(Error("One of the 'access_key_id' and 'kms_encrypted_access_key_id' should be set."))
	}
	if akId != "" {
		etlJob.Configuration.AccessKeyId = akId
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsAkId, d.Get("kms_encryption_access_key_id_context").(map[string]interface{}))
		if err != nil {
			return etlJob, WrapError(err)
		}
		etlJob.Configuration.AccessKeyId = decryptResp.Plaintext
	}
	akSecret := d.Get("access_key_secret").(string)
	kmsAkSecret := d.Get("kms_encrypted_access_key_secret").(string)
	if akSecret == "" && kmsAkSecret == "" {
		return etlJob, WrapError(Error("One of the 'access_key_secret' and 'kms_encrypted_access_key_secret' should be set."))
	}
	if akSecret != "" {
		etlJob.Configuration.AccessKeySecret = akSecret
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsAkSecret, d.Get("kms_encryption_access_key_secret_context").(map[string]interface{}))
		if err != nil {
			return etlJob, WrapError(err)
		}
		etlJob.Configuration.AccessKeySecret = decryptResp.Plaintext
	}
	return etlJob, nil
}
