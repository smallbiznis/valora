import { create } from "zustand"
import { persist } from "zustand/middleware"

import { auth } from "@/api/client"

type User = {
  id: string
  username: string
  email?: string
}

type AuthState = {
  user: User | null
  isAuthenticated: boolean
  login: (payload: { username: string; password: string }) => Promise<void>
  signup: (payload: { username: string; password: string }) => Promise<void>
  logout: () => Promise<void>
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      login: async (payload) => {
        try {
          const res = await auth.post("/login", payload)
          const user = res.data?.metadata
          set({ user, isAuthenticated: true })
        } catch (error) {
          console.error("error: ", error)
        }
      },
      signup: async (payload) => {
        const res = await auth.post("/signup", payload)
        const user = res.data?.metadata
        set({ user, isAuthenticated: true })
      },
      logout: async () => {
        try {
          await auth.post("/logout")
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
