import { useState } from "react";
import { useNavigate, Link } from "@tanstack/react-router";
import { useRegisterMutation } from "../hooks/useAuth";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  Shield,
  Key,
  Mail,
  User as UserIcon,
  Loader2,
  AlertCircle,
} from "lucide-react";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../../../components/ui/form";

const registerSchema = z.object({
  username: z
    .string()
    .min(3, { message: "Tên người dùng phải chứa ít nhất 3 ký tự" }),
  email: z.string().email({ message: "Vui lòng nhập địa chỉ email hợp lệ" }),
  password: z
    .string()
    .min(8, { message: "Mật khẩu phải chứa ít nhất 8 ký tự" }),
});

type RegisterValues = z.infer<typeof registerSchema>;

export function Register() {
  const navigate = useNavigate();
  const [errorMsg, setErrorMsg] = useState("");
  const [successMsg, setSuccessMsg] = useState("");

  const form = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
    },
  });

  const registerMutation = useRegisterMutation(
    () => {
      setSuccessMsg(
        "Đăng ký tài khoản thành công! Đang chuyển hướng sang đăng nhập...",
      );
      setTimeout(() => {
        navigate({ to: "/login" });
      }, 1500);
    },
    (err: any) => {
      setErrorMsg(err.message || "Đăng ký thất bại. Vui lòng thử lại.");
    },
  );

  const onSubmit = (values: RegisterValues) => {
    setErrorMsg("");
    setSuccessMsg("");
    registerMutation.mutate({
      username: values.username,
      email: values.email,
      password: values.password,
    });
  };

  return (
    <div className="min-h-screen bg-slate-950 flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden font-sans">
      {/* Background glow effects */}
      <div className="absolute top-1/4 left-1/4 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-indigo-500/10 rounded-full blur-3xl" />
      <div className="absolute bottom-1/4 right-1/4 translate-x-1/2 translate-y-1/2 w-96 h-96 bg-teal-500/10 rounded-full blur-3xl" />

      <div className="sm:mx-auto sm:w-full sm:max-w-md z-10">
        <div className="flex justify-center items-center gap-3">
          <div className="p-3 bg-gradient-to-tr from-indigo-500 to-teal-500 rounded-2xl shadow-lg shadow-indigo-500/20">
            <Shield className="h-8 w-8 text-white animate-pulse" />
          </div>
          <span className="text-3xl font-extrabold tracking-tight bg-gradient-to-r from-white via-slate-200 to-slate-400 bg-clip-text text-transparent">
            LokiForce
          </span>
        </div>
        <h2 className="mt-6 text-center text-3xl font-extrabold text-white tracking-tight">
          Create your account
        </h2>
        <p className="mt-2 text-center text-sm text-slate-400">
          Already have an account?{" "}
          <Link
            to="/login"
            className="font-medium text-indigo-400 hover:text-indigo-300 transition-colors"
          >
            Sign in instead
          </Link>
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md z-10 px-4">
        <div className="bg-slate-900/60 backdrop-blur-xl py-8 px-4 shadow-2xl rounded-3xl border border-slate-800 sm:px-10">
          <Form {...form}>
            <form className="space-y-6" onSubmit={form.handleSubmit(onSubmit)}>
              {errorMsg && (
                <div className="bg-rose-500/10 border border-rose-500/20 rounded-2xl p-4 text-sm text-rose-400 flex gap-2">
                  <AlertCircle className="h-5 w-5 shrink-0" />
                  {errorMsg}
                </div>
              )}
              {successMsg && (
                <div className="bg-teal-500/10 border border-teal-500/20 rounded-2xl p-4 text-sm text-teal-400">
                  {successMsg}
                </div>
              )}

              <FormField
                control={form.control}
                name="username"
                render={({ field }: any) => (
                  <FormItem>
                    <FormLabel>Username</FormLabel>
                    <FormControl>
                      <div className="relative rounded-2xl shadow-sm">
                        <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                          <UserIcon className="h-5 w-5 text-slate-500" />
                        </div>
                        <Input
                          {...field}
                          type="text"
                          className="pl-11"
                          placeholder="johndoe"
                        />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="email"
                render={({ field }: any) => (
                  <FormItem>
                    <FormLabel>Email address</FormLabel>
                    <FormControl>
                      <div className="relative rounded-2xl shadow-sm">
                        <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                          <Mail className="h-5 w-5 text-slate-500" />
                        </div>
                        <Input
                          {...field}
                          type="email"
                          className="pl-11"
                          placeholder="name@company.com"
                        />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="password"
                render={({ field }: any) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <div className="relative rounded-2xl shadow-sm">
                        <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                          <Key className="h-5 w-5 text-slate-500" />
                        </div>
                        <Input
                          {...field}
                          type="password"
                          className="pl-11"
                          placeholder="••••••••"
                        />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button
                type="submit"
                disabled={registerMutation.isPending}
                className="w-full h-12"
              >
                {registerMutation.isPending ? (
                  <Loader2 className="animate-spin h-5 w-5 text-white" />
                ) : (
                  "Create account"
                )}
              </Button>
            </form>
          </Form>
        </div>
      </div>
    </div>
  );
}
