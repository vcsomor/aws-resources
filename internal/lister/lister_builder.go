package lister

import (
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
)

type Builder struct {
	clientFactory conn.ClientFactory
	executor      executor.SynchronousExecutor
	logger        *logrus.Logger

	regions   []string
	resources []string
}

func NewLister(logger *logrus.Logger, clientFactory conn.ClientFactory, executor executor.SynchronousExecutor) *Builder {
	return &Builder{
		clientFactory: clientFactory,
		executor:      executor,
		logger:        logger,
	}
}

func (b *Builder) WithRegions(regions []string) *Builder {
	b.regions = regions
	return b
}

func (b *Builder) WithResources(resources []string) *Builder {
	b.resources = resources
	return b
}

func (b *Builder) Build() Lister {
	return &taskBasedLister{
		clientFactory: b.clientFactory,
		executor:      b.executor,
		logger:        b.logger,

		regions:   b.regions,
		resources: b.resources,
	}
}
