package controllers

import "github.com/astaxie/beego/httplib"

type ImageController struct {
	BaseController
}

type Images struct {
	ImageName string   `json:"imageName,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

func (c ImageController) GetImages() {
	var imagesName struct {
		Repositories []string `json:"repositories,omitempty"`
	}

	httplib.Get("http://114.116.173.97:5000/v2/_catalog").ToJSON(&imagesName)
	imagesNum := len(imagesName.Repositories)

	tagsName := make([]struct {
		Name string   `json:"name,omitempty"`
		Tags []string `json:"tags,omitempty"`
	}, imagesNum)
	images := make([]Images, imagesNum)

	for k, v := range imagesName.Repositories {
		httplib.Get("http://114.116.173.97:5000/v2/" + v + "/tags/list").ToJSON(&tagsName[k])
		images[k].ImageName = v
		images[k].Tags = tagsName[k].Tags
	}
	c.ApiResponse(200, "获取成功", map[string]interface{}{
		"images": images,
	})
}
