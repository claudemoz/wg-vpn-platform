declare module 'react-native-wireguard-vpn' {
  export interface WireGuardConfig {
    privateKey: string;
    publicKey: string;
    serverAddress: string;
    serverPort: number;
    address?: string | string[];
    allowedIPs: string[];
    dns?: string[];
    mtu?: number;
    presharedKey?: string;
  }

  export interface WireGuardStatus {
    isConnected: boolean;
    tunnelState:
      | 'ACTIVE'
      | 'INACTIVE'
      | 'CONNECTING'
      | 'DISCONNECTING'
      | 'ERROR'
      | 'UNKNOWN';
    status:
      | 'CONNECTED'
      | 'DISCONNECTED'
      | 'CONNECTING'
      | 'DISCONNECTING'
      | 'ERROR'
      | 'UNKNOWN';
    error?: string;
  }

  export interface VpnStateChangedEvent {
    isConnected: boolean;
    tunnelState: WireGuardStatus['tunnelState'];
    status: WireGuardStatus['status'];
  }

  const WireGuardVpnModule: {
    initialize(): Promise<void>;
    connect(config: WireGuardConfig): Promise<void>;
    disconnect(): Promise<void>;
    getStatus(): Promise<WireGuardStatus>;
    isSupported(): Promise<boolean>;
  };

  export default WireGuardVpnModule;
}
