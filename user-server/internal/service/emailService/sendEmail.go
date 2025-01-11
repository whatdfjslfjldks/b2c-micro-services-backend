package emailService

import (
	"gopkg.in/gomail.v2"
	logServerProto "micro-services/pkg/proto/log-server"
	"micro-services/pkg/utils"
	userPkg "micro-services/user-server/pkg"
	"micro-services/user-server/pkg/config"
	"micro-services/user-server/pkg/instance"
)

func SendEmail(email string, subject string, content string) {
	isValid := userPkg.IsEmailValid(email)
	if !isValid {
		return
	}
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
	emailMsg.SetHeader("Subject", subject)
	emailBody := content

	emailMsg.SetBody("text/html", emailBody)
	err := dialer.DialAndSend(emailMsg)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/sendEmail",
			Source:      "user-server",
			StatusCode:  "USR-001",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
	}
	return
}
