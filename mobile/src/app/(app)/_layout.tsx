import { Redirect, Stack } from 'expo-router';

import { useAuth } from '@/contexts/AuthContext';
import { VpnProvider } from '@/contexts/VpnContext';

export default function AppLayout() {
  const { token } = useAuth();

  if (!token) {
    return <Redirect href="/(auth)/login" />;
  }

  return (
    <VpnProvider>
      <Stack screenOptions={{ headerShown: false }}>
        <Stack.Screen name="index" />
      </Stack>
    </VpnProvider>
  );
}
