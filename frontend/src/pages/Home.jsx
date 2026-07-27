import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useDocumentTitle } from "../hooks/useDocumentTitle.js";
import { login } from "../services/auth.js";

export default function Home() {
  useDocumentTitle("LMSVR SM - Login");
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    document.documentElement.lang = "en";
  }, []);

  const handleSubmit = async (e) => {
  e.preventDefault();
  setError("");
  setLoading(true);

  try {
    const data = await login({
      useremail: email,
      password,
    });

    localStorage.setItem("auth_token", data.token);

    console.log("Login successful:", data);

    navigate("/admissions");
  } catch (err) {
    console.error(err);

    setError(
      err.response?.data?.message ||
      err.message ||
      "Login failed. Please try again."
    );
  } finally {
    setLoading(false);
  }
};
  return (
    <div className="min-h-screen w-full bg-gradient-to-br from-indigo-50 to-blue-50 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="bg-white rounded-2xl shadow-lg p-8">
          <div className="mb-8 text-center">
            <h1 className="text-2xl font-semibold text-slate-900 mb-2">LMSVR SM</h1>
            <p className="text-sm text-slate-600">Learning Management System</p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-slate-700 mb-1">
                Email Address
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Enter your email"
                required
                className="w-full px-4 py-2 rounded-lg border border-slate-300 bg-slate-50 text-slate-900 placeholder-slate-500 focus:outline-none focus:bg-white focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-slate-700 mb-1">
                Password
              </label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter your password"
                required
                className="w-full px-4 py-2 rounded-lg border border-slate-300 bg-slate-50 text-slate-900 placeholder-slate-500 focus:outline-none focus:bg-white focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
              />
            </div>

            {error && (
              <div className="p-3 rounded-lg bg-red-50 border border-red-200 text-red-700 text-sm">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full py-2 px-4 mt-6 bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-400 text-white font-medium rounded-lg transition duration-200"
            >
              {loading ? "Logging in..." : "Login"}
            </button>
          </form>

          <p className="mt-6 text-center text-sm text-slate-600">
            Don't have an account?{" "}
            <a href="#" className="text-indigo-600 hover:text-indigo-700 font-medium">
              Sign up
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}