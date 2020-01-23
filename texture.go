package main

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func loadTexture(path string) *sdl.Texture {
	image, err := img.Load(path)
	if err != nil {
		panic(err)
	}
	defer image.Free()
	image.SetColorKey(true, 0xffff00ff)
	texture, err := gc.renderer.CreateTextureFromSurface(image)
	if err != nil {
		panic(err)
	}
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	return texture
}

func texturePath(name string) string {
	return "assets/" + name + ".png"
}
