package alicloud

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaRoutine() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaRoutineCreate,
		Read:   resourceAliCloudEsaRoutineRead,
		Update: resourceAliCloudEsaRoutineUpdate,
		Delete: resourceAliCloudEsaRoutineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deploy_env": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"staging", "production"}, false),
			},
			"latest_code_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaRoutineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRoutine"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "LockFailed"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_routine", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Name"]))

	// Upload and commit the routine code if code content is provided.
	if code := d.Get("code").(string); code != "" {
		if err := esaRoutineUploadAndCommit(d, meta, code); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliCloudEsaRoutineRead(d, meta)
}

func resourceAliCloudEsaRoutineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaRoutine(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_routine DescribeEsaRoutine Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])

	d.Set("name", d.Id())

	// Refresh the latest committed code version (best-effort).
	if version, err := esaServiceV2.DescribeEsaRoutineLatestCodeVersion(d.Id()); err == nil && version != "" {
		d.Set("latest_code_version", version)
	}

	return nil
}

func resourceAliCloudEsaRoutineUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("code") || d.HasChange("code_description") || d.HasChange("deploy_env") {
		if code := d.Get("code").(string); code != "" {
			if err := esaRoutineUploadAndCommit(d, meta, code); err != nil {
				return WrapError(err)
			}
		}
	}

	return resourceAliCloudEsaRoutineRead(d, meta)
}

func resourceAliCloudEsaRoutineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRoutine"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Name"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "LockFailed"}) || NeedRetry(err) {
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

	return nil
}

// esaRoutineUploadAndCommit runs the staging code lifecycle: fetch an OSS upload
// credential (GetRoutineStagingCodeUploadInfo), upload the code to OSS via a form
// POST, then commit the staged code into a formal code version
// (CommitRoutineStagingCode). The committed version is stored in latest_code_version.
func esaRoutineUploadAndCommit(d *schema.ResourceData, meta interface{}, code string) error {
	client := meta.(*connectivity.AliyunClient)
	name := d.Id()

	// Step 1: obtain the OSS upload credential for the staging code.
	uploadAction := "GetRoutineStagingCodeUploadInfo"
	uploadRequest := map[string]interface{}{
		"Name": name,
	}
	if v, ok := d.GetOk("code_description"); ok {
		uploadRequest["CodeDescription"] = v
	}
	var uploadResponse map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		uploadResponse, err = client.RpcPost("ESA", "2024-09-10", uploadAction, nil, uploadRequest, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "LockFailed"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(uploadAction, uploadResponse, uploadRequest)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, name, uploadAction, AlibabaCloudSdkGoERROR)
	}

	ossPostConfigRaw, ok := uploadResponse["OssPostConfig"].(map[string]interface{})
	if !ok {
		return WrapError(fmt.Errorf("%s did not return a valid OssPostConfig for routine %q", uploadAction, name))
	}

	// Step 2: upload the code content to OSS with the returned form credential.
	if err := esaRoutinePostCodeToOss(ossPostConfigRaw, code); err != nil {
		return WrapError(err)
	}

	// Step 3: commit the staged code into a formal code version.
	commitAction := "CommitRoutineStagingCode"
	commitRequest := map[string]interface{}{
		"Name": name,
	}
	if v, ok := d.GetOk("code_description"); ok {
		commitRequest["CodeDescription"] = v
	}
	if v, ok := d.GetOk("deploy_env"); ok {
		commitRequest["DeployEnv"] = v
	}
	var commitResponse map[string]interface{}
	wait = incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		commitResponse, err = client.RpcPost("ESA", "2024-09-10", commitAction, nil, commitRequest, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "LockFailed"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(commitAction, commitResponse, commitRequest)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, name, commitAction, AlibabaCloudSdkGoERROR)
	}

	if v, ok := commitResponse["CodeVersion"]; ok {
		d.Set("latest_code_version", fmt.Sprint(v))
	}

	return nil
}

// esaRoutinePostCodeToOss uploads the routine code to OSS using the browser-style
// form POST credential returned by GetRoutineStagingCodeUploadInfo. The "file"
// part MUST be written last, as required by OSS PostObject.
func esaRoutinePostCodeToOss(ossPostConfig map[string]interface{}, code string) error {
	uploadURL, ok := ossPostConfig["Url"].(string)
	if !ok || uploadURL == "" {
		return fmt.Errorf("OssPostConfig is missing the upload Url")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Write all credential form fields except the endpoint Url. The security
	// token field must be sent as the OSS form field name.
	for k, v := range ossPostConfig {
		if k == "Url" {
			continue
		}
		fieldName := k
		if k == "XOssSecurityToken" {
			fieldName = "x-oss-security-token"
		}
		value := fmt.Sprint(v)
		if fieldName == "x-oss-security-token" && value == "" {
			continue
		}
		if err := writer.WriteField(fieldName, value); err != nil {
			return fmt.Errorf("writing OSS form field %q: %w", fieldName, err)
		}
	}
	if err := writer.WriteField("success_action_status", "200"); err != nil {
		return fmt.Errorf("writing OSS form field success_action_status: %w", err)
	}

	// The file part must be the last field in the form.
	part, err := writer.CreateFormFile("file", "index.js")
	if err != nil {
		return fmt.Errorf("creating OSS file form part: %w", err)
	}
	if _, err := part.Write([]byte(code)); err != nil {
		return fmt.Errorf("writing OSS file form part: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("closing OSS multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL, body)
	if err != nil {
		return fmt.Errorf("building OSS upload request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	httpClient := &http.Client{Timeout: 60 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("uploading routine code to OSS: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("OSS upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
