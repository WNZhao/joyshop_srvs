package model

// UserFav 用户收藏表
type UserFav struct {
	BaseModel
	User  int32 `gorm:"type:int;index:idx_user_goods,unique"`
	Goods int32 `gorm:"type:int;index:idx_user_goods,unique"`
}

func (UserFav) TableName() string {
	return "userfav"
}

// Address 收货地址表
type Address struct {
	BaseModel
	User         int32  `gorm:"type:int;index"`
	Province     string `gorm:"type:varchar(10)"`
	City         string `gorm:"type:varchar(10)"`
	District     string `gorm:"type:varchar(20)"`
	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"`
	SignerMobile string `gorm:"type:varchar(11)"`
}

func (Address) TableName() string {
	return "address"
}

// LeavingMessage 留言表
type LeavingMessage struct {
	BaseModel
	User        int32  `gorm:"type:int;index"`
	MessageType int32  `gorm:"type:int comment '留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)'"`
	Subject     string `gorm:"type:varchar(100)"`
	Message     string `gorm:"type:text"`
	File        string `gorm:"type:varchar(200) comment '上传的文件'"`
}

func (LeavingMessage) TableName() string {
	return "leavingmessages"
}