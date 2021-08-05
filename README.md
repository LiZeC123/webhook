WebHook
==============







### 交叉编译

在Window平台交叉编译Linux平台可执行文件

```bash
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o bin
```


