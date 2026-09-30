import { execFileSync } from "node:child_process";

// Returns JSON log lines emitted by a compose service since `since` (RFC3339).
export function jsonLogs(service: string, since: string): Record<string, unknown>[] {
  const out = execFileSync("docker", ["compose", "logs", "--no-color", "--since", since, service], {
    encoding: "utf8",
  });
  return out.split("\n").flatMap((l) => {
    const i = l.indexOf("{"); // strip the "service-1  | " prefix
    if (i < 0) return [];
    try {
      return [JSON.parse(l.slice(i))];
    } catch {
      return [];
    }
  });
}
