# Utiliser l'image officielle de MySQL
FROM mysql:latest

# Définir des variables d'environnement pour la configuration de MySQL
ENV MYSQL_ROOT_PASSWORD=root_password
ENV MYSQL_DATABASE=my_database

# Exposer le port MySQL
EXPOSE 3306

# Lancer MySQL
CMD ["mysqld"]
