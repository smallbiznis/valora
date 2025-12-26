describe("Settings page", () => {
  const orgId = Cypress.env("E2E_ORG_ID") || "2002990275537932288"
  const settingsPath = `/orgs/${orgId}/settings`

  beforeEach(() => {
    cy.intercept("POST", `/api/user/using/${orgId}`).as("useOrg")
    cy.intercept("GET", `/api/orgs/${orgId}`).as("getOrg")
    cy.loginAsAdmin({ orgId })
  })

  it("loads workspace settings", () => {
    cy.wait("@useOrg").its("response.statusCode").should("eq", 200)
    cy.wait("@getOrg").its("response.statusCode").should("eq", 200)

    cy.findByRole("link", { name: /^Settings$/ }).click()
    cy.location("pathname").should("eq", settingsPath)
    cy.findByRole("heading", { name: /^Settings$/, level: 1 }).should(
      "be.visible"
    )
    cy.findByText(/^Workspace$/).should("be.visible")
  })
})
