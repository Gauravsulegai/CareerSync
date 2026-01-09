import React from 'react';
import { Link } from 'react-router-dom';
import { Shield, Zap, Globe, ArrowRight, Terminal } from 'lucide-react';

function Landing() {
  return (
    <div className="min-h-screen bg-background text-primary selection:bg-white selection:text-black flex flex-col">
      
      {/* Navbar */}
      <nav className="w-full px-6 py-6 flex justify-between items-center max-w-7xl mx-auto">
        <div className="text-xl font-bold tracking-tighter flex items-center gap-2">
          <Terminal size={24} /> CAREERSYNC
        </div>
        <div className="flex gap-6 text-sm font-medium text-secondary">
          <Link to="/login" className="hover:text-white transition-colors">Login</Link>
          <Link to="/signup" className="px-4 py-2 bg-white text-black rounded-full hover:bg-gray-200 transition-colors">
            Get Started
          </Link>
        </div>
      </nav>

      {/* Hero Section */}
      <main className="flex-1 flex flex-col justify-center px-6 max-w-5xl mx-auto w-full pt-10 pb-20">
        
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-white/10 bg-white/5 w-fit mb-8 animate-fade-in">
          <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
          <span className="text-xs font-mono text-secondary tracking-wide">SYSTEM ONLINE</span>
        </div>

        <h1 className="text-5xl md:text-7xl font-bold tracking-tight leading-[1.1] mb-6 animate-fade-in">
          Break into <br />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-white to-gray-500">
            the fortress.
          </span>
        </h1>

        <p className="text-secondary text-lg md:text-xl max-w-2xl leading-relaxed mb-10 animate-fade-in">
          The referral network that connects top talent with verified insiders. 
          Skip the queue. Get the interview.
        </p>

        <div className="flex flex-col sm:flex-row gap-4 animate-fade-in">
          <Link to="/login" className="px-8 py-4 bg-white text-black font-bold rounded-lg hover:bg-gray-200 transition-all active:scale-95 flex items-center justify-center gap-2">
            Enter Portal <ArrowRight size={18} />
          </Link>
          <Link to="/signup" className="px-8 py-4 border border-white/10 text-white font-medium rounded-lg hover:bg-white/5 transition-all active:scale-95 text-center">
            Create Account
          </Link>
        </div>

      </main>

      {/* Features */}
      <section className="border-t border-white/10 bg-surface/30 py-20">
        <div className="max-w-7xl mx-auto px-6 grid grid-cols-1 md:grid-cols-3 gap-12">
          {[
            { icon: Shield, title: "Verified Insiders", desc: "Every employee is verified via valid work email protocols." },
            { icon: Zap, title: "Direct Uplink", desc: "Bypass ATS filters. Send your resume directly to the inbox." },
            { icon: Globe, title: "Global Network", desc: "Access referral pools for top tech giants worldwide." }
          ].map((item, i) => (
            <div key={i} className="group">
              <div className="w-12 h-12 bg-white/5 rounded-lg flex items-center justify-center mb-4 group-hover:bg-white/10 transition-colors">
                <item.icon className="text-white" size={24} />
              </div>
              <h3 className="text-lg font-bold mb-2">{item.title}</h3>
              <p className="text-secondary leading-relaxed">{item.desc}</p>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}

export default Landing;