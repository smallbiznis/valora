describe("Products create flow", () => {
  const orgId = Cypress.env("E2E_ORG_ID") || "2002990275537932288"
  const createPath = `/orgs/${orgId}/products/create`

  beforeEach(() => {
    cy.loginAsAdmin({ orgId })
  })

  it("creates a product with flat pricing", () => {
    const token = Date.now()
    const productName = `E2E Product ${token}`
    const productCode = `e2e-product-${token}`
    const priceName = `E2E Price ${token}`

    cy.intercept("GET", "/api/meters*").as("getMeters")
    cy.intercept("POST", "/api/products").as("createProduct")
    cy.intercept("POST", "/api/prices").as("createPrice")
    cy.intercept("POST", "/api/price_amounts").as("createPriceAmount")

    cy.visit(createPath)

    cy.findByRole("heading", { name: /create product/i, level: 1 }).should(
      "be.visible"
    )
    cy.wait("@getMeters").then((interception) => {
      expect(interception.response?.statusCode).to.eq(200)
      const meters = interception.response?.body?.data ?? []
      expect(meters).to.be.an("array")
    })

    cy.findByTestId("product-name").type(productName)
    cy.findByTestId("product-code").type(productCode)
    cy.findByTestId("price-name").type(priceName)
    cy.findByTestId("product-amount").clear().type("5000")
    cy.findByTestId("product-submit").click()

    cy.wait("@createProduct").then((interception) => {
      expect(interception.response?.statusCode).to.eq(200)
      const productId = interception.response?.body?.data?.id
      expect(productId).to.exist
      cy.wrap(productId).as("productId")
    })
    cy.wait("@createPrice").its("response.statusCode").should("eq", 200)
    cy.wait("@createPriceAmount").its("response.statusCode").should("eq", 200)

    cy.get("@productId").then((productId) => {
      cy.location("pathname", { timeout: 10000 }).should(
        "eq",
        `/orgs/${orgId}/products/${productId}`
      )
    })
    cy.findByRole("heading", { name: productName, level: 1 }).should(
      "be.visible"
    )
  })
})
