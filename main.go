package main

import (
	"database/sql"
	"encoding/json"

	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

var err error

type SuccessResponse struct {
	Message string `json:"message"`
}

type Brands struct {
	ID         int    `json:"ID"`
	Name       string `json:"name" validate:"required,name"`
	Total_item int    `json:"total_item" validate:"required,total_item"`
}
type CartItem struct {
	Name  string  `json:"name" validate:"required,name"`
	Price float32 `json:"price" validate:"required,price"`
}

type product struct {
	Product_Name string `json:"product_name" validate:"required,product_name"`
	Size         string `json:"size" validate:"required,size"`
}

type userProfile struct {
	Name    string `json:"Name" validate:"required,name"`
	Email   string `json:email" validate:"required,email"`
	Address string `json:address" validate:"required,address"`
}

type userinfo struct {
	Name    string `json:"Name" validate:"required,name"`
	Email   string `json:email" validate:"required,email"`
	Address string `json:address" validate:"required,address"`
}
type User struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
}

// type PostCartitem struct {
// }

// type envconfig struct {
// 	dbname   string
// 	password string
// }

//	func handler(w http.ResponseWriter, r *http.Request) {
//		return
//	}
func main() {

	// envconfig := os.Getenv(".env")
	// fmt.Print(envconfig)

	// connstr := "host=localhost port=8081 dbname=envconfig user=postgres password=envconfig sslmode=disable connect_timeout=10"
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error .env file")
	}
	connstr := os.Getenv("DB_CONNECTION_STRING")
	db, err = sql.Open("postgres", connstr)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("connected !")

	http.HandleFunc("/postcartitem", Postcartitem) //post
	http.HandleFunc("/postbrands", Postbrands)
	http.HandleFunc("/postproduct", Postproduct)
	http.HandleFunc("/userprofile", getuserProfile) //GET
	http.HandleFunc("/userinfo", postuserinfo)      //POST
	http.HandleFunc("/cartitem", getcartitem)       //GET
	http.HandleFunc("/brands", getbrands)           //GET
	http.HandleFunc("/products", getproduct)        //GET
	http.HandleFunc("/deleteUser", deleteuser)      //DELETE

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Post cart item Function
func Postcartitem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
		return
	}
	// decoder := json.NewDecoder(r.Body)
	// decoder.DisallowUnknownFields()

	var item CartItem
	// err := json.NewDecoder(r.Body).Decode(&item)
	// if err != nil {
	// 	http.Error(w, "Invalid request payload", http.StatusBadRequest)
	// 	return
	// }

	var payload map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//Iterate through payload keys to detect unrecognized fields
	for key := range payload {
		if key != "name" && key != "price" {
			http.Error(w, "Invalid entry: unrecognized field(s) in payload", http.StatusBadRequest)
			return
		}
	}
	_, err = db.Exec("INSERT INTO cartitem (name, price) VALUES($1, $2)", item.Name, item.Price) //Exec executes a query without returning any rows.
	if err != nil {
		http.Error(w, "Unable to save the cart item", http.StatusInternalServerError)
		return
	}
	response := SuccessResponse{Message: "Cart item created Successfully"}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// post Brands function
func Postproduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
		return
	}
	var product product
	// err := json.NewDecoder(r.Body).Decode(&product)
	// if err != nil {
	// 	http.Error(w, "Invalid request payload", http.StatusBadRequest)
	// 	return
	// }
	//map is collection of key value pair where each value is unique and map to corresponding value.

	//variable  for handling the validation.
	var payloads map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&payloads)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//Iterate through payload keys to detect unrecognized fields
	for key := range payloads {
		if key != "product_name" && key != "size" {
			http.Error(w, "Invalid entry: unrecognized field(s) in payload", http.StatusBadRequest)
			return
		}
	}
	_, err = db.Exec("INSERT INTO product (product_name, size) VALUES($1, $2)", product.Product_Name, product.Size) //Exec executes a query without returning any rows.
	if err != nil {
		http.Error(w, "Unable to save the product", http.StatusInternalServerError)
		return
	}
	response := SuccessResponse{Message: "Products are created Successfully"}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

//func for the post the products

func Postbrands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
		return
	}
	var brand Brands
	var payloads map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&brand)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	for key := range payloads {
		if key != "product_name" && key != "size" {
			http.Error(w, "Invalid entry: unrecognized field(s) in payload", http.StatusBadRequest)
			return
		}
	}
	_, err = db.Exec("INSERT INTO brands (name, total_item ) VALUES($1, $2)", brand.Name, brand.Total_item) //Exec executes a query without returning any rows.
	if err != nil {
		http.Error(w, "Unable to save the brand", http.StatusInternalServerError)
		return
	}
	response := SuccessResponse{Message: "Brands are created Successfully"}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// func to get the cart item
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

// function for get the user profile data
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

// function for the post the user information
func postuserinfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
		return
	}

	var user userinfo
	var payloads map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	for key := range payloads {
		if key != "product_name" && key != "size" {
			http.Error(w, "Invalid entry: unrecognized field(s) in payload", http.StatusBadRequest)
			return
		}
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
