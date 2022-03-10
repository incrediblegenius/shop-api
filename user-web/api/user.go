package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"shop-api/user-web/global/response"
	"shop-api/user-web/proto"
	"strconv"
	"strings"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的错误码转换为http的
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
					"msg": "用户服务不可用",
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(ctx *gin.Context) {
	//ip := "127.0.0.1"
	//port := 8088
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrcConfig.Host, global.ServerConfig.UserSrcConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error())
	}

	userSrvClient := proto.NewUserClient(userConn)

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, v := range rsp.Data {
		//data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       v.Id,
			Nickname: v.NickName,
			Birthday: response.JsonTime(time.Unix(int64(v.BirthDay), 0)),
			Gender:   v.Gender,
			Mobile:   v.Mobile,
		}
		//data["id"] = v.Id
		//data["name"] = v.NickName
		//data["birthday"] = v.BirthDay
		//data["gender"] = v.Gender
		//data["mobile"] = v.Mobile

		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for filed, err := range fields {
		rsp[filed[strings.Index(filed, ".")+1:]] = err
	}
	return rsp
}

func PassWordLogin(c *gin.Context) {
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBindJSON(&passwordLoginForm); err != nil {
		// 如何返回错误信息
		HandleValidatorError(c, err)
		return
	}
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrcConfig.Host, global.ServerConfig.UserSrcConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[PassWordLogin] 连接 【用户服务失败】",
			"msg", err.Error())
	}
	userSrvClient := proto.NewUserClient(userConn)

	// 登陆逻辑

	if rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登陆失败",
				})
			}
		}
	} else {
		// 查询到用户，并未检查密码
		//zap.S().Debugf("error: %s %s ", passwordLoginForm.PassWord, rsp.PassWord)
		if passRsp, passErr := userSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); passErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登陆失败",
			})
		} else {
			//fmt.Println(passRsp.Success)
			if passRsp.Success {
				c.JSON(http.StatusOK, map[string]string{
					"msg": "登陆成功",
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登陆失败",
				})
			}

		}
	}
}
