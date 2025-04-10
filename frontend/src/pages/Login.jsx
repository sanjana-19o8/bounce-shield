import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { loginUser } from '../services/api';
import AuthForm from '../components/AuthForm';

export default function Login({ setToken }) {
  const [form, setForm] = useState({ username: '', password: '' });
  const navigate = useNavigate();

  const handleLogin = async () => {
    try {
      const res = await loginUser(form.username, form.password);
      localStorage.setItem('token', res.data.token);
      setToken(res.data.token);
      navigate('/dashboard');
    } catch {
      alert('Login failed');
    }
  };

  return <AuthForm form={form} setForm={setForm} onSubmit={handleLogin} />;
}
