import { Link, router } from 'expo-router';
import { useRef, useState } from 'react';
import {
  Alert,
  Keyboard,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  StyleSheet,
  Text,
  View,
  type ScrollView as ScrollViewType,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { useAuth } from '@/contexts/AuthContext';
import { getHttpBaseURL } from '@/constants/config';
import { ApiClientError } from '@/services/api';
import { Colors, Spacing } from '@/constants/theme';

function showError(message: string) {
  if (Platform.OS === 'web') {
    window.alert(message);
  } else {
    Alert.alert('Inscription', message);
  }
}

export default function RegisterScreen() {
  const scrollRef = useRef<ScrollViewType>(null);
  const { register } = useAuth();
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleRegister() {
    Keyboard.dismiss();
    setError(null);

    if (!firstName.trim() || !lastName.trim() || !email.trim() || !password) {
      const msg = 'Tous les champs obligatoires doivent être remplis';
      setError(msg);
      showError(msg);
      return;
    }
    if (password.length < 8) {
      const msg = 'Le mot de passe doit contenir au moins 8 caractères';
      setError(msg);
      showError(msg);
      return;
    }

    setLoading(true);
    try {
      await register({
        first_name: firstName,
        last_name: lastName,
        email,
        phone: phone || undefined,
        password,
      });
      router.replace({ pathname: '/(auth)/login', params: { registered: '1' } });
    } catch (e) {
      const msg =
        e instanceof ApiClientError
          ? e.message
          : e instanceof Error
            ? e.message
            : 'Inscription impossible. Vérifiez l’URL de l’API.';
      setError(msg);
      showError(msg);
      scrollRef.current?.scrollToEnd({ animated: true });
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
          ref={scrollRef}
          contentContainerStyle={styles.content}
          keyboardShouldPersistTaps="always"
          keyboardDismissMode="on-drag">
          <View style={styles.header}>
            <Text style={styles.title}>Créer un compte</Text>
            <Text style={styles.subtitle}>Rejoignez la plateforme VPN</Text>
          </View>

          <View style={styles.form}>
            <View style={styles.row}>
              <View style={styles.half}>
                <Input
                  label="Prénom"
                  value={firstName}
                  onChangeText={setFirstName}
                  autoComplete="given-name"
                  placeholder="Jean"
                />
              </View>
              <View style={styles.half}>
                <Input
                  label="Nom"
                  value={lastName}
                  onChangeText={setLastName}
                  autoComplete="family-name"
                  placeholder="Dupont"
                />
              </View>
            </View>

            <Input
              label="Email"
              value={email}
              onChangeText={setEmail}
              keyboardType="email-address"
              autoComplete="email"
              placeholder="vous@exemple.com"
            />
            <Input
              label="Téléphone (optionnel)"
              value={phone}
              onChangeText={setPhone}
              keyboardType="phone-pad"
              placeholder="+33612345678"
            />
            <Input
              label="Mot de passe"
              value={password}
              onChangeText={setPassword}
              secureTextEntry
              autoComplete="new-password"
              placeholder="8 caractères minimum"
            />

            {error ? <Text style={styles.error}>{error}</Text> : null}

            <Button title="S'inscrire" onPress={handleRegister} loading={loading} />

            {__DEV__ ? (
              <Text style={styles.devHint}>API : {getHttpBaseURL()}</Text>
            ) : null}
          </View>

          <View style={styles.footer}>
            <Text style={styles.footerText}>Déjà un compte ?</Text>
            <Link href="/(auth)/login" asChild>
              <Text style={styles.link}>Se connecter</Text>
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
    paddingTop: Spacing.five,
    gap: Spacing.four,
  },
  header: {
    gap: Spacing.one,
  },
  title: {
    fontSize: 28,
    fontWeight: '700',
    color: Colors.dark.text,
  },
  subtitle: {
    fontSize: 15,
    color: Colors.dark.textSecondary,
  },
  form: {
    gap: Spacing.three,
  },
  row: {
    flexDirection: 'row',
    gap: Spacing.two,
  },
  half: {
    flex: 1,
  },
  error: {
    color: '#E5484D',
    fontSize: 14,
    textAlign: 'center',
    backgroundColor: '#E5484D18',
    padding: Spacing.two,
    borderRadius: 10,
  },
  devHint: {
    fontSize: 11,
    color: Colors.dark.textSecondary,
    textAlign: 'center',
    fontFamily: Platform.OS === 'ios' ? 'Menlo' : 'monospace',
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'center',
    gap: Spacing.one,
    paddingBottom: Spacing.four,
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
