describe("OSS mode login and signup", () => {
  const envMode = Cypress.env("APP_MODE")
  const mode: "oss" | "cloud" = envMode === "cloud" ? "cloud" : "oss"

  before(function (this: Mocha.Context) {
    if (mode !== "oss") {
      this.skip()
    }
  })

  it("hides sign up entry on login", () => {
    cy.visit("/login")
    cy.get('a[href="/signup"]').should("not.be.visible")
  })

  it("redirects direct signup visits to login", () => {
    cy.visit("/signup")
    const username = `e2e-oss-${Date.now()}`
    cy.get('input[id="username"]').type(username)
    cy.get('input[id="password"]').type("password123")
    cy.contains("button", "Create Account").click()

    cy.location("pathname", { timeout: 5000 }).should("eq", "/login")
  })
})

describe("Cloud mode login and signup", () => {
  const envMode = Cypress.env("APP_MODE")
  const mode: "oss" | "cloud" = envMode === "cloud" ? "cloud" : "oss"

  before(function (this: Mocha.Context) {
    if (mode !== "cloud") {
      this.skip()
    }
  })

  it("shows sign up entry on login", () => {
    cy.visit("/login")
    cy.get('a[href="/signup"]').should("be.visible")
  })

  it("allows signup and redirects to organizations", () => {
    cy.visit("/signup")

    const username = `e2e-${Date.now()}`
    cy.get('input[id="username"]').type(username)
    cy.get('input[id="password"]').type("password123")
    cy.contains("button", "Create Account").click()

    cy.location("pathname", { timeout: 10000 }).should("match", /^\/orgs/)
  })
})
