export interface User {
  id: number;
  name: string;
  email: string;
}

export interface Room {
  id: number;
  name: string;
  capacity: number;
  subject: string;
  description?: string;
  created_by: number;
  members?: RoomMember[];
  notes?: Note[];
  created_at: string;
  updated_at: string;
}

export interface Reservation {
  id: number;
  user_id: number;
  room_id: number;
  start_time: string;
  end_time: string;
  user?: User;
  room?: Room;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface RoomMember {
  user_id: number;
  user_name: string;
  user_email: string;
  role: string;
  joined_at: string;
}

export interface Note {
  id: number;
  user_id: number;
  room_id: number;
  title: string;
  content: string;
  user_name?: string;
  user_email?: string;
  room_name?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateNoteRequest {
  room_id: number;
  title: string;
  content: string;
}

export interface UpdateNoteRequest {
  title: string;
  content: string;
}

export interface ApiError {
  message: string;
}
