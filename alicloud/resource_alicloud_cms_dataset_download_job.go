// Package alicloud. This file is hand-written from the Cms 2024-03-30 CloudSpec DatasetDownloadJob definition.
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsDatasetDownloadJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsDatasetDownloadJobCreate,
		Read:   resourceAliCloudCmsDatasetDownloadJobRead,
		Delete: resourceAliCloudCmsDatasetDownloadJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"compression_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataset_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"job_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsDatasetDownloadJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	workspace := d.Get("workspace")
	datasetName := d.Get("dataset_name")
	action := fmt.Sprintf("/workspace/%s/dataset/%s/downloadjob", workspace, datasetName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["jobName"] = d.Get("job_name")
	request["query"] = d.Get("query")
	if v, ok := d.GetOk("compression_type"); ok {
		request["compressionType"] = v
	}
	if v, ok := d.GetOk("content_type"); ok {
		request["contentType"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_dataset_download_job", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", workspace, datasetName, d.Get("job_name")))

	return resourceAliCloudCmsDatasetDownloadJobRead(d, meta)
}

func resourceAliCloudCmsDatasetDownloadJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsDatasetDownloadJob(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_dataset_download_job DescribeCmsDatasetDownloadJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("compression_type", objectRaw["compressionType"])
	d.Set("content_type", objectRaw["contentType"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("dataset_name", objectRaw["datasetName"])
	d.Set("job_name", objectRaw["jobName"])
	d.Set("query", objectRaw["query"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("update_time", objectRaw["updateTime"])
	d.Set("workspace", objectRaw["workspace"])

	return nil
}

func resourceAliCloudCmsDatasetDownloadJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 3, len(parts)))
	}
	workspace := parts[0]
	datasetName := parts[1]
	jobName := parts[2]
	action := fmt.Sprintf("/workspace/%s/dataset/%s/downloadjob/%s", workspace, datasetName, jobName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"WorkspaceNotExist", "DatasetNotExist", "DatasetDownloadJobNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
