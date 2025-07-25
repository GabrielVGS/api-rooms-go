import axios from 'axios';
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  Room,
  Reservation,
  Note,
  CreateNoteRequest,
  UpdateNoteRequest,
  User,
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

export const usersApi = {
  getUsers: async (): Promise<User[]> => {
    try {
      console.log("Getting all users")
      const response = await api.get("/users")
      return response.data.users
    } catch (error) {
      console.error("Error getting users");
      throw error;
    }
  }
}

export const roomApi = {
  getRooms: async (): Promise<Room[]> => {
    try {
      console.log('Fetching all rooms');
      const response = await api.get('/rooms');
      console.log(`Successfully fetched ${response.data.length} rooms`);
      return response.data ?? [];
    } catch (error) {
      console.error('Failed to fetch rooms:', error);
      throw error;
    }
  },

  getRoomById: async (id: number): Promise<Room> => {
    try {
      console.log(`Fetching room ID ${id}`);
      const response = await api.get(`/rooms/${id}`);
      console.log(`Successfully fetched room ID ${id}:`, response.data.name);
      return response.data;
    } catch (error) {
      console.error(`Failed to fetch room ID ${id}:`, error);
      throw error;
    }
  },

  getUserRooms: async (): Promise<Room[]> => {
    try {
      console.log('Fetching user rooms');
      const response = await api.get('/rooms/my-rooms');
      console.log(`Successfully fetched ${response.data.length} user rooms`);
      return response.data ?? [];
    } catch (error) {
      console.error('Failed to fetch user rooms:', error);
      throw error;
    }
  },

  joinRoom: async (roomId: number): Promise<void> => {
    try {
      console.log(`Joining room ID ${roomId}`);
      await api.post(`/rooms/${roomId}/join`);
      console.log(`Successfully joined room ID ${roomId}`);
    } catch (error) {
      console.error(`Failed to join room ID ${roomId}:`, error);
      throw error;
    }
  },

  leaveRoom: async (roomId: number): Promise<void> => {
    try {
      console.log(`Leaving room ID ${roomId}`);
      await api.delete(`/rooms/${roomId}/leave`);
      console.log(`Successfully left room ID ${roomId}`);
    } catch (error) {
      console.error(`Failed to leave room ID ${roomId}:`, error);
      throw error;
    }
  },

  createRoom: async (room: Omit<Room, 'id' | 'created_by' | 'members' | 'notes' | 'created_at' | 'updated_at'>): Promise<Room> => {
    try {
      console.log('Creating new room:', room.name);
      const createRequest = {
        name: room.name,
        description: room.description || '',
        subject: room.subject,
        capacity: room.capacity,
      };
      const response = await api.post('/rooms', createRequest);
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
      const updateRequest = {
        name: room.name || '',
        description: room.description || '',
        subject: room.subject || '',
        capacity: room.capacity || 0,
      };
      const response = await api.put(`/rooms/${id}`, updateRequest);
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

export const notesApi = {
  getNotesByRoom: async (roomId: number): Promise<Note[]> => {
    try {
      console.log(`Fetching notes for room ID ${roomId}`);
      const response = await api.get(`/notes/room/${roomId}`);
      console.log(`Successfully fetched ${response.data.length} notes for room ID ${roomId}`);
      return response.data ?? [];
    } catch (error) {
      console.error(`Failed to fetch notes for room ID ${roomId}:`, error);
      throw error;
    }
  },

  getUserNotes: async (): Promise<Note[]> => {
    try {
      console.log('Fetching user notes');
      const response = await api.get('/notes/my-notes');
      console.log(`Successfully fetched ${response.data.length} user notes`);
      return response.data ?? [];
    } catch (error) {
      console.error('Failed to fetch user notes:', error);
      throw error;
    }
  },

  getNoteById: async (id: number): Promise<Note> => {
    try {
      console.log(`Fetching note ID ${id}`);
      const response = await api.get(`/notes/${id}`);
      console.log(`Successfully fetched note ID ${id}:`, response.data.title);
      return response.data;
    } catch (error) {
      console.error(`Failed to fetch note ID ${id}:`, error);
      throw error;
    }
  },

  getAllNotes: async (): Promise<Note[]> => {
    try {

      console.log("Getting all notes")
      const response = await api.get("/notes")
      console.log("Succesfully fetched all notes")
      return response.data
    } catch (error) {
      console.error(`Failed to fetch all notes`, error);
      throw error;
    }
  },

  createNote: async (note: CreateNoteRequest): Promise<Note> => {
    try {
      console.log('Creating new note:', note.title, 'for room ID', note.room_id);
      const response = await api.post('/notes', note);
      console.log('Successfully created note ID:', response.data.id);
      return response.data;
    } catch (error) {
      console.error('Failed to create note:', note, error);
      throw error;
    }
  },

  updateNote: async (id: number, note: UpdateNoteRequest): Promise<Note> => {
    try {
      console.log(`Updating note ID ${id}:`, note);
      const response = await api.put(`/notes/${id}`, note);
      console.log(`Successfully updated note ID ${id}`);
      return response.data;
    } catch (error) {
      console.error(`Failed to update note ID ${id}:`, error);
      throw error;
    }
  },

  deleteNote: async (id: number): Promise<void> => {
    try {
      console.log(`Deleting note ID ${id}`);
      await api.delete(`/notes/${id}`);
      console.log(`Successfully deleted note ID ${id}`);
    } catch (error) {
      console.error(`Failed to delete note ID ${id}:`, error);
      throw error;
    }
  },
};

export { api };