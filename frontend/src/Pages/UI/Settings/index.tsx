import { useState } from "react";
import { Link } from "react-router-dom";
import Cookies from "universal-cookie";
import RadioToggler from "@Features/DarkMode/RadioToggler";
import CheckBoxToggler from "@Features/DarkMode/CheckBoxToggler";
import { useAuth } from "@Hooks/useAuth";
import { KBMap } from "./KBMap";
import "./settings.scss";

const Settings = () => {
  const { setIsAuthenticated } = useAuth();

  // const handleLogout = () => {
  //   setIsAuthenticated(false);
  //   cookies.remove("Authorization");
  //   setTimeout(() => {
  //     window.location.href = "/login";
  //   }, 0);
  // };

  const cookies = new Cookies();

  const [user, setUser] = useState(null)

  function logout() {
    setUser(null);
    cookies.remove("Authorization", { path: "/" });
    setIsAuthenticated(false);
  }



  return (
    <>
      <div className="settings">
        <div className="settings__title">
          <h1>Настройки</h1>
          <p>{user}</p>
        </div>
        <div className="settings__section">
          <div className="setion__item">
            <div className="item__title">
              <h2>Тема</h2>
              <div className="togglers-box">
                <CheckBoxToggler />
              </div>
            </div>
          </div>
          <div className="setion__item">
            <div className="item__title">
              <h2>Сочитания клавишь</h2>
              <div className="item-box">
                <div className="p-20">
                  <Link to="/kbmap">Keyboard Map</Link>
                </div>
              </div>
            </div>
          </div>
          <div className="setion__item">
            <div className="item__title">
              <h2>Аккаунт</h2>
              <div className="item-box">
                <div className="p-20">
                  <Link to="/profile">Настройки аккаунта</Link>
                </div>
                <button type="button" className="button" onClick={logout}>
                  Выйти
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default Settings;
