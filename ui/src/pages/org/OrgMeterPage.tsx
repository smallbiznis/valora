import { useEffect, useMemo, useState } from "react"
import { Link, useParams } from "react-router-dom"

import { Plus, MoreHorizontal } from "lucide-react"

import { api } from "@/api/client"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

type Meter = {
  id: string
  code: string
  name: string
  aggregation: string
  unit: string
  active: boolean
  created_at: string
  updated_at: string
}

const formatAggregation = (value: string) => {
  const trimmed = value.trim()
  if (!trimmed) return "-"
  return trimmed.charAt(0).toUpperCase() + trimmed.slice(1).toLowerCase()
}

const formatTimestamp = (value: string) => {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return "-"
  return new Intl.DateTimeFormat("en-US", {
    month: "numeric",
    day: "numeric",
    year: "2-digit",
    hour: "numeric",
    minute: "2-digit",
  }).format(date)
}

export default function OrgMeterPage() {
  const { orgId } = useParams()
  const [meters, setMeters] = useState<Meter[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const rows = useMemo(() => meters, [meters])

  useEffect(() => {
    if (!orgId) {
      setIsLoading(false)
      return
    }

    let isMounted = true
    setIsLoading(true)
    setError(null)

    api
      .get("/meters", { params: { organization_id: orgId } })
      .then((response) => {
        if (!isMounted) return
        setMeters(response.data?.data ?? [])
      })
      .catch((err) => {
        if (!isMounted) return
        setError(err?.message ?? "Unable to load meters.")
      })
      .finally(() => {
        if (!isMounted) return
        setIsLoading(false)
      })

    return () => {
      isMounted = false
    }
  }, [orgId])

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold">Meters</h1>
          <p className="text-muted-foreground text-sm">
            Define usage meters for this organization.
          </p>
        </div>
        <div className="flex items-center gap-2">
          {orgId && (
            <Button size="sm" asChild>
              <Link to={`/orgs/${orgId}/meter/create`}>
                <Plus className="size-4" />
                Create meter
              </Link>
            </Button>
          )}
          <Button size="icon-sm" variant="outline" aria-label="More actions">
            <MoreHorizontal className="size-4" />
          </Button>
        </div>
      </div>

      <div className="border-b">
        <nav className="flex items-center gap-6 text-sm font-medium">
          <span className="border-b-2 border-primary pb-3 text-primary">Meters</span>
          <span className="pb-3 text-muted-foreground">Alerts</span>
          <span className="pb-3 text-muted-foreground">Credits</span>
        </nav>
      </div>

      <div className="rounded-lg border bg-white">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Display name</TableHead>
              <TableHead>Event name</TableHead>
              <TableHead>Aggregation method</TableHead>
              <TableHead>Event ingestion</TableHead>
              <TableHead>Created</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading && (
              <TableRow>
                <TableCell colSpan={5} className="text-muted-foreground">
                  Loading meters...
                </TableCell>
              </TableRow>
            )}
            {error && (
              <TableRow>
                <TableCell colSpan={5} className="text-destructive">
                  {error}
                </TableCell>
              </TableRow>
            )}
            {!isLoading && !error && rows.length === 0 && (
              <TableRow>
                <TableCell colSpan={5} className="text-muted-foreground">
                  No meters yet.
                </TableCell>
              </TableRow>
            )}
            {!isLoading &&
              !error &&
              rows.map((meter) => (
                <TableRow key={meter.id}>
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <Link
                        className="font-medium text-foreground hover:underline"
                        to={`/orgs/${orgId}/meter/${meter.id}`}
                      >
                        {meter.name}
                      </Link>
                      <Badge
                        variant={meter.active ? "secondary" : "outline"}
                        className={
                          meter.active
                            ? "border-green-200 bg-green-50 text-green-700"
                            : "text-muted-foreground"
                        }
                      >
                        {meter.active ? "Active" : "Deactivated"}
                      </Badge>
                    </div>
                  </TableCell>
                  <TableCell className="text-muted-foreground">{meter.code}</TableCell>
                  <TableCell className="text-muted-foreground">
                    {formatAggregation(meter.aggregation)}
                  </TableCell>
                  <TableCell className="text-muted-foreground">Raw</TableCell>
                  <TableCell className="text-muted-foreground">
                    {formatTimestamp(meter.created_at)}
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
