package entitys

type Auth struct {
	Username string
	Pass     string
	Port     int
	DBFile   string
	Actions  []string
	RTPro    string
	RTStage  string
}
