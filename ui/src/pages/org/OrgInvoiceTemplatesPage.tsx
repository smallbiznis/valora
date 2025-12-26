import { useMemo, useState } from "react"
import { Link, useParams } from "react-router-dom"
import { IconPlus } from "@tabler/icons-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyTitle,
} from "@/components/ui/empty"
import { Input } from "@/components/ui/input"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  findInvoiceTemplate,
  invoiceTemplates,
  type InvoiceTemplate,
} from "@/pages/org/invoice-template-data"

const formatDate = (value?: string) => {
  if (!value) return "-"
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return "-"
  return new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "2-digit",
    year: "numeric",
  }).format(date)
}

const formatStatus = (status: InvoiceTemplate["status"]) => {
  switch (status) {
    case "ACTIVE":
      return "Active"
    case "DRAFT":
      return "Draft"
    default:
      return status
  }
}

const statusVariant = (status: InvoiceTemplate["status"]) => {
  switch (status) {
    case "ACTIVE":
      return "secondary"
    case "DRAFT":
      return "outline"
    default:
      return "secondary"
  }
}

export default function OrgInvoiceTemplatesPage() {
  const { orgId } = useParams()
  const [searchQuery, setSearchQuery] = useState("")
  const orgBasePath = orgId ? `/orgs/${orgId}` : "/orgs"

  const filteredTemplates = useMemo(() => {
    const query = searchQuery.trim().toLowerCase()
    if (!query) return invoiceTemplates
    return invoiceTemplates.filter((template) => {
      const memo = template.memo?.toLowerCase() ?? ""
      const footer = template.footer?.toLowerCase() ?? ""
      return (
        template.name.toLowerCase().includes(query) ||
        memo.includes(query) ||
        footer.includes(query)
      )
    })
  }, [searchQuery])

  const countLabel = useMemo(() => {
    const total = filteredTemplates.length
    return `${total} template${total === 1 ? "" : "s"}`
  }, [filteredTemplates.length])

  const defaultTemplate = findInvoiceTemplate("default")

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold">Invoice templates</h1>
          <p className="text-text-muted text-sm">
            Set reusable layouts for invoices, notes, and line item grouping.
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          {orgId && (
            <Button asChild size="sm">
              <Link to={`${orgBasePath}/invoice-templates/create`}>
                <IconPlus />
                Create template
              </Link>
            </Button>
          )}
          <Button asChild variant="outline" size="sm">
            <Link to={`${orgBasePath}/invoices`}>View invoices</Link>
          </Button>
        </div>
      </div>

      {defaultTemplate && (
        <div className="rounded-lg border bg-bg-subtle/30 p-4 text-sm">
          <div className="flex flex-wrap items-center gap-2">
            <Badge variant="secondary">Default template</Badge>
            <span className="font-medium">{defaultTemplate.name}</span>
          </div>
          <div className="text-text-muted mt-1">
            Used for new subscriptions and invoices unless another template is selected.
          </div>
        </div>
      )}

      <div className="flex flex-wrap items-center gap-3">
        <Input
          className="w-full max-w-md"
          placeholder="Search templates"
          value={searchQuery}
          onChange={(event) => setSearchQuery(event.target.value)}
        />
        <div className="text-text-muted text-sm">{countLabel}</div>
      </div>

      {filteredTemplates.length === 0 ? (
        <Empty>
          <EmptyHeader>
            <EmptyTitle>No templates yet</EmptyTitle>
            <EmptyDescription>
              Create an invoice template to standardize memo, footer, and line item grouping.
            </EmptyDescription>
          </EmptyHeader>
          <EmptyContent>
            {orgId && (
              <Button asChild>
                <Link to={`${orgBasePath}/invoice-templates/create`}>
                  Create template
                </Link>
              </Button>
            )}
          </EmptyContent>
        </Empty>
      ) : (
        <div className="rounded-lg border">
          <Table className="min-w-[760px]">
            <TableHeader>
              <TableRow>
                <TableHead>Template</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Custom fields</TableHead>
                <TableHead>Grouping</TableHead>
                <TableHead>Updated</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredTemplates.map((template) => (
                <TableRow key={template.id}>
                  <TableCell>
                    <div className="flex flex-col gap-1">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="font-medium">{template.name}</span>
                        {template.isDefault && (
                          <Badge variant="secondary">Default</Badge>
                        )}
                      </div>
                      <span className="text-text-muted text-sm">
                        {template.memo}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant={statusVariant(template.status)}>
                      {formatStatus(template.status)}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    {template.customFields?.length ?? 0}
                  </TableCell>
                  <TableCell>{template.lineItemGrouping ?? "-"}</TableCell>
                  <TableCell>{formatDate(template.updatedAt)}</TableCell>
                  <TableCell className="text-right">
                    {orgId ? (
                      <Button asChild variant="ghost" size="sm">
                        <Link to={`${orgBasePath}/invoice-templates/${template.id}`}>
                          Edit
                        </Link>
                      </Button>
                    ) : (
                      <Button variant="ghost" size="sm" disabled>
                        Edit
                      </Button>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  )
}
