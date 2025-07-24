import axios from 'axios';
import type { 
  LoginRequest, 
  RegisterRequest, 
  AuthResponse, 
  Room, 
  Reservation 
} from '../types/api';

const API_BASE_URL = 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authApi = {
  login: async (credentials: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post('/auth/login', credentials);
    return response.data;
  },

  register: async (userData: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post('/auth/register', userData);
    return response.data;
  },
};

export const roomApi = {
  getRooms: async (): Promise<Room[]> => {
    const response = await api.get('/rooms');
    return response.data;
  },

  createRoom: async (room: Omit<Room, 'id'>): Promise<Room> => {
    const response = await api.post('/rooms', room);
    return response.data;
  },

  updateRoom: async (id: number, room: Partial<Room>): Promise<Room> => {
    const response = await api.put(`/rooms/${id}`, room);
    return response.data;
  },

  deleteRoom: async (id: number): Promise<void> => {
    await api.delete(`/rooms/${id}`);
  },
};

export const reservationApi = {
  getReservations: async (): Promise<Reservation[]> => {
    const response = await api.get('/reservations');
    return response.data;
  },

  createReservation: async (reservation: Omit<Reservation, 'id' | 'user' | 'room'>): Promise<Reservation> => {
    const response = await api.post('/reservations', reservation);
    return response.data;
  },

  updateReservation: async (id: number, reservation: Partial<Reservation>): Promise<Reservation> => {
    const response = await api.put(`/reservations/${id}`, reservation);
    return response.data;
  },

  deleteReservation: async (id: number): Promise<void> => {
    await api.delete(`/reservations/${id}`);
  },
};

export { api };