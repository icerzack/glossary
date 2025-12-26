import axios from 'axios';
import type {
  Term,
  Relationship,
  Graph,
  CreateTermRequest,
  UpdateTermRequest,
  CreateRelationshipRequest,
} from '../types';

const API_BASE_URL = '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Terms API
export const termsApi = {
  getAll: async (search?: string, category?: string): Promise<Term[]> => {
    const params = new URLSearchParams();
    if (search) params.append('search', search);
    if (category) params.append('category', category);
    const response = await api.get<Term[]>(`/terms?${params.toString()}`);
    return response.data;
  },

  getById: async (id: number): Promise<Term> => {
    const response = await api.get<Term>(`/terms/${id}`);
    return response.data;
  },

  create: async (data: CreateTermRequest): Promise<Term> => {
    const response = await api.post<Term>('/terms', data);
    return response.data;
  },

  update: async (id: number, data: UpdateTermRequest): Promise<Term> => {
    const response = await api.put<Term>(`/terms/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/terms/${id}`);
  },
};

// Relationships API
export const relationshipsApi = {
  getAll: async (): Promise<Relationship[]> => {
    const response = await api.get<Relationship[]>('/relationships');
    return response.data;
  },

  create: async (data: CreateRelationshipRequest): Promise<Relationship> => {
    const response = await api.post<Relationship>('/relationships', data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/relationships/${id}`);
  },
};

// Graph API
export const graphApi = {
  get: async (): Promise<Graph> => {
    const response = await api.get<Graph>('/graph');
    return response.data;
  },
};

