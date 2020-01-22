package semver

import "encoding/json"

func (v Version) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *Version) UnmarshalJSON(data []byte) (err error) {
	var decoded string
	if err = json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	*v, err = Parse(decoded)
	return err
}
