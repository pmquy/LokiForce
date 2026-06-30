import { apiRequest } from "../../../services/api";

export interface Organization {
  ID: string;
  Name: string;
  Description: string;
  OwnerID: string;
}

export async function fetchMyOrganizations(): Promise<Organization[]> {
  return apiRequest<Organization[]>('/organizations');
}

export async function createOrganization(
  name: string,
  description: string,
): Promise<any> {
  return apiRequest<any>("/organizations", {
    method: "POST",
    body: JSON.stringify({ name, description }),
  });
}
