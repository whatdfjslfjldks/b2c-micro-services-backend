package emailService

import "C"
import (
	"gopkg.in/gomail.v2"
	logServerProto "micro-services/pkg/proto/log-server"
	"micro-services/pkg/utils"
	"micro-services/user-server/internal/repository"
	"micro-services/user-server/pkg"
	"micro-services/user-server/pkg/config"
	"micro-services/user-server/pkg/instance"
)

func SendEmailCode(email string) (
	msg string,
	err error,
	httpCode int32,
	statusCode string) {
	isValid := userPkg.IsEmailValid(email)
	if !isValid {
		msg = "非法输入！邮箱格式不正确！"
		return msg, err, 400, "GLB-001"
	}
	code := userPkg.GenerateVerifyCode(6)

	dialer := gomail.NewDialer(
		config.EmailSender.Email.Host,     //SMTP服务器地址
		config.EmailSender.Email.Port,     //端口
		config.EmailSender.Email.Sender,   //发件人邮箱
		config.EmailSender.Email.Password, //发件人授权码
	)
	// 创建一个邮件消息
	emailMsg := gomail.NewMessage()
	emailMsg.SetHeader("From", config.EmailSender.Email.Sender)
	emailMsg.SetHeader("To", email)
	emailMsg.SetHeader("Subject", "b2c电商平台-登录验证码")
	emailBody := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>b2c电商平台-登录验证码</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f5f5f5;
					color: #333;
					margin: 0;
					padding: 0;
				}
				.container {
					width: 100%;
					max-width: 600px;
					margin: 30px auto;
					background-color: #fff;
					border-radius: 8px;
					box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
					padding: 20px;
					text-align: center;
				}
				h2 {
					color: #1a73e8;
				}
				.verification-code {
					font-size: 36px;
					font-weight: bold;
					color: #fff;
					background-color: #1a73e8;
					padding: 10px 20px;
					border-radius: 4px;
					display: inline-block;
					margin: 20px 0;
				}
				.footer {
					font-size: 14px;
					color: #888;
					margin-top: 30px;
				}
				.button {
					background-color: #1a73e8;
					color: white;
					border: none;
					border-radius: 4px;
					padding: 10px 20px;
					cursor: pointer;
					text-decoration: none;
					display: inline-block;
					font-size: 16px;
					margin-top: 20px;
				}
				.button:hover {
					background-color: #0c5b9c;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h2>b2c电商平台-登录验证码</h2>
				<p>亲爱的用户，</p>
				<p>请使用以下验证码完成您的登录：</p>
				<div class="verification-code">` + code + `</div>
				<p>该验证码将在 2 分钟内有效。</p>
				<div class="footer">
					<p>如果您没有请求此操作，请忽略此邮件。</p>
					<p>b2c电商平台</p>
				</div>
			</div>
		</body>
		</html>
	`

	emailMsg.SetBody("text/html", emailBody)
	repository.StoreCodeInRedis(email, code)
	err = dialer.DialAndSend(emailMsg)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/sendVerifyCode",
			Source:      "user-server",
			StatusCode:  "USR-001",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
		return "验证码发送失败", err, 500, "USR-001"
	} else {
		return "验证码已发送", nil, 200, "GLB-000"
	}
}
