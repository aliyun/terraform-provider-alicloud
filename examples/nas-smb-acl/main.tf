resource "alicloud_nas_file_system" "foo" {
  protocol_type    = "SMB"
  storage_type     = "Performance"
  file_system_type = "standard"
  zone_id          = "cn-hangzhou-g"
}

resource "alicloud_nas_smb_acl_attachment" "main" {
  file_system_id = alicloud_nas_file_system.foo.id
  keytab         = "BQIAAABHAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAAQAIqIx6v7p11oUAAABHAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAAwAIqIx6v7p11oUAAABPAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAFwAQnQZWB3RAPHU7PMIJyBWePAAAAF8AAgANQUxJQURURVNULkNPTQAEY2lmcwAZc21ic2VydmVyMjQuYWxpYWR0ZXN0LmNvbQAAAAEAAAAAAQASACAGJ7F0s+bcBjf6jD5HlvlRLmPSOW+qDZe0Qk0lQcf8WwAAAE8AAgANQUxJQURURVNULkNPTQAEY2lmcwAZc21ic2VydmVyMjQuYWxpYWR0ZXN0LmNvbQAAAAEAAAAAAQARABDdFmanrSIatnDDhoOXYadj"
  keytab_md5     = "E3CCF7E2416DF04FA958AA4513EA29E8"
}