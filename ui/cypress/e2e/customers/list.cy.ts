describe("Customers list", () => {
  const orgId = Cypress.env("E2E_ORG_ID") || "2002990275537932288"
  const customersPath = `/orgs/${orgId}/customers`

  beforeEach(() => {
    cy.loginAsAdmin({ orgId })
  })

  it("renders customers list and empty state", () => {
    cy.intercept("GET", "/api/customers*").as("getCustomers")

    cy.visit(customersPath)

    cy.findByRole("heading", { name: /^Customers$/, level: 1 }).should(
      "be.visible"
    )
    cy.findByRole("button", { name: /add customer/i }).should("be.visible")

    cy.wait("@getCustomers").then((interception) => {
      expect(interception.response?.statusCode).to.eq(200)
      const payload = interception.response?.body?.data ?? {}
      const customers = payload.customers ?? []
      expect(customers).to.be.an("array")

      if (customers.length > 0) {
        cy.findByRole("table").within(() => {
          cy.findAllByRole("row").should("have.length.greaterThan", 1)
        })
        const firstName = customers[0]?.name
        if (firstName) {
          cy.findByText(firstName).should("be.visible")
        }
      } else {
        cy.findByText(/no customers yet/i).should("be.visible")
      }
    })
  })
})
