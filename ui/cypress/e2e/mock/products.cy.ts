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
    cy.contains(/no products yet/i).should("be.visible")
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
    cy.contains("Starter").should("be.visible")
    cy.contains("Growth").should("be.visible")
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

    cy.intercept("GET", "/api/meters*", {
      statusCode: 200,
      body: { data: [] },
    }).as("getMeters")

    cy.intercept("GET", "/api/prices*", {
      statusCode: 200,
      body: { data: [] },
    }).as("getPrices")

    cy.intercept("GET", "/api/price_amounts*", {
      statusCode: 200,
      body: { data: [] },
    }).as("getAmounts")

    mockLogin(orgId)
    navigateToProducts()
    cy.wait("@getProducts")
    cy.get('[data-testid="products-create"]').click()
    cy.wait("@getMeters")

    cy.contains("h1", "Create product").should("be.visible")
    fillByTestId("product-name", "Starter")
    fillByTestId("product-code", "starter")
    fillByTestId("product-description", "E2E product description.")
    fillByTestId("price-name", "Starter monthly")
    fillByTestId("product-amount", "5000")

    cy.get('[data-testid="product-submit"]').click()
    cy.wait("@createProduct")
    cy.wait("@createPrice")
    cy.wait("@createAmount")
    cy.wait("@getProduct")
    cy.wait(["@getPrices", "@getAmounts"])
    cy.location("pathname", { timeout: 10000 }).should(
      "eq",
      `/orgs/${orgId}/products/prod-123`
    )
  })
})
