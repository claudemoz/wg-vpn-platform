import { Link, useLocalSearchParams } from 'expo-router';
import { useState } from 'react';
import {
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { useAuth } from '@/contexts/AuthContext';
import { ApiClientError } from '@/services/api';
import { Colors, Spacing } from '@/constants/theme';

export default function LoginScreen() {
  const { registered } = useLocalSearchParams<{ registered?: string }>();
  const { login } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleLogin() {
    setError(null);
    if (!email.trim() || !password) {
      setError('Email et mot de passe requis');
      return;
    }

    setLoading(true);
    try {
      await login(email, password);
    } catch (e) {
      const msg =
        e instanceof ApiClientError
          ? e.message
          : 'Connexion impossible. Vérifiez l’URL de l’API.';
      setError(msg);
    } finally {
      setLoading(false);
    }
  }

  return (
    <SafeAreaView style={styles.safe}>
      <KeyboardAvoidingView
        style={styles.flex}
        behavior={Platform.OS === 'ios' ? 'padding' : undefined}>
        <ScrollView
          contentContainerStyle={styles.content}
          keyboardShouldPersistTaps="handled">
          <View style={styles.header}>
            <Text style={styles.logo}>WG VPN</Text>
            <Text style={styles.subtitle}>Connectez-vous à votre compte</Text>
          </View>

          {registered === '1' ? (
            <Text style={styles.success}>
              Compte créé. Connectez-vous pour continuer.
            </Text>
          ) : null}

          <View style={styles.form}>
            <Input
              label="Email"
              value={email}
              onChangeText={setEmail}
              keyboardType="email-address"
              autoComplete="email"
              placeholder="vous@exemple.com"
            />
            <Input
              label="Mot de passe"
              value={password}
              onChangeText={setPassword}
              secureTextEntry
              autoComplete="password"
              placeholder="••••••••"
            />

            {error ? <Text style={styles.error}>{error}</Text> : null}

            <Button title="Se connecter" onPress={handleLogin} loading={loading} />
          </View>

          <View style={styles.footer}>
            <Text style={styles.footerText}>Pas encore de compte ?</Text>
            <Link href="/(auth)/register" asChild>
              <Text style={styles.link}>Créer un compte</Text>
            </Link>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: {
    flex: 1,
    backgroundColor: Colors.dark.background,
  },
  flex: {
    flex: 1,
  },
  content: {
    flexGrow: 1,
    padding: Spacing.four,
    justifyContent: 'center',
    gap: Spacing.five,
  },
  header: {
    alignItems: 'center',
    gap: Spacing.two,
  },
  logo: {
    fontSize: 36,
    fontWeight: '800',
    color: Colors.dark.text,
    letterSpacing: -1,
  },
  subtitle: {
    fontSize: 16,
    color: Colors.dark.textSecondary,
  },
  form: {
    gap: Spacing.three,
  },
  success: {
    color: '#30A46C',
    fontSize: 14,
    textAlign: 'center',
    backgroundColor: '#30A46C18',
    padding: Spacing.two,
    borderRadius: 10,
  },
  error: {
    color: '#E5484D',
    fontSize: 14,
    textAlign: 'center',
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'center',
    gap: Spacing.one,
  },
  footerText: {
    color: Colors.dark.textSecondary,
    fontSize: 15,
  },
  link: {
    color: '#208AEF',
    fontSize: 15,
    fontWeight: '600',
  },
});
