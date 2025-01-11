package service

import (
	pb "micro-services/pkg/proto/user-server"
	"micro-services/pkg/utils"
	"micro-services/risk-server/pkg/instance"
)

func SendEmail(email string, subject string, ip string, agent string) {
	content := `<!DOCTYPE html>
<html lang="zh">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>b2c电商平台 - 异常登录提醒</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #fff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        h2 {
            color: #333;
        }
        p {
            color: #555;
        }
        .highlight {
            color: #e74c3c;
            font-weight: bold;
        }
        .button {
            display: inline-block;
            background-color: #3498db;
            color: white;
            padding: 10px 20px;
            border-radius: 4px;
            text-decoration: none;
            text-align: center;
            margin-top: 20px;
        }
        .footer {
            text-align: center;
            font-size: 12px;
            color: #aaa;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>b2c电商平台 - 异常登录提醒</h2>
        <p>尊敬的用户，您好！</p>
        <p>您的账户在以下地点登录：<br>
            <strong>IP 地址:</strong> <span class="highlight">` + ip + `</span><br>
            <strong>设备信息:</strong> <span class="highlight">` + agent + `</span><br>
            <strong>登录时间:</strong> <span class="highlight">` + utils.GetTime() + `</span><br>
        </p>
        <p>此登录地点和设备与您平时使用的不同，可能存在账户被他人恶意访问的风险。</p>
        <p>如果这不是您本人操作，请尽快登录平台修改密码，确保账户安全。</p>
        <a href="[修改密码链接]" class="button">立即修改密码</a>
        <div class="footer">
            <p>如果您有任何疑问，欢迎联系 b2c电商平台客服。</p>
        </div>
    </div>
</body>
</html>
`

	instance.GrpcClient.SendEmail(&pb.SendEmailRequest{
		Email:   email,
		Subject: subject,
		Content: content,
	})

}
