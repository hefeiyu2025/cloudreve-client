package cloudreve

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 汇总方法，一键操作

func md5Hash(key string) string {
	// 计算JSON数据的MD5
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func isEmpty(dirPath string) (bool, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer dir.Close()
	//如果目录不为空，Readdirnames 会返回至少一个文件名
	_, err = dir.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// UploadPath 一键上传路径
func (c *CloudreveClient) UploadPath(req OneStepUploadPathReq) error {
	// 遍历目录
	err := filepath.Walk(req.LocalPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			for _, ignorePath := range req.IgnorePaths {
				if filepath.Base(path) == ignorePath {
					return filepath.SkipDir
				}
			}
		} else {
			// 获取相对于root的相对路径
			relPath, _ := filepath.Rel(req.LocalPath, path)
			relPath = strings.Replace(relPath, "\\", "/", -1)
			relPath = strings.Replace(relPath, info.Name(), "", 1)
			NotUpload := false
			for _, ignoreFile := range req.IgnoreFiles {
				if info.Name() == ignoreFile {
					NotUpload = true
					break
				}
			}
			for _, extension := range req.IgnoreExtensions {
				if strings.HasSuffix(info.Name(), extension) {
					NotUpload = true
					break
				}
			}
			for _, extension := range req.Extensions {
				if strings.HasSuffix(info.Name(), extension) {
					NotUpload = false
					break
				}
				NotUpload = true
			}
			if !NotUpload {
				err = c.UploadFile(OneStepUploadFileReq{
					LocalFile:  path,
					RemotePath: strings.TrimRight(req.RemotePath, "/") + "/" + relPath,
					PolicyId:   req.PolicyId,
					Resumable:  req.Resumable,
					SuccessDel: req.SuccessDel,
				})
				if err == nil {
					if req.SuccessDel {
						dir := filepath.Dir(path)
						if dir != "." {
							empty, _ := isEmpty(dir)
							if empty {
								_ = os.Remove(dir)
								fmt.Println("uploaded success and delete", dir)
							}
						}
					}
				} else {
					if !req.SkipFileErr {
						return err
					} else {
						fmt.Println("upload err", err)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// UploadFile 一键上传文件
func (c *CloudreveClient) UploadFile(req OneStepUploadFileReq) error {
	file, err := os.Open(req.LocalFile)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	remotePath := strings.TrimLeft(req.RemotePath, "/")
	md5Key := md5Hash(req.LocalFile + remotePath + req.PolicyId)
	var session UploadCredential
	if req.Resumable {
		cacheErr := GetCache("session_"+md5Key, &session)
		if cacheErr != nil {
			fmt.Println("cache err:", cacheErr)
		}
	}
	if session.SessionID == "" {
		resp, err := c.FileUploadGetUploadSession(CreateUploadSessionReq{
			Path:         "/" + remotePath,
			Size:         uint64(stat.Size()),
			Name:         stat.Name(),
			PolicyID:     req.PolicyId,
			LastModified: stat.ModTime().UnixMilli(),
		})
		if err != nil {
			return err
		}
		session = resp.Data
		cacheErr := SetCache("session_"+md5Key, session)
		if cacheErr != nil {
			fmt.Println("cache err:", cacheErr)
		}
	}

	uploadedSize := 0
	if req.Resumable {
		cacheErr := GetCache("chunk_"+md5Key, &uploadedSize)
		if cacheErr != nil {
			fmt.Println("cache err:", cacheErr)
		}
	}
	uploaded, err := c.OneDriveUpload(OneDriveUploadReq{
		UploadUrl:    session.UploadURLs[0],
		LocalFile:    file,
		UploadedSize: int64(uploadedSize),
		ChunkSize:    int64(session.ChunkSize),
	})
	if err != nil {
		dealError(req.Resumable, md5Key, session.SessionID, uploaded, c)
		return err
	}

	_, err = c.OneDriveCallback(session.SessionID)
	if err != nil {
		dealError(req.Resumable, md5Key, session.SessionID, uploaded, c)
		return err
	}
	if req.Resumable {
		_ = DelCache("session_" + md5Key)
		_ = DelCache("chunk_" + md5Key)
	}
	// 上传成功则移除文件了
	if req.SuccessDel {
		_ = os.Remove(req.LocalFile)
		fmt.Println("uploaded success and delete", req.LocalFile)
	}
	return nil
}

func dealError(resumable bool, md5Key, sessionId string, uploaded int64, c *CloudreveClient) {
	needDelSession := true
	if resumable {
		cacheErr := SetCache("chunk_"+md5Key, uploaded)
		if cacheErr != nil {
			fmt.Println("cache err:", cacheErr)
		} else {
			needDelSession = false
		}
	}
	if needDelSession {
		_, delErr := c.FileUploadDeleteUploadSession(sessionId)
		if delErr != nil {
			fmt.Println(delErr)
		} else {
			_ = DelCache("session_" + md5Key)
			_ = DelCache("chunk_" + md5Key)
		}
	}
}

type RenameDealFunc func(obj Object) string

// RenameDeal 一键重命名
func (c *CloudreveClient) RenameDeal(path string, fn RenameDealFunc) error {
	directory, err := c.ListDirectory("/" + strings.TrimLeft(path, "/"))
	if err != nil {
		return err
	}
	data := directory.Data
	for _, object := range data.Objects {
		if object.Type == "dir" {
			err := c.RenameDeal(object.Path+"/"+object.Name, fn)
			if err != nil {
				return err
			}
		}
		newName := fn(object)
		if newName != object.Name {
			item := Item{}
			if object.Type == "dir" {
				item.Dirs = []string{object.ID}
			} else {
				item.Items = []string{object.ID}
			}
			_, err := c.ObjectRename(ItemRenameReq{
				Src:     item,
				NewName: newName,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
