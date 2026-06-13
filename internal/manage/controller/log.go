package controller

import (
	conf "gin-web/init/config"
	"gin-web/internal/manage/models"
	"gin-web/internal/manage/types"
	"gin-web/pkg"
	"gin-web/pkg/extendController"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type LogHandlerController struct {
	extendController.BaseController
}

func (l *LogHandlerController) GetListByOperation() gin.HandlerFunc {
	return func(c *gin.Context) {
		//接收参数
		logReq := &types.LogByOperationGetListReq{
			Account:   c.Query("account"),
			CurrPage:  c.DefaultQuery("currPage", "1"),
			PageSize:  c.DefaultQuery("pageSize", "10"),
			StartTime: c.Query("startTime"),
			EndTime:   c.Query("endTime"),
		}
		skip, limit, err := pkg.GetPage(logReq.CurrPage, logReq.PageSize)
		if err != nil {
			l.SendServerErrorResponse(c, 5131, err)
			return
		}
		//DB操作
		logDB := models.OperationLog{
			Account: logReq.Account,
		}
		resDB, total, err := logDB.GetList(skip, limit, logReq.StartTime, logReq.EndTime)
		if err != nil {
			l.SendServerErrorResponse(c, 5130, err)
			return
		}
		l.SendSuccessResponse(c, types.LogByOperationGetListResp{
			Logs:  resDB,
			Total: total,
		})
	}
}

func (l *LogHandlerController) GetListByRun() gin.HandlerFunc {
	return func(c *gin.Context) {
		//接收参数
		logReq := &types.LogByRunningGetListReq{
			CurrPage:  c.DefaultQuery("currPage", "1"),
			PageSize:  c.DefaultQuery("pageSize", "10"),
			StartTime: c.Query("startTime"),
			EndTime:   c.Query("endTime"),
		}
		skip, limit, err := pkg.GetPage(logReq.CurrPage, logReq.PageSize)
		if err != nil {
			l.SendServerErrorResponse(c, 5131, err)
			return
		}
		resDB, total, err := new(models.RunningLog).GetList(
			conf.SystemConfig.APP.RunLog,
			skip,
			limit,
			logReq.StartTime,
			logReq.EndTime)
		if err != nil {
			l.SendServerErrorResponse(c, 5130, err)
			return
		}
		l.SendSuccessResponse(c, types.LogByRunningGetListResp{
			Logs:  resDB,
			Total: total,
		})
	}
}

func (l *LogHandlerController) GetOneByRun() gin.HandlerFunc {
	return func(c *gin.Context) {
		//接收参数
		logReq := c.Query("path")
		if logReq == "" {
			l.SendParameterErrorResponse(c, 4001, nil)
			return
		}
		logReq = strings.ReplaceAll(logReq, "\\\\", "\\") // 把 \\ 变成 \
		// 判断文件是否存在
		if _, err := os.Stat(logReq); os.IsNotExist(err) {
			l.SendServerErrorResponse(c, 5104, err)
			return
		}
		// 读取文件
		content, err := ioutil.ReadFile(logReq)
		if err != nil {
			l.SendServerErrorResponse(c, 5130, err)
			return
		}
		//直接返回数据
		l.SendSuccessResponse(c, string(content))
	}
}
