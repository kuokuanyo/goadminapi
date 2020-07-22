package service

// db.Connection(interface)也具有Service(interface)型態，因為都具有Name方法
// Service(interface)也具有db.Connection(interface)型態，因為都具有Name方法
type Service interface {
	Name() string
}

type List map[string]Service

type Generator func() (Service, error)
type Generators map[string]Generator

var services = make(Generators)


// 初始化List(map[string]Service)，Service是interface(Name方法)
func GetServices() List {
	var (
		l   = make(List)
		err error
	)
	for k, gen := range services {
		if l[k], err = gen(); err != nil {
			panic("service initialize fail")
		}
	}
	return l
}

// 將參數k、gen將入services(map[string]Generator)中
func Register(k string, gen Generator) {
	if _, ok := services[k]; ok {
		panic("service has been registered")
	}
	services[k] = gen
}

// 透過參數(k)取得匹配的Service(interface)
func (g List) Get(k string) Service {
	if v, ok := g[k]; ok {
		return v
	}
	panic("service not found")
}

// 判斷是否有取得符合參數(k)的Service
func (g List) GetOrNot(k string) (Service, bool) {
	v, ok := g[k]
	return v, ok
}

// 藉由參數新增List(map[string]Service)
func (g List) Add(k string, service Service) {
	if _, ok := g[k]; ok {
		panic("service exist")
	}
	g[k] = service
}

