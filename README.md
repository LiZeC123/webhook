WebHook
==============
![](https://img.shields.io/github/license/LiZeC123/webhook)
![](https://img.shields.io/github/issues/LiZeC123/webhook)
![](https://img.shields.io/github/v/tag/LiZeC123/webhook)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/LiZeC123/webhook)
![GitHub last commit](https://img.shields.io/github/last-commit/LiZeC123/webhook)


一个基于Go语言的Web Hook


### 配置Config文件

项目启动时读取项目根路径上的`config.json`文件，此文件的实例如下所示

```json
{
  "Token": "fHxs3dsA",
  "Config": [
    {
      "appName": "<自定义应用名>",
      "type": "<自定义类型>",
      "template": "<Shell文件名>"
    },
    {
      "appName": "Blog",
      "type": "GithubHook",
      "template": "Update.sh"
    },
    {
      "appName": "System",
      "type": "User",
      "template": "SystemInfo.sh"
    }
  ]
}
```

其中Token是一个用来增加路径复杂程度的字符串，可以是任意的随机字符。appName是用户自定义的应用名称，type是用户自定义的操作类型。

type字段用来区分来源， 例如要同时处理来自Github和Gitlab的WebHook，则可以分别配置不同的type字符串来区分不同的请求。

访问路径为`/<Token>/<type>/<appName>`，例如按照上述的配置，如果访问URL`/fHxs3dsA/GithubHook/Blog`，则会调用位于`command`路径中的`Test.sh`脚本，并将appName作为第一个参数传入。



