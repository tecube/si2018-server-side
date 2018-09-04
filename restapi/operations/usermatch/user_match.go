package usermatch

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/restapi/operations/util"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/runtime/middleware"
	"fmt"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	limit := int(p.Limit)
	offset := int(p.Offset)
	token := p.Token
	id, err := util.GetIDByToken(token)
	//id := int64(1)
	if err != nil {
		fmt.Print("Get id err: ")
		fmt.Println(err)
	}

	ru := repositories.NewUserRepository()
	rm := repositories.NewUserMatchRepository()
	
	matches, err := rm.FindByUserIDWithLimitOffset(id, limit, offset)
	if err != nil {
		fmt.Print("Find matches err: ")
		fmt.Println(err)
	}

	ids := make([]int64, 0) /* TODO can use map's key as ids slice? */
	for _, mat := range matches {
		if mat.UserID == id {
			ids = append(ids, mat.PartnerID)
		} else if mat.PartnerID == id {
			ids = append(ids, mat.UserID)
		}
	}

	users, err := ru.FindByIDs(ids)
	if err != nil {
		fmt.Print("Find users by ids err: ")
		fmt.Println(err)
	}

	um := make(map[int64]entities.User)
	for _, u := range users {
		um[u.ID] = u
	}

	res := make([]entities.MatchUserResponse, 0)
	
	for _, mat := range matches{
		r := entities.MatchUserResponse{}
		r.MatchedAt =  mat.CreatedAt

		var u entities.User
		if mat.UserID == id {
			u = um[mat.PartnerID]
		} else if mat.PartnerID == id {
			u = um[mat.UserID]
		} 

		r.ApplyUser(u)
		res = append(res, r)
	}

	var reses entities.MatchUserResponses
	reses = res

	return si.NewGetMatchesOK().WithPayload(reses.Build())
}
