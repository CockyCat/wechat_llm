package martin

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
	configQuant "wechat_llm/quant/config"

	"github.com/iaping/go-okx/rest"
	"github.com/iaping/go-okx/rest/api"
	"github.com/iaping/go-okx/rest/api/market"
	"github.com/iaping/go-okx/rest/api/trade"
)

const (
	symbol          = "BTC-USDT" // 交易对
	initialAmount   = 10.0       // 初始交易金额
	maxRetries      = 5          // 最大加倍次数
	priceDifference = 100        // 每档价格差值
)

type priceLevel struct {
	lowerBound float64
	upperBound float64
}

func Martin() {
	client := NewOkxClient(OkxApiConf{
		ApiHost:     configQuant.LoadOKXConfig().ApiHost,
		ApiKey:      configQuant.LoadOKXConfig().ApiKey,
		SecretKey:   configQuant.LoadOKXConfig().SecretKey,
		Passphrase:  configQuant.LoadOKXConfig().Passphrase,
		IsSimulated: configQuant.LoadOKXConfig().IsSimulated,
	})

	ctx := context.Background()
	baseAmount := initialAmount

	// 获取初始市场价格
	initialPrice, err := getMarketPrice(ctx, client, symbol)
	if err != nil {
		log.Fatalf("获取市场价格失败: %v", err)
	}

	for level := 1; level <= maxRetries; level++ {
		// 计算当前档位的价格
		targetPrice := initialPrice - float64((level-1)*priceDifference)

		// 获取当前市场价格
		currentPrice, err := getMarketPrice(ctx, client, symbol)
		if err != nil {
			log.Fatalf("获取市场价格失败: %v", err)
		}

		// 检查当前价格是否在本档位的价格范围内
		if currentPrice > targetPrice {
			fmt.Printf("当前价格超过第 %d 档位的价格范围，等待下一次检查\n", level)
			time.Sleep(1 * time.Minute) // 暂停1分钟后再尝试
			level--                     // 重新检查当前档位
			continue
		}

		// 计算本档位的交易量
		adjustedAmount := baseAmount * float64(level)

		// 挂限价单
		orderID, err := placeLimitOrder(ctx, client, symbol, targetPrice, adjustedAmount)
		if err != nil {
			log.Fatalf("挂单失败: %v", err)
		}

		// 检查订单状态
		success, err := checkOrderStatus(ctx, client, symbol, orderID)
		if err != nil {
			log.Fatalf("检查订单状态失败: %v", err)
		}

		if success {
			fmt.Printf("第 %d 档位交易成功!\n", level)
		} else {
			fmt.Printf("第 %d 档位交易失败，调整下注金额\n", level)
			baseAmount *= 2 // 加倍下注金额
			level--         // 重新尝试当前档位
		}

		time.Sleep(1 * time.Minute) // 暂停1分钟后再尝试下一档
	}

	fmt.Println("所有档位交易完成!")
}

func getMarketPrice(ctx context.Context, client *rest.Client, symbol string) (float64, error) {
	param := &market.GetTickersParam{
		InstType: api.InstTypeSPOT,
	}
	req, resp := market.NewGetTickers(param)
	if err := client.Do(req, resp); err != nil {
		panic(err)
	}
	log.Println(req, resp.(*market.GetTickersResponse))

	for _, ticker := range resp.(*market.GetTickersResponse).Data {
		if ticker.InstId == symbol {
			return strconv.ParseFloat(ticker.Last, 64)
		}
	}
	return 0, fmt.Errorf("未找到交易对: %s 的市场价格", symbol)
}

func placeOrder(ctx context.Context, client *rest.Client, symbol string, amount float64) (string, error) {
	param := &trade.PostOrderParam{
		InstId:  symbol,
		TdMode:  api.TdModeCash,
		Side:    api.SideBuy,
		OrdType: api.OrdTypeLimit,
		Sz:      "9",
		Px:      "5",
	}
	req, resp := trade.NewPostOrder(param)
	if err := client.Do(req, resp); err != nil {
		panic(err)
	}
	log.Println(req, resp.(*trade.PostOrderResponse))
	return resp.(*trade.PostOrderResponse).Data[0].OrdId, nil
}

func placeLimitOrder(ctx context.Context, client *rest.Client, symbol string, price, amount float64) (string, error) {
	param := &trade.PostOrderParam{
		InstId:  symbol,
		TdMode:  api.TdModeCash,
		Side:    api.SideBuy,
		OrdType: api.OrdTypeLimit,
		Sz:      strconv.FormatFloat(amount, 'f', -1, 64),
		Px:      strconv.FormatFloat(price, 'f', -1, 64),
	}
	req, resp := trade.NewPostOrder(param)
	if err := client.Do(req, resp); err != nil {
		return "", err
	}
	return resp.(*trade.PostOrderResponse).Data[0].OrdId, nil
}

func checkOrderStatus(ctx context.Context, client *rest.Client, symbol string, orderID string) (bool, error) {
	param := &trade.GetOrderParam{
		InstId: symbol,
		OrdId:  orderID,
	}
	req, resp := trade.NewGetOrder(param)
	if err := client.Do(req, resp); err != nil {
		return false, nil
	}
	isFilled := resp.(*trade.GetOrderResponse).Data[0].IsFilled()
	return isFilled, nil
}
