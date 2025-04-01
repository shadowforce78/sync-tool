# 📂 SushiSync - Synchronisation de Fichiers

## 🚀 Description

Cet outil ultra-léger permet de synchroniser des fichiers entre plusieurs machines via un serveur auto-hébergé. Pas de dépendance à un cloud externe, juste **performance et contrôle total**.

## 🎯 Fonctionnalités

- 📡 **Client CLI** : Interface en ligne de commande pour gérer la synchro.
- 🔄 **Détection des changements** : Ajouts, suppressions, modifications détectés automatiquement.
- 📥 **Stockage minimaliste** : Gestion optimisée de l'espace disque (~10GB max).
- 🔒 **Sécurité** : Authentification par clé API, chiffrement des fichiers.
- ⚡ **Performances optimisées** : Compression des fichiers, communication rapide avec le serveur.

## 🏗️ Architecture

1. **Client CLI** : Interface utilisateur et envoi des fichiers.
2. **Serveur API** : Réception et stockage sécurisé des fichiers.
3. **Stockage** : Gestion optimisée avec compression et suppression automatique.
4. **Auth & Sécurité** : Clé API, chiffrement, communication sécurisée.

## 🛠️ Technologies

- **Go** 🟡 : Backend & client CLI
- **SQLite/PostgreSQL** 🗄️ : Stockage des métadonnées
- **REST/gRPC** 🔌 : Communication entre client & serveur
- **TLS & Chiffrement** 🔐 : Sécurisation des échanges

## 📌 Installation

### 1️⃣ Cloner le projet

```sh
git clone https://github.com/shadowforce78/sync-tool.git
cd sync-tool
```

### 2️⃣ Installer Go (si ce n'est pas déjà fait)

[🔗 Télécharger Go](https://golang.org/dl/)

### 3️⃣ Compiler et exécuter

#### Serveur :

```sh
go run server/main.go
```

#### Client :

```sh
go run client/main.go
```

## Composants

- **Serveur** : Service GO qui gère l'upload, le download et la liste des fichiers
- **Client CLI** : Interface en ligne de commande pour interagir avec le serveur
- **Client GUI** : Interface graphique pour interagir avec le serveur (nécessite compilation manuelle)

## Déploiement

### Serveur

Pour démarrer le serveur :

```bash
go run main.go
```

Le serveur écoute sur http://localhost:8080

### Client CLI

Le client CLI est compilé automatiquement par GitHub Actions et disponible en téléchargement dans les releases.

Pour l'exécuter :

```bash
# Windows
SushiSyncCLI-windows-amd64.exe

# Linux
chmod +x SushiSyncCLI-linux-amd64
./SushiSyncCLI-linux-amd64

# macOS
chmod +x SushiSyncCLI-darwin-amd64
./SushiSyncCLI-darwin-amd64
```

### Client GUI

Le client GUI doit être compilé manuellement en raison des dépendances graphiques.

#### Sur Windows:
```bash
cd client
set CGO_ENABLED=0
go build -tags nocgo -o SushiSyncGUI.exe
```

#### Sur Linux:
```bash
cd client
CGO_ENABLED=0 go build -tags nocgo -o SushiSyncGUI
```

#### Sur macOS:
```bash
cd client
CGO_ENABLED=0 go build -tags nocgo -o SushiSyncGUI
```

Vous pouvez aussi utiliser le script `build.bat` fourni sur Windows pour compiler les deux clients:
```
client\build.bat
```

## Création d'une nouvelle version

Pour créer une nouvelle version et déclencher le workflow de release :

```bash
git tag v1.0.0
git push origin v1.0.0
```

Les versions suivantes incrémenteront le numéro (v1.0.1, v1.1.0, v2.0.0, etc.)

## 🛣️ Roadmap

- [x] **Phase 1** : Mise en place du projet et POC
- [ ] **Phase 2** : Détection des changements & synchro de fichiers
- [ ] **Phase 3** : Sécurisation & optimisation
- [ ] **Phase 4** : UI avancée et fonctionnalités bonus
