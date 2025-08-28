/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { registerUser, setToken } from "../api/api";
import { useAuth } from "../context/AuthContext"; // 👈 import context
import "../styles/index.scss";

export default function Register() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [name, setName] = useState("");
    const [error, setError] = useState("");
    const navigate = useNavigate();
    const { login } = useAuth(); // 👈 ใช้ login() จาก context

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const res = await registerUser({ email, password, name });
            const { access_token, user } = res.data; // 👈 backend ควรส่ง user กลับมาด้วย

            setToken(access_token);
            login(user, access_token); // 👈 อัปเดต context + localStorage

            navigate("/todos");
        } catch (err: any) {
            setError(err.response?.data?.error || "Register failed");
        }
    };

    return (
        <div className="register__container">
            <h2 className="register__title">Register</h2>
            <form onSubmit={handleSubmit} className="register__form">
                <input
                    type="text"
                    placeholder="Name"
                    value={name}
                    onChange={e => setName(e.target.value)}
                />
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
                <button type="submit" className="register__form--btn">Register</button>
            </form>
            {error && <p className="register__error">{error}</p>}
            <a href="/login" className="register__already">Already have an account?</a>
        </div>
    );
}
