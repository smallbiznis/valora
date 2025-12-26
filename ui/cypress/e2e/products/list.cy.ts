describe("Products list", () => {
  const orgId = Cypress.env("E2E_ORG_ID") || "2002990275537932288"
  const productsPath = `/orgs/${orgId}/products`

  beforeEach(() => {
    cy.loginAsAdmin({ orgId })
  })

  it("renders products, filters, and empty state", () => {
    cy.intercept("GET", "/api/products*").as("getProducts")

    cy.visit(productsPath)

    cy.findByRole("heading", { name: /^Products$/, level: 1 }).should(
      "be.visible"
    )
    cy.findByRole("textbox", { name: /search products/i }).should("be.visible")
    cy.findByText("All").should("be.visible")
    cy.findByText("Active").should("be.visible")
    cy.findByText("Archived").should("be.visible")
    cy.findAllByRole("link", { name: /create product/i })
      .first()
      .should("be.visible")

    cy.wait("@getProducts").then((interception) => {
      expect(interception.response?.statusCode).to.eq(200)
      const products = interception.response?.body?.data ?? []
      expect(products).to.be.an("array")

      if (products.length > 0) {
        cy.findByRole("table").within(() => {
          cy.findAllByRole("row").should("have.length.greaterThan", 1)
        })
      } else {
        cy.findByText(/no products yet/i).should("be.visible")
      }
    })
  })
})
