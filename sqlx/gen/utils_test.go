package gen

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/artisanhe/tools/sqlx"
)

func init() {
	os.Chdir("./test")
}

func TestParseIndexesFromDoc(t *testing.T) {
	tt := assert.New(t)

	tt.Equal(&Keys{
		Primary: []string{"ID"},
	}, parseKeysFromDoc(`
	@def primary ID
	`))

	tt.Equal(&Keys{
		Indexes: sqlx.Indexes{
			"I_name":     []string{"Name"},
			"I_nickname": []string{"Nickname", "Name"},
		},
	}, parseKeysFromDoc(`
	@def index I_name   Name
	@def index I_nickname   Nickname Name
	`))

	tt.Equal(&Keys{
		Primary: []string{"ID"},
		Indexes: sqlx.Indexes{
			"I_nickname": []string{"Nickname", "Name"},
		},
		UniqueIndexes: sqlx.Indexes{
			"I_name": []string{"Name"},
		},
	}, parseKeysFromDoc(`
	@def primary ID
	@def index I_nickname Nickname Name
	@def unique_index I_name Name
	`))
}
