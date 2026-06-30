import { useEffect, useState } from 'react';
import { Outlet, Link, useNavigate, useRouterState } from '@tanstack/react-router';
import { Shield, LayoutDashboard, Building2, FolderKanban, Cpu, LogOut, Menu, X, User } from 'lucide-react';
import { apiRequest } from '../services/api';

interface UserProfile {
  ID: string;
  Username: string;
  Email: string;
  Role: string;
}

export function DashboardLayout() {
  const navigate = useNavigate();
  const routerState = useRouterState();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('lokiforce_token');
    if (!token) {
      navigate({ to: '/login' });
      return;
    }

    // Fetch user profile
    apiRequest<UserProfile>('/users/profile')
      .then(setProfile)
      .catch(() => {
        localStorage.removeItem('lokiforce_token');
        navigate({ to: '/login' });
      });
  }, [navigate]);

  const handleLogout = () => {
    localStorage.removeItem('lokiforce_token');
    navigate({ to: '/login' });
  };

  const navLinks = [
    { name: 'Dashboard', to: '/', icon: LayoutDashboard },
    { name: 'Organizations', to: '/organizations', icon: Building2 },
    { name: 'Projects', to: '/projects', icon: FolderKanban },
    { name: 'Services', to: '/services', icon: Cpu },
  ];

  if (!profile) {
    return (
      <div className="min-h-screen bg-slate-950 flex items-center justify-center">
        <div className="w-10 h-10 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex font-sans">
      {/* Sidebar - Desktop */}
      <aside className="hidden md:flex flex-col w-64 bg-slate-900 border-r border-slate-800">
        <div className="h-16 flex items-center gap-3 px-6 border-b border-slate-800">
          <div className="p-2 bg-gradient-to-tr from-indigo-500 to-teal-500 rounded-xl">
            <Shield className="h-6 w-6 text-white" />
          </div>
          <span className="font-extrabold text-xl bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent">
            LokiForce
          </span>
        </div>

        <nav className="flex-1 px-4 py-6 space-y-1.5">
          {navLinks.map((link) => {
            const Icon = link.icon;
            const isActive = routerState.location.pathname === link.to;
            return (
              <Link
                key={link.name}
                to={link.to}
                className={`flex items-center gap-3 px-4 py-3 rounded-2xl text-sm font-semibold transition-all ${
                  isActive
                    ? 'bg-gradient-to-r from-indigo-500/20 to-teal-500/10 text-white border border-indigo-500/20'
                    : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200'
                }`}
              >
                <Icon className={`h-5 w-5 ${isActive ? 'text-indigo-400' : 'text-slate-500'}`} />
                {link.name}
              </Link>
            );
          })}
        </nav>

        <div className="p-4 border-t border-slate-800 space-y-3">
          <div className="flex items-center gap-3 px-2 py-1.5">
            <div className="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center text-slate-300 font-bold border border-slate-700">
              {profile.Username.substring(0, 2).toUpperCase()}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-bold text-slate-200 truncate">{profile.Username}</p>
              <p className="text-xs text-slate-500 truncate">{profile.Email}</p>
            </div>
          </div>
          <button
            onClick={handleLogout}
            className="w-full flex items-center justify-center gap-2.5 px-4 py-2.5 border border-slate-800 hover:border-rose-500/30 rounded-2xl text-sm font-bold text-slate-400 hover:text-rose-400 hover:bg-rose-500/5 transition-all cursor-pointer"
          >
            <LogOut className="h-4.5 w-4.5" />
            Logout
          </button>
        </div>
      </aside>

      {/* Mobile Sidebar Overlay */}
      {sidebarOpen && (
        <div className="fixed inset-0 bg-slate-950/60 backdrop-blur-sm z-40 md:hidden" onClick={() => setSidebarOpen(false)} />
      )}

      {/* Sidebar - Mobile */}
      <aside
        className={`fixed inset-y-0 left-0 w-64 bg-slate-900 border-r border-slate-800 z-50 transform md:hidden transition-transform duration-300 ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        <div className="h-16 flex items-center justify-between px-6 border-b border-slate-800">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-gradient-to-tr from-indigo-500 to-teal-500 rounded-xl">
              <Shield className="h-6 w-6 text-white" />
            </div>
            <span className="font-extrabold text-xl text-white">LokiForce</span>
          </div>
          <button onClick={() => setSidebarOpen(false)} className="text-slate-400 hover:text-slate-200 cursor-pointer">
            <X className="h-6 w-6" />
          </button>
        </div>

        <nav className="flex-1 px-4 py-6 space-y-1.5">
          {navLinks.map((link) => {
            const Icon = link.icon;
            const isActive = routerState.location.pathname === link.to;
            return (
              <Link
                key={link.name}
                to={link.to}
                onClick={() => setSidebarOpen(false)}
                className={`flex items-center gap-3 px-4 py-3 rounded-2xl text-sm font-semibold transition-all ${
                  isActive
                    ? 'bg-gradient-to-r from-indigo-500/20 to-teal-500/10 text-white border border-indigo-500/20'
                    : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200'
                }`}
              >
                <Icon className="h-5 w-5" />
                {link.name}
              </Link>
            );
          })}
        </nav>

        <div className="p-4 border-t border-slate-800 space-y-3 absolute bottom-0 w-full bg-slate-900">
          <div className="flex items-center gap-3 px-2 py-1.5">
            <div className="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center text-slate-300 font-bold border border-slate-700">
              {profile.Username.substring(0, 2).toUpperCase()}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-bold text-slate-200 truncate">{profile.Username}</p>
              <p className="text-xs text-slate-500 truncate">{profile.Email}</p>
            </div>
          </div>
          <button
            onClick={handleLogout}
            className="w-full flex items-center justify-center gap-2.5 px-4 py-2.5 border border-slate-800 hover:border-rose-500/30 rounded-2xl text-sm font-bold text-slate-400 hover:text-rose-400 hover:bg-rose-500/5 transition-all cursor-pointer"
          >
            <LogOut className="h-4.5 w-4.5" />
            Logout
          </button>
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col min-w-0 relative">
        <header className="h-16 bg-slate-900/40 backdrop-blur-md border-b border-slate-800 flex items-center justify-between px-6 z-10">
          <div className="flex items-center gap-3">
            <button
              onClick={() => setSidebarOpen(true)}
              className="md:hidden p-2 -ml-2 text-slate-400 hover:text-slate-200 cursor-pointer"
            >
              <Menu className="h-6 w-6" />
            </button>
            <h1 className="text-lg font-bold text-white tracking-wide">
              {navLinks.find((l) => l.to === routerState.location.pathname)?.name || 'Dashboard'}
            </h1>
          </div>

          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 px-3 py-1.5 bg-slate-900 border border-slate-800 rounded-full text-xs font-semibold text-slate-300">
              <User className="h-4 w-4 text-indigo-400" />
              <span>{profile.Role.toUpperCase()}</span>
            </div>
          </div>
        </header>

        <main className="flex-1 overflow-y-auto p-6 md:p-8 bg-slate-950">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
