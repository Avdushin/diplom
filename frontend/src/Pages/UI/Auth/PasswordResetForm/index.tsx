import React, { FC, FormEvent, useState } from "react";
import { useForm } from "@Hooks/useForm";
import useApi from "@Hooks/useAPI";

const PasswordResetForm: FC = () => {
  const { values, handleChange } = useForm({
    email: "",
  });

  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handlePasswordReset = async (event: FormEvent) => {
    event.preventDefault();

    try {
      const response = await useApi("forgot-password", "POST", values);
      console.log("Password reset link sent successfully!", response);
      setSuccessMessage("Password reset link sent successfully!");
    } catch (error) {
      console.error("Error sending password reset link:", error);
    }
  };

  return (
    <div style={{ display: "flex", justifyContent: "center", alignItems: "center", height: "100vh" }}>
      <form onSubmit={handlePasswordReset}>
        <div>
          <input
            type="text"
            name="email"
            value={values.email}
            onChange={handleChange}
            placeholder="Enter your email"
          />
        </div>
        <button type="submit">Send Reset Link</button>
      </form>

      {successMessage && <div>{successMessage}</div>}
    </div>
  );
};

export default PasswordResetForm;
