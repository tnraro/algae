package alga

import (
	"os"
	"slices"
	"testing"
	"tnraro/algae/internal/util"
)

func TestGetEmptyAlga(t *testing.T) {
	util.SetupDataDir(t)
	err0 := os.MkdirAll(util.DataDir("algae"), 0755)
	if err0 != nil {
		t.Fatal(err0)
	}
	t.Parallel()
	t.Run("empty algae", func(t *testing.T) {
		algae, err1 := GetAlgae()
		if err1 != nil {
			t.Fatal(err1)
		}
		util.AssertEq(t, len(algae), 0)
	})
	t.Run("empty alga", func(t *testing.T) {
		alga, err1 := GetAlga("hello_world")
		if alga != nil {
			t.Fatal(err1)
		}

		util.AssertEq(t, err1.msg, `The alga "hello_world" not exists`)
	})
}

func TestAlga(t *testing.T) {
	util.SetupDataDir(t)
	name := "hello_world"
	compose := "services:\n  app:\n    image: hello-world:latest"
	env := ""
	t.Run("create alga", func(t *testing.T) {
		if _, err := CreateAlga(name, compose, env); err != nil {
			t.Fatal(err)
		}
		algae, err := GetAlgae()
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Contains(algae, name) {
			t.Fatalf("no hello_world in %v", algae)
		}
	})

	t.Run("get alga", func(t *testing.T) {
		if alga, err := GetAlga(name); err != nil {
			t.Fatal(err)
		} else {
			util.AssertEq(t, alga.Name, name)
			util.AssertEq(t, alga.Compose, compose)
			util.AssertEq(t, alga.Env, env)
		}
	})

	t.Run("update alga", func(t *testing.T) {
		compose := "services:\n  app:\n    image: hello-world:latest\n    environment:\n    - TZ=${TZ:-Asia/Seoul}"
		env := "TZ=Asia/Seoul"
		if _, err := UpdateAlga(name, compose, env); err != nil {
			t.Fatal(err)
		}

		if alga, err := GetAlga(name); err != nil {
			t.Fatal(err)
		} else {
			util.AssertEq(t, alga.Name, name)
			util.AssertEq(t, alga.Compose, compose)
			util.AssertEq(t, alga.Env, env)
		}
	})

	t.Run("update alga config", func(t *testing.T) {
		t.Run("compose", func(t *testing.T) {
			compose := "services:\n  app:\n    image: hello-world:latest\n    environment:\n    - FOO=${FOO:-foo}"
			if _, err := UpdateAlgaConfig(name, "compose.yml", compose); err != nil {
				t.Fatal(err)
			}

			if alga, err := GetAlga(name); err != nil {
				t.Fatal(err)
			} else {
				util.AssertEq(t, alga.Compose, compose)
			}
		})
		t.Run("env", func(t *testing.T) {
			env := "FOO=bar"
			if _, err := UpdateAlgaConfig(name, ".env", env); err != nil {
				t.Fatal(err)
			}

			if alga, err := GetAlga(name); err != nil {
				t.Fatal(err)
			} else {
				util.AssertEq(t, alga.Env, env)
			}
		})
	})

	t.Run("delete alga", func(t *testing.T) {
		if _, err := DeleteAlga(name); err != nil {
			t.Fatal(err)
		}
		algae, err := GetAlgae()
		if err != nil {
			t.Fatal(err)
		}
		if slices.Contains(algae, name) {
			t.Fatalf("%s to be deleted", name)
		}
	})
}
