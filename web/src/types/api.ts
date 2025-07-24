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

export interface ApiError {
  message: string;
}
