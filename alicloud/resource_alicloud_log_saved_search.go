package alicloud

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudLogSavedSearch() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogSavedSearchCreate,
		Read:   resourceAlicloudLogSavedSearchRead,
		Update: resourceAlicloudLogSavedSearchUpdate,
		Delete: resourceAlicloudLogSavedSearchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"search_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"search_query": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudLogSavedSearchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *sls.Client
	projectName := d.Get("project_name").(string)
	logstoreName := d.Get("logstore_name").(string)
	searchName := d.Get("search_name").(string)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	if err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			project, _ := sls.NewLogProject(projectName, slsClient.Endpoint, slsClient.AccessKeyID, slsClient.AccessKeySecret)
			project, _ = project.WithToken(slsClient.SecurityToken)
			_, _ = sls.NewLogStore(logstoreName, project)
			return nil, slsClient.CreateSavedSearch(projectName, &sls.SavedSearch{
				SavedSearchName: d.Get("search_name").(string),
				SearchQuery:     d.Get("search_query").(string),
				Logstore:        d.Get("logstore_name").(string),
				Topic:           d.Get("topic").(string),
				DisplayName:     d.Get("display_name").(string),
			})
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateSavedSearch", raw, requestInfo, map[string]string{
			"project_name":  projectName,
			"logstore_name": logstoreName,
			"search_name":   searchName,
		})
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_saved_search", "CreateSavedSearch", AliyunLogGoSdkERROR)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", projectName, COLON_SEPARATED, logstoreName, COLON_SEPARATED, searchName))
	return resourceAlicloudLogSavedSearchRead(d, meta)
}

func resourceAlicloudLogSavedSearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	object, err := logService.DescribeLogSavedSearch(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("project_name", parts[0])
	d.Set("search_name", object.SavedSearchName)
	d.Set("search_query", object.SearchQuery)
	d.Set("logstore_name", object.Logstore)
	d.Set("topic", object.Topic)
	d.Set("display_name", object.DisplayName)
	return nil
}

func resourceAlicloudLogSavedSearchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.UpdateSavedSearch(d.Get("project_name").(string), &sls.SavedSearch{
				SavedSearchName: d.Get("search_name").(string),
				SearchQuery:     d.Get("search_query").(string),
				Logstore:        d.Get("logstore_name").(string),
				Topic:           d.Get("topic").(string),
				DisplayName:     d.Get("display_name").(string),
			})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateSavedSearch", AliyunLogGoSdkERROR)
	}
	return resourceAlicloudLogSavedSearchRead(d, meta)

}

func resourceAlicloudLogSavedSearchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	var requestInfo *sls.Client
	request := map[string]string{
		"projectName": d.Get("project_name").(string),
		"searchName":  d.Get("search_name").(string),
	}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteSavedSearch(request["projectName"], request["searchName"])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout, "RequestTimeout"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteSavedSearch", raw, requestInfo, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ProjectNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteSavedSearch", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogSavedSearch(d.Id(), Deleted, DefaultTimeout))
}
