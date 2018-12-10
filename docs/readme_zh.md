#  apiServer

README: [ENGLISGH](../README.md) | [中文](./readme_zh.md)

## 概述

BOX（Enterprise Token Safe Box）是一个企业级数字资产保险柜应
用，它利用区块链、密码学、通信安全等领域的公理性技术对各类数字资
产的私钥、操作指令进行保护，从原理上解决了私钥、指令的盗取、篡改
等问题。更多详细信息，请参考[BOX官网](https://box.la/zh)及[BOX项目白皮书](https://box.la/static/BOX_white_paper_zh.pdf)。

### 目录

- [项目特点](#项目特点)
- [安装部署](#安装部署)
- [快速开始](#快速开始)
- [开发](#开发)
- [授权许可](#授权许可)

### 项目特点

- 审批流自动转币，账⽬规范明晰
- 私钥多人共管，自主加密保存
- 私钥，审批流指令，通讯全⽅位顶级安全
- 多币种一站式管理
- 代码开源免费，官⽅团队提供技术支持

### 安装部署

- 安装部署apiServer之前，您需要先安装基本环境。
- 生成加密通信所需的证书：

    $sh scripts/keys.sh

- 修改config/config.toml.example文件，正确填写数据库配置信息，以及与签名机加密通信的证书目录，默认为项目根目录的certs文件夹下。修改次文件名为config.toml
- 配置wallet模块
  - 修改项目根目录path.toml，用以指定wallet模块运行目录
  - 进入配置好的wallet运行目录，修改config.yml文件，正确配置MySQL
  - 本地全节点配置请修改local目录下相关文件

#### 快速开始

- 获取代码

    $ git clone git@gitlab.2se.com:boxproject/apiServer.git

- 编译

    $ make build

- 重新编译

    $ make rebuild

- 启动

    $ ./apiServer start

### 开发

相关接口文档，请参考[API.md](./api.md)

### 授权许可

Licensed under the Apache License, Version 2.0, Copyright 2018. bolaxy.org authors.

     Copyright 2018. bolaxy.org authors.
    
     Licensed under the Apache License, Version 2.0 (the "License");
     you may not use this file except in compliance with the License.
     You may obtain a copy of the License at
    
          http://www.apache.org/licenses/LICENSE-2.0
    
     Unless required by applicable law or agreed to in writing, software
     distributed under the License is distributed on an "AS IS" BASIS,
     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
     See the License for the specific language governing permissions and
     limitations under the License.


