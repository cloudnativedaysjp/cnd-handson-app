// Contract for handson-idp (#104). Login/Register are gRPC and have no proto in proto/ yet
// (proposed/idp.proto), so those tests are fixme until the proto is approved.
import { expect, test } from "@playwright/test";
import { createRemoteJWKSet, decodeProtectedHeader, jwtVerify } from "jose";
import { cfg } from "./lib/config";
import { client } from "./lib/grpc";

const idp = () => client("idp.proto", "idp", "IdpService", cfg.idpGrpc, true);
const email = `contract-${Date.now()}@example.com`;
const password = "correct-horse-battery";

test("openid-configuration has jwks_uri", async ({ request }) => {
  const res = await request.get(`${cfg.idpUrl}/.well-known/openid-configuration`);
  expect(res.status()).toBe(200);
  const body = await res.json();
  expect(body.issuer).toBeTruthy();
  expect(body.jwks_uri).toMatch(/\/\.well-known\/jwks\.json$/);
});

test("jwks.json publishes an RS256 RSA key", async ({ request }) => {
  const res = await request.get(`${cfg.idpUrl}/.well-known/jwks.json`);
  expect(res.status()).toBe(200);
  const { keys } = await res.json();
  expect(keys.length).toBeGreaterThan(0);
  expect(keys[0]).toMatchObject({ kty: "RSA", alg: "RS256", use: "sig" });
  expect(keys[0].kid).toBeTruthy();
});

// fixme: needs idp.IdpService (proposed/idp.proto), see #104.
test.fixme("login returns an RS256 JWT verifiable via JWKS", async () => {
  const call = idp();
  await call("Register", { name: "contract", email, password });
  const { access_token } = await call("Login", { email, password });

  expect(decodeProtectedHeader(access_token)).toMatchObject({ alg: "RS256" });
  const jwks = createRemoteJWKSet(new URL(`${cfg.idpUrl}/.well-known/jwks.json`));
  const { payload } = await jwtVerify(access_token, jwks, {
    algorithms: ["RS256"],
    issuer: cfg.expectedIss,
    audience: cfg.expectedAud,
  });
  for (const k of ["iss", "aud", "sub", "exp"]) expect(payload).toHaveProperty(k);
});

// fixme: needs idp.IdpService (proposed/idp.proto), see #104.
test.fixme("wrong password is rejected", async () => {
  const call = idp();
  await call("Register", { name: "contract", email: `bad-${email}`, password });
  await expect(call("Login", { email: `bad-${email}`, password: "wrong" })).rejects.toMatchObject({
    code: 16, // UNAUTHENTICATED
  });
});
