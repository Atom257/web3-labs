package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// ⚠️ 这里只需要填你正在运行的服务的地址，不需要连库
const BaseURL = "http://127.0.0.1:8080"

// ⚠️ 填入你想查询的真实参数
const (
	TargetChainID  = "11155111"
	TargetContract = "0xBEfe9d9726c3BFD513b6aDd74B243a82b272C073" // 你的真实合约地址
	TargetAccount  = "0x569744F510D38d7e8E68829f284AAa7F07611552" // 你的真实用户地址
)

func TestLiveServer_OutputOnly(t *testing.T) {
	// 这是一个纯客户端测试，它假设你的服务已经在别处启动了

	callAndPrint := func(name, url string) {
		fullURL := BaseURL + url
		fmt.Printf("\n====== 测试接口: %s ======\n", name)
		fmt.Printf("GET 请求: %s\n", fullURL)

		// 发起真实的 HTTP 网络请求
		resp, err := http.Get(fullURL)
		if err != nil {
			t.Errorf("请求失败 (请确认服务是否已启动?): %v", err)
			return
		}
		defer resp.Body.Close()

		fmt.Printf("状态码: %d\n", resp.StatusCode)

		// 读取并格式化 Body
		bodyBytes, _ := io.ReadAll(resp.Body)
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
			fmt.Printf("返回数据:\n%s\n", prettyJSON.String())
		} else {
			fmt.Printf("返回数据 (原始): %s\n", string(bodyBytes))
		}
	}

	// 1. 测试 /head
	callAndPrint("GetHead", fmt.Sprintf("/head?chain_id=%s&contract=%s", TargetChainID, TargetContract))

	// 2. 测试 /user/points
	callAndPrint("GetUserPoints", fmt.Sprintf("/user/points?chain_id=%s&contract=%s&account=%s", TargetChainID, TargetContract, TargetAccount))

	// 3. 测试 /rate/current
	callAndPrint("GetCurrentRate", fmt.Sprintf("/rate/current?chain_id=%s&contract=%s", TargetChainID, TargetContract))

	// 4. 测试 /user/point_logs
	callAndPrint("GetUserPointLogs", fmt.Sprintf("/user/point_logs?chain_id=%s&contract=%s&account=%s&limit=5", TargetChainID, TargetContract, TargetAccount))
}
