import { FC, FormEvent, useState } from "react";
import { useForm } from "@Hooks/useForm";
import useApi from "@Hooks/useAPI";
import { useNavigate } from "react-router-dom";

const PasswordResetForm: FC = () => {
  const { values, handleChange } = useForm({
    email: "",
  });
  const navigate = useNavigate(); // Хук для навигации

  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handlePasswordReset = async (event: FormEvent) => {
    event.preventDefault();

    try {
      const response = await useApi("forgot-password", "POST", values);
      console.log("Password reset link sent successfully!", response);
      setSuccessMessage("Password reset link sent successfully!");

      // Переходим на новый маршрут после успешной отправки запроса
      navigate("/setup-new-password");
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
