db_cluster_nodes_configs = {
  node_reader_1 = {
    db_node_class = "polar.mysql.x4.medium.c"
    db_node_role  = "Reader"
  }
  node_reader_2 = {
    db_node_class = "polar.mysql.x4.medium.c"
    db_node_role  = "Writer"
  }
}