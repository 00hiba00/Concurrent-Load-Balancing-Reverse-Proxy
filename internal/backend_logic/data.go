package backendlogic

type Book struct{
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Price  float64 `json:"price"`
}

var Books []Book

//Go has a special function called init(). It runs automatically
//once when the package is first loaded, before the main() function even starts.
//This is useful if you want to generate IDs or do more complex setup.
func init() {
    Books = append(Books, Book{ID: "1", Title: "The Go Programming Language", Author: "Alan Donovan", Year: 2015, Price: 34.99})
}