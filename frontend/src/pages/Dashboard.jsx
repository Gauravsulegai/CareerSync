import React, { useContext, useState, useEffect } from 'react';
import { AuthContext } from '../context/AuthContext';
import API from '../api';
import { Search, Building, LogOut, X, Send, AlertTriangle, ExternalLink, Link as LinkIcon, Copy, Check } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const Dashboard = () => {
  const { user, logout } = useContext(AuthContext);
  const navigate = useNavigate();
  const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);

  const handleLogoutClick = () => {
    setShowLogoutConfirm(true);
  };

  const confirmLogout = () => {
    logout();
    navigate('/');
  };

  if (!user) return <div className="p-10 text-center animate-pulse tracking-widest text-xs">LOADING_IDENTITY...</div>;

  return (
    <div className="min-h-screen bg-background text-primary p-6 md:p-12 relative">
      {/* Header */}
      <div className="flex justify-between items-center mb-12 pb-6 border-b border-white/5">
        <div>
          <h1 className="text-3xl font-bold tracking-widest">Dashboard</h1>
          <p className="text-secondary text-[10px] font-mono mt-2 tracking-widest">
            Name: <span className="text-white">{user.name.toUpperCase()}</span> ||
            Role: <span className="text-white">{user.role.toUpperCase()}</span>
          </p>
        </div>
        <button 
          onClick={handleLogoutClick} 
          className="group flex items-center gap-2 text-[10px] font-bold text-zinc-500 hover:text-red-400 tracking-widest transition-colors px-4 py-2 rounded-full hover:bg-red-500/5 border border-transparent hover:border-red-500/20"
        >
          <LogOut className="w-3 h-3 group-hover:-translate-x-1 transition-transform" /> LOGOUT
        </button>
      </div>

      {user.role === 'student' ? <StudentView user={user} /> : <EmployeeView />}

      {/* Logout Modal */}
      {showLogoutConfirm && (
        <div className="fixed inset-0 bg-black/90 backdrop-blur-sm flex items-center justify-center z-50 p-4 animate-fade-in-up">
          <div className="bg-zinc-900/50 backdrop-blur-xl border border-white/10 p-8 max-w-sm w-full rounded-2xl text-center shadow-2xl">
            <div className="flex justify-center mb-6">
              <div className="p-4 bg-red-500/10 rounded-full border border-red-500/20 shadow-[0_0_20px_rgba(239,68,68,0.2)]">
                <AlertTriangle className="w-8 h-8 text-red-500" />
              </div>
            </div>
            <h3 className="text-lg font-bold mb-2 tracking-widest text-white">TERMINATE SESSION?</h3>
            <p className="text-zinc-400 text-xs font-mono mb-8 leading-relaxed">
              CONFIRM DISCONNECTION FROM THE SECURE NETWORK.
            </p>
            <div className="flex gap-3">
              <button onClick={() => setShowLogoutConfirm(false)} className="flex-1 py-3 border border-white/10 hover:bg-white/5 text-xs font-bold tracking-widest transition-colors rounded-xl text-zinc-300">CANCEL</button>
              <button onClick={confirmLogout} className="flex-1 py-3 bg-red-600 hover:bg-red-500 text-white text-xs font-bold tracking-widest transition-colors rounded-xl shadow-lg shadow-red-900/20">CONFIRM</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

// --- HELPER: Handles Job ID (Copy) vs Job Link (Click + Copy) ---
const JobIdentifier = ({ value }) => {
    const [copied, setCopied] = useState(false);
    const isUrl = value.startsWith('http');

    const handleCopy = () => {
        navigator.clipboard.writeText(value);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <div className="flex items-center gap-3 mb-3">
            <span className="text-zinc-500 text-xs font-mono tracking-wide uppercase">
                {isUrl ? 'TARGET:' : 'JOB ID:'}
            </span>
            
            <div className="flex items-center gap-2 group/link">
                {/* The actual link/ID text */}
                <span 
                    className="text-yellow-400 font-bold text-sm tracking-wide truncate max-w-[200px] md:max-w-[300px] cursor-default" 
                    title={value}
                >
                    {value}
                </span>

                {/* Copy Button */}
                <button 
                    onClick={handleCopy}
                    className="p-1.5 rounded-md hover:bg-white/10 text-zinc-500 hover:text-white transition-all active:scale-95"
                    title="Copy to Clipboard"
                >
                    {copied ? <Check className="w-3.5 h-3.5 text-green-500" /> : <Copy className="w-3.5 h-3.5" />}
                </button>

                {/* External Link Icon (Only if URL) - Optional but helpful */}
                {isUrl && (
                    <a 
                        href={value} 
                        target="_blank" 
                        rel="noreferrer" 
                        className="p-1.5 rounded-md hover:bg-white/10 text-zinc-500 hover:text-blue-400 transition-all opacity-0 group-hover/link:opacity-100"
                        title="Open Link in New Tab"
                    >
                        <ExternalLink className="w-3.5 h-3.5" />
                    </a>
                )}
            </div>
        </div>
    );
};

// --- SUB COMPONENT: STUDENT VIEW ---
const StudentView = ({ user }) => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [searching, setSearching] = useState(false);
  const [selectedCompany, setSelectedCompany] = useState(null); 

  useEffect(() => {
    if (!query.trim()) {
      setResults([]);
      setSearching(false);
      return;
    }
    setSearching(true);
    const delayDebounceFn = setTimeout(async () => {
      try {
        const { data } = await API.get(`/companies/search?query=${query}`);
        setResults(data);
      } catch (err) {
        console.error("Search failed", err);
      } finally {
        setSearching(false);
      }
    }, 500);
    return () => clearTimeout(delayDebounceFn);
  }, [query]);

  return (
    <div className="max-w-4xl mx-auto">
      <div className="glass-panel p-8 rounded-3xl mb-12 shadow-2xl relative overflow-hidden group">
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 blur-[80px] rounded-full pointer-events-none -translate-y-1/2 translate-x-1/2"></div>
        <h2 className="text-lg font-bold mb-6 flex items-center gap-3 tracking-widest text-white relative z-10">
          <Search className="w-5 h-5 text-zinc-400" /> TARGET ACQUISITION
        </h2>
        <form onSubmit={(e) => e.preventDefault()} className="flex gap-4 relative z-10">
          <div className="relative flex-1 group/input">
            <input 
              type="text" 
              placeholder="ENTER COMPANY NAME (e.g. Google)" 
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              className="w-full bg-black/40 border border-white/10 p-5 rounded-2xl focus:outline-none focus:border-white/30 focus:ring-1 focus:ring-white/20 transition-all tracking-wide text-white placeholder-zinc-600 shadow-inner"
            />
          </div>
          <button type="submit" className="bg-white text-black px-8 rounded-2xl font-bold tracking-widest hover:bg-zinc-200 transition-all shadow-[0_0_20px_rgba(255,255,255,0.1)] active:scale-95">
            {searching ? 'SCANNING...' : 'SCAN'}
          </button>
        </form>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {results.map((company) => (
          <div key={company.id} className="p-6 border border-white/5 bg-zinc-900/30 hover:bg-zinc-900/50 hover:border-white/20 transition-all duration-300 rounded-2xl animate-fade-in-up group shadow-lg backdrop-blur-sm">
            <div className="flex items-start justify-between mb-8">
              <div>
                <h3 className="text-xl font-bold text-white tracking-wide group-hover:text-white transition-colors">{company.name}</h3>
                <p className="text-[10px] text-zinc-500 font-mono mt-1 tracking-widest uppercase">{company.domain}</p>
              </div>
              <div className="p-3 bg-white/5 rounded-xl text-zinc-400 group-hover:text-white group-hover:bg-white/10 transition-all">
                <Building className="w-6 h-6" />
              </div>
            </div>
            <div className="flex justify-end">
               <button onClick={() => setSelectedCompany(company)} className="text-[10px] font-bold border border-white/10 bg-black/20 hover:bg-white hover:text-black hover:border-white text-zinc-300 px-6 py-3 rounded-xl transition-all tracking-widest uppercase shadow-sm">
                 REQ_REFERRAL
               </button>
            </div>
          </div>
        ))}
        {results.length === 0 && !searching && query && (
          <div className="col-span-1 md:col-span-2 text-center py-16 border border-dashed border-white/10 rounded-3xl bg-white/5">
             <p className="text-zinc-500 font-mono text-xs tracking-widest">NO TARGETS FOUND IN DATABASE.</p>
          </div>
        )}
      </div>

      {selectedCompany && <ReferralModal company={selectedCompany} user={user} onClose={() => setSelectedCompany(null)} />}
    </div>
  );
};

// --- SUB COMPONENT: REFERRAL MODAL ---
const ReferralModal = ({ company, user, onClose }) => {
  const [formData, setFormData] = useState({
    first_name: user.name.split(' ')[0] || '', last_name: user.name.split(' ')[1] || '',
    email: user.email || '', mobile: '', linkedin_url: '', resume_url: '', job_link: '', motivation: ''
  });
  const [status, setStatus] = useState('IDLE'); 

  const handleSubmit = async (e) => {
    e.preventDefault();
    setStatus('SENDING');
    try {
      const payload = { employee_id: 2, ...formData }; // NOTE: Update ID if needed
      await API.post('/request/referral', payload);
      setStatus('SUCCESS');
      setTimeout(onClose, 2000); 
    } catch (err) {
      console.error(err);
      setStatus('ERROR');
    }
  };

  return (
    <div className="fixed inset-0 bg-black/90 backdrop-blur-md flex items-center justify-center z-50 p-4">
      <div className="bg-black border border-white/10 p-8 md:p-10 max-w-lg w-full relative animate-fade-in-up rounded-3xl shadow-2xl">
        <button onClick={onClose} className="absolute top-6 right-6 text-zinc-500 hover:text-red-500 transition-colors p-2 hover:bg-red-500/10 rounded-full"><X className="w-5 h-5" /></button>
        <div className="mb-8">
            <h3 className="text-2xl font-bold mb-2 tracking-wide">INITIATE UPLINK</h3>
            <div className="flex items-center gap-2">
                <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                <p className="text-[10px] text-zinc-500 font-mono tracking-widest uppercase">TARGET: <span className="text-white font-bold">{company.name}</span></p>
            </div>
        </div>
        {status === 'SUCCESS' ? (
          <div className="text-green-400 text-center py-16 font-mono border border-green-500/20 bg-green-500/5 rounded-2xl flex flex-col items-center justify-center gap-4">
            <Send className="w-10 h-10 mb-2" />
            <div>TRANSMISSION SUCCESSFUL. <br/> <span className="text-xs text-green-600/70 mt-2 block">STAND BY FOR RESPONSE.</span></div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <input placeholder="FIRST NAME" className="bg-zinc-900/50 border border-white/10 p-4 rounded-xl text-sm focus:border-white/40 focus:bg-zinc-900 outline-none transition-all placeholder:text-zinc-600" value={formData.first_name} onChange={e => setFormData({...formData, first_name: e.target.value})} required />
              <input placeholder="LAST NAME" className="bg-zinc-900/50 border border-white/10 p-4 rounded-xl text-sm focus:border-white/40 focus:bg-zinc-900 outline-none transition-all placeholder:text-zinc-600" value={formData.last_name} onChange={e => setFormData({...formData, last_name: e.target.value})} required />
            </div>
            <input placeholder="EMAIL CONTACT" className="w-full bg-zinc-900/50 border border-white/10 p-4 rounded-xl text-sm focus:border-white/40 focus:bg-zinc-900 outline-none transition-all placeholder:text-zinc-600" value={formData.email} onChange={e => setFormData({...formData, email: e.target.value})} required />
            <input placeholder="MOBILE NUMBER" className="w-full bg-zinc-900/50 border border-white/10 p-4 rounded-xl text-sm focus:border-white/40 focus:bg-zinc-900 outline-none transition-all placeholder:text-zinc-600" value={formData.mobile} onChange={e => setFormData({...formData, mobile: e.target.value})} required />
            <input placeholder="LINKEDIN URL" className="w-full bg-zinc-900/50 border border-white/10 p-4 rounded-xl text-sm focus:border-white/40 focus:bg-zinc-900 outline-none transition-all placeholder:text-zinc-600" value={formData.linkedin_url} onChange={e => setFormData({...formData, linkedin_url: e.target.value})} required />
            <input placeholder="RESUME URL (Drive/Dropbox)" className="w-full bg-zinc-900/50 border border-white/10 p-4 rounded-xl text-sm focus:border-white/40 focus:bg-zinc-900 outline-none transition-all placeholder:text-zinc-600" value={formData.resume_url} onChange={e => setFormData({...formData, resume_url: e.target.value})} required />
            <div className="relative group">
                <LinkIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-500 group-focus-within:text-white transition-colors" />
                <input placeholder="JOB ID / LINK (Required)" className="w-full bg-zinc-900/50 border border-white/10 p-4 pl-12 rounded-xl text-sm focus:border-white/40 focus:bg-zinc-900 outline-none transition-all placeholder:text-zinc-600" value={formData.job_link} onChange={e => setFormData({...formData, job_link: e.target.value})} required />
            </div>
            <textarea placeholder="Describe yourself (Max 500 words)..." className="w-full bg-zinc-900/50 border border-white/10 p-4 rounded-xl text-sm focus:border-white/40 focus:bg-zinc-900 outline-none h-32 resize-none transition-all placeholder:text-zinc-600" value={formData.motivation} onChange={e => setFormData({...formData, motivation: e.target.value})} required />
            <button disabled={status === 'SENDING'} className="w-full bg-white text-black font-bold py-4 rounded-xl mt-4 hover:bg-zinc-200 flex justify-center items-center gap-2 transition-all active:scale-95 shadow-[0_0_20px_rgba(255,255,255,0.15)]">{status === 'SENDING' ? 'TRANSMITTING...' : <><Send className="w-4 h-4" /> SEND REQUEST</>}</button>
            {status === 'ERROR' && <p className="text-red-500 text-xs text-center mt-2 font-mono">TRANSMISSION FAILED. CHECK CONNECTION.</p>}
          </form>
        )}
      </div>
    </div>
  );
};

// --- SUB COMPONENT: EMPLOYEE VIEW ---
const EmployeeView = () => {
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchRequests = async () => {
        try {
          const { data } = await API.get('/requests');
          setRequests(data);
        } catch (err) { console.error("Failed to load inbox", err); } 
        finally { setLoading(false); }
    };
    fetchRequests();
  }, []);

  const handleStatusUpdate = async (id, newStatus) => {
    try {
      await API.put(`/request/${id}/status`, { status: newStatus });
      // Quick Optimistic Update
      setRequests(requests.map(req => req.id === id ? {...req, status: newStatus} : req));
    } catch (err) { alert("Failed to update status"); }
  };

  if (loading) return <div className="p-8 text-center text-zinc-500 animate-pulse tracking-widest text-xs">DECRYPTING INBOX...</div>;

  return (
    <div className="max-w-5xl mx-auto">
      <h2 className="text-lg font-bold mb-8 tracking-widest flex items-center gap-3">
        <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
        INCOMING TRANSMISSIONS ({requests.length})
      </h2>

      <div className="space-y-4">
        {requests.map((req) => (
          <div key={req.id} className={`bg-zinc-900/40 border p-8 rounded-3xl relative overflow-hidden group transition-all backdrop-blur-sm
            ${req.status === 'Accepted' ? 'border-green-500/30 shadow-[0_0_30px_rgba(34,197,94,0.1)]' : 'border-white/5'} 
            ${req.status === 'Rejected' ? 'border-red-500/30' : ''}
          `}>
            
            <div className="flex flex-col md:flex-row gap-8">
              {/* Left: Candidate Info */}
              <div className="flex-1">
                <h3 className="text-xl font-bold text-white mb-2">{req.first_name} {req.last_name}</h3>
                
                {/* ✨ UPDATED: Clean text display with Copy Icon */}
                {req.job_link && <JobIdentifier value={req.job_link} />}

                <div className="text-xs text-zinc-400 font-mono space-y-3 mt-4">
                  <div className="flex items-center gap-2">
                    <span className="text-zinc-600">EMAIL:</span> 
                    <span className="text-white select-all">{req.email}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-zinc-600">MOBILE:</span> 
                    <span className="text-white select-all">{req.mobile}</span>
                  </div>
                  <div className="flex gap-4 mt-6">
                    <a href={req.linkedin_url} target="_blank" rel="noreferrer" className="text-blue-400 hover:text-blue-300 hover:underline underline-offset-4 transition-colors">LINKEDIN</a>
                    <a href={req.resume_url} target="_blank" rel="noreferrer" className="text-blue-400 hover:text-blue-300 hover:underline underline-offset-4 transition-colors">RESUME</a>
                  </div>
                </div>
              </div>

              {/* Right: Motivation, Badge & Actions */}
              <div className="flex-1 border-l border-white/5 pl-8 border-dashed flex flex-col">
                <div className={`self-end mb-4 px-4 py-1.5 text-[10px] font-bold tracking-widest rounded-full border
                  ${req.status === 'Pending' ? 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20' : ''}
                  ${req.status === 'Accepted' ? 'bg-green-500/10 text-green-500 border-green-500/20' : ''} 
                  ${req.status === 'Rejected' ? 'bg-red-500/10 text-red-500 border-red-500/20' : ''}
                `}>
                  {req.status.toUpperCase()}
                </div>

                <p className="text-sm italic text-zinc-400 mb-8 leading-relaxed">"{req.motivation}"</p>
                
                {req.status === 'Pending' && (
                  <div className="flex gap-3">
                    <button onClick={() => handleStatusUpdate(req.id, 'Accepted')} className="bg-green-600 hover:bg-green-500 text-white px-6 py-3 rounded-xl text-xs font-bold tracking-widest transition-all shadow-lg shadow-green-900/20 hover:scale-105 active:scale-95">ACCEPT</button>
                    <button onClick={() => handleStatusUpdate(req.id, 'Rejected')} className="border border-red-500/50 text-red-500 hover:bg-red-500 hover:text-white px-6 py-3 rounded-xl text-xs font-bold tracking-widest transition-all hover:scale-105 active:scale-95">REJECT</button>
                  </div>
                )}
                
                {req.status === 'Accepted' && (
                  <div className="flex items-center gap-2 text-green-500 text-xs font-bold tracking-widest bg-green-500/10 p-3 rounded-xl border border-green-500/20 inline-block">✓ CANDIDATE ENDORSED</div>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Dashboard;