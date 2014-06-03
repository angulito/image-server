package mapper

import (
	"fmt"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/encoders/base62"
)

type SourceMapper struct {
	MapperConfiguration *core.MapperConfiguration
}

// RemoteImageURL returns a URL string for original image
func (m *SourceMapper) RemoteImageURL(ic *core.ImageConfiguration) string {
	if ic.Source != "" {
		return ic.Source
	}
	url := ic.ServerConfiguration.SourceDomain + "/" + m.imageDirectory(ic) + "/original.jpg"
	return url
}

func (m *SourceMapper) imageDirectory(ic *core.ImageConfiguration) string {
	id := base62.Decode(ic.ID)
	return fmt.Sprintf("%s/%d", m.MapperConfiguration.NamespaceMappings[ic.Namespace], id)
}