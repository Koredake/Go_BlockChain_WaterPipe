# go程序设计-测试水管

## 步骤：
* 更新module依赖
* 更改配置文件
* 编译运行

## 1 更新module
在命令行进入todo所在的目录，执行如下命令，更新module依赖的代码包。
```
go mod tidy
```

## 2 更改配置文件

在cfg.json中将配置改为本机配置

**注：不更改以下配置程序将不能正常运行！！！**

```
{
    "gethPort": "http://localhost:7545",                                    //以太坊RPC端口
    "privateKey" : "34880979e84a928e04a031a3d1eb590cbde...",              //挖矿收益地址私钥
    "dbUserName": "root",                                                    //mysql用户名
    "dbPassword": "12345678",                                                  //mysql密码
    "DbIp": "127.0.0.1",                                                       //mysql ip
    "DbPort": "3306",                                                         //mysql 端口
    "dbName": "chaozhichang"                                           //mysql 本地数据库名
}

```

### 3 编译运行

执行如下命令，编译todo代码，并运行程序。
```
go build -o waterPipe.exe .
waterPipe.exe cfg.json
```

在浏览器中打开 http://127.0.0.1:12345/web ，可以看到我们创建好的测试币水管项目，可以输入账户地址进行测试币申请。


### 4 编译运行

重复第3部的操作，编译代码，并运行。

