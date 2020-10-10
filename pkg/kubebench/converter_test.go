package kubebench

import (
	"encoding/json"
	"errors"
	starboard2 "github.com/aquasecurity/starboard/pkg/starboard"
	"os"
	"testing"

	starboard "github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConverter_Convert(t *testing.T) {
	config := starboard2.ConfigData{
		"kube-bench.imageRef": "aquasec/kube-bench:0.3.1",
	}
	var testcases = []struct {
		name string
		in   string // input File
		op   string // golden file
		err  error  // expected error
	}{
		{
			name: "Valid single json object in array",
			in:   "testdata/valid.json",
			op:   "testdata/goldenSingle.json",
			err:  nil,
		},
		{
			name: "invalid json object",
			in:   "testdata/invalid.json",
			err:  errors.New("json: cannot unmarshal object into Go value of type []v1alpha1.CISKubeBenchSection"),
		},
		{
			name: "Valid multiple json object in array",
			in:   "testdata/multiObjects.json",
			op:   "testdata/goldenMultiple.json",
			err:  nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			inFile, err := os.Open(tc.in)
			require.NoError(t, err)
			defer func() {
				_ = inFile.Close()
			}()

			var r starboard.CISKubeBenchOutput
			r, err = DefaultConverter.Convert(config, inFile)

			switch {
			case tc.err == nil:
				require.NoError(t, err)
				gFile, err := os.Open(tc.op)
				require.NoError(t, err)
				dec := json.NewDecoder(gFile)
				var kbop starboard.CISKubeBenchOutput
				err = dec.Decode(&kbop)
				require.NoError(t, err)
				defer func() {
					_ = gFile.Close()
				}()

				assert.Equal(t, kbop, r, "Converted report does not match expected report")
			default:
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}
