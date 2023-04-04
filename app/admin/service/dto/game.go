package dto

import "mango-admin/common/dto"

type GameIdListReq struct {
	dto.Pagination `search:"-"`
	IsSpecial      bool `form:"is_special"`
	UserId         int  `form:"user_id"`
	GameId         int  `form:"game_id"`
}
