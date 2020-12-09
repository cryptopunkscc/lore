package server

type Config struct {
	Admin     bool
	Content   bool
	Transport string
	Address   string
}

var UnixAdminConfig = Config{
	Admin:     true,
	Content:   true,
	Transport: "unix",
	Address:   "/tmp/lore.sock",
}

var TCPContentConfig = Config{
	Admin:     false,
	Content:   true,
	Transport: "tcp",
	Address:   ":10768",
}
