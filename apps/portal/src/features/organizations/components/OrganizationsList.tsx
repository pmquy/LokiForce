import React, { useState } from "react";
import {
  useOrganizationsQuery,
  useCreateOrgMutation,
} from "../hooks/useOrganizations";
import { Building2, Plus, Loader2, X, AlertCircle } from "lucide-react";

export function OrganizationsList() {
  const [modalOpen, setModalOpen] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [errorMsg, setErrorMsg] = useState("");

  const { data: orgs, isLoading } = useOrganizationsQuery();

  const createOrgMutation = useCreateOrgMutation(
    () => {
      setModalOpen(false);
      setName("");
      setDescription("");
    },
    (err: any) => {
      setErrorMsg(err.message || "Failed to create organization");
    },
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg("");
    createOrgMutation.mutate({ name, description });
  };

  return (
    <div className="space-y-6 max-w-7xl mx-auto">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-extrabold text-white tracking-tight">
            Organizations
          </h2>
          <p className="text-sm text-slate-400">
            Manage your business domains and namespace boundaries.
          </p>
        </div>
        <button
          onClick={() => setModalOpen(true)}
          className="flex items-center gap-2 px-4 py-2.5 rounded-2xl bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-sm transition-all cursor-pointer shadow-lg shadow-indigo-600/25"
        >
          <Plus className="h-4.5 w-4.5" />
          Create Organization
        </button>
      </div>

      {isLoading ? (
        <div className="h-64 flex items-center justify-center">
          <Loader2 className="h-8 w-8 text-indigo-500 animate-spin" />
        </div>
      ) : orgs && orgs.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {orgs.map((org) => (
            <div
              key={org.ID}
              className="bg-slate-900 border border-slate-800 rounded-3xl p-6 flex flex-col justify-between hover:border-slate-700 transition-colors"
            >
              <div className="space-y-3">
                <div className="p-3 w-fit bg-slate-800 border border-slate-700 rounded-2xl text-indigo-400">
                  <Building2 className="h-6 w-6" />
                </div>
                <h3 className="text-lg font-bold text-white tracking-tight">
                  {org.Name}
                </h3>
                <p className="text-sm text-slate-400 line-clamp-2">
                  {org.Description || "No description provided."}
                </p>
              </div>
              <div className="mt-6 pt-4 border-t border-slate-800/60 flex items-center text-xs text-slate-500">
                <span>Owner ID: {org.OwnerID.substring(0, 8)}...</span>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="bg-slate-900 border border-slate-800 rounded-3xl p-12 text-center max-w-xl mx-auto space-y-4">
          <div className="p-4 bg-slate-800 border border-slate-700 rounded-full w-fit mx-auto text-slate-400">
            <Building2 className="h-8 w-8" />
          </div>
          <h3 className="text-lg font-bold text-white">
            No Organizations Found
          </h3>
          <p className="text-sm text-slate-400">
            You don't belong to any organizations yet. Create one to start
            hosting projects and repositories.
          </p>
          <button
            onClick={() => setModalOpen(true)}
            className="px-5 py-2.5 rounded-2xl bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-sm transition-all cursor-pointer inline-block"
          >
            Create your first Organization
          </button>
        </div>
      )}

      {/* Create Modal */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
          <div className="bg-slate-900 border border-slate-800 rounded-3xl w-full max-w-lg overflow-hidden shadow-2xl relative">
            <div className="p-6 border-b border-slate-800 flex justify-between items-center">
              <h3 className="text-lg font-bold text-white">New Organization</h3>
              <button
                onClick={() => setModalOpen(false)}
                className="text-slate-400 hover:text-slate-200 cursor-pointer"
              >
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
                <label className="block text-sm font-semibold text-slate-300">
                  Name
                </label>
                <input
                  type="text"
                  required
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="mt-1 block w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-2xl text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all text-sm"
                  placeholder="e.g. LokiForce Engineering"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-slate-300">
                  Description
                </label>
                <textarea
                  rows={3}
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="mt-1 block w-full px-4 py-3 bg-slate-950 border border-slate-800 rounded-2xl text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all text-sm"
                  placeholder="Tell us what this organization does..."
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
                  disabled={createOrgMutation.isPending}
                  className="flex-1 py-3 px-4 rounded-2xl text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 transition-all cursor-pointer shadow-lg shadow-indigo-600/20 disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center"
                >
                  {createOrgMutation.isPending ? (
                    <Loader2 className="animate-spin h-5 w-5 text-white" />
                  ) : (
                    "Create"
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
