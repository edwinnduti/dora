package models

//message form
type Message struct {
	TO_EMAIL	[]string
	To_Name		string
	From_Email	string
	From_Name	string
	Subject		string
	Body		[]byte
}

//smtp FQDN
type SmtpServer struct {
	Host string
	Port string
}
