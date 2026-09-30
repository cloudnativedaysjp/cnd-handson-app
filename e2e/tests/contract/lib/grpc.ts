import path from "node:path";
import * as grpc from "@grpc/grpc-js";
import * as loader from "@grpc/proto-loader";

const protoRoot = path.resolve(import.meta.dirname, "../../../../proto");
const proposedRoot = path.resolve(import.meta.dirname, "../proposed");

// Loads protos straight from the root proto/ (no codegen). `proposed: true` loads from ./proposed.
export function client<T = Record<string, (...a: any[]) => void>>(
  file: string,
  pkg: string,
  service: string,
  addr: string,
  proposed = false,
) {
  const def = loader.loadSync(file, {
    includeDirs: [proposed ? proposedRoot : protoRoot],
    keepCase: true,
    defaults: true,
    longs: String,
  });
  let ns: any = grpc.loadPackageDefinition(def);
  for (const p of pkg.split(".")) ns = ns[p];
  const c = new ns[service](addr, grpc.credentials.createInsecure());
  // Promisified unary call; every request carries the user id as `x-user-id` (docs: #72/#105).
  return (method: string, req: object, userId = USER_ID) =>
    new Promise<any>((resolve, reject) => {
      const md = new grpc.Metadata();
      md.set("x-user-id", userId);
      (c as any)[method](req, md, { deadline: Date.now() + 5000 }, (e: Error | null, r: any) =>
        e ? reject(e) : resolve(r),
      );
    });
}

export { grpc };

// Project ids are UUID columns, so the caller id must be a UUID too.
export const USER_ID = "00000000-0000-4000-8000-000000000071";
