import { useMutation } from "@tanstack/react-query";
import { loginUser, registerUser } from "../services/auth";

export function useLoginMutation(
  onSuccess: (data: { token: string }) => void,
  onError: (err: any) => void,
) {
  return useMutation({
    mutationFn: ({ email, password }: any) => loginUser(email, password),
    onSuccess,
    onError,
  });
}

export function useRegisterMutation(
  onSuccess: () => void,
  onError: (err: any) => void,
) {
  return useMutation({
    mutationFn: ({ username, email, password }: any) =>
      registerUser(username, email, password),
    onSuccess,
    onError,
  });
}
