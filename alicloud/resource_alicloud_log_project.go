package alicloud

import (
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudLogProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogProjectCreate,
		Read:   resourceAlicloudLogProjectRead,
		Update: resourceAlicloudLogProjectUpdate,
		Delete: resourceAlicloudLogProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudLogProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *sls.Client
	request := map[string]string{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
	}
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.CreateProject(request["name"], request["description"])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateProject", raw, requestInfo, request)
		response, _ := raw.(*sls.LogProject)
		d.SetId(response.Name)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_project", "CreateProject", AliyunLogGoSdkERROR)
	}

	return resourceAlicloudLogProjectUpdate(d, meta)
}

func resourceAlicloudLogProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	object, err := logService.DescribeLogProject(d.Id())

	projectTags, err := logService.DescribeLogProjectTags(object.Name)
	tags := map[string]string{}
	for _, tag := range projectTags {
		tags[tag.TagKey] = tag.TagValue
	}
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("tags", tags)
	return nil
}

func buildTags(projectName string, tags map[string]interface{}) *sls.ResourceTags{
	slsTags := []sls.ResourceTag{}

	for key, value := range tags {
		tag := sls.ResourceTag{key, value.(string)}
		slsTags = append(slsTags, tag)
	}
	projectTags := sls.NewProjectTags(projectName, slsTags)
	return projectTags
}


func resourceAlicloudLogProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *sls.Client
	logService := LogService{client}
	projectName := d.Get("name").(string)
	request := map[string]string{
		"name":        projectName,
		"description": d.Get("description").(string),
	}
	if d.HasChange("description") {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.UpdateProject(request["name"], request["description"])
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateProject", AliyunLogGoSdkERROR)
		}
		addDebug("UpdateProject", raw, requestInfo, request)
	}
	if d.HasChange("tags") {
		projectTags, err := logService.DescribeLogProjectTags(projectName)
		if err != nil {
			return err
		}
		tags := d.Get("tags").(map[string]interface{})
		if len(tags) == 0 {
			slsTags := []string{}
			for _, value := range projectTags {
				slsTags = append(slsTags, value.TagKey)
			}
			projectUnTags := sls.NewProjectUnTags(projectName, slsTags)
			raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				requestInfo = slsClient
				return nil, slsClient.UnTagResources(projectName, projectUnTags)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeletaTags", AliyunLogGoSdkERROR)
			}
			addDebug("DeletaTags", raw, requestInfo, request)

		}else{
			slsTags := []string{}
			for _, value := range projectTags {
				slsTags = append(slsTags, value.TagKey)
			}
			projectUnTags := sls.NewProjectUnTags(projectName, slsTags)
			raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				requestInfo = slsClient
				return nil, slsClient.UnTagResources(projectName, projectUnTags)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeletaTags", AliyunLogGoSdkERROR)
			}
			addDebug("DeletaTags", raw, requestInfo, request)
			projectNewTags := buildTags(projectName, tags)
			raw, err = client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				requestInfo = slsClient
				return nil, slsClient.TagResources(projectName, projectNewTags)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateTags", AliyunLogGoSdkERROR)
			}
			addDebug("UpdateTags", raw, requestInfo, request)

		}
	}
	return resourceAlicloudLogProjectRead(d, meta)
}

func resourceAlicloudLogProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	var requestInfo *sls.Client
	request := map[string]string{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
	}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteProject(request["name"])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout, "RequestTimeout"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteProject", raw, requestInfo, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ProjectNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteProject", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogProject(d.Id(), Deleted, DefaultTimeout))
}
