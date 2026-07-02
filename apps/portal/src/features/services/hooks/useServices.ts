import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchServicesByProject,
  fetchTemplates,
  createService,
  deleteService,
  updateService,
  fetchAccessPolicies,
  createAccessPolicy,
  deleteAccessPolicy,
  type Service,
  type Template,
  type CreateServiceOutput,
  type AccessPolicy,
} from "../services/services";

export function useServicesQuery(projectId: string) {
  return useQuery<Service[]>({
    queryKey: ["services", projectId],
    queryFn: () => fetchServicesByProject(projectId),
    enabled: !!projectId,
  });
}

export function useTemplatesQuery() {
  return useQuery<Template[]>({
    queryKey: ["templates"],
    queryFn: fetchTemplates,
  });
}

export function useCreateServiceMutation(
  projectId: string,
  onSuccess: (data: CreateServiceOutput) => void,
  onError: (err: any) => void,
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      name,
      description,
      templateId,
    }: {
      name: string;
      description: string;
      templateId: string;
    }) => createService(name, description, projectId, templateId),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["services", projectId] });
      onSuccess(data);
    },
    onError,
  });
}

export function useDeleteServiceMutation(
  projectId: string,
  onSuccess: () => void,
  onError: (err: any) => void,
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteService(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["services", projectId] });
      onSuccess();
    },
    onError,
  });
}

export function useUpdateServiceMutation(
  projectId: string,
  onSuccess: () => void,
  onError: (err: any) => void,
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      name,
      description,
    }: {
      id: string;
      name: string;
      description: string;
    }) => updateService(id, name, description),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["services", projectId] });
      onSuccess();
    },
    onError,
  });
}

export function useAccessPoliciesQuery(serviceId: string) {
  return useQuery<AccessPolicy[]>({
    queryKey: ["accessPolicies", serviceId],
    queryFn: () => fetchAccessPolicies(serviceId),
    enabled: !!serviceId,
  });
}

export function useCreateAccessPolicyMutation(
  serviceId: string,
  onSuccess: () => void,
  onError: (err: any) => void,
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      clientId,
      targetPort,
      projectId,
    }: {
      clientId: string;
      targetPort: string;
      projectId: string;
    }) => createAccessPolicy(clientId, serviceId, targetPort, projectId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["accessPolicies", serviceId] });
      onSuccess();
    },
    onError,
  });
}

export function useDeleteAccessPolicyMutation(
  serviceId: string,
  onSuccess: () => void,
  onError: (err: any) => void,
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (policyId: string) => deleteAccessPolicy(policyId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["accessPolicies", serviceId] });
      onSuccess();
    },
    onError,
  });
}

