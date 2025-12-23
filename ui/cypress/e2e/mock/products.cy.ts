import { defaultOrg, mockLogin, mockOrgEndpoints } from "../../support/helpers"

const orgId = defaultOrg.id

const navigateToProducts = () => {
  cy.contains("a", "Products").click()
  cy.location("pathname", { timeout: 10000 }).should(
    "eq",
    `/orgs/${orgId}/products`
  )
}

const fillByTestId = (testId: string, value: string) => {
  cy.get(`[data-testid="${testId}"]`).clear().type(value)
}

describe("Products page", () => {
  beforeEach(() => {
    mockOrgEndpoints()
  })

  it("shows empty state", () => {
    cy.intercept({ method: "GET", pathname: "/api/products" }, {
      statusCode: 200,
      body: { data: [] },
    }).as("getProducts")

    mockLogin(orgId)
    navigateToProducts()
    cy.wait("@getProducts")
    cy.get('[data-testid="products-json"]').should("contain.text", "[]")
  })

  it("lists products", () => {
    cy.intercept({ method: "GET", pathname: "/api/products" }, {
      statusCode: 200,
      body: {
        data: [
          { id: "prod-1", name: "Starter", code: "starter", active: true },
          { id: "prod-2", name: "Growth", code: "growth", active: false },
        ],
      },
    }).as("getProducts")

    mockLogin(orgId)
    navigateToProducts()
    cy.wait("@getProducts")
    cy.get('[data-testid="products-json"]').should("contain.text", "Starter")
    cy.get('[data-testid="products-json"]').should("contain.text", "Growth")
  })

  it("creates a product", () => {
    cy.intercept({ method: "GET", pathname: "/api/products" }, {
      statusCode: 200,
      body: { data: [] },
    }).as("getProducts")

    cy.intercept("POST", "/api/products", {
      statusCode: 200,
      body: { data: { id: "prod-123", name: "Starter", code: "starter" } },
    }).as("createProduct")

    cy.intercept("POST", "/api/prices", {
      statusCode: 200,
      body: { data: { id: "price-123" } },
    }).as("createPrice")

    cy.intercept("POST", "/api/price_amounts", {
      statusCode: 200,
      body: { data: { id: "amount-123" } },
    }).as("createAmount")

    cy.intercept("GET", "/api/products/prod-123*", {
      statusCode: 200,
      body: { data: { id: "prod-123", name: "Starter", code: "starter" } },
    }).as("getProduct")

    mockLogin(orgId)
    navigateToProducts()
    cy.wait("@getProducts")
    cy.get('[data-testid="products-create"]').click()

    cy.contains("h1", "Create product").should("be.visible")
    fillByTestId("product-name", "Starter")
    fillByTestId("product-code", "starter")
    fillByTestId("product-description", "E2E product description.")
    fillByTestId("product-amount", "5000")

    cy.get('[data-testid="product-submit"]').click()
    cy.wait("@createProduct")
    cy.wait("@createPrice")
    cy.wait("@createAmount")
    cy.wait("@getProduct")
    cy.location("pathname", { timeout: 10000 }).should(
      "eq",
      `/orgs/${orgId}/products/prod-123`
    )
  })
})
