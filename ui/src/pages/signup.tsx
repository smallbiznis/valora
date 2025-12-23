import { useState } from "react"
import { useNavigate } from "react-router-dom"

import { api } from "@/api/client"
import { SignupForm } from "@/components/signup-form"
import { useAuthStore } from "@/stores/authStore"

export default function SignupPage() {
  const navigate = useNavigate()
  const signup = useAuthStore((s) => s.signup)

  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (payload: { username: string; password: string }) => {
    setError(null)
    setIsLoading(true)
    try {
      await signup(payload)
      await api.get("/me")
      navigate("/orgs", { replace: true })
    } catch (err: any) {
      const status = err?.response?.status
      if (status === 404) {
        navigate("/login", { replace: true })
        return
      }
      setError(err?.message ?? "Unable to create account.")
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-muted/40 px-4 py-12">
      <SignupForm onSubmit={handleSubmit} isLoading={isLoading} error={error} className="w-full max-w-md" />
    </div>
  )
}
