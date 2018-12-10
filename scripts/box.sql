# Copyright 2018. bolaxy.org authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ************************************************************
# Database: boxapi
# Generation Time: 2018-10-24 01:47:30 +0000
# ************************************************************
# 数据库更新log
# 修改人      修改内容       更新时间
# ************************************************************

# Dump of table account
# 用户账户表
# ------------------------------------------------------------

#1.bruce    transferReview 表  reason  字段 长度加到100    2018-10-30

#2.bruce    configs表 加默认数据address_sign 2018-10-31

#3.eileen   ownerReg表 添加 msg 字段  2018-11-18
#4.david    message表 添加 padding字段 2018-12-06

DROP TABLE IF EXISTS `account`;

CREATE TABLE `account` (
  `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `appId`         varchar(100) NOT NULL                       COMMENT 'appid',
  `name`          varchar(20) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT '姓名',
  `pwd`           varchar(100) CHARACTER SET utf8 COLLATE utf8_bin  NULL COMMENT '密码',
  `salt`          varchar(100) DEFAULT NULL                   COMMENT '密码盐',
  `pubKey`        varchar(1000) NOT NULL                      COMMENT '公钥',
  `isDeleted`     int(2) DEFAULT '0'                          COMMENT '0.未删除 1.删除',
  `msg`           text                                        COMMENT '签名信息',
  `departmentId`  int(2) DEFAULT '1'                          COMMENT '部门id',
  `userType`      int(2) DEFAULT '0'                          COMMENT '0.普通员工 1.管理员 2.股东',
  `frozen`        int(2) DEFAULT '0'                          COMMENT '1.冻结 0.未冻结',
  `attempts`      int(2) DEFAULT '0'                          COMMENT '输错密码次数',
  `frozenTo`      timestamp NULL DEFAULT NULL                 COMMENT '冻结到什么时候',
  `sourceAppId`   varchar(100) DEFAULT NULL                   COMMENT '上级appid',
  `level`         int(10) DEFAULT '0'                         COMMENT '第几层',
  `createdAt`     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt`     timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`appId`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table address
# 地址表
# ------------------------------------------------------------

DROP TABLE IF EXISTS `address`;

CREATE TABLE `address` (
  `id`        bigint(20) NOT NULL AUTO_INCREMENT,
  `address`   varchar(50) NOT NULL DEFAULT ''                   COMMENT '地址',
  `type`      int(1) NOT NULL DEFAULT '0'                       COMMENT '地址类型 0.主账户地址 1.普通地址 2.合约地址',
  `coinId`    int(10) NOT NULL                                  COMMENT '所属币种类别，对应coin表ID',
  `coinName`  varchar(10) DEFAULT NULL                          COMMENT '币种名称',
  `tag`       varchar(50) DEFAULT NULL                          COMMENT '别名',
  `tagIndex`  bigint(10) NOT NULL                               COMMENT '别名索引',
  `isDeleted` int(2) DEFAULT '0'                                COMMENT '0.未删除 1.删除',
  `usedBy`    int(10) DEFAULT NULL                              COMMENT '使用者账号ID',
  `deep`      varchar(100) DEFAULT NULL                         COMMENT '深度',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table auth
# 权限表
# ------------------------------------------------------------

DROP TABLE IF EXISTS `auth`;

CREATE TABLE `auth` (
  `id`        int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name`      varchar(50) DEFAULT NULL                              COMMENT '权限名称',
  `authType`  varchar(50) DEFAULT NULL                              COMMENT '权限标识',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Records of `auth`
-- ----------------------------
BEGIN;
INSERT INTO `auth` VALUES ('1', '资产查询', 'balance_auth'), ('2', '多地址管理', 'address_auth');
COMMIT;



# Dump of table authmap
# 账号关联权限表
# ------------------------------------------------------------

DROP TABLE IF EXISTS `authmap`;

CREATE TABLE `authmap` (
  `id`            int(11) NOT NULL AUTO_INCREMENT,
  `authId`        int(5) DEFAULT NULL                           COMMENT '权限id',
  `accountId`     int(10) DEFAULT NULL                          COMMENT '账号id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table capital
# 资金表
# ------------------------------------------------------------

DROP TABLE IF EXISTS `capital`;

CREATE TABLE `capital` (
  `id`          int(10) NOT NULL AUTO_INCREMENT                   COMMENT '资金ID,自增',
  `addressId`   int(10) NOT NULL                                  COMMENT '地址ID',
  `coinName`    varchar(10) NOT NULL DEFAULT ''                   COMMENT '币种名称',
  `balance`     varchar(100) DEFAULT NULL                         COMMENT '余额',
  `coinId`      int(10) DEFAULT NULL                              COMMENT '币种id',
  `address`     varchar(100) NOT NULL DEFAULT ''                   COMMENT '地址',
  `createdAt`   timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table coin
# 币种表
# ------------------------------------------------------------

DROP TABLE IF EXISTS `coin`;
CREATE TABLE `coin` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) DEFAULT NULL COMMENT '币种名称',
  `precise` int(3) DEFAULT '18' COMMENT '精度',
  `balance` varchar(100) DEFAULT '0' COMMENT '余额，总余额',
  `tokenType` int(3) DEFAULT NULL COMMENT '0.主链币  其他，代币',
  `tokenAddress` varchar(80) DEFAULT NULL COMMENT '代币合约地址',
  `available` int(2) DEFAULT '0' COMMENT '0.启用 1.禁用',
  `fullName` varchar(30) DEFAULT NULL COMMENT '币种全称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;


-- ----------------------------
--  Records of `coin`
-- ----------------------------
BEGIN;
INSERT INTO `coin` VALUES ('1', 'BTC', '8', '0', '0', null, '0', '比特币'), ('2', 'ETH', '18', '0', '0', null, '0', '以太坊'), ('3', 'LTC', '8', '0', '0', null, '0', '莱特币'), ('4', 'USDT', '8', '0', '0', '31', '0', '泰达币');
COMMIT;


# Dump of table configs
# 配置表
# ------------------------------------------------------------

DROP TABLE IF EXISTS `configs`;

CREATE TABLE `configs` (
  `id`        int(10) NOT NULL AUTO_INCREMENT,
  `con_key`   varchar(100) DEFAULT NULL                           COMMENT '配置的key',
  `con_value` text DEFAULT NULL                                   COMMENT '配置的值',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;



-- ----------------------------
--  Records of `configs`
-- ----------------------------
BEGIN;
INSERT INTO `configs` VALUES ('1', 'getbalance_time', ''),('2', 'tree_version', '1'), ('3', 'address_sign',''),('4', 'combine_account','0');
COMMIT;



# Dump of table department
# 部门表
# ------------------------------------------------------------

DROP TABLE IF EXISTS `department`;
CREATE TABLE `department` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '部门ID,自增',
  `name` varchar(20) DEFAULT NULL COMMENT '部门名称',
  `creatorAccId` int(10) DEFAULT NULL COMMENT '部门创建者账号ID',
  `available` int(2) NOT NULL DEFAULT '0' COMMENT '是否被删除 0.启用 1.禁用',
  `createdAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `order` int(10) DEFAULT '0' COMMENT '部门排序',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Records of `department`
-- ----------------------------
BEGIN;
INSERT INTO `department` VALUES ('1', '其他', '0', '0', now(), now(), '1');
COMMIT;



# Dump of table deposit
# 充值表
# ------------------------------------------------------------

DROP TABLE IF EXISTS `deposit`;

CREATE TABLE `deposit` (
  `id`            varchar(40) NOT NULL DEFAULT '',
  `coinId`        int(10) NOT NULL                        COMMENT '币种ID',
  `coinName`      varchar(10) NOT NULL DEFAULT ''         COMMENT '币种名称',
  `precise`       int(10) DEFAULT NULL                    COMMENT '币种精度',
  `fromAddr`      varchar(50) NOT NULL DEFAULT ''         COMMENT '付款地址',
  `toAddr`        varchar(50) NOT NULL DEFAULT ''         COMMENT '收款地址',
  `amount`        varchar(100) NOT NULL DEFAULT ''         COMMENT '充值金额',
  `txId`          varchar(100) NOT NULL DEFAULT ''        COMMENT '交易哈希',
  `blockNumber`   bigint(20) NOT NULL                     COMMENT '交易存储的区块编号',
  `confirm`       bigint(20) NOT NULL DEFAULT '0'         COMMENT '被确认数',
  `isUpdate`      int(2) DEFAULT '0'                      COMMENT '0.代币初始化 1.未更新 2.更新过',
  `tokenAddress`  varchar(100) DEFAULT NULL               COMMENT 'token地址',
  `createdAt`     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt`     timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table log
# 日志表
# ------------------------------------------------------------

DROP TABLE IF EXISTS `log`;

CREATE TABLE `log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `operator` varchar(20) DEFAULT NULL COMMENT '操作者',
  `detail` varchar(200) DEFAULT NULL COMMENT '操作内容',
  `note` char(200) DEFAULT NULL COMMENT '备注',
  `createdAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '日志记录时间',
  `logType` varchar(100) DEFAULT NULL COMMENT '日志类型',
  `pos` varchar(40) DEFAULT NULL COMMENT '次数标记',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table message
# 站内信
# ------------------------------------------------------------

DROP TABLE IF EXISTS `message`;

CREATE TABLE `message` (
  `id`        int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title`     varchar(100) NOT NULL DEFAULT ''                    COMMENT '标题',
  `content`   varchar(1000) NOT NULL DEFAULT ''                   COMMENT '内容模板',
  `padding`   varchar(1000) NULL                                   COMMENT '填充内容',
  `receiver`  varchar(1000) DEFAULT NULL                          COMMENT '接收人',
  `reader`    varchar(1000) DEFAULT NULL                          COMMENT '已阅读人员',
  `type`      int(11) DEFAULT '0'                                 COMMENT '要跳转到的页面  1.审批流详情页 2.转账申请详情页',
  `param`     varchar(200) DEFAULT NULL                           COMMENT '参数',
  `warnType`  int(2) DEFAULT 0                                    COMMENT '站内信类型 0.普通 1.签名机返回错误 2.签名机报警',
  `createdAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table ownerOper
# 股东恢复操作
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ownerOper`;

CREATE TABLE `ownerOper` (
  `id`                int(10) NOT NULL AUTO_INCREMENT,
  `account`           varchar(100) DEFAULT NULL                   COMMENT '用户名',
  `operatedAccount`   varchar(100) DEFAULT NULL                   COMMENT '被操作的用户',
  `status`            int(2) DEFAULT '0'                          COMMENT '0.审核中 1.同意 2.拒绝',
  `operatedAppid`     varchar(100) DEFAULT NULL                   COMMENT '被操作用户 appid',
  `regId`             int(10) DEFAULT NULL                        COMMENT '股东注册表关联 id',
  `updatedAt`         timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table ownerReg
# 股东注册
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ownerReg`;

CREATE TABLE `ownerReg` (
  `id`            int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name`          varchar(20) DEFAULT NULL                          COMMENT '申请人姓名',
  `pwd`           varchar(100) DEFAULT NULL                         COMMENT '登录密码',
  `salt`          varchar(100) DEFAULT NULL                         COMMENT '登录密码盐',
  `appId`         varchar(100) DEFAULT NULL                         COMMENT 'appid',
  `pubKey`        varchar(1000) DEFAULT NULL                        COMMENT '公钥',
  `status`        int(2) DEFAULT '0'                                COMMENT '0.注册 1.扫过码，等待审核 2.同意 3.拒绝 4.签名机未通过 5.失效',
  `refuseReason`  varchar(200) DEFAULT NULL                         COMMENT '拒绝原因',
  `msg`           text                                              COMMENT '签名信息',
  `createdAt`     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt`     timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table registration
# 普通账户注册
# ------------------------------------------------------------

DROP TABLE IF EXISTS `registration`;

CREATE TABLE `registration` (
  `id`              int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name`            varchar(20) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL                        COMMENT '申请人姓名',
  `pwd`             varchar(100) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL                       COMMENT '登录密码',
  `salt`            varchar(100) DEFAULT NULL                       COMMENT '登录密码盐',
  `sourceAppId`     varchar(100) DEFAULT NULL                       COMMENT '上级账号id',
  `sourceAccount`   varchar(100) DEFAULT NULL                       COMMENT '上级账号',
  `appId`           varchar(100) DEFAULT NULL                       COMMENT 'appid',
  `pubKey`          varchar(1000) DEFAULT NULL                      COMMENT '公钥',
  `status`          int(2) DEFAULT '0'                              COMMENT '0.注册 1.扫过码，等待审核 2.同意 3.拒绝',
  `msg`             varchar(1000) DEFAULT NULL                      COMMENT '上级对下级公钥签名',
  `refuseReason`    varchar(200) DEFAULT NULL                       COMMENT '拒绝原因',
  `level`           int(10) DEFAULT '0'                             COMMENT '第几层',
  `createdAt`       timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt`       timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table template
# 审批流模板
# ------------------------------------------------------------

DROP TABLE IF EXISTS `template`;

CREATE TABLE `template` (
  `id`          varchar(40) NOT NULL DEFAULT ''         COMMENT '业务流模板ID',
  `hash`        varchar(100) DEFAULT NULL               COMMENT '上链模板哈希值',
  `name`        varchar(100) NOT NULL DEFAULT ''        COMMENT '业务结构名称',
  `creatorId`   int(10) NOT NULL                        COMMENT '创建者账号ID',
  `content`     text                                    COMMENT '模板内容',
  `status`      int(1) NOT NULL DEFAULT '0'             COMMENT '审批流模板审批进度 0.待审批 1.审批通过 2.审批拒绝 3.申请禁用 4.上链 5.禁用成功 6签名机连接异常，需要重新上链',
  `period`      int(10) NOT NULL DEFAULT '0'            COMMENT '预设的额度恢复时间',
  `createdAt`   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt`   timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table templateOper
# 模板审批
# ------------------------------------------------------------

DROP TABLE IF EXISTS `templateOper`;

CREATE TABLE `templateOper` (
  `id`            int(10) NOT NULL AUTO_INCREMENT,
  `templateId`    varchar(100) NOT NULL                   COMMENT '模板id',
  `accountId`     int(10) DEFAULT NULL                    COMMENT '股东账号ID',
  `status`        int(2) DEFAULT NULL                     COMMENT '1.同意 2.拒绝',
  `sign`          varchar(1000) DEFAULT NULL              COMMENT 'app对内容模板哈希的签名',
  `appId`         varchar(100) DEFAULT NULL               COMMENT 'appid',
  `appName`       varchar(100) DEFAULT NULL               COMMENT '账户名',
  `reason`        varchar(20) DEFAULT NULL                COMMENT '拒绝原因',
  `createAt`      timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table templateView
# 模板信息
# ------------------------------------------------------------

DROP TABLE IF EXISTS `templateView`;

CREATE TABLE `templateView` (
  `templateId`    varchar(40) DEFAULT NULL              COMMENT '业务流模板ID',
  `coinId`        int(10) NOT NULL                      COMMENT '币种ID',
  `referAmount`   varchar(100) NOT NULL DEFAULT ''       COMMENT '当前可用转账额度',
  `amountLimit`   varchar(100) NOT NULL DEFAULT ''       COMMENT '预设的转账额度上限值',
  `referspan`     int(10) NOT NULL                      COMMENT '预设的冻结周期',
  `frozenTo`      timestamp NULL DEFAULT NULL           COMMENT '额度解冻时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table transfer
# 小订单
# ------------------------------------------------------------

DROP TABLE IF EXISTS `transfer`;

CREATE TABLE `transfer` (
  `id`          varchar(40) NOT NULL,
  `coinName`    varchar(10) NOT NULL                  COMMENT '币种名称',
  `amount`      varchar(100) NOT NULL DEFAULT ''       COMMENT '转账金额',
  `status`      int(2) NOT NULL DEFAULT '0'           COMMENT '转账状态 0.审批中 1.审批成功 2.审批失败 3.转账中 4.转账失败 5.转账成功 6.撤回 7.非法 8.审批过期 9.转账过期 10.员工作废 11.模板停用作废 12.余额不足',
  `txId`        varchar(100) DEFAULT NULL             COMMENT '交易哈希',
  `toAddress`   varchar(100) DEFAULT NULL             COMMENT '账户地址',
  `orderId`     varchar(40) DEFAULT NULL              COMMENT '大订单的id',
  `tag`         varchar(100) DEFAULT NULL             COMMENT '备注',
  `applyReason` varchar(200) DEFAULT NULL             COMMENT '申请原因',
  `amountIndex` int(10) DEFAULT NULL                  COMMENT 'amount索引',
  `types`       int(11) DEFAULT '1'                   COMMENT '1.普通转账 2.内部转账',
  `msg`         varchar(1000) DEFAULT NULL            COMMENT '转账信息',
  `accepted`    int(1) NOT NULL DEFAULT '0'           COMMENT '是否被SE记录 0.未被记录 1.已被记录',
  `createdAt`   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt`   timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `fromAddress`  varchar(100) DEFAULT NULL            COMMENT '发送地址',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table transferOrder
# 大订单
# ------------------------------------------------------------

DROP TABLE IF EXISTS `transferOrder`;

CREATE TABLE `transferOrder` (
  `id`            varchar(40) NOT NULL,
  `coinName`      varchar(40) DEFAULT NULL            COMMENT '币种名称',
  `hash`          varchar(1000) DEFAULT NULL          COMMENT '模板哈希',
  `applyReason`   varchar(200) DEFAULT NULL           COMMENT '申请理由',
  `miner`         varchar(30) DEFAULT NULL            COMMENT '矿工费',
  `sign`          varchar(1000) DEFAULT NULL          COMMENT '申请人对转账内容的签名',
  `status`        int(2) DEFAULT '0'                  COMMENT '0.审批中 1.审批通过，转账中 2.拒绝 3.部分转账成功 4.转账失败 5.全部成功 6.撤回 7.非法 8.审批过期 9.转账过期 10.员工作废 11.模板停用作废',
  `content`       text                                COMMENT '转账内容',
  `applyerId`     int(11) DEFAULT NULL                COMMENT '申请人账号id',
  `amount`        varchar(100) DEFAULT NULL            COMMENT '转账金额',
  `nowLevel`      int(10) DEFAULT '1'                 COMMENT '目前到第几层审批',
  `approversSign` text                                COMMENT '所有审批人的签名',
  `deadline`      timestamp NULL DEFAULT NULL         COMMENT '截止时间',
  `createdAt`     timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt`     timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table transferReview
# 转账审批
# ------------------------------------------------------------

DROP TABLE IF EXISTS `transferReview`;

CREATE TABLE `transferReview` (
  `id`          bigint(20) NOT NULL AUTO_INCREMENT,
  `status`      int(1) NOT NULL DEFAULT '0'             COMMENT '审批意见 0.未审批 1.同意 2.拒绝 3.无需操作',
  `encode`      varchar(1000) DEFAULT NULL              COMMENT '审批者对该笔订单的签名值',
  `reason`      varchar(100) DEFAULT NULL                COMMENT '撤回或拒绝需要填写原因',
  `orderNum`    varchar(40) DEFAULT NULL                COMMENT '订单id',
  `accountName` varchar(100) DEFAULT NULL               COMMENT '审批者账户名',
  `level`       int(11) DEFAULT '0'                     COMMENT '第几级',
  `updatedAt`   timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



-- ----------------------------
--  Table structure for `tx_infos`
-- ----------------------------
DROP TABLE IF EXISTS `tx_infos`;
CREATE TABLE `tx_infos` (
  `txId` varchar(255) COLLATE utf8_bin NOT NULL,
  `confirm` bigint(20) unsigned DEFAULT NULL,
  `token` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `height` bigint(20) unsigned DEFAULT NULL,
  `target` bigint(20) unsigned DEFAULT NULL,
  `fee` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `tx_obj` json DEFAULT NULL,
  `ext_valid` tinyint(1) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deletedd_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`txId`),
  KEY `idx_tx_infos_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- ----------------------------
--  Table structure for `webTransfers`
-- ----------------------------
DROP TABLE IF EXISTS `webTransfers`;
CREATE TABLE `webTransfers` (
  `id` bigint(10) NOT NULL AUTO_INCREMENT,
  `transferId` varchar(50) DEFAULT NULL COMMENT '订单id',
  `msg` text COMMENT '内容',
  `createdTime` bigint(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniid` (`transferId`) USING HASH
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



