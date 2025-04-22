package main

// func main() {
// 	provider := services.NewSparkProvider(
// 		"0b45810b",
// 		"ZTlmZTNiMDkwYzVkMzk4YTM0ZWRmMTI1",
// 		"3c0127462241e39ed723627e4509ba53",
// 		"spark-api.xf-yun.com",
// 		"wss://spark-api.xf-yun.com/v1/x1",
// 		"spark-x1",
// 	)

// 	history := []models.Message{
// 		{Role: "user", Content: "你好"},
// 		{Role: "assistant", Content: "你好呀！今天过得怎么样？"},
// 	}

// 	reply, err := provider.Chat("我今天有点难过", history)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("✨ AI 回复：", reply)
// }