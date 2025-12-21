package api

import (
	"net/http"
	"play-tools/database"
	"play-tools/model"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	Name       string `json:"name" binding:"required"`
	MaxPlayers int    `json:"maxPlayers" binding:"required,min=2,max=10"`
}

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	RoomID string `json:"roomId" binding:"required"`
}

// AddTransactionRequest 添加交易请求
type AddTransactionRequest struct {
	RoomID   string `json:"roomId" binding:"required"`
	ToUserID uint   `json:"toUserId" binding:"required"`
	Amount   int    `json:"amount" binding:"required"`
}

// FinishGameRequest 结束游戏请求
type FinishGameRequest struct {
	RoomID string `json:"roomId" binding:"required"`
}

// RoomDetail 房间详情
type RoomDetail struct {
	model.GameRoom
	Members      []MemberInfo      `json:"members"`
	Transactions []model.Transaction `json:"transactions"`
	IsOwner      bool              `json:"isOwner"`
	IsMember     bool              `json:"isMember"`
}

// MemberInfo 成员信息
type MemberInfo struct {
	UserID   uint   `json:"userId"`
	Nickname string `json:"nickname"`
	IsReady  bool   `json:"isReady"`
	Balance  int    `json:"balance"` // 计算得出的余额
}

// Settlement 结算结果
type Settlement struct {
	FromUserID   uint   `json:"fromUserId"`
	FromNickname string `json:"fromNickname"`
	ToUserID     uint   `json:"toUserId"`
	ToNickname   string `json:"toNickname"`
	Amount       int    `json:"amount"`
}

// CreateRoom 创建房间
func CreateRoom(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 创建房间
	roomID := uuid.New().String()[:8] // 使用短ID方便分享
	room := model.GameRoom{
		ID:         roomID,
		Name:       req.Name,
		OwnerID:    userID.(uint),
		MaxPlayers: req.MaxPlayers,
		Status:     "waiting",
	}

	if err := database.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建房间失败",
		})
		return
	}

	// 房主自动加入房间
	var user model.User
	database.DB.First(&user, userID)

	member := model.RoomMember{
		RoomID:   roomID,
		UserID:   userID.(uint),
		Nickname: user.Nickname,
		IsReady:  true, // 房主自动准备
	}
	database.DB.Create(&member)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": room,
	})
}

// JoinRoom 加入房间
func JoinRoom(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 检查房间是否存在
	var room model.GameRoom
	if err := database.DB.First(&room, "id = ?", req.RoomID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "房间不存在",
		})
		return
	}

	// 检查房间是否已满
	var memberCount int64
	database.DB.Model(&model.RoomMember{}).Where("room_id = ?", req.RoomID).Count(&memberCount)
	if int(memberCount) >= room.MaxPlayers {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "房间已满",
		})
		return
	}

	// 检查是否已经加入
	var existMember model.RoomMember
	err := database.DB.Where("room_id = ? AND user_id = ?", req.RoomID, userID).First(&existMember).Error
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "已经加入该房间",
		})
		return
	}

	// 加入房间
	var user model.User
	database.DB.First(&user, userID)

	member := model.RoomMember{
		RoomID:   req.RoomID,
		UserID:   userID.(uint),
		Nickname: user.Nickname,
		IsReady:  false, // 新加入的成员需要准备
	}

	if err := database.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "加入房间失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "加入成功",
		"data": room,
	})
}

// GetRoomInfo 获取房间信息
func GetRoomInfo(c *gin.Context) {
	userID, _ := c.Get("userId")
	roomID := c.Param("roomId")

	var room model.GameRoom
	if err := database.DB.First(&room, "id = ?", roomID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "房间不存在",
		})
		return
	}

	// 获取所有成员
	var members []model.RoomMember
	database.DB.Where("room_id = ?", roomID).Find(&members)

	// 获取所有交易记录
	var transactions []model.Transaction
	database.DB.Where("room_id = ?", roomID).Order("created_at").Find(&transactions)

	// 计算每个成员的余额
	balanceMap := make(map[uint]int)
	for _, tx := range transactions {
		balanceMap[tx.FromUserID] -= tx.Amount
		balanceMap[tx.ToUserID] += tx.Amount
	}

	// 构建成员信息列表
	memberInfos := make([]MemberInfo, 0, len(members))
	for _, member := range members {
		memberInfos = append(memberInfos, MemberInfo{
			UserID:   member.UserID,
			Nickname: member.Nickname,
			IsReady:  member.IsReady,
			Balance:  balanceMap[member.UserID],
		})
	}

	// 检查当前用户是否是房主和成员
	isOwner := room.OwnerID == userID.(uint)
	isMember := false
	for _, member := range members {
		if member.UserID == userID.(uint) {
			isMember = true
			break
		}
	}

	detail := RoomDetail{
		GameRoom:     room,
		Members:      memberInfos,
		Transactions: transactions,
		IsOwner:      isOwner,
		IsMember:     isMember,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": detail,
	})
}

// AddTransaction 添加交易记录
func AddTransaction(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req AddTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 检查房间是否存在
	var room model.GameRoom
	if err := database.DB.First(&room, "id = ?", req.RoomID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "房间不存在",
		})
		return
	}

	// 检查是否是房间成员
	var member model.RoomMember
	if err := database.DB.Where("room_id = ? AND user_id = ?", req.RoomID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 403,
			"msg":  "你不是房间成员",
		})
		return
	}

	// 检查收款人是否是房间成员
	var toMember model.RoomMember
	if err := database.DB.Where("room_id = ? AND user_id = ?", req.RoomID, req.ToUserID).First(&toMember).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "收款人不是房间成员",
		})
		return
	}

	// 不能给自己转账
	if userID.(uint) == req.ToUserID {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "不能给自己转账",
		})
		return
	}

	// 更新房间状态为进行中
	if room.Status == "waiting" {
		database.DB.Model(&room).Update("status", "playing")
	}

	// 创建交易记录
	transaction := model.Transaction{
		RoomID:       req.RoomID,
		FromUserID:   userID.(uint),
		ToUserID:     req.ToUserID,
		Amount:       req.Amount,
		FromNickname: member.Nickname,
		ToNickname:   toMember.Nickname,
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建交易记录失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
		"data": transaction,
	})
}

// FinishGame 结束游戏并结算
func FinishGame(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req FinishGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 获取房间
	var room model.GameRoom
	if err := database.DB.First(&room, "id = ?", req.RoomID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "房间不存在",
		})
		return
	}

	// 只有房主可以结束游戏
	if room.OwnerID != userID.(uint) {
		c.JSON(http.StatusOK, gin.H{
			"code": 403,
			"msg":  "只有房主可以结束游戏",
		})
		return
	}

	// 更新房间状态
	database.DB.Model(&room).Update("status", "finished")

	// 获取所有成员和交易记录
	var members []model.RoomMember
	database.DB.Where("room_id = ?", req.RoomID).Find(&members)

	var transactions []model.Transaction
	database.DB.Where("room_id = ?", req.RoomID).Find(&transactions)

	// 计算结算方案
	settlements := calculateSettlements(members, transactions)

	// 保存结算结果到数据库
	for _, s := range settlements {
		settlement := model.Settlement{
			RoomID:       req.RoomID,
			FromUserID:   s.FromUserID,
			FromNickname: s.FromNickname,
			ToUserID:     s.ToUserID,
			ToNickname:   s.ToNickname,
			Amount:       s.Amount,
		}
		database.DB.Create(&settlement)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "结算成功",
		"data": settlements,
	})
}

// calculateSettlements 计算最优结算方案
func calculateSettlements(members []model.RoomMember, transactions []model.Transaction) []Settlement {
	// 计算每个人的余额
	balanceMap := make(map[uint]int)
	nicknameMap := make(map[uint]string)

	for _, member := range members {
		balanceMap[member.UserID] = 0
		nicknameMap[member.UserID] = member.Nickname
	}

	for _, tx := range transactions {
		balanceMap[tx.FromUserID] -= tx.Amount
		balanceMap[tx.ToUserID] += tx.Amount
	}

	// 分离欠款人和收款人
	type BalanceInfo struct {
		UserID   uint
		Nickname string
		Balance  int
	}

	var debtors []BalanceInfo   // 欠款人（余额为负）
	var creditors []BalanceInfo // 收款人（余额为正）

	for userID, balance := range balanceMap {
		if balance < 0 {
			debtors = append(debtors, BalanceInfo{
				UserID:   userID,
				Nickname: nicknameMap[userID],
				Balance:  -balance, // 转为正数
			})
		} else if balance > 0 {
			creditors = append(creditors, BalanceInfo{
				UserID:   userID,
				Nickname: nicknameMap[userID],
				Balance:  balance,
			})
		}
	}

	// 排序
	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].Balance > debtors[j].Balance
	})
	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].Balance > creditors[j].Balance
	})

	// 贪心算法匹配
	settlements := make([]Settlement, 0)
	i, j := 0, 0

	for i < len(debtors) && j < len(creditors) {
		amount := debtors[i].Balance
		if creditors[j].Balance < amount {
			amount = creditors[j].Balance
		}

		settlements = append(settlements, Settlement{
			FromUserID:   debtors[i].UserID,
			FromNickname: debtors[i].Nickname,
			ToUserID:     creditors[j].UserID,
			ToNickname:   creditors[j].Nickname,
			Amount:       amount,
		})

		debtors[i].Balance -= amount
		creditors[j].Balance -= amount

		if debtors[i].Balance == 0 {
			i++
		}
		if creditors[j].Balance == 0 {
			j++
		}
	}

	return settlements
}

// GetMyRooms 获取我的房间列表
func GetMyRooms(c *gin.Context) {
	userID, _ := c.Get("userId")

	// 获取我加入的所有房间
	var members []model.RoomMember
	database.DB.Where("user_id = ?", userID).Find(&members)

	roomIDs := make([]string, 0, len(members))
	for _, member := range members {
		roomIDs = append(roomIDs, member.RoomID)
	}

	if len(roomIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": []model.GameRoom{},
		})
		return
	}

	var rooms []model.GameRoom
	database.DB.Where("id IN ?", roomIDs).Order("created_at DESC").Find(&rooms)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": rooms,
	})
}

// GetRoomSettlements 获取房间结算结果
func GetRoomSettlements(c *gin.Context) {
	roomID := c.Param("roomId")

	// 检查房间是否存在
	var room model.GameRoom
	if err := database.DB.First(&room, "id = ?", roomID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "房间不存在",
		})
		return
	}

	// 检查房间是否已结束
	if room.Status != "finished" {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "房间未结束",
		})
		return
	}

	// 获取结算记录
	var settlements []model.Settlement
	database.DB.Where("room_id = ?", roomID).Find(&settlements)

	// 转换为API响应格式
	result := make([]Settlement, 0, len(settlements))
	for _, s := range settlements {
		result = append(result, Settlement{
			FromUserID:   s.FromUserID,
			FromNickname: s.FromNickname,
			ToUserID:     s.ToUserID,
			ToNickname:   s.ToNickname,
			Amount:       s.Amount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": result,
	})
}

// ToggleReady 切换准备状态
func ToggleReady(c *gin.Context) {
	userID, _ := c.Get("userId")
	roomID := c.Param("roomId")

	// 检查房间是否存在
	var room model.GameRoom
	if err := database.DB.First(&room, "id = ?", roomID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "房间不存在",
		})
		return
	}

	// 检查是否是房间成员
	var member model.RoomMember
	if err := database.DB.Where("room_id = ? AND user_id = ?", roomID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 403,
			"msg":  "你不是房间成员",
		})
		return
	}

	// 房主不能取消准备
	if room.OwnerID == userID.(uint) {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "房主始终处于准备状态",
		})
		return
	}

	// 切换准备状态
	database.DB.Model(&member).Update("is_ready", !member.IsReady)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "状态已更新",
		"data": gin.H{"isReady": !member.IsReady},
	})
}

// StartGame 开始游戏
func StartGame(c *gin.Context) {
	userID, _ := c.Get("userId")
	roomID := c.Param("roomId")

	// 检查房间是否存在
	var room model.GameRoom
	if err := database.DB.First(&room, "id = ?", roomID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "房间不存在",
		})
		return
	}

	// 只有房主可以开始游戏
	if room.OwnerID != userID.(uint) {
		c.JSON(http.StatusOK, gin.H{
			"code": 403,
			"msg":  "只有房主可以开始游戏",
		})
		return
	}

	// 检查所有成员是否准备
	var members []model.RoomMember
	database.DB.Where("room_id = ?", roomID).Find(&members)

	for _, member := range members {
		if !member.IsReady {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "还有成员未准备",
			})
			return
		}
	}

	// 更新房间状态为游戏中
	database.DB.Model(&room).Update("status", "playing")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "游戏已开始",
	})
}
