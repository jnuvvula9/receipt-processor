package main

func main() {
	store := NewReceiptStore()
	server := NewAPIServer(":63342", store)
	server.Run()
}
