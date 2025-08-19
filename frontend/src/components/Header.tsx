import "../styles/header.scss";

export default function Header() {
    return (
        <div className="header-container">
            <div className="header">
                <div className="header-title">
                    <a href="/" className="title">canwejust</a>
                </div>

                <div className="header-button">
                    <a href="/register" className="reg-button">Register</a>
                    <a href="/login" className="login-button">Login</a>
                </div>
            </div>
        </div>
    );
}
