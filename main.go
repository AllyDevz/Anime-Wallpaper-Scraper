package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	maxPages := 200
	downloadedImages := make(map[int]bool)
	fmt.Println("Created by AllyDevz\n Website scraped:\nhttps://wallhaven.cc/")
	// Verifique se a pasta "wallpapers" existe, se não, crie-a
	err := createWallpapersFolder()
	if err != nil {
		fmt.Println("Erro ao criar a pasta wallpapers:", err)
		return
	}

	for page := 1; page <= maxPages; page++ {
		url := fmt.Sprintf("https://wallhaven.cc/search?q=anime&categories=010&purity=100&atleast=3840x1600&sorting=views&order=desc&ai_art_filter=1&page=%d", page)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Erro ao fazer solicitação HTTP:", err)
			return
		}
		defer resp.Body.Close()

		doc, err := html.Parse(resp.Body)
		if err != nil {
			fmt.Println("Erro ao fazer parsing do HTML:", err)
			return
		}

		var imageURLs []string
		findImageURLs(doc, &imageURLs)

		// Baixe e salve as imagens
		for i, imageURL := range imageURLs {
			// Certifique-se de que a URL é válida (pode ser necessário ajustar isso dependendo da estrutura do site)
			if !strings.HasPrefix(imageURL, "http") {
				continue
			}

			// Verifique se o número já foi baixado
			imageNumber := i + 1 + (page-1)*len(imageURLs) // Cálculo para obter um número único para cada imagem em todas as páginas
			if downloadedImages[imageNumber] {
				fmt.Printf("Imagem %d já foi baixada. Ignorando.\n", imageNumber)
				continue
			}

			// Marque o número como baixado
			downloadedImages[imageNumber] = true

			imageResp, err := http.Get(imageURL)
			if err != nil {
				fmt.Printf("Erro ao baixar imagem %d: %v\n", imageNumber, err)
				continue
			}
			defer imageResp.Body.Close()

			imageData, err := ioutil.ReadAll(imageResp.Body)
			if err != nil {
				fmt.Printf("Erro ao ler dados da imagem %d: %v\n", imageNumber, err)
				continue
			}

			// Salve a imagem localmente
			err = ioutil.WriteFile(fmt.Sprintf("./wallpapers/imagem_%d.jpg", imageNumber), imageData, 0644)
			if err != nil {
				fmt.Printf("Erro ao salvar imagem %d: %v\n", imageNumber, err)
			} else {
				fmt.Printf("Imagem %d baixada e salva.\n", imageNumber)
			}
		}
	}
}

func findImageURLs(n *html.Node, imageURLs *[]string) {
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, attr := range n.Attr {
			if attr.Key == "data-src" {
				*imageURLs = append(*imageURLs, attr.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findImageURLs(c, imageURLs)
	}
}

func createWallpapersFolder() error {
	_, err := os.Stat("./wallpapers")
	if os.IsNotExist(err) {
		err := os.Mkdir("./wallpapers", 0755)
		if err != nil {
			return err
		}
		fmt.Println("Pasta wallpapers criada.")
	}
	return nil
}
