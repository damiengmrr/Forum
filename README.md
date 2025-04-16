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

## ScreenShot

-- **Ecran d'accueil**


![ecran-dacceuil](https://github.com/user-attachments/assets/6213ab19-27f7-466a-aba5-fb25a2683c24)


-- **Ecran de connexion**


![connexion](https://github.com/user-attachments/assets/04c3174c-5635-4742-8dbc-c5626df48309)


-- **Ecran d'inscription**


![inscription](https://github.com/user-attachments/assets/5e93b71d-2f48-41ec-bf7f-d5ad91911132)


-- **Choix des catégories**


![categorie](https://github.com/user-attachments/assets/8c499d45-e510-41ff-8d3d-5da52c8ff083)


-- **Contact**


![contact](https://github.com/user-attachments/assets/9e2ec97f-9074-4f53-bf9e-6ff5c539f88b)


-- **Compte**


![compte](https://github.com/user-attachments/assets/44ce38f1-7a4d-4934-a996-265adfffc951)


-- **Création d'un post**


![post](https://github.com/user-attachments/assets/429ee428-2fa1-43dd-8f9c-640556503c9c)


-- **Détail d'un post**


![detailspost](https://github.com/user-attachments/assets/0b124352-bf46-4e29-b812-4988b06844ce)


-- **Modification du profil**


![profil](https://github.com/user-attachments/assets/022a298d-cdbf-4163-b459-35f8bc85c12e)



Bonne visite !

