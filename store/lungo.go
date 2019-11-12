package store

import (
	"github.com/256dpi/lungo"
	"github.com/apex/log"
	"go.mongodb.org/mongo-driver/bson"
)

type lungoBE struct {
	clnt lungo.IClient
	eng  *lungo.Engine
	db   lungo.IDatabase
	col  lungo.ICollection
}

//Doc is what we will store the k/v as.
type Doc struct {
	Key   string
	Value []byte
}

//NewLungo gets you a new fake mongodb
func NewLungo(dbname, colln string) Store {
	var err error
	be := &lungoBE{}

	// open database
	be.clnt, be.eng, err = lungo.Open(nil,
		lungo.Options{
			Store: lungo.NewMemoryStore(),
		})
	if err != nil {
		panic(err)
	}

	be.db = be.clnt.Database(dbname)
	be.col = be.db.Collection(colln)

	return be
}

func (r *lungoBE) Get(key string) (results []byte, err error) {
	var item Doc
	log.Debugf("Getting key: %s", key)
	found := r.col.FindOne(nil, bson.M{"key": key})

	if err = found.Err(); err != nil {
		log.Debugf("FindOne error: %s", err)
		return results, err
	}
	err = found.Decode(&item)

	return item.Value, nil
}

func (r *lungoBE) Put(key string, val []byte) error {
	item := &Doc{
		Key:   key,
		Value: val,
	}
	_, err := r.col.InsertOne(nil, item)
	return err
}

func (r *lungoBE) Delete(key string) error {
	return nil
}

func (r *lungoBE) Keys() (results []string, err error) {
	var items []Doc

	found, err := r.col.Find(nil, bson.M{})
	if err != nil {
		return results, err
	}

	err = found.All(nil, &items)
	if err != nil {
		return results, err
	}
	for _, item := range items {
		log.Debugf("item: %+#v", item)
		results = append(results, item.Key)
	}
	return results, nil
}
