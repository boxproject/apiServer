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
package routers

import (
	"net/http"

	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/controllers"
	middle "github.com/boxproject/apiServer/middleware"
	"github.com/boxproject/apiServer/versionlog"
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	e "github.com/boxproject/apiServer/errors"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/js", "./static/dist/js")
	router.Static("/css", "./static/dist/css")
	router.Static("/img", "./static/dist/img")
	router.StaticFile("/webpage", "./static/dist/index.html")

	// 处理未被捕捉的错误
	router.Use(nice.Recovery(recoveryHandler))
	app := router.Group("/api/")
	app.GET("version", controllers.GetVersion)
	app = router.Group("/api/" + versionlog.LogMap["version"].(string))

	// web transfer group
	webTransfer := app.Group("/webtransfer")
	{ // 获取币种列表
		webTransfer.GET("/getcoins", controllers.GetCoins)
		// 获取模板列表
		webTransfer.GET("/gettemplate", controllers.GetTemplate)
		// 提交转账
		webTransfer.POST("/transfercommit", controllers.TransferCommit)
		// 获取网页模板链接
		webTransfer.GET("/getwebrouter", controllers.GetWebRouter)
		// 获取提交后的转账状态
		webTransfer.GET("/getcommitstatus", controllers.GetCommitStatus)
	}
	// App and the voucher interaction routing group
	voucher := app.Group("/voucher")
	{
		// 查询连接状态
		voucher.GET("/getconnectstatus", controllers.GetConnectStatus)
		// 轮询签名机状态
		voucher.GET("/getstatus", controllers.GetStatus)
		// 扫码提交
		voucher.POST("/qrcommit", controllers.QrCommit)
		// voucher.POST("/qrcommit", middle.VerifyParamSignMiddleWare(), controllers.QrCommit)
		//提交备份口令（全部一样）
		voucher.POST("/commandcommit", controllers.CommandCommit)
		//各个私钥app同步其他公钥
		voucher.GET("/getotherpubkey", controllers.GetOtherPubKey)
		//备份关键句
		voucher.POST("/backupkey", controllers.BackupKey)
		//传对其他公钥的签名值
		//voucher.POST("/savesign", controllers.QrCode)
		//启动签名机
		voucher.POST("/start", controllers.StartVoucher)
		//关闭签名机
		voucher.POST("/stop", controllers.StopVoucher)
		//对其他股东app签名保存
		voucher.POST("/savepubkeysign", controllers.SavePubkeySign)
		//获取token
		voucher.GET("/gettoken", controllers.GetToken)
		//获取各个主地址
		voucher.GET("/getmasteraddress", controllers.GetMasterAddress)
		//保存对主地址的签名
		voucher.POST("/savemasteraddress", controllers.SaveMasterAddress)
		//交换秘钥
		voucher.GET("/operatekeyexchange", controllers.OperateKeyExchange)
	}

	// manage account group
	account := app.Group("/account")
	{
		// 注册
		account.POST("/signup", controllers.Signup)
		// 登录
		account.POST("/login", controllers.Login)
		// 获取用户状态
		account.GET("/finduserstatus", controllers.FindUserStatus)
		// 重新注册
		account.POST("/resignup", controllers.ReSignUp)
		// 恢复股东账户
		account.POST("/recoveryowner", controllers.RecoveryOwner)
		// 重新恢复
		account.POST("/resetrecovery", controllers.ResetRecovery)
		// 获取恢复股东认证结果
		account.GET("/recoveryresult", controllers.RecoveryResult)
		// 激活股东账户
		account.POST("/activerecovery", controllers.ActiveRecovery)
		// 获取以被恢复股东签名的用户公钥
		account.GET("/recpubkeys", controllers.GetPubKeys)
		// 是否有进行中的
		account.GET("/hasrecovery", controllers.HasRecovery)
	}

	// get log group
	logger := app.Group("/logger")
	{
		// 获取日志
		logger.GET("/logs", controllers.GetLogs)
	}
	// 校验token
	app.Use(middle.JWTAuth())
	// 校验参数签名
	app.Use(middle.VerifyParamSignMiddleWare())

	// manage user group
	user := app.Group("/user")
	{
		// 获取账号组织信息
		user.GET("/usertree", controllers.UserTree)
		// 审核
		user.POST("/verifyuser", middle.AdminAuth("admin"), controllers.VerifyUser)
		// 扫码提交认证
		user.POST("/scan", controllers.Scan)
		// 修改密码
		user.POST("/modifypassword", controllers.ModifyPassword)
		// 获取指定角色用户列表
		user.GET("/listbytype", controllers.GetAccountsByType)
		// 添加管理员
		user.POST("/addadmin", middle.AdminAuth("owner"), controllers.AddAdmin)
		// 删除管理员
		user.POST("/deladmin", middle.AdminAuth("owner"), controllers.DelAdmin)
		// 获取所有用户列表
		user.GET("/list", middle.AdminAuth("admin"), controllers.GetAllUsers)
		// 获取单个用户信息
		user.GET("/userbyid", middle.AdminAuth("admin"), controllers.GetUserByID)
		// 设置用户信息
		user.POST("/setuser", middle.AdminAuth("admin"), controllers.SetUser)
		// 停用账户
		user.POST("/disableacc", middle.AdminAuth("admin"), controllers.DisableAcc)
		// 获取注册扫码后的用户列表
		user.GET("/reglist", middle.AdminAuth("admin"), controllers.GetRegList)
		// 获取待恢复用户列表
		user.GET("/recoverylist", middle.AdminAuth("owner"), controllers.RecoveryList)
		// 验证密码
		user.POST("/verifypwd", controllers.VerifyPwd)
		// 扫码提交恢复股东信息
		user.POST("/subrecovery", middle.AdminAuth("owner"), controllers.SubRecovery)
		// 认证股东
		user.POST("/verifyrecovery", middle.AdminAuth("owner"), controllers.VerifyRecovery)
		// 待办任务统计
		user.GET("/delaytasknum", controllers.DelayTaskNum)
		//消息列表
		user.GET("insideletter", controllers.InsideLetter)
		//读消息
		user.POST("readletter", controllers.ReadLetter)
		//获取块高
		user.GET("/blockheight", controllers.GetBlockHeight)
	}

	// department group
	department := app.Group("/department")
	{
		// 添加部门
		department.POST("/add", middle.AdminAuth("admin"), controllers.AddDepartment)
		// 获取部门列表
		department.GET("/list", controllers.DepartmentList)
		// 编辑部门
		department.POST("/edit", middle.AdminAuth("admin"), controllers.EditDepartment)
		// 删除部门
		department.POST("/delete", middle.AdminAuth("admin"), controllers.DelDepartment)
		// 获取指定权限指定部门下成员列表
		department.GET("/accountsbydepid", controllers.DepartmentAccounts)
		// 部门排序
		department.POST("/sortdep", middle.AdminAuth("admin"), controllers.SortDep)
	}

	// auth group
	auth := app.Group("/auth", middle.AdminAuth("admin"))
	{
		// 获取权限列表
		auth.GET("/list", controllers.AuthList)
		// 获取指定权限下用户列表
		auth.GET("/accountsbyauthid", controllers.AuthAccounts)
		// 用户添加权限
		auth.POST("/addauthtoaccount", middle.AdminAuth("admin"), controllers.AddAuthToAccount)
		// 用户删除权限
		auth.POST("/delauthfromaccount", middle.AdminAuth("admin"), controllers.DelAuthFromAccount)
	}

	// template group
	template := app.Group("/template")
	{
		// 新建审批流模板
		template.POST("/new", middle.AdminAuth("admin"), controllers.CreateTemplate)
		// 获取审批流模板列表
		template.GET("/list", controllers.TemplateList)
		// 获取审批流模板详情
		template.GET("/findtembyid", controllers.FindTemplateById)
		// 审批审批流模板
		template.POST("/verify", middle.AdminAuth("owner"), middle.VerifyKeyWord(), controllers.VerifyTemplate)
		// 作废审批流模板
		template.POST("/cancel", middle.AdminAuth("admin"), controllers.CancelTemplate)
		// 统计审批流模板
		template.GET("/tempstatistic", middle.AdminAuth("admin"), controllers.TempStatistics)
		// 根据模板统计转账数目
		template.GET("/transfers", middle.AdminAuth("admin"), controllers.TxNumByTemplate)
	}

	// coin group
	coin := app.Group("/coin")
	{
		// 新增币种
		coin.POST("/add", middle.AdminAuth("admin"), controllers.AddCoin)
		// 校验地址
		coin.POST("/verifyaddress", controllers.VerifyAddress)
		// 获取币种列表
		coin.GET("/list", controllers.CoinList)
		// 获取指定币种状态
		coin.POST("/status", controllers.CoinStauts)
		// 根据币种获取余额
		coin.GET("/balance", middle.AuthCheck("balance"), controllers.CoinBalance)
		// 付款码
		coin.GET("/qrcode", controllers.QRcode)
	}

	// account manage group
	accountmanage := app.Group("/accmanage") //多账户管理
	{
		//主链币多账户统计
		accountmanage.GET("/statistical", controllers.Statistical)
		//子账户统计
		accountmanage.GET("/childstical", controllers.ChildStatistical)
		//计算子账户数量和总额
		accountmanage.GET("/countallchild", controllers.CountAllChild)
		//主链币对应账户分页列表
		accountmanage.GET("/listbypage", controllers.AddressList)
		//账户明细
		accountmanage.GET("/detail", controllers.AccountDetail)
		//转账记录
		accountmanage.GET("/transferrecord", controllers.TransferRecord)
		//账户金额top10
		accountmanage.GET("/topten", controllers.AddressTopTen)
		//代币转账记录
		accountmanage.GET("/tokenrecord", controllers.TokenRecord)
		//代币明细
		accountmanage.GET("/tokenlist", controllers.TokenList)
		//设置别名
		accountmanage.POST("/settag", controllers.SetTag)
		//新增子账户
		accountmanage.POST("/addchildaccount", controllers.AddChildAccount)
		//合并子账户
		accountmanage.POST("/merge", controllers.VerifyPassword(), controllers.MergeAccount)
		//创建合约账户
		accountmanage.POST("/add/contract", middle.AdminAuth("admin"), controllers.GenContractAddr)
		//获取合约账户
		accountmanage.GET("/contract", controllers.ContractAddr)
	}

	transfer := app.Group("/transfer") //转账流
	{
		//获取币种列表
		transfer.GET("/getcoins", controllers.GetCoins)
		//获取模板列表
		transfer.GET("/gettemplate", controllers.GetTemplate)
		//申请转账
		transfer.POST("/apply", controllers.VerifyPassword(), controllers.Apply)
		//申请列表
		transfer.GET("/findapplylist", controllers.FindApplyList)
		//通过id查询转账申请详情
		transfer.GET("/findapplybyid", controllers.FindApplyById)
		//通过id查询转账日志
		transfer.GET("/findapplylog", controllers.FindApplyLog)
		//审批
		transfer.POST("/verify", controllers.VerifyPassword(), controllers.Verify)
		//取消申请
		transfer.POST("/cancel", controllers.VerifyPassword(), controllers.Cancel)
		//待审批列表
		transfer.GET("/batchlist", controllers.BatchList)
		//批量审批
		transfer.POST("/batchverify", controllers.VerifyPassword(), controllers.BatchVerify)
		//通过模板hash查找模板
		transfer.GET("/gettemplatebyhash", controllers.GetTemplatebyHash)
		//获取转账信息
		transfer.GET("/findtranfersbyid", controllers.FindTranfersById)
	}
	return router
}

// recoveryHandler catch the unknown error to log and then
// throw out the 'system error' error type.
func recoveryHandler(ctx *gin.Context, err interface{}) {
	log.Error("Got Panic", err)
	ctx.JSON(http.StatusInternalServerError, gin.H{"Code": e.Failed, "Msg": e.SYSTEM_ERROR})
	return
}
