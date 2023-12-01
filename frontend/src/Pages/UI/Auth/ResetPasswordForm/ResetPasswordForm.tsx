import React, { useState } from 'react';
import { useLocation } from 'react-router-dom';
import axios from 'axios';
import './resetPass.scss';

const ResetPassword: React.FC = () => {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [message, setMessage] = useState('');
  const location = useLocation();

  const resetPassword = async (token: string) => {
    try {
      const response = await axios.post(
        `${import.meta.env.VITE_REACT_API_URL}reset-password`,
        { token, newPassword: password }
      );
      setMessage(response.data.message);
    } catch (error) {
      setMessage('Произошла ошибка при сбросе пароля');
    }
  };

  const handleFormSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      setMessage('Пароли не совпадают');
      return;
    }

    const searchParams = new URLSearchParams(location.search);
    const token = searchParams.get('token');
    if (token) {
      resetPassword(token);
    } else {
      setMessage('Токен не найден в URL');
    }
  };

  return (
    <form onSubmit={handleFormSubmit} className='pass-reset'>
      <input
        type='password'
        placeholder='Новый пароль'
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <input
        type='password'
        placeholder='Подтвердите пароль'
        value={confirmPassword}
        onChange={(e) => setConfirmPassword(e.target.value)}
      />
      <button type='submit'>Сбросить пароль</button>
      {message && <p>{message}</p>}
    </form>
  );
};

export default ResetPassword;
