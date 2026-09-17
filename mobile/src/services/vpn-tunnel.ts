import { NativeEventEmitter, NativeModules, Platform } from 'react-native';
import type {
  VpnStateChangedEvent,
  WireGuardConfig,
  WireGuardStatus,
} from 'react-native-wireguard-vpn';

import { wgQuickToNativeConfig } from '@/services/wg-config-parser';

export type TunnelStateListener = (event: VpnStateChangedEvent) => void;

function isNativePlatform(): boolean {
  return Platform.OS === 'ios' || Platform.OS === 'android';
}

async function loadModule() {
  if (!isNativePlatform()) return null;
  try {
    const mod = await import('react-native-wireguard-vpn');
    return mod.default;
  } catch {
    return null;
  }
}

/** Vrai si le module natif WireGuard est disponible (pas Expo Go / pas web). */
export async function isVpnTunnelAvailable(): Promise<boolean> {
  const mod = await loadModule();
  if (!mod) return false;
  try {
    return await mod.isSupported();
  } catch {
    return false;
  }
}

let initialized = false;

export async function initializeTunnel(): Promise<void> {
  const mod = await loadModule();
  if (!mod) {
    throw new Error('Tunnel VPN indisponible sur cette plateforme (utilisez un build natif)');
  }
  if (!initialized) {
    await mod.initialize();
    initialized = true;
  }
}

export async function connectTunnelFromConfig(
  configText: string,
  privateKey?: string,
): Promise<void> {
  const nativeConfig = wgQuickToNativeConfig(configText, privateKey);
  await connectTunnel(nativeConfig);
}

export async function connectTunnel(config: WireGuardConfig): Promise<void> {
  const mod = await loadModule();
  if (!mod) {
    throw new Error('Tunnel VPN indisponible. Lancez avec : npx expo run:ios --device');
  }
  await initializeTunnel();
  await mod.connect(config);
}

export async function disconnectTunnel(): Promise<void> {
  const mod = await loadModule();
  if (!mod) return;
  await mod.disconnect();
}

export async function getTunnelStatus(): Promise<WireGuardStatus | null> {
  const mod = await loadModule();
  if (!mod) return null;
  try {
    return await mod.getStatus();
  } catch {
    return null;
  }
}

/** Écoute les changements d'état du tunnel (Network Extension / VpnService). */
export function subscribeTunnelState(listener: TunnelStateListener): () => void {
  if (!isNativePlatform()) {
    return () => {};
  }

  const moduleRef = NativeModules.WireGuardVpnModule;
  if (!moduleRef) {
    return () => {};
  }

  const emitter = new NativeEventEmitter(moduleRef);
  const sub = emitter.addListener('vpnStateChanged', listener);
  return () => sub.remove();
}

export function mapTunnelStatusToConnectionStatus(
  status: WireGuardStatus['status'] | WireGuardStatus['tunnelState'],
): 'disconnected' | 'connecting' | 'connected' {
  switch (status) {
    case 'CONNECTED':
    case 'ACTIVE':
      return 'connected';
    case 'CONNECTING':
      return 'connecting';
    default:
      return 'disconnected';
  }
}
