# checks: import-vpc (fixture)

importer persona:state rm → import(id 取自 output vpc_id)→ plan 应空 diff(import 还原完整 state)。
