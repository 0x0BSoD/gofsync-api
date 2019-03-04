package entitys

type Auth struct {
	Username string `json:"username"`
	Pass     string `json:"password"`
	Port     int    `json:"port"`
	DBFile   string `json:"db_file"`
}
