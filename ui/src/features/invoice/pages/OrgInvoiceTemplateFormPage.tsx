import { useEffect, useMemo, useState } from "react"
import { Link, useParams } from "react-router-dom"
import { IconPlus } from "@tabler/icons-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Separator } from "@/components/ui/separator"
// import { Skeleton } from "@/components/ui/skeleton"
import { Textarea } from "@/components/ui/textarea"
import { findInvoiceTemplate, type InvoiceTemplate } from "@/features/invoice/data/invoice-template-data"
import { useOrgStore } from "@/stores/orgStore"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

type CustomField = {
  label: string
  value: string
}

const groupingOptions = [
  "None",
  "Service category",
  "Subscription",
  "Meter",
  "Product",
]

const previewLineItems = [
  { description: "Starter plan", quantity: 1, unitPrice: 0, amount: 0 },
  { description: "Usage overage", quantity: 0, unitPrice: 0, amount: 0 },
]

const formatCurrency = (amount: number, currency = "USD") => {
  try {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency,
    }).format(amount / 100)
  } catch {
    return `${(amount / 100).toFixed(2)} ${currency}`
  }
}

const buildState = (template: InvoiceTemplate | null) => ({
  name: template?.name ?? "",
  memo: template?.memo ?? "",
  footer: template?.footer ?? "",
  customFields: template?.customFields ?? [],
  grouping: template?.lineItemGrouping ?? "None",
})

const formatStatus = (status?: InvoiceTemplate["status"]) => {
  switch (status) {
    case "ACTIVE":
      return "Active"
    case "DRAFT":
      return "Draft"
    default:
      return status ?? "Draft"
  }
}

const statusVariant = (status?: InvoiceTemplate["status"]) => {
  switch (status) {
    case "ACTIVE":
      return "secondary"
    case "DRAFT":
      return "outline"
    default:
      return "secondary"
  }
}

export default function OrgInvoiceTemplateFormPage() {
  const { orgId, templateId } = useParams()
  const orgName = useOrgStore((state) => state.currentOrg?.name)
  const template = useMemo(
    () => findInvoiceTemplate(templateId),
    [templateId]
  )
  const [name, setName] = useState("")
  const [memo, setMemo] = useState("")
  const [footer, setFooter] = useState("")
  const [grouping, setGrouping] = useState("None")
  const [customFields, setCustomFields] = useState<CustomField[]>([])
  const [previewInvoiceId, setPreviewInvoiceId] = useState("")
  const orgBasePath = orgId ? `/orgs/${orgId}` : "/orgs"

  useEffect(() => {
    const nextState = buildState(template)
    setName(nextState.name)
    setMemo(nextState.memo)
    setFooter(nextState.footer)
    setGrouping(nextState.grouping)
    setCustomFields(nextState.customFields)
  }, [template])

  const isEditing = Boolean(templateId)

  const handleAddCustomField = () => {
    setCustomFields((prev) => [...prev, { label: "", value: "" }])
  }

  const handleCustomFieldChange = (
    index: number,
    key: keyof CustomField,
    value: string
  ) => {
    setCustomFields((prev) =>
      prev.map((field, fieldIndex) =>
        fieldIndex === index ? { ...field, [key]: value } : field
      )
    )
  }

  const handleRemoveCustomField = (index: number) => {
    setCustomFields((prev) => prev.filter((_, fieldIndex) => fieldIndex !== index))
  }

  const totals = previewLineItems.reduce(
    (acc, item) => acc + item.amount,
    0
  )

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold">
            {isEditing ? "Update template" : "Create template"}
          </h1>
          <p className="text-text-muted text-sm">
            Set invoice defaults and preview how the template will render.
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Button asChild variant="outline" size="sm">
            <Link to={`${orgBasePath}/invoice-templates`}>Back to templates</Link>
          </Button>
          {isEditing && (
            <Badge variant={statusVariant(template?.status)}>
              {formatStatus(template?.status)}
            </Badge>
          )}
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
        <Card>
          <CardContent className="space-y-8 pt-6">
            <section className="space-y-3">
              <div>
                <h2 className="text-base font-semibold">Template name</h2>
                <p className="text-text-muted text-sm">
                  Only you and your team will see this.
                </p>
              </div>
              <Input
                placeholder="Enter name"
                value={name}
                onChange={(event) => setName(event.target.value)}
              />
            </section>

            <Separator />

            <section className="space-y-4">
              <div>
                <h2 className="text-base font-semibold">Invoice details</h2>
                <p className="text-text-muted text-sm">
                  Memo and footer appear on every invoice generated from this template.
                </p>
              </div>
              <div className="space-y-2">
                <Label htmlFor="template-memo">Memo</Label>
                <Textarea
                  id="template-memo"
                  placeholder="Have questions? Call us at 555-123-4567"
                  value={memo}
                  onChange={(event) => setMemo(event.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="template-footer">Footer</Label>
                <Textarea
                  id="template-footer"
                  placeholder="Thank you for your business!"
                  value={footer}
                  onChange={(event) => setFooter(event.target.value)}
                />
              </div>
            </section>

            <Separator />

            <section className="space-y-4">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div>
                  <h2 className="text-base font-semibold">Custom fields</h2>
                  <p className="text-text-muted text-sm">
                    Add invoice fields that your team uses for references.
                  </p>
                </div>
                <Button variant="link" size="sm" onClick={handleAddCustomField}>
                  <IconPlus />
                  Add custom field
                </Button>
              </div>
              {customFields.length === 0 ? (
                <div className="text-text-muted text-sm">
                  No custom fields added yet.
                </div>
              ) : (
                <div className="space-y-3">
                  {customFields.map((field, index) => (
                    <div
                      key={`${field.label}-${index}`}
                      className="grid gap-3 md:grid-cols-[1fr_1fr_auto]"
                    >
                      <Input
                        placeholder="Label"
                        value={field.label}
                        onChange={(event) =>
                          handleCustomFieldChange(index, "label", event.target.value)
                        }
                      />
                      <Input
                        placeholder="Value"
                        value={field.value}
                        onChange={(event) =>
                          handleCustomFieldChange(index, "value", event.target.value)
                        }
                      />
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleRemoveCustomField(index)}
                      >
                        Remove
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </section>

            <Separator />

            <section className="space-y-4">
              <div>
                <h2 className="text-base font-semibold">Line item grouping</h2>
                <p className="text-text-muted text-sm">
                  Determine how to group and hide line items on the invoice.{" "}
                  <Button variant="link" size="sm" className="p-0 h-auto">
                    Learn more
                  </Button>
                </p>
              </div>
              <div className="space-y-2">
                <Label htmlFor="grouping-select">Grouping rule</Label>
                <Select value={grouping} onValueChange={setGrouping}>
                  <SelectTrigger id="grouping-select">
                    <SelectValue placeholder="Select a grouping rule" />
                  </SelectTrigger>
                  <SelectContent>
                    {groupingOptions.map((option) => (
                      <SelectItem key={option} value={option}>
                        {option}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </section>

            <div className="flex flex-wrap items-center justify-end gap-2 border-t pt-4">
              <Button asChild variant="outline">
                <Link to={`${orgBasePath}/invoice-templates`}>Cancel</Link>
              </Button>
              <Button>
                {isEditing ? "Update template" : "Create template"}
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card className="h-fit">
          <CardContent className="space-y-5 pt-6">
            <div className="space-y-1">
              <h2 className="text-lg font-semibold">Preview</h2>
              <p className="text-text-muted text-sm">
                Use an invoice that matches the template parameters to see a preview.{" "}
                <Link
                  to={`${orgBasePath}/invoices`}
                  className="text-accent-primary underline underline-offset-4"
                >
                  View all invoices
                </Link>
              </p>
            </div>
            <div className="space-y-2">
              <Label htmlFor="preview-invoice-id">Invoice ID</Label>
              <Input
                id="preview-invoice-id"
                placeholder="ex: in_1Ns8S2DCA5oQnOCe3900uq11R"
                value={previewInvoiceId}
                onChange={(event) => setPreviewInvoiceId(event.target.value)}
              />
            </div>
            <div className="rounded-xl border bg-bg-subtle/30 p-4">
              <div className="rounded-xl bg-bg-primary p-6 shadow-sm">
                <div className="flex items-center justify-between">
                  <div className="text-2xl font-semibold">Invoice</div>
                  <div className="text-lg font-semibold text-text-muted">
                    {orgName ?? "Organization"}
                  </div>
                </div>

                <div className="mt-6 grid gap-4 text-sm md:grid-cols-2">
                  <div className="space-y-1">
                    <div className="text-text-muted">Invoice number</div>
                    <div className="font-medium">EXAMPLE-0001</div>
                    <div className="text-text-muted mt-3">Date of issue</div>
                    <div className="font-medium">January 24, 2026</div>
                    <div className="text-text-muted mt-3">Due date</div>
                    <div className="font-medium">January 24, 2026</div>
                  </div>
                  <div className="space-y-4">
                    <div>
                      <div className="text-text-muted">Bill to</div>
                      <div className="font-medium">Example Customer</div>
                      <div className="text-text-muted">United States</div>
                    </div>
                    <div>
                      <div className="text-text-muted">Line item grouping</div>
                      <div className="font-medium">{grouping}</div>
                    </div>
                  </div>
                </div>

                <div className="mt-6 text-lg font-semibold">
                  {formatCurrency(totals)} due January 24, 2026
                </div>

                <div className="mt-6 rounded-lg border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Description</TableHead>
                        <TableHead>Qty</TableHead>
                        <TableHead>Unit price</TableHead>
                        <TableHead className="text-right">Amount</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {previewLineItems.map((item) => (
                        <TableRow key={item.description}>
                          <TableCell className="font-medium">
                            {item.description}
                          </TableCell>
                          <TableCell>{item.quantity}</TableCell>
                          <TableCell>{formatCurrency(item.unitPrice)}</TableCell>
                          <TableCell className="text-right">
                            {formatCurrency(item.amount)}
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>

                <div className="mt-4 flex justify-end">
                  <div className="grid w-full max-w-xs gap-2 text-sm">
                    <div className="flex items-center justify-between">
                      <span className="text-text-muted">Subtotal</span>
                      <span className="font-medium">{formatCurrency(totals)}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-text-muted">Total</span>
                      <span className="font-semibold">{formatCurrency(totals)}</span>
                    </div>
                  </div>
                </div>

                <Separator className="my-4" />

                <div className="space-y-3 text-sm">
                  <div>
                    <div className="text-text-muted">Memo</div>
                    <div className="font-medium">
                      {memo || "Add a memo to show here."}
                    </div>
                  </div>
                  <div>
                    <div className="text-text-muted">Footer</div>
                    <div className="font-medium">
                      {footer || "Add a footer message for customers."}
                    </div>
                  </div>
                </div>

                <Separator className="my-4" />

                <div className="space-y-2 text-sm">
                  <div className="text-text-muted">Custom fields</div>
                  {customFields.length === 0 ? (
                    <div className="font-medium">No custom fields</div>
                  ) : (
                    <div className="space-y-1">
                      {customFields.map((field, index) => (
                        <div
                          key={`${field.label}-${index}`}
                          className="flex items-center justify-between"
                        >
                          <span className="text-text-muted">
                            {field.label || "Field label"}
                          </span>
                          <span className="font-medium">
                            {field.value || "-"}
                          </span>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
