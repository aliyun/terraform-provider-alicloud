package alicloud

import (
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"time"
)

func resourceAlicloudLogStoreIndex() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogStoreIndexCreate,
		Read:   resourceAlicloudLogStoreIndexRead,
		Update: resourceAlicloudLogStoreIndexUpdate,
		Delete: resourceAlicloudLogStoreIndexDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"full_text": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"case_sensitive": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"include_chinese": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"token": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
			},
			//field search
			"field_search": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  LongType,
							ValidateFunc: validateAllowedStringValue([]string{string(TextType), string(LongType),
								string(DoubleType), string(JsonType)}),
						},
						"alias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"case_sensitive": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"include_chinese": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"token": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable_analytics": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
				MinItems: 1,
			},
		},
	}
}

func resourceAlicloudLogStoreIndexCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}

	_, fullOk := d.GetOk("full_text")
	_, fieldOk := d.GetOk("field_search")
	if !fullOk && !fieldOk {
		return WrapError(Error("At least one of the 'full_text' and 'field_search' should be specified."))
	}

	project := d.Get("project").(string)
	store, err := logService.DescribeLogStore(fmt.Sprintf("%s%s%s", project, COLON_SEPARATED, d.Get("logstore").(string)))
	if err != nil {
		return WrapError(err)
	}

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := store.GetIndex()
		if err != nil {
			if IsExceptedError(err, LogClientTimeout) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "GetIndex", AliyunLogGoSdkERROR))
			}
			if !IsExceptedErrors(err, []string{IndexConfigNotExist}) {
				return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "GetIndex", AliyunLogGoSdkERROR))
			}
		}
		if raw != nil {
			return resource.NonRetryableError(WrapError(Error("There is aleady existing an index in the store %s. Please import it using id '%s%s%s'.",
				store.Name, project, COLON_SEPARATED, store.Name)))
		}
		addDebug("GetIndex", raw)
		return nil
	}); err != nil {
		return err
	}

	var index sls.Index
	if fullOk {
		index.Line = buildIndexLine(d)
	}
	if fieldOk {
		index.Keys = buildIndexKeys(d)
	}

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {

		if e := store.CreateIndex(index); e != nil {
			if IsExceptedErrors(e, []string{InternalServerError, LogClientTimeout}) {
				return resource.RetryableError(e)
			}
			return resource.NonRetryableError(e)
		}
		addDebug("CreateIndex", nil)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store_index", "CreateIndex", AliyunLogGoSdkERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", project, COLON_SEPARATED, store.Name))

	return resourceAlicloudLogStoreIndexRead(d, meta)
}

func resourceAlicloudLogStoreIndexRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	index, err := logService.DescribeLogStoreIndex(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if line := index.Line; line != nil {
		mapping := map[string]interface{}{
			"case_sensitive":  line.CaseSensitive,
			"include_chinese": line.Chn,
			"token":           strings.Join(line.Token, ""),
		}
		if err := d.Set("full_text", []map[string]interface{}{mapping}); err != nil {
			return WrapError(err)
		}
	}
	if keys := index.Keys; keys != nil {
		var keySet []map[string]interface{}
		for k, v := range keys {
			mapping := map[string]interface{}{
				"name":             k,
				"type":             v.Type,
				"alias":            v.Alias,
				"case_sensitive":   v.CaseSensitive,
				"include_chinese":  v.Chn,
				"token":            strings.Join(v.Token, ""),
				"enable_analytics": v.DocValue,
			}
			keySet = append(keySet, mapping)
		}
		if err := d.Set("field_search", keySet); err != nil {
			return WrapError(err)
		}
	}

	d.Set("project", parts[0])
	d.Set("logstore", parts[1])

	return nil
}

func resourceAlicloudLogStoreIndexUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	logService := LogService{client}
	index, err := logService.DescribeLogStoreIndex(d.Id())
	if err != nil {
		return WrapError(err)
	}
	update := false
	if d.HasChange("full_text") {
		index.Line = buildIndexLine(d)
		update = true
	}
	if d.HasChange("field_search") {
		index.Keys = buildIndexKeys(d)
		update = true
	}

	if update {
		if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				return nil, slsClient.UpdateIndex(parts[0], parts[1], *index)
			})
			if err != nil {
				if IsExceptedError(err, LogClientTimeout) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug("UpdateIndex", raw)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateIndex", AliyunLogGoSdkERROR)
		}
	}

	return resourceAlicloudLogStoreIndexRead(d, meta)
}

func resourceAlicloudLogStoreIndexDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if _, err := logService.DescribeLogStoreIndex(d.Id()); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteIndex(parts[0], parts[1])
		})
		if err != nil {
			if IsExceptedError(err, LogClientTimeout) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteIndex", raw)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteIndex", AliyunLogGoSdkERROR)
	}
	return nil
}

func buildIndexLine(d *schema.ResourceData) *sls.IndexLine {
	if fullText, ok := d.GetOk("full_text"); ok {
		value := fullText.(*schema.Set).List()[0].(map[string]interface{})
		return &sls.IndexLine{
			CaseSensitive: value["case_sensitive"].(bool),
			Chn:           value["include_chinese"].(bool),
			Token:         strings.Split(value["token"].(string), ""),
		}
	}
	return nil
}

func buildIndexKeys(d *schema.ResourceData) map[string]sls.IndexKey {
	keys := make(map[string]sls.IndexKey)
	if field, ok := d.GetOk("field_search"); ok {
		for _, f := range field.(*schema.Set).List() {
			v := f.(map[string]interface{})
			keys[v["name"].(string)] = sls.IndexKey{
				Type:          v["type"].(string),
				Alias:         v["alias"].(string),
				DocValue:      v["enable_analytics"].(bool),
				Token:         strings.Split(v["token"].(string), ""),
				CaseSensitive: v["case_sensitive"].(bool),
				Chn:           v["include_chinese"].(bool),
			}
		}
	}
	return keys
}
