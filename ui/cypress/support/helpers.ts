type OrgFixture = {
  id: string
  name: string
}

export const defaultOrg: OrgFixture = {
  id: "2003007679928209408",
  name: "Acme Org",
}

export const mockOrgEndpoints = (org: OrgFixture = defaultOrg) => {
  cy.intercept("GET", `/api/orgs/${org.id}`, {
    statusCode: 200,
    body: { org },
  }).as("getOrg")

  cy.intercept("GET", "/api/me/orgs", {
    statusCode: 200,
    body: { orgs: [org] },
  }).as("getOrgs")

  return org
}

export const mockLogin = (orgId: string = defaultOrg.id) => {
  cy.intercept("POST", "/auth/login", {
    statusCode: 200,
    body: {
      metadata: {
        user_id: "1",
        username: "cypress",
        password_state: "rotated",
        auth_provider: "local",
      },
    },
  }).as("login")

  cy.visit("/login")
  cy.get('[data-testid="login-username"]').type("cypress")
  cy.get('[data-testid="login-password"]').type("password123")
  cy.get('[data-testid="login-submit"]').click()
  cy.wait("@login")
  cy.wait("@getOrgs")
  cy.wait("@getOrg")
  cy.location("pathname", { timeout: 10000 }).should(
    "match",
    new RegExp(`^/orgs/${orgId}/dashboard$`)
  )
}
