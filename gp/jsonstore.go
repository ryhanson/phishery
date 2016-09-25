package gp

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
)

type JsonStore struct {
	filename string;
}

func (store *JsonStore) AddObject(obj interface{}) (bool, error) {
	objs := make(map[string]interface{})
	authJson, _ := json.Marshal(obj)
	authSum := getMD5(string(authJson))

	objs, err := store.loadStore()
	if err != nil {
		return false, err
	}

	if _, ok := objs[authSum]; ok {
		return false, nil
	}

	objs[authSum] = obj
	out, err := json.MarshalIndent(objs, "", "\t")
	if err != nil {
		return false, err
	}

	return true, ioutil.WriteFile(store.filename, out, 0755)
}

func (store *JsonStore) loadStore() (map[string]interface{}, error) {
	objs := make(map[string]interface{})

	file, err := os.Open(store.filename)
	if err != nil {
		return objs, err
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&objs); err != nil && err.Error() != "EOF" {
		return objs, err
	}

	return objs, nil
}

func getMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}