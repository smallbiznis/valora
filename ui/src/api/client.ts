import axios from "axios";

import { useOrgStore } from "@/stores/orgStore";

export const auth = axios.create({
  baseURL: "/auth",
  withCredentials: true,
});

export const authLocal = axios.create({
  baseURL: "/internal/auth/local",
  withCredentials: true,
});

export const api = axios.create({
  baseURL: "/api",
  withCredentials: true,
});

api.interceptors.request.use((config) => {
  const orgId = useOrgStore.getState().currentOrg?.id;
  if (orgId) {
    config.headers = config.headers ?? {};
    config.headers["X-Org-Id"] = orgId;
  }
  return config;
});
