type SeedUserOptions = {
  prefix?: string
}

type SeededUser = {
  username: string
  email: string
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

type LoginAsAdminOptions = {
  orgId?: string
  username?: string
  email?: string
  password?: string
  sessionToken?: string
}

const buildPrefix = (prefix?: string) =>
  prefix?.trim() || Cypress.env("E2E_PREFIX") || "e2e"

const uniqueToken = () => `${Date.now()}-${Cypress._.random(0, 1_000_000)}`

const buildAuthState = (user: { id: string; username: string; email?: string }) => {
  return JSON.stringify({
    state: {
      user,
      isAuthenticated: true,
    },
    version: 0,
  })
}

Cypress.Commands.add("seedUser", (options: SeedUserOptions = {}) => {
  const prefix = buildPrefix(options.prefix)
  const username = `${prefix}-user-${uniqueToken()}`
  const email = `${username}@example.test`
  const password = "password123"

  return cy
    .request("POST", "/internal/auth/local/signup", {
      email,
      password,
      display_name: username,
    })
    .then(() => cy.wrap({ username, email, password, prefix }))
})

Cypress.Commands.add("loginViaUI", (creds: SeededUser) => {
  cy.visit("/login")
  cy.get('[data-testid="login-email"]').type(creds.email)
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

Cypress.Commands.add(
  "createCustomer",
  (orgId: string, options: CreateCustomerOptions = {}) => {
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
  }
)

Cypress.Commands.add("cleanupTestData", (prefix?: string) => {
  const value = buildPrefix(prefix)
  return cy.request("POST", "/api/test/cleanup", { prefix: value })
})

Cypress.Commands.add("loginAsAdmin", (options: LoginAsAdminOptions = {}) => {
  const orgId =
    options.orgId || Cypress.env("E2E_ORG_ID") || "2002990275537932288"
  const username =
    options.username || Cypress.env("E2E_ADMIN_USERNAME") || "admin"
  const email =
    options.email ||
    Cypress.env("E2E_ADMIN_EMAIL") ||
    (username.includes("@") ? username : `${username}@example.com`)
  const password =
    options.password || Cypress.env("E2E_ADMIN_PASSWORD") || "ChangeMe123!"
  const sessionToken =
    options.sessionToken || Cypress.env("E2E_SESSION_TOKEN") || ""

  const startSession = (user: { id: string; username: string; email?: string }) => {
    return cy.request("POST", `/api/user/using/${orgId}`).then((useOrg) => {
      expect(useOrg.status).to.eq(200)
      cy.visit(`/orgs/${orgId}/dashboard`, {
        onBeforeLoad(win) {
          win.localStorage.setItem("valora-auth", buildAuthState(user))
        },
      })
    })
  }

  if (sessionToken) {
    cy.setCookie("_sid", sessionToken, { httpOnly: true, path: "/" })
    const user = {
      id: Cypress.env("E2E_ADMIN_USER_ID") || "admin",
      username,
      email,
    }
    return startSession(user)
  }

  return cy
    .request({
      method: "POST",
      url: "/internal/auth/local/login",
      body: { email, password },
      failOnStatusCode: false,
    })
    .then((response) => {
      if (response.status !== 200) {
        throw new Error(
          `Admin login failed (${response.status}). Set CYPRESS_E2E_ADMIN_EMAIL/CYPRESS_E2E_ADMIN_PASSWORD or CYPRESS_E2E_SESSION_TOKEN.`
        )
      }
      const metadata = response.body?.metadata ?? {}
      const displayName = metadata.display_name || metadata.username || email
      const user = {
        id:
          metadata.user_id ||
          Cypress.env("E2E_ADMIN_USER_ID") ||
          "admin",
        username: displayName || username,
        email: metadata.email || email,
      }

      return startSession(user)
    })
})

declare global {
  namespace Cypress {
    interface Chainable {
      seedUser(options?: SeedUserOptions): Chainable<SeededUser>
      loginViaUI(creds: SeededUser): Chainable<void>
      createOrg(options?: CreateOrgOptions): Chainable<CreatedOrg>
      createCustomer(
        orgId: string,
        options?: CreateCustomerOptions
      ): Chainable<CreatedCustomer>
      cleanupTestData(prefix?: string): Chainable<void>
      loginAsAdmin(options?: LoginAsAdminOptions): Chainable<void>
    }
  }
}

export {}
