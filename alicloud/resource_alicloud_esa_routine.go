package alicloud

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
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
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			// Derive a content fingerprint from the local code file so that changing
			// the file content triggers a new staging upload + commit on the next apply.
			if filename, ok := diff.GetOk("filename"); ok && filename.(string) != "" {
				content, err := loadFileContent(filename.(string))
				if err != nil {
					return WrapError(err)
				}
				sum := sha256.Sum256(content)
				checksum := base64.StdEncoding.EncodeToString(sum[:])
				oldChecksum, _ := diff.GetChange("code_checksum")
				if oldChecksum.(string) != checksum {
					if err := diff.SetNew("code_checksum", checksum); err != nil {
						return WrapError(err)
					}
					// A code change produces a new committed version during apply.
					// Mark latest_code_version as known-after-apply so resources that
					// reference it (e.g. alicloud_esa_routine_code_deployment) are
					// re-planned within the same apply and the plan converges.
					if err := diff.SetNewComputed("latest_code_version"); err != nil {
						return WrapError(err)
					}
				}
			}
			return nil
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
			"filename": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code_checksum": {
				Type:     schema.TypeString,
				Computed: true,
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

	if filename, ok := d.GetOk("filename"); ok && filename.(string) != "" {
		codeVersion, err := esaRoutineUploadAndCommitCode(client, d.Id(), filename.(string), d.Get("code_description").(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return WrapError(err)
		}
		d.Set("latest_code_version", codeVersion)
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

	return nil
}

func resourceAliCloudEsaRoutineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// The routine metadata (name/description) is immutable (no UpdateRoutine API);
	// those attributes are marked ForceNew. The only in-place update is the code
	// lifecycle: a changed local file (or code_checksum) uploads a new staging
	// version and commits it into a new immutable code version.
	if d.HasChange("filename") || d.HasChange("code_checksum") {
		if filename, ok := d.GetOk("filename"); ok && filename.(string) != "" {
			codeVersion, err := esaRoutineUploadAndCommitCode(client, d.Id(), filename.(string), d.Get("code_description").(string), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return WrapError(err)
			}
			d.Set("latest_code_version", codeVersion)
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

// esaRoutineUploadAndCommitCode runs the three-step ESA Routine staging code flow:
//  1. GetRoutineStagingCodeUploadInfo -> obtain the OSS PostObject credentials
//  2. Upload the local code file to OSS with a multipart/form-data POST
//  3. CommitRoutineStagingCode -> generate an immutable code version
//
// It returns the committed code version.
func esaRoutineUploadAndCommitCode(client *connectivity.AliyunClient, name, filename, codeDescription string, timeout time.Duration) (string, error) {
	content, err := loadFileContent(filename)
	if err != nil {
		return "", WrapError(err)
	}

	// Step 1: obtain the OSS upload credentials.
	uploadAction := "GetRoutineStagingCodeUploadInfo"
	uploadRequest := map[string]interface{}{
		"Name": name,
	}
	if codeDescription != "" {
		uploadRequest["CodeDescription"] = codeDescription
	}
	var uploadResponse map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(timeout, func() *resource.RetryError {
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
		return "", WrapErrorf(err, DefaultErrorMsg, name, uploadAction, AlibabaCloudSdkGoERROR)
	}

	ossPostConfigRaw, ok := uploadResponse["OssPostConfig"].(map[string]interface{})
	if !ok {
		return "", WrapErrorf(Error("OssPostConfig is missing in %s response", uploadAction), DefaultErrorMsg, name, uploadAction, AlibabaCloudSdkGoERROR)
	}

	// Step 2: upload the code file to OSS via a form POST.
	if err := esaRoutinePostCodeToOSS(ossPostConfigRaw, content); err != nil {
		return "", WrapError(err)
	}

	// Step 3: commit the staging code, producing an immutable code version.
	commitAction := "CommitRoutineStagingCode"
	commitRequest := map[string]interface{}{
		"Name": name,
	}
	if codeDescription != "" {
		commitRequest["CodeDescription"] = codeDescription
	}
	var commitResponse map[string]interface{}
	wait = incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(timeout, func() *resource.RetryError {
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
		return "", WrapErrorf(err, DefaultErrorMsg, name, commitAction, AlibabaCloudSdkGoERROR)
	}

	return fmt.Sprint(commitResponse["CodeVersion"]), nil
}

// esaRoutinePostCodeToOSS performs the OSS PostObject multipart upload using the
// credentials returned by GetRoutineStagingCodeUploadInfo. All form fields must be
// written before the file part per the OSS PostObject contract.
func esaRoutinePostCodeToOSS(ossPostConfig map[string]interface{}, content []byte) error {
	uploadURL := fmt.Sprint(ossPostConfig["Url"])
	if uploadURL == "" || uploadURL == "<nil>" {
		return Error("OssPostConfig.Url is empty")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Map the OSS PostObject form fields. "key" must match the object key the
	// backend expects; the remaining fields carry the signed policy.
	formFields := map[string]string{
		"key":               fmt.Sprint(ossPostConfig["key"]),
		"OSSAccessKeyId":    fmt.Sprint(ossPostConfig["OSSAccessKeyId"]),
		"policy":            fmt.Sprint(ossPostConfig["policy"]),
		"Signature":         fmt.Sprint(ossPostConfig["Signature"]),
		"callback":          fmt.Sprint(ossPostConfig["callback"]),
		"x:codeDescription": fmt.Sprint(ossPostConfig["x:codeDescription"]),
	}
	if token, ok := ossPostConfig["XOssSecurityToken"]; ok && fmt.Sprint(token) != "" && fmt.Sprint(token) != "<nil>" {
		formFields["x-oss-security-token"] = fmt.Sprint(token)
	}
	for field, value := range formFields {
		if value == "" || value == "<nil>" {
			continue
		}
		if err := writer.WriteField(field, value); err != nil {
			return WrapError(err)
		}
	}

	part, err := writer.CreateFormFile("file", "index.js")
	if err != nil {
		return WrapError(err)
	}
	if _, err := part.Write(content); err != nil {
		return WrapError(err)
	}
	if err := writer.Close(); err != nil {
		return WrapError(err)
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL, body)
	if err != nil {
		return WrapError(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	httpClient := &http.Client{Timeout: 5 * time.Minute}
	resp, err := httpClient.Do(req)
	if err != nil {
		return WrapError(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return Error("failed to upload ESA Routine code to OSS, status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
