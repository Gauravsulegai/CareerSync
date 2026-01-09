import React, { useState } from 'react';
import API from '../api';
import { useNavigate, Link } from 'react-router-dom';
import { User, Briefcase, Mail, Lock, Building, ArrowRight } from 'lucide-react';

const Signup = () => {
  const navigate = useNavigate();
  const [role, setRole] = useState('student'); // 'student' or 'employee'
  const [formData, setFormData] = useState({
    name: '', email: '', password: '', 
    company_name: '', work_email: '', position: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      // Prepare payload based on role
      const payload = {
        name: formData.name,
        email: formData.email,
        password: formData.password,
        role: role,
        // Only include these if employee
        ...(role === 'employee' && {
          company_name: formData.company_name,
          work_email: formData.work_email,
          position: formData.position
        })
      };

      await API.post('/signup', payload);
      // If successful, go to login
      navigate('/login'); 
    } catch (err) {
      setError(err.response?.data?.error || "Signup failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-6 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-gray-900 via-background to-background">
      
      <div className="w-full max-w-md animate-fade-in">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold tracking-tight mb-2">Initialize Protocol</h1>
          <p className="text-secondary text-sm">Create your identity to access the network.</p>
        </div>

        {/* Role Toggle */}
        <div className="grid grid-cols-2 gap-2 p-1 bg-surface rounded-lg mb-8 border border-white/10">
          <button 
            onClick={() => setRole('student')}
            className={`flex items-center justify-center gap-2 py-2 text-sm font-medium rounded-md transition-all ${role === 'student' ? 'bg-white text-black shadow-lg' : 'text-secondary hover:text-white'}`}
          >
            <User size={16} /> Student
          </button>
          <button 
            onClick={() => setRole('employee')}
            className={`flex items-center justify-center gap-2 py-2 text-sm font-medium rounded-md transition-all ${role === 'employee' ? 'bg-white text-black shadow-lg' : 'text-secondary hover:text-white'}`}
          >
            <Briefcase size={16} /> Employee
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          
          <div className="relative group">
            <User className="absolute left-3 top-3 text-secondary group-focus-within:text-white transition-colors" size={18} />
            <input 
              type="text" placeholder="Full Name" required
              className="w-full bg-surface border border-white/10 rounded-lg py-3 pl-10 pr-4 text-sm focus:outline-none focus:border-white/40 transition-all placeholder:text-gray-600"
              value={formData.name}
              onChange={e => setFormData({...formData, name: e.target.value})}
            />
          </div>

          <div className="relative group">
            <Mail className="absolute left-3 top-3 text-secondary group-focus-within:text-white transition-colors" size={18} />
            <input 
              type="email" placeholder="Personal Email" required
              className="w-full bg-surface border border-white/10 rounded-lg py-3 pl-10 pr-4 text-sm focus:outline-none focus:border-white/40 transition-all placeholder:text-gray-600"
              value={formData.email}
              onChange={e => setFormData({...formData, email: e.target.value})}
            />
          </div>

          <div className="relative group">
            <Lock className="absolute left-3 top-3 text-secondary group-focus-within:text-white transition-colors" size={18} />
            <input 
              type="password" placeholder="Password" required
              className="w-full bg-surface border border-white/10 rounded-lg py-3 pl-10 pr-4 text-sm focus:outline-none focus:border-white/40 transition-all placeholder:text-gray-600"
              value={formData.password}
              onChange={e => setFormData({...formData, password: e.target.value})}
            />
          </div>

          {/* Employee Specific Fields */}
          {role === 'employee' && (
            <div className="space-y-4 pt-4 border-t border-white/10 animate-fade-in">
              <p className="text-xs font-mono text-accent uppercase tracking-widest">Verification Details</p>
              
              <div className="relative group">
                <Building className="absolute left-3 top-3 text-secondary group-focus-within:text-white transition-colors" size={18} />
                <input 
                  type="text" placeholder="Company Name (e.g. Google)" required
                  className="w-full bg-surface border border-white/10 rounded-lg py-3 pl-10 pr-4 text-sm focus:outline-none focus:border-white/40 transition-all placeholder:text-gray-600"
                  value={formData.company_name}
                  onChange={e => setFormData({...formData, company_name: e.target.value})}
                />
              </div>

              <div className="relative group">
                <Mail className="absolute left-3 top-3 text-secondary group-focus-within:text-white transition-colors" size={18} />
                <input 
                  type="email" placeholder="Work Email (@company.com)" required
                  className="w-full bg-surface border border-white/10 rounded-lg py-3 pl-10 pr-4 text-sm focus:outline-none focus:border-white/40 transition-all placeholder:text-gray-600"
                  value={formData.work_email}
                  onChange={e => setFormData({...formData, work_email: e.target.value})}
                />
              </div>

               <div className="relative group">
                <Briefcase className="absolute left-3 top-3 text-secondary group-focus-within:text-white transition-colors" size={18} />
                <input 
                  type="text" placeholder="Job Position (e.g. Senior SDE)" required
                  className="w-full bg-surface border border-white/10 rounded-lg py-3 pl-10 pr-4 text-sm focus:outline-none focus:border-white/40 transition-all placeholder:text-gray-600"
                  value={formData.position}
                  onChange={e => setFormData({...formData, position: e.target.value})}
                />
              </div>
            </div>
          )}

          {error && <div className="text-red-400 text-xs text-center font-mono bg-red-500/10 py-2 rounded">{error}</div>}

          <button 
            type="submit" disabled={loading}
            className="w-full bg-white text-black font-bold py-3 rounded-lg hover:bg-gray-200 transition-all active:scale-95 mt-4 flex items-center justify-center gap-2"
          >
            {loading ? 'Processing...' : <>Create Account <ArrowRight size={18} /></>}
          </button>
        </form>

        <p className="text-center text-xs text-secondary mt-8">
          Already have an identity? <Link to="/login" className="text-white underline hover:text-accent">Login here</Link>
        </p>
      </div>
    </div>
  );
};

export default Signup;