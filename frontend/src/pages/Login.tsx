/* eslint-disable @typescript-eslint/no-explicit-any */
import '../styles/index.scss'
import { useState } from 'react'
import { loginUser, setToken } from '../api/api'
import { useNavigate } from 'react-router-dom';

export default function Login() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const res = await loginUser({ email, password});
            const token = res.data.access_token;
            setToken(token);
            localStorage.setItem("token", token);
            navigate("/todos");
        } catch (err: any) {
            setError(err.response?.data?.error || "Login failed");
        }
    };
 
    return (
        <div>
            <h2>Login</h2>
            <form onSubmit={handleSubmit}>
                <input type="email" placeholder="Email" value={email} onChange={ e => setEmail(e.target.value)} />
                <input type="password" placeholder="Password" value={password} onChange={ e => setPassword(e.target.value)} />
                <button type="submit">Login</button>
            </form>
            {error && <p style={{ color: "red"}}>{error}</p>}
        </div>
    )
}
