import { create } from "zustand"
import { persist } from "zustand/middleware"

import { auth, authLocal } from "@/api/client"

type User = {
  id: string
  username: string
  email?: string
  displayName?: string
  provider?: string
  externalId?: string
}

type AuthState = {
  user: User | null
  isAuthenticated: boolean
  login: (payload: { email: string; password: string }) => Promise<void>
  signup: (payload: { email: string; password: string; displayName?: string; orgName?: string }) => Promise<void>
  logout: () => Promise<void>
}

const buildUser = (payload: any): User | null => {
  const metadata = payload?.metadata ?? payload
  if (!metadata) {
    return null
  }
  const id = metadata.user_id ?? metadata.id
  if (!id) {
    return null
  }
  const email = metadata.email
  const displayName = metadata.display_name ?? metadata.username ?? email
  const username = displayName ?? email ?? String(id)

  return {
    id: String(id),
    username,
    email,
    displayName,
    provider: metadata.provider ?? metadata.auth_provider,
    externalId: metadata.external_id,
  }
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      login: async (payload) => {
        try {
          const email = payload.email.trim().toLowerCase()
          const res = await authLocal.post("/login", {
            email,
            password: payload.password,
          })
          const user = buildUser(res.data)
          if (!user) {
            throw new Error("invalid_login_response")
          }
          set({ user, isAuthenticated: true })
        } catch (err) {
          set({ user: null, isAuthenticated: false })
          throw err
        }
      },
      signup: async (payload) => {
        try {
          const displayName = payload.displayName?.trim()
          const email = payload.email.trim().toLowerCase()
          const username =
            displayName ||
            email.split("@")[0] ||
            "user"
          const body: Record<string, string> = {
            username,
            email,
            password: payload.password,
          }
          if (payload.orgName?.trim()) {
            body.org_name = payload.orgName.trim()
          }
          const res = await auth.post("/signup", body)
          const user = buildUser(res.data)
          if (!user) {
            throw new Error("invalid_signup_response")
          }
          set({ user, isAuthenticated: true })
        } catch (err) {
          set({ user: null, isAuthenticated: false })
          throw err
        }
      },
      logout: async () => {
        try {
          await authLocal.post("/logout")
        } catch (err) {
          console.warn("logout failed", err)
        }
        set({ user: null, isAuthenticated: false })
      },
    }),
    {
      name: "valora-auth",
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
