package notify

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"io/ioutil"
	"log"
	"micro-services/order-server/pkg/kafka"
	"micro-services/pay-server/internal/repository"
	"micro-services/pay-server/pkg/ali"
	"net/http"
)

// AlipayNotify 处理支付宝异步通知 TODO 测试
func AlipayNotify() {

	data, err := ioutil.ReadFile("pay-server/pkg/ali/privateKey.pem")
	if err != nil {
		log.Fatalf("读取私钥失败: %v", err)
	}
	// 初始化支付宝客户端
	client, err := alipay.New(ali.AppID, string(data), false)
	if err != nil {
		log.Fatalf("初始化支付宝客户端失败: %v", err)
	}
	// 加载支付宝公钥
	err = client.LoadAliPayPublicKey(ali.AlipayPublicKey)
	if err != nil {
		log.Fatalf("加载支付宝公钥失败: %v", err)
	}

	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		// 获取异步通知参数
		err := r.ParseForm()
		if err != nil {
			log.Printf("解析请求参数失败: %v", err)
			return
		}
		noti, err := client.DecodeNotification(r.Form)
		if err != nil {
			log.Printf("解析支付宝通知失败: %v", err)
			http.Error(w, "解析通知失败", http.StatusBadRequest)
			return
		}
		// 验证签名
		err = client.VerifySign(r.Form)
		if err != nil {
			log.Printf("支付宝通知签名验证失败: %v", err)
			http.Error(w, "验证失败", http.StatusBadRequest)
			return
		}
		fmt.Println("noti:", noti)
		// 处理不同的交易状态
		switch noti.TradeStatus {
		case "WAIT_BUYER_PAY":
			log.Printf("交易创建，等待买家付款，订单号: %s\n", noti.OutTradeNo)
		case "TRADE_SUCCESS":
			// 消费掉未支付（初始状态下）Kafka的消息
			err := kafka.ConsumePartition(0, noti.OutTradeNo)
			if err != nil {
				log.Printf("消费消息失败: %v", err)
				return
			}
			// 发送新消息到已支付状态的Kafka
			err = kafka.SendMessageToPartition(1, noti.OutTradeNo)
			if err != nil {
				log.Printf("发送消息到Kafka失败: %v", err)
				return
			}
			// 修改数据库中订单状态
			err = repository.ReverseOrderStatus(noti.OutTradeNo, 1)
			if err != nil {
				log.Println("修改数据库订单状态失败")
				return
			}
			// 删除redis里的订单二维码
			err = repository.DeleteAliPayQRCode(noti.OutTradeNo)
			if err != nil {
				log.Println("删除redis订单二维码失败")
				return
			}
			log.Printf("支付成功，订单号: %s\n", noti.OutTradeNo)
			// 在这里可以进行支付成功的业务处理，例如更新数据库订单状态
		case "TRADE_CLOSED":
			log.Printf("交易关闭，订单号: %s\n", noti.OutTradeNo)
		default:
			log.Printf("未知交易状态: %s，订单号: %s\n", noti.TradeStatus, noti.OutTradeNo)
		}

		// 响应支付宝，告知已经收到通知
		_, err = w.Write([]byte("success"))
		if err != nil {
			log.Printf("响应支付宝失败: %v", err)
			return
		}
	})

	// 启动 HTTP 服务
	go func() {
		if err := http.ListenAndServe(":8888", nil); err != nil {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

}
