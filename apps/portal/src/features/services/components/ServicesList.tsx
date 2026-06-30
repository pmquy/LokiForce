import React, { useState, useEffect } from 'react';
import { useOrganizationsQuery } from '../../organizations/hooks/useOrganizations';
import { useProjectsQuery } from '../../projects/hooks/useProjects';
import { useServicesQuery, useTemplatesQuery, useCreateServiceMutation } from '../hooks/useServices';
import { Cpu, Plus, Loader2, X, AlertCircle, ChevronDown, GitBranch, Copy, Check, Terminal, ExternalLink } from 'lucide-react';

export function ServicesList() {
  const [selectedOrgId, setSelectedOrgId] = useState<string>('');
  const [selectedProjId, setSelectedProjId] = useState<string>('');
  const [modalOpen, setModalOpen] = useState(false);
  const [successModalOpen, setSuccessModalOpen] = useState(false);

  // Form State
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [templateId, setTemplateId] = useState('golang');
  const [errorMsg, setErrorMsg] = useState('');
  const [generatedRepoUrl, setGeneratedRepoUrl] = useState('');
  const [copied, setCopied] = useState(false);

  // 1. Fetch Organizations
  const { data: orgs, isLoading: orgsLoading } = useOrganizationsQuery();

  useEffect(() => {
    if (orgs && orgs.length > 0 && !selectedOrgId) {
      setSelectedOrgId(orgs[0].ID);
    }
  }, [orgs, selectedOrgId]);

  // 2. Fetch Projects under Selected Org
  const { data: projects, isLoading: projectsLoading } = useProjectsQuery(selectedOrgId);

  useEffect(() => {
    if (projects && projects.length > 0) {
      setSelectedProjId(projects[0].ID);
    } else {
      setSelectedProjId('');
    }
  }, [projects, selectedOrgId]);

  // 3. Fetch Services under Selected Project
  const { data: services, isLoading: servicesLoading } = useServicesQuery(selectedProjId);

  // 4. Fetch Templates list
  const { data: templates } = useTemplatesQuery();

  // Create Service & Scaffold Mutation
  const createServiceMutation = useCreateServiceMutation(
    selectedProjId,
    (data) => {
      setGeneratedRepoUrl(data.RepositoryURL);
      setModalOpen(false);
      setSuccessModalOpen(true);
      setName('');
      setDescription('');
    },
    (err: any) => {
      setErrorMsg(err.message || 'Failed to scaffold service');
    }
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg('');
    createServiceMutation.mutate({ name, description, templateId });
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(`git clone ${generatedRepoUrl}`);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col lg:flex-row justify-between lg:items-center gap-4">
        <div>
          <h2 className="text-2xl font-extrabold text-white tracking-tight">Service Catalog</h2>
          <p className="text-sm text-slate-400">Scaffold and bootstrap standardization-compliant microservices using Golden Paths.</p>
        </div>

        <div className="flex flex-wrap items-center gap-3">
          {/* Org Selector */}
          {orgs && orgs.length > 0 && (
            <div className="relative">
              <select
                value={selectedOrgId}
                onChange={(e) => setSelectedOrgId(e.target.value)}
                className="appearance-none bg-slate-900 border border-slate-800 rounded-2xl pl-4 pr-10 py-2.5 text-sm font-semibold text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all cursor-pointer min-w-[160px]"
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

          {/* Project Selector */}
          {projects && projects.length > 0 && (
            <div className="relative">
              <select
                value={selectedProjId}
                onChange={(e) => setSelectedProjId(e.target.value)}
                className="appearance-none bg-slate-900 border border-slate-800 rounded-2xl pl-4 pr-10 py-2.5 text-sm font-semibold text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all cursor-pointer min-w-[160px]"
              >
                {projects.map((proj) => (
                  <option key={proj.ID} value={proj.ID}>
                    {proj.Name}
                  </option>
                ))}
              </select>
              <ChevronDown className="absolute right-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-500 pointer-events-none" />
            </div>
          )}

          {selectedProjId && (
            <button
              onClick={() => setModalOpen(true)}
              className="flex items-center gap-2 px-4 py-2.5 rounded-2xl bg-gradient-to-r from-indigo-500 to-teal-500 text-white font-bold text-sm transition-all cursor-pointer shadow-lg shadow-indigo-500/20"
            >
              <Plus className="h-4.5 w-4.5" />
              Scaffold New Service
            </button>
          )}
        </div>
      </div>

      {orgsLoading || projectsLoading ? (
        <div className="h-64 flex items-center justify-center">
          <Loader2 className="h-8 w-8 text-indigo-500 animate-spin" />
        </div>
      ) : !orgs || orgs.length === 0 ? (
        <div className="bg-slate-900 border border-slate-800 rounded-3xl p-12 text-center max-w-xl mx-auto space-y-4">
          <h3 className="text-lg font-bold text-white">No Organizations Available</h3>
          <p className="text-sm text-slate-400">Please create an organization first before creating services.</p>
        </div>
      ) : !projects || projects.length === 0 ? (
        <div className="bg-slate-900 border border-slate-800 rounded-3xl p-12 text-center max-w-xl mx-auto space-y-4">
          <h3 className="text-lg font-bold text-white">No Projects Available</h3>
          <p className="text-sm text-slate-400">Please create a project inside the selected organization first.</p>
        </div>
      ) : servicesLoading ? (
        <div className="h-64 flex items-center justify-center">
          <Loader2 className="h-8 w-8 text-indigo-500 animate-spin" />
        </div>
      ) : services && services.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {services.map((svc) => (
            <div key={svc.ID} className="bg-slate-900 border border-slate-800 rounded-3xl p-6 flex flex-col justify-between hover:border-slate-700 transition-colors">
              <div className="space-y-4">
                <div className="flex justify-between items-start">
                  <div className="p-3 bg-slate-800 border border-slate-700 rounded-2xl text-indigo-400">
                    <Cpu className="h-6 w-6" />
                  </div>
                  <span className="px-3 py-1 bg-slate-800 border border-slate-750 rounded-full text-xs font-semibold text-slate-400 capitalize">
                    {svc.TemplateID}
                  </span>
                </div>
                <div>
                  <h3 className="text-lg font-bold text-white tracking-tight">{svc.Name}</h3>
                  <p className="text-sm text-slate-400 mt-1 line-clamp-2">{svc.Description || 'No description provided.'}</p>
                </div>
              </div>

              {svc.Repository && (
                <div className="mt-6 pt-4 border-t border-slate-800/60 flex items-center justify-between">
                  <span className="text-xs font-semibold text-slate-500 flex items-center gap-1.5">
                    <GitBranch className="h-4 w-4" />
                    GitHub Repository
                  </span>
                  <a
                    href={svc.Repository}
                    target="_blank"
                    rel="noreferrer"
                    className="text-xs font-bold text-indigo-400 hover:text-indigo-300 transition-colors flex items-center gap-1 cursor-pointer"
                  >
                    View Code
                    <ExternalLink className="h-3 w-3" />
                  </a>
                </div>
              )}
            </div>
          ))}
        </div>
      ) : (
        <div className="bg-slate-900 border border-slate-800 rounded-3xl p-12 text-center max-w-xl mx-auto space-y-4">
          <div className="p-4 bg-slate-800 border border-slate-700 rounded-full w-fit mx-auto text-slate-400">
            <Cpu className="h-8 w-8" />
          </div>
          <h3 className="text-lg font-bold text-white">No Services Found</h3>
          <p className="text-sm text-slate-400">
            No services have been scaffolded in this project yet. Start a new Golden Path now!
          </p>
          <button
            onClick={() => setModalOpen(true)}
            className="px-5 py-2.5 rounded-2xl bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-sm transition-all cursor-pointer inline-block"
          >
            Scaffold your first Service
          </button>
        </div>
      )}

      {/* Scaffold Modal */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
          <div className="bg-slate-900 border border-slate-800 rounded-3xl w-full max-w-xl overflow-hidden shadow-2xl relative">
            <div className="p-6 border-b border-slate-800 flex justify-between items-center">
              <h3 className="text-lg font-bold text-white">Scaffold New Microservice</h3>
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
                <label className="block text-sm font-semibold text-slate-300">Service Name</label>
                <input
                  type="text"
                  required
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="mt-1 block w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-2xl text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all text-sm"
                  placeholder="e.g. user-payment-service"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-slate-300">Description</label>
                <textarea
                  rows={2}
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="mt-1 block w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-2xl text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all text-sm"
                  placeholder="Description of the service responsibilities..."
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-slate-300 mb-3">Select Golden Path Template</label>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  {templates &&
                    templates.map((tmpl) => (
                      <div
                        key={tmpl.id}
                        onClick={() => setTemplateId(tmpl.id)}
                        className={`p-4 border rounded-2xl cursor-pointer flex flex-col justify-between gap-2 transition-all ${
                          templateId === tmpl.id
                            ? 'border-indigo-500 bg-indigo-500/5'
                            : 'border-slate-800 bg-slate-950 hover:border-slate-700'
                        }`}
                      >
                        <div>
                          <h4 className="font-bold text-sm text-slate-200">{tmpl.name}</h4>
                          <p className="text-xs text-slate-400 mt-1">{tmpl.description}</p>
                        </div>
                      </div>
                    ))}
                </div>
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
                  disabled={createServiceMutation.isPending}
                  className="flex-1 py-3 px-4 rounded-2xl text-sm font-bold text-white bg-gradient-to-r from-indigo-500 to-teal-500 hover:from-indigo-600 hover:to-teal-600 transition-all cursor-pointer shadow-lg shadow-indigo-500/20 disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center"
                >
                  {createServiceMutation.isPending ? (
                    <Loader2 className="animate-spin h-5 w-5 text-white" />
                  ) : (
                    'Scaffold & Push'
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Success Modal */}
      {successModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
          <div className="bg-slate-900 border border-slate-800 rounded-3xl w-full max-w-xl overflow-hidden shadow-2xl relative p-6 space-y-6 text-center">
            <div className="w-16 h-16 rounded-full bg-teal-500/10 border border-teal-500/20 flex items-center justify-center mx-auto text-teal-400">
              <Check className="h-8 w-8" />
            </div>

            <div className="space-y-2">
              <h3 className="text-2xl font-extrabold text-white tracking-tight">Microservice Scaffolded!</h3>
              <p className="text-sm text-slate-400 max-w-md mx-auto">
                Your new service codebase has been successfully created and pushed to your remote GitHub repository.
              </p>
            </div>

            <div className="bg-slate-950 border border-slate-800 rounded-2xl p-4 space-y-3 text-left">
              <span className="text-xs font-bold text-slate-500 uppercase tracking-wider block">Repository URL</span>
              <a
                href={generatedRepoUrl}
                target="_blank"
                rel="noreferrer"
                className="text-sm font-semibold text-indigo-400 hover:text-indigo-300 transition-colors flex items-center gap-1.5 truncate"
              >
                <GitBranch className="h-4 w-4 text-slate-400 shrink-0" />
                {generatedRepoUrl}
                <ExternalLink className="h-3.5 w-3.5 shrink-0" />
              </a>
            </div>

            <div className="bg-slate-950 border border-slate-800 rounded-2xl p-4 space-y-2 text-left">
              <span className="text-xs font-bold text-slate-500 uppercase tracking-wider flex items-center gap-1">
                <Terminal className="h-3.5 w-3.5 text-teal-400" />
                Clone Command
              </span>
              <div className="flex items-center justify-between gap-3 bg-slate-900 px-3.5 py-2.5 rounded-xl border border-slate-800 font-mono text-xs text-slate-300">
                <span className="truncate">git clone {generatedRepoUrl}</span>
                <button
                  onClick={handleCopy}
                  className="p-1.5 rounded-lg border border-slate-700 hover:bg-slate-800 text-slate-400 hover:text-slate-200 cursor-pointer transition-colors shrink-0"
                >
                  {copied ? <Check className="h-4 w-4 text-teal-400" /> : <Copy className="h-4 w-4" />}
                </button>
              </div>
            </div>

            <button
              onClick={() => setSuccessModalOpen(false)}
              className="w-full py-3 px-4 rounded-2xl text-sm font-bold text-white bg-slate-800 hover:bg-slate-750 border border-slate-700 transition-all cursor-pointer"
            >
              Done
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
