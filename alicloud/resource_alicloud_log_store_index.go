package alicloud

import (
	"fmt"
	"strings"

	"time"

	"bytes"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"full_text": &schema.Schema{
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
				Set: func(v interface{}) int {
					return 1
				},
				MaxItems: 1,
			},
			//field search
			"field_search": &schema.Schema{
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
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["name"]))
					return hashcode.String(buf.String())
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
		return fmt.Errorf("At least one of the 'full_text' and 'field_search' should be specified.")
	}

	project := d.Get("project").(string)
	store, err := logService.DescribeLogStore(project, d.Get("logstore").(string))
	if err != nil {
		return fmt.Errorf("DescribeLogStore got an error: %#v.", err)
	}
	exist, err := store.GetIndex()
	if err != nil && !IsExceptedErrors(err, []string{IndexConfigNotExist}) {
		return fmt.Errorf("While Creating index, GetIndex got an error: %#v.", err)
	}

	if exist != nil {
		return fmt.Errorf("There is aleady existing an index in the store %s. Please import it using id '%s%s%s'.",
			store.Name, project, COLON_SEPARATED, store.Name)
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
			if IsExceptedErrors(err, []string{InternalServerError}) {
				return resource.RetryableError(fmt.Errorf("CreateLogStoreIndex timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("CreateLogStoreIndex got an error: %#v.", err))
		}
		return nil
	}); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", project, COLON_SEPARATED, store.Name))

	return resourceAlicloudLogStoreIndexRead(d, meta)
}

func resourceAlicloudLogStoreIndexRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)

	index, err := logService.DescribeLogStoreIndex(split[0], split[1])

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("GetIndex got an error: %#v.", err)
	}

	if line := index.Line; line != nil {
		mapping := map[string]interface{}{
			"case_sensitive":  line.CaseSensitive,
			"include_chinese": line.Chn,
			"token":           strings.Join(line.Token, ""),
		}
		if err := d.Set("full_text", []map[string]interface{}{mapping}); err != nil {
			return err
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
			return err
		}
	}

	d.Set("project", split[0])
	d.Set("logstore", split[1])

	return nil
}

func resourceAlicloudLogStoreIndexUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	split := strings.Split(d.Id(), COLON_SEPARATED)
	d.Partial(true)

	update := false
	var index sls.Index
	if d.HasChange("full_text") {
		index.Line = buildIndexLine(d)
		update = true
		d.SetPartial("full_text")
	}
	if d.HasChange("field_search") {
		index.Keys = buildIndexKeys(d)
		update = true
		d.SetPartial("field_search")
	}

	if update {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.UpdateIndex(split[0], split[1], index)
		})
		if err != nil {
			return fmt.Errorf("UpdateLogStoreIndex got an error: %#v.", err)
		}
	}
	d.Partial(false)

	return resourceAlicloudLogStoreIndexRead(d, meta)
}

func resourceAlicloudLogStoreIndexDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}

	split := strings.Split(d.Id(), COLON_SEPARATED)

	if _, err := logService.DescribeLogStoreIndex(split[0], split[1]); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return fmt.Errorf("While deleting index, GetIndex got an error: %#v.", err)
	}

	_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return nil, slsClient.DeleteIndex(split[0], split[1])
	})
	if err != nil {
		return fmt.Errorf("DeleteIndex got an error: %#v.", err)
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
