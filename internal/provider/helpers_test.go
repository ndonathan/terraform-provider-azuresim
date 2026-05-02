package provider

import (
	"encoding/base64"
	"net"
	"regexp"
	"strings"
	"testing"
)

// uuidV4Regex matches the canonical 8-4-4-4-12 lowercase hex layout with the
// version (4) and variant (8/9/a/b) bits enforced.
var uuidV4Regex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestSimulatedUUID_DeterministicAndValidV4(t *testing.T) {
	a := simulatedUUID("seed")
	b := simulatedUUID("seed")
	c := simulatedUUID("other")

	if a != b {
		t.Errorf("simulatedUUID is not deterministic: %q vs %q", a, b)
	}
	if a == c {
		t.Errorf("simulatedUUID collided across distinct seeds")
	}
	if !uuidV4Regex.MatchString(a) {
		t.Errorf("simulatedUUID returned %q which does not match v4 layout", a)
	}
}

func TestSimulatedIPv4_InTestNetRange(t *testing.T) {
	ip := simulatedIPv4("pip", "rg")
	parsed := net.ParseIP(ip)
	if parsed == nil || parsed.To4() == nil {
		t.Fatalf("simulatedIPv4 returned non-IPv4: %q", ip)
	}
	if !strings.HasPrefix(ip, "203.0.113.") {
		t.Errorf("simulatedIPv4 returned %q outside TEST-NET-3 (203.0.113.0/24)", ip)
	}
	if simulatedIPv4("pip", "rg") != ip {
		t.Errorf("simulatedIPv4 not deterministic")
	}
	if simulatedIPv4("pip2", "rg") == ip {
		// Hash collisions are possible but very unlikely for distinct inputs.
		t.Errorf("simulatedIPv4 unexpectedly collided")
	}
}

func TestSimulatedIPv6_InDocumentationRange(t *testing.T) {
	ip := simulatedIPv6("pip", "rg")
	parsed := net.ParseIP(ip)
	if parsed == nil || parsed.To4() != nil {
		t.Fatalf("simulatedIPv6 returned non-IPv6: %q", ip)
	}
	if !strings.HasPrefix(ip, "2001:db8::") {
		t.Errorf("simulatedIPv6 returned %q outside 2001:db8::/32 docs range", ip)
	}
}

func TestSimulatedMAC_HyperVOUI(t *testing.T) {
	mac := simulatedMAC("nic", "rg")
	if !strings.HasPrefix(mac, "00-15-5D-") {
		t.Errorf("simulatedMAC %q is missing Hyper-V OUI prefix", mac)
	}
	matched, err := regexp.MatchString(`^[0-9A-F]{2}(-[0-9A-F]{2}){5}$`, mac)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Errorf("simulatedMAC %q is not a valid 6-octet MAC", mac)
	}
}

func TestSimulatedPrivateIPv4_AvoidsReservedHostBits(t *testing.T) {
	for i := 0; i < 20; i++ {
		ip := simulatedPrivateIPv4("nic", "ipconf", i)
		if !strings.HasPrefix(ip, "10.0.0.") {
			t.Errorf("simulatedPrivateIPv4(%d) returned %q outside expected /24", i, ip)
		}
		// Azure reserves the first 4 host addresses (.0-.3); we should never hit them.
		// Parse the last octet via net.ParseIP to be safe.
		parsed := net.ParseIP(ip).To4()
		if parsed == nil {
			t.Fatalf("simulatedPrivateIPv4 returned %q which won't parse as IPv4", ip)
		}
		last := int(parsed[3])
		if last < 4 || last > 254 {
			t.Errorf("simulatedPrivateIPv4(%d) octet %d is in reserved/broadcast range", i, last)
		}
	}
}

func TestSimulatedSharedKey_Base64AndDeterministic(t *testing.T) {
	a := simulatedSharedKey("law/rg/name")
	b := simulatedSharedKey("law/rg/name")
	if a != b {
		t.Errorf("simulatedSharedKey not deterministic")
	}
	decoded, err := base64.StdEncoding.DecodeString(a)
	if err != nil {
		t.Fatalf("simulatedSharedKey is not valid base64: %v", err)
	}
	if len(decoded) != 64 {
		t.Errorf("simulatedSharedKey decoded length = %d, want 64", len(decoded))
	}
}

func TestSimulatedKVVersion_LengthAndHex(t *testing.T) {
	v := simulatedKVVersion("kv/secret/name")
	if len(v) != 32 {
		t.Errorf("simulatedKVVersion length = %d, want 32", len(v))
	}
	matched, err := regexp.MatchString(`^[0-9a-f]{32}$`, v)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Errorf("simulatedKVVersion = %q, want 32 lowercase hex characters", v)
	}
}

func TestVaultNameFromID(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "valid full ARM ID",
			in:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.KeyVault/vaults/myvault",
			want: "myvault",
		},
		{
			name: "trailing path segment",
			in:   "/subscriptions/x/resourceGroups/rg/providers/Microsoft.KeyVault/vaults/myvault/secrets/foo",
			want: "myvault",
		},
		{
			name: "missing marker",
			in:   "/subscriptions/x/resourceGroups/rg",
			want: "",
		},
		{
			name: "empty input",
			in:   "",
			want: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := vaultNameFromID(tc.in)
			if got != tc.want {
				t.Errorf("vaultNameFromID(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestRedisAndServiceBusKeys_Deterministic(t *testing.T) {
	if mkRedisKey("seed") != mkRedisKey("seed") {
		t.Errorf("mkRedisKey not deterministic")
	}
	if mkRedisKey("seed") == mkRedisKey("other") {
		t.Errorf("mkRedisKey collided")
	}
	if _, err := base64.StdEncoding.DecodeString(mkRedisKey("seed")); err != nil {
		t.Errorf("mkRedisKey is not valid base64: %v", err)
	}
	if mkSBKey("seed") != mkSBKey("seed") {
		t.Errorf("mkSBKey not deterministic")
	}
	if _, err := base64.StdEncoding.DecodeString(mkSBKey("seed")); err != nil {
		t.Errorf("mkSBKey is not valid base64: %v", err)
	}
}

func TestSBChildID(t *testing.T) {
	got := sbChildID("/subscriptions/x/.../namespaces/ns/", "queues", "myqueue")
	want := "/subscriptions/x/.../namespaces/ns/queues/myqueue"
	if got != want {
		t.Errorf("sbChildID = %q, want %q", got, want)
	}
}
