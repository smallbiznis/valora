describe("Sidebar navigation", () => {
  const orgId = Cypress.env("E2E_ORG_ID") || "2002990275537932288"
  const orgPath = (path: string) => `/orgs/${orgId}/${path}`

  const menuItems = [
    { label: "Products", path: "products", heading: "Products" },
    { label: "Meters", path: "meter", heading: "Meters" },
    { label: "Customers", path: "customers", heading: "Customers" },
    { label: "Subscriptions", path: "subscriptions", heading: "Subscriptions" },
    { label: "Invoices", path: "invoices", heading: "Invoices" },
    { label: "Settings", path: "settings", heading: "Settings" },
  ]

  beforeEach(() => {
    cy.loginAsAdmin({ orgId })
  })

  menuItems.forEach((item) => {
    it(`navigates to ${item.label}`, () => {
      cy.findByRole("link", {
        name: new RegExp(`^${item.label}$`, "i"),
      }).click()

      cy.location("pathname").should("eq", orgPath(item.path))
      cy.findByRole("heading", {
        name: new RegExp(`^${item.heading}$`, "i"),
        level: 1,
      }).should("be.visible")
    })
  })
})
