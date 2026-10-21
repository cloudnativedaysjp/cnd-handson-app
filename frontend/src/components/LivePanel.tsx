import type React from "react";
import { useSyncExternalStore } from "react";
import { recentRequests } from "../utils/api";
import "./live.css";

// legacy と modern は /color と同じ青と緑にする。ほかは名前から紫の帯の中で決め、エラーの赤と紛れないようにする
function hue(version: string) {
  if (version === "legacy") return 210;
  if (version === "modern") return 140;
  let h = 0;
  for (const c of version) h = (h * 31 + c.charCodeAt(0)) % 60;
  return 260 + h;
}

const LivePanel: React.FC = () => {
  const logs = useSyncExternalStore(
    recentRequests.subscribe,
    recentRequests.get,
  );
  return (
    <details className="live-panel" open>
      <summary>応答したバージョン</summary>
      {logs.length === 0 ? (
        <p>画面を操作すると、API の応答がここに並びます。</p>
      ) : (
        <ol>
          {logs.map((l, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: 同じリクエストが何度も並ぶので、順番が唯一の区別
            <li key={i} className={l.status >= 500 ? "is-error" : undefined}>
              <span
                className="live-version"
                style={{ "--hue": hue(l.version) } as React.CSSProperties}
              >
                {l.version || "不明"}
              </span>
              <code>
                {l.method} {l.path}
              </code>
              <span>{l.status >= 400 ? l.status : `${l.ms} ms`}</span>
            </li>
          ))}
        </ol>
      )}
    </details>
  );
};

export default LivePanel;
