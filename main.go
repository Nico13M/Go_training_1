package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rivo/tview"
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
var app *tview.Application

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
}

func sauvegarderContactsDansXML() error {
	contactsWrapper := Contacts{List: contacts}
	data, err := xml.MarshalIndent(contactsWrapper, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur lors de la conversion en XML : %v", err)
	}

	data = append([]byte(xml.Header), data...)

	err = ioutil.WriteFile("contact.xml", data, 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture dans le fichier XML : %v", err)
	}

	return nil
}

func afficherMenu() {
	fmt.Println("\nMenu:")
	fmt.Println("1. Ajouter un contact")
	fmt.Println("2. Lister les contacts")
	fmt.Println("3. Rechercher un contact")
	fmt.Println("4. Supprimer un contact")
	fmt.Println("5. Quitter")
}

func ajouterContact(reader *bufio.Reader) {
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

	err := sauvegarderContactsDansXML()
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement du contact :", err)
		return
	}

	fmt.Println("Contact ajouté et enregistré avec succès !")
}

func listerContacts() {
	if len(contacts) == 0 {
		modal := tview.NewModal().
			SetText("Aucun contact trouvé.").
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.Stop()
			})

		app.SetRoot(modal, true)
		return
	}

	list := tview.NewList()
	for _, contact := range contacts {
		list.AddItem(fmt.Sprintf("%s | %s | %s", contact.Nom, contact.Telephone, contact.Email), "", 0, nil)
	}

	list.AddItem("Quitter", "", 'q', func() {
		app.Stop()
	})

	list.SetBorder(true).SetTitle("Liste des Contacts").SetTitleAlign(tview.AlignLeft)

	app.SetRoot(list, true)
	app.Run()
}

func rechercherContact(reader *bufio.Reader) {
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

func supprimerContact() {
	if len(contacts) == 0 {
		fmt.Println("Aucun contact à supprimer.")
		return
	}

	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(true).SetTitle("Supprimer un Contact").SetTitleAlign(tview.AlignLeft)

	for i, contact := range contacts {
		contact := contact // Pour éviter des problèmes de capture dans les closures
		index := i         // Capturer l'index
		list.AddItem(fmt.Sprintf("%s | %s | %s", contact.Nom, contact.Telephone, contact.Email), "", 0, func() {
			contacts = append(contacts[:index], contacts[index+1:]...)
			err := sauvegarderContactsDansXML()
			if err != nil {
				fmt.Println("Erreur lors de la sauvegarde après suppression :", err)
			} else {
				fmt.Println("Contact supprimé avec succès !")
			}
			app.Stop()
		})
	}

	list.AddItem("Retour", "", 'r', func() {
		app.Stop()
	})

	app.SetRoot(list, true)
	app.Run()
}

func main() {
	chargerContactsDepuisXML()
	app = tview.NewApplication()
	reader := bufio.NewReader(os.Stdin)

	for {
		afficherMenu()

		fmt.Print("Choisissez une option: ")
		choixStr, _ := reader.ReadString('\n')
		choixStr = strings.TrimSpace(choixStr)

		choix := 0
		fmt.Sscanf(choixStr, "%d", &choix)

		switch choix {
		case 1:
			ajouterContact(reader)
		case 2:
			listerContacts()
		case 3:
			rechercherContact(reader)
		case 4:
			supprimerContact()
		case 5:
			fmt.Println("Au revoir!")
			return
		default:
			fmt.Println("Option invalide. Essayez encore.")
		}
	}
}
