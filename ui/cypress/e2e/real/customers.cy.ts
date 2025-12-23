describe("Customers (real)", () => {
  const envMode = Cypress.env("APP_MODE")
  const mode: "oss" | "cloud" = envMode === "cloud" ? "cloud" : "oss"

  let cleanupPrefix = "e2e"

  before(function (this: Mocha.Context) {
    if (mode !== "cloud") {
      this.skip()
    }
  })

  beforeEach(() => {
    cleanupPrefix = `e2e-${Date.now()}`
  })

  afterEach(() => {
    cy.cleanupTestData(cleanupPrefix)
  })

  it("creates and lists customers", () => {
    cy.seedUser({ prefix: cleanupPrefix }).then((creds) => {
      cy.loginViaUI(creds)

      cy.createOrg({ prefix: cleanupPrefix }).then((org) => {
        const customerName = `${cleanupPrefix}-customer`
        const customerEmail = `${cleanupPrefix}-customer@example.test`

        cy.visit(`/orgs/${org.id}/customers`)
        cy.contains("button", "Add customer").click()
        cy.get("#customer-name").type(customerName)
        cy.get("#customer-email").type(customerEmail)
        cy.get('[data-slot="dialog-content"]').contains("button", "Add customer").click()

        cy.contains(customerName).should("be.visible")
        cy.contains(customerEmail).should("be.visible")
      })
    })
  })
})
