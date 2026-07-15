命名空间（namespace）是模型的空间名称，在镇元（公有云）环境下，namespace的组成如下：

```bash
alicloud.{PRODUCT}.{POP-CODE}.v{POP-VERSION}
```

例如：

```bash
alicloud.ACVS.acvs.v20210910
```

其中：

+ alicloud 为固定的字段，表示阿里云公有云；
+ {PRODUCT}为产品名称；
+ {POP-CODE} 为POP code；
+ {POP-VERSION} 为POP version，通常是日期版本。



**_需要注意的是，部分inner API没有关联到任何product，对于这部分API，{PRODUCT} 固定为inner。_**

