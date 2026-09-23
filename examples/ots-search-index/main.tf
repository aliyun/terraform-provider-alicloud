data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "example-ots-table"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "example-ots-table"
  cidr_block   = "172.16.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_ots_instance" "default" {
  name        = var.ots_instance_name
  description = "TF ots instance example"
}

resource "alicloud_ots_instance_attachment" "default" {
  instance_name = alicloud_ots_instance.default.id
  vswitch_id    = alicloud_vswitch.default.id
  vpc_name      = "table"
}

resource "alicloud_ots_table" "table" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = var.table_name

  primary_key {
    name = var.primary_key_1_name
    type = var.integer_type
  }
  primary_key {
    name = var.primary_key_2_name
    type = var.string_type
  }


  defined_column {
    name = var.defined_column_1_name
    type = var.string_type
  }

  defined_column {
    name = var.defined_column_2_name
    type = var.integer_type
  }

  time_to_live = var.time_to_live
  max_version  = var.max_version
}

resource "alicloud_ots_search_index" "index1" {
  instance_name =  alicloud_ots_instance.default.name
  table_name = alicloud_ots_table.table.table_name

  index_name = var.search_index_name
  time_to_live = var.search_index_ttl
  schema {
    field_schema {
      field_name = var.defined_column_1_name
      field_type = "Text"
      is_array = false
      index = true
      analyzer = "Split"
      store = true
    }
    field_schema {
      field_name =  var.defined_column_2_name
      field_type = "Long"
      enable_sort_and_agg = true
    }


    field_schema {
      field_name =  var.primary_key_1_name
      field_type = "Long"

    }
    field_schema {
      field_name =  var.primary_key_2_name
      field_type = "Text"

    }


    index_setting {
      routing_fields = [ var.primary_key_1_name, var.primary_key_2_name]
    }

    index_sort {
      sorter {
        sorter_type = "PrimaryKeySort"
        order = "Asc"
      }
      sorter {
        sorter_type = "FieldSort"
        order = "Desc"
        field_name =  var.defined_column_2_name
        mode = "Max"
      }
    }
  }
}