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

api.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`, {
      baseURL: config.baseURL,
      headers: config.headers,
      data: config.data
    });
    
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => {
    console.log(`API Response: ${response.status} ${response.config.method?.toUpperCase()} ${response.config.url}`, {
      status: response.status,
      statusText: response.statusText,
      data: response.data
    });
    return response;
  },
  (error) => {
    console.error('API Response Error:', {
      status: error.response?.status,
      statusText: error.response?.statusText,
      data: error.response?.data,
      message: error.message,
      config: {
        method: error.config?.method,
        url: error.config?.url
      }
    });
    return Promise.reject(error);
  }
);

export const authApi = {
  login: async (credentials: LoginRequest): Promise<AuthResponse> => {
    try {
      console.log('Attempting login for user:', credentials.email);
      const response = await api.post('/auth/login', credentials);
      console.log('Login successful for user:', credentials.email);
      return response.data;
    } catch (error) {
      console.error('Login failed for user:', credentials.email, error);
      throw error;
    }
  },

  register: async (userData: RegisterRequest): Promise<AuthResponse> => {
    try {
      console.log('Attempting registration for user:', userData.email);
      const response = await api.post('/auth/register', userData);
      console.log('Registration successful for user:', userData.email);
      return response.data;
    } catch (error) {
      console.error('Registration failed for user:', userData.email, error);
      throw error;
    }
  },
};

export const roomApi = {
  getRooms: async (): Promise<Room[]> => {
    try {
      console.log('Fetching all rooms');
      const response = await api.get('/rooms');
      console.log(`Successfully fetched ${response.data.length} rooms`);
      return response.data?? [];
    } catch (error) {
      console.error('Failed to fetch rooms:', error);
      throw error;
    }
  },

  createRoom: async (room: Omit<Room, 'id'>): Promise<Room> => {
    try {
      console.log('Creating new room:', room.name);
      const response = await api.post('/rooms', room);
      console.log('Successfully created room:', response.data.name);
      return response.data;
    } catch (error) {
      console.error('Failed to create room:', room.name, error);
      throw error;
    }
  },

  updateRoom: async (id: number, room: Partial<Room>): Promise<Room> => {
    try {
      console.log(`Updating room ID ${id}:`, room);
      const response = await api.put(`/rooms/${id}`, room);
      console.log(`Successfully updated room ID ${id}`);
      return response.data;
    } catch (error) {
      console.error(`Failed to update room ID ${id}:`, error);
      throw error;
    }
  },

  deleteRoom: async (id: number): Promise<void> => {
    try {
      console.log(`Deleting room ID ${id}`);
      await api.delete(`/rooms/${id}`);
      console.log(`Successfully deleted room ID ${id}`);
    } catch (error) {
      console.error(`Failed to delete room ID ${id}:`, error);
      throw error;
    }
  },
};

export const reservationApi = {
  getReservations: async (): Promise<Reservation[]> => {
    try {
      console.log('Fetching all reservations');
      const response = await api.get('/reservations');
      console.log(`Successfully fetched ${response.data.length} reservations`);
      return response.data;
    } catch (error) {
      console.error('Failed to fetch reservations:', error);
      throw error;
    }
  },

  createReservation: async (reservation: Omit<Reservation, 'id' | 'user' | 'room'>): Promise<Reservation> => {
    try {
      console.log('Creating new reservation:', {
        roomId: reservation.room_id,
        startTime: reservation.start_time,
        endTime: reservation.end_time
      });
      const response = await api.post('/reservations', reservation);
      console.log('Successfully created reservation ID:', response.data.id);
      return response.data;
    } catch (error) {
      console.error('Failed to create reservation:', reservation, error);
      throw error;
    }
  },

  updateReservation: async (id: number, reservation: Partial<Reservation>): Promise<Reservation> => {
    try {
      console.log(`Updating reservation ID ${id}:`, reservation);
      const response = await api.put(`/reservations/${id}`, reservation);
      console.log(`Successfully updated reservation ID ${id}`);
      return response.data;
    } catch (error) {
      console.error(`Failed to update reservation ID ${id}:`, error);
      throw error;
    }
  },

  deleteReservation: async (id: number): Promise<void> => {
    try {
      console.log(`Deleting reservation ID ${id}`);
      await api.delete(`/reservations/${id}`);
      console.log(`Successfully deleted reservation ID ${id}`);
    } catch (error) {
      console.error(`Failed to delete reservation ID ${id}:`, error);
      throw error;
    }
  },
};

export { api };