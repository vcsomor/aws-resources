package lister

import (
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
)

type listerBuilder struct {
	clientFactory conn.ClientFactory
	executor      executor.SynchronousExecutor
	logger        *logrus.Logger

	regions   []string
	resources []string
}

func newLister(logger *logrus.Logger, clientFactory conn.ClientFactory, executor executor.SynchronousExecutor) *listerBuilder {
	return &listerBuilder{
		clientFactory: clientFactory,
		executor:      executor,
		logger:        logger,
	}
}

func (b *listerBuilder) withRegions(regions []string) *listerBuilder {
	b.regions = regions
	return b
}

func (b *listerBuilder) withResources(resources []string) *listerBuilder {
	b.resources = resources
	return b
}

func (b *listerBuilder) build() Lister {
	return &defaultLister{
		clientFactory: b.clientFactory,
		executor:      b.executor,
		logger:        b.logger,

		regions:   b.regions,
		resources: b.resources,
	}
}
