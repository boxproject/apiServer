// Copyright 2018. bolaxy.org authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package errors

import "errors"

type CodeType int

// 各种errModel的code码
const (
	Success         CodeType = 0
	ErrToken        CodeType = 1 //非法tokens
	Code_2          CodeType = 2 //签名错误
	Code_3          CodeType = 3 //无管理员权限
	VerifyParamFail CodeType = 4 // 参数签名错误
	Code_5          CodeType = 5 // 时间戳超时，需要和服务器校准时间
	Code_6          CodeType = 6 // 请求签名机失败

	Code_102       CodeType = 102 //Verify error 				【验证签名错误】
	Code_105       CodeType = 105 //password error				【密码错误】
	Code_106       CodeType = 106 //app name not math			【app名称不匹配】
	Code_107       CodeType = 107 //same app name				【相同的app名称】
	VoucherTimeout CodeType = 121 //签名机超时

	Failed        CodeType = 1000 // 系统错误
	ParamsNil     CodeType = 1001 // 参数不能为空
	Duplicate     CodeType = 1002 // 数据重复
	AddDepFailed  CodeType = 1003 // 创建部门失败
	EditDepFailed CodeType = 1004 // 更新部门失败
	SetDepFailed  CodeType = 1005 // 部门不存在或设置失败
	VoucherFail   CodeType = 1006 // 签名机状态异常
	LogFailed     CodeType = 1007 // 操作日志添加失败
	DelDepFailed  CodeType = 1008 // 部门下还有成员，不能删除
	WallertFailed CodeType = 1009 // 调用钱包服务出错
	MessageFailed CodeType = 1010 // 插入站内信消息失败
	MessageRead   CodeType = 1011 // 消息已读

	User_2001 CodeType = 2001 // 注册用户失败
	User_2002 CodeType = 2002 // 用户名已存在
	User_2004 CodeType = 2004 // appid已存在
	User_2005 CodeType = 2005 // 重新注册失败
	User_2006 CodeType = 2006 // 获取注册用户列表失败
	User_2007 CodeType = 2007 // 账户恢复失败
	User_2008 CodeType = 2008 // 股东不存在，请核对姓名
	User_2009 CodeType = 2009 // 提交审批失败
	User_2010 CodeType = 2010 // 不能重复提交
	User_2011 CodeType = 2011 // 需重新注册或扫码提交认证
	User_2012 CodeType = 2012 // 认证股东失败
	User_2013 CodeType = 2013 // 部门排序参数有误
	User_2014 CodeType = 2014 // 不能提交自己的恢复申请
	User_2015 CodeType = 2015 // 该申请已被拒绝

	User_3001 CodeType = 3001 // 账号不存在
	User_3002 CodeType = 3002 // 密码不正确
	User_3003 CodeType = 3003 // 生成token错误
	User_3004 CodeType = 3004 // 校验用户不存在
	User_3005 CodeType = 3005 // 账户停用
	User_3006 CodeType = 3006 // 账户未解冻
	User_3007 CodeType = 3007 // 账户被锁定
	User_3008 CodeType = 3008 //审核，数据库状态不对
	User_3010 CodeType = 3010 //修改密码失败
	User_3011 CodeType = 3011 //添加管理员失败
	User_3012 CodeType = 3012 //删除管理员失败
	User_3013 CodeType = 3013 //已提交审核，无需重复提交
	User_3014 CodeType = 3014 //获取用户列表失败
	User_3015 CodeType = 3015 //停用账户失败
	User_3016 CodeType = 3016 //账户审批已通过，无需重复提交
	User_3017 CodeType = 3017
	User_3018 CodeType = 3018 // 股东恢复已交换公钥

	Auth_5002 CodeType = 5002 // 参数类型错误
	Auth_5003 CodeType = 5003 // 权限不存在
	Auth_5004 CodeType = 5004 // 账户不存在
	Auth_5005 CodeType = 5005 // 账户权限添加失败
	Auth_5006 CodeType = 5006 // 数据库操作错误
	Auth_5007 CodeType = 5007 // 无操作权限

	Coin_6001 CodeType = 6001 //币种已存在
	Coin_6003 CodeType = 6003 //地址错误
	Coin_6005 CodeType = 6005 //币种不存在
	Coin_6006 CodeType = 6006 //更新币种状态失败

	Acc_Manage_6008 CodeType = 6008 //添加子账户失败
	Acc_Manage_6009 CodeType = 6009 //为查到主账户
	Acc_Manage_6010 CodeType = 6010 //没有子账户
	Acc_Manage_6011 CodeType = 6011 //合并子账户失败
	Acc_Manage_6012 CodeType = 6012 //有账户正在进行交易，请完成后再尝试
	Acc_Manage_6013 CodeType = 6013 // 余额不足
	Acc_Manage_6014 CodeType = 6014 // 数量超过可添加子账号数量
	Acc_Manage_6016 CodeType = 6016 // 没有查询到地址信息
	Acc_Manage_6017 CodeType = 6017 // 没有查询到地址信息
	Acc_Manage_6018 CodeType = 6018 // 无法更改主账户地址别名
	Acc_Manage_6020 CodeType = 6020 // 正在合并
	Acc_Manage_6021 CodeType = 6021 //不支持该币种

	Template_9001 CodeType = 9001 // 审批流已存在

	Template_9002 CodeType = 9002 // 审批流不存在
	Template_9003 CodeType = 9003 // 审批审批流状态不对
	Template_9007 CodeType = 9007 // 创建合约失败
	Template_9008 CodeType = 9008 // 重复创建
	Template_9009 CodeType = 9009 // 审批流被篡改
	Template_9010 CodeType = 9010 // 该审批流已经被拒绝

	Transfer_10001 CodeType = 10001
	Transfer_10002 CodeType = 10002
	Transfer_10003 CodeType = 10003
	Transfer_10005 CodeType = 10005
	Transfer_10010 CodeType = 10010
	Transfer_10011 CodeType = 10011
	Transfer_10013 CodeType = 10013
	Transfer_10014 CodeType = 10014
	Transfer_10015 CodeType = 10015
	Transfer_10016 CodeType = 10016
	Transfer_10017 CodeType = 10017
	Transfer_10018 CodeType = 10018
	Transfer_10019 CodeType = 10019
	Transfer_10020 CodeType = 10020
	Transfer_10021 CodeType = 10021
	Transfer_10022 CodeType = 10022
	Transfer_10023 CodeType = 10023
	Transfer_10026 CodeType = 10026
	Transfer_10027 CodeType = 10027
	Transfer_10028 CodeType = 10028
	Transfer_10029 CodeType = 10029
	Transfer_10032 CodeType = 10032
	Transfer_10033 CodeType = 10033
	Transfer_10034 CodeType = 10034
	Transfer_10035 CodeType = 10035
)

var (
	MSG_2   = errors.New("签名错误.")
	MSG_3   = errors.New("无管理员权限.")
	MSG_4   = errors.New("参数签名错误.")
	MSG_5   = errors.New("请求超时，请校准时间.")
	MSG_6   = errors.New("请求签名机失败")
	MSG_102 = errors.New("验证签名错误")
	MSG_105 = errors.New("密码错误")
	MSG_106 = errors.New("app名称不匹配")
	MSG_107 = errors.New("相同的app名称")

	PARAMS_NULL  = errors.New("参数错误")
	DEP_EXISTS   = errors.New("部门已存在.")
	SYSTEM_ERROR = errors.New("系统错误.")
	MSG_ErrToken = errors.New("非法token.")
	MSG_1003     = errors.New("部门创建失败")
	MSG_1004     = errors.New("部门更新失败")
	MSG_1005     = errors.New("部门不存在或部门设置失败")
	MSG_1006     = errors.New("签名机状态异常")
	MSG_1007     = errors.New("操作日志添加失败")
	MSG_1008     = errors.New("部门下还有成员，无法删除")
	MSG_1009     = errors.New("调用钱包服务出错")
	MSG_1010     = errors.New("插入站内信消息失败")
	MSG_1011     = errors.New("该消息已读")

	MSG_2001 = errors.New("注册用户失败")     // 注册用户失败
	MSG_2002 = errors.New("用户名已存在.")    // 用户名已存在
	MSG_2004 = errors.New("appid已存在.")  // appid已存在
	MSG_2005 = errors.New("重新注册失败")     // 重新注册失败
	MSG_2006 = errors.New("获取注册用户列表失败") // 获取注册用户列表失败
	MSG_2007 = errors.New("账户恢复失败")
	MSG_2008 = errors.New("股东不存在，请核对姓名")
	MSG_2009 = errors.New("提交审批失败")
	MSG_2010 = errors.New("不能重复提交")
	MSG_2011 = errors.New("需重新注册或扫码提交认证")
	MSG_2012 = errors.New("认证股东失败")
	MSG_2013 = errors.New("部门排序参数有误")
	MSG_2014 = errors.New("不能提交自己的恢复申请")
	MSG_2015 = errors.New("该申请已被拒绝")

	MSG_3001 = errors.New("账号不存在.")
	MSG_3002 = errors.New("密码不正确.")
	MSG_3003 = errors.New("生成token错误.")

	MSG_3004 = errors.New("用户不存在.")
	MSG_3005 = errors.New("账户停用.")
	MSG_3006 = errors.New("账户未解除冻结")
	MSG_3008 = errors.New("审核失败.")
	MSG_3010 = errors.New("修改密码失败.")
	MSG_3011 = errors.New("管理员添加失败")
	MSG_3012 = errors.New("管理员删除失败")
	MSG_3013 = errors.New("已提交审核，无需重复提交")
	MSG_3014 = errors.New("获取用户列表失败")
	MSG_3015 = errors.New("停用账户失败")
	MSG_3016 = errors.New("账户审批已通过，无需重复提交")
	MSG_3017 = errors.New("该账户已被其他人审批")
	MSG_3018 = errors.New("股东恢复已交换公钥")

	MSG_5002 = errors.New("参数类型错误")
	MSG_5003 = errors.New("权限不存在")
	MSG_5004 = errors.New("账户不存在")
	MSG_5005 = errors.New("账户权限添加失败")
	MSG_5006 = errors.New("数据库操作错误")
	MSG_5007 = errors.New("无操作权限")

	MSG_6001 = errors.New("该币种已存在.")
	MSG_6003 = errors.New("地址不合法.")
	MSG_6005 = errors.New("该币种不存在.")
	MSG_6006 = errors.New("更新币种状态失败.")
	MSG_6008 = errors.New("添加子账户失败.")

	MSG_6009 = errors.New("未查询到主账户.")
	MSG_6010 = errors.New("合并子账户失败.")
	MSG_6011 = errors.New("没有可合并的子账户")
	MSG_6012 = errors.New("有账户正在进行交易，请完成后再尝试")
	MSG_6013 = errors.New("余额不足")
	MSG_6014 = errors.New("数量超过可添加子账号数量")
	MSG_6016 = errors.New("没有查询到地址信息.")
	MSG_6017 = errors.New("目前账户合并仅支持ETH.")
	MSG_6018 = errors.New("无法修改主账户地址别名.")

	MSG_6020 = errors.New("账户正在合并.")
	MSG_6021 = errors.New("不支持该币种.")

	MSG_9001            = errors.New("审批流模板已存在.")
	MSG_9002            = errors.New("该模板不存在.")
	MSG_9003            = errors.New("模板状态不对")
	MSG_9007            = errors.New("创建合约失败")
	MSG_9008            = errors.New("重复创建合约")
	MSG_9009            = errors.New("审批流模板被篡改")
	MSG_9010            = errors.New("审核流已被其他股东拒绝")
	MSG_Voucher_Timeout = errors.New("签名机获取状态超时")

	MSG_10001 = errors.New("参数错误")
	MSG_10002 = errors.New("查询列表错误")
	MSG_10003 = errors.New("查询错误")
	MSG_10005 = errors.New("查询明细错误")
	MSG_10010 = errors.New("查找模板hash错误")
	MSG_10011 = errors.New("解析模板内容错误")
	MSG_10013 = errors.New("模板不存在或者状态不正确")
	MSG_10014 = errors.New("解析转账内容错误")
	MSG_10015 = errors.New("查找订单错误")
	MSG_10016 = errors.New("该申请已被其他用户审核")
	MSG_10017 = errors.New("查询该层审批人错误")
	MSG_10018 = errors.New("该审批状态不对")
	MSG_10019 = errors.New("该用户不能操作")
	MSG_10020 = errors.New("该审批流已经通过，无需审批")
	MSG_10021 = errors.New("拒绝人数已经超过，结束")
	MSG_10022 = errors.New("操作数据库失败")
	MSG_10023 = errors.New("审批不通过")
	MSG_10026 = errors.New("通过orderid查找转账信息错误")
	MSG_10027 = errors.New("状态不对")
	MSG_10028 = errors.New("该用户不能操作")
	MSG_10029 = errors.New("数据库错误")
	MSG_10032 = errors.New("该订单已过期")
	MSG_10033 = errors.New("查找模板额度失败")
	MSG_10034 = errors.New("当前申请额度大于剩余额度")
	MSG_10035 = errors.New("该笔转账已过期，重新扫码")
)
