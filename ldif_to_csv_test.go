package ldiftocsv

import (
	"testing"
	"bytes"
	"strings"
)

func TestReadLdifFile(t *testing.T) {
	t.Run("parses single record", func(t *testing.T) {
		ldif := `# extended LDIF
#
# LDAPv3
# base <ou=groups,dc=example,dc=com> with scope subtree
# filter: (objectclass=*)
# requesting: cn owner
#

# test-group, groups, example.com
dn: cn=test-group,ou=groups,dc=example,dc=com
cn: test-group
owner: uid=bob,ou=users,dc=example,dc=com
`

		headers := []string{
			"cn", "owner",
		}

		ldifEntries := ReadLdifFile(bytes.NewBufferString(ldif), headers)

		if len(ldifEntries) != 1 {
			t.Errorf("got %d entries, wanted %d", len(ldifEntries), 1)
		}
	})

	t.Run("parses multiple records", func(t *testing.T) {
		ldif := `# extended LDIF
#
# LDAPv3
# base <ou=groups,dc=example,dc=com> with scope subtree
# filter: (objectclass=*)
# requesting: cn owner
#

# test-group, groups, example.com
dn: cn=test-group,ou=groups,dc=example,dc=com
cn: test-group
owner: uid=bob,ou=users,dc=example,dc=com

# test-group-2, groups, example.com
dn: cn=test-group-2,ou=groups,dc=example,dc=com
cn: test-group-2
owner: uid=bob,ou=users,dc=example,dc=com
`

		headers := []string{
			"cn", "owner",
		}

		ldifEntries := ReadLdifFile(bytes.NewBufferString(ldif), headers)

		if len(ldifEntries) != 2 {
			t.Errorf("got %d entries, wanted %d", len(ldifEntries), 1)
		}
	})

	t.Run("handles multiple values in single field", func(t *testing.T) {
		ldif := `# extended LDIF
#
# LDAPv3
# base <ou=groups,dc=example,dc=com> with scope subtree
# filter: (objectclass=*)
# requesting: cn owner
#

# test-group, groups, example.com
dn: cn=test-group,ou=groups,dc=example,dc=com
cn: test-group
owner: uid=bob,ou=users,dc=example,dc=com
owner: uid=frank,ou=users,dc=example,dc=com
`

		headers := []string{
			"cn", "owner",
		}

		ldifEntries := ReadLdifFile(bytes.NewBufferString(ldif), headers)

		if len(ldifEntries) != 1 {
			t.Fatalf("got %d entries, wanted %d", len(ldifEntries), 1)
		}

		owners := []string{}
		for _, prop := range ldifEntries[0].properties {
			if prop.name == "owner" {
				owners = strings.Split(prop.value, "\n")
			}
		}
		if len(owners) != 2 {
			t.Errorf("found %d owners, wanted %d", len(owners), 2)
		}
	})
}

func TestWriteToCsv(t *testing.T) {
	t.Run("writes single record", func(t *testing.T) {
		out := bytes.Buffer{}
		entries := []LdifEntry{
			LdifEntry{
				properties: []NameValue{
					NameValue{name: "cn", value: "test-group"},
					NameValue{name: "owner", value: "bob"},
				},
			},
		}
		expected := "test-group,bob\n"

		WriteToCsv(&out, entries)

		if out.String() != expected {
			t.Errorf("wanted '%s', got '%s'", expected, out.String())
		}
	})

	t.Run("writes multilpe record", func(t *testing.T) {
		out := bytes.Buffer{}
		entries := []LdifEntry{
			LdifEntry{
				properties: []NameValue{
					NameValue{name: "cn", value: "test-group"},
					NameValue{name: "owner", value: "bob"},
				},
			},
			LdifEntry{
				properties: []NameValue{
					NameValue{name: "cn", value: "test-group-2"},
					NameValue{name: "owner", value: "bob"},
				},
			},
		}
		expected := "test-group,bob\ntest-group-2,bob\n"

		WriteToCsv(&out, entries)

		if out.String() != expected {
			t.Errorf("wanted '%s', got '%s'", expected, out.String())
		}
	})
}
