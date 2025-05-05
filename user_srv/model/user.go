package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32          `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
	IsDeleted bool           `json:"is_deleted" gorm:"column:is_deleted"`
}

/*
*
 1. 用户密码保存肯定不是明文的，采用md5加密 (对称加密：加载密钥，解密和加密都用同一个密钥)
 2. 用户密码加密后存储到数据库中
*/
type User struct {
	BaseModel
	Mobile   string     `json:"mobile" gorm:"column:mobile;unique;index:idx_mobile;not null"`
	Email    string     `json:"email" gorm:"column:email;unique;index:idx_email;not null"`
	Password string     `json:"password" gorm:"column:password;not null;type:varchar(100)"`
	NickName string     `json:"nick_name" gorm:"column:nick_name;type:varchar(20)"`
	UserName string     `json:"user_name" gorm:"column:user_name;type:varchar(20);not null;unique;index:idx_user_name"`
	Birthday *time.Time `json:"birthday" gorm:"column:birthday;type:datetime"`
	Gender   string     `json:"gender" gorm:"column:gender;type:varchar(10);default:'unknown';comment:'unknown:未知,male:男,female:女'"`
	Avatar   string     `json:"avatar" gorm:"column:avatar;type:varchar(100);default:'https://cdn.jsdelivr.net/gh/zhangguanzhang/picgo@master/img/20230920161420.png'"`
	Role     int        `json:"role" gorm:"column:role;type:int;default:1;comment:'1:普通用户,2:管理员'"`
}
