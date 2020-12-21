package alicloud

import (
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceAlicloudLogProjectTags() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogProjectTagsCreate,
		Read:   resourceAlicloudLogProjecTagstRead,
		Delete: resourceAlicloudLogProjectTagsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: schema.TypeString,
				ForceNew: true,
			},

		},
	}
}


func resourceAlicloudLogProjectTagsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *sls.Client
	projectName := d.Get("project_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	slsTags := []sls.ResourceTag{}
	for key, value := range(tags) {
		tag := sls.ResourceTag{key,value.(string)}
		slsTags = append(slsTags, tag)
	}
	projectTags := sls.NewProjectTags(projectName, slsTags)
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.TagResources(projectName, projectTags)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout, "ProjectNotExist"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateProjectTags", raw, requestInfo, tags)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_project_tags", "CreateProjectTags", AliyunLogGoSdkERROR)
	}
	var keyList = ""
	for k,_ := range(tags) {
		keyList += k
	}
	d.SetId(fmt.Sprintf("%s%s%s", projectName, COLON_SEPARATED, keyList))
	return resourceAlicloudLogProjecTagstRead(d, meta)
}


func resourceAlicloudLogProjecTagstRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	projectTags, err := logService.DescribeLogProjectTags(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("project_name", projectTags[0].ResourceID)
	tags := map[string]string{}
	for _, tag := range projectTags {
		tags[tag.TagKey] = tag.TagValue
	}
	d.Set("tags",tags)
	return nil
}

func resourceAlicloudLogProjectTagsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	var requestInfo *sls.Client
	projectName := d.Get("project_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	slsTags := []string{}
	for key, _ := range(tags) {
		slsTags = append(slsTags, key)
	}
	projectUnTags := sls.NewProjectUnTags(projectName, slsTags)
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.UnTagResources(projectName, projectUnTags)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout, "RequestTimeout"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteProjectTags", DefaultErrorMsg, requestInfo, slsTags)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteProjectTagsError", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogProjectTags(d.Id(), Deleted, DefaultTimeout))

}
