import React, { useState } from 'react';

interface SetupNewPasswordProps {
  onPasswordSetup: (password: string) => void;
}

const SetupNewPassword: React.FC<SetupNewPasswordProps> = ({ onPasswordSetup }) => {
  const [password, setPassword] = useState('');

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // Передаем пароль на сервер
    onPasswordSetup(password);
  };

  return (
    <div>
      <h2>Setup New Password</h2>
      <form onSubmit={handleSubmit}>
        <label>
          New Password:
          <input type="password" value={password} onChange={handlePasswordChange} />
        </label>
        <button type="submit">Set Password</button>
      </form>
    </div>
  );
};

export default SetupNewPassword;