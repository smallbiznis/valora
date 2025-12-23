type SeedUserOptions = {
  prefix?: string
}

type SeededUser = {
  username: string
  password: string
  prefix: string
}

type CreateOrgOptions = {
  prefix?: string
  name?: string
}

type CreatedOrg = {
  id: string
  name: string
}

type CreateCustomerOptions = {
  prefix?: string
  name?: string
  email?: string
}

type CreatedCustomer = {
  id: string
  name: string
  email: string
}

const buildPrefix = (prefix?: string) => prefix?.trim() || Cypress.env("E2E_PREFIX") || "e2e"

const uniqueToken = () => `${Date.now()}-${Cypress._.random(0, 1_000_000)}`

Cypress.Commands.add("seedUser", (options: SeedUserOptions = {}) => {
  const prefix = buildPrefix(options.prefix)
  const username = `${prefix}-user-${uniqueToken()}`
  const password = "password123"

  return cy
    .request("POST", "/auth/signup", { username, password })
    .then(() => cy.wrap({ username, password, prefix }))
})

Cypress.Commands.add("loginViaUI", (creds: SeededUser) => {
  cy.visit("/login")
  cy.get('[data-testid="login-username"]').type(creds.username)
  cy.get('[data-testid="login-password"]').type(creds.password)
  cy.get('[data-testid="login-submit"]').click()
  cy.location("pathname", { timeout: 10000 }).should("match", /^\/orgs(\/|$)/)
})

Cypress.Commands.add("createOrg", (options: CreateOrgOptions = {}) => {
  const prefix = buildPrefix(options.prefix)
  const name = options.name ?? `${prefix}-org-${uniqueToken()}`

  return cy
    .request("POST", "/api/organizations", {
      name,
      country_code: "ID",
      timezone_name: "Asia/Jakarta",
      default_currency: "IDR",
    })
    .then((res) => cy.wrap({ id: res.body?.id, name }))
})

Cypress.Commands.add("createCustomer", (orgId: string, options: CreateCustomerOptions = {}) => {
  const prefix = buildPrefix(options.prefix)
  const name = options.name ?? `${prefix}-customer-${uniqueToken()}`
  const email = options.email ?? `${name}@example.test`

  return cy
    .request("POST", "/api/customers", {
      organization_id: orgId,
      name,
      email,
    })
    .then((res) => cy.wrap({ id: res.body?.data?.id, name, email }))
})

Cypress.Commands.add("cleanupTestData", (prefix?: string) => {
  const value = buildPrefix(prefix)
  return cy.request("POST", "/api/test/cleanup", { prefix: value })
})

declare global {
  namespace Cypress {
    interface Chainable {
      seedUser(options?: SeedUserOptions): Chainable<SeededUser>
      loginViaUI(creds: SeededUser): Chainable<void>
      createOrg(options?: CreateOrgOptions): Chainable<CreatedOrg>
      createCustomer(orgId: string, options?: CreateCustomerOptions): Chainable<CreatedCustomer>
      cleanupTestData(prefix?: string): Chainable<void>
    }
  }
}

export {}
