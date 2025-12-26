import { defaultOrg, mockLogin, mockOrgEndpoints } from "../../support/helpers"

const orgId = defaultOrg.id

describe("Customers page", () => {
  beforeEach(() => {
    mockOrgEndpoints()
  })

  it("shows empty state", () => {
    cy.intercept("GET", "/api/customers*", {
      statusCode: 200,
      body: { data: { customers: [], has_more: false } },
    }).as("getCustomers")

    mockLogin(orgId)
    cy.visit(`/orgs/${orgId}/customers`)
    cy.wait(["@useOrg", "@getOrg", "@getOrgs", "@getCustomers"])

    cy.contains("No customers yet.").should("be.visible")
  })

  it("lists customers", () => {
    cy.intercept("GET", "/api/customers*", {
      statusCode: 200,
      body: {
        data: {
          customers: [
            { id: "1", name: "Acme Co", email: "billing@acme.test" },
            { id: "2", name: "Globex", email: "finance@globex.test" },
          ],
          has_more: false,
        },
      },
    }).as("getCustomers")

    mockLogin(orgId)
    cy.visit(`/orgs/${orgId}/customers`)
    cy.wait(["@useOrg", "@getOrg", "@getOrgs", "@getCustomers"])

    cy.contains("Acme Co").should("be.visible")
    cy.contains("billing@acme.test").should("be.visible")
    cy.contains("Globex").should("be.visible")
    cy.contains("finance@globex.test").should("be.visible")
  })

  it("creates a customer", () => {
    cy.intercept("GET", "/api/customers*", {
      statusCode: 200,
      body: { data: { customers: [], has_more: false } },
    }).as("getCustomers")

    cy.intercept("POST", "/api/customers", {
      statusCode: 200,
      body: {
        data: {
          id: "3",
          name: "Umbrella Corp",
          email: "billing@umbrella.test",
        },
      },
    }).as("createCustomer")

    mockLogin(orgId)
    cy.visit(`/orgs/${orgId}/customers`)
    cy.wait(["@useOrg", "@getOrg", "@getOrgs", "@getCustomers"])

    cy.contains("button", "Add customer").click()
    cy.get("#customer-name").type("Umbrella Corp")
    cy.get("#customer-email").type("billing@umbrella.test")
    cy.get('[data-slot="dialog-content"]').contains("button", "Add customer").click()
    cy.wait("@createCustomer")

    cy.contains("Umbrella Corp").should("be.visible")
    cy.contains("billing@umbrella.test").should("be.visible")
  })

  it("shows validation error on create", () => {
    cy.intercept("GET", "/api/customers*", {
      statusCode: 200,
      body: { data: { customers: [], has_more: false } },
    }).as("getCustomers")

    cy.intercept("POST", "/api/customers", {
      statusCode: 400,
      body: { error: { type: "validation_error", message: "invalid email" } },
    }).as("createCustomer")

    mockLogin(orgId)
    cy.visit(`/orgs/${orgId}/customers`)
    cy.wait(["@useOrg", "@getOrg", "@getOrgs", "@getCustomers"])

    cy.contains("button", "Add customer").click()
    cy.get("#customer-name").type("Bad Email")
    cy.get("#customer-email").type("bad@example.com")
    cy.get('[data-slot="dialog-content"]').contains("button", "Add customer").click()
    cy.wait("@createCustomer")

    cy.get('[data-slot="alert"]').should("be.visible")
  })
})
