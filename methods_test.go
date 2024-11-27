package cloudreve

import (
	"fmt"
	"testing"
)

var cloudreveSession = "11MTczMDY4MjkxN3xOd3dBTkVGTVZsVXpWelZVV1RkRFJWUkVRa3BJVVZNM1NsTkJTRVUxU1ZrM1QwSkpWbFJSVjBsVVNqSkVXakpMUzFsTVZ6TkhNMUU9fGn6FJQILddlrz5oIVWWphbWmtexz6f3zCc_zGiTjmq2"
var cloudreveUrl = "http://localhost:8080"

func beforeClient() *CloudreveClient {
	return NewClient(cloudreveUrl, cloudreveSession)
}

func TestConfig(t *testing.T) {
	client := beforeClient()
	resp, err := client.Config()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileUploadGetUploadSession(t *testing.T) {
	client := beforeClient()

	resp, err := client.FileUploadGetUploadSession(CreateUploadSessionReq{
		// TODO
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileUploadDeleteUploadSession(t *testing.T) {
	client := beforeClient()

	resp, err := client.FileUploadDeleteUploadSession("")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileUploadDeleteAllUploadSession(t *testing.T) {
	client := beforeClient()

	resp, err := client.FileUploadDeleteAllUploadSession()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileCreateFile(t *testing.T) {
	client := beforeClient()

	resp, err := client.FileCreateFile("/demo/33.txt")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileCreateDownloadSession(t *testing.T) {
	client := beforeClient()

	resp, err := client.FileCreateDownloadSession("OX9B2Vuz")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

//func TestFilePreview(t *testing.T) {
//	client := beforeClient()
//
//	resp, err := client.FilePreview("mqoRMnTX")
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//	}
//	fmt.Println(resp)
//}

func TestFileGetSource(t *testing.T) {
	client := beforeClient()
	dirs := make([]string, 0)
	dirs = append(dirs, "")
	items := make([]string, 0)
	items = append(items, "mqoRMnTX")

	resp, err := client.FileGetSource(ItemReq{
		Item: Item{
			Dirs:  dirs,
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestFileArchive(t *testing.T) {
	client := beforeClient()
	dirs := make([]string, 0)
	dirs = append(dirs, "DVBmxvCo")
	items := make([]string, 0)
	items = append(items, "")

	resp, err := client.FileArchive(ItemReq{
		Item: Item{
			Dirs:  dirs,
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestCreateDirectory(t *testing.T) {
	client := beforeClient()
	resp, err := client.CreateDirectory("/demo")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestListDirectory(t *testing.T) {
	client := beforeClient()
	resp, err := client.ListDirectory("/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestObjectDelete(t *testing.T) {
	client := beforeClient()

	dirs := make([]string, 0)
	dirs = append(dirs, "")
	items := make([]string, 0)
	items = append(items, "6KZbgYC2")

	resp, err := client.ObjectDelete(ItemReq{
		Item: Item{
			Dirs:  dirs,
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestObjectMove(t *testing.T) {
	client := beforeClient()

	dirs := make([]string, 0)
	dirs = append(dirs, "")
	items := make([]string, 0)
	items = append(items, "6KZbgYC2")

	resp, err := client.ObjectMove(ItemMoveReq{
		SrcDir: "/",
		Dst:    "/11",
		Src: Item{
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestObjectCopy(t *testing.T) {
	client := beforeClient()

	dirs := make([]string, 0)
	dirs = append(dirs, "")
	items := make([]string, 0)
	items = append(items, "D1lD4aHB")

	resp, err := client.ObjectCopy(ItemMoveReq{
		SrcDir: "/11",
		Dst:    "/22",
		Src: Item{
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestObjectRename(t *testing.T) {
	client := beforeClient()

	dirs := make([]string, 0)
	dirs = append(dirs, "")
	items := make([]string, 0)
	items = append(items, "6KZLBVT2")

	resp, err := client.ObjectRename(ItemRenameReq{
		NewName: "22.txt",
		Src: Item{
			Items: items,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestObjectGetProperty(t *testing.T) {
	client := beforeClient()

	resp, err := client.ObjectGetProperty(ItemPropertyReq{
		Id:        "E7qQZnHK",
		IsFolder:  true,
		TraceRoot: true,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(resp)
}

func TestOneStepRename(t *testing.T) {
	client := beforeClient()
	err := client.Rename("/demo", func(obj Object) string {
		convert := "11" + obj.Name
		fmt.Printf("%s --> %s\n", obj.Name, convert)
		return convert
	})
	fmt.Println(err)
}

func TestOneStepUploadPath(t *testing.T) {
	client := beforeClient()
	directory, err := client.ListDirectory("/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	data := directory.Data
	err = client.UploadPath(OneStepUploadPathReq{
		LocalPath:   "./demo",
		RemotePath:  "demo",
		PolicyId:    data.Policy.ID,
		Resumable:   false,
		SkipFileErr: true,
		SuccessDel:  false,
	})
	fmt.Println(err)
}

func TestDownload(t *testing.T) {
	client := beforeClient()

	err := client.Download(OneStepDownloadReq{
		Remote: "/aa", LocalPath: "./aa", IsParallel: false,
		DownloadCallback: func(localPath, localFile string) {
			fmt.Println(localFile)
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
