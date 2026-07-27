import { http } from "../services/http.js";

export async function login(payload) {
  const response = await http.post("/auth/login", payload);
  return response.data;
}