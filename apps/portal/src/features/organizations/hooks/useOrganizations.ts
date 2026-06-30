import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { fetchMyOrganizations, createOrganization, type Organization } from '../services/organizations';

export function useOrganizationsQuery() {
  return useQuery<Organization[]>({
    queryKey: ['organizations'],
    queryFn: fetchMyOrganizations,
  });
}

export function useCreateOrgMutation(onSuccess: () => void, onError: (err: any) => void) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ name, description }: { name: string; description: string }) =>
      createOrganization(name, description),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['organizations'] });
      onSuccess();
    },
    onError,
  });
}
