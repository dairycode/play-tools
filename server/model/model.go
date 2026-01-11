package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"userId"`   // 用户唯一ID，8位正整数
	Nickname  string    `gorm:"type:varchar(50);not null" json:"nickname"` // 昵称，必填
	Avatar    string    `gorm:"type:varchar(500)" json:"avatar"`      // 头像URL
	WxOpenID  string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"-"` // 微信OpenID，不返回给前端
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GameRoom 游戏房间模型
type GameRoom struct {
	ID         string    `gorm:"type:varchar(8);primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`        // 房间名称
	OwnerID    uint      `gorm:"index;not null" json:"ownerId"`                 // 房主ID
	MaxPlayers int       `gorm:"not null" json:"maxPlayers"`                    // 房间人数上限
	Status     string    `gorm:"type:varchar(20);index;not null" json:"status"` // waiting, playing, finished
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// RoomMember 房间成员模型
type RoomMember struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	RoomID   string `gorm:"type:varchar(8);index;not null" json:"roomId"`
	UserID   uint   `gorm:"index;not null" json:"userId"`
	Nickname string `gorm:"type:varchar(50);not null" json:"nickname"` // 冗余存储，方便查询
	Avatar   string `gorm:"type:varchar(500)" json:"avatar"`           // 冗余存储用户头像
	IsReady  bool   `gorm:"default:false" json:"isReady"`              // 是否准备
}

// Transaction 交易记录模型
type Transaction struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoomID       string    `gorm:"type:varchar(8);index;not null" json:"roomId"`
	FromUserID   uint      `gorm:"index;not null" json:"fromUserId"` // 付款人用户ID
	ToUserID     uint      `gorm:"index;not null" json:"toUserId"`   // 收款人用户ID
	Amount       int       `gorm:"not null" json:"amount"`
	FromNickname string    `gorm:"type:varchar(50);not null" json:"fromNickname"` // 冗余存储
	ToNickname   string    `gorm:"type:varchar(50);not null" json:"toNickname"`   // 冗余存储
	CreatedAt    time.Time `json:"createdAt"`
}

// Settlement 结算记录模型
type Settlement struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoomID       string    `gorm:"type:varchar(8);index;not null" json:"roomId"`
	FromUserID   uint      `gorm:"not null" json:"fromUserId"`
	FromNickname string    `gorm:"type:varchar(50);not null" json:"fromNickname"`
	ToUserID     uint      `gorm:"not null" json:"toUserId"`
	ToNickname   string    `gorm:"type:varchar(50);not null" json:"toNickname"`
	Amount       int       `gorm:"not null" json:"amount"`
	CreatedAt    time.Time `json:"createdAt"`
}
