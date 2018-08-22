package tfx_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/izumin5210/tfx"
)

func TestFactory_Load_Object(t *testing.T) {
	cases := []struct {
		test string
		opts []tfx.LoadOption
		want interface{}
	}{
		{
			test: "simple case",
			want: &User{
				ID:         1,
				Name:       "user-1",
				Preference: &Preference{Searchable: true, Reviewable: false},
			},
		},
		{
			test: "with params",
			opts: []tfx.LoadOption{tfx.WithParams(tfx.Params{"registered": true})},
			want: &User{
				ID:         1,
				Name:       "user-1",
				Registered: true,
				Preference: &Preference{Searchable: true, Reviewable: false},
			},
		},
		{
			test: "with loop",
			opts: []tfx.LoadOption{
				tfx.WithParams(tfx.Params{"registered": true}),
				tfx.WithLoop("postCount", 3),
			},
			want: &User{
				ID:         1,
				Name:       "user-1",
				Registered: true,
				Posts: []*Post{
					{ID: 1, UserID: 1, Title: "This is a post 1 by user-1", Tags: []string{"foo", "bar"}},
					{ID: 2, UserID: 2, Title: "This is a post 2 by user-1", Tags: []string{"foo", "bar"}},
					{ID: 3, UserID: 3, Title: "This is a post 3 by user-1", Tags: []string{"foo", "bar"}},
				},
				Preference: &Preference{Searchable: true, Reviewable: false},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			f := tfx.New()
			tb := &FakeTB{TB: t}

			var dest *User

			f.Load(tb, "object", dest, tc.opts...)

			if got, want := dest, tc.want; !reflect.DeepEqual(got, want) {
				t.Errorf("Loaded object is %#v, want %#v", got, want)
			}
		})
	}
}

//  Structs
//================================================================

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
}

type Preference struct {
	UserID     int  `json:"user_id"`
	Searchable bool `json:"searchable"`
	Reviewable bool `json:"reviewable"`
}

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Registered bool   `json:"registered"`

	Posts      []*Post     `json:"post"`
	Preference *Preference `json:"preference"`
}

//  Fake implementations
//================================================================

type FakeTB struct {
	testing.TB
}

func (FakeTB) Error(args ...interface{})                 {}
func (FakeTB) Errorf(format string, args ...interface{}) {}
func (FakeTB) Fail()                                     {}
func (FakeTB) FailNow()                                  {}
func (FakeTB) Failed() bool                              { return false }
func (FakeTB) Fatal(args ...interface{})                 {}
func (FakeTB) Fatalf(format string, args ...interface{}) {}
func (FakeTB) Log(args ...interface{})                   {}
func (FakeTB) Logf(format string, args ...interface{})   {}
func (FakeTB) Name() string                              { return "" }
func (FakeTB) Skip(args ...interface{})                  {}
func (FakeTB) SkipNow()                                  {}
func (FakeTB) Skipf(format string, args ...interface{})  {}
func (FakeTB) Skipped() bool                             { return false }
func (FakeTB) Helper()                                   {}
