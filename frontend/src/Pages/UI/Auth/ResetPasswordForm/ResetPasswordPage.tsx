import React from 'react';
import { useLocation } from 'react-router-dom';
import ResetPasswordForm from './ResetPasswordForm';

const ResetPasswordPage: React.FC = () => {
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const token = searchParams.get('token');

  if (!token) {
    return <div>Invalid token</div>;
  }

  return (
    <div>
      <ResetPasswordForm token={token} />
    </div>
  );
};

export default ResetPasswordPage;