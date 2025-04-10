import axios from 'axios';


const API = 'http://localhost:8080/api';
export const api = axios.create({ baseURL: API });

export const loginUser = (username, password) =>
  api.post('/login', { username, password });

export const registerUser = (username, password) =>
  api.post('/register', { username, password });

export const saveJob = (filename, token) =>
  api.post('/api/save-job', { filename, status: 'Pending...' }, {
    headers: { Authorization: `Bearer ${token}` }
  });

export const getJobHistory = (token) =>
  api.get('/api/job-history', {
    headers: { Authorization: `Bearer ${token}` }
  });
