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

func ajouterContact() {
	form := tview.NewForm()

	form.AddInputField("Nom", "", 30, nil, nil).
		AddInputField("Téléphone", "", 30, nil, nil).
		AddInputField("Email", "", 30, nil, nil).
		AddButton("Ajouter", func() {
			nom := form.GetFormItemByLabel("Nom").(*tview.InputField).GetText()
			telephone := form.GetFormItemByLabel("Téléphone").(*tview.InputField).GetText()
			email := form.GetFormItemByLabel("Email").(*tview.InputField).GetText()

			if nom == "" || telephone == "" || email == "" {
				modal := tview.NewModal().
					SetText("Tous les champs sont requis !").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						app.SetRoot(form, true)
					})
				app.SetRoot(modal, true)
				return
			}

			contact := Contact{Nom: nom, Telephone: telephone, Email: email}
			contacts = append(contacts, contact)

			err := sauvegarderContactsDansXML()
			if err != nil {
				modal := tview.NewModal().
					SetText(fmt.Sprintf("Erreur lors de l'enregistrement : %v", err)).
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						app.SetRoot(form, true)
					})
				app.SetRoot(modal, true)
			} else {
				modal := tview.NewModal().
					SetText("Contact ajouté avec succès !").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						app.Stop()
					})
				app.SetRoot(modal, true)
			}
		}).
		AddButton("Annuler", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Ajouter un Contact").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
	app.Run()
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

func rechercherContact() {
	form := tview.NewForm()
	form.AddInputField("Nom", "", 30, nil, nil).
		AddButton("Rechercher", func() {
			nomRecherche := form.GetFormItemByLabel("Nom").(*tview.InputField).GetText()
			nomRecherche = strings.TrimSpace(nomRecherche)

			if nomRecherche == "" {
				modal := tview.NewModal().
					SetText("Veuillez entrer un nom.").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						app.SetRoot(form, true)
					})
				app.SetRoot(modal, true)
				return
			}

			results := []Contact{}
			for _, contact := range contacts {
				if strings.Contains(strings.ToLower(contact.Nom), strings.ToLower(nomRecherche)) {
					results = append(results, contact)
				}
			}

			if len(results) == 0 {
				modal := tview.NewModal().
					SetText("Aucun contact trouvé.").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						app.SetRoot(form, true)
					})
				app.SetRoot(modal, true)
			} else {
				list := tview.NewList()
				for _, contact := range results {
					list.AddItem(fmt.Sprintf("%s | %s | %s", contact.Nom, contact.Telephone, contact.Email), "", 0, nil)
				}

				list.AddItem("Retour", "", 'r', func() {
					app.SetRoot(form, true)
				})

				list.SetBorder(true).SetTitle("Résultats de la Recherche").SetTitleAlign(tview.AlignLeft)
				app.SetRoot(list, true)
			}
		}).
		AddButton("Annuler", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Rechercher un Contact").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
	app.Run()
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
				modal := tview.NewModal().
					SetText("Contact supprimé avec succès !").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						app.SetRoot(list, true) 
					})
				app.SetRoot(modal, true)
			}
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
			ajouterContact()
		case 2:
			listerContacts()
		case 3:
			rechercherContact()
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
