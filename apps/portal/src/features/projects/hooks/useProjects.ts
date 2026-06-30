import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchProjectsByOrg,
  createProject,
  type Project,
} from "../services/projects";

export function useProjectsQuery(orgId: string) {
  return useQuery<Project[]>({
    queryKey: ["projects", orgId],
    queryFn: () => fetchProjectsByOrg(orgId),
    enabled: !!orgId,
  });
}

export function useCreateProjectMutation(
  orgId: string,
  onSuccess: () => void,
  onError: (err: any) => void,
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      name,
      description,
    }: {
      name: string;
      description: string;
    }) => createProject(name, description, orgId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", orgId] });
      onSuccess();
    },
    onError,
  });
}
