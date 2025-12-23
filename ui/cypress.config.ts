import { defineConfig } from "cypress"

const baseUrl = process.env.VITE_APP_BASE_URL || "http://localhost:8080"
const appMode =
  process.env.CYPRESS_APP_MODE ||
  process.env.APP_MODE ||
  process.env.VITE_APP_MODE
const testSuite = process.env.CYPRESS_TEST_SUITE || "all"
const specPattern =
  testSuite === "mock"
    ? "cypress/e2e/mock/**/*.cy.ts"
    : testSuite === "real"
      ? "cypress/e2e/real/**/*.cy.ts"
      : "cypress/e2e/**/*.cy.ts"

export default defineConfig({
  e2e: {
    baseUrl,
    env: {
      APP_MODE: appMode,
    },
    specPattern,
    supportFile: "cypress/support/e2e.ts",
    video: false,
  },
})
