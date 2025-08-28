/* eslint-disable @typescript-eslint/no-explicit-any */
import '../styles/index.scss'
import { useState } from 'react'
import { loginUser, setToken } from '../api/api'
import { useNavigate } from 'react-router-dom';
import { useAuth } from "../context/AuthContext";

export default function Login() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const navigate = useNavigate();
    const { login } = useAuth();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const res = await loginUser({ email, password });
            const { access_token, user } = res.data;

            setToken(access_token);
            login(user, access_token); // 👈 update context + localStorage

            navigate("/todos");
        } catch (err: any) {
            setError(err.response?.data?.error || "Login failed");
        }
    };

    return (
        <div className="login__container">
            <h2 className="login__title">Login</h2>
            <form onSubmit={handleSubmit} className="login__form">
                <input
                    type="email"
                    placeholder="Email"
                    value={email}
                    onChange={e => setEmail(e.target.value)}
                />
                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                />
                <button type="submit" className="login__form--btn">Login</button>
            </form>
            {error && <p className="login__error">{error}</p>}
            <a href="/register" className="login__create">Does not have an account?</a>
        </div>
    )
}
