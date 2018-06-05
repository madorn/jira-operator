package stub

import (
	"context"

	"github.com/jmckind/jira-operator/pkg/apis/jira/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewJiraHandler() sdk.Handler {
	return &JiraHandler{}
}

type JiraHandler struct {
	// Fill me
}

func (h *JiraHandler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Jira:
		err := sdk.Create(newJiraPod(o))
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create jira pod : %v", err)
			return err
		}
	}
	return nil
}

// newJiraPod will create a jira pod
func newJiraPod(cr *v1alpha1.Jira) *v1.Pod {
	labels := map[string]string{
		"app":     "jira",
		"cluster": cr.Name,
	}
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "Jira",
				}),
			},
			Labels: labels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "jira",
					Image: "cptactionhank/atlassian-jira:7.10.0",
				},
			},
		},
	}
}
