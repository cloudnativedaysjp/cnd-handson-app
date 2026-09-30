// Single source of addresses for contract tests. Override via env; defaults are the compose ports.
const env = (k: string, d: string) => process.env[k] || d;

export const cfg = {
  idpUrl: env("IDP_URL", "http://localhost:8082"), // HTTP 8080 in the container: discovery + JWKS
  idpGrpc: env("IDP_GRPC_ADDR", "localhost:50051"),
  projectGrpc: env("PROJECT_GRPC_ADDR", "localhost:50053"),
  taskGrpc: env("TASK_GRPC_ADDR", "localhost:50055"),
  // Seeded/registered by the idp implementation (#104) for the UI test.
  userEmail: env("CONTRACT_USER_EMAIL", "demo@example.com"),
  userPassword: env("CONTRACT_USER_PASSWORD", "demo-password"),
  // docs/conventions.md: iss/aud are env-configured on the idp; set these to assert exact values.
  expectedIss: process.env.IDP_ISS,
  expectedAud: process.env.IDP_AUD,
};

// ENTRY_URL is the legacy entry (kept for the single-entry case); modern defaults to 8081.
export const entries = [
  {
    variant: "legacy",
    color: "blue",
    url: env("ENTRY_LEGACY_URL", env("ENTRY_URL", "http://localhost:8080")),
    service: env("ENTRY_LEGACY_SERVICE", "handson-legacy"), // compose service name for `docker compose logs`
  },
  {
    variant: "modern",
    color: "green",
    url: env("ENTRY_MODERN_URL", "http://localhost:8081"),
    service: env("ENTRY_MODERN_SERVICE", "handson-modern"),
  },
] as const;
