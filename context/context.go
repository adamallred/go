package context

import (
	gocontext "context"
	"encoding/binary"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/firestore"
	//"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Route is the value part of a shortcut.
type Route struct {
	URL  string    `json:"url" firestore:"url"`
	Time time.Time `json:"time" firestore:"time"`
}

// NextID is the next numeric ID to use for auto-generated IDs
type NextID struct {
	ID uint32 `json:"id" firestore:"id"`
}

// Serialize this Route into the given Writer.
func (o *Route) write(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, o.Time.UnixNano()); err != nil {
		return err
	}

	if _, err := w.Write([]byte(o.URL)); err != nil {
		return err
	}

	return nil
}

// Deserialize this Route from the given Reader.
func (o *Route) read(r io.Reader) error {
	var t int64
	if err := binary.Read(r, binary.LittleEndian, &t); err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	o.URL = string(b)
	o.Time = time.Unix(0, t)
	return nil
}

// Context provides access to the data store.
type Context struct {
	db *firestore.Client
}

// Open the context. Instantiate a new firestore client
func Open() (*Context, error) {
	ctx := gocontext.Background()
	client, err := firestore.NewClient(ctx, getGoogleProject())
	if err != nil {
		return nil, err
	}

	return &Context{
		db: client,
	}, nil
}

// Close the resources associated with this context.
func (c *Context) Close() error {
	return c.db.Close()
}

// Get retreives a shortcut from the data store.
func (c *Context) Get(name string) (*Route, error) {
	ref := c.db.Doc("routes/" + name)

	ctx, cancel := gocontext.WithTimeout(gocontext.Background(), time.Minute)
	defer cancel()

	snap, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}

	var rt Route
	if err := snap.DataTo(&rt); err != nil {
		return nil, err
	}

	return &rt, nil
}

// Put stores a new shortcut in the data store.
func (c *Context) Put(key string, rt *Route) error {
	ref := c.db.Doc("routes/" + key)

	ctx, cancel := gocontext.WithTimeout(gocontext.Background(), time.Minute)
	defer cancel()

	_, err := ref.Create(ctx, rt)
	if err != nil {
		return err
	}

	return nil
}

// Del removes an existing shortcut from the data store.
func (c *Context) Del(key string) error {
	ref := c.db.Doc("routes/" + key)

	ctx, cancel := gocontext.WithTimeout(gocontext.Background(), time.Minute)
	defer cancel()

	_, err := ref.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

// List all routes in an iterator, starting with the key prefix of start (which can also be nil).
// func (c *Context) List(start []byte) *Iter {
// 	return &Iter{
// 		it: c.db.NewIterator(&util.Range{
// 			Start: start,
// 			Limit: nil,
// 		}, nil),
// 	}
// }

// GetAll gets everything in the db to dump it out for backup purposes
func (c *Context) GetAll() (map[string]Route, error) {
	golinks := map[string]Route{}
	col := c.db.Collection("routes")

	ctx, cancel := gocontext.WithTimeout(gocontext.Background(), time.Minute)
	defer cancel()

	routes, err := col.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range routes {
		var rt Route
		if err := doc.DataTo(&rt); err != nil {
			return nil, err
		}
		golinks[doc.Ref.ID] = rt
	}
	return golinks, nil
}

// NextID generates the next numeric ID to be used for an auto-named shortcut.
func (c *Context) NextID() (uint32, error) {
	ref := c.db.Doc("IDs/nextID")
	var nid uint32

	ctx, cancel := gocontext.WithTimeout(gocontext.Background(), time.Minute)
	defer cancel()

	err := c.db.RunTransaction(ctx, func(ctx gocontext.Context, tx *firestore.Transaction) error {
		var nextID *NextID

		doc, err := tx.Get(ref)
		if err != nil {
			if grpc.Code(err) == codes.NotFound {
				// this is the very first auto-generated ID, we can make it
				// as :1 and return it
				nextID = new(NextID)
				nextID.ID = 1
				nid = 1
				err := tx.Create(ref, nextID)
				if err != nil {
					return err
				}
				return nil
			}
			return err
		}

		if err := doc.DataTo(&nextID); err != nil {
			return err
		}
		nextID.ID += 1
		nid = nextID.ID

		return tx.Set(ref, &nextID)
	})
	if err != nil {
		return 0, err
	}

	return nid, nil
}

func getGoogleProject() string {
	// TODO: take this in as an flag
	return "aallred-sawa-poc"
	// ctx, cancel := gocontext.WithTimeout(gocontext.Background(), time.Second*5)
	// defer cancel()

	// creds, err := google.FindDefaultCredentials(ctx)
	// if err != nil {
	// 	return ""
	// }
	// return creds.ProjectID
}
