package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Contact struct {
	Nom       string `xml:"Nom"`
	Telephone string `xml:"Telephone"`
	Email     string `xml:"Email"`
}

type Contacts struct {
	List []Contact `xml:"Contact"`
}

var contacts []Contact

func chargerContactsDepuisXML() {
	file, err := os.Open("contact.xml")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier XML:", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier XML:", err)
		return
	}

	var parsedContacts Contacts
	err = xml.Unmarshal(data, &parsedContacts)
	if err != nil {
		fmt.Println("Erreur lors du parsing du fichier XML:", err)
		return
	}

	contacts = parsedContacts.List
	fmt.Println("Contacts chargés avec succès depuis le fichier XML !")
}

func afficherMenu() {
	fmt.Println("\nMenu:")
	fmt.Println("1. Ajouter un contact")
	fmt.Println("2. Lister les contacts")
	fmt.Println("3. Rechercher un contact")
	fmt.Println("4. Quitter")
}

func ajouterContact() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Nom: ")
	nom, _ := reader.ReadString('\n')
	nom = strings.TrimSpace(nom)

	fmt.Print("Téléphone: ")
	telephone, _ := reader.ReadString('\n')
	telephone = strings.TrimSpace(telephone)

	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	contact := Contact{Nom: nom, Telephone: telephone, Email: email}
	contacts = append(contacts, contact)

	fmt.Println("Contact ajouté avec succès !")
}

func listerContacts() {
	if len(contacts) == 0 {
		fmt.Println("Aucun contact trouvé.")
		return
	}

	for i, contact := range contacts {
		fmt.Printf("%d. Nom: %s | Téléphone: %s | Email: %s\n", i+1, contact.Nom, contact.Telephone, contact.Email)
	}
}

func rechercherContact() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Entrez le nom à rechercher: ")
	nomRecherche, _ := reader.ReadString('\n')
	nomRecherche = strings.TrimSpace(nomRecherche)

	trouve := false
	for _, contact := range contacts {
		if strings.Contains(strings.ToLower(contact.Nom), strings.ToLower(nomRecherche)) {
			fmt.Printf("Contact trouvé : Nom: %s | Téléphone: %s | Email: %s\n", contact.Nom, contact.Telephone, contact.Email)
			trouve = true
		}
	}

	if !trouve {
		fmt.Println("Aucun contact trouvé avec ce nom.")
	}
}

func main() {
	chargerContactsDepuisXML()

	for {
		afficherMenu()

		var choix int
		fmt.Print("Choisissez une option: ")
		fmt.Scan(&choix)

		switch choix {
		case 1:
			ajouterContact()
		case 2:
			listerContacts()
		case 3:
			rechercherContact()
		case 4:
			fmt.Println("Au revoir!")
			return
		default:
			fmt.Println("Option invalide. Essayez encore.")
		}
	}
}
