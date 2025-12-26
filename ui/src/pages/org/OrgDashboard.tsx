import { Navigate } from "react-router-dom"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { useOrgStore } from "@/stores/orgStore"

export default function OrgDashboard() {
  const org = useOrgStore((s) => s.currentOrg)

  if (!org) {
    return <Navigate to="/orgs" replace />
  }

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-2xl font-semibold">Welcome</h1>
        <p className="text-text-muted">Working in {org.name}</p>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Get started</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <p className="text-sm text-text-muted">
            Use the sidebar to manage products, meters, customers, subscriptions, invoices, and settings. Pricing lives inside each product.
          </p>
          <Separator />
        </CardContent>
      </Card>
    </div>
  )
}
