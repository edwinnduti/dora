package helpers

//message form
type Message struct {
	Name    string
	From    string
	Subject []byte
	Message []byte
}

//smtp FQDN
type SmtpServer struct {
	Host string
	Port string
}
