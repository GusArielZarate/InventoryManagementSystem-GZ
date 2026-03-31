package providers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Provider struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var (
	providersMap = make(map[int]Provider)
	nextID       = 1
	mu           sync.RWMutex
)

func init() {
	providersMap[nextID] = Provider{ID: nextID, Name: "Carnes S.A.", Phone: "555-0001"}
	nextID++
	providersMap[nextID] = Provider{ID: nextID, Name: "Verduras Frescas", Phone: "555-0002"}
	nextID++
}

func GetAll() []Provider {
	mu.RLock()
	defer mu.RUnlock()
	var providers []Provider
	for _, p := range providersMap {
		providers = append(providers, p)
	}
	return providers
}

func GetByID(id int) *Provider {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := providersMap[id]
	if !ok {
		return nil
	}
	return &p
}

func HandleFormCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/providers", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	phone := r.FormValue("phone")

	mu.Lock()
	p := Provider{ID: nextID, Name: name, Phone: phone}
	providersMap[nextID] = p
	nextID++
	mu.Unlock()

	http.Redirect(w, r, "/providers", http.StatusSeeOther)
}

func HandleFormUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/providers", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)
	name := r.FormValue("name")
	phone := r.FormValue("phone")

	mu.Lock()
	if _, ok := providersMap[id]; ok {
		providersMap[id] = Provider{ID: id, Name: name, Phone: phone}
	}
	mu.Unlock()

	http.Redirect(w, r, "/providers", http.StatusSeeOther)
}

func HandleFormDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/providers", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)

	mu.Lock()
	delete(providersMap, id)
	mu.Unlock()

	http.Redirect(w, r, "/providers", http.StatusSeeOther)
}

func HandleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(GetAll())
	}
}
