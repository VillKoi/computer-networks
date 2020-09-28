package client

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/jlaffaye/ftp"
)

type FTPClient struct {
	ftp *ftp.ServerConn
}

// put local/example.txt server/new_file.txt
func (c FTPClient) StoreFile(path, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(data)
	err = c.ftp.Stor(path, buffer)
	return err
}

// get example.txt new_file.txt
func (c FTPClient) ReadFile(path, newPath string) error {
	r, err := c.ftp.Retr(path)
	if err != nil {
		return err
	}
	defer r.Close()

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}

	err = ioutil.WriteFile(newPath, buf, 0644)
	return err
}

// pwd
func (c FTPClient) CurrentDir() error {
	pwd, err := c.ftp.CurrentDir()
	if err != nil {
		return err
	}
	fmt.Println(pwd)
	return nil
}

// ls
func (c FTPClient) List(path string) error {
	list, err := c.ftp.NameList(path)
	if err != nil {
		return err
	}
	fmt.Println(list)
	return nil
}

// cd
func (c FTPClient) ChangeDir(path string) error {
	err := c.ftp.ChangeDir(path)
	return err
}

// mkdir
func (c FTPClient) MakeDir(path string) error {
	err := c.ftp.MakeDir(path)
	return err
}
