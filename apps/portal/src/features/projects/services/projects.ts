import { apiRequest } from "../../../services/api";

export interface Project {
  ID: string;
  Name: string;
  Description: string;
  OrgID: string;
}

export async function fetchProjectsByOrg(orgId: string): Promise<Project[]> {
  return apiRequest<Project[]>(`/projects?org_id=${orgId}`);
}

export async function createProject(
  name: string,
  description: string,
  orgId: string,
): Promise<any> {
  return apiRequest<any>("/projects", {
    method: "POST",
    body: JSON.stringify({ name, description, org_id: orgId }),
  });
}
