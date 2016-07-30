package redis

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/redis.v3"

	"github.com/containerops/arkor/utils/db/factory"
)

var (
	Client *redis.Client
)

type redisdrv struct{}

func init() {
	factory.RegisterKV("redis", &redisdrv{})
}

func (r *redisdrv) RegisterModel(models ...interface{}) error {
	return nil
}

func (r *redisdrv) InitDB(driver, user, passwd, uri, name string, partition int64) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: passwd,
		DB:       partition,
	})

	if _, err := Client.Ping().Result(); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *redisdrv) Save(obj interface{}) error {
	// Generate a key and save the key to the key-list,
	// the key of the key-list is the object type.
	key := getKey(obj)
	objectType := reflect.TypeOf(obj).Elem().Name()
	Client.SAdd(strings.ToLower(objectType), key)
	// Save object to a redis key
	result, err := json.Marshal(&obj)
	if err != nil {
		return err
	}
	if _, err := Client.Set(key, string(result), 0).Result(); err != nil {
		return err
	}

	return nil
}

func (r *redisdrv) Create(obj interface{}) error {
	return r.Save(obj)
}

func (r *redisdrv) Delete(obj interface{}) error {
	// Get key of the object
	key := getKey(obj)
	// Delete object hash
	if Client.Exists(key).Val() {
		Client.Del(key)
	} else {
		return fmt.Errorf("object not exist")
	}
	// Delete key in key-list
	objectType := reflect.TypeOf(obj).Elem().Name()
	Client.SRem(strings.ToLower(objectType), key)

	return nil
}

func (r *redisdrv) Query(obj interface{}) (bool, error) {
	key := getKey(obj)
	result, err := Client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		} else {
			return false, err
		}
	}

	if err = json.Unmarshal([]byte(result), &obj); err != nil {
		return false, err
	}

	return true, nil
}

func (r *redisdrv) QueryMulti(condition interface{}, value interface{}) (bool, error) {
	return true, nil
}

func getKey(obj interface{}) (result string) {
	objectType := reflect.TypeOf(obj).Elem().Name()
	s := reflect.ValueOf(obj).Elem()
	typeOfS := s.Type()
	keys := []string{}

	for k := 0; k < s.NumField(); k++ {
		t := typeOfS.Field(k).Name
		if t == "ID" {
			keys = append(keys, s.Field(k).Interface().(string))
		}
	}
	result = fmt.Sprintf("%s-%s", strings.ToLower(objectType), keys[0])
	return result
}

/* Currently the redis driver does not provide a db reference */
func (r *redisdrv) GetDB() interface{} {
	return nil
}
