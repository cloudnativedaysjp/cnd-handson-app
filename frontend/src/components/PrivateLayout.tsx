import type React from "react";
import { useSelector } from "react-redux";
import { Navigate, Outlet } from "react-router-dom";
import type { RootState } from "../store";
import LivePanel from "./LivePanel";

// 画面の枠（ヘッダーやメニュー）は、legacy と modern の画面がそれぞれ持つ
const PrivateLayout: React.FC = () => {
  const isLoggedIn = useSelector(
    (state: RootState) => state.session.isLoggedIn,
  );

  if (!isLoggedIn) return <Navigate to="/login" replace />;

  return (
    <>
      <Outlet />
      <LivePanel />
    </>
  );
};

export default PrivateLayout;
