package flags_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/util/flags"
	"github.com/stretchr/testify/require"
)

func TestFlags(t *testing.T) {
	fs := flags.Flags{"--admin-username=foofoo", "--san foo", "--ucp-insecure-tls"}
	require.Equal(t, "--ucp-insecure-tls", fs[2])
	require.Equal(t, 0, fs.Index("--admin-username"))
	require.Equal(t, 1, fs.Index("--san"))
	require.Equal(t, 2, fs.Index("--ucp-insecure-tls"))
	require.True(t, fs.Include("--san"))

	fs.Delete("--san")
	require.Equal(t, 1, fs.Index("--ucp-insecure-tls"))
	require.False(t, fs.Include("--san"))

	fs.AddOrReplace("--san 10.0.0.1")
	require.Equal(t, 2, fs.Index("--san"))
	require.Equal(t, "--san 10.0.0.1", fs.Get("--san"))
	require.Equal(t, "10.0.0.1", fs.GetValue("--san"))
	require.Equal(t, "foofoo", fs.GetValue("--admin-username"))

	require.Len(t, fs, 3)
	fs.AddOrReplace("--admin-password=barbar")
	require.Equal(t, 3, fs.Index("--admin-password"))
	require.Equal(t, "barbar", fs.GetValue("--admin-password"))

	require.Len(t, fs, 4)
	fs.AddUnlessExist("--admin-password=borbor")
	require.Len(t, fs, 4)
	require.Equal(t, "barbar", fs.GetValue("--admin-password"))

	fs.AddUnlessExist("--help")
	require.Len(t, fs, 5)
	require.True(t, fs.Include("--help"))
}

func TestFlagsWithQuotes(t *testing.T) {
	fs := flags.Flags{"--admin-username \"foofoo\"", "--admin-password=\"foobar\""}
	require.Equal(t, "foofoo", fs.GetValue("--admin-username"))
	require.Equal(t, "foobar", fs.GetValue("--admin-password"))
}

func TestString(t *testing.T) {
	fs := flags.Flags{"--help", "--setting=false"}
	require.Equal(t, "--help --setting=false", fs.Join())
}

func TestGetBoolean(t *testing.T) {
	t.Run("Valid flags", func(t *testing.T) {
		testsValid := []struct {
			flag   string
			expect bool
		}{
			{"--flag", true},
			{"--flag=true", true},
			{"--flag=false", false},
			{"--flag=1", true},
			{"--flag=TRUE", true},
		}
		for _, test := range testsValid {
			fs := flags.Flags{test.flag}
			result, err := fs.GetBoolean(test.flag)
			require.NoError(t, err)
			require.Equal(t, test.expect, result)

			fs = flags.Flags{"--unrelated-flag1", "--unrelated-flag2=foo", test.flag}
			result, err = fs.GetBoolean(test.flag)
			require.NoError(t, err)
			require.Equal(t, test.expect, result)
		}
	})

	t.Run("Invalid flags", func(t *testing.T) {
		testsInvalid := []string{
			"--flag=foo",
			"--flag=2",
			"--flag=TrUe",
			"--flag=-4",
			"--flag=FalSe",
		}
		for _, test := range testsInvalid {
			fs := flags.Flags{test}
			_, err := fs.GetBoolean(test)
			require.Error(t, err)

			fs = flags.Flags{"--unrelated-flag1", "--unrelated-flag2=foo", test}
			_, err = fs.GetBoolean(test)
			require.Error(t, err)
		}
	})

	t.Run("Unknown flags", func(t *testing.T) {
		fs := flags.Flags{"--flag1=1", "--flag2"}
		result, err := fs.GetBoolean("--flag3")
		require.NoError(t, err)
		require.Equal(t, result, false)

	})
}
