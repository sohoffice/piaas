package piaas

import "testing"

func TestConfig_GetApp(t *testing.T) {
	conf := Config{
		Apps: []App{
			{
				Name: "foo",
			},
			{
				Name: "bar",
			},
		},
	}

	_, err := conf.GetApp("baz")
	if err == nil {
		t.Errorf("Should return error for baz can not be found.")
	}

	_, err = conf.GetApp("foo")
	if err != nil {
		t.Errorf("Should not return error for foo can be found.")
	}

	_, err = conf.GetApp("")
	if err == nil {
		t.Errorf("Should return error for having multiple app")
	}

	conf2 := Config{
		Apps: []App{
			{
				Name: "foo",
			},
		},
	}
	_, err = conf2.GetApp("")
	if err != nil {
		t.Errorf("Should find the sole app")
	}
}
