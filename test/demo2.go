package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func main() {
	//connectdocker()
	//containersstatus()
	//ListImages()
	//StopContainer("9c14134da48d")
	query()
}

func connectdocker() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://192.168.10.168:2375"), client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		panic(err)
	}
	//imageName := "bfirsh/reticulate-splines"
	//out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//io.Copy(os.Stdout, out)

	imageName := "httpd:latest"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

	//rsp, err := cli.ContainerCreate(ctx, &container.Config{
	//	Image: imageName,
	//}, nil, nil, nil, "")
	//if err != nil {
	//	panic(err)
	//}
	//if cli.ContainerStart(ctx, rsp.ID, types.ContainerStartOptions{}); err != nil {
	//	panic(err)
	//}
	//fmt.Println(rsp.ID)
}

func containersstatus() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://192.168.10.168:2375"), client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		fmt.Println(container.Names, container.Image, container.Ports, container.Status, container.ImageID)
	}
}

func ListImages() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://192.168.10.168:2375"), client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		panic(err)
	}
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	var imagesLib []string
	for _, image := range images {
		imagesLib = append(imagesLib, image.RepoTags[0])
	}
	fmt.Println(imagesLib)
}

func StopContainer(containerId string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://192.168.10.168:2375"), client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		panic(err)
	}
	err = cli.ContainerStop(ctx, containerId, nil)
	if err != nil {
		panic(err)
	}
}

func query() {
	//Default返回一个默认的路由引擎
	r := gin.Default()
	r.GET("/user/search", func(c *gin.Context) {
		//username := c.DefaultQuery("username", "tom")
		username := c.Query("username")
		//address := c.Query("address")
		//输出json结果给调用方
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			//"address":  address,
		})
	})
	r.Run()

}
