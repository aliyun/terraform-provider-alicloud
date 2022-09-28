package alicloud

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOtsSearchIndexes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOtsSearchIndexesRead,

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSInstanceName,
			},
			"table_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSTableName,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"indexes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"time_to_live": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sync_phase": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"current_sync_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"schema": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"row_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"reserved_read_cu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"metering_last_update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func parseSearchIndexDataSourceArgs(d *schema.ResourceData) *SearchIndexDataSourceArgs {
	args := &SearchIndexDataSourceArgs{
		instanceName: d.Get("instance_name").(string),
		tableName:    d.Get("table_name").(string),
	}

	if ids, ok := d.GetOk("ids"); ok && len(ids.([]interface{})) > 0 {
		args.ids = Interface2StrSlice(ids.([]interface{}))
	}
	if regx, ok := d.GetOk("name_regex"); ok && regx.(string) != "" {
		args.nameRegex = regx.(string)
	}
	return args
}

type SearchIndexDataSourceArgs struct {
	instanceName string
	tableName    string
	ids          []string
	nameRegex    string
}

type SearchIndexDataSource struct {
	ids     []string
	names   []string
	indexes []map[string]interface{}
}

func (s *SearchIndexDataSource) export(d *schema.ResourceData) error {
	d.SetId(dataResourceIdHash(s.ids))
	if err := d.Set("indexes", s.indexes); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", s.names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", s.ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if filepath, ok := d.GetOk("output_file"); ok && filepath.(string) != "" {
		err := writeToFile(filepath.(string), s.indexes)
		if err != nil {
			return err
		}
	}
	return nil
}

func dataSourceAlicloudOtsSearchIndexesRead(d *schema.ResourceData, meta interface{}) error {
	args := parseSearchIndexDataSourceArgs(d)

	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	totalIndexes, err := otsService.ListOtsSearchIndex(args.instanceName, args.tableName)
	if err != nil {
		return WrapError(err)
	}

	filteredIndexes := args.doFilters(totalIndexes)
	source, err := genSearchIndexDataSource(&otsService, filteredIndexes, args)
	if err != nil {
		return WrapError(err)
	}
	if err := source.export(d); err != nil {
		return WrapError(err)
	}
	return nil
}

func genSearchIndexDataSource(otsService *OtsService, filteredIndexes []interface{}, args *SearchIndexDataSourceArgs) (*SearchIndexDataSource, error) {
	size := len(filteredIndexes)
	ids := make([]string, 0, size)
	names := make([]string, 0, size)
	indexes := make([]map[string]interface{}, 0, size)

	for _, idx := range filteredIndexes {
		idxInfo := idx.(*tablestore.IndexInfo)
		id := ID(args.instanceName, args.tableName, idxInfo.IndexName, SearchIndexTypeHolder)

		indexResp, err := otsService.DescribeOtsSearchIndex(id)
		if err != nil {
			return nil, WrapError(err)
		}
		phase, err := ConvertSearchIndexSyncPhase(indexResp.SyncStat.SyncPhase)
		if err != nil {
			return nil, WrapError(err)
		}

		b, err := json.MarshalIndent(indexResp.Schema, "", "  ")
		if err != nil {
			return nil, WrapError(err)
		}
		schemaJSON := string(b)

		index := map[string]interface{}{
			"id":            id,
			"instance_name": args.instanceName,
			"table_name":    args.tableName,
			"index_name":    idxInfo.IndexName,

			"create_time":               indexResp.CreateTime,
			"time_to_live":              indexResp.TimeToLive,
			"sync_phase":                phase,
			"current_sync_timestamp":    *indexResp.SyncStat.CurrentSyncTimestamp,
			"storage_size":              indexResp.MeteringInfo.StorageSize,
			"row_count":                 indexResp.MeteringInfo.RowCount,
			"reserved_read_cu":          indexResp.MeteringInfo.ReservedReadCU,
			"metering_last_update_time": indexResp.MeteringInfo.LastUpdateTime,
			"schema":                    schemaJSON,
		}

		names = append(names, idxInfo.IndexName)
		ids = append(ids, id)
		indexes = append(indexes, index)
	}

	return &SearchIndexDataSource{
		ids:     ids,
		names:   names,
		indexes: indexes,
	}, nil
}

func (args *SearchIndexDataSourceArgs) doFilters(total []*tablestore.IndexInfo) []interface{} {
	slice := make([]interface{}, len(total))
	for i, t := range total {
		slice[i] = t
	}
	ds := InputDataSource{
		inputs: slice,
	}
	// add filter: index id
	if args.ids != nil {
		idFilter := &ValuesFilter{
			allowedValues: Str2InterfaceSlice(args.ids),
			getSourceValue: func(sourceObj interface{}) interface{} {
				idx := sourceObj.(*tablestore.IndexInfo)
				// search index and search index can have same IndexName, so joined IndexType for index id
				return ID(args.instanceName, args.tableName, idx.IndexName, SearchIndexTypeHolder)
			},
		}
		ds.filters = append(ds.filters, idFilter)
	}
	// add filter: index name regex
	if args.nameRegex != "" {
		regxFilter := &RegxFilter{
			regx: regexp.MustCompile(args.nameRegex),
			getSourceValue: func(sourceObj interface{}) interface{} {
				return sourceObj.(*tablestore.IndexInfo).IndexName
			},
		}
		ds.filters = append(ds.filters, regxFilter)
	}

	filtered := ds.doFilters()
	return filtered
}
