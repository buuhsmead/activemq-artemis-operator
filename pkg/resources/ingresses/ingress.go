package ingresses

import (
	"context"
	"github.com/rh-messaging/activemq-artemis-operator/pkg/utils/selectors"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/rh-messaging/activemq-artemis-operator/pkg/apis/broker/v2alpha1"
	extv1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("package ingresses")

// Create newIngressForCR method to create exposed ingress
func NewIngressForCR(cr *v2alpha1.ActiveMQArtemis, target string) *extv1b1.Ingress {

	ingress := &extv1b1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    selectors.LabelBuilder.Labels(),
			Name:      cr.Name + "-" + target,
			Namespace: cr.Namespace,
		},
		Spec: extv1b1.IngressSpec{
			Rules: []extv1b1.IngressRule{
				{
					Host: os.Getenv("KUBERNETES_SERVICE_HOST"),
					IngressRuleValue: extv1b1.IngressRuleValue{
						HTTP: &extv1b1.HTTPIngressRuleValue{
							Paths: []extv1b1.HTTPIngressPath{
								extv1b1.HTTPIngressPath{
									Path: "/",
									Backend: extv1b1.IngressBackend{
										ServiceName: "hs",
										ServicePort: intstr.FromString(target),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return ingress
}

func Create(cr *v2alpha1.ActiveMQArtemis, client client.Client, scheme *runtime.Scheme, ingress *extv1b1.Ingress) error {

	reqLogger := log.WithValues("ActiveMQArtemis Name", cr.Name)
	reqLogger.Info("Creating new ingress")

	var err error = nil
	// Set ActiveMQArtemis instance as the owner and controller
	reqLogger.Info("Set controller reference for new  ingress")
	if err = controllerutil.SetControllerReference(cr, ingress, scheme); err != nil {
		reqLogger.Error(err, "Failed to set controller reference for new ingress")
	}

	// Call k8s create for ingress
	if err = client.Create(context.TODO(), ingress); err != nil {
		reqLogger.Error(err, "Failed to creating new ingress")
	}
	reqLogger.Info("End of ingress Creation")

	return err
}

func Retrieve(cr *v2alpha1.ActiveMQArtemis, namespacedName types.NamespacedName, client client.Client, ingress *extv1b1.Ingress) error {

	// Log where we are and what we're doing
	reqLogger := log.WithValues("ActiveMQArtemis Name", cr.Name)
	reqLogger.Info("Retrieving the ingress ")

	var err error = nil

	// Check if the headless ingress already exists
	if err = client.Get(context.TODO(), namespacedName, ingress); err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Error(err, "ingress Not Found", "Namespace", cr.Namespace, "Name", cr.Name)
		} else {
			reqLogger.Error(err, "ingress found", "Namespace", cr.Namespace, "Name", cr.Name)
		}
	}

	return err
}
