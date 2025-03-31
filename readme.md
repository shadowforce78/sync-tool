# 📂 Projet de Synchronisation de Fichiers

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

## 🛣️ Roadmap

- [ ] **Phase 1** : Mise en place du projet et POC
- [ ] **Phase 2** : Détection des changements & synchro de fichiers
- [ ] **Phase 3** : Sécurisation & optimisation
- [ ] **Phase 4** : UI avancée et fonctionnalités bonus
