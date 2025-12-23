const appMode = Cypress.env("APP_MODE")
const isCloud = appMode === "cloud"

const fillByTestId = (testId: string, value: string) =>
  cy.get(`[data-testid="${testId}"]`).clear().type(value)

describe("Login, onboarding, and product flows", () => {
  let cleanupPrefix = "e2e"

  before(function () {
    if (!isCloud) {
      this.skip()
    }
  })

  beforeEach(() => {
    cleanupPrefix = `e2e-${Date.now()}`
  })

  afterEach(() => {
    cy.cleanupTestData(cleanupPrefix)
  })

  it("logs in with an existing user", () => {
    cy.seedUser({ prefix: cleanupPrefix }).then((creds) => {
      cy.loginViaUI(creds)
      cy.createOrg({ prefix: cleanupPrefix }).then((org) => {
        cy.visit(`/orgs/${org.id}/dashboard`)
        cy.contains("h1", "Welcome", { timeout: 10000 }).should("be.visible")
      })
    })
  })

  it("creates a new organization via onboarding", () => {
    const orgName = `${cleanupPrefix}-org-${Date.now()}`

    cy.seedUser({ prefix: cleanupPrefix }).then((creds) => {
      cy.loginViaUI(creds)

      cy.visit("/onboarding")
      cy.get('[data-testid="onboarding-org-name"]').should("be.visible")

      cy.get('[data-testid="onboarding-org-name"]').type(orgName)
      cy.get('[data-testid="onboarding-continue"]').click()
      cy.get('[data-testid="onboarding-skip-invites"]').click()
      cy.get('[data-testid="onboarding-skip-billing"]').click()
      cy.get('[data-testid="onboarding-finish"]').click()

      cy.location("pathname", { timeout: 10000 }).should(
        "match",
        /^\/orgs\/\d+\/dashboard$/
      )
      cy.contains(orgName, { timeout: 10000 }).should("be.visible")
    })
  })

  it("creates a product from the org workspace", () => {
    const productName = `Starter ${Date.now()}`
    const productCode = `starter-${Date.now()}`

    cy.seedUser({ prefix: cleanupPrefix }).then((creds) => {
      cy.loginViaUI(creds)
      cy.createOrg({ prefix: cleanupPrefix }).then((org) => {
        cy.visit(`/orgs/${org.id}/products`)
        cy.contains("h1", "Products").should("be.visible")
        cy.get('[data-testid="products-create"]').click()

        cy.contains("h1", "Create product").should("be.visible")
        fillByTestId("product-name", productName)
        fillByTestId("product-code", productCode)
        fillByTestId("product-description", "E2E product description.")
        fillByTestId("product-currency", "USD")
        fillByTestId("product-amount", "5000")

        cy.get('[data-testid="product-submit"]').click()
        cy.location("pathname", { timeout: 20000 }).should(
          "match",
          /^\/orgs\/\d+\/products\/\d+$/
        )
      })
    })
  })
})
