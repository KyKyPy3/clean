import { useEffect } from "react";
import { Outlet, useNavigate } from "react-router-dom";

import Layout from "@components/Layout";
import { useSessionStore } from "@session/infrastructure/controller/http/v1/store";

export const ProtectedRoute = () => {
  const navigate = useNavigate();
  const { token } = useSessionStore();

  useEffect(() => {
    if (token === '') {
      navigate("/signin");
    }
  });

  return (
    <Layout>
      <Outlet />
    </Layout>
  );
}
