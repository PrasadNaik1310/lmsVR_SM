import axios from "axios";

console.log(import.meta.env.VITE_API_BASE_URL)
export const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});