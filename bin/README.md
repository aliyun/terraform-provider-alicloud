# bin/a1id —— a1 多身份切换

a1 凭据固定在 `~/.config/a1/auth.yaml`,**只能持一个身份**(且 a1 写入用临时文件 rename,符号链接会被替换)。`a1id` 把各身份的 auth.yaml 各存一份真文件,切换=复制覆盖 live,串行切换(a1 本就不支持并发)。

## 四个身份
- **jarvis**(默认,open_jarvis):Jarvis 全程用它。
- **chenyi**(仓库主人辰羿本人):Jarvis 不得擅用,仅当面授权时临时切。无代码硬门,纪律见 CLAUDE.md「身份纪律」。
- **guozai**(仓库主人过载本人):同 chenyi,Jarvis 不得擅用,仅当面授权时临时切。
- **linjun**(李超林,工号 429768):同 chenyi,Jarvis 不得擅用,仅当面授权时临时切。

## 一次性设置
```bash
bin/a1id login jarvis     # 已有 auth.yaml 首跑自动收编为 jarvis
bin/a1id login chenyi     # 辰羿本人登录一次,存为 chenyi 身份
bin/a1id login guozai     # 过载本人登录一次,存为 guozai 身份
bin/a1id login linjun     # 李超林登录一次,存为 linjun 身份
```

## 用法
- 默认身份(jarvis)跑:`a1id -- <a1 args>`
- 临时以辰羿跑一条(用完自动还原):`a1id as chenyi -- <args>`
- 临时以过载跑一条(用完自动还原):`a1id as guozai -- <args>`
- 临时以李超林跑一条(用完自动还原):`a1id as linjun -- <args>`
- 持久切:`a1id use <jarvis|chenyi|guozai|linjun>`
- 看状态:`a1id status` / `a1id who`(a1 实际账号)
