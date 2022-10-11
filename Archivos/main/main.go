package main

func main() {
	//IngresarEmails()
	mux := Routes()
	server := NewServer(mux)
	server.Run()
}
