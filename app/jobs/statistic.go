package jobs

import (
	"mango-admin/app/admin/service"
	"time"
)

// 新添加的job 必须按照以下格式定义，并实现Exec函数
type StatisticJob struct {
}

func (t StatisticJob) Exec(arg interface{}) error {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	service.User{}.Statistic(yesterday, true)

	time.Sleep(time.Second * 2)
	for i := 1; i <= 7; i++ {
		service.User{}.StatisticRemainCount(yesterday, i)
	}
	return nil
}
