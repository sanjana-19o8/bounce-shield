import React, { useState } from "react";
import AuthForm from "./AuthForm";
import EmailVerifyLite from "./EmailVerifyLite";

const Welcome = ({ setToken }) => {
  const [showLogin, setShowLogin] = useState(false);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-r from-blue-500 to-purple-600 text-white px-4">
      <h1 className="text-4xl font-extrabold mb-4">Welcome to BounceShield 🛡️</h1>
      <p className="text-lg mb-6">Your trusted shield for secure operations.</p>

      {showLogin ? (
        <div className="w-full max-w-md">
          <AuthForm setToken={setToken} />
        </div>
      ) : (
        <>
          <button
            onClick={() => setShowLogin(true)}
            className="bg-white text-blue-600 px-6 py-3 rounded-full shadow-lg hover:bg-gray-100 transition transform hover:scale-105 mb-4"
          >
            Login
          </button>
          <EmailVerifyLite onForceLogin={() => setShowLogin(true)} />
        </>
      )}
    </div>
  );
};

export default Welcome;
