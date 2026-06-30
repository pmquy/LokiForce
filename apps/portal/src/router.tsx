import {
  createRootRoute,
  createRoute,
  createRouter,
  Outlet,
} from "@tanstack/react-router";
import { Login } from "./features/auth/components/Login";
import { Register } from "./features/auth/components/Register";
import { DashboardLayout } from "./components/DashboardLayout";
import { DashboardHome } from "./components/DashboardHome";
import { OrganizationsList } from "./features/organizations/components/OrganizationsList";
import { ProjectsList } from "./features/projects/components/ProjectsList";
import { ServicesList } from "./features/services/components/ServicesList";

const rootRoute = createRootRoute({
  component: () => <Outlet />,
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/login",
  component: Login,
});

const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/register",
  component: Register,
});

const protectedRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: DashboardLayout,
});

const indexRoute = createRoute({
  getParentRoute: () => protectedRoute,
  path: "/",
  component: DashboardHome,
});

const organizationsRoute = createRoute({
  getParentRoute: () => protectedRoute,
  path: "/organizations",
  component: OrganizationsList,
});

const projectsRoute = createRoute({
  getParentRoute: () => protectedRoute,
  path: "/projects",
  component: ProjectsList,
});

const servicesRoute = createRoute({
  getParentRoute: () => protectedRoute,
  path: "/services",
  component: ServicesList,
});

const routeTree = rootRoute.addChildren([
  loginRoute,
  registerRoute,
  protectedRoute.addChildren([
    indexRoute,
    organizationsRoute,
    projectsRoute,
    servicesRoute,
  ]),
]);

export const router = createRouter({ routeTree });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
