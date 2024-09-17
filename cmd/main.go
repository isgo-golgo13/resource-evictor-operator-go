package main

import (
	"os"

	"github.com/isgo-golgo13/resource-evictor-operator/internal/controller"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()        // Initialize scheme
	setupLog = ctrl.Log.WithName("setup") // Initialize logger
)

func init() {
	// Add client-go (core Kubernetes types) to the scheme
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	// If you have other CRDs (e.g., custom resources), you should add them to the scheme here
	// utilruntime.Must(myCustomResourceV1.AddToScheme(scheme))
}

func main() {
	var enableLeaderElection bool
	var metricsAddr string

	// Initialize controller-runtime logger (using Zap logger)
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Manager setup with controller-runtime
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,      // Pass the initialized scheme
		MetricsBindAddress: metricsAddr, // Address for metrics
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "resource-evictor-operator-lock",
	})

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controller.WorkloadScanReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "WorkloadScan")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
