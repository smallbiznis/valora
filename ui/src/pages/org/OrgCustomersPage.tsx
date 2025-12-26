import { useCallback, useEffect, useMemo, useState } from "react"
import { useParams } from "react-router-dom"
import { IconPlus } from "@tabler/icons-react"

import { api } from "@/api/client"
import { Alert } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Separator } from "@/components/ui/separator"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

type Customer = {
  id: string | number
  name: string
  email: string
  created_at?: string
}

type CustomerListResponse = {
  customers: Customer[]
  next_page_token?: string
  previous_page_token?: string
  has_more?: boolean
}

const PAGE_SIZE = 25

const formatDateTime = (value?: string) => {
  if (!value) return "-"
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return "-"
  return new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "2-digit",
    hour: "numeric",
    minute: "2-digit",
  }).format(date)
}

export default function OrgCustomersPage() {
  const { orgId } = useParams()
  const [customers, setCustomers] = useState<Customer[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isLoadingMore, setIsLoadingMore] = useState(false)
  const [listError, setListError] = useState<string | null>(null)
  const [createError, setCreateError] = useState<string | null>(null)
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [language, setLanguage] = useState("en-US")
  const [isCreating, setIsCreating] = useState(false)
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [nextPageToken, setNextPageToken] = useState<string | null>(null)
  const [hasMore, setHasMore] = useState(false)

  const loadCustomers = useCallback(async (options?: { pageToken?: string; append?: boolean }) => {
    if (!orgId) {
      setIsLoading(false)
      return
    }

    const { pageToken, append } = options ?? {}
    if (append) {
      setIsLoadingMore(true)
    } else {
      setIsLoading(true)
    }
    setListError(null)

    try {
      const res = await api.get("/customers", {
        params: {
          page_token: pageToken,
          page_size: PAGE_SIZE,
        },
      })
      const payload: CustomerListResponse = res.data?.data ?? { customers: [] }
      const list = Array.isArray(payload.customers) ? payload.customers : []
      if (append) {
        setCustomers((prev) => [...prev, ...list])
      } else {
        setCustomers(list)
      }
      setNextPageToken(payload.next_page_token ?? null)
      setHasMore(Boolean(payload.has_more))
    } catch (err: any) {
      setListError(err?.message ?? "Unable to load customers.")
      if (!append) {
        setCustomers([])
      }
      setNextPageToken(null)
      setHasMore(false)
    } finally {
      if (append) {
        setIsLoadingMore(false)
      } else {
        setIsLoading(false)
      }
    }
  }, [orgId])

  useEffect(() => {
    void loadCustomers()
  }, [loadCustomers])

  const handleCreate = async (event: React.FormEvent) => {
    event.preventDefault()
    if (!orgId) return
    setIsCreating(true)
    setCreateError(null)
    try {
      const res = await api.post("/customers", {
        organization_id: orgId,
        name,
        email,
      })
      const created = res.data?.data
      if (created) {
        setCustomers((prev) => [created, ...prev])
      } else {
        await loadCustomers()
      }
      setName("")
      setEmail("")
      setLanguage("en-US")
      setIsDialogOpen(false)
    } catch (err: any) {
      setCreateError(err?.message ?? "Unable to create customer.")
    } finally {
      setIsCreating(false)
    }
  }

  const countLabel = useMemo(() => {
    const total = customers.length
    return `${total} item${total === 1 ? "" : "s"}`
  }, [customers.length])

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold">Customers</h1>
          <p className="text-text-muted text-sm">
            Manage customers, payment methods, and billing contact data.
          </p>
        </div>
        <Dialog
          open={isDialogOpen}
          onOpenChange={(open) => {
            setIsDialogOpen(open)
            if (!open) {
              setCreateError(null)
              setName("")
              setEmail("")
              setLanguage("en-US")
            }
          }}
        >
          <DialogTrigger asChild>
            <Button size="sm">
              <IconPlus />
              Add customer
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-xl">
            <DialogHeader>
              <DialogTitle>Create customer</DialogTitle>
              <DialogDescription>
                Add a customer profile for billing and invoicing.
              </DialogDescription>
            </DialogHeader>
            <form className="space-y-4" onSubmit={handleCreate}>
              {createError && <Alert variant="destructive">{createError}</Alert>}
              <div className="space-y-2">
                <Label htmlFor="customer-name">Name</Label>
                <Input
                  id="customer-name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Customer display name"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="customer-email">Email</Label>
                <Input
                  id="customer-email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="billing@company.com"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="customer-language">Language</Label>
                <Select value={language} onValueChange={setLanguage}>
                  <SelectTrigger id="customer-language">
                    <SelectValue placeholder="Select a language" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="en-US">English (United States)</SelectItem>
                    <SelectItem value="en-GB">English (United Kingdom)</SelectItem>
                    <SelectItem value="id-ID">Bahasa Indonesia</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="px-0 text-accent-primary"
              >
                More options
              </Button>
              <Separator />
              <DialogFooter>
                <DialogClose asChild>
                  <Button type="button" variant="outline">
                    Cancel
                  </Button>
                </DialogClose>
                <Button type="submit" disabled={isCreating || !orgId}>
                  {isCreating ? "Saving..." : "Add customer"}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid gap-3 md:grid-cols-[2fr_1fr]">
        <Input placeholder="All customers" aria-label="Search customers" />
        <Input placeholder="Remaining balances" aria-label="Filter by balance" />
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap items-center gap-2">
          <Button variant="outline" size="sm">
            <IconPlus />
            Email
          </Button>
          <Button variant="outline" size="sm">
            <IconPlus />
            Name
          </Button>
          <Button variant="outline" size="sm">
            <IconPlus />
            Created date
          </Button>
          <Button variant="outline" size="sm">
            <IconPlus />
            Type
          </Button>
          <Button variant="outline" size="sm">
            <IconPlus />
            More filters
          </Button>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Button variant="outline" size="sm">
            Copy
          </Button>
          <Button variant="outline" size="sm">
            Export
          </Button>
          <Button variant="outline" size="sm">
            Analyze
          </Button>
          <Button variant="outline" size="sm">
            Edit columns
          </Button>
        </div>
      </div>

      {listError && <Alert variant="destructive">{listError}</Alert>}

      <div className="rounded-lg border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-10">
                <Checkbox aria-label="Select all customers" />
              </TableHead>
              <TableHead>Customer</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Primary payment method</TableHead>
              <TableHead>Created</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading && (
              <TableRow>
                <TableCell colSpan={6} className="text-text-muted">
                  Loading customers...
                </TableCell>
              </TableRow>
            )}
            {!isLoading && customers.length === 0 && !listError && (
              <TableRow>
                <TableCell colSpan={6} className="text-text-muted">
                  No customers yet.
                </TableCell>
              </TableRow>
            )}
            {!isLoading &&
              customers.length > 0 &&
              customers.map((customer) => (
                <TableRow key={customer.id}>
                  <TableCell>
                    <Checkbox aria-label={`Select ${customer.name}`} />
                  </TableCell>
                  <TableCell className="font-medium">{customer.name}</TableCell>
                  <TableCell className="text-text-muted">
                    {customer.email}
                  </TableCell>
                  <TableCell className="text-text-muted">-</TableCell>
                  <TableCell className="text-text-muted">
                    {formatDateTime(customer.created_at)}
                  </TableCell>
                  <TableCell className="text-right">
                    <Button variant="ghost" size="icon-sm" aria-label="Open actions">
                      ...
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3 text-sm">
        <span className="text-text-muted">{countLabel}</span>
        <Button
          variant="outline"
          size="sm"
          disabled={!hasMore || isLoadingMore}
          onClick={() => {
            if (nextPageToken) {
              void loadCustomers({ pageToken: nextPageToken, append: true })
            }
          }}
        >
          {isLoadingMore ? "Loading..." : "Load more"}
        </Button>
      </div>
    </div>
  )
}
