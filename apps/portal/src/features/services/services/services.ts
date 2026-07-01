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
