package notify

import (
	"github.com/smartwalle/alipay/v3"
	"io/ioutil"
	"log"
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
		r.ParseForm()
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

		// 处理不同的交易状态
		switch noti.TradeStatus {
		case "WAIT_BUYER_PAY":
			log.Printf("交易创建，等待买家付款，订单号: %s\n", noti.OutTradeNo)
		case "TRADE_SUCCESS":
			log.Printf("支付成功，订单号: %s\n", noti.OutTradeNo)
			// 在这里可以进行支付成功的业务处理，例如更新数据库订单状态
		case "TRADE_CLOSED":
			log.Printf("交易关闭，订单号: %s\n", noti.OutTradeNo)
		default:
			log.Printf("未知交易状态: %s，订单号: %s\n", noti.TradeStatus, noti.OutTradeNo)
		}

		// 响应支付宝，告知已经收到通知
		w.Write([]byte("success"))
	})

	// 启动 HTTP 服务
	go func() {
		if err := http.ListenAndServe(":8888", nil); err != nil {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

}
