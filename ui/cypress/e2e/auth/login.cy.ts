describe("Admin session bootstrap", () => {
  const orgId = Cypress.env("E2E_ORG_ID") || "2002990275537932288"

  it("loads the dashboard with org context", () => {
    cy.intercept("POST", `/api/user/using/${orgId}`).as("useOrg")
    cy.intercept("GET", `/api/orgs/${orgId}`).as("getOrg")

    cy.loginAsAdmin({ orgId })

    cy.location("pathname", { timeout: 10000 }).should(
      "eq",
      `/orgs/${orgId}/dashboard`
    )
    cy.wait("@useOrg")
      .its("response.statusCode")
      .should("eq", 200)

    cy.wait("@getOrg").then((interception) => {
      expect(interception.response?.statusCode).to.eq(200)
      const orgName =
        interception.response?.body?.org?.name || `Org ${orgId}`
      cy.findByRole("button", { name: orgName }).should("be.visible")
    })

    cy.findByRole("link", { name: /^Products$/ }).should("be.visible")
  })
})
