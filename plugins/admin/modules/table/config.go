package table

// Config 為引擎、是否可以增加、編輯....等資訊
type Config struct {
	Driver         string
	Connection     string
	CanAdd         bool
	Editable       bool
	Deletable      bool
	Exportable     bool
	PrimaryKey     PrimaryKey
	SourceURL      string
	GetDataFun     GetDataFun
	OnlyInfo       bool
	OnlyNewForm    bool
	OnlyUpdateForm bool
	OnlyDetail     bool
}

// 預設Config(struct)，引擎mysql主鍵id
func DefaultConfig() Config {
	return Config{
		Driver:     "mysql",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		Connection: "default",
		PrimaryKey: PrimaryKey{
			Type: "INT",
			Name: "id",
		},
	}
}

// 預設Config(struct)，driver設為參數，主鍵id
func DefaultConfigWithDriver(driver string) Config {
	return Config{
		Driver:     driver,
		Connection: "default",
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		PrimaryKey: PrimaryKey{
			Type: "INT",
			Name: "id",
		},
	}
}
