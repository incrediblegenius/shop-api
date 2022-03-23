package api

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/http"
	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"strings"
	"time"
)

func GenerateSmsCode(width int) string {
	numeric := [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, err := fmt.Fprintf(&sb, "%c", numeric[rand.Intn(r)])
		if err != nil {
			return ""
		}
	}
	return sb.String()

}

func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecret)
	/* use STS Token
	client, err := dysmsapi.NewClientWithStsToken("cn-hangzhou", "<your-access-key-id>", "<your-access-key-secret>", "<your-sts-token>")
	*/
	mobile := sendSmsForm.Mobile
	smsCode := GenerateSmsCode(6)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.Method = "POST"
	request.SignName = "阿里云短信测试"
	request.TemplateCode = "SMS_154950909"
	request.PhoneNumbers = mobile
	request.TemplateParam = "{\"code\":\"" + smsCode + "\"}"

	response, err := client.SendSms(request)
	fmt.Printf("response is %#v\n", response)
	if err != nil {
		fmt.Print(err.Error())
	}
	//保存验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(), mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})

}
