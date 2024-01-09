package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// Structure pour stocker le token d'accès Spotify
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Duration    int    `json:"expires_in"`
}

// Structure pour représenter les données des albums
type AlbumsData struct {
	Name        string
	Image       string
	ReleaseDate string
	Tracks      int
}

// Structure pour représenter les données des titres
type TrackData struct {
	Title       string
	AlbumCover  string
	Album       string
	Artist      string
	ReleaseDate string
	SpotifyLink string
}

// Fonction pour obtenir le token d'accès Spotify
func getAccessToken() (string, error) {
	clientID := "a64cd1f4ec7a4ff39957dc6e42f1c0e5"
	clientSecret := "1a0108afa8974e7683f022621db9effd"

	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	tokenURL := "https://accounts.spotify.com/api/token"
	payload := strings.NewReader("grant_type=client_credentials")

	req, err := http.NewRequest("POST", tokenURL, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Basic "+authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	accessToken, ok := tokenResp["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("Token d'accès introuvable")
	}

	return accessToken, nil
}

func main() {
	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatal("Impossible d'obtenir un token: ", err)
	}

	// Initialisation des templates
	temp, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Erreur dans la récupération des templates : ", err)
		return
	}

	// Serveur de fichiers statiques (pour les styles CSS)
	fileserver := http.FileServer(http.Dir("styles"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	// Gestion des routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "index.html", nil)
	})

	http.HandleFunc("/album/jul", func(w http.ResponseWriter, r *http.Request) {
		apiURL := "https://api.spotify.com/v1/artists/3IW7ScrzXmPvZhB27hmfgy/albums"

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Println("Erreur dans la requête d'albums JUL :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		req.Header.Add("Authorization", "Bearer "+accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Erreur dans la requête HTTP pour les albums JUL :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var albumsData struct {
			Items []struct {
				Name   string `json:"name"`
				Images []struct {
					URL string `json:"url"`
				} `json:"images"`
				ReleaseDate string `json:"release_date"`
				TotalTracks int    `json:"total_tracks"`
			} `json:"items"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&albumsData); err != nil {
			log.Println("Erreur dans la lecture des données d'albums JUL :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		var decodeData []AlbumsData

		for _, album := range albumsData.Items {
			data := AlbumsData{
				Name:        album.Name,
				Image:       album.Images[0].URL,
				ReleaseDate: album.ReleaseDate,
				Tracks:      album.TotalTracks,
			}
			decodeData = append(decodeData, data)
		}

		temp.ExecuteTemplate(w, "jul.html", decodeData)
	})

	http.HandleFunc("/track/sdm", func(w http.ResponseWriter, r *http.Request) {
		apiURL := "https://api.spotify.com/v1/tracks/0EzNyXyU7gHzj2TN8qYThj"

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Println("Erreur dans la requête de Bolide allemand :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		req.Header.Add("Authorization", "Bearer "+accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Erreur dans la requête HTTP pour Bolide allemand :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var trackData TrackData

		if err := json.NewDecoder(resp.Body).Decode(&trackData); err != nil {
			log.Println("Erreur dans la lecture des données de Bolide allemand :", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		temp.ExecuteTemplate(w, "bolideallemand.html", trackData)
	})

	fmt.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
