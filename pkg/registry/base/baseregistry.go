package baseregistry

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
	"github.com/apache/dubbo-go/common"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

type FacadeRegistry interface {
	// LoadInterfaces loads the dubbo services from interfaces level
	LoadInterfaces() ([]router.API, []error)
	// LoadApplication loads the dubbo services from application level
	LoadApplications() ([]router.API, []error)
	// DoSubscribe subscribes the registries to monitor the changes.
	DoSubscribe(*common.URL) error
	// DoUnsubscribe unsubscribes the registries.
	DoUnsubscribe(*common.URL) error
}

type BaseRegistry struct {
	listeners      []registry.Listener
	facadeRegistry FacadeRegistry
}

func NewBaseRegistry(facade FacadeRegistry) *BaseRegistry {
	return &BaseRegistry{
		listeners:      []registry.Listener{},
		facadeRegistry: facade,
	}
}

// LoadServices loads all the registered Dubbo services from registry
func (r *BaseRegistry) LoadServices() error {
	interfaceAPIs, _ := r.facadeRegistry.LoadInterfaces()
	applicationAPIs, _ := r.facadeRegistry.LoadApplications()
	apis := r.deduplication(append(interfaceAPIs, applicationAPIs...))
	localAPIDiscSrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	for i := range apis {
		localAPIDiscSrv.AddAPI(apis[i])
		// r.facadeRegistry.DoSubscribe()
	}
	return nil
}

func (r *BaseRegistry) deduplication(apis []router.API) []router.API {
	urlPatternMap := make(map[string]router.API)
	for i := range apis {
		urlPatternMap[apis[i].URLPattern] = apis[i]
	}
	rstAPIs := []router.API{}
	for i := range urlPatternMap {
		rstAPIs = append(rstAPIs, urlPatternMap[i])
	}
	return rstAPIs
}

// Subscribe monitors the target registry.
func (r *BaseRegistry) Subscribe(url *common.URL) error {
	return r.facadeRegistry.DoSubscribe(url)
}

// Unsubscribe stops monitoring the target registry.
func (r *BaseRegistry) Unsubscribe(url *common.URL) error {
	return r.facadeRegistry.DoUnsubscribe(url )
}