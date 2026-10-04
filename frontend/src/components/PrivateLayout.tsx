import type React from "react";
import { useSelector } from "react-redux";
import { Navigate, Outlet } from "react-router-dom";
import type { RootState } from "../store";
import Header from "./Header";
import { legacy } from "./Legacy";

const PrivateLayout: React.FC = () => {
  const isLoggedIn = useSelector(
    (state: RootState) => state.session.isLoggedIn,
  );

  if (!isLoggedIn) return <Navigate to="/login" replace />;

  return (
    <>
      {/* legacy の画面は、それぞれが自分の枠（帯とメニュー）を持つ */}
      {!legacy && <Header />}
      <Outlet />
    </>
  );
};

export default PrivateLayout;
