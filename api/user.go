package api

import (
	"context"
	"fmt"
	"mxshop_web/forms"
	"mxshop_web/global"
	"mxshop_web/global/response"
	"mxshop_web/proto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用" + e.Message(),
				})
			case codes.Unknown:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "未知错误",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
		}
	}
}

func HandleValidatorError(ctx *gin.Context, err error) {
	// TODO 英文提示转成中文
	ctx.JSON(http.StatusBadRequest, gin.H{
		"msg": err.Error(),
	})
}

func GetUserList(ctx *gin.Context) {
	zap.S().Debugf("获取用户列表")
	zap.S().Debugf(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port))
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	//调用接口
	userSrvClient := proto.NewUserClient(userConn)
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn) // 字符串到int的强转
	pSize := ctx.DefaultQuery("psize", "0")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【获取用户列表失败】")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2000-1-11"),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		result = append(result, user)

	}
	ctx.JSON(http.StatusOK, result)
}

func PasswordLogin(ctx *gin.Context) {
	//表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		// TODO 英文提示转成中文
		HandleValidatorError(ctx, err)
		return
	}
	//获取用户验证密码
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	//调用接口
	userSrvClient := proto.NewUserClient(userConn)
	if rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		//验证密码
		if passRsp, passErr := userSrvClient.CheckPassWord(context.Background(), &proto.PassWordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); passErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"password": "登录失败",
			})
		} else {
			if passRsp.Success {
				ctx.JSON(http.StatusOK, gin.H{
					"msg": "登录成功",
				})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "密码错误",
				})
			}
		}
	}
}
