import axios from 'axios';

// 1. Define the Backend Address (The "Post Office")
const API = axios.create({
  baseURL: 'https://careersync-backend-ntwb.onrender.com',
});

// 2. The Automatic Stamper (Interceptor)
// Before any request is sent, this code runs.
API.interceptors.request.use((req) => {
  // Check if we have a token saved in the browser's storage
  const token = localStorage.getItem('token');
  
  // If yes, attach it to the header (Authorization: Bearer <token>)
  if (token) {
    req.headers.Authorization = `Bearer ${token}`;
  }
  return req;
});

export default API;