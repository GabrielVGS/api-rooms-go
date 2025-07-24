export interface User {
  id: number;
  name: string;
  email: string;
}

export interface Room {
  id: number;
  name: string;
  capacity: number;
  description?: string;
}

export interface Reservation {
  id: number;
  userId: number;
  roomId: number;
  startTime: string;
  endTime: string;
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