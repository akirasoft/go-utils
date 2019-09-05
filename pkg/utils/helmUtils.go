package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/keptn/go-utils/pkg/models"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func getHelmChartURI(chartName string) string {
	return "helm/" + chartName + ".tgz"
}

// StoreChart stores a chart in the configuration service
func StoreChart(project string, service string, stage string, chartName string, helmChart []byte, configServiceURL string) error {
	resourceHandler := NewResourceHandler(configServiceURL)

	uri := getHelmChartURI(chartName)
	resource := models.Resource{ResourceURI: &uri, ResourceContent: string(helmChart)}

	_, err := resourceHandler.CreateServiceResources(project, stage, service, []*models.Resource{&resource})
	if err != nil {
		return fmt.Errorf("Error when storing chart %s of service %s in project %s: %s",
			chartName, service, project, err.Error())
	}
	return nil
}

// GetChart reads the chart from the configuration service
func GetChart(project string, service string, stage string, chartName string, configServiceURL string) (*chart.Chart, error) {
	resourceHandler := NewResourceHandler(configServiceURL)

	resource, err := resourceHandler.GetServiceResource(project, stage, service, getHelmChartURI(chartName))
	if err != nil {
		return nil, fmt.Errorf("Error when reading chart %s from project %s: %s",
			chartName, project, err.Error())
	}

	ch, err := LoadChart([]byte(resource.ResourceContent))
	if err != nil {
		return nil, fmt.Errorf("Error when reading chart %s from project %s: %s",
			chartName, project, err.Error())
	}
	return ch, nil
}

// LoadChart converts a byte array into a Chart
func LoadChart(data []byte) (*chart.Chart, error) {
	return chartutil.LoadArchive(bytes.NewReader(data))
}

// PackageChart packages the chart and returns it
func PackageChart(ch *chart.Chart) ([]byte, error) {
	helmPackage, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, fmt.Errorf("Error when packaging chart: %s", err.Error())
	}
	defer os.RemoveAll(helmPackage)

	name, err := chartutil.Save(ch, helmPackage)
	if err != nil {
		return nil, fmt.Errorf("Error when packaging chart: %s", err.Error())
	}

	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Error when packaging chart: %s", err.Error())
	}
	return data, nil
}

// GetDeployments returns all deployments contained in the provided chart
func GetDeployments(ch *chart.Chart) []*appsv1.Deployment {

	deployments := make([]*appsv1.Deployment, 0, 0)

	for _, templateFile := range ch.Templates {
		dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(templateFile.Data))
		for {
			var dpl appsv1.Deployment
			err := dec.Decode(&dpl)
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			if IsDeployment(&dpl) {
				deployments = append(deployments, &dpl)
			}
		}
	}

	return deployments
}

// GetServices returns all services contained in the provided chart
func GetServices(ch *chart.Chart) []*corev1.Service {

	services := make([]*corev1.Service, 0, 0)

	for _, templateFile := range ch.Templates {
		dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(templateFile.Data))
		for {
			var svc corev1.Service
			err := dec.Decode(&svc)
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			if IsService(&svc) {
				services = append(services, &svc)
			}
		}
	}

	return services
}

// IsService tests whether the provided struct is a service
func IsService(svc *corev1.Service) bool {
	return strings.ToLower(svc.Kind) == "service"
}

// IsDeployment tests whether the provided struct is a deployment
func IsDeployment(dpl *appsv1.Deployment) bool {
	return strings.ToLower(dpl.Kind) == "deployment"
}
