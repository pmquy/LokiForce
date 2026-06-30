import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { fetchServicesByProject, fetchTemplates, createService, type Service, type Template, type CreateServiceOutput } from '../services/services';

export function useServicesQuery(projectId: string) {
  return useQuery<Service[]>({
    queryKey: ['services', projectId],
    queryFn: () => fetchServicesByProject(projectId),
    enabled: !!projectId,
  });
}

export function useTemplatesQuery() {
  return useQuery<Template[]>({
    queryKey: ['templates'],
    queryFn: fetchTemplates,
  });
}

export function useCreateServiceMutation(projectId: string, onSuccess: (data: CreateServiceOutput) => void, onError: (err: any) => void) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ name, description, templateId }: { name: string; description: string; templateId: string }) =>
      createService(name, description, projectId, templateId),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['services', projectId] });
      onSuccess(data);
    },
    onError,
  });
}
