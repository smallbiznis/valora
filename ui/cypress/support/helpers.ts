type OrgFixture = {
  id: string
  name: string
}

export const defaultOrg: OrgFixture = {
  id: "2003007679928209408",
  name: "Acme Org",
}

export const mockOrgEndpoints = (org: OrgFixture = defaultOrg) => {
  cy.intercept("POST", `/api/user/using/${org.id}`, {
    statusCode: 200,
    body: {
      metadata: {
        active_org_id: org.id,
        org_ids: [org.id],
      },
    },
  }).as("useOrg")

  cy.intercept("GET", `/api/orgs/${org.id}`, {
    statusCode: 200,
    body: { org },
  }).as("getOrg")

  cy.intercept("GET", "/api/user/orgs", {
    statusCode: 200,
    body: { orgs: [org] },
  }).as("getOrgs")

  cy.intercept("GET", "/api/me/orgs", {
    statusCode: 200,
    body: { orgs: [org] },
  }).as("getMeOrgs")

  return org
}

export const mockLogin = (orgId: string = defaultOrg.id) => {
  const authOrgId = Cypress.env("E2E_ORG_ID") || orgId
  cy.loginAsAdmin({ orgId: authOrgId })
}
