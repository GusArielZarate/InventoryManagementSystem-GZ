package inventory

import (
	"encoding/json"
	"inventory-fridays/products"
	"net/http"
	"strconv"
	"sync"
)

type StockView struct {
	ProductID   int
	ProductName string
	Quantity    int
}

type StockItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

var (
	stockMap = make(map[int]int)
	mu       sync.RWMutex
)

func init() {
	stockMap[1] = 50
	stockMap[2] = 100
}

func GetFullInventory() []StockView {
	mu.RLock()
	defer mu.RUnlock()

	allProducts := products.GetAll()
	var views []StockView

	for _, p := range allProducts {
		qty := stockMap[p.ID]
		views = append(views, StockView{
			ProductID:   p.ID,
			ProductName: p.Name,
			Quantity:    qty,
		})
	}
	return views
}

func HandleFormUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/inventory", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("product_id")
	qtyStr := r.FormValue("quantity")

	id, _ := strconv.Atoi(idStr)
	qty, err := strconv.Atoi(qtyStr)

	if err != nil || qty < 0 {
		http.Redirect(w, r, "/inventory", http.StatusSeeOther)
		return
	}

	mu.Lock()
	stockMap[id] = qty
	mu.Unlock()

	http.Redirect(w, r, "/inventory", http.StatusSeeOther)
}

func HandleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		mu.RLock()
		json.NewEncoder(w).Encode(stockMap)
		mu.RUnlock()
	} else if r.Method == http.MethodPost {
		var item StockItem
		json.NewDecoder(r.Body).Decode(&item)
		mu.Lock()
		stockMap[item.ProductID] = item.Quantity
		mu.Unlock()
		json.NewEncoder(w).Encode(item)
	}
}
