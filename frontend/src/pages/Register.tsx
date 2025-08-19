import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { registerUser, setToken } from "../api/api";

export default function Register() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [name, setName] = useState("");
    const [error, setError] = useState("");
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const res = await registerUser({ email, password, name});
            const token = res.data.access_token;
            setToken(token);
            localStorage.setItem("token", token);
            navigate("/todos");
        } catch (err: any) {
            setError(err.response?.data?.error || "Register failed")
        }
    }

    return (
        <div>
            <h2>Register</h2>
            <form onSubmit={handleSubmit}>
                <input type="text" placeholder="Name" value={name} onChange={ e => setName(e.target.value)} />
                <input type="email" placeholder="Email" value={email} onChange={ e => setEmail(e.target.value)} />
                <input type="password" placeholder="Password" value={password} onChange={ e => setPassword(e.target.value)} />
                <button type="submit">Register</button>
            </form>
            {error && <p style={{ color: "red"}}>{error}</p>}
        </div>
    );
}