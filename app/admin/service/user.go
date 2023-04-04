package service

import (
	"encoding/json"
	"errors"
	"fmt"
	log "mango-admin/pkg/logger"
	"mango-admin/app/admin/models"
	"mango-admin/app/admin/service/dto"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct{}

func (User) List(req dto.UserListReq, resp *dto.UserListResp) error {
	var ua *models.UserAccount
	offset := (req.Pagination.GetPageIndex() - 1) * req.Pagination.GetPageSize()
	ua.DB().Count(&resp.Count).Offset(offset).Limit(req.Pagination.GetPageSize()).Find(&resp.List)

	// 在线状态
	today := time.Now().Format("2006-01-02")
	var uo *models.UserOnline
	for i := range resp.List {
		var online models.UserOnline
		uo.DB().Where("user_id = ?", resp.List[i].UserId, today).Last(&online)
		if online.SysId > 0 {
			resp.List[i].IsOnline = online.OfflineTime == 0
			resp.List[i].ChangeTime = online.ChangeTime
		}
	}

	return nil
}

func (u User) Dashboard(resp *dto.DashboardResp) error {
	resp.Cards = make([]dto.DashboardRespCard, 4)

	var (
		now     = time.Now()
		weekday = [...]string{"日", "一", "二", "三", "四", "五", "六"}[now.Weekday()]
		// 今天
		today    = now.Format("2006-01-02")
		respDate = today + "(" + weekday + ")"
		// 昨天
		yesterday = now.AddDate(0, 0, -1).Format("2006-01-02")
		// 上周
		lastMonday = now.AddDate(0, 0, -1*int(now.Weekday())-6).Format("2006-01-02")
		lastSunday = now.AddDate(0, 0, -1*int(now.Weekday())).Format("2006-01-02")
		// 统计
		todayStatistic     = u.Statistic(today, false)
		yesterdayStatistic = u.Statistic(yesterday, false)
		lastweekStatistic  = u.rangeStatistic(lastMonday, lastSunday)
	)

	resp.Cards[0].Date = respDate
	resp.Cards[0].Count = todayStatistic.LoginCount
	resp.Cards[0].Compare1 = lastweekStatistic.LoginCount / 7
	resp.Cards[0].Compare2 = yesterdayStatistic.LoginCount

	resp.Cards[1].Date = respDate
	resp.Cards[1].Count = todayStatistic.SignupCount
	resp.Cards[1].Compare1 = lastweekStatistic.SignupCount / 7
	resp.Cards[1].Compare2 = yesterdayStatistic.SignupCount

	resp.Cards[2].Date = respDate[0:7]

	resp.Cards[3].Date = respDate
	resp.Cards[3].Count = todayStatistic.RechargeCount

	// 氧气、叶子🍃
	resp.Oxygen.Total = todayStatistic.OxygenCount
	resp.Oxygen.Consume = todayStatistic.OxygenGenerateCount // 用户的获取是后台意义的消耗
	resp.Leaf.Total = todayStatistic.LeafCount
	resp.Leaf.Consume = todayStatistic.LeafGenerateCount
	resp.Leaf.Recovery = todayStatistic.LeafConsumeCount

	// 表格数据
	resp.Table, _ = u.DashboardTable(dto.DashboardTableReq{LastDays: 7})

	return nil
}

// 首页表格
func (User) DashboardTable(req dto.DashboardTableReq) (resp dto.DashboardTableResp, err error) {
	start := time.Now().AddDate(0, 0, -req.LastDays).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	models.StatisticPtr.DB().Where("the_date >= ? and the_date <= ?", start, today).
		Order("the_date DESC").Find(&resp)
	return
}

// 首页用户趋势
func (User) DashboardOnline(req dto.DashboardOnlineReq) (resp dto.DashboardOnlineResp, err error) {
	if req.Start == "" {
		req.Start = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	}
	if req.End == "" {
		req.End = time.Now().Format("2006-01-02")
	}

	// 每个用户每天总计的在线时长
	var rows []struct {
		OnlineDate string
		Sum        int64
	}
	models.UserOnlinePtr.DB().Where("create_time between ? and ? ", req.Start+" 00:00:00", req.End+" 23:59:59").
		Group("user_id").Group("online_date").Select("user_id, SUBSTRING(create_time, 1, 10) AS online_date, SUM(online_second) as sum").Find(&rows)

	var total int64
	counter := make([]int, 4)
	resp.List = append(resp.List, []interface{}{"日期", "0-2分钟", "2-5分钟", "5-10分钟", "10+分钟"})
	resp.YAxios = make([][]int, 4)
	for i, v := range rows {
		total += v.Sum
		if v.Sum < 120 {
			counter[0]++
		} else if v.Sum < 300 {
			counter[1]++
		} else if v.Sum < 600 {
			counter[2]++
		} else {
			counter[3]++
		}

		if i == len(rows)-1 || v.OnlineDate != rows[i+1].OnlineDate {
			resp.XAxios = append(resp.XAxios, v.OnlineDate)
			resp.YAxios[0] = append(resp.YAxios[0], counter[0])
			resp.YAxios[1] = append(resp.YAxios[1], counter[1])
			resp.YAxios[2] = append(resp.YAxios[2], counter[2])
			resp.YAxios[3] = append(resp.YAxios[3], counter[3])
			resp.List = append(resp.List, []interface{}{v.OnlineDate, counter[0], counter[1], counter[2], counter[3]})
			counter = make([]int, 4)
		}
	}
	if len(rows) > 0 {
		resp.Avg = total / int64(len(rows))
	}
	return
}

// 整体数据
func (User) WholeData(req dto.WholeDataReq) (resp dto.WholeDataResp, err error) {
	if req.Start == "" {
		req.Start = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	}
	if req.End == "" {
		req.End = time.Now().Format("2006-01-02")
	}
	models.StatisticPtr.DB().Where("the_date between ? and ?", req.Start, req.End).
		Order("the_date ASC").Find(&resp.Remain)

	var statistics []struct {
		TheDate     string
		RemainCount models.IntArray
	}
	models.StatisticPtr.DB().Where("the_date between ? and ?", req.Start, req.End).
		Order("the_date ASC").Find(&statistics)

	for i, v := range statistics {
		if i == 0 {
			continue
		}
		var p float32
		// 昨日注册数
		if len(statistics[i-1].RemainCount) > 0 && statistics[i-1].RemainCount[0] > 0 {
			p = float32(v.RemainCount[0]*10000/statistics[i-1].RemainCount[0]) / 100
		}
		resp.Loss.Avg += p
		resp.Loss.List = append(resp.Loss.List, []interface{}{statistics[i].TheDate, p})
	}

	if len(statistics)-1 > 0 {
		resp.Loss.Avg = float32(int(resp.Loss.Avg*100)/(len(statistics)-1)) / 100
	}

	return
}

// 在线数据 - 在线人数
func (User) OnlineData(req dto.OnlineDataReq, resp *dto.OnlineDataResp) error {
	q := models.UserOnlinePtr.DB()
	var err error
	*resp, err = User{}.activeData(q, req)

	return err
}

// 在线数据 - 充值金额
func (User) RechargeData(req dto.RechargeDataReq, resp *dto.RechargeDataResp) error {
	q := models.UserOnlinePtr.DB()
	var err error
	*resp, err = User{}.activeData(q, req)

	return err
}

// 实时动态数据
func (User) activeData(q *gorm.DB, req dto.ActiveDataReq) (dto.ActiveDataResp, error) {
	var resp dto.ActiveDataResp

	minute := time.Now().Format("15:04")
	start := fmt.Sprintf("%s %s:00", req.Date, req.StartMinute)
	end := fmt.Sprintf("%s %s:59", req.Date, minute)

	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	if err != nil {
		return resp, err
	}
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	if err != nil {
		return resp, err
	}

	var list dto.ActiveQueryList
	q.Select("COUNT(DISTINCT(user_id)) as n, SUBSTRING(create_time, 12, 5) AS t").
		Where("create_time >= ? and create_time <= ?", start, end).Group("t").Find(&list)

	m := make(map[string]int, len(list))
	for i := range list {
		m[list[i].T] = list[i].N
		resp.Total += list[i].N
	}

	for {
		if t1.After(t2) {
			break
		}
		s := t1.Format("15:04")
		i, _ := m[s]
		// 假数据
		// if i == 0 {
		// 	i = rand.Intn(1000) + 100
		// 	resp.Total += i
		// }
		resp.List = append(resp.List, [2]string{s, strconv.Itoa(i)})
		t1 = t1.Add(time.Minute)
	}

	if len(resp.List) == 0 {
		resp.List = append(resp.List, [2]string{minute, "0"})
	}

	return resp, nil
}

func (u User) rangeStatistic(start, end string) models.Statistic {
	var rows []models.Statistic
	models.StatisticPtr.DB().Where("the_date >= ? and the_date <= ?", start, end).Find(&rows)

	var row models.Statistic
	for i := range rows {
		row.SignupCount += rows[i].SignupCount
		row.LoginCount += rows[i].LoginCount
		row.RechargeCount += rows[i].RechargeCount
		row.OxygenGenerateCount += rows[i].OxygenGenerateCount
		row.OxygenConsumeCount += rows[i].OxygenConsumeCount
		row.LeafGenerateCount += rows[i].LeafGenerateCount
		row.LeafConsumeCount += rows[i].LeafConsumeCount
	}

	return row
}

// 统计数据写入数据库
func (u User) Statistic(date string, force bool) (row models.Statistic) {
	isToday := date == time.Now().Format("2006-01-02")
	models.StatisticPtr.DB().Where("the_date = ?", date).Take(&row)
	// 非当日数据，直接返回
	if !force && !isToday {
		return
	}
	// 今日数据，减少频繁更新
	// if row.IsValid() && time.Since(row.UpdatedAt.ToTime()).Minutes() < 2 {
	// 	return
	// }

	var (
		sum struct {
			Sum int64
		}
		dateStart = date + " 00:00:00"
		dateEnd   = date + " 23:59:59"
	)

	// 注册
	models.UserAccountPtr.DB().Where("create_time BETWEEN ? and ?", date+" 00:00:00", date+" 23:59:59").Count(&row.SignupCount)
	row.RemainCount = []int{0: int(row.SignupCount), 7: 0}
	// 登录
	models.UserOnlinePtr.DB().Where("create_time BETWEEN ? and ?", dateStart, dateEnd).Select("COUNT(DISTINCT(user_id)) as sum").Take(&sum)
	row.LoginCount = sum.Sum

	// 氧气总量，过后不再统计
	models.WealthPtr.DB().Select("SUM(oxygen) as sum").Take(&sum)
	row.OxygenCount = sum.Sum

	// 叶子总量
	row.LeafCount = 0

	// 氧气产生
	models.WealthChangeLogPtr.DB().Where("change_time BETWEEN ? and ?", dateStart, dateEnd).
		Where("change_id = 1 and change_count > 0").Select("SUM(change_count) as sum").Take(&sum)
	row.OxygenGenerateCount = sum.Sum

	// 叶子产生
	models.WealthChangeLogPtr.DB().Where("change_time BETWEEN ? and ?", dateStart, dateEnd).
		Where("change_id = 2 and change_count > 0").Select("SUM(change_count) as sum").Take(&sum)
	row.LeafGenerateCount = sum.Sum

	// 叶子消耗
	models.WealthChangeLogPtr.DB().Where("change_time BETWEEN ? and ?", dateStart, dateEnd).
		Where("change_id = 2 and change_count < 0").Select("SUM(change_count) as sum").Take(&sum)
	row.LeafConsumeCount = -sum.Sum

	if row.IsValid() {
		models.StatisticPtr.DB().Model(&row).Updates(&row)
	} else {
		row.TheDate = date
		models.StatisticPtr.DB().Create(&row)
	}

	return
}

// date 在 subDay 天前注册的用户， 在 date 日期登录的人数
func (User) StatisticRemainCount(date string, subDay int) error {
	if subDay <= 0 || subDay > 7 {
		return errors.New("subDay 参数错误")
	}

	t, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return err
	}

	d1 := t.AddDate(0, 0, -subDay).Format("2006-01-02")
	// d1注册用户ID
	var uids []int
	models.UserAccountPtr.DB().Where("create_time >= ? and create_time <= ?", d1+" 00:00:00", d1+" 23:59:59").Pluck("user_id", &uids)

	if len(uids) == 0 {
		return nil
	}

	// 在d2的登录人数
	var count int64
	models.UserOnlinePtr.DB().Where("user_id in (?) and create_time between ? and ?", uids, date+" 00:00:00", date+" 23:59:59").
		Pluck("COUNT(DISTINCT(user_id))", &count)

	// 更改d1的统计数据
	var row models.Statistic
	row.DB().Where("the_date = ?", d1).Take(&row)

	if !row.IsValid() {
		return errors.New("not found")
	}

	// 当日占第0位。第7天占第7位
	row.RemainCount[subDay] = int(count)
	b, _ := json.Marshal(row.RemainCount)
	row.DB().Model(&row).UpdateColumn("remain_count", string(b))

	return nil
}

func (User) DailyStatistics(req dto.DailyStatisticsListReq, resp *dto.DailyStatisticsListResq) error {
	var sumRegisterAccountCount, sumRegisterDayCount, sumLoginAccountCount int64
	//总注册量
	var ua *models.UserAccount
	ua.DB().Count(&sumRegisterAccountCount)

	log.Infof("%+v", req)

	if sumRegisterAccountCount <= 0 {
		return nil
	}

	//第一个注册
	var uaObj models.UserAccount
	ua.DB().First(&uaObj)

	firstDay, _ := time.Parse("2006-01-02", strings.Split(uaObj.CreateTime, " ")[0])
	sumRegisterDayCount = int64(time.Now().Sub(firstDay).Hours())/24 + 1

	endOffset := int(math.Min(float64((req.Pagination.GetPageIndex()-1)*req.Pagination.GetPageSize()), float64(sumRegisterDayCount)))
	startOffset := int(math.Min(float64(endOffset+req.Pagination.GetPageSize()-1), float64(sumRegisterDayCount)))

	startDay, endDay := time.Now().AddDate(0, 0, -startOffset), time.Now().AddDate(0, 0, -endOffset)

	//总登录人数
	var ull *models.UserLoginLog
	ull.DB().Raw("SELECT user_id FROM login_log GROUP BY user_id;").Count(&sumLoginAccountCount)

	//每日注册人数
	sql := fmt.Sprintf("SELECT date(create_time) as date, count(create_time) as register_count FROM forest_user_center.user_account "+
		"WHERE date(create_time) BETWEEN \"%v\" AND \"%v\" GROUP BY date ORDER BY date DESC;",
		startDay.Format("2006-01-02"), endDay.Format("2006-01-02"))
	var registerList []*dto.DailyStatistics
	ua.DB().Raw(sql).Find(&registerList)

	//每日登录人数及活跃时长
	sql = fmt.Sprintf("SELECT user_id, reason, login_time, logout_time, date(login_time) as date FROM forest_logic.login_log "+
		"WHERE date(login_time) BETWEEN \"%v\" AND \"%v\";",
		startDay.Format("2006-01-02"), endDay.Format("2006-01-02"))
	var userLoginLogs []*models.UserLoginLog
	ull.DB().Raw(sql).Find(&userLoginLogs)
	loginLogs := calcUserLogin(userLoginLogs)

	//每日首次登录人数
	var ula *models.UserLogicAccount
	sql = fmt.Sprintf("SELECT user_id, date(reg_time) as reg_time FROM forest_logic.account WHERE date(reg_time) BETWEEN \"%v\" AND \"%v\";", startDay.Format("2006-01-02"), endDay.Format("2006-01-02"))
	var userLogicAccounts []*models.UserLogicAccount
	ula.DB().Raw(sql).Find(&userLogicAccounts)
	userLogicAccountsMap := make(map[string]map[int]*models.UserLogicAccount, len(userLogicAccounts))
	for _, account := range userLogicAccounts {
		if _, ok := userLogicAccountsMap[account.RegTime]; !ok {
			userLogicAccountsMap[account.RegTime] = make(map[int]*models.UserLogicAccount, len(userLogicAccounts))
		}
		userLogicAccountsMap[account.RegTime][account.UserId] = account
	}
	//log.Infof("userLogicAccounts=%+v,\nuserLogicAccountsMap=%+v", userLogicAccounts[0], userLogicAccountsMap)

	for endDay.Unix() > startDay.Unix() {
		bHas := false
		day := endDay.Format("2006-01-02")
		endDay = endDay.AddDate(0, 0, -1)
		var register *dto.DailyStatistics
		for _, registerTmp := range registerList {
			if registerTmp.Date == day {
				bHas = true
				register = registerTmp
				break
			}
		}
		if !bHas {
			register = &dto.DailyStatistics{
				SysId:           0,
				Date:            day,
				RegisterCount:   0,
				LoginCount:      0,
				FirstLoginCount: 0,
				MaxOnlineCount:  0,
				OnlineCount1:    0,
				OnlineCount2:    0,
				OnlineCount3:    0,
			}
		}
		if dayLoginLogs, ok := loginLogs[day]; ok {
			register.LoginCount = len(dayLoginLogs)
			for _, loginLog := range dayLoginLogs {
				//log.Warnf("loginLog.ActiveDuration=%v\n", loginLog.ActiveDuration)
				if loginLog.ActiveDuration >= int64(30*60) {
					register.OnlineCount3++
				} else if loginLog.ActiveDuration >= int64(15*60) {
					register.OnlineCount2++
				} else if loginLog.ActiveDuration >= int64(3*60) {
					register.OnlineCount1++
				}
				if dayUserLogicAccountsMap, ok2 := userLogicAccountsMap[day]; ok2 {
					if _, ok3 := dayUserLogicAccountsMap[loginLog.UserId]; ok3 {
						register.FirstLoginCount++
					}
				}
			}
		}
		resp.List = append(resp.List, *register)
	}

	resp.Count = sumRegisterDayCount
	resp.SumRegisterCount = sumRegisterAccountCount
	resp.SumLoginCount = sumLoginAccountCount

	log.Infof("%+v\n", resp)

	return nil
}

func calcUserLogin(loginLogs []*models.UserLoginLog) map[string]map[int]*models.UserLoginLog {
	userLoginLogsMap := make(map[string]map[int][]*models.UserLoginLog, len(loginLogs))
	for _, loginLog := range loginLogs {
		loginLog.Date = strings.Split(loginLog.Date, " ")[0]
		if _, ok := userLoginLogsMap[loginLog.Date]; !ok {
			userLoginLogsMap[loginLog.Date] = make(map[int][]*models.UserLoginLog, len(loginLogs))
		}
		if _, ok := userLoginLogsMap[loginLog.Date][loginLog.UserId]; !ok {
			userLoginLogsMap[loginLog.Date][loginLog.UserId] = make([]*models.UserLoginLog, 0, len(loginLogs))
		}
		userLoginLogsMap[loginLog.Date][loginLog.UserId] = append(userLoginLogsMap[loginLog.Date][loginLog.UserId], loginLog)
	}
	userLoginLogMap := make(map[string]map[int]*models.UserLoginLog, len(userLoginLogsMap))
	for _, userLoginLogMap2 := range userLoginLogsMap {
		for _, loginLogSl := range userLoginLogMap2 {
			for idx, loginLog := range loginLogSl {
				if _, ok := userLoginLogMap[loginLog.Date]; !ok {
					userLoginLogMap[loginLog.Date] = make(map[int]*models.UserLoginLog, len(userLoginLogsMap))
				}
				if _, ok := userLoginLogMap[loginLog.Date][loginLog.UserId]; !ok {
					userLoginLogMap[loginLog.Date][loginLog.UserId] = &models.UserLoginLog{
						UserId: loginLog.UserId,
					}
				}
				loginTime, err1 := time.Parse("2006-01-02 15:04:05", loginLog.LoginTime)
				logoutTime, err2 := time.Parse("2006-01-02 15:04:05", loginLog.LogoutTime)
				if err1 != nil || err2 != nil {
					log.Errorf("time Parse failed, err1=%v, err2=%v", err1, err2)
					continue
				}
				duration := logoutTime.Unix() - loginTime.Unix()
				//log.Info("duration=", duration, "logoutTime=", logoutTime, "loginTime", loginTime)
				if loginLog.Reason == 0 && idx == len(loginLogSl) && duration <= 0 {
					dayStartTime, _ := time.Parse("2006-01-02", loginLog.Date)
					duration = dayStartTime.AddDate(0, 0, 1).Unix() - loginTime.Unix()
				}
				if duration < 0 {
					duration = 0
				}
				userLoginLogMap[loginLog.Date][loginLog.UserId].ActiveDuration += duration
			}
		}
	}
	log.Infof("userLoginLogMap=%+v", userLoginLogMap)
	return userLoginLogMap
}
