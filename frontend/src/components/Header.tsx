import { LuLogIn } from "react-icons/lu";
import "../styles/header.scss";

export default function Header() {
    return (
        <div className="header__container">
            <div className="header__content">
                <div className="header__title">
                    <a href="/" className="title">tomomo</a>
                </div>

                <div className="header__button">
                    <a href="/register" className="header__button--reg">Register</a>
                    <a href="/login" className="header__button--login">Login <LuLogIn/></a>
                </div>
            </div>
        </div>
    );
}
