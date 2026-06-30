import React, { useState, useEffect } from 'react';
import { useOrganizationsQuery } from '../../organizations/hooks/useOrganizations';
import { useProjectsQuery, useCreateProjectMutation } from '../hooks/useProjects';
import { FolderKanban, Plus, Loader2, X, AlertCircle, ChevronDown, Building2 } from 'lucide-react';

export function ProjectsList() {
  const [selectedOrgId, setSelectedOrgId] = useState<string>('');
  const [modalOpen, setModalOpen] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [errorMsg, setErrorMsg] = useState('');

  // 1. Fetch organizations
  const { data: orgs, isLoading: orgsLoading } = useOrganizationsQuery();

  // Set default organization once loaded
  useEffect(() => {
    if (orgs && orgs.length > 0 && !selectedOrgId) {
      setSelectedOrgId(orgs[0].ID);
    }
  }, [orgs, selectedOrgId]);

  // 2. Fetch projects under the selected organization
  const { data: projects, isLoading: projectsLoading } = useProjectsQuery(selectedOrgId);

  // Create project mutation
  const createProjMutation = useCreateProjectMutation(
    selectedOrgId,
    () => {
      setModalOpen(false);
      setName('');
      setDescription('');
    },
    (err: any) => {
      setErrorMsg(err.message || 'Failed to create project');
    }
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg('');
    createProjMutation.mutate({ name, description });
  };

  return (
    <div className="space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4">
        <div>
          <h2 className="text-2xl font-extrabold text-white tracking-tight">Projects</h2>
          <p className="text-sm text-slate-400">Manage code repositories, services, and environments grouped under projects.</p>
        </div>

        <div className="flex items-center gap-3">
          {orgs && orgs.length > 0 && (
            <div className="relative">
              <select
                value={selectedOrgId}
                onChange={(e) => setSelectedOrgId(e.target.value)}
                className="appearance-none bg-slate-900 border border-slate-800 rounded-2xl pl-4 pr-10 py-2.5 text-sm font-semibold text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all cursor-pointer min-w-[200px]"
              >
                {orgs.map((org) => (
                  <option key={org.ID} value={org.ID}>
                    {org.Name}
                  </option>
                ))}
              </select>
              <ChevronDown className="absolute right-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-500 pointer-events-none" />
            </div>
          )}

          {selectedOrgId && (
            <button
              onClick={() => setModalOpen(true)}
              className="flex items-center gap-2 px-4 py-2.5 rounded-2xl bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-sm transition-all cursor-pointer shadow-lg shadow-indigo-600/25"
            >
              <Plus className="h-4.5 w-4.5" />
              New Project
            </button>
          )}
        </div>
      </div>

      {orgsLoading ? (
        <div className="h-64 flex items-center justify-center">
          <Loader2 className="h-8 w-8 text-indigo-500 animate-spin" />
        </div>
      ) : !orgs || orgs.length === 0 ? (
        <div className="bg-slate-900 border border-slate-800 rounded-3xl p-12 text-center max-w-xl mx-auto space-y-4">
          <div className="p-4 bg-slate-800 border border-slate-700 rounded-full w-fit mx-auto text-slate-400">
            <Building2 className="h-8 w-8" />
          </div>
          <h3 className="text-lg font-bold text-white">No Organizations Available</h3>
          <p className="text-sm text-slate-400">
            You must create or join an organization before you can manage projects.
          </p>
        </div>
      ) : projectsLoading ? (
        <div className="h-64 flex items-center justify-center">
          <Loader2 className="h-8 w-8 text-indigo-500 animate-spin" />
        </div>
      ) : projects && projects.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projects.map((proj) => (
            <div key={proj.ID} className="bg-slate-900 border border-slate-800 rounded-3xl p-6 flex flex-col justify-between hover:border-slate-700 transition-colors">
              <div className="space-y-3">
                <div className="p-3 w-fit bg-slate-800 border border-slate-700 rounded-2xl text-teal-400">
                  <FolderKanban className="h-6 w-6" />
                </div>
                <h3 className="text-lg font-bold text-white tracking-tight">{proj.Name}</h3>
                <p className="text-sm text-slate-400 line-clamp-2">{proj.Description || 'No description provided.'}</p>
              </div>
              <div className="mt-6 pt-4 border-t border-slate-800/60 flex items-center justify-between text-xs text-slate-500">
                <span>Org ID: {proj.OrgID.substring(0, 8)}...</span>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="bg-slate-900 border border-slate-800 rounded-3xl p-12 text-center max-w-xl mx-auto space-y-4">
          <div className="p-4 bg-slate-800 border border-slate-700 rounded-full w-fit mx-auto text-slate-400">
            <FolderKanban className="h-8 w-8" />
          </div>
          <h3 className="text-lg font-bold text-white">No Projects Found</h3>
          <p className="text-sm text-slate-400">
            No projects exist under the selected organization. Create one to start scaffolding microservices!
          </p>
          <button
            onClick={() => setModalOpen(true)}
            className="px-5 py-2.5 rounded-2xl bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-sm transition-all cursor-pointer inline-block"
          >
            Create first Project
          </button>
        </div>
      )}

      {/* Create Modal */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
          <div className="bg-slate-900 border border-slate-800 rounded-3xl w-full max-w-lg overflow-hidden shadow-2xl relative">
            <div className="p-6 border-b border-slate-800 flex justify-between items-center">
              <h3 className="text-lg font-bold text-white">New Project</h3>
              <button onClick={() => setModalOpen(false)} className="text-slate-400 hover:text-slate-200 cursor-pointer">
                <X className="h-6 w-6" />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-6">
              {errorMsg && (
                <div className="bg-rose-500/10 border border-rose-500/20 rounded-2xl p-4 text-sm text-rose-400 flex gap-2">
                  <AlertCircle className="h-5 w-5 shrink-0" />
                  {errorMsg}
                </div>
              )}

              <div>
                <label className="block text-sm font-semibold text-slate-300">Name</label>
                <input
                  type="text"
                  required
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="mt-1 block w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-2xl text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all text-sm"
                  placeholder="e.g. core-services"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-slate-300">Description</label>
                <textarea
                  rows={3}
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="mt-1 block w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-2xl text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all text-sm"
                  placeholder="What is the purpose of this project..."
                />
              </div>

              <div className="flex gap-4 pt-2">
                <button
                  type="button"
                  onClick={() => setModalOpen(false)}
                  className="flex-1 py-3 px-4 rounded-2xl text-sm font-bold text-slate-400 bg-slate-850 hover:bg-slate-800 border border-slate-800 transition-all cursor-pointer"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={createProjMutation.isPending}
                  className="flex-1 py-3 px-4 rounded-2xl text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 transition-all cursor-pointer shadow-lg shadow-indigo-600/20 disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center"
                >
                  {createProjMutation.isPending ? (
                    <Loader2 className="animate-spin h-5 w-5 text-white" />
                  ) : (
                    'Create'
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
