import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"

import {api} from "@/api/client"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Alert } from "@/components/ui/alert"
import { Separator } from "@/components/ui/separator"
import { useOrgStore } from "@/stores/orgStore"

type OrgResponse = {
  id: string
  name: string
}

export default function OrgResolverPage() {
  const navigate = useNavigate()
  const { setOrgs, setCurrentOrg, orgs } = useOrgStore()
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    let isMounted = true
    setIsLoading(true)
    setError(null)
    api
      .get("/me/orgs")
      .then((res) => {
        if (!isMounted) return
        const orgList: OrgResponse[] = res.data?.orgs ?? []
        setOrgs(orgList)

        if (orgList.length === 0) {
          navigate("/onboarding", { replace: true })
          return
        }

        if (orgList.length === 1) {
          setCurrentOrg(orgList[0])
          navigate(`/orgs/${orgList[0].id}`, { replace: true })
          return
        }

        setIsLoading(false)
      })
      .catch((err: any) => {
        if (!isMounted) return
        setError(err?.message ?? "Unable to load organizations.")
        setIsLoading(false)
      })
    return () => {
      isMounted = false
    }
  }, [navigate, setCurrentOrg, setOrgs])

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Card className="w-full max-w-lg">
          <CardHeader>
            <CardTitle>Resolving workspace</CardTitle>
            <CardDescription>Please wait...</CardDescription>
          </CardHeader>
          <CardContent>
            <Separator />
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-muted/40 px-4">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Select a workspace</CardTitle>
          <CardDescription>Choose which organization to open.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          {error && <Alert variant="destructive">{error}</Alert>}
          {orgs.map((org) => (
            <Button
              key={org.id}
              variant="outline"
              className="w-full justify-start"
              onClick={() => {
                setCurrentOrg(org)
                navigate(`/orgs/${org.id}/dashboard`)
              }}
            >
              {org.name}
            </Button>
          ))}
        </CardContent>
      </Card>
    </div>
  )
}
