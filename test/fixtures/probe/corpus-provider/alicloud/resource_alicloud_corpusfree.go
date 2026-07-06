package alicloud

// Fixture source for corpus-gen product-derivation tests. 不对应真实资源。
// _source_pv 应从 RpcPost 抽出 (Corpusfree, 2023-01-01) → 产品目录 corpusfree。

func resourceAliCloudCorpusfreeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateCorpusfree"
	request := make(map[string]interface{})
	response, err := client.RpcPost("Corpusfree", "2023-01-01", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_corpusfree", action, AlibabaCloudSdkGoERROR)
	}
	_ = response
	return nil
}
