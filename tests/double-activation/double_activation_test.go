package e2e

import (
	"context"
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"sigs.k8s.io/e2e-framework/third_party/helm"
)

var (
	testEnv         env.Environment
	namespace       string
	kindClusterName string
)

func TestMain(m *testing.M) {
	cfg, _ := envconf.NewFromFlags()
	testEnv = env.NewWithConfig(cfg)
	kindClusterName = envconf.RandomName("double-activation", 16)
	namespace = envconf.RandomName("dapr-tests", 16)

	testEnv.Setup(
		envfuncs.CreateKindCluster(kindClusterName),
		envfuncs.CreateNamespace(namespace),
	)

	testEnv.Finish(
		envfuncs.DeleteNamespace(namespace),
		envfuncs.DestroyKindCluster(kindClusterName),
	)
	os.Exit(testEnv.Run(m))
}

func TestDoubleActivation(t *testing.T) {
	feature := features.New("Double activation test").Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
		manager := helm.New(config.KubeconfigFile())
		err := manager.RunRepo(helm.WithArgs("add", "dapr", "https://dapr.github.io/helm-charts/"))
		if err != nil {
			t.Fatal("failed to add dapr helm chart repo")
		}
		err = manager.RunRepo(helm.WithArgs("update"))
		if err != nil {
			t.Fatal("failed to upgrade helm repo")
		}
		err = manager.RunInstall(helm.WithName("dapr"), helm.WithWait(), helm.WithNamespace(namespace), helm.WithReleaseName("dapr/dapr"))
		if err != nil {
			t.Fatal("failed to install dapr Helm chart")
		}
		return ctx
	}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			manager := helm.New(config.KubeconfigFile())
			err := manager.RunRepo(helm.WithArgs("remove", "dapr"))
			if err != nil {
				t.Fatal("cleanup of the helm repo failed")
			}
			return ctx
		}).Feature()

	testEnv.Test(t, feature)
}
