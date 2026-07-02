import { apiRequest } from "../../../services/api";

export interface Service {
  ID: string;
  Name: string;
  Description: string;
  ProjectID: string;
  TemplateID: string;
  Repository: string;
}

export interface Template {
  id: string;
  name: string;
  description: string;
}

export interface CreateServiceOutput {
  ServiceID: string;
  RepositoryURL: string;
}

export async function fetchServicesByProject(
  projectId: string,
): Promise<Service[]> {
  return apiRequest<Service[]>(`/services?project_id=${projectId}`);
}

export async function fetchTemplates(): Promise<Template[]> {
  return apiRequest<Template[]>("/services/templates");
}

export async function createService(
  name: string,
  description: string,
  projectId: string,
  templateId: string,
): Promise<CreateServiceOutput> {
  return apiRequest<CreateServiceOutput>("/services", {
    method: "POST",
    body: JSON.stringify({
      name,
      description,
      project_id: projectId,
      template_id: templateId,
    }),
  });
}

export async function deleteService(id: string): Promise<any> {
  return apiRequest<any>(`/services/${id}`, {
    method: "DELETE",
  });
}

export async function updateService(
  id: string,
  name: string,
  description: string,
): Promise<any> {
  return apiRequest<any>(`/services/${id}`, {
    method: "PUT",
    body: JSON.stringify({ name, description }),
  });
}

export interface AccessPolicy {
  id: string;
  client_id: string;
  target_id: string;
  target_port: string;
  project_id: string;
}

export async function fetchAccessPolicies(
  serviceId: string,
): Promise<AccessPolicy[]> {
  return apiRequest<AccessPolicy[]>(`/services/${serviceId}/policies`);
}

export async function createAccessPolicy(
  clientId: string,
  targetId: string,
  targetPort: string,
  projectId: string,
): Promise<{ id: string }> {
  return apiRequest<{ id: string }>("/services/policies", {
    method: "POST",
    body: JSON.stringify({
      client_id: clientId,
      target_id: targetId,
      target_port: targetPort,
      project_id: projectId,
    }),
  });
}

export async function deleteAccessPolicy(policyId: string): Promise<any> {
  return apiRequest<any>(`/services/policies/${policyId}`, {
    method: "DELETE",
  });
}
