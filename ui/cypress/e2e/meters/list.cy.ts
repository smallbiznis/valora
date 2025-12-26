describe("Meters list", () => {
  const orgId = Cypress.env("E2E_ORG_ID") || "2002990275537932288"
  const metersPath = `/orgs/${orgId}/meter`

  beforeEach(() => {
    cy.loginAsAdmin({ orgId })
  })

  it("renders meters list and empty state", () => {
    cy.intercept("GET", "/api/meters*").as("getMeters")

    cy.visit(metersPath)

    cy.findByRole("heading", { name: /^Meters$/, level: 1 }).should(
      "be.visible"
    )
    cy.findByRole("link", { name: /create meter/i }).should("be.visible")

    cy.wait("@getMeters").then((interception) => {
      expect(interception.response?.statusCode).to.eq(200)
      const meters = interception.response?.body?.data ?? []
      expect(meters).to.be.an("array")

      if (meters.length > 0) {
        cy.findByRole("table").within(() => {
          cy.findAllByRole("row").should("have.length.greaterThan", 1)
        })
      } else {
        cy.findByText(/no meters yet/i).should("be.visible")
      }
    })
  })
})
