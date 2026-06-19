package game

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
)

//go:embed assets/appicon/AppIcon-16.png assets/appicon/AppIcon-32.png assets/appicon/AppIcon-64.png assets/appicon/AppIcon-128.png assets/appicon/AppIcon-256.png assets/appicon/AppIcon-512.png assets/appicon/AppIcon-1024.png
var iconFiles embed.FS

var appIconPaths = []string{
	"assets/appicon/AppIcon-16.png",
	"assets/appicon/AppIcon-32.png",
	"assets/appicon/AppIcon-64.png",
	"assets/appicon/AppIcon-128.png",
	"assets/appicon/AppIcon-256.png",
	"assets/appicon/AppIcon-512.png",
	"assets/appicon/AppIcon-1024.png",
}

func WindowIcons() []image.Image {
	icons, err := loadWindowIcons()
	if err != nil {
		panic(err)
	}
	return icons
}

func loadWindowIcons() ([]image.Image, error) {
	icons := make([]image.Image, 0, len(appIconPaths))
	for _, path := range appIconPaths {
		data, err := iconFiles.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read app icon %s: %w", path, err)
		}
		icon, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("decode app icon %s: %w", path, err)
		}
		icons = append(icons, icon)
	}
	return icons, nil
}
