package test

import (
	"container/list"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 时间轮层数
const timeWheelsCount = 5

// 每层时间轮槽的数量
const timeWheelsSize = 60

// 时间轮层
type TimeWheel struct {
	// 当前节点指针
	currentSlot int
	// 槽数组
	slots []*list.List
	// 下一层时间轮
	nextWheel *TimeWheel
	// 每个槽代表的单位时间
	slotSpan time.Duration
	// 缓存间隔 time.NewTicker 频率
	tickInterval time.Duration
	// 最近一次调用 time.Now 的结果
	lastCalledTime time.Time
}

func NewTimeWheel(span time.Duration, wheels int) *TimeWheel {
	// 创建时间轮实例
	tw := &TimeWheel{
		slots:        make([]*list.List, timeWheelsSize),
		slotSpan:     span,
		tickInterval: span / time.Duration(timeWheelsSize),
	}
	// 初始化第一层
	current := tw
	for i := 1; i < wheels; i++ {
		next := &TimeWheel{
			slots:        make([]*list.List, timeWheelsSize),
			slotSpan:     current.slotSpan * time.Duration(timeWheelsSize),
			tickInterval: current.slotSpan * time.Duration(timeWheelsSize),
		}
		current.nextWheel = next
		current = next
	}
	return tw
}

// 根据超时时间戳计算通知应该被添加到哪一层时间轮和该层的哪个槽位
func (tw *TimeWheel) getWheelAndSlotIndex(expiration int64) []*list.List {
	wheels := make([]*list.List, timeWheelsCount)
	currentTime := tw.lastCalledTime.UnixMilli()
	if currentTime >= expiration {
		wheels[0] = tw.slots[0]
		return wheels
	}

	offsetMilli := expiration - currentTime
	wheelIndex := 0
	slotIndex := 0
	for offsetMilli >= int64(tw.tickInterval) {
		if wheelIndex+1 < timeWheelsCount {
			offsetMilli = offsetMilli / int64(tw.nextWheel.slotSpan/tw.slotSpan)
			wheelIndex++
			tw = tw.nextWheel
		} else {
			break
		}
	}

	slotIndex = int(offsetMilli / int64(tw.tickInterval))
	wheels[wheelIndex] = tw.slots[slotIndex]
	for i := wheelIndex - 1; i >= 0; i-- {
		wheels[i] = tw.nextWheel.slots[0]
		tw = tw.nextWheel
	}

	return wheels
}

// 添加定时任务到时间轮
func (tw *TimeWheel) addNotification(n *Notification) {
	// 计算超时的时间戳
	expiration := n.NotifyAt.UnixMilli()
	// 根据超时时间计算在哪一层哪个槽
	wheels := tw.getWheelAndSlotIndex(expiration)
	wheels[0].PushBack(n)
	// 向高层时间轮添加节点
	for i := 1; i < len(wheels); i++ {
		if wheels[i] != nil {
			wheels[i].PushBack(n)
		}
	}
}

// 获取定时通知
func (tw *TimeWheel) getNotifications() []*Notification {
	var notifications []*Notification
	for e := tw.slots[tw.currentSlot].Front(); e != nil; e = e.Next() {
		n := e.Value.(*Notification)
		notifications = append(notifications, n)
	}
	return notifications
}

// 时间轮前进逻辑
func (tw *TimeWheel) moveForward() {
	// 当前槽的任务复制到下一层时间轮
	tasks := tw.slots[tw.currentSlot]
	tw.slots[tw.currentSlot] = nil
	if tasks != nil {
		for e := tasks.Front(); e != nil; e = e.Next() {
			n := e.Value.(*Notification)
			tw.nextWheel.addNotification(n)
		}
	}
	// 前进当前指针
	tw.currentSlot = (tw.currentSlot + 1) % timeWheelsSize
	if tw.nextWheel != nil {
		tw.nextWheel.moveForward()
	}
}

// 表示一条通知记录
type Notification struct {
	ID       int
	Content  string
	NotifyAt time.Time
}

// 多层时间轮实例
var tw *TimeWheel

func init() {
	// 创建5层时间轮,最小精度为1毫秒
	tw = NewTimeWheel(1*time.Millisecond, 5)
}

func BenchmarkTimeWheel(b *testing.B) {
	// 生成随机的通知记录
	notifications := generateRandomNotifications(b.N)

	// 重置计时器
	b.ResetTimer()

	// 启动一个goroutine来驱动时间轮前进
	go func() {
		ticker := time.NewTicker(1 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			tw.moveForward()
			sendNotifications(tw.getNotifications())
		}
	}()

	// 将通知插入时间轮
	for _, n := range notifications {
		tw.addNotification(n)
	}

	// 让基准测试运行一段时间
	time.Sleep(5 * time.Second)
}

func generateRandomNotifications(n int) []*Notification {
	var notifications []*Notification
	for i := 0; i < n; i++ {
		delay := rand.Intn(5000) // 模拟0~5秒的随机延迟
		content := fmt.Sprintf("This is notification %d", i)
		notifications = append(notifications, &Notification{
			ID:       i,
			Content:  content,
			NotifyAt: time.Now().Add(time.Duration(delay) * time.Millisecond),
		})
	}
	return notifications
}

func sendNotifications(notifications []*Notification) {
	for _, n := range notifications {
		fmt.Printf("Notification %d: %s (expired at %v)\n", n.ID, n.Content, n.NotifyAt)
	}
}
