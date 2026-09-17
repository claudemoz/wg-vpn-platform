import { ActivityIndicator, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';

import { Colors, Spacing } from '@/constants/theme';
import type { Server } from '@/types/api';

type Props = {
  servers: Server[];
  selected: Server | null;
  onSelect: (server: Server) => void;
  disabled?: boolean;
  loading?: boolean;
};

export function ServerPicker({ servers, selected, onSelect, disabled, loading }: Props) {
  if (loading) {
    return (
      <View style={styles.loading}>
        <ActivityIndicator color="#208AEF" />
        <Text style={styles.loadingText}>Chargement des serveurs…</Text>
      </View>
    );
  }

  if (servers.length === 0) {
    return (
      <View style={styles.empty}>
        <Text style={styles.emptyText}>Aucun serveur disponible</Text>
      </View>
    );
  }

  return (
    <ScrollView
      horizontal
      showsHorizontalScrollIndicator={false}
      contentContainerStyle={styles.list}>
      {servers.map((server) => {
        const isSelected = selected?.id === server.id;
        return (
          <Pressable
            key={server.id}
            disabled={disabled}
            onPress={() => onSelect(server)}
            style={[
              styles.card,
              isSelected && styles.cardSelected,
              disabled && styles.cardDisabled,
            ]}>
            <Text style={styles.region}>{server.region.toUpperCase()}</Text>
            <Text style={styles.name} numberOfLines={1}>
              {server.name}
            </Text>
            <Text style={styles.peers}>
              {server.current_peers}/{server.max_peers} peers
            </Text>
          </Pressable>
        );
      })}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  list: {
    gap: Spacing.two,
    paddingHorizontal: Spacing.one,
  },
  card: {
    backgroundColor: Colors.dark.backgroundElement,
    borderRadius: 14,
    padding: Spacing.three,
    minWidth: 140,
    borderWidth: 2,
    borderColor: 'transparent',
  },
  cardSelected: {
    borderColor: '#208AEF',
    backgroundColor: '#208AEF18',
  },
  cardDisabled: {
    opacity: 0.6,
  },
  region: {
    fontSize: 12,
    fontWeight: '700',
    color: '#208AEF',
    letterSpacing: 1,
  },
  name: {
    fontSize: 15,
    fontWeight: '600',
    color: Colors.dark.text,
    marginTop: 4,
  },
  peers: {
    fontSize: 12,
    color: Colors.dark.textSecondary,
    marginTop: 4,
  },
  loading: {
    alignItems: 'center',
    gap: Spacing.two,
    padding: Spacing.four,
  },
  loadingText: {
    color: Colors.dark.textSecondary,
    fontSize: 14,
  },
  empty: {
    padding: Spacing.four,
    alignItems: 'center',
  },
  emptyText: {
    color: Colors.dark.textSecondary,
    fontSize: 14,
  },
});
