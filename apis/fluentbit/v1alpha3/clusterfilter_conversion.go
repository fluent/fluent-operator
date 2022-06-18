package v1alpha3

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this ClusterFilter to the Hub version (v1alpha2).
func (src *ClusterFilter) ConvertTo(dstRaw conversion.Hub) error {
	return nil
}

// ConvertFrom converts from the Hub version (v1alpha2) to this version.
func (dst *ClusterFilter) ConvertFrom(srcRaw conversion.Hub) error {
	return nil
}
