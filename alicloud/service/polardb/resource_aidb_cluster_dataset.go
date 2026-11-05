package polardb

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAIDBClusterDataset() *schema.Resource {
	return &schema.Resource{
		Create: resourceAIDBClusterDatasetCreate,
		Read:   resourceAIDBClusterDatasetRead,
		Delete: resourceAIDBClusterDatasetDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), ":")
				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("invalid import ID %q for alicloud_polardb_aidb_cluster_dataset, expected <db_cluster_id>:<dataset_id>", d.Id())
				}
				d.Set("db_cluster_id", parts[0])
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dataset_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dataset_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"train_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"import_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"relative_db_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAIDBClusterDatasetCreate(d *schema.ResourceData, meta interface{}) error {
	service := NewService(meta.(*connectivity.AliyunClient))
	request := map[string]interface{}{
		"DBClusterId": d.Get("db_cluster_id"),
		"DatasetName": d.Get("dataset_name"),
	}
	if value, ok := d.GetOk("dataset_type"); ok {
		request["DatasetType"] = value
	}
	if value, ok := d.GetOk("train_mode"); ok {
		request["TrainMode"] = value
	}
	if value, ok := d.GetOk("import_mode"); ok {
		request["ImportMode"] = value
	}
	response, err := service.rpcPost("CreateAIDBClusterDataset", request, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_aidb_cluster_dataset", "CreateAIDBClusterDataset", errs.AlibabaCloudSdkGoERROR)
	}
	datasetId := fmt.Sprint(response["DatasetId"])
	if datasetId == "" || datasetId == "<nil>" {
		return errs.Error("creating alicloud_polardb_aidb_cluster_dataset: empty DatasetId in response: %v", response)
	}
	d.SetId(datasetId)

	return resourceAIDBClusterDatasetRead(d, meta)
}

func resourceAIDBClusterDatasetRead(d *schema.ResourceData, meta interface{}) error {
	service := NewService(meta.(*connectivity.AliyunClient))
	object, err := service.DescribeAIDBClusterDataset(d.Get("db_cluster_id").(string), d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_aidb_cluster_dataset not found, removing from state: %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("dataset_type", object["DatasetType"])
	d.Set("train_mode", object["TrainMode"])
	d.Set("relative_db_cluster_id", object["RelativeDBClusterId"])
	d.Set("path", object["Path"])
	d.Set("data_service_id", object["DataServiceId"])
	d.Set("file_name", object["FileName"])
	d.Set("bucket_name", object["BucketName"])
	d.Set("storage_type", object["StorageType"])
	d.Set("capacity", object["Capacity"])
	d.Set("last_modified", object["LastModified"])
	d.Set("creation_time", object["CreationTime"])
	if service.client.RegionId != "" {
		d.Set("region_id", service.client.RegionId)
	}
	return nil
}

func resourceAIDBClusterDatasetDelete(d *schema.ResourceData, meta interface{}) error {
	service := NewService(meta.(*connectivity.AliyunClient))
	request := map[string]interface{}{
		"DBClusterId": d.Get("db_cluster_id"),
		"DatasetId":   d.Id(),
	}
	if service.client.RegionId != "" {
		request["RegionId"] = service.client.RegionId
	}
	if _, err := service.rpcPost("DeleteAIDBClusterDataset", request, d.Timeout(schema.TimeoutDelete)); err != nil && !errs.NotFoundError(err) {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteAIDBClusterDataset", errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("")
	return nil
}
