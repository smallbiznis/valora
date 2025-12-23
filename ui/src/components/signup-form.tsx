import { useState, type ComponentPropsWithoutRef } from "react"
import { Link } from "react-router-dom"

import { Alert } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"

type SignupFormProps = {
  onSubmit: (payload: { username: string; password: string }) => Promise<void>
  isLoading?: boolean
  error?: string | null
} & Omit<ComponentPropsWithoutRef<typeof Card>, "onSubmit">

export function SignupForm({
  onSubmit,
  isLoading,
  error,
  ...props
}: SignupFormProps) {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    await onSubmit({ username, password })
  }

  return (
    <Card {...props}>
      <CardHeader>
        <CardTitle>Create an account</CardTitle>
        <CardDescription>
          Enter your information below to create your account
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <FieldGroup>
            {error && <Alert variant="destructive">{error}</Alert>}
            <Field>
              <FieldLabel htmlFor="username">Username</FieldLabel>
              <Input
                id="username"
                type="text"
                placeholder="admin"
                required
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
              <FieldDescription>
                Use a unique username for your account.
              </FieldDescription>
            </Field>
            <Field>
              <FieldLabel htmlFor="password">Password</FieldLabel>
              <Input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              <FieldDescription>Must be at least 8 characters long.</FieldDescription>
            </Field>
            <FieldGroup>
              <Field>
                <Button type="submit" disabled={isLoading}>
                  {isLoading ? "Creating..." : "Create Account"}
                </Button>
                <FieldDescription className="px-6 text-center">
                  Already have an account?{" "}
                  <Link to="/login" className="text-primary hover:underline">
                    Sign in
                  </Link>
                </FieldDescription>
              </Field>
            </FieldGroup>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  )
}
