import * as Device from 'expo-device';
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';

import { useAuth } from '@/contexts/AuthContext';
import { deviceApi, serverApi } from '@/services/api';
import * as storage from '@/services/storage';
import type { Server } from '@/types/api';

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected';

type VpnContextValue = {
  servers: Server[];
  selectedServer: Server | null;
  status: ConnectionStatus;
  assignedIp: string | null;
  isLoadingServers: boolean;
  error: string | null;
  selectServer: (server: Server) => void;
  connect: () => Promise<void>;
  disconnect: () => Promise<void>;
  refreshServers: () => Promise<void>;
};

const VpnContext = createContext<VpnContextValue | null>(null);

function deviceLabel(): string {
  const name = Device.deviceName ?? Device.modelName ?? 'Mobile';
  return name.slice(0, 100);
}

export function VpnProvider({ children }: { children: ReactNode }) {
  const { token } = useAuth();
  const [servers, setServers] = useState<Server[]>([]);
  const [selectedServer, setSelectedServer] = useState<Server | null>(null);
  const [status, setStatus] = useState<ConnectionStatus>('disconnected');
  const [assignedIp, setAssignedIp] = useState<string | null>(null);
  const [isLoadingServers, setIsLoadingServers] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const refreshServers = useCallback(async () => {
    if (!token) return;
    setIsLoadingServers(true);
    setError(null);
    try {
      const list = await serverApi.list(token);
      const active = list.filter((s) => s.is_active);
      setServers(active);

      const session = await storage.getVpnSession();
      setSelectedServer((prev) => {
        if (prev && active.some((s) => s.id === prev.id)) return prev;
        if (session.serverId) {
          const saved = active.find((s) => s.id === session.serverId);
          if (saved) return saved;
        }
        return active[0] ?? null;
      });
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Impossible de charger les serveurs');
    } finally {
      setIsLoadingServers(false);
    }
  }, [token]);

  useEffect(() => {
    refreshServers();
  }, [token]);

  useEffect(() => {
    (async () => {
      const session = await storage.getVpnSession();
      if (session.deviceId && session.config) {
        setStatus('connected');
        const match = session.config.match(/Address = ([\d.]+)/);
        if (match) setAssignedIp(match[1]);
      }
    })();
  }, []);

  const selectServer = useCallback((server: Server) => {
    if (status === 'connected') return;
    setSelectedServer(server);
    setError(null);
  }, [status]);

  const connect = useCallback(async () => {
    if (!token || !selectedServer) {
      setError('Sélectionnez un serveur');
      return;
    }

    setStatus('connecting');
    setError(null);

    try {
      const session = await storage.getVpnSession();
      if (session.deviceId && session.serverId === selectedServer.id) {
        setAssignedIp(session.config?.match(/Address = ([\d.]+)/)?.[1] ?? null);
        setStatus('connected');
        return;
      }

      const device = await deviceApi.create(token, {
        server_id: selectedServer.id,
        name: deviceLabel(),
      });

      await storage.saveVpnSession({
        deviceId: device.id,
        serverId: selectedServer.id,
        config: device.config,
        privateKey: device.private_key,
      });

      setAssignedIp(device.assigned_ip);
      setStatus('connected');

      // Le tunnel WireGuard natif nécessite un module natif (Network Extension iOS /
      // VpnService Android). La config est prête côté API ; l'intégration native
      // sera branchée ici via l'agent WireGuard.
    } catch (e) {
      setStatus('disconnected');
      setError(e instanceof Error ? e.message : 'Échec de la connexion');
    }
  }, [token, selectedServer]);

  const disconnect = useCallback(async () => {
    setStatus('disconnected');
    setAssignedIp(null);
    await storage.clearVpnSession();
  }, []);

  const value = useMemo(
    () => ({
      servers,
      selectedServer,
      status,
      assignedIp,
      isLoadingServers,
      error,
      selectServer,
      connect,
      disconnect,
      refreshServers,
    }),
    [
      servers,
      selectedServer,
      status,
      assignedIp,
      isLoadingServers,
      error,
      selectServer,
      connect,
      disconnect,
      refreshServers,
    ],
  );

  return <VpnContext.Provider value={value}>{children}</VpnContext.Provider>;
}

export function useVpn() {
  const ctx = useContext(VpnContext);
  if (!ctx) throw new Error('useVpn must be used within VpnProvider');
  return ctx;
}
