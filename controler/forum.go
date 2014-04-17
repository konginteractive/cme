package controler

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	M "github.com/konginteractive/cme/model"
	"log"
	"math"
	. "strconv"
	"strings"
)

var ForumTempl = "forum"
var ForumAddTempl = "forum_add"

var maxElementsInPage = 10

// permet d'afficher la liste des questions du forum
func ForumView() M.Page {

	log.Println("ForumView appelé")

	p := new(M.PageForum)
	p.Title = "Forum"
	p.MainClass = "forum"
	p.PageLevel = ""
	p.Forums = getListForums()
	p.Categories = getAllFormCategories()
	p.PagesList = createPaginate()

	injectDataForumToDisplay(p.Forums)

	return p
}

// permet d'afficher la liste des questions du forum avec la catégorie correspondante
func FormViewCategory(cat string) M.Page {
	log.Println("ForumView appelé")

	// récupère l'id de la catégorie
	idCat := getIdFromCatName(cat)

	p := new(M.PageForum)
	p.Title = "Forum " + cat
	p.MainClass = "forum"
	p.PageLevel = "../"
	p.Forums = getListFormusFromCat(idCat)
	p.Categories = getAllFormCategories()
	p.PagesList = createPaginateFromIdCat(idCat)

	injectDataForumToDisplay(p.Forums)

	return p
}

// permet d'afficher le formulaire de création d'une question du formulaire
func ForumAddView() M.Page {

	log.Println("ForumAddView appelé")

	p := new(M.PageForum)
	p.Title = "Titre du sujet"
	p.MainClass = "nouveaututo"
	p.PageLevel = "../"

	return p
}

// permet de donner l'id d'une question à partir de son titre
func getIdFromCatName(cat string) int64 {
	catList := getAllFormCategories()
	lenCat := len(cat)
	for i := 0; i < lenCat; i++ {
		// vérifie si la cat est celle cherchée
		if strings.ToLower(catList[i].Title) == cat {
			return catList[i].Id
		}
	}
	// par défaut retourne la cat 1
	return 1
}

// permet de retourner toutes les catégories
// permet aussi de créer les liens pour les catégories
func getAllFormCategories() []M.ForumCategory {
	db := connectToDatabase()
	var cat []M.ForumCategory
	db.Find(&cat)

	// create dynamic links
	lenCat := len(cat)
	for i := 0; i < lenCat; i++ {
		//cat[i].Url = "/forum/categorie/" + strings.ToLower(cat[i].Title) + "/" + Itoa(int(cat[i].Id))
		cat[i].Url = "/forum/" + strings.ToLower(cat[i].Title)
	}

	return cat
}

// Permet de retrouver le nombre de réponses pour chaque post
// Permet aussi de réduire la description du texte de desc à 250 caractères
func injectDataForumToDisplay(forums []M.Forum) []M.Forum {
	lenForum := len(forums)

	for i := 0; i < lenForum; i++ {
		id := forums[i].Id
		text := forums[i].Text[0:250]
		forums[i].PostNumb = getNumPostForum(id)
		forums[i].Text = text
	}

	return forums
}

// permet de récupérer toute la listes des questions du forum
// en fonction de la limite affichable par page
func getListForums() []M.Forum {
	db := connectToDatabase()
	var forums []M.Forum
	db.Limit(maxElementsInPage).Where("is_online = ?", "1").Find(&forums)
	return forums
}

// permet de récupérer des questions du forum à partir de l'id de la catégorie
func getListFormusFromCat(id int64) []M.Forum {
	db := connectToDatabase()
	var forums []M.Forum
	db.Limit(maxElementsInPage).Where("is_online = ? and forum_category_id = ?", "1", Itoa(int(id))).Find(&forums)
	return forums
}

// permet de récupérer le nombre de forums de question total de la base de donnée
func getNumForms() int {
	db := connectToDatabase()
	var forums []M.Forum
	var num int
	db.Where("is_online = ?", "1").Find(&forums).Count(&num)
	return num
}

// permet de récupérer le nombre de forums de question total de la base de donnée
// en fonction de l'id de la catégorie
func getNumFormsFromIdCat(id int64) int {
	db := connectToDatabase()
	var forums []M.Forum
	var num int
	db.Where("is_online = ? and forum_category_id = ?", "1", Itoa(int(id))).Find(&forums).Count(&num)
	return num
}

// permet de récupérer les posts d'un forum
// à partir de l'id d'un forum
func getPostForum(id int) []M.ForumPost {
	idForum := Itoa(id)
	db := connectToDatabase()
	var posts []M.ForumPost
	db.Where("is_online = ? and forum_id = ?", "1", idForum).Find(&posts)
	return posts
}

// permet de récupérer le nombre de posts d'un forum
// à partir de l'id d'un forum
func getNumPostForum(id int64) int64 {
	i := int(id)
	idForum := Itoa(i)
	db := connectToDatabase()
	var posts []M.ForumPost
	var num int64
	db.Where("is_online = ? and forum_id = ?", "1", idForum).Find(&posts).Count(&num)
	return num
}

// fonction permettant de se connecter à la base de donnée
func connectToDatabase() gorm.DB {
	db, _ := gorm.Open("mysql", "root:root@tcp(127.0.0.1:8889)/cme_test?charset=utf8&parseTime=True")
	db.SingularTable(true)
	return db
}

// fonction pour créer la pagination
func createPaginate() []M.Paginate {
	elTotal := getNumForms()

	nb := elTotal / maxElementsInPage
	mf := int(math.Floor(float64(nb)))
	p := make([]M.Paginate, nb)

	for i := 0; i < nb; i++ {
		t := Itoa(i + 1)
		p[i].Title = t
		p[i].Url = "/forum/p/" + t
	}
	return p
}

// fonction pour créer la pagination à partir d'une catégorie sélectionnée
func createPaginateFromIdCat(id int64) []M.Paginate {
	elTotal := getNumFormsFromIdCat(id)

	nb := elTotal / maxElementsInPage
	mf := int(math.Floor(float64(nb)))
	p := make([]M.Paginate, nb)

	for i := 0; i < nb; i++ {
		t := Itoa(i + 1)
		p[i].Title = t
		p[i].Url = "/forum/p/" + t
	}
	return p
}

/*
// fonction permetttant de rechercher dans les titres des questions
func searchInTitle(s string) []M.Forum {
	db := connectToDatabase()
	var forums []M.Forum
	db.Table(&forums).Where("title = ?", s)
	return forums
}
*/
