package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Nickname  string    `gorm:"not null" json:"nickname"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GameRoom 游戏房间模型
type GameRoom struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`                 // 房间名称
	OwnerID     uint      `gorm:"index" json:"ownerId"` // 房主ID
	MaxPlayers  int       `json:"maxPlayers"`           // 房间人数上限
	Status      string    `gorm:"index" json:"status"`  // waiting, playing, finished
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// RoomMember 房间成员模型
type RoomMember struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	RoomID   string    `gorm:"index" json:"roomId"`
	UserID   uint      `gorm:"index" json:"userId"`
	Nickname string    `json:"nickname"` // 冗余存储，方便查询
	IsReady  bool      `json:"isReady"`  // 是否准备
	JoinedAt time.Time `json:"joinedAt"`
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Transaction 交易记录模型
type Transaction struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoomID       string    `gorm:"index" json:"roomId"`
	FromUserID   uint      `gorm:"index" json:"fromUserId"` // 付款人用户ID
	ToUserID     uint      `gorm:"index" json:"toUserId"`   // 收款人用户ID
	Amount       int       `json:"amount"`
	FromNickname string    `json:"fromNickname"` // 冗余存储
	ToNickname   string    `json:"toNickname"`   // 冗余存储
	CreatedAt    time.Time `json:"createdAt"`
}

// Settlement 结算记录模型
type Settlement struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoomID       string    `gorm:"index" json:"roomId"`
	FromUserID   uint      `json:"fromUserId"`
	FromNickname string    `json:"fromNickname"`
	ToUserID     uint      `json:"toUserId"`
	ToNickname   string    `json:"toNickname"`
	Amount       int       `json:"amount"`
	CreatedAt    time.Time `json:"createdAt"`
}
