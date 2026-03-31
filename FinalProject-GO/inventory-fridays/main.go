package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"inventory-fridays/inventory"
	"inventory-fridays/products"
	"inventory-fridays/providers"
)

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	templateDir := "templates"
	tmpl, err := template.ParseFiles(
		filepath.Join(templateDir, "header.html"),
		filepath.Join(templateDir, "footer.html"),
		filepath.Join(templateDir, tmplName),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, tmplName, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "index.html", nil)
}

func productsViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "products.html", products.GetAll())
}

func productsEditViewHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	p := products.GetByID(id)
	if p == nil {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}
	renderTemplate(w, "product_edit.html", p)
}

func providersViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "providers.html", providers.GetAll())
}

func providersEditViewHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	p := providers.GetByID(id)
	if p == nil {
		http.Redirect(w, r, "/providers", http.StatusSeeOther)
		return
	}
	renderTemplate(w, "provider_edit.html", p)
}

func inventoryViewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "inventory.html", inventory.GetFullInventory())
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/products", productsViewHandler)
	mux.HandleFunc("/products/create", products.HandleFormCreate)
	mux.HandleFunc("/products/edit", productsEditViewHandler)
	mux.HandleFunc("/products/update", products.HandleFormUpdate)
	mux.HandleFunc("/products/delete", products.HandleFormDelete)

	mux.HandleFunc("/providers", providersViewHandler)
	mux.HandleFunc("/providers/create", providers.HandleFormCreate)
	mux.HandleFunc("/providers/edit", providersEditViewHandler)
	mux.HandleFunc("/providers/update", providers.HandleFormUpdate)
	mux.HandleFunc("/providers/delete", providers.HandleFormDelete)

	mux.HandleFunc("/inventory", inventoryViewHandler)
	mux.HandleFunc("/inventory/update", inventory.HandleFormUpdate)

	mux.HandleFunc("/api/products", products.HandleAPI)
	mux.HandleFunc("/api/inventory", inventory.HandleAPI)
	mux.HandleFunc("/api/providers", providers.HandleAPI)

	log.Println("Servidor iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
