import { apiRequest } from "../../../services/api";

export async function loginUser(
  email: string,
  password: string,
): Promise<{ token: string }> {
  return apiRequest<{ token: string }>("/users/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export async function registerUser(
  username: string,
  email: string,
  password: string,
): Promise<{ user_id: string }> {
  return apiRequest<{ user_id: string }>("/users/register", {
    method: "POST",
    body: JSON.stringify({ username, email, password }),
  });
}
