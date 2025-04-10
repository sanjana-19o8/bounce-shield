import React, { useEffect, useState } from "react";
import { BrowserRouter as Router, Routes, Route, Navigate, useLocation } from "react-router-dom";
import axios from "axios";
import { Toaster } from "react-hot-toast";

import Dashboard from "./pages/Dashboard";
import Verify from "./pages/Verify";
import JobHistory from "./pages/JobHistory";
import AuthForm from "./components/AuthForm";
import PrivateRoute from "./components/PrivateRoute";
import Welcome from "./components/Welcome";


const AppRoutes = ({ token, handleLogout, setToken }) => {
  const location = useLocation();

  return (
    <Routes>
      <Route
        path="/"
        element={
          !token ? (
            <Welcome setToken={setToken} />
          ) : (
            <Navigate to="/dashboard" />
          )
        } />

      <Route
        path="/login"
        element={!token ? <AuthForm setToken={setToken} /> : <Navigate to="/dashboard" />}
      />

      <Route
        path="/dashboard"
        element={
          <PrivateRoute token={token}>
            <Dashboard onLogout={handleLogout} />
          </PrivateRoute>
        }
      />

      <Route
        path="/verify"
        element={
          <PrivateRoute token={token}>
            <Verify />
          </PrivateRoute>
        }
      />

      <Route
        path="/history"
        element={
          <PrivateRoute token={token}>
            <JobHistory />
          </PrivateRoute>
        }
      />
      <Route path="*" element={<Navigate to={token ? "/dashboard" : "/"} />} />
    </Routes>
  );
};

const App = () => {
  const [token, setToken] = useState(null);
  const [checkingAuth, setCheckingAuth] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
      const saved = localStorage.getItem("token");
      if (!saved) return setCheckingAuth(false);

      try {
        const res = await axios.get("/api/verify-token", {
          headers: { Authorization: `Bearer ${saved}` },
        });
        if (res.data.status === "valid") setToken(saved);
        else localStorage.removeItem("token");
      } catch {
        localStorage.removeItem("token");
      } finally {
        setCheckingAuth(false);
      }
    };
    checkAuth();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("token");
    setToken(null);
  };

  if (checkingAuth) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-100 text-gray-700">
        <div className="text-center">
          <div className="loader mb-4"></div>
          <p className="text-xl font-medium">Checking authentication...</p>
        </div>
      </div>
    );
  }

  return (
    <Router>
      <Toaster position="top-center" reverseOrder={false} />
      <AppRoutes token={token} handleLogout={handleLogout} setToken={setToken} />
    </Router>
  );
};

// Loader CSS
const style = document.createElement("style");
style.innerHTML = `
  .loader {
    border: 4px solid #f3f3f3;
    border-top: 4px solid #3498db;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
`;
document.head.appendChild(style);

export default App;
