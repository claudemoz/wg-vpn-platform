import {
  ActivityIndicator,
  Platform,
  Pressable,
  StyleSheet,
  Text,
  type PressableProps,
} from 'react-native';

import { Colors, Spacing } from '@/constants/theme';

type Props = PressableProps & {
  title: string;
  variant?: 'primary' | 'secondary' | 'ghost';
  loading?: boolean;
};

export function Button({
  title,
  variant = 'primary',
  loading,
  disabled,
  style,
  ...rest
}: Props) {
  const isPrimary = variant === 'primary';
  const isGhost = variant === 'ghost';

  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [
        styles.base,
        isPrimary && styles.primary,
        variant === 'secondary' && styles.secondary,
        isGhost && styles.ghost,
        (disabled || loading) && styles.disabled,
        pressed && !disabled && !loading && styles.pressed,
        Platform.OS === 'web' && styles.web,
        style,
      ]}
      disabled={disabled || loading}
      {...rest}>
      {loading ? (
        <ActivityIndicator color={isPrimary ? '#fff' : Colors.light.text} />
      ) : (
        <Text
          style={[
            styles.label,
            isPrimary && styles.labelPrimary,
            isGhost && styles.labelGhost,
          ]}>
          {title}
        </Text>
      )}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  base: {
    paddingVertical: Spacing.three,
    paddingHorizontal: Spacing.four,
    borderRadius: 14,
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: 52,
  },
  primary: {
    backgroundColor: '#208AEF',
  },
  secondary: {
    backgroundColor: Colors.dark.backgroundElement,
  },
  ghost: {
    backgroundColor: 'transparent',
  },
  disabled: {
    opacity: 0.5,
  },
  pressed: {
    opacity: 0.85,
  },
  web: {
    cursor: 'pointer',
  } as object,
  label: {
    fontSize: 16,
    fontWeight: '600',
    color: Colors.light.text,
  },
  labelPrimary: {
    color: '#fff',
  },
  labelGhost: {
    color: '#208AEF',
  },
});
