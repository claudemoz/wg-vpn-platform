import { Pressable, StyleSheet, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { ConnectButton } from '@/components/vpn/ConnectButton';
import { ServerPicker } from '@/components/vpn/ServerPicker';
import { useAuth } from '@/contexts/AuthContext';
import { useVpn } from '@/contexts/VpnContext';
import { Colors, Spacing } from '@/constants/theme';

const STATUS_TEXT = {
  disconnected: 'Non connecté',
  connecting: 'Connexion en cours…',
  connected: 'Connecté',
} as const;

export default function HomeScreen() {
  const { user, logout } = useAuth();
  const {
    servers,
    selectedServer,
    status,
    assignedIp,
    isLoadingServers,
    tunnelAvailable,
    error,
    selectServer,
    connect,
    disconnect,
  } = useVpn();

  const isConnected = status === 'connected';

  function handleToggle() {
    if (isConnected) {
      disconnect();
    } else {
      connect();
    }
  }

  return (
    <SafeAreaView style={styles.safe}>
      <View style={styles.header}>
        <View>
          <Text style={styles.greeting}>Bonjour,</Text>
          <Text style={styles.name}>
            {user?.first_name} {user?.last_name}
          </Text>
        </View>
        <Pressable onPress={logout} style={styles.logout}>
          <Text style={styles.logoutText}>Déconnexion</Text>
        </Pressable>
      </View>

      <View style={styles.statusCard}>
        <View style={[styles.statusDot, isConnected && styles.statusDotOn]} />
        <Text style={styles.statusLabel}>{STATUS_TEXT[status]}</Text>
        {assignedIp ? (
          <Text style={styles.ip}>IP : {assignedIp}</Text>
        ) : null}
        {selectedServer && isConnected ? (
          <Text style={styles.serverInfo}>
            {selectedServer.name} ({selectedServer.region.toUpperCase()})
          </Text>
        ) : null}
      </View>

      <View style={styles.connectSection}>
        <ConnectButton
          status={status}
          onPress={handleToggle}
          disabled={!selectedServer && !isConnected}
        />
      </View>

      <View style={styles.serversSection}>
        <Text style={styles.sectionTitle}>Serveur VPN</Text>
        <ServerPicker
          servers={servers}
          selected={selectedServer}
          onSelect={selectServer}
          disabled={isConnected || status === 'connecting'}
          loading={isLoadingServers}
        />
      </View>

      {error ? <Text style={styles.error}>{error}</Text> : null}

      <Text style={styles.hint}>
        {tunnelAvailable
          ? 'Tunnel WireGuard natif actif (Network Extension iOS / VpnService Android).'
          : 'Expo Go ne supporte pas le VPN natif. Build requis : npm run ios ou npm run android.'}
      </Text>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: {
    flex: 1,
    backgroundColor: Colors.dark.background,
    paddingHorizontal: Spacing.four,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    paddingTop: Spacing.two,
    marginBottom: Spacing.four,
  },
  greeting: {
    fontSize: 14,
    color: Colors.dark.textSecondary,
  },
  name: {
    fontSize: 22,
    fontWeight: '700',
    color: Colors.dark.text,
    marginTop: 2,
  },
  logout: {
    padding: Spacing.two,
  },
  logoutText: {
    color: Colors.dark.textSecondary,
    fontSize: 14,
  },
  statusCard: {
    backgroundColor: Colors.dark.backgroundElement,
    borderRadius: 16,
    padding: Spacing.three,
    alignItems: 'center',
    gap: Spacing.one,
    marginBottom: Spacing.four,
  },
  statusDot: {
    width: 10,
    height: 10,
    borderRadius: 5,
    backgroundColor: Colors.dark.textSecondary,
  },
  statusDotOn: {
    backgroundColor: '#30A46C',
  },
  statusLabel: {
    fontSize: 18,
    fontWeight: '600',
    color: Colors.dark.text,
  },
  ip: {
    fontSize: 14,
    color: '#30A46C',
    fontFamily: 'monospace',
  },
  serverInfo: {
    fontSize: 13,
    color: Colors.dark.textSecondary,
  },
  connectSection: {
    alignItems: 'center',
    paddingVertical: Spacing.five,
  },
  serversSection: {
    gap: Spacing.two,
    marginBottom: Spacing.three,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: Colors.dark.text,
    paddingHorizontal: Spacing.one,
  },
  error: {
    color: '#E5484D',
    textAlign: 'center',
    fontSize: 14,
    marginBottom: Spacing.two,
  },
  hint: {
    fontSize: 12,
    color: Colors.dark.textSecondary,
    textAlign: 'center',
    lineHeight: 18,
    marginTop: 'auto',
    paddingBottom: Spacing.three,
  },
});
