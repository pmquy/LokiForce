import { Sparkles, Cpu, FolderKanban, Building2, Terminal } from 'lucide-react';
import { Link } from '@tanstack/react-router';

export function DashboardHome() {
  return (
    <div className="space-y-8 max-w-7xl mx-auto">
      {/* Welcome Banner */}
      <div className="relative rounded-3xl overflow-hidden bg-gradient-to-r from-slate-900 to-indigo-950/40 border border-slate-800 p-6 md:p-8">
        <div className="absolute top-0 right-0 w-80 h-80 bg-indigo-500/10 rounded-full blur-3xl -translate-y-12 translate-x-12" />
        <div className="relative z-10 max-w-2xl">
          <div className="flex items-center gap-2 text-indigo-400 text-sm font-bold tracking-wide uppercase">
            <Sparkles className="h-4.5 w-4.5 text-indigo-400 animate-spin" />
            Platform Engineering Portal
          </div>
          <h2 className="mt-3 text-3xl md:text-4xl font-extrabold text-white tracking-tight">
            Welcome to LokiForce IDP
          </h2>
          <p className="mt-4 text-base text-slate-400 leading-relaxed">
            LokiForce is an Internal Developer Platform that standardizes infrastructure, scaffolding, and services across all engineering teams. Deploy standard services with GitHub integration instantly.
          </p>
          <div className="mt-6 flex flex-wrap gap-4">
            <Link
              to="/services"
              className="inline-flex items-center justify-center px-5 py-3 border border-transparent text-sm font-bold rounded-2xl text-white bg-indigo-600 hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-600/25 cursor-pointer"
            >
              Go to Service Catalog
            </Link>
          </div>
        </div>
      </div>

      {/* Quick Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-slate-900 border border-slate-800 rounded-3xl p-6 relative overflow-hidden group">
          <div className="flex justify-between items-start">
            <div className="space-y-2">
              <span className="text-slate-400 text-sm font-semibold block">Organizations</span>
              <span className="text-3xl font-extrabold text-white">Active</span>
            </div>
            <div className="p-3 bg-slate-800 border border-slate-700 rounded-2xl text-indigo-400 group-hover:scale-110 transition-transform">
              <Building2 className="h-6 w-6" />
            </div>
          </div>
          <Link to="/organizations" className="mt-4 text-xs font-bold text-indigo-400 hover:text-indigo-300 transition-colors block cursor-pointer">
            Manage Organizations &rarr;
          </Link>
        </div>

        <div className="bg-slate-900 border border-slate-800 rounded-3xl p-6 relative overflow-hidden group">
          <div className="flex justify-between items-start">
            <div className="space-y-2">
              <span className="text-slate-400 text-sm font-semibold block">Projects</span>
              <span className="text-3xl font-extrabold text-white">Active</span>
            </div>
            <div className="p-3 bg-slate-800 border border-slate-700 rounded-2xl text-teal-400 group-hover:scale-110 transition-transform">
              <FolderKanban className="h-6 w-6" />
            </div>
          </div>
          <Link to="/projects" className="mt-4 text-xs font-bold text-teal-400 hover:text-teal-300 transition-colors block cursor-pointer">
            Manage Projects &rarr;
          </Link>
        </div>

        <div className="bg-slate-900 border border-slate-800 rounded-3xl p-6 relative overflow-hidden group">
          <div className="flex justify-between items-start">
            <div className="space-y-2">
              <span className="text-slate-400 text-sm font-semibold block">Services Catalog</span>
              <span className="text-3xl font-extrabold text-white">Scaffolder</span>
            </div>
            <div className="p-3 bg-slate-800 border border-slate-700 rounded-2xl text-indigo-400 group-hover:scale-110 transition-transform">
              <Cpu className="h-6 w-6" />
            </div>
          </div>
          <Link to="/services" className="mt-4 text-xs font-bold text-indigo-400 hover:text-indigo-300 transition-colors block cursor-pointer">
            Scaffold Service &rarr;
          </Link>
        </div>
      </div>

      {/* Developer Getting Started */}
      <div className="bg-slate-900 border border-slate-800 rounded-3xl p-6 md:p-8">
        <h3 className="text-xl font-bold text-white flex items-center gap-3">
          <Terminal className="h-5 w-5 text-teal-400" />
          Quick Start Guide for Developers
        </h3>
        <div className="mt-6 grid grid-cols-1 md:grid-cols-3 gap-6 text-sm">
          <div className="space-y-2">
            <div className="w-8 h-8 rounded-full bg-slate-800 flex items-center justify-center font-bold text-slate-300 border border-slate-700">1</div>
            <h4 className="font-bold text-slate-200">Register Organization</h4>
            <p className="text-slate-400 text-xs">Create a logical scope namespace for your team and services catalog.</p>
          </div>
          <div className="space-y-2">
            <div className="w-8 h-8 rounded-full bg-slate-800 flex items-center justify-center font-bold text-slate-300 border border-slate-700">2</div>
            <h4 className="font-bold text-slate-200">Define a Project</h4>
            <p className="text-slate-400 text-xs">Bundle your related Microservices together under a Project name.</p>
          </div>
          <div className="space-y-2">
            <div className="w-8 h-8 rounded-full bg-slate-800 flex items-center justify-center font-bold text-slate-300 border border-slate-700">3</div>
            <h4 className="font-bold text-slate-200">Scaffold Service</h4>
            <p className="text-slate-400 text-xs">Choose a template (Go/Node) and create a repository on GitHub instantly.</p>
          </div>
        </div>
      </div>
    </div>
  );
}
