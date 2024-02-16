package lister

import (
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
)

type Dependencies struct {
	clientFactory conn.ClientFactory
	executor      executor.SynchronousExecutor
	logger        *logrus.Logger
	writerFactory ResultBasedWriterFactory
}

type DependencyFn func(d *Dependencies)

func WithClientFactory(cf conn.ClientFactory) DependencyFn {
	return func(d *Dependencies) {
		d.clientFactory = cf
	}
}

func WithExecutor(e executor.SynchronousExecutor) DependencyFn {
	return func(d *Dependencies) {
		d.executor = e
	}
}

func WithLogger(l *logrus.Logger) DependencyFn {
	return func(d *Dependencies) {
		d.logger = l
	}
}

func WithWriterFactory(f ResultBasedWriterFactory) DependencyFn {
	return func(d *Dependencies) {
		d.writerFactory = f
	}
}

type Parameters struct {
	regions   []string
	resources []string
}

type ParametersFn func(p *Parameters)

func WithRegions(regions []string) ParametersFn {
	return func(p *Parameters) {
		p.regions = append([]string{}, regions...)
	}
}

func WithResources(resources []string) ParametersFn {
	return func(p *Parameters) {
		p.resources = append([]string{}, resources...)
	}
}

type Builder struct {
	depFns   []DependencyFn
	paramFns []ParametersFn
}

func NewLister() *Builder {
	return &Builder{}
}

func (b *Builder) Dependencies(depFns ...DependencyFn) *Builder {
	b.depFns = append(b.depFns, depFns...)
	return b
}

func (b *Builder) Parameters(paramFns ...ParametersFn) *Builder {
	b.paramFns = append(b.paramFns, paramFns...)
	return b
}

func (b *Builder) Build() Lister {
	var deps Dependencies
	for _, fnc := range b.depFns {
		fnc(&deps)
	}

	var params Parameters
	for _, fnc := range b.paramFns {
		fnc(&params)
	}

	return &taskBasedLister{
		clientFactory: deps.clientFactory,
		executor:      deps.executor,
		logger:        deps.logger,
		writerFactory: deps.writerFactory,

		regions:   params.regions,
		resources: params.resources,
	}
}
