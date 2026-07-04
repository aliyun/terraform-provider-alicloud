# checks: net-vpc-basic (fixture)

最小 VPC + VSwitch 组网,供 probe_test.sh hermetic 遍历,不真实 apply。

## 期望
1. validate 通过。
2. apply 后立即 plan 空 diff。
3. destroy 干净、state 清空。
