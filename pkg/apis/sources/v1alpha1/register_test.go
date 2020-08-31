package v1alpha1

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

func TestRegisterHelpers(t *testing.T) {
	if got, want := Kind("Foo"), "Foo.sources.muru.dev"; got.String() != want {
		t.Errorf("Kind(Foo) = %v, want %v", got.String(), want)
	}

	if got, want := Resource("Foo"), "Foo.sources.muru.dev"; got.String() != want {
		t.Errorf("Resource(Foo) = %v, want %v", got.String(), want)
	}

	if got, want := SchemeGroupVersion.String(), "sources.muru.dev/v1alpha1"; got != want {
		t.Errorf("SchemeGroupVersion() = %v, want %v", got, want)
	}

	scheme := runtime.NewScheme()
	if err := addKnownTypes(scheme); err != nil {
		t.Errorf("addKnownTypes() = %v", err)
	}
}
