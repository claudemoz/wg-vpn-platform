import { router } from 'expo-router';
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';

import { authApi, userApi } from '@/services/api';
import * as storage from '@/services/storage';
import type { User } from '@/types/api';

type AuthContextValue = {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (data: {
    first_name: string;
    last_name: string;
    email: string;
    phone?: string;
    password: string;
  }) => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    (async () => {
      try {
        const savedToken = await storage.getToken();
        if (!savedToken) return;

        const savedUser = await storage.getUser<User>();
        if (savedUser) {
          setToken(savedToken);
          setUser(savedUser);
          return;
        }

        const me = await userApi.me(savedToken);
        setToken(savedToken);
        setUser(me);
        await storage.setUser(me);
      } catch {
        await storage.clearAuth();
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const res = await authApi.login(email.trim(), password);
    await storage.setToken(res.token);
    await storage.setUser(res.user);
    setToken(res.token);
    setUser(res.user);
    router.replace('/(app)');
  }, []);

  const register = useCallback(
    async (data: {
      first_name: string;
      last_name: string;
      email: string;
      phone?: string;
      password: string;
    }) => {
      await authApi.register({
        ...data,
        first_name: data.first_name.trim(),
        last_name: data.last_name.trim(),
        email: data.email.trim(),
        phone: data.phone?.trim() || undefined,
      });
    },
    [],
  );

  const logout = useCallback(async () => {
    await storage.clearAll();
    setToken(null);
    setUser(null);
    router.replace('/(auth)/login');
  }, []);

  const value = useMemo(
    () => ({ user, token, isLoading, login, register, logout }),
    [user, token, isLoading, login, register, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
