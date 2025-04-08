# 🧵 Forum Prime -

Bienvenue sur **Forum Prime**, un forum web développé dans le cadre du module *Infrastructure & Système d’Information* à Ynov Montpellier (Bachelor 1).

---

## 🧠 Objectif du projet

Créer une **application web complète** avec :
- Un backend codé en **Go**
- Un front dynamique en **HTML/CSS**
- Une structuration **POO** claire
- Des interactions : création de post, likes, commentaires, réponses...
- Une base **SQLite** pour la persistance des comptes

---

## ✨ Fonctionnalités principales

| 🧩 Fonction | ✅ Description |
|------------|----------------|
| 📝 Création de posts | Ajout de titre, contenu, image, catégorie |
| ❤️ Likes / Dislikes | Sur les posts ET les commentaires |
| 💬 Commentaires & Réponses | Affichage clair avec indentation et avatars |
| 📁 Tri par catégories | Discussion, Technologie, Jeux Vidéo, Littérature |
| 👤 Authentification | Inscription, Connexion, Sessions |
| 📆 Dates lisibles | Format simplifié `26 Mar 2025 à 10:15` |

---

## 🛠️ Structure du projet

```
MonForum/
│
├── static/              # CSS, images, icônes
├── templates/           # Fichiers HTML (home, post, etc.)
├── models/              # Structures Go : Post, Comment, User...
├── handlers/            # Fonctions Go (PostHandler, HomeHandler...)
├── main.go              # Point d’entrée du serveur
└── database.go          # Connexion SQLite
```

---

## 🧱 Technologies utilisées

- **Go (Golang)** : serveur, handlers, POO
- **HTML/CSS** : structure et design
- **SQLite** : base de données légère
- **Net/http** : gestion des routes
- **Sessions personnalisées** : pour les comptes utilisateurs

---

## 🎨 Design

- Design propre et clair
- Espacement visuel entre les posts
- Avatars arrondis, structure des commentaires indentée
- Style inspiré de **Twitter** / **Threads** avec une touche minimaliste

---

## 📌 Infos supplémentaires

- Aucune bibliothèque JS externe
- Pas de framework CSS : full custom
- Conformité avec le sujet demandé
- Séparation claire **frontend / backend**
- Prêt pour la démonstration orale

---

## 🚀 Lancer le projet en local

```bash
# 1. Cloner le dépôt
git clone https://github.com/votre-utilisateur/MonForum.git
cd MonForum

# 2. Lancer le serveur (Go doit être installé)
go run main.go
```

👉 Le forum sera disponible sur : [http://localhost:8080/home](http://localhost:8080/home)

---

📂 Si vous utilisez SQLite :
- Assurez-vous que le fichier `forum.db` est bien présent à la racine
- Le script de création est dans `init.sql` (exécutable avec `sqlite3`)

```bash
sqlite3 forum.db < init.sql
```
---

## 👨‍💻 Auteurs

- **Damien**
- **Noah**
- **Théo**
- **Guilhem**

---

Bonne visite !