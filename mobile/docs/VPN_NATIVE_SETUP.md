# Tunnel VPN natif — configuration

L'app utilise [`react-native-wireguard-vpn`](https://github.com/usama7365/react-native-wireguard-vpn) pour activer le tunnel via :

- **iOS** : Network Extension (Packet Tunnel Provider)
- **Android** : `VpnService` + consentement système

**Expo Go ne fonctionne pas.** Il faut un **development build** (`expo run:ios` / `expo run:android`).

---

## 1. Installation

```bash
cd mobile
npm install
npx expo prebuild --clean
```

---

## 2. Android

Le plugin Expo ajoute automatiquement les permissions VPN dans `AndroidManifest.xml`.

```bash
npm run android
```

Au premier lancement, Android affiche la **demande de consentement VPN** — acceptez-la.

---

## 3. iOS (étapes manuelles Xcode)

### Prérequis

- Compte **Apple Developer payant** ($99/an) — les Network Extensions ne fonctionnent pas avec un compte gratuit
- **Appareil physique** (le simulateur ne supporte pas le VPN)

### A. Prebuild

```bash
npm run ios
# ou
npx expo run:ios --device
```

### B. Créer l'extension Packet Tunnel

1. Ouvrir `ios/WGVPN.xcworkspace` dans Xcode
2. **File → New → Target → Network Extension → Packet Tunnel Provider**
3. Nom : `WireGuardTunnel`
4. **Bundle Identifier** : `com.wireguardvpn.tunnel`  
   (obligatoire — c'est l'ID attendu par la librairie native)

### C. Capabilities

Sur la **cible principale** (app) :

- Signing & Capabilities → **+ Capability** → **Network Extensions**
- Cocher **Packet Tunnel**

Répéter sur la cible **WireGuardTunnel** si nécessaire.

### D. App IDs (Developer Portal)

Pour `com.wgvpn.platform` et `com.wireguardvpn.tunnel` :

- Activer **Network Extensions** → **Packet Tunnel**
- Regénérer les provisioning profiles

### E. Build & run

```bash
npx expo run:ios --device
```

---

## 4. Flux dans l'app

1. Login → sélection serveur → **Se connecter**
2. L'API crée un `device` et renvoie la config wg-quick + clé privée
3. `vpn-tunnel.ts` parse la config et appelle `WireGuardVpnModule.connect()`
4. iOS/Android établit le tunnel WireGuard
5. **Déconnecter** appelle `WireGuardVpnModule.disconnect()`

---

## 5. Dépannage

| Problème | Solution |
|----------|----------|
| Module non lié | `cd ios && pod install`, rebuild complet |
| iOS Code=10 permission denied | Compte dev payant + entitlements Network Extension |
| Android Bad address | Vérifier `Address = x.x.x.x/32` dans la config |
| Tunnel indisponible dans l'app | Pas Expo Go → `npm run ios` |

---

## 6. Fichiers concernés

| Fichier | Rôle |
|---------|------|
| `src/services/vpn-tunnel.ts` | Wrapper natif initialize/connect/disconnect |
| `src/services/wg-config-parser.ts` | `.conf` → `WireGuardConfig` |
| `src/contexts/VpnContext.tsx` | Orchestration API + tunnel |
| `app.json` | Plugin + entitlements iOS + permissions Android |
