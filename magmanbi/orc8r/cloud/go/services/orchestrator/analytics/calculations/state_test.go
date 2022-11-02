package calculations_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/lib/go/metrics"
	"middlewareApp/magmanbi/orc8r/cloud/go/orc8r"
	"middlewareApp/magmanbi/orc8r/cloud/go/serdes"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/analytics/calculations"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/configurator"
	configurator_test_init "middlewareApp/magmanbi/orc8r/cloud/go/services/configurator/test_init"
	orchestrator_calcs "middlewareApp/magmanbi/orc8r/cloud/go/services/orchestrator/analytics/calculations"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/orchestrator/obsidian/models"
	state_test_init "middlewareApp/magmanbi/orc8r/cloud/go/services/state/test_init"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/state/test_utils"
)

func TestSiteCalculations(t *testing.T) {
	configurator_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity(context.Background(), "n0", configurator.NetworkEntity{
		Type:       orc8r.MagmadGatewayType,
		Key:        "g0",
		Config:     &models.MagmadGatewayConfigs{},
		PhysicalID: "hw0"}, serdes.Entity)
	assert.NoError(t, err)

	ctx := test_utils.GetContextWithCertificate(t, "hw0")
	test_utils.ReportGatewayStatus(t, ctx, models.NewDefaultGatewayStatus("hw0"))
	analyticsConfig := &calculations.AnalyticsConfig{
		Metrics: map[string]calculations.MetricConfig{
			metrics.GatewayMagmaVersionMetric: {
				Export:   true,
				Register: true,
			},
		},
	}
	siteMetricsCalculation := orchestrator_calcs.SiteMetricsCalculation{
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: calculations.CalculationParams{
				AnalyticsConfig: analyticsConfig,
			},
		},
	}
	results, err := siteMetricsCalculation.Calculate(nil)
	assert.NoError(t, err)
	t.Log(results)
	resultMetricMap := make(map[string]string)
	for _, result := range results {
		resultMetricMap[result.GetLabels()[metrics.GatewayLabelName]] = result.GetLabels()[metrics.GatewayMagmaVersionLabel]
	}
	assert.Equal(t, resultMetricMap["hw0"], "0.0.0.0")
}

func TestNetworkCalculations(t *testing.T) {
	configurator_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n0_1", Type: "LTE"}, serdes.Network)
	configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1", Type: "FEG_LTE"}, serdes.Network)
	configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n2_0", Type: "FEG"}, serdes.Network)
	configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n2_2", Type: "FEG"}, serdes.Network)
	analyticsConfig := &calculations.AnalyticsConfig{
		Metrics: map[string]calculations.MetricConfig{
			metrics.NetworkTypeMetric: {
				Export:   true,
				Register: true,
			},
		},
	}
	generalCalculation := orchestrator_calcs.NetworkMetricsCalculation{
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: calculations.CalculationParams{
				AnalyticsConfig: analyticsConfig,
			},
		},
	}
	results, err := generalCalculation.Calculate(nil)
	resultMetricMap := make(map[string]string)
	for _, result := range results {
		resultMetricMap[result.GetLabels()[metrics.NetworkLabelName]] = result.GetLabels()[metrics.NetworkTypeLabel]
	}
	assert.NoError(t, err)
	t.Log(results)
	assert.Equal(t, resultMetricMap["n0_1"], "LTE")
	assert.Equal(t, resultMetricMap["n1"], "FEG_LTE")
	assert.Equal(t, resultMetricMap["n2_0"], "FEG")
	assert.Equal(t, resultMetricMap["n2_2"], "FEG")
}
