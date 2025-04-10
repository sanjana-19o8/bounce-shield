import React, { useState } from "react";
import toast from "react-hot-toast";

const EmailVerifyLite = ({ onForceLogin }) => {
  const [email, setEmail] = useState("");

  const handleVerify = () => {
    const attemptCount = parseInt(localStorage.getItem("emailAttemptCount") || "0");

    if (attemptCount >= 3) {
      onForceLogin();
      return;
    }

    localStorage.setItem("emailAttemptCount", attemptCount + 1);

    if (!/\S+@\S+\.\S+/.test(email)) {
      toast.error("Invalid email format");
    } else {
      toast.success(`Email ${email} verification attempted`);
    }

    setEmail("");
  };

  return (
    <div className="mt-6 bg-white p-4 rounded-xl shadow-md w-full max-w-md">
      <h3 className="text-lg font-semibold mb-2 text-gray-800">Quick Email Check</h3>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        className="w-full p-2 border rounded mb-2"
        placeholder="Enter email to verify"
      />
      <button
        onClick={handleVerify}
        className="bg-blue-600 text-white px-4 py-2 rounded w-full hover:bg-blue-700"
      >
        Verify
      </button>
      <p className="text-sm text-gray-500 mt-2">
        Only 3 attempts allowed without login (valid or invalid).
      </p>
    </div>
  );
};

export default EmailVerifyLite;
