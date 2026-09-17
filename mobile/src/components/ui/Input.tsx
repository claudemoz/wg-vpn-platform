import { StyleSheet, Text, TextInput, View, type TextInputProps } from 'react-native';

import { Colors, Spacing } from '@/constants/theme';

type Props = TextInputProps & {
  label: string;
  error?: string;
};

export function Input({ label, error, style, ...rest }: Props) {
  return (
    <View style={styles.wrapper}>
      <Text style={styles.label}>{label}</Text>
      <TextInput
        style={[styles.input, error && styles.inputError, style]}
        placeholderTextColor={Colors.dark.textSecondary}
        autoCapitalize="none"
        {...rest}
      />
      {error ? <Text style={styles.error}>{error}</Text> : null}
    </View>
  );
}

const styles = StyleSheet.create({
  wrapper: {
    gap: Spacing.one,
  },
  label: {
    fontSize: 14,
    fontWeight: '500',
    color: Colors.dark.textSecondary,
  },
  input: {
    backgroundColor: Colors.dark.backgroundElement,
    borderRadius: 12,
    paddingHorizontal: Spacing.three,
    paddingVertical: Spacing.two + 4,
    fontSize: 16,
    color: Colors.dark.text,
    borderWidth: 1,
    borderColor: 'transparent',
  },
  inputError: {
    borderColor: '#E5484D',
  },
  error: {
    fontSize: 13,
    color: '#E5484D',
  },
});
