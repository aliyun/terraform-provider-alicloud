package alicloud

import (
	"regexp"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOtsSecondaryIndexes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOtsSecondaryIndexesRead,

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
						"index_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"defined_columns": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func parseSecIndexDataSourceArgs(d *schema.ResourceData) *SecIndexDataSourceArgs {
	args := &SecIndexDataSourceArgs{
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

type SecIndexDataSourceArgs struct {
	instanceName string
	tableName    string
	ids          []string
	nameRegex    string
}

type SecIndexDataSource struct {
	ids     []string
	names   []string
	indexes []map[string]interface{}
}

func (s *SecIndexDataSource) export(d *schema.ResourceData) error {
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

func dataSourceAlicloudOtsSecondaryIndexesRead(d *schema.ResourceData, meta interface{}) error {
	args := parseSecIndexDataSourceArgs(d)

	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	total, err := otsService.ListOtsSecondaryIndex(args.instanceName, args.tableName)
	if err != nil {
		return WrapError(err)
	}

	filteredIndexes := args.doFilters(total)
	source, err := genSecIndexDataSource(filteredIndexes, args)
	if err != nil {
		return WrapError(err)
	}
	if err := source.export(d); err != nil {
		return WrapError(err)
	}
	return nil
}

func genSecIndexDataSource(filteredIndexes []interface{}, args *SecIndexDataSourceArgs) (*SecIndexDataSource, error) {
	size := len(filteredIndexes)
	ids := make([]string, 0, size)
	names := make([]string, 0, size)
	indexes := make([]map[string]interface{}, 0, size)

	for _, idx := range filteredIndexes {
		idxMeta := idx.(*tablestore.IndexMeta)
		typeName, err := ConvertSecIndexType(idxMeta.IndexType)
		if err != nil {
			return nil, WrapError(err)
		}
		id := ID(args.instanceName, args.tableName, idxMeta.IndexName, string(typeName))

		index := map[string]interface{}{
			"id":              id,
			"instance_name":   args.instanceName,
			"table_name":      args.tableName,
			"index_name":      idxMeta.IndexName,
			"index_type":      typeName,
			"primary_keys":    idxMeta.Primarykey,
			"defined_columns": idxMeta.DefinedColumns,
		}

		names = append(names, idxMeta.IndexName)
		ids = append(ids, id)
		indexes = append(indexes, index)
	}

	return &SecIndexDataSource{
		ids:     ids,
		names:   names,
		indexes: indexes,
	}, nil
}

func (args *SecIndexDataSourceArgs) doFilters(total []*tablestore.IndexMeta) []interface{} {
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
				idx := sourceObj.(*tablestore.IndexMeta)
				// secondary index and search index can have same IndexName, so joined IndexType for index id
				indexType, err := ConvertSecIndexType(idx.IndexType)
				if err != nil {
					return nil
				}
				return ID(args.instanceName, args.tableName, idx.IndexName, string(indexType))
			},
		}
		ds.filters = append(ds.filters, idFilter)
	}
	// add filter: index name regex
	if args.nameRegex != "" {
		regxFilter := &RegxFilter{
			regx: regexp.MustCompile(args.nameRegex),
			getSourceValue: func(sourceObj interface{}) interface{} {
				return sourceObj.(*tablestore.IndexMeta).IndexName
			},
		}
		ds.filters = append(ds.filters, regxFilter)
	}

	filtered := ds.doFilters()
	return filtered
}
