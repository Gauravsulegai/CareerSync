import React, { useState, useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { useNavigate, Link } from 'react-router-dom';
import { ArrowLeft, Lock, Mail } from 'lucide-react';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login } = useContext(AuthContext); // Access the global brain
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    
    // Call the backend
    const result = await login(email, password);
    
    if (result.success) {
      navigate('/dashboard'); // Redirect to dashboard on success
    } else {
      setError(result.message); // Show error message
    }
  };

  return (
    <div className="min-h-screen bg-background text-primary flex items-center justify-center p-6">
      
      {/* Back Button */}
      <Link to="/" className="absolute top-8 left-8 text-secondary hover:text-white transition-colors flex items-center gap-2 tracking-widest text-xs">
        <ArrowLeft className="w-4 h-4" /> RETURN TO BASE
      </Link>

      {/* Login Card */}
      <div className="w-full max-w-md animate-fade-in-up">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold tracking-widest mb-2">ACCESS POINT</h1>
          <p className="text-secondary text-xs font-mono tracking-widest">ENTER CREDENTIALS TO PROCEED</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          
          {/* Email Input */}
          <div className="group relative">
            <Mail className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-secondary group-focus-within:text-white transition-colors" />
            <input
              type="email"
              placeholder="IDENTITY (EMAIL)"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-surface border border-white/10 px-12 py-4 text-sm tracking-widest focus:outline-none focus:border-white/40 transition-all placeholder:text-white/20"
              required
            />
          </div>

          {/* Password Input */}
          <div className="group relative">
            <Lock className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-secondary group-focus-within:text-white transition-colors" />
            <input
              type="password"
              placeholder="PASSPHRASE"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full bg-surface border border-white/10 px-12 py-4 text-sm tracking-widest focus:outline-none focus:border-white/40 transition-all placeholder:text-white/20"
              required
            />
          </div>

          {/* Error Message */}
          {error && (
            <div className="text-red-500 text-xs font-mono tracking-widest text-center animate-pulse">
              ERROR: {error.toUpperCase()}
            </div>
          )}

          {/* Submit Button */}
          <button
            type="submit"
            className="w-full bg-white text-black font-bold py-4 tracking-[0.2em] hover:bg-gray-200 transition-all active:scale-95"
          >
            AUTHENTICATE
          </button>

          <div className="text-center mt-6">
             <Link to="/signup" className="text-xs text-secondary hover:text-white underline decoration-white/20 underline-offset-4 tracking-widest">
                REQ_NEW_ACCOUNT
             </Link>
          </div>

        </form>
      </div>
    </div>
  );
};

export default Login;