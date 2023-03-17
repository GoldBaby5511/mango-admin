package service

import (
	"errors"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/service"
)

const (
	debug = false
)

var (
	game = Game{}
)

type Game struct {
	service.Service
}

var (
	mu     sync.Mutex
	start0 = [...]string{
		// 想念你 一厢情愿 爱我就了解我 其实不想走 爱你久久久 被爱就是幸福
		"360", "1372", "259695", "74839", "20999", "829475",
		"668", "518", "588", "886", "881", "800", "400", "600",
	}
	start1 = [...]string{
		"00188", "22188", "33188", "00168", "22168", "33168", "00186", "22186", "33186",
	}
	contains = [...]string{
		"000", "111", "222", "333", "444", "555", "666", "777", "888", "999",
		"110", "119", "120", "123", "520", "521", "527",
		// 一生一世 我想亲亲你 我深情依旧 爱我一辈子 一生就爱你一人 爱爱你爱爱我 我发誓我爱你 爱是如此神奇 我就是爱想你
		"1314", "53770", "53719", "25184", "1392010", "220250", "584520", "246437", "594230",
		"4321", "806", "668", "986", "998",
	}
	regPhone    = regexp.MustCompile(`^1[3456789]\d{9}$`)
	regBirthday = regexp.MustCompile(`(19[4-9]|20[0-4])[0-9](0[1-9]|1[0-2])`) // 1940-01 ~ 2049-12
)

func init() {
	// a := regBirthday.MatchString("xx193912xx")
	// b := regBirthday.MatchString("xx194001xx")
	// c := regBirthday.MatchString("xx204912xx")
	// d := regBirthday.MatchString("xx205001xx")
	// fmt.Println(a, b, c, d)
}

func (g *Game) GetMax() (max int64) {
	g.Orm.Model(&models.GameIdExcellent{}).Pluck("MAX(game_id)", &max)
	return max
}

func (g *Game) GetPage(req *dto.GameIdListReq, list *[]models.GameIdNormal, count *int64) error {
	var model interface{}
	if req.IsSpecial {
		model = &models.GameIdExcellent{}
	} else {
		model = &models.GameIdNormal{}
	}

	q := g.Orm.Model(model).Scopes(
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
	)

	if req.UserId > 0 {
		q = q.Where("user_id = ?", req.UserId)
	}
	if req.GameId > 0 {
		q = q.Where("game_id = ?", req.GameId)
	}

	return q.Find(list).Count(count).Error
}

func (g *Game) GenerateId(start, end int) error {
	mu.Lock()
	defer mu.Unlock()

	if start < 1 {
		return errors.New("start must be greater than 10000")
	}

	if start > end {
		return errors.New("start must less than end")
	}

	if end > 100000 {
		return errors.New("end must less than 1000000000")
	}

	if end-start > 30 {
		return errors.New("最大每次生成30万个")
	}

	start *= 10000
	end *= 10000

	var row models.GameIdExcellent
	g.Orm.Model(models.GameIdExcellent{}).Where("game_id >= ? and game_id < ?", start, end).Take(&row)
	if row.SysId != 0 {
		return errors.New("包含已存在数字，请根据最大值设定开始和结束值")
	}

	normarl := make([]int, 0, end-start)
	special := make([]int, 0, end-start)

	t := time.Now()

	for i := start; i < end; i++ {
		n := strconv.Itoa(i)

		flag := true

		for j := range start0 {
			if !flag {
				break
			}
			if strings.HasPrefix(n, start0[j]) {
				special = append(special, i)
				flag = false
			}
		}

		for j := range start1 {
			if !flag {
				break
			}
			if strings.HasPrefix(n[1:], start1[j]) {
				special = append(special, i)
				flag = false
			}
		}

		for j := range contains {
			if !flag {
				break
			}
			if strings.Contains(n, contains[j]) {
				special = append(special, i)
				flag = false
			}
		}

		// 手机号
		if flag && len(n) == 11 && regPhone.MatchString(n) {
			special = append(special, i)
			continue
		}

		// 出生年月
		if flag && len(n) >= 6 && regBirthday.MatchString(n) {
			special = append(special, i)
			continue
		}

		// _AA_AA
		x := n[1]
		if flag && len(n) >= 6 && n[2] == x && n[4] == x && n[5] == x {
			special = append(special, i)
			continue
		}

		if flag {
			normarl = append(normarl, i)
		}
	}

	Shuffle(normarl)
	Shuffle(special)

	log.Info("生成游戏ID:", start, "-", end, ", 耗时:", time.Since(t).Milliseconds(), "ms")
	log.Info("normarl:", len(normarl), "special:", len(special), "total:", len(normarl)+len(special))

	m1 := make([]models.GameIdNormal, len(normarl))
	m2 := make([]models.GameIdExcellent, len(special))

	for i := 0; i < len(normarl); i++ {
		m1[i].GameId = uint64(normarl[i])
	}

	for i := 0; i < len(special); i++ {
		m2[i].GameId = uint64(special[i])
	}

	if err := g.Orm.CreateInBatches(m1, 2000).Error; err != nil {
		return err
	}
	if err := g.Orm.CreateInBatches(m2, 2000).Error; err != nil {
		return err
	}

	return nil
}

func Shuffle(slice []int) {
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := rander.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}
