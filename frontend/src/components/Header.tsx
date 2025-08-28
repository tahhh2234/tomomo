import { useAuth } from "../context/AuthContext";
import { LuLogIn, LuLogOut } from "react-icons/lu";
import "../styles/index.scss";

export default function Header() {
  const { user, logout } = useAuth();

  return (
    <div className="header__container">
      <div className="header__content">
        <div className="header__title">
          <a href="/" className="title">tomomo</a>
        </div>

        <div className="header__button">
          {user ? (
            <>
              <span className="header__username">{user.name}</span>
              <button onClick={logout} className="header__button--logout">
                Logout <LuLogOut />
              </button>
            </>
          ) : (
            <>
              <a href="/register" className="header__button--reg">Register</a>
              <a href="/login" className="header__button--login">
                Login <LuLogIn />
              </a>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
