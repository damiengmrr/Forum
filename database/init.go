package database

import (
    "forum/models"
    "time"
)

// Stockage en mémoire pour simplifier
var Posts = []models.Post{
    {ID: 1, Author: "Alice", Content: "Bienvenue sur le forum !", Category: "Discussion", Date: time.Now()},
    {ID: 2, Author: "Bob", Content: "Les nouvelles technologies avancent vite.", Category: "Technologie", Date: time.Now()},
    {ID: 3, Author: "Charlie", Content: "Quel est votre jeu préféré ?", Category: "Jeux vidéo", Date: time.Now()},
    {ID: 4, Author: "Dave", Content: "Quel est votre livre du moment ?", Category: "Littérature", Date: time.Now()},
}