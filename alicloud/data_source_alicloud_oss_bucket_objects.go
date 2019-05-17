package alicloud

import (
	"log"
	"regexp"
	"strings"
	"time"

	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudOssBucketObjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOssBucketObjectsRead,

		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_mode": {
				Type:     schema.TypeString,
				Default:  "Default",
				Optional: true,
				ValidateFunc: validateOssBucketObjectQueryMode,
			},

			// Computed values
			"objects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cache_control": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_disposition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_encoding": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_md5": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expires": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_side_encryption": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sse_kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"etag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modification_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_latest": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"delete_marker": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOssBucketObjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	bucketName := d.Get("bucket_name").(string)
	
	queryMode  := d.Get("query_mode").(string)
	if queryMode != "Default" {
		return dataSourceAlicloudOssBucketObjectVersionsRead(d, meta)
	}

	// List bucket objects
	var initialOptions []oss.Option
	if v, ok := d.GetOk("key_prefix"); ok && v.(string) != "" {
		keyPrefix := v.(string)
		initialOptions = append(initialOptions, oss.Prefix(keyPrefix))
	}

	var allObjects []oss.ObjectProperties
	nextMarker := ""
	for {
		var options []oss.Option
		options = append(options, initialOptions...)
		if nextMarker != "" {
			options = append(options, oss.Marker(nextMarker))
		}

		raw, err := client.WithOssBucketByName(bucketName, func(bucket *oss.Bucket) (interface{}, error) {
			return bucket.ListObjects(options...)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(oss.ListObjectsResult)

		if resp.Objects == nil || len(resp.Objects) < 1 {
			break
		}

		allObjects = append(allObjects, resp.Objects...)

		nextMarker = resp.NextMarker
		if nextMarker == "" {
			break
		}
	}

	var filteredObjectsTemp []oss.ObjectProperties
	keyRegex, ok := d.GetOk("key_regex")
	if ok && keyRegex.(string) != "" {
		var r *regexp.Regexp
		if keyRegex != "" {
			r = regexp.MustCompile(keyRegex.(string))
		}
		for _, object := range allObjects {
			if r != nil && !r.MatchString(object.Key) {
				continue
			}
			filteredObjectsTemp = append(filteredObjectsTemp, object)
		}
	} else {
		filteredObjectsTemp = allObjects
	}

	return bucketObjectsDescriptionAttributes(d, bucketName, filteredObjectsTemp, meta)
}

func bucketObjectsDescriptionAttributes(d *schema.ResourceData, bucketName string, objects []oss.ObjectProperties, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var ids []string
	var s []map[string]interface{}
	for _, object := range objects {
		mapping := map[string]interface{}{
			"key":                    object.Key,
			"etag":                   strings.Trim(object.ETag, `"`),
			"storage_class":          object.StorageClass,
			"last_modification_time": object.LastModified.Format(time.RFC3339),
		}

		// Add metadata information
		raw, err := client.WithOssBucketByName(bucketName, func(bucket *oss.Bucket) (interface{}, error) {
			return bucket.GetObjectDetailedMeta(object.Key)
		})
		if err != nil {
			log.Printf("[ERROR] Unable to get metadata for the object %s: %v", object.Key, err)
		} else {
			objectHeader, _ := raw.(http.Header)
			mapping["content_type"] = objectHeader.Get("Content-Type")
			mapping["content_length"] = objectHeader.Get("Content-Length")
			mapping["cache_control"] = objectHeader.Get("Cache-Control")
			mapping["content_disposition"] = objectHeader.Get("Content-Disposition")
			mapping["content_encoding"] = objectHeader.Get("Content-Encoding")
			mapping["content_md5"] = objectHeader.Get("Content-Md5")
			mapping["expires"] = objectHeader.Get("Expires")
			mapping["server_side_encryption"] = objectHeader.Get(oss.HTTPHeaderOssServerSideEncryption)
			mapping["sse_kms_key_id"] = objectHeader.Get(oss.HTTPHeaderOssServerSideEncryptionKeyID)
		}
		// Add ACL information
		raw, err = client.WithOssBucketByName(bucketName, func(bucket *oss.Bucket) (interface{}, error) {
			return bucket.GetObjectACL(object.Key)
		})
		if err != nil {
			log.Printf("[ERROR] Unable to get ACL for the object %s: %v", object.Key, err)
		} else {
			objectACL, _ := raw.(oss.GetObjectACLResult)
			mapping["acl"] = objectACL.ACL
		}

		ids = append(ids, object.Key)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("objects", s); err != nil {
		return err
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAlicloudOssBucketObjectVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	bucketName := d.Get("bucket_name").(string)
	queryMode  := d.Get("query_mode").(string)

	// List bucket objectVersions
	var initialOptions []oss.Option
	if v, ok := d.GetOk("key_prefix"); ok && v.(string) != "" {
		keyPrefix := v.(string)
		initialOptions = append(initialOptions, oss.Prefix(keyPrefix))
	}

	var allObjectVersions []oss.ObjectVersionProperties
	var allObjectDeleteds []oss.ObjectDeleteMarkerProperties

	var flagVersions = false
	var flagDeleted  = false

	if queryMode == "All" || queryMode == "VersionOnly" {
		flagVersions = true
	}

	if queryMode == "All" || queryMode == "DeleteMarkerOnly" {
		flagDeleted = true
	}

	nextKeyMarker := ""
	for {
		var options []oss.Option
		options = append(options, initialOptions...)
		if nextKeyMarker != "" {
			options = append(options, oss.KeyMarker(nextKeyMarker))
		}

		raw, err := client.WithOssBucketByName(bucketName, func(bucket *oss.Bucket) (interface{}, error) {
			return bucket.ListObjectVersions(options...)
		})
		if err != nil {
			return err
		}
		
		resp, _ := raw.(oss.ListObjectVersionsResult)

		if flagVersions && resp.ObjectVersions != nil  {
			allObjectVersions = append(allObjectVersions, resp.ObjectVersions...)
		}

		if flagDeleted && resp.ObjectDeleteMarkers != nil {
			allObjectDeleteds = append(allObjectDeleteds, resp.ObjectDeleteMarkers...)
		}

		nextKeyMarker = resp.NextKeyMarker
		if nextKeyMarker == "" {
			break
		}
	}


	//Filter
	keyRegex, ok := d.GetOk("key_regex")
	if ok && keyRegex.(string) != "" {
		var r *regexp.Regexp
		r = regexp.MustCompile(keyRegex.(string))
		for i := len(allObjectVersions) - 1; i >= 0; i-- {
			if r != nil && !r.MatchString(allObjectVersions[i].Key) {
				allObjectVersions = append(allObjectVersions[:i], allObjectVersions[i+1:]...)
			}
		}

		for i := len(allObjectDeleteds) - 1; i >= 0; i-- {
			if r != nil && !r.MatchString(allObjectDeleteds[i].Key) {
				allObjectDeleteds = append(allObjectDeleteds[:i], allObjectDeleteds[i+1:]...)
			}
		}
	} 

	return bucketObjectsDescriptionAttributesVersions(d, bucketName, allObjectVersions, allObjectDeleteds, meta)
}

func bucketObjectsDescriptionAttributesVersions(d *schema.ResourceData, bucketName string, objectVersions []oss.ObjectVersionProperties, objectDeleteds []oss.ObjectDeleteMarkerProperties, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var ids []string
	var s []map[string]interface{}
	for _, object := range objectVersions {
		mapping := map[string]interface{}{
			"key":                    object.Key,
			"etag":                   strings.Trim(object.ETag, `"`),
			"storage_class":          object.StorageClass,
			"last_modification_time": object.LastModified.Format(time.RFC3339),
			"version_id":             object.VersionId,
			"is_latest":             object.IsLatest,
		}

		var options []oss.Option
		options=append(options, oss.VersionId(object.VersionId))

		// Add metadata information
		raw, err := client.WithOssBucketByName(bucketName, func(bucket *oss.Bucket) (interface{}, error) {
			return bucket.GetObjectDetailedMeta(object.Key, options...)
		})
		if err != nil {
			log.Printf("[ERROR] Unable to get metadata for the object %s: %v", object.Key, err)
		} else {
			objectHeader, _ := raw.(http.Header)
			mapping["content_type"] = objectHeader.Get("Content-Type")
			mapping["content_length"] = objectHeader.Get("Content-Length")
			mapping["cache_control"] = objectHeader.Get("Cache-Control")
			mapping["content_disposition"] = objectHeader.Get("Content-Disposition")
			mapping["content_encoding"] = objectHeader.Get("Content-Encoding")
			mapping["content_md5"] = objectHeader.Get("Content-Md5")
			mapping["expires"] = objectHeader.Get("Expires")
			mapping["server_side_encryption"] = objectHeader.Get(oss.HTTPHeaderOssServerSideEncryption)
			mapping["sse_kms_key_id"] = objectHeader.Get(oss.HTTPHeaderOssServerSideEncryptionKeyID)
		}
		// Add ACL information
		raw, err = client.WithOssBucketByName(bucketName, func(bucket *oss.Bucket) (interface{}, error) {
			return bucket.GetObjectACL(object.Key, options...)
		})
		if err != nil {
			log.Printf("[ERROR] Unable to get ACL for the object %s: %v", object.Key, err)
		} else {
			objectACL, _ := raw.(oss.GetObjectACLResult)
			mapping["acl"] = objectACL.ACL
		}

		ids = append(ids, object.Key + "_" + object.VersionId)
		s = append(s, mapping)
	}

	for _, object := range objectDeleteds {
		mapping := map[string]interface{}{
			"key":                    object.Key,
			"last_modification_time": object.LastModified.Format(time.RFC3339),
			"version_id":             object.VersionId,
			"is_latest":             object.IsLatest,
			"delete_marker":          true,
		}

		ids = append(ids, object.Key + "_" + object.VersionId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("objects", s); err != nil {
		return err
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}