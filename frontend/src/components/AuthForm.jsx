import React, { useState } from "react";
import axios from "axios";

const AuthForm = ({ setToken }) => {
  const [form, setForm] = useState({ username: "", password: "" });
  const [error, setError] = useState("");

  const handleLogin = async () => {
    try {
      const res = await axios.post("/api/login", form);
      console.log("LOGIN SUCCESS", res.data); // 👈 See what's coming
  
      const token = res.data.token;
      localStorage.setItem("token", token);
      setToken(token);
      setError("");
    } catch (err) {
      console.error("LOGIN FAILED", err.response?.data || err.message); // 👈 Useful debug info
      setError("Invalid credentials or server not reachable.");
    }
  };  

  return (
    <div className="jusitfy-center max-w-sm mx-auto mt-20 bg-white rounded-2xl shadow-md">
      <h2 className="text-2xl font-bold text-center mb-4">Login</h2>

      {error && <p className="text-red-600 text-sm mb-4 text-center">{error}</p>}

      <input
        type="text"
        placeholder="Username"
        value={form.username}
        onChange={(e) => setForm({ ...form, username: e.target.value })}
        className="w-full p-2 border rounded mb-3 focus:outline-none focus:ring focus:ring-blue-300"
      />

      <input
        type="password"
        placeholder="Password"
        value={form.password}
        onChange={(e) => setForm({ ...form, password: e.target.value })}
        className="w-full p-2 border rounded mb-4 focus:outline-none focus:ring focus:ring-blue-300"
      />

      <button
        onClick={handleLogin}
        className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
      >
        Sign In
      </button>
    </div>
  );
};

export default AuthForm;
