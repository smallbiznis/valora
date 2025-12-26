type Invoice = {
  status?: string
  Status?: string
  paid_at?: string
  PaidAt?: string
  due_at?: string
  DueAt?: string
}

const readField = (invoice: Invoice, keys: (keyof Invoice)[]) => {
  for (const key of keys) {
    const value = invoice[key]
    if (typeof value === "string" && value.trim()) {
      return value
    }
  }
  return ""
}

const deriveStatus = (invoice: Invoice) => {
  const raw = readField(invoice, ["status", "Status"]).toUpperCase() || "UNKNOWN"
  const paidAt = readField(invoice, ["paid_at", "PaidAt"])
  const dueAt = readField(invoice, ["due_at", "DueAt"])
  const dueDate = dueAt ? new Date(dueAt) : null
  if (paidAt || raw === "PAID") return "PAID"
  if (raw === "VOID") return "VOID"
  if (raw === "UNCOLLECTIBLE") return "UNCOLLECTIBLE"
  if (raw === "DRAFT") return "DRAFT"
  if (raw === "OPEN" || raw === "ISSUED") {
    if (dueDate && dueDate.getTime() < Date.now()) return "PAST_DUE"
    return "OPEN"
  }
  if (dueDate && dueDate.getTime() < Date.now()) return "PAST_DUE"
  return raw
}

const formatStatus = (status?: string) => {
  switch (status) {
    case "PAST_DUE":
      return "Past due"
    case "DRAFT":
      return "Draft"
    case "OPEN":
      return "Open"
    case "PAID":
      return "Paid"
    case "VOID":
      return "Void"
    case "UNCOLLECTIBLE":
      return "Uncollectible"
    default:
      return status || "-"
  }
}

describe("Invoices list", () => {
  const orgId = Cypress.env("E2E_ORG_ID") || "2002990275537932288"
  const invoicesPath = `/orgs/${orgId}/invoices`

  beforeEach(() => {
    cy.loginAsAdmin({ orgId })
  })

  it("renders invoices list and empty state", () => {
    cy.intercept("GET", "/api/invoices*").as("getInvoices")
    cy.intercept("GET", "/api/customers*").as("getCustomers")

    cy.visit(invoicesPath)

    cy.findByRole("heading", { name: /^Invoices$/, level: 1 }).should(
      "be.visible"
    )
    cy.findAllByRole("button", { name: /create invoice/i })
      .first()
      .should("be.visible")

    cy.wait("@getCustomers").its("response.statusCode").should("eq", 200)
    cy.wait("@getInvoices").then((interception) => {
      expect(interception.response?.statusCode).to.eq(200)
      const invoices = interception.response?.body?.data ?? []
      expect(invoices).to.be.an("array")

      if (invoices.length > 0) {
        cy.findByRole("table").within(() => {
          cy.findByRole("columnheader", { name: /total/i }).should("be.visible")
          cy.findByRole("columnheader", {
            name: /invoice number/i,
          }).should("be.visible")
          cy.findByRole("columnheader", {
            name: /customer/i,
          }).should("be.visible")
        })

        const statusLabel = formatStatus(deriveStatus(invoices[0] as Invoice))
        if (statusLabel && statusLabel !== "-") {
          cy.findByText(statusLabel).should("be.visible")
        }
      } else {
        cy.findByText(/no invoices yet/i).should("be.visible")
      }
    })
  })
})
