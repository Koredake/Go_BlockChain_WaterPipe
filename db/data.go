package db

type Admin struct {
	ID         int    `sql:"id" json:"id" form:"id"`
	MaxEther   int    `sql:"maxEther" json:"maxEther" form:"maxEther"`
	TransEther int    `sql:"TransEther" json:"TransEther" form:"TransEther"`
	AddrAmount int    `sql:"AddrAmount" json:"AddrAmount" form:"AddrAmount"`
	Date       string `sql:"Date" json:"Date" form:"Date"`
}
type User struct {
	ID          int    `sql:"id" json:"id" form:"id"`
	Address     string `sql:"address" json:"address" form:"address"`
	ApplyCount  int    `sql:"applyCount" json:"applyCount" form:"applyCount"`
	ApplyTimes  int    `sql:"ApplyTimes" json:"ApplyTimes" form:"ApplyTimes"`
	LatestTrans int64  `sql:"LatestTrans" json:"LatestTrans" form:"LatestTrans"`
	Dates       string `sql:"dates" json:"dates" form:"dates"`
}
