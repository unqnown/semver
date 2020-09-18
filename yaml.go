package semver

func (v *Version) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var decoded string
	if err := unmarshal(&decoded); err != nil {
		return err
	}
	*v, err = Parse(decoded)
	return err
}

func (v Version) MarshalYAML() (interface{}, error) { return v.String(), nil }
