import { ActivityIndicator, Pressable, StyleSheet, Text, View } from 'react-native';

import type { ConnectionStatus } from '@/contexts/VpnContext';

type Props = {
  status: ConnectionStatus;
  onPress: () => void;
  disabled?: boolean;
};

const STATUS_LABEL: Record<ConnectionStatus, string> = {
  disconnected: 'Se connecter',
  connecting: 'Connexion…',
  connected: 'Déconnecter',
};

const STATUS_COLOR: Record<ConnectionStatus, string> = {
  disconnected: '#208AEF',
  connecting: '#F5A623',
  connected: '#30A46C',
};

export function ConnectButton({ status, onPress, disabled }: Props) {
  const color = STATUS_COLOR[status];
  const isConnecting = status === 'connecting';

  return (
    <Pressable
      onPress={onPress}
      disabled={disabled || isConnecting}
      style={({ pressed }) => [styles.wrapper, pressed && styles.pressed]}>
      <View style={[styles.ring, { borderColor: color + '40' }]}>
        <View style={[styles.circle, { backgroundColor: color }]}>
          {isConnecting ? (
            <ActivityIndicator size="large" color="#fff" />
          ) : (
            <View style={styles.powerIcon}>
              <View style={styles.powerStem} />
              <View style={styles.powerArc} />
            </View>
          )}
        </View>
      </View>
      <Text style={[styles.label, { color }]}>{STATUS_LABEL[status]}</Text>
    </Pressable>
  );
}

const SIZE = 160;

const styles = StyleSheet.create({
  wrapper: {
    alignItems: 'center',
    gap: 16,
  },
  pressed: {
    opacity: 0.9,
    transform: [{ scale: 0.98 }],
  },
  ring: {
    width: SIZE + 24,
    height: SIZE + 24,
    borderRadius: (SIZE + 24) / 2,
    borderWidth: 3,
    alignItems: 'center',
    justifyContent: 'center',
  },
  circle: {
    width: SIZE,
    height: SIZE,
    borderRadius: SIZE / 2,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.25,
    shadowRadius: 16,
    elevation: 8,
  },
  powerIcon: {
    width: 48,
    height: 48,
    alignItems: 'center',
    justifyContent: 'center',
  },
  powerStem: {
    width: 4,
    height: 22,
    backgroundColor: '#fff',
    borderRadius: 2,
    marginBottom: 2,
  },
  powerArc: {
    width: 32,
    height: 32,
    borderWidth: 4,
    borderColor: '#fff',
    borderTopColor: 'transparent',
    borderRadius: 16,
    position: 'absolute',
    top: 8,
  },
  label: {
    fontSize: 18,
    fontWeight: '600',
  },
});
