// 类型定义

// 用户类型
export interface User {
  userId: number      // 用户唯一ID，8位正整数
  nickname: string
  avatar?: string
  createdAt: string
}

// 房间成员类型
export interface Member {
  userId: number
  nickname: string
  avatar?: string  // 用户头像
  isReady: boolean
  balance: number
}

// 游戏房间类型
export interface GameRoom {
  id: string
  name: string
  ownerId: number
  maxPlayers: number
  status: 'waiting' | 'playing' | 'finished'
  createdAt: string
  updatedAt: string
}

// 房间详情类型
export interface RoomDetail {
  id: string
  name: string
  ownerId: number
  maxPlayers: number
  status: 'waiting' | 'playing' | 'finished'
  createdAt: string
  updatedAt: string
  members: Member[]
  transactions: Transaction[]
  isOwner: boolean
  isMember: boolean
}

// 交易记录类型
export interface Transaction {
  id: number
  roomId: string
  fromUserId: number
  toUserId: number
  amount: number
  fromNickname: string
  toNickname: string
  createdAt: string
}

// 结算结果类型
export interface Settlement {
  fromUserId: number
  fromNickname: string
  toUserId: number
  toNickname: string
  amount: number
}
