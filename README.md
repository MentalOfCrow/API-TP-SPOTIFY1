﻿# TP API SPOTIFY 

//Creation dans le Terminal (Rappel : echo "# Nom Du Projet" > README.md) = echo "# TP API SPOTIFY" > README.md
N'oubliez pas d'installer : go get -u github.com/gin-gonic/gin = Bibliotheque ("net/http")
Cela peut etre utile : go get github.com/joho/godotenv = Pour dans le main.go : "godotenv" l-13
Ca aussi c'est utile : go get golang.org/x/oauth2 (dans le terminal l'installer)
Important installer : go get github.com/zmb3/spotify


Structure Projet :
TP_API_SPOTIFY/
project_root/
|-- requetes/
|   |-- requetes.http   // Fichier contenant les requêtes HTTP pour la partie I
|-- site web/
|   |-- main.go         // Code principal pour la partie II
|   |-- templates/
|   |   |-- index.html
|   |   |-- jul_albums.html
|   |   |-- sdm_track_details.html
|   |-- styles/
|   |   |-- style.css

Seulement 1 Terminal Pour le Lancer : 
cd .\API-TP-SPOTIFY1\
cd '.\site web\'
go run .\main.go
Puis http://localhost:8080
