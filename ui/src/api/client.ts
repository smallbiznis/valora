import axios from "axios";

export const auth = axios.create({
  baseURL: "/auth",
  withCredentials: true,
});

export const api = axios.create({
  baseURL: "/api",
  withCredentials: true,
});