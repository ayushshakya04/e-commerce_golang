package main

import (
	"database/sql"
	"encoding/json"

	// "encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

var err error

type SuccessResponse struct {
	Message string `json:"message"`
}

type Brands struct {
	ID         int
	Name       string
	Total_item int
}
type CartItem struct {
	Name  string
	Price float32
}

type product struct {
	Product_Name string
	Size         string
}

type userProfile struct {
	Name    string
	Email   string
	Address string
}

type userinfo struct {
	Name    string
	Email   string
	Address string
}

//	func handler(w http.ResponseWriter, r *http.Request) {
//		return
//	}
func main() {
	connstr := "host=localhost port=8081 dbname=postgres user=postgres password=Bimal@1998 sslmode=disable connect_timeout=10"
	db, err = sql.Open("postgres", connstr)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("connected !")

	http.HandleFunc("/userprofile", getuserProfile)
	http.HandleFunc("/userinfo", postuserinfo)
	http.HandleFunc("/cartitem", getcartitem)
	http.HandleFunc("/brands", getbrands)
	http.HandleFunc("/products", getproduct)
	http.HandleFunc("/deleteUser", deleteuser)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func getcartitem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT * FROM cartitem")
	if err != nil {
		http.Error(w, "get error", http.StatusInternalServerError)
		return
	}
	response := struct {
		Message string `json:"message"`
	}{
		Message: "show the data",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	defer rows.Close()

	var cartitems []CartItem
	for rows.Next() {
		var cartitem CartItem
		err := rows.Scan(&cartitem.Name, &cartitem.Price)
		if err != nil {
			http.Error(w, "Internal server an errors", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		cartitems = append(cartitems, cartitem)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Internal server errors", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Types", "applications/json")

	// Encode the cartitems as JSON and send the response
	err = json.NewEncoder(w).Encode(cartitems)
	if err != nil {
		http.Error(w, "Internal servers error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func getbrands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Methods not allowed", http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT * From brands")
	if err != nil {
		http.Error(w, "get an error", http.StatusInternalServerError)
		return
	}
	response := struct {
		Message string `json:"message"`
	}{
		Message: "show the brands data",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	defer rows.Close()

	var brand []Brands
	for rows.Next() {
		var brands Brands
		err := rows.Scan(&brands.Name, &brands.Total_item)
		if err != nil {
			http.Error(w, "Internal server errors", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		brand = append(brand, brands)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "Internals server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Contents-Type", "application/jsom")

	err = json.NewEncoder(w).Encode(brand)
	if err != nil {
		http.Error(w, "Internal server an error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

}
func getproduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
		return
	}
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		http.Error(w, "get error", http.StatusInternalServerError)
		return
	}
	response := struct {
		Message string `json:"message"`
	}{
		Message: "show the product data",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	defer rows.Close()

	var Products []product
	for rows.Next() {
		var products product
		err := rows.Scan(&products.Product_Name, &products.Size)
		if err != nil {
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		Products = append(Products, products)

	}
	if err = rows.Err(); err != nil {
		http.Error(w, "Internal server an error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Types", "application/jsom")
	err = json.NewEncoder(w).Encode(Products)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

// var getcartitem cartItem
// err := json.NewDecoder(r.Body).Decode(&getcartitem)
// if err != nil {
// 	http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 	return
// }

// _, err = db.Exec("SELECT * FROM cartitem")
// if err !=nil{
// 	http.Error(w, "Database error", http.StatusInternalServerError)
// 	return
// }
// defer rows.Close()

// cartItem := []cartItem{}
// for rows.Next() {
// 	var getcartitem cartItem
// 	err := rows.Scan(&getcartitem.ID, &getcartitem.Name, &getcartitem.Price)

// 	if err != nil {
// 		http.Error(w, "Database error", http.StatusInternalServerError)
// 		return
// 	}
// 	cartItem = append(cartItem, getcartitem)
// }

//	if err := rows.Err(); err != nil {
//		http.Error(w, "Database error", http.StatusInternalServerError)
//		return
//	}
func getuserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT * FROM userprofile")
	if err != nil {
		http.Error(w, "get an error", http.StatusInternalServerError)
		return
	}
	//printing the message in postman terminal

	// response := struct {
	// 	Message string `json:"message"`
	// }{
	// 	Message: "show the product data",
	// }
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(response)

	defer rows.Close()

	var profile []userProfile
	for rows.Next() {
		var userprofile userProfile
		err := rows.Scan(&userprofile.Name, &userprofile.Email, &userprofile.Address)
		if err != nil {
			http.Error(w, "Internal server an errors", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		profile = append(profile, userprofile)
	}
	w.Header().Set("COntent-Type", "application/json")

	err = json.NewEncoder(w).Encode(profile)
	if err != nil {
		http.Error(w, "Internal servers error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

}

func postuserinfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
		return
	}

	var user userinfo
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO userprofile (name, email, address) VALUES($1, $2, $3)", user.Name, user.Email, user.Address)
	if err != nil {
		http.Error(w, "Unable to save user info.", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Data created successfully",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func deleteuser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method is not allowed", http.StatusInternalServerError)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	// rows, err = db.Query("DELETE FROM userprofile WHERE id = $1")

	_, err = db.Exec("DELETE FROM userprofile WHERE id =$1", id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

}

// UPDATE users
// SET first_name = 'Johnny', last_name = 'Appleseed'
// WHERE id = 1;

// Parse the request body
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 		return
// 	}

// 	// Extract user profile data from the form
// 	name := r.Form.Get("name")
// 	email := r.Form.Get("email")
// 	address := r.Form.Get("address")

// 	// Check if the user profile already exists
// 	var count int
// 	err = db.QueryRow("SELECT COUNT(*) FROM userprofile WHERE name = $1 AND email = $2", name, email).Scan(&count)
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		log.Println(err)
// 		return
// 	}

// 	if count > 0 {
// 		// User profile already exists
// 		response := struct {
// 			Message string `json:"message"`
// 		}{
// 			Message: "User profile already exists",
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusConflict)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	// Insert the user profile into the database
// 	_, err = db.Exec("INSERT INTO userprofile (name, email, address) VALUES ($1, $2, $3)", name, email, address)
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		log.Println(err)
// 		return
// 	}

// 	// User profile added successfully
// 	response := struct {
// 		Message string `json:"message"`
// 	}{
// 		Message: "User profile added successfully",
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(response)
