package alicloud

import (
	"fmt"
	"time"

	"strings"

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
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention_period": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validateIntegerInRange(1, 3650),
			},
			"shard_count": &schema.Schema{
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
			"shards": &schema.Schema{
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
		},
	}
}

func resourceAlicloudLogStoreCreate(d *schema.ResourceData, meta interface{}) error {
	if err := meta.(*AliyunClient).logconn.CreateLogStore(d.Get("project").(string), d.Get("name").(string),
		d.Get("retention_period").(int), d.Get("shard_count").(int)); err != nil {
		return fmt.Errorf("CreateLogStore got an error: %#v.", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))

	return resourceAlicloudLogStoreUpdate(d, meta)
}

func resourceAlicloudLogStoreRead(d *schema.ResourceData, meta interface{}) error {
	split := strings.Split(d.Id(), COLON_SEPARATED)

	store, err := meta.(*AliyunClient).DescribeLogStore(split[0], split[1])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project", split[0])
	d.Set("name", store.Name)
	d.Set("retention_period", store.TTL)
	d.Set("shard_count", store.ShardCount)
	shards, err := store.ListShards()
	if err != nil {
		return fmt.Errorf("ListShards got an error: %#v.", err)
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

	return nil
}

func resourceAlicloudLogStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	if d.IsNewResource() {
		return resourceAlicloudLogStoreRead(d, meta)
	}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	d.Partial(true)

	update := false
	if d.HasChange("retention_period") {
		update = true
		d.SetPartial("retention_period")
	}

	if update {
		store, err := client.DescribeLogStore(split[0], split[1])
		if err != nil {
			return err
		}
		if err = client.logconn.UpdateLogStore(split[0], split[1], d.Get("retention_period").(int), store.ShardCount); err != nil {
			return fmt.Errorf("UpdateLogStore %s got an error: %#v.", split[1], err)
		}
	}
	d.Partial(false)

	return resourceAlicloudLogStoreRead(d, meta)
}

func resourceAlicloudLogStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	split := strings.Split(d.Id(), COLON_SEPARATED)

	project, err := client.DescribeLogProject(split[0])
	if err != nil {
		return err
	}
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		if err := project.DeleteLogStore(split[1]); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting log store %s got an error: %#v", split[1], err))
		}

		store, err := client.DescribeLogStore(split[0], split[1])
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if store.Name == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting log store %s got an error: %#v", split[1], err))
	})
}
