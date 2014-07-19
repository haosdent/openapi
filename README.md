阿里云API工具
----------

### 免责声明：
写这个工具仅仅只是自己为图方便使用阿里云API，与阿里云或任何团体无关。如果在使用过程中遇到任何后果，请自行承担！

### 简介：
参照[阿里云帮助中心][1]的文档写成，为啥要自己写的原因是阿里云没有提供给我等吊丝开发者像AWS一样完善方便的控制台工具，所以只能自己写。主要提供了更方便的帮助，避免在使用工具时查阅阿里云API文档而浪费时间。

### 安装：

可手工从源码安装或者直接下载已预先编译好的包。

#### 直接下载预先编译好的包（选择对应的操作系统）：

- [Windows 64位下的版本][2]
- [Linux下的版本][3]
- [OSX下的版本][4]

#### 从源码编译安装：

从Github上clone源码到本地路径后，

```shell
$ cd src
$ go build openapi.go
```

然后执行生成的openapi文件即可

```shell
$ ls
help.json   openapi    openapi.go
$ ./openapi help
Usage: openapi [production] [action] [property=value]
Allow production:
        ecs
        rds
        slb
$ ./openapi ecs DescribeRegion
Config your accessId and accessKey through 'openapi config --id=accessId --secret=accessKey' firstly.
```

### 使用：

#### 配置AccessId和AccessKey：

```
$ ./openapi config --id=accessId --secret=accessKey
```

#### 帮助

```
$ ./openapi help ecs StartInstance
Usage: openapi ecs StartInstance [property=value]
Allow property:
        RegionId
        InstanceId
        ForceStop
```

```
$ ./openapi help rds
Usage: openapi rds [action] [property=value]
Allow action:
        ModifySecurityIps
        CreateBackup
        DescribeResourceUsage
        DescribeDBInstancePerformance
        CreateDatabase
        DescribeDatabases
        GrantAccountPrivilege
        DescribeAccounts
        RestoreDBInstance
        RestartDBInstance
        DeleteDatabase
        RevokeAccountPrivilege
        DescribeSecurityIps
        DescribeDBInstances
        CreateAccount
        ModifyAccountAttribute
        DescribeBackups
        SwitchDBInstanceNetType
        CreateImportDataUpload
        ImportData
        DescribeDataFiles
```

```shell
$ ./openapi help
Usage: openapi [production] [action] [property=value]
Allow production:
        ecs
        rds
        slb
```

### 开源License

```
The MIT License (MIT)

Copyright (c) 2014 haosdent

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```


  [1]: http://dev.aliyun.com/thread.php?fid=8
  [2]: https://github.com/haosdent/openapi/raw/master/bin/windows/openapi.exe
  [3]: https://github.com/haosdent/openapi/raw/master/bin/linux/openapi
  [4]: https://github.com/haosdent/openapi/raw/master/bin/osx/openapi
  [5]: https://raw.githubusercontent.com/haosdent/openapi/master/LICENSE
