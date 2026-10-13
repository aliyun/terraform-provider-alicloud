package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudMaxComputePackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudMaxComputePackagesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"packages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"install_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudMaxComputePackagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	projectName := d.Get("project_name").(string)
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	action := fmt.Sprintf("/api/v1/projects/%s/packages", projectName)
	query := make(map[string]*string)
	query["maxItem"] = StringPointer(strconv.Itoa(PageSizeLarge))
	var response map[string]interface{}
	var err error
	objects := make([]map[string]interface{}, 0)
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("MaxCompute", "2022-01-04", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, projectName)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		resp, _ := jsonpath.Get("$.data.installedPackages[*]", response)
		if resp != nil {
			if list, ok := resp.([]interface{}); ok {
				for _, v := range list {
					item, ok := v.(map[string]interface{})
					if !ok {
						continue
					}
					name := fmt.Sprint(item["name"])
					if name == "" {
						continue
					}
					if nameRegex != nil && !nameRegex.MatchString(name) {
						continue
					}
					packageId := fmt.Sprintf("%s:%s", projectName, name)
					if len(idsMap) > 0 {
						if _, ok := idsMap[packageId]; !ok {
							continue
						}
					}
					objects = append(objects, item)
				}
			}
		}
		marker, _ := jsonpath.Get("$.data.marker", response)
		if token, ok := marker.(string); ok && token != "" {
			query["marker"] = StringPointer(token)
			continue
		}
		break
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	packages := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		name := fmt.Sprint(objectRaw["name"])
		mapping := map[string]interface{}{
			"project_name":   projectName,
			"package_name":   name,
			"source_project": objectRaw["sourceProject"],
			"status":         objectRaw["status"],
			"install_time":   objectRaw["installTime"],
		}
		ids = append(ids, fmt.Sprintf("%s:%s", projectName, name))
		names = append(names, name)
		packages = append(packages, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("packages", packages); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), packages)
	}
	return nil
}
