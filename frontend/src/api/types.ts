export type User = {
  id: number
  username: string
  role: string
}

export type Room = {
  id: number
  block: string
  name: string
  capacity: number
  status: string
  allowed_purposes?: string
  notes?: string
}

export type Booking = {
  id: number
  room_id: number
  room?: Room
  purpose: string
  participants: number
  start_time: string
  end_time: string
  status: string
}
