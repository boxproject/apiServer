# API 文档

## 0 初始化签名机

```bash
Msg 结构体
type Msg struct {
    Password  string    // 关键句 2,3,4,5,6
    BakPasswd string    // 备份密码 2,
    PublicKey string    // 公钥 1,
    ...
}
```

### 0.1 【私钥 APP-->voucher】添加公钥【经 RSA 加密】（初次建立联系）

+ router:/api/v1/voucher/qrcommit
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| app_id    |   string   | appid  |
| name      |   string  | 私钥 app 用户名  |
| sign      |  string   |sign   |
| timestamp |   int  | 时间戳  |
| aeskey    | string | 加密 key|
| Msg       | string  | 签名机整体加密参数 PublicKey  |

+ 返回值

返回：status

```json
{
  "Code": 0,
  "Msg": "success",
  "Data": {
    "Type": "string", //req type
    "AppId": "string", //app id
    "AppName": "string", //app name
    "PublicKeys": "string", //app public key;map[appId]pubKeyBytes
    "Sign": "string", //app msg sign
    "ContractAddr": "string", //contract address
    "D": "string", //voucher random D value
    "Status": "int", //server req status  0. 成功  1. 失败
    "Other": "string", //other
    "Timestamp": "int64" //Timestamp, 时间戳
  }
}
```

### 0.2 【私钥 APP-->voucher】私钥口令 + 口令备份密码接收（N 个每轮, 限制接收指定 N 个 [口令 + ID]，抢占式）

+ router:/api/v1/voucher/commandcommit
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| app_id    |   string | appid |
| name      |   string  | 私钥 app 用户名 |
| sign      |  string |sign  |
| timestamp |   int  | 时间戳  |
| aeskey    | string | 加密 key |
| Msg       | string | 签名机整体加密参数 BakPasswd,Password  |

+ 返回值

返回：status

```json
{
  "Code": 0,
  "Msg": "success",
  "Data": {
    "Type": "string", //req type
    "AppId": "string", //app id
    "AppName": "string", //app name
    "PublicKeys": "string", //app public key;map[appId]pubKeyBytes
    "Sign": "string", //app msg sign
    "ContractAddr": "string", //contract address
    "D": "string", //voucher random D value
    "Status": "int", //server req status
    "Other": "string", //other
    "Timestamp": "int64" //Timestamp, 时间戳
  }
}
```

### 0.3 【私钥 APP-->voucher】签名机随机口令 + 私钥 APP 公钥获取请求（签名机需要返回加密过后的随机口令, 私钥 APP 在备份前向 voucher 请求）

+ router:/api/v1/voucher/getotherpubkey
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| app_id    |   string  | appid |
| name      |   string   | 私钥 app 用户名 |
| sign      |  string |sign  |
| timestamp |   int | 时间戳  |
| aeskey    | string  | 加密 key  |
| msg       | string  | 签名机整体加密参数 Password |
| password    | string  |  app 登陆密码   |

+ 返回值

返回：status，以及 pubkey+D

```json
{
  "Code": 0,
  "Msg": "success",
  "Data": {
    "Type": "string", //req type
    "AppId": "string", //app id
    "AppName": "string", //app name
    "PublicKeys": [ //app 所有股东公钥（密文）
      {
        "AppName": "app1",
        "AppId": "appid1",
        "AppPubKey": "aaa"
      },
      {
        "AppName": "app2",
        "AppId": "appid2",
        "AppPubKey": "bbb"
      }
    ],
    "Sign": "string", //app msg sign
    "ContractAddr": "string", //contract address
    "D": "string", //voucher random D value
    "Status": "int", //server req status
    "Other": "string", //other
    "Timestamp": "int64", //Timestamp, 时间戳
    "AesKey": "string", //aes key
  }
}
```

### 0.4 【私钥 APP-->voucher】私钥 APP 硬件【备份请求 + 备份状态查询请求】（签名机返回当前备份状态：是否备份 ing、是否有就绪设备接入、备份成功次数等……）

+ router:/api/v1/voucher/backupkey
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| app_id    |   string  | appid  |
| name      |   string  | 私钥 app 用户名  |
| sign      |  string  |sign   |
| timestamp |   int | 时间戳 |
| aeskey    | string  | 加密 key   |
| Msg       | string  | 签名机整体加密参数 Password |
| bak_action | string |1. 添加 app 信息等待备份  2. 备份信息  3. 验证信息  4. 备份取消 |

+ 返回值

```json
{
  "Code": 0,
  "Msg": "success",
  "Data": {
    "Type": "string", //req type
    "AppId": "string", //app id
    "AppName": "string", //app name
    "PublicKeys": "string", //app public key;map[appId]pubKeyBytes
    "Sign": "string", //app msg sign
    "ContractAddr": "string", //contract address
    "D": "string", //voucher random D value
    "Status": "int", //server req status
    "Other": "string", //other
    "Timestamp": "int64", //Timestamp, 时间戳
  }
}
```

### 0.5 【私钥 APP-->voucher】签名机服务启动请求

+ router:/api/v1/voucher/start
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| app_id    |   string          | appid  |
| name      |   string          | 私钥 app 用户名  |
| sign      |  string           |sign    |
| timestamp |   int             | 时间戳   |
| aeskey    | string            | 加密 key  |
| msg       | string            | 签名机整体加密参数 Password  |

+ 返回值

```json
{
  "Code": 0,
  "Msg": "success",
  "Data": {
    "Type": "string", //req type
    "AppId": "string", //app id
    "AppName": "string", //app name
    "PublicKeys": "string", //app public key;map[appId]pubKeyBytes
    "Sign": "string", //app msg sign
    "ContractAddr": "string", //contract address
    "D": "string", //voucher random D value
    "Status": "int", //server req status
    "Other": "string", //other
    "Timestamp": "int64", //Timestamp, 时间戳
  }
}
```

### 0.6 【私钥 APP-->voucher】签名机服务停止请求

+ router:/api/v1/voucher/stop
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| app_id    |   string | appid |
| name      |   string  | 私钥 app 用户名  |
| sign      |  string  |sign  |
| timestamp |   int  | 时间戳 |
| aeskey    | string  | 加密 key  |
| msg       | string  | 签名机整体加密参数 Password  |

+ 返回值

```json
{
  "Code": 0,
  "Msg": "success",
  "Data": {
    "Type": "string", //req type
    "AppId": "string", //app id
    "AppName": "string", //app name
    "PublicKeys": "string", //app public key;map[appId]pubKeyBytes
    "Sign": "string", //app msg sign
    "ContractAddr": "string", //contract address
    "D": "string", //voucher random D value
    "Status": "int", //server req status
    "Other": "string", //other
    "Timestamp": "int64", //Timestamp, 时间戳
  }
}
```

### 0.7 签名机状态

+ router:/api/v1/voucher/getstatus
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| 无    |       |                             |

+ 返回值

```json
{
  "Code": 0,
  "Msg": "success",
  "Data": {
    "ServerStatus": "int", // 系统状态     0 // 未连接 1// 未创建 2 // 已创建 3 // 已备份 4 // 已启动  5 // 已停止
    "Total": "int", // 密钥数量
    "BakDevStatus": 0, //0. 未连接  1. 设备就绪  2. 设备备份完成  3. 设备本份失败
    "BakStatus": {
      "app_1": 2, //account  , 几次
      "app_2": 5, //account  , 几次
      "app_3": 5, //account  , 几次
    },
    "EthNetState": 0, //0. 正常  -1. 不正常
    "BaseAddress": "string", // 主账户 eth 地址
    "ContractAddress": "string", // 合约地址
    "NodesAuthorized": { //create：备份口令轮询     start: 启动轮询
      "Timestamp": 1536890629666, // 时间戳
      "create": [{
          "ApplyerId": "1534152759518", //appid，
          "Authorized": true
        },
        {
          "ApplyerId": "1534152716649",
          "Authorized": true
        },
        {
          "ApplyerId": "1534152742892",
          "Authorized": true
        }
      ],
      "start": [{
          "ApplyerId": "1534152759518", //appid
          "Authorized": true
        },
        {
          "ApplyerId": "1534152716649",
          "Authorized": true
        },
        {
          "ApplyerId": "1534152742892",
          "Authorized": true
        }
      ],
      "backup": [{
          "ApplyerId": "1534152759518", //appid
          "Authorized": true
        },
        {
          "ApplyerId": "1534152716649",
          "Authorized": true
        },
        {
          "ApplyerId": "1534152742892",
          "Authorized": true
        }
      ]
    }, // 授权情况
    "KeyStoreStatus": {
      "app-001": {
        "ApplyerId": "001", //appid
        "ApplyerName": "app-001" //appname
      },
      "app-002": {
        "ApplyerId": "002",
        "ApplyerName": "app-002"
      }
    } // 公约添加状态
  }
}
```

### 0.8 获取 token

+ router:/api/v1/voucher/gettoken
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| app_id    |   string    |  appid   |
| name    |   string    | app 用户名   |
| sign    |   string    |   对时间戳签名  |
| timestamp    |   int    |        时间戳             |

+ 返回值

```json
{
  "Code": 0,
  "Msg": "success",
  "Data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBpZCI6ImRkZGQiLCJhY2NvdW50Ijoid3d3IiwiZXhwIjoxNTM1MzU0NjUzfQ.8yUC4bqL0g44dRSgWmocX4cf4Yg_m5FfGUOegJO61dI",
    "name": "appname1",
    "userType": 2, // 用户类型  0 普通用户 1 管理员 2 股东
  }
}
```

### 0.9 获取签名机稳定状态

+ router:/api/v1/voucher/getconnectstatus
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| 无    |       |     |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "status": 2,  // 0.稳定  1.不稳定 2.异常 3.等待启动 4。关停
    "start" :[
        {
            ApplyerId:"sdfdsfsfsdfsfs"  //已经启动的人的appid
        },
        {
            ApplyerId:"fefefefefe"
        }
    ]
  },
  "Msg": "success"
}
```

### 0.10 获取主地址

+ router:/api/v1/voucher/getmasteraddress
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| 无    |       |     |

+ 返回值

```json
{
    "Code": 0,
    "Data": [
        {
            "Address": "n2Jw5VwQYvxnHreuihxcH2TfrBGtswdwB7",
            "CoinName": "BTC"
        },
        {
            "Address": "0x2882A4C5Fc11EE72B2A1004aa8EBc4fEa119e1Cf",
            "CoinName": "ETH"
        },
        {
            "Address": "n2Jw5VwQYvxnHreuihxcH2TfrBGtswdwB7",
            "CoinName": "LTC"
        },
        {
            "Address": "n2Jw5VwQYvxnHreuihxcH2TfrBGtswdwB7",
            "CoinName": "USDT"
        }
    ],
    "Msg": "success"
}
```

### 0.11 保存主地址签名

+ router:/api/v1/voucher/savemasteraddress
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| msg    |    string   |   看注释  |

注释：msg :{
​        "account": "zhangsan",
​        "addressSign": [{
​        "coinName": "ETH",
​        "masterAddress": "fdfdfdfdfd",
​        "sign": "fdfdfdfdfdf"
​      },{
​        "coinName": "BTC",
​        "masterAddress": "fdfdfdfdfd",
​        "sign": "fdfdfdfdfdf"
​      },{
​        "coinName": "USDT",
​        "masterAddress": "fdfdfdfdfd",
​        "sign": "fdfdfdfdfdf"
​      }]
​    }

+ 返回值

```json
{
    "Code": 0,
    "Data": {
        "BTC": "n2Jw5VwQYvxnHreuihxcH2TfrBGtswdwB7",
        "ETH": "0x2882A4C5Fc11EE72B2A1004aa8EBc4fEa119e1Cf",
        "LTC": "n2Jw5VwQYvxnHreuihxcH2TfrBGtswdwB7",
        "USDT": "n2Jw5VwQYvxnHreuihxcH2TfrBGtswdwB7"
    },
    "Msg": "success"
}
```

### 0.12 【私钥 APP-->voucher】签名机与app交换密钥

+ router:/api/v1/voucher/operatekeyexchange
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| app_id    |   string | appid |
| name      |   string  | 私钥 app 用户名  |
| sign      |  string  |sign  |
| timestamp |   int  | 时间戳 |
| aeskey    | string  | 加密 key  |
| msg       | string  | 签名机整体加密参数 Password  |

+ 返回值

```json
{
 "Code": 0,
 "Data": {
  "Type": "15",
  "AppId": "C8BD6CFC-1F70-4DBE-A886-E77FA93980BA",
  "AppName": "上世纪",
  "AesKey": "BUQU49h9YI4QTd/1w9fPmzE2Np7ISoyyeFra0w7QeuDFlYFJtzFXPBcPg1UYt0rr9JA9WlRubLWcTo0ClWvk4priSXp6ij+QKliSSef0J9uU/LraMCAWBSxhfFQBWyLSeW/dEp0m1ysZWBoIUDvjP/pPIZkOTNP3rMRNBH6wPIlv/oI1hJvFC7zfdFgipx54RQtbGdbJ9P/uMjp9TBgEWzqZyVE93Te5yx1/e2vF1Q4Ve0RpzVtgldlD7LLGrhZdjAFBp51XoXxURmNDU3BlkBmX6WEHJ9ST7RtflhakXp37T4sAAetie/9HZAdKFknEz8ZiIx3rKuLDYFa5rVLZ0w==",
  "PublicKeys": "OsS7Gtuahq4Uktk6nnYNp0O6ukHW/3mCstn101PkrpCh1RAIiDYQ+t71ml0fMsxcq9HSgOvR99yixph/",
  "Sign": "V0ZIcTlBdkVIYWk2WkoxdHJCbVg1b2RNeUl6VHRnRWRvU1A0Z25BWFF3eEtwMmRlT3dFVFNQQldXN0hHMDI3Y2hGTzhZWWhjcUN1dWNtS1dadzNlYXpPSWRJS1N3aGFMNlVqUzdHd1U2ZWtKVmJCazB4U3RkZklwWk1wSFpJYm1VWGQ5NFR2VEc1d0NqOG9maDY3R0lWY3NFQjNpK1dMbVE3N2NrdnZLV3pySStaeGhMK0l0SjZIdHZhK3JhcnBTZzFiK1E5eTFwd0JYZXVVbkdpL25VY1ZQTnhQSS9VbHNCM0IrL2cyOTZ4UHdvYUhkQjY0UFF5ODVsbEpRdG81d3hwNGRyUll6SFB2NXMwTmVBUU9kUTZ4K2JRc3RSWERmR3I3dWRtWnJBaFhsaVQ0NmNRZXUwSHloSW9UVENab0JkKzRiS20xNDRaaDIzQmQ2TGJDMFp3PT0=",
  "ContractAddr": "",
  "D": "",
  "Status": 0,
  "SignTx": null,
  "Other": null,
  "Timestamp": 1541559619233,
  "Logs": "W3siVGltZXN0YW1wIjoxNTQxNTU5NjE5MjMzLCJBcHBJZCI6IkM4QkQ2Q0ZDLTFGNzAtNERCRS1BODg2LUU3N0ZBOTM5ODBCQSIsIkFwcE5hbWUiOiLkuIrkuJbnuqoiLCJPcGVyYXRlVHlwZSI6IjE1IiwiUmVzdWx0IjoiWzAg5q2j56GuXSIsIk90aGVyIjoiIn1d"
 },
 "Msg": "success"
}
```

## 1 用户

### 1.1 注册

+ router:  /api/v1/account/signup
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |   备注  |
| :----: | :----: | :--------: |
| name  | string |  账户名 |
| app_id | string | appid |
| pub_key | string | 公钥 |
| pwd | string | 密码 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,// 成功
  "Data": {},
  "Msg": "success"
}
```

### 1.2 登陆

+ router:  /api/v1/account/login
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |  备注    |
| :----: | :----: | :----------: |
| name  | string |    账户名  |
| pwd | string | 密码 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "name": "admin2",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MCwiYXBwaWQiOiIyIiwiYWNjb3VudCI6ImFkbWluMiIsInVzZXJUeXBlIjowLCJQdWJLZXkiOiIiLCJMZXZlbCI6MCwiZXhwIjoxNTM2NzI3MDk4wQ.KrxH7speWiFTya4DiUMgHPRGQUpymKl8Ldf0vf6mkCs",
    "userType": 2 // 0. 普通用户，1. 管理员，2. 股东
  },
  "Msg": "success"
}
```

### 1.3 新员工查询注册结果

+ router:  /api/v1/account/finduserstatus
+ 请求方式： GET
+ 参数：

|   字段   |   类型   |     备注    |
| :----: | :----: | :------------: |
| app_id  | string |    appid    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "refuse_reason": "",
    "status": 2, // 0. 注册 1. 已扫码，等待审核 2. 同意 3. 拒绝
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MCwiYXBwaWQiOiI0NmFiZTk4Mi1iIiwiYWNjb3VudCI6IjQ2YWJlOTgyLWIiLwJ1c2VyVHlwZSI6MCwiUHViS2V5IjoiIiwiTGV2ZWwiOjAsImV4cCI6MTUzNjcyNTgzM30.R4-Sc0cHaYX6vQs_Tw970Laj9B055_Zg_6nDh8f5ptU"
  },
  "Msg": "success"
}
```

### 1.4 重新注册

+ router:  /api/v1/account/resignup
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |     备注    |
| :----: | :----: | :------------: |
| app_id  | string |    appid    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 1.5 上级扫码

+ router:  /api/v1/user/scan
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |     备注    |
| :----: | :----: | :------------: |
| new_app_id  | string |    新员工的 appid    |
| msg  | string |   签名值   |
| token  | string |  token   |
| sign  | string |    参数签名   |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}

```

### 1.6 审核

+ router:  /api/v1/user/verifyuser
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |     备注     |
| :----: | :----: | :--------------: |
| new_app_id  | string |  新员工的 appid    |
| dep_id  | int |  部门 id，拒绝时非必须   |
| status  | int |  2. 同意  3. 拒绝    |
| refuse_reason  | string |    拒绝原因    |
| token  | string |    token     |
| sign  | string |    参数签名  |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 1.7 添加管理员

+ router:  /api/v1/user/addadmin
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |    备注     |
| :----: | :----: | :---------------: |
| add  | string |     需要添加为管理员的账户 id，多人以【,】分割             |
| token  | string |    token     |
| sign  | string |     参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 1.8 删除管理员

+ router:  /api/v1/user/deladmin
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |   备注    |
| :----: | :----: | :-------------: |
| account_id  | int |   需要删除管理员权限的账户 id    |
| token  | string |   token      |
| sign  | string |    参数签名     |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 1.9 通过用户类型获取用户列表

+ router:  /api/v1/user/listbytype
+ 请求方式： GET
+ 参数：

|   字段   |   类型   |     备注     |
| :----: | :----: | :--------------: |
| user_type  | int |     用户类型 0: 普通用户， 1: 管理员， 2: 股东       |
| token  | string |    token      |
| sign  | string |     参数签名      |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
    "ID": 1,
    "AppId":"123334",
    "Name": "test",
    "UserType": 1
  }],
  "Msg": "success"
}
```

### 1.10  修改密码

+ router: /api/v1/user/modifypassword
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   old_psd     | string |     旧密码     |
|   new_psd     | string |     新密码    |
| token  | string |     token     |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Msg": "success"
}
```

### 1.11  获取所有用户列表

+ router: /api/v1/user/list
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| token  | string |     token     |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
    "AppId": "123334",
    "Frozen": 0, //0. 未冻结，1. 冻结
    "ID": 1,
    "Name": "admin1",
    "UserType": 2
  }],
  "Msg": "success"
}
```

### 1.12  获取单个用户信息

+ router: /api/v1/user/userbyid
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| account_id  | int |     账户 id     |
| token  | string |     token     |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "Auths": [{ // 用户有的权限列表
      "ID": 2,
      "Name": "地址管理"
    }, {
      "ID": 1,
      "Name": "资产查询"
    }],
    "DepartmentId": 1,
    "DepartmentName": "技术部"
    "ID": 1,
    "AppId": "123334",
    "Name": "admin1"
  },
  "Msg": "success"
}
```

### 1.13 设置用户信息

+ router: /api/v1/user/setuser
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| account_id  | int |     账户 id     |
| auth_id  | string |  权限 id，多个以【,】分割 |
| dep_id  | int |    部门 id     |
| user_type  | int |   用户角色    |
| token  | string |     token     |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 1.14 停用账户

+ router: /api/v1/user/disableacc
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| account_id  | int |     账户 id     |
| token  | string |     token     |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 1.15 获取已注册用户列表

+ router: /api/v1/user/reglist
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| token  | string |     token     |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
    "AppID": "1122333"
    "CreatedAt": "2018-08-28T19:18:12+08:00",
    "Name": "ddd323",
    "SourceAccount": "test", // 扫码用户
    "Status": 1 // 1. 待审核 2. 已同意 3. 已拒绝
  }],
  "Msg": "success"
}
```

### 1.16 用户列表树

+ router: /api/v1/user/usertree
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| token  | string |     token     |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "trees":[{
    "Name": "admin1",
    "AppId": "1",
    "PubKey": "eeeefef",
    "Msg": "",
    "SourceAppId": "0",
    "Children": [{
      "Name": "A",
      "AppId": "11",
      "PubKey": "c86e37ba-a",
      "Msg": "%$1ss",
      "SourceAppId": "1",
      "Children": [{
        "Name": "B",
        "AppId": "21",
        "PubKey": "80677a6d-a",
        "Msg": "%$1ss",
        "SourceAppId": "11",
        "Children": null
        }]
      }]
    }],
    "voucher":{
      "ENPublickey":"eefefefefefefefefefefefefefefefefefefefeff",
      "ENKey":"swsfefefefeffdsfsfdsfdsfdfsfdsfsf"
    }
  },
  "Msg":"success"
}
```

### ~~1.17 校验密码并返回当前用户信息~~

+ router: /api/v1/user/verifypwd
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| pwd  | string |      密码     |
| token  | string |     token     |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "Auths": [{ // 用户有的权限列表
      "AuthId": 2,
      "AuthName": "地址管理"
    }, {
      "AuthId": 1,
      "AuthName": "资产查询"
    }],
    "Department": [{ // 用户所在部门
      "Department": "test",
      "ID": 3
    }],
    "ID": 1,
    "AppId": "123334",
    "Name": "admin1"
  },
  "Msg": "success"
}
```

### 1.18 申请恢复股东

+ router: /api/v1/account/recoveryowner
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| name  | string |    账户名    |
| pwd  | int |     密码    |
| pub_key  | string |    公钥    |
| app_id  | string |    app ID    |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "RegID": 12,
  },
  "Msg": "success"
}
```

### 1.19 重新注册申请恢复股东

+ router: /api/v1/account/resetrecovery
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| name  | string |    账户名    |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 1.20 扫码提交恢复股东申请

+ router: /api/v1/user/subrecovery
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| name  | string |    账户名    |
| app_id  | string |    app ID    |
| reg_id  | int |  注册 id    |
| token  | string |  token  |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 1.21 获取申请恢复股东账户列表

+ router: /api/v1/user/recoverylist
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| token  | string |  token  |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
    "AppId": "gudong_appid_2",
    "ID": 5,
    "Name": "gudong_2",
    "Status": 1,// 1. 待我认证，2. 完全通过, 3. 已拒绝，4. 我已认证
    "UpdatedAt": "2018-09-27T14:39:51+08:00"
  }],
  "Msg": "success"
}
```

### 1.22 认证股东

+ router: /api/v1/user/verifyrecovery
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| status  | int |  状态 1. 同意，2. 拒绝  |
| reg_id  | int | 申请列表单条记录 id  |
| token  | string |  token  |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |
| app_id    |   string    | appid    |
| name |   string  | 私钥 app 用户名 |
| aeskey | string | 加密 key|
| msg  | string  | 签名机整体加密参数 |

```bash
备注
msg:
Password  string //password      【股东密码，非登录口令】
BakPasswd string //backup password    【股东备份密码】
PublicKey string //app public key    【股东 app 公钥】
RecID     string //recovery app id    【申请恢复的股东 APP ID】
RecName   string //recovery app name   【申请恢复的股东 APP 名称】
RecPubKey string //recovery app pubKey   【申请恢复的股东 APP 公钥】
RecDecide int    //recovery decide     【恢复意见：(1)= 同意恢复， (非 1)= 拒绝恢复】
```

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 1.23 获取恢复股东认证结果

+ router: /api/v1/account/recoveryresult
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| app_id  | string |  申请人 app id  |
| reg_id  | int | 申请列表单条记录 id  |
| token  | string |  token  |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "List": [{
      "ID": 4,
      "Account": "gudong_3", // 认证者账户
      "OpAccount": "gudong_2", // 被认证者账户
      "Status": 2, // 0. 未审核 1. 同意 2. 拒绝
      "OpAppId": "gudong_appid_2",
      "RegId": 4,
      "UpdatedAt": "2018-09-27T14:25:28+08:00"
    }, {
      "ID": 5,
      "Account": "gudong_4",
      "OpAccount": "gudong_2",
      "Status": 0,
      "OpAppId": "gudong_appid_2",
      "RegId": 4,
      "UpdatedAt": "2018-09-27T17:51:49+08:00"
    }],
    "RecoveryStatus": 1 // 0. 注册 1. 已扫码审核中 2. 全部通过 3. 拒绝
  },
  "Msg": "success"
}
```

### 1.24 激活股东账户

+ router: /api/v1/account/activerecovery
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| reg_id  | int | 注册 id  |
| app_id  | string |  app id  |
| name  | string |  账户  |
| replaced  | string |  需要替换 msg 用户名和 msg  |
| token  | string |  token  |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |
| msg  | string  | 签名机整体加密参数 |

```bash
备注
msg:
Password  string //password      【股东密码，非登录口令】
BakPasswd string //backup password    【股东备份密码】
PublicKey string //app public key    【股东 app 公钥】
RecID     string //recovery app id    【申请恢复的股东 APP ID】
RecName   string //recovery app name   【申请恢复的股东 APP 名称】
RecPubKey string //recovery app pubKey   【申请恢复的股东 APP 公钥】
RecDecide int    //recovery decide     【恢复意见：(1)= 同意恢复， (非 1)= 拒绝恢复】
```

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 1.25 获取被恢复股东签名员工公钥

+ router: /api/v1/account/recpubkeys
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| name  | string | 账户  |
| sign  | string |    参数签名    |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
    "Name": "test12334",
    "Pubkey": "234489353jojergr"
  }],
  "Msg": "success"
}
```

### 1.26 待办任务统计

+ router: /api/v1/user/delaytasknum

+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| sign      |  string        | sign   |
| timestamp |   int     | 时间戳  |
| token      |  string        |token   |

+ 返回值

```json
{
  "Code": 0,
  "Msg": "success",
  "Data": {
    "Template": {
      "Number": 10,
      "Reason": ".. 审批流审批"
    },
    "Account": {
      "Number": 10,
      "Reason": "员工名称"
    },
    "Transfer": {
      "Number": 10,
      "Reason": ".. 员工申请转出 100BOX"
    }
  }
}
```

### 1.27 获取块高

+ router: /api/v1/user/blockheight

+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| sign      |  string        | sign   |
| timestamp |   int     | 时间戳  |
| token      |  string        |token   |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "AgrStatus": 0, // 0.正常 1.正在同步
    "List": [{
      "CurHeight": 0, // 当前
      "PubHeight": 0, // 最新
      "Status": 0, // 0.正常 1.正在同步
      "Name": "BTC"
    }, {
      "CurHeight": 1446344,
      "PubHeight": 1446344,
      "Status": 0,
      "Name": "ETH"
    }, {
      "CurHeight": 0,
      "PubHeight": 0,
      "Status": 0,
      "Name": "LTC"
    }]
  },
  "Msg": "success"
}
```

### 1.28 获取账户正在进行股东恢复的数量

+ router: /api/v1/account/hasrecovery

+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| sign      |  string        | sign   |
| timestamp |   int     | 时间戳  |
| name      |  string        | 恢复的股东账户名   |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "Count": 1
  },
  "Msg": "success"
}
```

## 2 部门管理

### 2.1 添加部门

+ router:  /api/v1/department/add
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |   备注     |
| :----: | :----: | :-----------------: |
| dep_name  | string |     部门名称       |
| token  | string |    token     |
| sign  | string |    参数签名     |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,// 成功
  "Data": {},
  "Msg":"success"
}
```

### 2.2 修改部门

+ router:  /api/v1/department/edit
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |    备注    |
| :----: | :----: | :------------: |

| id  | int |    部门 id     |
| dep_name | string | 部门名称 |
| token  | string |    token     |
| sign  | string |    参数签名     |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,// 成功
  "Data": {},
  "Msg":"success"
}
```

### 2.3 部门列表

+ router:  /api/v1/department/list
+ 请求方式： GET
+ 参数：

|   字段   |   类型   |    备注    |
| :----: | :----: | :--------------: |
| token  | string |      token      |
| sign  | string |     参数签名     |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": [
   {
      "Count": 1,
      "ID": 5,
      "Name": "部门 1"
    }
  ]
}
```

### 2.4 删除部门

+ router:  /api/v1/department/delete
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |     备注     |
| :----: | :----: | :--------------: |
| id  | int |    部门 id    |
| token  | string |   token       |
| sign  | string |      参数签名     |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,// 成功
  "Data": {},
  "Msg":"success"
}
```

### 2.5 部门排序

+ router:  /api/v1/department/sortdep
+ 请求方式： POST
+ 参数：

|   字段   |   类型   |     备注     |
| :----: | :----: | :--------------: |
| dep_id  | string | 排序后的部门 id，以【,】分隔   |
| token  | string |   token       |
| sign  | string |      参数签名     |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,// 成功
  "Data": {},
  "Msg":"success"
}
```

## 3 权限管理

### 3.1 获取权限列表

+ router: api/v1/auth/list`
+ 请求方式：GET
+ 参数：

| 字段 | 类型 | 备注 |
| ---- | ---- | ---- |
| token  | string | token |
| sign  | string | 参数签名 |
| timestamp | int | 时间戳 |

+ 返回:

```json
  {
    "Code": 0,
    "Data": [{
      "ID": 1,
      "Name": "查询",
      "Count": 2
    }],
    "Msg": "success"
  }
```

### 3.2 获取指定权限下用户列表

+ router: api/v1/auth/accountsbyauthid`
+ 请求方式：GET
+ 参数：

| 字段 | 类型 | 备注 |
| ---- | ---- | ---- |
| token  | string | token |
| auth_id  | int | 权限 id |
| sign  | string | 参数签名 |
| timestamp | int | 时间戳 |

+ 返回:

```json
  {
    "Code": 0,
    "Data": [{
      "ID": 1,
      "AppId:"123334",
      "Name": "admin"
    }],
    "Msg": "success"
  }
```

### 3.3 给指定用户添加权限

+ router: api/v1/auth/addauthtoaccount`
+ 请求方式：POST
+ 参数：

| 字段 | 类型 | 备注 |
| ---- | ---- | ---- |
| token  | string | token |
| auth_id  | int | 权限 id |
| add  | string | 添加权限的账户 id 以【,】分割 |
| sign  | string | 参数签名 |
| timestamp | int | 时间戳 |

+ 返回:

```json
  {
    "Code":0,
    "Data":{},
    "Msg":"success"
  }
```

### 3.4 删除指定用户权限

+ router: api/v1/auth/delauthfromaccount
+ 请求方式：POST
+ 参数：

| 字段 | 类型 | 备注 |
| ---- | ---- | ---- |
| token  | string | token |
| auth_id  | int | 权限 id |
| account_id  | int | 删除权限的 id |
| sign  | string | 参数签名 |
| timestamp | int | 时间戳 |

+ 返回:

```json
  {
    "Code":0,
    "Data": {},
    "Msg":"success"
  }
```

### 3.5 获取指定部门和权限下用户列表

+ router: api/v1/department/accountsbydepid
+ 请求方式：GET
+ 参数：

| 字段 | 类型 | 备注 |
| ---- | ---- | ---- |
| id  | int |  部门 id |
| auth_id  | int | 权限 id |
| token  | string | token |
| sign  | string | 参数签名 |
| timestamp | int | 时间戳 |

+ 返回:

```json
  {
    "Code": 0,
    "Data": [{
      "ID": 1,
      "AppId:"123334",
      "Name": "ll",
      "HasAuth": true,
      "UserType": 1, // 0. 普通成员 1. 管理员 2. 股东
    }],
    "Msg":"success"
  }
```

## 4 币种相关

### 4.1 币种列表

+ router: /api/v1/coin/list
+ 请求方式：GET
+ 参数：

|   字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|    type      | number |    0 主链币   1 代币  2 全部 |
| token  | string |                token                 |
| sign  | string |                参数签名                 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Data": [
           {
            "ID": 123,
            "Name": "BLS",// 名称
            "TokenType": 0, // 0 主链币  非零为代币
            "TokenAddress": "01XQE12EDQWE12631213",// 代币合约地址
            "FullName":"fffdx",// 币种全称
            "Status": 0
        }
    ],
    "Msg": "success"
}
```

### 4.2 添加代币

+ router: /api/v1/coin/add
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| name         | string |    全称     |
| symbol       | string |    简称     |
| decimals     |   int  |    精度     |
| address      |  string|    地址     |
| token        | string |   token     |
| sign         | string |   参数签名   |
| timestamp    |  int   | 时间戳 |

+ 返回值

```json
{
    "Code":0,
    "Data": {},
    "Msg":"success"
  }
```

### 4.3  查询余额

+ router: /api/v1/coin/balance
+ 请求方式：GET
+ 参数：

|    字段   |   类型     |   备注   |
| :------: | :--------------:| :------------: |
| token  | string |  token |
| sign  | string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Msg": "success",
    "data":[
        {
          "Name": "BTC",
          "Balance":300.22
        },
        {
          "Name": "ETH",
          "Balance":123
        }
    ]
  }
}
```

### 4.4  启用 / 禁用币种

+ router: /api/v1/coin/status

+ 请求方式：POST

+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|    id        | int |        id       |
|   status     | int |   0 启用  1 禁用 |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code":0,
    "Data": {},
    "Msg":"success"
  }
```

### 4.5 校验币种地址

+ router:/api/v1/coin/verifyaddress
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|    address|  string    |    地址   |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Msg": "success",
    "Data":{
        "Name" :"ethreum",
    "Symbol": "ETH",
    "Decimals": 18
    }
}
```

### 4.6  收付款码

+ router:/api/v1/coin/qrcode

+ 请求方式：GET

+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|    id      |       int       |   币种 ID             |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data":{
    "MainAddress":"0x123e1e1gr18f1fgg1f83g91",// 主账户地址
    "Sign" : "fdfdfdfdfdf",
    "Account":"zhangsan",
    "RandomAddress":"0x34e123rf1k3fb3ui1gr19eg 9"// 随机账户地址
    "RandomIndex":"[1,2]"
  },
  "Msg": "success"
}
```

## 5 多账户管理

### 5.1 主链币多账户统计

+ router: /api/v1/accmanage/statistical
+ 请求方式：GET
+ 返回值
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

```json
{
    "Code": 0,
    "Msg": "success",
    "Data": [
        {
            "ID": 124,
            "Name": "BTC",
            "Main": 3,// 主账户数量
            "Child": 2, // 子账户数量
            "HasToken": false// 是否含有代币，目前只有 ETH 有
        },
        {
            "ID": 125,
            "Name": "ETH",
            "Main": 2,
            "Child": 1,
            "HasToken": true// 是否含有代币，目前只有 ETH 有
        },
        {
            "ID": 126,
            "Name": "LTC",
            "Main": 1,
            "Child": 0,
            "HasToken": false// 是否含有代币，目前只有 ETH 有
        }
    ]
}
```

### 5.2 主链币对应账户分页列表

+ router: /api/v1/accmanage/listbypage
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   coin_id     |    int     |        币种 ID        |
|   condition   |   string   |   模糊查询 (地址或别名) |
|   page        |    int     |     页数  每页 20 条    |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Data": {
        "rows": [
            {
                "ID": 3,
                "Address": "0xhFad5d6d3569DF641570DEd06cb7A1b2Ccd112cG",// 地址
                "Tag": "v",// 别名
                "Balance": "11", // 余额
                "Type": 1 // 帐户类型  0 主账户  1 子账户
            },
            {
                "ID": 4,
                "Address": "0xqFfs5d6d3569DF678570DEd06cb7A1b2Ccd112cG",
                "Tag": "v",
                "Balance": "11",
                "Type": 0
            }
        ],
        "total": 22
    },
    "Msg": "success"
}
```

### 5.3 账户明细

+ router: /api/v1/accmanage/detail
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   id     |    int     |        账户 ID        |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Data": {
        "Total": 11,// 总额
        "Tag": "s", // 别名
        "Address": "0xSafs5d6d3569DF676530DEd06cb7A1b2Ccd112cG" // 地址
    },
    "Msg": "success"
}
```

### 5.4 转账记录

+ router: /api/v1/accmanage/transferrecord
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   id     |    int     |        账户 ID        |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Data": [
       {
            "ApplyReason": " ",// 交易信息
            "Status": 2,// 转账状态 0 审批中 1 转账中 2 转出成功 3 被驳回 4 转账失败 5 转入成功
            "Amount": "50",// 金额
            "CreatedAt": "2018-09-06 05:35:39 PM",// 时间
            "Type": 0, //  0  转入  1 转出
            "Link": "" // 外链查询转账信息
        }
    ],
    "Msg": "success"
}
```

### 5.5 代币明细

+ router: /api/v1/accmanage/tokenlist
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   id     |    int     |   账户 ID  |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Data": [
        {
            "ID": 1,
            "AddressId": 1,
            "CoinName": "BOX",// 代币名称
            "Balance": "12"// 余额
        },
        {
            "ID": 5,
            "AddressId": 1,
            "CoinName": "DFG",
            "Balance": "12"
        }
    ],
    "Msg": "success"
}
```

### 5.6 设置别名

+ router:  /api/v1/accmanage/settag
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   id     |    int        |        账户 ID        |
|   tag    |    string     |         别名         |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code":0,
    "Data": {},
    "Msg":"success"
  }
```

### 5.7 新增子账户

+ router: /api/v1/accmanage/addchildaccount
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   coin_id  |    int        |        币种 ID        |
|   tag      |    string     |         别名         |
|   amount   |    int        |         数量         |
|   psw      |    string     |         密码         |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code":0,
    "Data": {},
    "Msg":"success"
  }
```

### 5.8 指定代币转账记录

+ router: /api/v1/accmanage/tokenrecord
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   id       |    int     |        账户 ID        |
|   coinId   |    int     |        币种 ID        |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Data": [
        {
            "ApplyReason": "50",// 交易信息
            "Status": 2,// 转账状态 0 审批中 1 转账中 2 转出成功 3 被驳回 4 转账失败 5 转入成功
            "Amount": "50",// 金额
            "CreatedAt": "2018-09-06 05:35:39 PM",// 时间
            "Type": 0 //  0  转入  1 转出
        }
    ],
    "Msg": "success"
}
```

### 5.9 合并账户

+ router: /api/v1/accmanage/merge
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   coin_id   |    int     |        币种 ID             |
|   type      |    int     |  0 合并所有子账户 1 合并指定子账户 |
|   ids       |    string  |   type 为 1 时, 子账户的 id 集合 以逗号分开  |
|   psw      |    string     |         密码         |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

```json
{
    "Code": 0,
    "Msg": "success"
}
```

### 5.10 账户 top10 列表

+ router: /api/v1/accmanage/topten
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   coin_id     |    int     |        币种 ID        |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Data": {
        "rows": [
            {
                "ID": 3,
                "Address": "0xhFad5d6d3569DF641570DEd06cb7A1b2Ccd112c",// 地址
                "Tag": "v",// 别名
                "Balance": "11", // 余额
                "Type": 1 // 帐户类型  0 主账户  1 子账户
            },
            {
                "ID": 4,
                "Address": "0xqFfs5d6d3569DF678570DEd06cb7A1b2Ccd112c",
                "Tag": "v",
                "Balance": "11",
                "Type": 0
            }
        ],
        "total": 22
    },
    "Msg": "success"
}
```

### 5.11 子账户统计

+ router: /api/v1/accmanage/childstical
+ 请求方式：GET

+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   coin_id     |    int     |        币种 ID        |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Data": {
        "ChildCount":100, // 已建子账户数量
        "OtherCount":100,// 其他已建子账户数量
        "UsableCount":9800,// 剩余可建数量
        "MaxCount":10000  // 子账户数量上限
    },
    "Msg": "success"
}
```

### 5.12 计算子账户数量和总额

+ router: /api/v1/accmanage/countallchild
+ 请求方式：GET

+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   coin_id     |    int     |        币种 ID        |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Data": {
        "Balance":"633", // 总额
        "Amount":100,// 个数
    },
    "Msg": "success"
}
```

### 5.13 创建合约地址

+ router:/api/v1/accmanage/add/contract
+ 请求方式：POST
+ 参数：

|   字段    |  类型  |     备注     |
| :-------: | :----: | :----------: |
| timestamp |  int   |    时间戳    |
|   sign    | string | 参数签名信息 |
|   token   | string |    token     |

+ 返回值

```json
{
    "Code": 0,
    "Data": [
        {}
    ],
    "Msg": "success"
}
```

### 5.14 获取合约账户地址

+ router:/api/v1/accmanage/contract
+ 请求方式：GET
+ 参数：

|   字段    |  类型  |     备注     |
| :-------: | :----: | :----------: |
| timestamp |  int   |    时间戳    |
|   sign    | string | 参数签名信息 |
|   token   | string |    token     |

+ 返回值

```json
{
    "Code": 0,
    "Data": {
      "type": 1,// 0 - 未申请创建合约地址 1 - 合约地址创建中 2 - 合约地址已生成
      "result": {
        "address": "0xhFad5d6d3569DF655070DEd06cb7A1b2Ccd1D3AF",
        "balance": "422",
        "type": 0
      }
      "account": "sss", //对主地址签名的股东
      "sign":"dfdsafdafdsaf",//对主地址的签名
      "masterAddress" :"oxfdsafdsfasfadsfdasfasfsaf" //ETH主地址
    },
    "Msg": "success"
}
```

## 6 转账

### 6.1 转账申请明细

+ router: /api/v1/transfer/findapplybyid
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |
|   order_id     |    string     |        id        |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "CoinName": "ETH",
    "ApplyerName": "admin_3",
    "ApplyReason": "真正的交易",
    "Deadline": "2018-10-01T21:21:21+08:00",
    "CreatedAt": "2018-09-28T16:22:37+08:00",
    "OrderAddress": [{
      "Address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
      "Amount": [{
          "amount": "300",
          "status": 0
        },
        {
          "amount": "200",
          "status": 0
        },
        {
          "amount": "100",
          "status": 0
        }
      ],
      "Tag": "梭哈第一个"
    },{
      "Address": "0x39ed0b367a03a4c854dec433c98f15510183324b",
      "Amount": [{
        "amount": "50",
        "status": 0
      }, {
        "amount": "200",
        "status": 0
      }],
      "Tag": "梭哈第二个"
    }],
    "StrTemplate": "",
    "Amount": "1000",
    "Status": 8,
    "NowLevel": 1,
    "Msg": "{\n\t\"coin_name\": \"ETH\",\n\t\"coin_fullname\": \" 以太坊 \",\n\t\"precise\": 18,\n\t\"token\": \"0x2a65aca4d5fc5b5c859090a6c34d164135398226\",\n\t\"t_hash\": \"a4e4eabf44be356bad7d0f3e69409fd512262708886bde072a6c524265a526d8\",\n\t\"reason\": \" 真正的交易 \",\n\t\"deadline\": \"2018-10-01 21:21:21\",\n\t\"miner\": \"0.09\",\n\t\"timestamp\": 1923232328808,\n\t\"amount\": \"1000\",\n\t\"applys\": [{\n\t\t\t\"to_address\": \"0x2a65aca4d5fc5b5c859090a6c34d164135398226\",\n\t\t\t\"tag\": \" 梭哈第一个 \",\n\t\t\t\"amount\": [\"100\", \"200\", \"300\"]\n\t\t},\n\t\t{\n\t\t\t\"to_address\": \"0x39ed0b367a03a4c854dec433c98f15510183324b\",\n\t\t\t\"tag\": \" 梭哈第二个 \",\n\t\t\t\"amount\": [\"50\", \"150\", \"200\"]\n\t\t}\n\t]\n}",
    "ApprovalContent": {
      "admin_3": {
        "Status": 0,
        "sign": ""
      },
      "employee_2-1": {
        "Status": 0,
        "sign": ""
      }
    },
    "Template": {
      "name": "测试",
      "period": 0,
      "approvalInfo": [{
        "require": 1,
        "approvers": [{
          "account": "employee_2-1",
          "pubkey": "MIIBCgKCAQEAyRjJ2hJWv4uZpqcmpRKrz+KWhsF0FGc9oWFbSOr6sp+hrOqy0Ezb+osTGsvqQlvKTPZvLgB0tcsmu3z3IktwrirMn/q29iIBFLSdE5cKME2Hm/S6Z/oZYpEaz5ss56ADAXZT8FlB72NYOnPDcOXtiHRIDS2JGeyQtxAGSOJDIcGw+8qlIeKP0XpK4m8eQwlab672pQ3K55G5NtBg2j8OJawo0uswNMmKY/FoB6elrLUZpO9fObzCCHSZjTPOXi279uLI5u8ZnD8LmQ9T7q7GmoDn8Ds5rRpXT2JzNV4w8a7cz8oG9WXRhJGzj+6ujYXiJJLUFGffwacMh4bp8wcfVwIDAQAB"
        }, {
          "account": "admin_3",
          "pubkey": "MIIBCgKCAQEA2lbCdurvknHEAcq8D76m2Ax7zkO5PPQ25f0Et9n1AfKLNSDTlv9jecPsRLMQQtwxL7hKDE59wgPPoyjOVmzVseCnhsRTFL1pxjcsqMWqUOmYRJCMDxat2tKly3FSHOutlVYN2xipTxa/xZWxSSqlD1+12CgAahFor+VlzEmBqp8XjSrLkV9R+j6Imn+yV3Lsoma72LS9Tu2ulFGYLqRe29GJJpaCcwgfeISWUnVlYomUfNHGZb3xAx/RBQlFHSLcbElzyHEDzN1K9WLRCjPaYVSWLCryLQgFgJbvQ7OAvRf1+LdCJyKIfcE/YIYF7oJQXhxhS9z8KxHZw9FGoDYGYwIDAQAB"
        }]
      }, {
        "require": 2,
        "approvers": [{
          "account": "gudong_2",
          "pubkey": ""
        }, {
          "account": "admin_3",
          "pubkey": ""
        }]
      }],
      "limitInfo": [{
        "tokenAddress": "0xe1A178B681BD05964d3e3Ed33AE731577d9d96dD",
        "symbol": "BOX",
        "name": "BOX Token",
        "precise": 18,
        "limit": "10000"
      }]
    }
  },
  "Msg": "success"
}
```

### 6.2  转账申请列表

+ router: /api/v1/transfer/findapplylist
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   sort     |    int     | 1. 全部  2. 发起的 3. 参与的  |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
    "Amount": "1000",
    "ApplyReason": "真正的交易",
    "CoinName": "ETH",
    "CreatedAt": "2018-09-28T16:22:37+08:00",
    "OrderId": "2bf963d5-61fa-4c35-8233-0ee0527ec546",
    "Status": 0,
    "ApplyName": "bruce"
  }, {
    "Amount": "1000",
    "ApplyReason": "真正的交易",
    "CoinName": "ETH",
    "CreatedAt": "2018-09-28T16:46:22+08:00",
    "OrderId": "8ff6a491-7d4d-466f-9cd8-728788810f5d",
    "Status": 0 //0. 审批中。1 审批通过，转账中。2. 拒绝。 3 部分转账成功 4. 转账失败。5。全部成功 .6. 撤回 7. 非法 8. 审批过期 9. 转账过期 。10 作废
  }],
  "Msg": "success"
}
```

### 6.3  转账申请日志

+ router: /api/v1/transfer/findapplylog
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| sign      |  string     | sign |
| timestamp |   int    | 时间戳 |
| token      |  string       |token |
| order_id      |  string     | 订单 id  |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "ApproversOper": [{
      "OrderNum": "12",
      "Status": 0, // 审批意见 1 同意 2 拒绝 3 撤回 '
      "Encode": "", // 审批者对该笔订单的签名值
      "Reason": "test", // 撤回或拒绝需要填写原因
      "UpdatedAt": "2018-09-18T09:35:16+08:00", // 时间
      "AccountName": "sun" // 审批人
    }],
    "LastStatus": {
      "status": 1, //0. 审批中。1 审批通过，转账中。2. 拒绝。 3 部分转账成功 4. 转账失败。
      //5。全部成功 .6. 撤回 7. 非法 8. 审批过期 9. 转账过期 。10。员工作废 11. 模板停用作废
      "UpdatedAt": "2018-09-18T09:35:16+08:00", // 时间
    }
  },
  "Msg": "success"
}
```

```bash
转账的状态 ：//0. 审批中。1 审批通过，转账中。2. 拒绝。 3 部分转账成功 4. 转账失败。5。全部成功 .6. 撤回 7. 非法 8. 审批过期 9. 转账过期 。10。员工作废 11. 模板停用作废
```

### 6.4 获取所有币种

+ router:/api/v1/transfer/getcoins
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| sign      |  string      | sign |
| timestamp |   int             | 时间戳    |
| token      |  string   |token |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
    "ID": 127,
    "Name": "BOX",
    "FullName": "BOX Token",
    "Precise": 18,
    "Balance": "120",
    "TokenType": 125,
    "TokenAddress": "0xe1A178B681BD05964d3e3Ed33AE731577d9d96dD",
    "Available": 0,
    "Currency": "ETH"
  }, {
    "ID": 129,
    "Name": "USDT",
    "FullName": "泰达币",
    "Precise": 0,
    "Balance": "100",
    "TokenType": 0,
    "TokenAddress": "",
    "Available": 0,
    "Currency": "ETH"
  }],
  "Msg": "success"
}
```

### 6.5 通过币种 id 获取所有模板

+ router:/api/v1/transfer/gettemplate
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| sign      |  string   | sign |
| timestamp |   int      | 时间戳  |
| token      |  string    |token  |
| coin_id      |  string    | 币种 id  |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
    "Id": "aaaaa", // 模板 id
    "Name": "aaaa", // 模板名称
    "Hash": "fdsfsdf" // 模板哈希
  },{
    "Id": "bbbb",
    "Name": "vvv",
    "Hash": "ffefefe"
  }],
  "Msg": "success"
}
```

### 6.6 转账申请

+ router:/api/v1/transfer/apply
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| sign      |  string  |sign |
| timestamp |   int             | 时间戳   |
| token      |  string     | token  |
| apply_msg      |  string  | 申请转账的信息 |
| apply_sign      |  string   | 对申请转账的信息的签名 |
| password      |  string  | 密码 |

+ 注释 applys 的数据格式

```json
{
"apply_msg":{
  "currency": "ETH", // 主链币名称
  "coin_name": "ETH",
  "coin_fullname": "etheratum",
  "precise": 18,
  "token": "x2a65aca4d5fc5b5c859090a6c34d164135398226", // 代币必须填   主链不需要填
  "t_hash": "ec0910844cfc95fc99af4ca9208e777d19f07c97652c06f0ea9e62a29bbc5c57",
  "reason": "梭哈",
  "deadline": "",
  "miner": "0.09",
  "timestamp": 1923232323223,
  "amount": "1000",
  "applys": [{
    "to_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
    "tag": "梭哈第一个",
    "amount": ["100", "200", "300"]
  }, {
    "to_address": "0x39ed0b367a03a4c854dec433c98f15510183324b",
    "tag": "梭哈第二个",
    "amount": ["50", "150", "200"]
  }]
}
```

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 6.7 审批

+ router:/api/v1/transfer/verify
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| sign      |  string  | sign |
| timestamp |   int             | 时间戳  |
| token      |  string   |token |
| order_id |string  | 订单 id |
| status      |  int |1. 同意  2. 拒绝 |
| transfer_sign      |  string   | 对转账流的签名 |
| reason      |  string   | 拒绝原因  （选填）|
| password      |  string    | 密码 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 6.8 取消

+ router:/api/v1/transfer/cancel
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| sign      |  string  |sign |
| timestamp |   int   | 时间戳   |
| token      |  string  |token |
|order_id |string | 订单 id |
|password |string | 密码 |
|reason | string  | 撤销原因 |

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

### 6.9  批量审批的列表

+ router: /api/v1/transfer/batchlist
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| sign      |  string |  sign  |
| timestamp |   int             | 时间戳  |
| token      |  string  | token |
| order_ids |  string  |"fdsfad","fdsfdafdasf" 多个用，隔开 |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
    "CoinName": "ETH",
    "ApplyerName": "admin_3",
    "ApplyReason": "真正的交易",
    "Deadline": "2018-10-01T21:21:21+08:00",
    "CreatedAt": "2018-09-28T16:22:37+08:00",
    "OrderAddress": [{
        "Address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
        "Amount": [{
          "amount": "300",
          "status": 0
        }, {
          "amount": "200",
          "status": 0
        }],
        "Tag": "梭哈第一个"
      },
      {
        "Address": "0x39ed0b367a03a4c854dec433c98f15510183324b",
        "Amount": [{
          "amount": "50",
          "status": 0
        }, {
          "amount": "200",
          "status": 0
        }],
        "Tag": "梭哈第二个"
      }
    ],
    "StrTemplate": "",
    "Amount": "1000",
    "Status": 8,
    "NowLevel": 1,
    "Msg": "{\n\t\"coin_name\": \"ETH\",\n\t\"coin_fullname\": \" 以太坊 \",\n\t\"precise\": 18,\n\t\"token\": \"0x2a65aca4d5fc5b5c859090a6c34d164135398226\",\n\t\"t_hash\": \"a4e4eabf44be356bad7d0f3e69409fd512262708886bde072a6c524265a526d8\",\n\t\"reason\": \" 真正的交易 \",\n\t\"deadline\": \"2018-10-01 21:21:21\",\n\t\"miner\": \"0.09\",\n\t\"timestamp\": 1923232328808,\n\t\"amount\": \"1000\",\n\t\"applys\": [{\n\t\t\t\"to_address\": \"0x2a65aca4d5fc5b5c859090a6c34d164135398226\",\n\t\t\t\"tag\": \" 梭哈第一个 \",\n\t\t\t\"amount\": [\"100\", \"200\", \"300\"]\n\t\t},\n\t\t{\n\t\t\t\"to_address\": \"0x39ed0b367a03a4c854dec433c98f15510183324b\",\n\t\t\t\"tag\": \" 梭哈第二个 \",\n\t\t\t\"amount\": [\"50\", \"150\", \"200\"]\n\t\t}\n\t]\n}",
    "ApprovalContent": {
      "admin_3": {
        "Status": 0,
        "sign": ""
      },
      "employee_2-1": {
        "Status": 0,
        "sign": ""
      }
    },
    "Template": {
      "name": "测试",
      "period": 0,
      "approvalInfo": [{
          "require": 1,
          "approvers": [{
            "account": "employee_2-1",
            "pubkey": "MIIBCgKCAQEAyRjJ2hJWv4uZpqcmpRKrz+KWhsF0FGc9oWFbSOr6sp+hrOqy0Ezb+osTGsvqQlvKTPZvLgB0tcsmu3z3IktwrirMn/q29iIBFLSdE5cKME2Hm/S6Z/oZYpEaz5ss56ADAXZT8FlB72NYOnPDcOXtiHRIDS2JGeyQtxAGSOJDIcGw+8qlIeKP0XpK4m8eQwlab672pQ3K55G5NtBg2j8OJawo0uswNMmKY/FoB6elrLUZpO9fObzCCHSZjTPOXi279uLI5u8ZnD8LmQ9T7q7GmoDn8Ds5rRpXT2JzNV4w8a7cz8oG9WXRhJGzj+6ujYXiJJLUFGffwacMh4bp8wcfVwIDAQAB"
          }, {
            "account": "admin_3",
            "pubkey": "MIIBCgKCAQEA2lbCdurvknHEAcq8D76m2Ax7zkO5PPQ25f0Et9n1AfKLNSDTlv9jecPsRLMQQtwxL7hKDE59wgPPoyjOVmzVseCnhsRTFL1pxjcsqMWqUOmYRJCMDxat2tKly3FSHOutlVYN2xipTxa/xZWxSSqlD1+12CgAahFor+VlzEmBqp8XjSrLkV9R+j6Imn+yV3Lsoma72LS9Tu2ulFGYLqRe29GJJpaCcwgfeISWUnVlYomUfNHGZb3xAx/RBQlFHSLcbElzyHEDzN1K9WLRCjPaYVSWLCryLQgFgJbvQ7OAvRf1+LdCJyKIfcE/YIYF7oJQXhxhS9z8KxHZw9FGoDYGYwIDAQAB"
          }]
        },
        {
          "require": 2,
          "approvers": [{
            "account": "gudong_2",
            "pubkey": ""
          }, {
            "account": "admin_3",
            "pubkey": ""
          }]
        }
      ],
      "limitInfo": [{
        "tokenAddress": "0xe1A178B681BD05964d3e3Ed33AE731577d9d96dD",
        "symbol": "BOX",
        "name": "BOX Token",
        "precise": 18,
        "limit": "10000"
      }, {
        "tokenAddress": "",
        "symbol": "USDT",
        "name": "泰达币",
        "precise": 0,
        "limit": "10000"
      }]
    }
  }],
  "Msg": "success"
}
```

### 6.10  批量审批

+ router: /api/v1/transfer/batchverify
+ 请求方式：POST
+ 参数：

| sign      |  string |sign |
| timestamp |   int            | 时间戳   |
| token      |  string      | token |
| password      |  string      | 密码 |
| order_ids      |  string  | 描述一下 |
| reason      |  string  | 拒绝原因  （选填）|

```bash
注释：
order_ids:
[{order_id:"fdsafasf",status:1,transfer_sign:"ssssss"}, //status : 1. 同意  2. 拒绝  3. 验证签名不通过

{order_id:"fdsafasf",status:1,transfer_sign:"saaaaaa"}]
```

+ 返回值

```json
{
  "Code": 0,
  "Data": {},
  "Msg": "success"
}
```

## 7 审批流

### 7.1 新建审批流模板

+ router:/api/v1/template/new
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| content    |   string        | 审批流模板内容 |
| template_sign | string | 对审批流模板 content 的签名 |
| token  |   string         | token |
| timestamp | int | 时间戳 |
| sign |   string        | 参数签名信息   |

+ 返回值

```json
{
    "Code": 0,
    "Msg": "success",
    "Data": [""] // 审批流模板 ID string
}
```

+ 备注

>`content` 字段对应的结构为：

```json
{
  "name" :      // 审批流模板内容, string
  "period" :     // 预设的转账额度冻结周期, int 0~240
  "approvalInfo" : [  // 审批流层级信息
    {
      "require" :      // 本级所需最低审批通过人数, int
      "approvers" : [
        {
          "account" :       // 审批者账号, string
          "pubkey" :       // 审批者公钥, string(base64)
        },
  ...
      ]
    },
 ...
  ],
  "limitInfo" : [   // 该审批流对应的币种限额信息
    {
      "tokenAddress":   // 代币合约地址，主链币则为空，string
      "symbol" :     // 币种简称, string
      "name":      // 币种全称, string
      "precise":    // 精度，int
      "limit" :     // 额度, string
    },
 ...
  ]
}
```

### 7.2 审批审批流模板

+ router:/api/v1/template/verify
+ 请求方式：POST
+ 参数：

|     字段      |  类型  |         备注         |
| :-----------: | :----: | :------------------: |
|  template_id  | string |     审批流模板 ID     |
|    status     |  int   | 审批意见 1 同意 2 拒绝 |
| template_sign | string |  对审批流模板的签名  |
|    aeskey     | string |                      |
|   timestamp   |  int   |        时间戳        |
|     token     | string |        token         |
|     sign      | string |     参数签名信息     |

+ 返回值

```json
{
    "Code": 0,
    "Msg": "success"
}
```

### 7.3 作废审批流模板

+ router:/api/v1/template/cancel
+ 请求方式：POST
+ 参数：

|     字段      |  类型  |        备注        |
| :-----------: | :----: | :----------------: |
|  template_id  | string |    审批流模板 ID    |
| template_sign |  int   | 对审批流模板的签名 |
|   timestamp   |  int   |       时间戳       |
|     sign      |  int   |    参数签名信息    |
|     token     | string |       token        |

+ 返回值

```json
{
    "Code": 0,
    "Msg": "success"
}
```

### 7.4 通过审批流模版 hash 得到模版内容

+ router: /api/v1/template/findcontentbyhash
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
|   hash     |    string     | 1. 模版 hash  |
|   token  |  string |  token |
|   sign  |   string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
    "Code": 0,
    "Data": {
        "Content": {
            "name": "测试",
            "period": 0,
            "ApprovalInfo": [// 审批流数组
                { // 第一步
                    "Approvers": [// 审批人员
                        {
                            "Account": "testa",// 名称
                            "Pubkey": "MIIBCgKCAQEAmtiFfFYAzoVrx8Yeue4KWVjOoYO6aIDuYrjbxgZiy1YRiXZesXmraLUNB3gwqEH1C5NltHdaaRd/yUJHD1Tefm8RzLazpJ+OYYcFr89ocRyfk+hUX8LgJyzKcCAIsjLEUGKiCje2N+V+1UeB2X8F6r+fKI9FEDBy8zbR9eQdP9AQqDi3eQiQutBvrLtDWTjlUhV83I7X0UF9AObJS+I89c/GmOMMNztTT1dHNy5ML4QQrOEnEJqzu+ad5/tyAecqw/QSJycNPQOSfz/cIA+G1hBTgeT+R/ekdCuTBtkb2AnXmKiq5d7uA3PeHyP8uNEabCmdqikarJhc1aEGTW5QTQIDAQAB",// 公钥
                            "Status": 0 // 审批状态 0 未审批 1 同意 2 拒绝 3. 无需操作
                        }
                        ...
                    ],
                    "Require": 1// 需要审批的人数
                },
                { // 第二步
                    "Approvers": [
                        {
                            "Account": "tests",
                            "Pubkey": "MIIBCgKCAQEA0BQE1q7X1BGPBQDP/dC7Asfi4ex8tPdXelCumewb4W4CkPhkDrN3MggxYPt791WaXK7Nb/VEveZVwvITq4Zlpl11JW0umrdf+83rRVnhQJDsvfsyC41gBW7Aooa8ZpOX/Ronm9aXeanZQA0397Hmt234hc7dkn887c+qz5WS4VQ6nnAwKlBr4yPAJeGPbdIyaTN4IZekt2Mck0blS1QdyK7HKTqNtznq2QcJ51PhUfakLUbd/HwIIAxfyjghXkwzvQ8AU3jHTjDCQxUKd8RoDbN8szOXQbYZKID9kNdhivp3SHvdHeTkhF50kBVOYdhxTuSd8fVl79EFbYPr3X9WMwIDAQAB",
                            "Status": 0
                        }
                        ...
                    ],
                    "Require": 1
                }
                ...
            ],
            "limitInfo": [
                {
                    "coin": "ETH",
                    "limit": "10000000000"
                }
            ]
        },
        "Id": "aaaaa",
        "Status": 4 // 0 待审批 1 审批通过 2 审批拒绝 3. 禁用 4. 上链
    },
    "Msg": "success"
}
```

### 7.5 获取模板列表

+ router:/api/v1/template/list
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| token  |   string         | token                          |
| sign |   string        | 参数签名信息   |
| types |   int        | 0.全部  1.等待审批   |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
      "Id": "aaaaa", // 模板 id
      "Name": "aaaa" // 模板名称
    },
    {
      "Id": "bbbb",
      "Name": "vvv"
    }
  ],
  "Msg": "success"
}
```

### 7.6 通过模板 id 获取详情

+ router:/api/v1/template/findtembyid
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| token  |   string         | token                          |
| sign |   string        | 参数签名信息   |
| template_id |   string        | 模板 id  |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "ID": "aaaaa", // 模板 id
    "Hash": "", // 模板哈希
    "Name": "aaaa", // 模板名称
    "CreatorID": 1, // 创建者用户 id
    "Content": "{\n  \"name\": \" 测试 \",\n  \"period\": 0,\n  \"approvalInfo\": [\n    {\n      \"require\": 1,\n      \"approvers\": [\n        {\n          \"account\": \"testa\",\n          \"pubkey\": \"MIIBCgKCAQEAmtiFfFYAzoVrx8Yeue4KWVjOoYO6aIDuYrjbxgZiy1YRiXZesXmraLUNB3gwqEH1C5NltHdaaRd\\/yUJHD1Tefm8RzLazpJ+OYYcFr89ocRyfk+hUX8LgJyzKcCAIsjLEUGKiCje2N+V+1UeB2X8F6r+fKI9FEDBy8zbR9eQdP9AQqDi3eQiQutBvrLtDWTjlUhV83I7X0UF9AObJS+I89c\\/GmOMMNztTT1dHNy5ML4QQrOEnEJqzu+ad5\\/tyAecqw\\/QSJycNPQOSfz\\/cIA+G1hBTgeT+R\\/ekdCuTBtkb2AnXmKiq5d7uA3PeHyP8uNEabCmdqikarJhc1aEGTW5QTQIDAQAB\"\n}\n      ]\n    },\n    {\n      \"require\": 1,\n      \"approvers\": [\n        {\n          \"account\": \"tests\",\n          \"pubkey\": \"MIIBCgKCAQEA0BQE1q7X1BGPBQDP\\/dC7Asfi4ex8tPdXelCumewb4W4CkPhkDrN3MggxYPt791WaXK7Nb\\/VEveZVwvITq4Zlpl11JW0umrdf+83rRVnhQJDsvfsyC41gBW7Aooa8ZpOX\\/Ronm9aXeanZQA0397Hmt234hc7dkn887c+qz5WS4VQ6nnAwKlBr4yPAJeGPbdIyaTN4IZekt2Mck0blS1QdyK7HKTqNtznq2QcJ51PhUfakLUbd\\/HwIIAxfyjghXkwzvQ8AU3jHTjDCQxUKd8RoDbN8szOXQbYZKID9kNdhivp3SHvdHeTkhF50kBVOYdhxTuSd8fVl79EFbYPr3X9WMwIDAQAB\"\n        }\n      ]\n    }\n  ],\n  \"limitInfo\": [\n    {\n      \"coin\": \"ETH\",\n      \"limit\": \"10000000000\"\n    }\n  ]\n}", // 模板内容
    "Status": 0, // 审批流模板审批进度 0 待审批 1 审批通过 2 审批拒绝 3. 禁用 4. 上链
    "Period": 0, // 预设的额度恢复时间
    "CreatedAt": "2018-09-06T15:11:26+08:00",
    "UpdatedAt": "2018-09-06T16:20:03+08:00",
    "ApproveStatus" : //0.带审批   1.同意  2.拒绝
  },
  "Msg": "success"
}
```

### 7.7  签名流模板列表

+ router:/api/v1/template/list
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| token    |   string    |  token    |
| sign    |   string    |   对参数的签名  |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
      "Id": "aaaaa", // 模板 id
      "Name": "aaaa" // 模板名称
    },
    {
      "Id": "bbbb",
      "Name": "vvv"
    }
  ],
  "Msg": "success"
}
```

### 7.8 通过 id 查询模板

+ router:/api/v1/template/
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| template_id    |   string    |   模板 id  |
| token    |   string    |  token  |
| sign    |   string    |         对参数的签名   |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "ID": "aaaaa", // 模板 id
    "Hash": "", // 哈希
    "Name": "aaaa", // 模板名称
    "CreatorID": 1, // 创建人 id
    "Content": "{
      "name": "测试", // 模板名称
      "period": 0, //int  单位：小时  0～240
      "approvalInfo": [{
          "require": 1, // 要求几个人审核
          "approvers": [{
            "account": "testa", // 审批人名称
            "pubkey": "MIIBCgKCAQEAmtiFfFYAzoVrx8Yeue4KWVjOoYO6aIDuYrjbxgZiy1YRiXZesXmraLUNB3gwqEH1C5NltHdaaRd\\/yUJHD1Tefm8RzLazpJ+OYYcFr89ocRyfk+hUX8LgJyzKcCAIsjLEUGKiCje2N+V+1UeB2X8F6r+fKI9FEDBy8zbR9eQdP9AQqDi3eQiQutBvrLtDWTjlUhV83I7X0UF9AObJS+I89c\\/GmOMMNztTT1dHNy5ML4QQrOEnEJqzu+ad5\\/tyAecqw\\/QSJycNPQOSfz\\/cIA+G1hBTgeT+R\\/ekdCuTBtkb2AnXmKiq5d7uA3PeHyP8uNEabCmdqikarJhc1aEGTW5QTQIDAQAB" // 公钥
          }]
        },
        {
          "require": 1, // 要求几个人审核
          "approvers": [{
            "account": "tests", // 审批人名称
            "pubkey": "MIIBCgKCAQEA0BQE1q7X1BGPBQDP/dC7Asfi4ex8tPdXelCumewb4W4CkPhkDrN3MggxYPt791WaXK7Nb\\/VEveZVwvITq4Zlpl11JW0umrdf+83rRVnhQJDsvfsyC41gBW7Aooa8ZpOX\\/Ronm9aXeanZQA0397Hmt234hc7dkn887c+qz5WS4VQ6nnAwKlBr4yPAJeGPbdIyaTN4IZekt2Mck0blS1QdyK7HKTqNtznq2QcJ51PhUfakLUbd/HwIIAxfyjghXkwzvQ8AU3jHTjDCQxUKd8RoDbN8szOXQbYZKID9kNdhivp3SHvdHeTkhF50kBVOYdhxTuSd8fVl79EFbYPr3X9WMwIDAQAB" // 公钥
          }]
        }
      ],
      "limitInfo": [{
        "coin": "ETH", // 审批流允许的币种
        "limit": "10000000000" // 限制额度
      }]
    }",
    "Status": 0, //0 待审批 2 审批拒绝 3 审批通过 4. 禁用
    "CreatedAt": "2018-09-06T15:11:26+08:00",
    "UpdatedAt": "2018-09-06T16:20:03+08:00"
  },
  "Msg": "success"
}
```

## 8 消息通知、操作日志

### 8.1 消息列表

+ router: /api/v1/user/insideletter
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| page      |  int            |分页页数从0开始  |
| sign      |  string        |sign   |
| timestamp |   int             |时间戳          |
| token      |  string        |token   |

+ 返回值

```json
{
  "Code": 0,
  "Data": {
    "List": [{
      "ID": 6,
      "Title": "管理员任命通知", //标题
      "Status": 0, // 状态 0 未读 1已读
      "Time": "2018-10-12T11:32:25Z", //时间
      "Content": "股东gudong_3已将您设置为管理员", //内容
      "Type": 0, //要跳转到的页面 0 不用跳转  1 审批流详情页 2 转账申请详情页
      "Param": "{'id:23'}"
    }],
    "UnReadNumber": 10 ,//未读消息总数量
    "Total": 20,
    "AdminCount": 2 // 管理员数量
  },
  "Msg": "success"
}
```

### 8.2 消息标记为已读

+ router: /api/v1/user/readletter
+ 请求方式：POST
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| id        |  int      |ID   |
| sign      |  string        |sign   |
| timestamp |   int       | 时间戳  |
| token      |  string        |token   |

+ 返回值

```json
{
  "Code": 0,
  "Msg": "success"
}
```

### 8.3 获取日志

+ router: /api/v1/logger/logs
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| start  | int | 查询位置  |
| limit | int |  查询数量    |
| log_type | string | voucher：签名机， recovery: 股东恢复 |
| sign | string |  参数签名 |
| timestamp | int | 时间戳 |

+ 返回值

```json
{
  "Code": 0,
  "Data": [{
    "ID": 13,
    "Detail": "授权码被 cfo 扫描",
    "Note": "",
    "Operator": "cfo",
    "LogType": "voucher",
    "CreatedAt": "2018-10-11T13:29:18+08:00"
  }],
  "Msg": "success"
}
```

### 9.1 版本号

+ router: /api/version
+ 请求方式：GET
+ 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |
| 无  |

+ 返回值

```json
{
    "Code": 0,
    "Data": {
        "TreeVersion": "1",  //树的版本号
        "content": [        //更新内容
            "fdsfsdafsafdsafasf",
            "fdsfdsafsafsdf"
        ],
        "version": "1.0.0"  //接口版本号
    },
    "Msg": "success"
}
```


### 10.1 获取所有币种
- router:/api/v1/webtransfer/getcoins
- 请求方式：GET
- 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |

| 无		|




- 返回值

```json
{
    "Code": 0,
  "Data": [
        {
            "ID": 124,
            "Name": "BTC",
            "FullName": "比特币",
            "Precise": 8,
            "Balance": "300",
            "TokenType": 0,
            "TokenAddress": "",
            "Available": 0,
            "Currency":"BTC"
            "CurrencyId":1
        },
        {
            "ID": 125,
            "Name": "ETH",
            "FullName": "以太坊",
            "Precise": 18,
            "Balance": "499",
            "TokenType": 0,
            "TokenAddress": "",
            "Available": 0,
            "Currency":"ETH",
              "CurrencyId":2
        },
        {
            "ID": 126,
            "Name": "LTC",
            "FullName": "莱特币",
            "Precise": 18,
            "Balance": "200",
            "TokenType": 0,
            "TokenAddress": "",
            "Available": 0,
             "Currency":"ETH"
        },
        {
            "ID": 127,
            "Name": "BOX",
            "FullName": "BOX Token",
            "Precise": 18,
            "Balance": "120",
            "TokenType": 125,
            "TokenAddress": "0xe1A178B681BD05964d3e3Ed33AE731577d9d96dD",
            "Available": 0,
             "Currency":"ETH"
        },
        {
            "ID": 129,
            "Name": "USDT",
            "FullName": "泰达币",
            "Precise": 0,
            "Balance": "100",
            "TokenType": 0,
            "TokenAddress": "",
            "Available": 0,
             "Currency":"ETH"
        }
    ],
    "Msg": "success"
}
```


### 10.2 通过币种id获取所有模板
- router:/api/v1/webtransfer/gettemplate
- 请求方式：GET
- 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |

| coin_id      |  string       	|币种id			|

- 返回值

```json
{
    "Code": 0,
    "Data": [
        {
            "Id": "aaaaa",	//模板id
			"Name": "aaaa",	//模板名称
            "Hash": "fdsfsdf"	//模板哈希
        },
        {
            "Id": "bbbb",
            "Name": "vvv",
            "Hash": "ffefefe"
        }
    ],
    "Msg": "success"
}
```

### 10.3 提交转账信息
- router:/api/v1/webtransfer/transfercommit
- 请求方式：POST
- 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |

| order_id      |  string       	|id			|
| msg |   string             |加密后的转账信息                                 		|

注释：msg:
```json
"msg": {
    "currency": "ETH", //主链币名称
    "coin_name": "ETH",
    "coin_fullname":"etheratum",
    "precise":18,
    "token":"x2a65aca4d5fc5b5c859090a6c34d164135398226", //代币必须填   主链不需要填
	"t_hash": "ec0910844cfc95fc99af4ca9208e777d19f07c97652c06f0ea9e62a29bbc5c57",
	"reason": "梭哈",
	"deadline": "", //如 2018-09-09 11:11:11
	"miner": "0.09", //矿工费
	"timestamp": 1923232323223, //当前时间戳
	"amount": "1000", //总量
	"coin_id":1,   //币种id
	"applys": [{
			"to_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
			"tag": "",
			"amount": "100"
		},
		{
			"to_address": "0x39ed0b367a03a4c854dec433c98f15510183324b",
			"tag": "梭哈第二个",
			"amount": "50"
		}
	]
}
```


- 返回值

```json
{
    "Code": 0,
    "Msg": "success"
}
```


### 10.4 通过id查询转账信息
- router:/api/v1/transfer/findtranfersbyid
- 请求方式：GET
- 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |

| id      |  string       	|id			|



- 返回值

```json
{
    "Code": 0,
    "Data": {
        "ID": 1,
        "TransferId": "fdsafafs",
        "Msg": "fefefefef"
    },
    "Msg": "success"
}
```

### 10.5 查询web路由
- router:/api/v1/webtransfer/getwebrouter
- 请求方式：GET
- 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |

|无	|



- 返回值

```json
{
    "Code": 0,
    "Data": {
        "router": "aaa" //路由
    },
    "Msg": "success"
}
```

### 10.6 web查询提交状态
- router:/api/v1/webtransfer/getcommitstatus
- 请求方式：GET
- 参数：

|    字段   |        类型     |         备注         |
| :------: | :--------------:| :-----------------: |

|id	        |      string     |            id                   |



- 返回值

```json
{
    "Code": 0,
    "Data": {
        "status": 2    // 1.没有提交  2 提交
    },
    "Msg": "success"
}
```







