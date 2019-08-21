package alicloud

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudLogStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogStoreCreate,
		Read:   resourceAlicloudLogStoreRead,
		Update: resourceAlicloudLogStoreUpdate,
		Delete: resourceAlicloudLogStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validateIntegerInRange(1, 3650),
			},
			"shard_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == "" {
						return false
					}
					return true
				},
			},
			"shards": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"begin_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"auto_split": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"max_split_shard_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntegerInRange(0, 64),
			},
			"append_meta": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enable_web_tracking": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlicloudLogStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logstore := &sls.LogStore{
		Name:          d.Get("name").(string),
		TTL:           d.Get("retention_period").(int),
		ShardCount:    d.Get("shard_count").(int),
		WebTracking:   d.Get("enable_web_tracking").(bool),
		AutoSplit:     d.Get("auto_split").(bool),
		MaxSplitShard: d.Get("max_split_shard_count").(int),
		AppendMeta:    d.Get("append_meta").(bool),
	}
	var requestinfo *sls.Client
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {

		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestinfo = slsClient
			return nil, slsClient.CreateLogStoreV2(d.Get("project").(string), logstore)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InternalServerError, LogClientTimeout}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("CreateLogStoreV2", raw, requestinfo, map[string]interface{}{
				"project":  d.Get("project").(string),
				"logstore": logstore,
			})
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "CreateLogStoreV2", AliyunLogGoSdkERROR)
	}
	d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))

	return resourceAlicloudLogStoreUpdate(d, meta)
}

func resourceAlicloudLogStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := logService.DescribeLogStore(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project", parts[0])
	d.Set("name", object.Name)
	d.Set("retention_period", object.TTL)
	d.Set("shard_count", object.ShardCount)
	var shards []*sls.Shard
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		shards, err = object.ListShards()
		if err != nil {
			if IsExceptedError(err, InternalServerError) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("ListShards", shards)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "ListShards", AliyunLogGoSdkERROR)
	}
	var shardList []map[string]interface{}
	for _, s := range shards {
		mapping := map[string]interface{}{
			"id":        s.ShardID,
			"status":    s.Status,
			"begin_key": s.InclusiveBeginKey,
			"end_key":   s.ExclusiveBeginKey,
		}
		shardList = append(shardList, mapping)
	}
	d.Set("shards", shardList)
	d.Set("append_meta", object.AppendMeta)
	d.Set("auto_split", object.AutoSplit)
	d.Set("enable_web_tracking", object.WebTracking)
	d.Set("max_split_shard_count", object.MaxSplitShard)

	return nil
}

func resourceAlicloudLogStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}

	if d.IsNewResource() {
		return resourceAlicloudLogStoreRead(d, meta)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	if d.HasChange("retention_period") {
		update = true
		d.SetPartial("retention_period")
	}
	if d.HasChange("max_split_shard_count") {
		update = true
		d.SetPartial("max_split_shard_count")
	}
	if d.HasChange("enable_web_tracking") {
		update = true
		d.SetPartial("enable_web_tracking")
	}
	if d.HasChange("append_meta") {
		update = true
		d.SetPartial("append_meta")
	}
	if d.HasChange("auto_split") {
		update = true
		d.SetPartial("auto_split")
	}

	if update {
		store, err := logService.DescribeLogStore(d.Id())
		if err != nil {
			return WrapError(err)
		}
		store.MaxSplitShard = d.Get("max_split_shard_count").(int)
		store.TTL = d.Get("retention_period").(int)
		store.WebTracking = d.Get("enable_web_tracking").(bool)
		store.AppendMeta = d.Get("append_meta").(bool)
		store.AutoSplit = d.Get("auto_split").(bool)
		var requestInfo *sls.Client
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.UpdateLogStoreV2(parts[0], store)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateLogStoreV2", AliyunLogGoSdkERROR)
		}
		if debugOn() {
			addDebug("UpdateLogStoreV2", raw, requestInfo, map[string]interface{}{
				"project":  parts[0],
				"logstore": store,
			})
		}
	}
	d.Partial(false)

	return resourceAlicloudLogStoreRead(d, meta)
}

func resourceAlicloudLogStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	project, err := logService.DescribeLogProject(parts[0])
	if err != nil {
		return WrapError(err)
	}
	err = project.DeleteLogStore(parts[1])
	if err != nil {
		if IsExceptedErrors(err, []string{LogStoreNotExist}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteLogStore", AliyunLogGoSdkERROR)
	}
	addDebug("DeleteLogStore", nil)
	return WrapError(logService.WaitForLogStore(d.Id(), Deleted, DefaultTimeout))
}
