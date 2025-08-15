import axios from "axios";

const API_BASE = "http://localhost:8080";

const api = axios.create({
  baseURL: API_BASE,
});

export const setToken = (token: string) => {
  api.defaults.headers.common["Authorization"] = `Bearer ${token}`;
};

export const registerUser = (data: { email: string; password: string; name: string }) =>
  api.post("/auth/register", data);

export const loginUser = (data: { email: string; password: string }) =>
  api.post("/auth/login", data);

export const getTasks = () => api.get("/tasks");
export const createTask = (data: { title: string; priority?: number }) => api.post("/tasks", data);
export const updateTask = (id: number, data: never) => api.put(`/tasks/${id}`, data);
export const deleteTask = (id: number) => api.delete(`/tasks/${id}`);

export default api;
