package mysql

import (
	"blueself/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (CommunityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if err := db.Select(&CommunityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}


// 根据ID查询社区详情
func GetCommunityDetailByID(id int64)(community *models.CommunityDetail,err error){
	community = new(models.CommunityDetail)
	sqlStr := `select 
	community_id, community_name, introduction, create_time 
	from community 
	where community_id = ?`
	if er := db.Get(community, sqlStr, id); er != nil{
		if er == sql.ErrNoRows{
			err = ErrInvalidID
		}
	}
	return community,err
}