package products

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var (
	productsMap = make(map[int]Product)
	nextID      = 1
	mu          sync.RWMutex
)

func init() {
	productsMap[nextID] = Product{ID: nextID, Name: "Hamburguesa Clásica", Price: 8.50}
	nextID++
	productsMap[nextID] = Product{ID: nextID, Name: "Papas Fritas", Price: 3.00}
	nextID++
}

func GetAll() []Product {
	mu.RLock()
	defer mu.RUnlock()
	var products []Product
	for _, p := range productsMap {
		products = append(products, p)
	}
	return products
}

func GetByID(id int) *Product {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := productsMap[id]
	if !ok {
		return nil
	}
	return &p
}

func HandleFormCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)

	if err != nil || price < 0 {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	mu.Lock()
	newP := Product{ID: nextID, Name: name, Price: price}
	productsMap[nextID] = newP
	nextID++
	mu.Unlock()

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func HandleFormUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)
	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)

	if err != nil || price < 0 {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	mu.Lock()
	if _, ok := productsMap[id]; ok {
		productsMap[id] = Product{ID: id, Name: name, Price: price}
	}
	mu.Unlock()

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func HandleFormDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)

	mu.Lock()
	delete(productsMap, id)
	mu.Unlock()

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func HandleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(GetAll())
	} else if r.Method == http.MethodPost {
		var p Product
		json.NewDecoder(r.Body).Decode(&p)
		mu.Lock()
		p.ID = nextID
		nextID++
		productsMap[p.ID] = p
		mu.Unlock()
		json.NewEncoder(w).Encode(p)
	}
}
