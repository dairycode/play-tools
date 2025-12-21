// 游戏相关API
import { request } from '@/utils/request'
import type { GameRoom, RoomDetail, Transaction, Settlement } from '@/types'

// 创建游戏房间
export const createRoom = (data: { name: string; maxPlayers: number }) => {
  return request<GameRoom>({
    url: '/game/room/create',
    method: 'POST',
    data
  })
}

// 加入房间
export const joinRoom = (data: { roomId: string }) => {
  return request<GameRoom>({
    url: '/game/room/join',
    method: 'POST',
    data
  })
}

// 获取房间信息
export const getRoomInfo = (roomId: string) => {
  return request<RoomDetail>({
    url: `/game/room/${roomId}`,
    method: 'GET'
  })
}

// 添加交易记录
export const addTransaction = (data: {
  roomId: string
  toUserId: number
  amount: number
}) => {
  return request<Transaction>({
    url: '/game/transaction/add',
    method: 'POST',
    data
  })
}

// 结束游戏并结算
export const finishGame = (roomId: string) => {
  return request<Settlement[]>({
    url: '/game/room/finish',
    method: 'POST',
    data: { roomId }
  })
}

// 获取我的房间列表
export const getMyRooms = () => {
  return request<GameRoom[]>({
    url: '/game/rooms',
    method: 'GET'
  })
}

// 切换准备状态
export const toggleReady = (roomId: string) => {
  return request<{ isReady: boolean }>({
    url: `/game/room/${roomId}/ready`,
    method: 'POST'
  })
}

// 开始游戏
export const startGame = (roomId: string) => {
  return request<void>({
    url: `/game/room/${roomId}/start`,
    method: 'POST'
  })
}

// 获取房间结算结果
export const getRoomSettlements = (roomId: string) => {
  return request<Settlement[]>({
    url: `/game/room/${roomId}/settlements`,
    method: 'GET'
  })
}
