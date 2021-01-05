```
[server]
ServerType = "ssh"  可选local保存到本地 或ssh保存到另一个服务器 
SourceHost = ""     服务器ip 如:111.111.111.111
SourceUser = "root" 登陆账号
SourcePwd = ""      登陆密码
SourcePort = 22     ssh端口，服务器默认22
TargetHost = ""     目标服务器ip，如果是保存到本地则不需要填写
TargetUser = "root" 目标服务器登陆账号
TargetPwd = ""      目标服务器密码
TargetPort = 22     目标服务器端口

[file_info]
Source = "/home/go/src/gopkg.in"     源文件夹路径
Target = "/home/go/temp"             目标路径，保存本地可写相对路径，如./backup 或写绝对路径

需要排除的文件夹，以,分割开,请填写绝对路径
ExceptDir = "/www/wwwroot/default/gocode/src/github.com,/www/wwwroot/default/gocode/src/gitee.com"

MaxChannelNum = 10  并发传输数量
```
> 日志表格模板在 template 文件夹下 

> 如果不需要表格自动换行 把head里单独那行css去掉即可

> 调整输出信息，信息、日志放在runtime/logs下

> 本地文件迁移 服务器文件迁移

> 简单修复了并发问题，其他修改都放在v2版本