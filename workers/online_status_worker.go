package workers

import (
	"fmt"
	"time"

	"gin-doniai/handlers"
)
// 消息队列
// OnlineStatusUpdate 在线状态更新消息结构
type OnlineStatusUpdate struct {
	UserID    uint
	IP        string
	UserAgent string
}

// HandleOnlineStatusUpdates 处理在线状态更新的消息处理器
func HandleOnlineStatusUpdates(onlineStatusChan <-chan OnlineStatusUpdate) {
	// 批量处理间隔
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var updates []OnlineStatusUpdate
	batchSize := 50 // 批量处理大小

	for {
		select {
		case update := <-onlineStatusChan:
			updates = append(updates, update)

			// 达到批次大小时立即处理
			if len(updates) >= batchSize {
				processBatchOnlineStatus(updates)
				updates = updates[:0] // 清空切片
			}

		case <-ticker.C:
			// 定时处理剩余的更新
			if len(updates) > 0 {
				processBatchOnlineStatus(updates)
				updates = updates[:0]
			}
		}
	}
}

// processBatchOnlineStatus 批量处理在线状态更新
func processBatchOnlineStatus(updates []OnlineStatusUpdate) {
	// 这里调用 handlers 中的方法批量处理
	// 可以根据需要修改 handlers.UpdateUserOnlineStatus 来支持批量操作
	for _, update := range updates {
		// 创建模拟的 gin.Context 用于处理
		// 或者创建新的批量处理方法
		handlers.UpdateUserOnlineStatusWithInfo(update.UserID, update.IP, update.UserAgent)
	}
	fmt.Printf("批量处理了 %d 个在线状态更新\n", len(updates))
}
