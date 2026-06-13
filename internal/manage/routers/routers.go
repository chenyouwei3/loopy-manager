package routers

import (
	"gin-web/init/runLog"
	"gin-web/internal/manage/controller"
	"gin-web/internal/manage/middleware"
	"gin-web/pkg/extendController"
	publicMiddleware "gin-web/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	//store := cookie.NewStore([]byte("something-very-secret"))
	//r.Use(sessions.Sessions("mysession", store))
	r.Use(
		middleware.OperationLog(runLog.RunningLog),
		publicMiddleware.CorsMiddleware(),
	)
	//r := gin.New()
	//r.Use(gin.Logger(), gin.Recovery()) //动记录所有 HTTP 请求的详细信息，如请求方法、请求路径、状态码、响应时间等。
	userCH := &controller.UserHandlerController{
		extendController.BaseController{RunLog: runLog.RunningLog},
	}
	roleCH := &controller.RoleHandlerController{
		extendController.BaseController{RunLog: runLog.RunningLog},
	}
	logCH := &controller.LogHandlerController{
		extendController.BaseController{RunLog: runLog.RunningLog},
	}
	r.POST("/sign_in", userCH.Login())
	r.POST("/sign_up", userCH.Insert())
	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, "success")
	})
	//r.Use(publicMiddleware.AuthMiddleware())
	user := r.Group("user")
	{
		user.POST("/delete", userCH.Delete())                    //删除用户
		user.POST("/update", userCH.Update())                    //更新用户
		user.GET("/getList", userCH.GetList())                   //查询用户列表
		user.GET("/getRolesByUserId", userCH.GetRolesByUserID()) //根据用户id查询角色
	}
	role := r.Group("role")
	{
		role.POST("/insert", roleCH.Insert())  //添加角色
		role.POST("/delete", roleCH.Delete())  //删除角色
		role.POST("/update", roleCH.Update())  //更新角色
		role.GET("/getList", roleCH.GetList()) //查询角色列表
	}
	log := r.Group("log")
	{
		log.GET("/operation/getList", logCH.GetListByOperation()) //查询操作日志列表
		log.GET("/running/getList", logCH.GetListByRun())         //查询运行日志列表
		log.GET("/running/getOne", logCH.GetOneByRun())           //查询单个运行日志
	}
	//cfg := r.Group("config")
	{

	}
	return r
}
